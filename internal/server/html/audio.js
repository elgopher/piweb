// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

/** @type {AudioContext} */
var audioCtx;

/** @type {AudioWorkletNode} */
var audioWorkletNode;

const ctxSampleRate = 48000;

// TODO DUPLICATED CODE

class SharedRing {
    static HEADER_BYTES = 12; // Int32[3]: [capacity, head, tail]

    constructor(sab) {
        if (!(sab instanceof SharedArrayBuffer)) throw new TypeError("SharedArrayBuffer is required");
        this.ctrl = new Int32Array(sab, 0, 3);           // [cap, head, tail]
        this.cap  = sab.byteLength - SharedRing.HEADER_BYTES;
        if (this.cap <= 0) throw new RangeError("Za mały SAB");

        // Inicjalizacja nagłówka (idempotentna)
        if (Atomics.load(this.ctrl, 0) !== this.cap) {
            Atomics.store(this.ctrl, 0, this.cap); // capacity
            Atomics.store(this.ctrl, 1, 0);        // head
            Atomics.store(this.ctrl, 2, 0);        // tail
        }
        this.data = new Uint8Array(sab, SharedRing.HEADER_BYTES, this.cap);
    }

    // ile miejsca zostało w buforze
    freeSpace() {
        const H = Atomics.load(this.ctrl, 1) >>> 0;
        const T = Atomics.load(this.ctrl, 2) >>> 0;
        const used = (H - T) >>> 0;
        return this.cap - used;
    }

    // Zapisuje jak najwięcej bajtów z src; zwraca liczbę zapisanych.
    write(src) {
        const a = toU8(src);
        if (a.length === 0) return 0;

        const H = Atomics.load(this.ctrl, 1) >>> 0;
        const T = Atomics.load(this.ctrl, 2) >>> 0;
        const used = (H - T) >>> 0;
        const free = this.cap - used;
        if (free === 0) return 0;

        const n = Math.min(free, a.length);
        const w = H % this.cap;

        const first = Math.min(n, this.cap - w);
        this.data.set(a.subarray(0, first), w);
        if (n > first) this.data.set(a.subarray(first, n), 0);

        Atomics.store(this.ctrl, 1, (H + n) >>> 0); // publikacja head
        return n;
    }

    // Czyta do dst; zwraca liczbę odczytanych bajtów.
    read(dst) {
        const b = toU8(dst);
        if (b.length === 0) return 0;

        const H = Atomics.load(this.ctrl, 1) >>> 0;
        const T = Atomics.load(this.ctrl, 2) >>> 0;
        const avail = (H - T) >>> 0;
        if (avail === 0) return 0;

        const n = Math.min(avail, b.length);
        const r = T % this.cap;

        const first = Math.min(n, this.cap - r);
        b.set(this.data.subarray(r, r + first), 0);
        if (n > first) b.set(this.data.subarray(0, n - first), first);

        Atomics.store(this.ctrl, 2, (T + n) >>> 0); // publikacja tail
        return n;
    }
}

function toU8(x) {
    if (x instanceof Uint8Array) return x;
    if (ArrayBuffer.isView(x)) return new Uint8Array(x.buffer, x.byteOffset, x.byteLength);
    if (x instanceof ArrayBuffer) return new Uint8Array(x);
    throw new TypeError("Oczekiwano Uint8Array/ArrayBufferView/ArrayBuffer");
}

// block main thread until audio worklet node is created
(function startAudio() {
    (async () => {
        audioCtx = new AudioContext({
            latencyHint: "interactive", sampleRate: ctxSampleRate
        });
        await audioCtx.audioWorklet.addModule("audio-processor.js");
        audioWorkletNode = new AudioWorkletNode(
            audioCtx,
            "pi-audio-processor",
            {
                numberOfInputs: 0,
                numberOfOutputs: 1,
                outputChannelCount: [2],
                channelInterpretation: 'speakers',
                processorOptions: {
                    commandsBuffer: audioCommandsSharedArrayBuffer
                }
            })
        audioWorkletNode.connect(audioCtx.destination);
    })();
})();

function waitUntilAudioIsRunning() {
    return new Promise(resolve => {
        if (audioCtx.state === "suspended") {
            const onGesture = function () {
                window.removeEventListener("pointerdown", onGesture);
                window.removeEventListener("keydown", onGesture);
                console.log("Audio resumed");
                audioCtx.resume();
                resolve();
            }

            window.addEventListener("pointerdown", onGesture);
            window.addEventListener("keydown", onGesture);
        } else {
            resolve();
        }
    });
}

var temporaryBuffer = new Uint8Array(32768);
var temporaryBufferLength = 0;
var audioCommandsSharedArrayBuffer = new SharedArrayBuffer(32768);
var audioCommands = new SharedRing(audioCommandsSharedArrayBuffer);

/**
 * @param {string} id
 * @param {Uint8Array} data
 * @param {number} rate
 */
function loadSample(id, data, rate) {
    const msg = {
        type: "loadSample",
        id: id,
        data: data,
        rate: rate
    }
    audioWorkletNode.port.postMessage(msg, [data.buffer]); // zero-copy
}

/** @param {string} id */
function unloadSample(id) {
    const msg = {
        type: "unloadSample",
        id: id
    }
    audioWorkletNode.port.postMessage(msg);
}

function sendCommands() {
    audioCommands.write(temporaryBuffer.subarray(0, temporaryBufferLength));
    console.log("stored " + temporaryBufferLength + " in the audioCommands buffer");
    temporaryBufferLength = 0;
}