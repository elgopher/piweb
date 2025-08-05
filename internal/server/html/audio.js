// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

import {SharedRing} from './shared-ring.js';

window.temporaryBuffer = new Uint8Array(32768);
window.temporaryBufferLength = 0;

const audioCommandsSharedArrayBuffer = new SharedArrayBuffer(32768);
window.audioCommands = new SharedRing(audioCommandsSharedArrayBuffer);

// block main thread until audio worklet node is created
(async function startAudio() {
    window.audioCtx = new AudioContext({
        latencyHint: "interactive", sampleRate: 48000,
    });
    await audioCtx.audioWorklet.addModule("audio-processor.js");
    window.audioWorkletNode = new AudioWorkletNode(
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
        });
    audioWorkletNode.connect(audioCtx.destination);
})();

export function waitUntilAudioIsRunning() {
    return new Promise(resolve => {
        if (audioCtx.state === "suspended") {
            const onGesture = function () {
                window.removeEventListener("pointerdown", onGesture);
                window.removeEventListener("keydown", onGesture);
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

/**
 * @param {number} id
 * @param {Uint8Array} data
 * @param {number} rate
 */
window.loadSample = function (id, data, rate) {
    const msg = {
        type: "loadSample",
        id: id,
        data: data,
        rate: rate
    }
    audioWorkletNode.port.postMessage(msg, [data.buffer]); // zero-copy
}

/** @param {number} id */
window.unloadSample = function (id) {
    const msg = {
        type: "unloadSample",
        id: id
    }
    audioWorkletNode.port.postMessage(msg);
}

window.sendCommands = function () {
    audioCommands.write(temporaryBuffer.subarray(0, temporaryBufferLength));
    temporaryBufferLength = 0;
}
