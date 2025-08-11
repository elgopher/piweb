// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

var audioCtx;
var audioWorkletNode;

// block main thread until audio worklet node is created
(function startAudio() {
    (async () => {
        audioCtx = new AudioContext();
        await audioCtx.audioWorklet.addModule("audio-processor.js");
        audioWorkletNode = new AudioWorkletNode(audioCtx, "pi-audio-processor")
        audioWorkletNode.connect(audioCtx.destination);
    })();
})();

function onGesture() {
    window.removeEventListener("pointerdown", onGesture);
    window.removeEventListener("keydown", onGesture);
    if (audioCtx.state === "suspended") {
        console.log("Audio resumed");
        audioCtx.resume();
    }
}

window.addEventListener("pointerdown", onGesture);
window.addEventListener("keydown", onGesture);

//

var audioCommandsSharedArrayBuffer = new SharedArrayBuffer(1024);
var audioCommands = new Uint8Array(audioCommandsSharedArrayBuffer);

const cmdKindSetSample = 1
const cmdKindSetLoop = 2
const cmdKindClearChan = 3
const cmdKindSetPitch = 4
const cmdKindSetVolume  = 5

