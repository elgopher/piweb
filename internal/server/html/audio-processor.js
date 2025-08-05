// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

class PiAudioProcessor extends AudioWorkletProcessor {
    constructor() {
        super();
        console.log("PiAudioProcessor constructed");
    }

    process(inputs, outputs, parameters) {
        const output = outputs[0];
        for (let channel = 0; channel < output.length; ++channel) {
            const outputChannel = output[channel];
            for (let i = 0; i < outputChannel.length; ++i) {
                outputChannel[i] = Math.sin(2 * Math.PI * 440 * currentFrame / sampleRate);
            }
        }
        return true; // continue processing
    }
}

registerProcessor("pi-audio-processor", PiAudioProcessor);