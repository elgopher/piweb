// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// @ts-check
/// <reference path="./types/audioworklet.d.ts" />

import {Pool} from './pool.js';
import {SharedRing} from './shared-ring.js';

const cmdSetSample = 1;
const cmdSetLoop = 2;
const cmdClearChan = 3;
const cmdSetPitch = 4;
const cmdSetVolume = 5;

const loopNone = 0;
const loopForward = 1;

/**
 * @callback predicate
 * @param {number} i
 * @return {boolean}
 */

/**
 * Works like sort.Search(n, pred) in Go.
 *
 * binarySearch finds and return the smallest index i
 * in [0, n) at which pred(i) is true, assuming that on the range [0, n),
 * pred(i)  true implies pred(i+1)  true.
 *
 * @param n {number}
 * @param pred {predicate}
 * @return {number}
 */
function binarySearch(n, pred) {
    let i = 0, j = n;
    while (i < j) {
        const h = (i + j) >>> 1;
        if (pred(h)) j = h;
        else i = h + 1;
    }
    return i;
}

class Command {
    command = cmdSetSample
    ch = 0
    id = 0
    offset = 0
    pitch = 0.0
    time = 0.0
    volume = 0.0
    loopStart = 0
    loopStop = 0
    loopType = loopNone
    uses = 0 // the number of channels the command is currently used
}

/** @type {Pool<Command>} */
const commandPool = new Pool(function () {
    return new Command()
})

const chanLen = 4;
const commandSize = 19;

class Sample {
    /** @type Int8Array */
    data
    rate = 0
}

class SampleOutput {
    /** @type number from -1.0 to 1.0 */
    value
    /** @type boolean If true sample was returned, if false then no data */
    ok
}

class Channel {
    /** @type {Int8Array|null} */
    sampleData = null;

    active = false;
    position = 0.0;
    pitch = 1.0;
    sampleRate = 0;
    volume = 1.0;
    loopStart = 0;
    loopStop = (2 ** 31) - 1; // max int32
    loopType = loopNone;

    /** @param out {SampleOutput} */
    nextSample(out) {
        if (!this.active || this.sampleData === null || this.volume <= 0) {
            out.value = 0;
            out.ok = false;
            return;
        }

        let pos = Math.floor(this.position);
        if (pos >= Math.min(this.sampleData.length, this.loopStop)) {
            // End of sample
            if (this.loopType === loopForward) {
                this.position = this.loopStart;
                pos = this.loopStart;
            } else {
                this.active = false;
                out.value = 0;
                out.ok = false;
                return;
            }
        }

        let sample = this.sampleData[pos]; // int8
        sample /= 128; // convert int8 to -1.0 - 1.0

        // Advance position
        this.position += (this.sampleRate / sampleRate) * this.pitch;

        // Apply volume
        sample *= this.volume;

        out.value = sample;
        out.ok = true;
    }
}

class PiAudioProcessor extends AudioWorkletProcessor {
    /** @type {Map<number, Sample>} */
    #samplesByID = new Map();
    /** @type SharedRing - contains commands sent by main thread encoded as bytes */
    #commands;

    /** @type Uint8Array - buffer used as an output for reading commands */
    #buffer;
    /** @type DataView - DataView for #this.buffer */
    #bufferView;

    /** @type {Array<Channel>} - Channel states */
    #channels;
    /** @type {Array<Array<Command>>} - each channel's planned commands sorted by time */
    #commandsByTime;

    /** reusable output for channel.nextSample (to avoid allocation on each call) */
    #sample = new SampleOutput();

    constructor(options) {
        super()

        this.#commands = new SharedRing(options.processorOptions.commandsBuffer);

        const buffer = new ArrayBuffer(commandSize * 128); // 128 commands per 2.5ms -> 50k cmds per sec
        this.#buffer = new Uint8Array(buffer);
        this.#bufferView = new DataView(buffer);

        this.#channels = [];
        for (let i = 0; i < chanLen; i++) {
            const ch = new Channel();
            this.#channels.push(ch);
        }

        this.#commandsByTime = [];
        for (let i = 0; i < chanLen; i++) {
            this.#commandsByTime.push([]);
        }

        this.port.onmessage = (event) => {
            const msg = event.data;
            switch (msg.type) {
                case "loadSample":
                    // convert uint8 to int8 (without copying):
                    const data =
                        new Int8Array(msg.data.buffer, msg.data.byteOffset, msg.data.length);

                    const sample = new Sample();
                    sample.data = data;
                    sample.rate = msg.rate;

                    this.#samplesByID[msg.id] = sample;
                    break;
                case "unloadSample":
                    delete this.#samplesByID[msg.id];
                    break;
            }
        };
    }

    process(inputs, outputs, parameters) {
        this.#processCommands();

        const output = outputs[0];
        const frames = output[0].length;

        this.#runCommands();

        for (let i = 0; i < frames; i++) {
            this.#write(i, output);
        }

        return true; // continue processing
    }

    #processCommands() {
        const bytesRead = this.#commands.read(this.#buffer);

        if (bytesRead > 0) {
            let off = 0;
            while (off < bytesRead) {
                const command = this.#bufferView.getUint8(off);

                /** @type Command */
                const cmd = commandPool.get();

                cmd.command = command;
                cmd.ch = this.#bufferView.getUint8(off + 1);

                // command will be used by all channels and each channel is different array
                // therefore add the "uses" field which will track how many arrays still reference
                // the command
                cmd.uses = 0;
                for (let i = 0; i < chanLen; i++) {
                    const chanNum = 1 << i;
                    if ((cmd.ch & chanNum) !== 0) {
                        cmd.uses++;
                    }
                }

                switch (command) {
                    case cmdSetSample:
                        cmd.id = this.#bufferView.getUint32(off + 2, true);
                        cmd.offset = this.#bufferView.getInt32(off + 6, true);
                        cmd.time = this.#bufferView.getFloat64(off + 10, true);
                        break;
                    case cmdSetLoop:
                        cmd.loopStart = this.#bufferView.getInt32(off + 2, true);
                        cmd.loopStop = this.#bufferView.getInt32(off + 6, true);
                        cmd.loopType = this.#bufferView.getUint8(off + 10);
                        cmd.time = this.#bufferView.getFloat64(off + 11, true);
                        break;
                    case cmdClearChan:
                        cmd.time = this.#bufferView.getFloat64(off + 2, true);
                        break;
                    case cmdSetPitch:
                        cmd.pitch = this.#bufferView.getFloat64(off + 2, true);
                        cmd.time = this.#bufferView.getFloat64(off + 10, true);
                        break;
                    case cmdSetVolume:
                        cmd.volume = this.#bufferView.getFloat64(off + 2, true);
                        cmd.time = this.#bufferView.getFloat64(off + 10, true);
                        break;
                }

                off += commandSize;

                if (cmd.command !== cmdClearChan) {
                    for (let i = 0; i < chanLen; i++) {
                        const chanNum = 1 << i
                        // a single command can be executed on multiple channels at once
                        if ((cmd.ch & chanNum) !== 0) {
                            this.#commandsByTime[i].push(cmd);
                        }
                    }
                } else {
                    this.#clearChan(cmd.ch, cmd.time);
                    commandPool.put(cmd);
                }
            }

            // TODO use manual merge sort: all new commands must be sorted first then merge two sorted arrays

            for (let i = 0; i < chanLen; i++) {
                this.#commandsByTime[i].sort((a, b) => a.time - b.time);
            }
        }
    }

    #clearChan(ch, time) {
        for (let i = 0; i < chanLen; i++) {
            const chanNum = 1 << i
            if ((ch & chanNum) !== 0) {
                const chanCommands = this.#commandsByTime[i];
                const newLength = binarySearch(
                    chanCommands.length,
                    (idx) => chanCommands[idx].time >= time
                );
                // put removed commands back to the object pool
                for (let j = newLength; j < chanCommands.length; j++) {
                    chanCommands[j].uses--;
                    if (chanCommands[j].uses === 0) {
                        commandPool.put(chanCommands[j]);
                    }
                }
                chanCommands.length = newLength;
            }
        }
    }

    #runCommands() {
        for (let i = 0; i < chanLen; i++) {
            const selectedChan = this.#channels[i];
            let processed = 0;

            for (const cmd of this.#commandsByTime[i]) {
                if (cmd.time > currentTime) {
                    break;
                }

                switch (cmd.command) {
                    case cmdSetSample:
                        switch (true) {
                            case cmd.id === 0:
                                selectedChan.active = false;
                                selectedChan.sampleData = null;
                                break;
                            case this.#samplesByID[cmd.id] === undefined:
                                console.warn("[piaudio] SetSample failed: Sample not found, id:", cmd.id)
                                selectedChan.active = false;
                                selectedChan.sampleData = null;
                                break;
                            default:
                                selectedChan.active = true;
                                const sample = /** @type Sample */ this.#samplesByID[cmd.id];
                                selectedChan.sampleData = sample.data;
                                selectedChan.sampleRate = sample.rate;
                                break;
                        }
                        selectedChan.position = cmd.offset;
                        break;
                    case cmdSetLoop:
                        selectedChan.loopStart = cmd.loopStart;
                        selectedChan.loopStop = cmd.loopStop;
                        selectedChan.loopType = cmd.loopType;
                        break;
                    case cmdSetPitch:
                        selectedChan.pitch = cmd.pitch;
                        break;
                    case cmdSetVolume:
                        selectedChan.volume = cmd.volume;
                        break;
                }

                cmd.uses--;
                if (cmd.uses === 0) {
                    commandPool.put(cmd); // this cmd will be removed from array after the loop
                }

                processed++;
            }

            if (processed > 0) {
                this.#commandsByTime[i].copyWithin(0, processed);
                this.#commandsByTime[i].length -= processed;
            }
        }
    }

    #write(sampleIndex, output) {
        let mixL = 0;
        let mixR = 0;

        for (let ch = 0; ch < chanLen; ch++) {
            this.#channels[ch].nextSample(this.#sample);
            if (!this.#sample.ok) {
                continue;
            }

            if (ch === 0 || ch === 3) {
                mixL += this.#sample.value;
            } else {
                mixR += this.#sample.value;
            }
        }

        output[0][sampleIndex] = mixL;
        output[1][sampleIndex] = mixR;
    }
}

registerProcessor("pi-audio-processor", PiAudioProcessor);