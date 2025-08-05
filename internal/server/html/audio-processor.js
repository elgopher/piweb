// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// @ts-check
/// <reference path="./types/audioworklet.d.ts" />

import {createPool} from './pool.js';
import {SharedRing} from './sharedring.js';

const cmdSetSample = 1;
const cmdSetLoop = 2;
const cmdClearChan = 3;
const cmdSetPitch = 4;
const cmdSetVolume = 5;

const loopNone = 0;
const loopForward = 1;

const commandPool = createPool(function () {
    return {
        command: cmdSetSample,
        ch: 0,
        sampleID: 0,
        offset: 0,
        pitch: 0.0,
        time: 0.0,
        vol: 0.0,
        loopStart: 0,
        loopStop: 0,
        loopType: loopNone,
    }
})

const sampleTime = 1.0 / sampleRate;
const chanLen = 4;
const commandSize = 19;

class Sample {
    /** @type number -1.0 do 1.0 */
    value
    /** @type boolean */
    ok
}

class channel {
    /** @type boolean */
    active
    /** @type {Uint8Array} */
    sampleData
    /** @type number */
    position
    /** @type number */
    pitch
    /** @type number */
    sampleRate
    /** @type number */
    volume
    /** @type number */
    loopStart
    /** @type number */
    loopStop
    /** @type number */
    loopType

    // out is of class sample
    nextSample(out) {
        if (!this.active || this.sampleData === null || this.volume <= 0) {
            out.value = 0;
            out.ok = false;
            return;
        }

        let pos = Math.ceil(this.position);
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

        let sample = this.sampleData[pos]; // uint8
        sample = (sample - 128) / 128

        // Advance position
        this.position += (this.sampleRate / sampleRate) * this.pitch;

        // Apply volume
        sample *= this.volume

        out.value = sample;
        out.ok = true;
    }
}

class PiAudioProcessor extends AudioWorkletProcessor {
    samplesByID = {};
    #commands;

    #buffer;
    #bufferView;


    #channels; // [chanLen]channel
    #commandsByTime; // [chanLen][]command - each channel's planned commands sorted by time

    #preciseTime;

    /** @type Sample */
    #sample;

    constructor(options) {
        super()

        this.#commands = new SharedRing(options.processorOptions.commandsBuffer);

        const buffer = new ArrayBuffer(commandSize * 128); // 128 commands per 2.5ms -> 50k cmds per sec
        this.#buffer = new Uint8Array(buffer);
        this.#bufferView = new DataView(buffer);

        this.#channels = [];
        for (let i = 0; i < chanLen; i++) {
            const ch = new channel();
            ch.active = false;
            ch.sampleData = null;
            ch.position = 0.0;
            ch.pitch = 1.0;
            ch.sampleRate = 0;
            ch.volume = 1.0;
            ch.loopStart = 0;
            ch.loopStop = (2 ** 31) - 1; // max int32
            ch.loopType = loopNone;

            this.#channels.push(ch);
        }

        this.#commandsByTime = [];
        for (let i = 0; i < chanLen; i++) {
            this.#commandsByTime.push([]);
        }

        this.#sample = new Sample();

        this.port.onmessage = (event) => {
            const msg = event.data;
            switch (msg.type) {
                case "loadSample":
                    console.log(msg.data);
                    this.samplesByID[msg.id] = msg;
                    console.log("sample %d loaded: %s bytes, rate %d", msg.id, msg.data.length, msg.rate);
                    break;
                case "unloadSample":
                    delete this.samplesByID[msg.id];
                    console.log("sample %d unloaded", msg.id);
                    break;
            }
        };
    }

    process(inputs, outputs, parameters) {
        this.#preciseTime = currentTime;

        this.processCommands();

        const output = outputs[0];

        for (let i = 0; i < output.length; i++) {
            this.runCommands()
            this.#preciseTime += sampleTime

            let mixL= 0;
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

            console.log(mixL, mixR);

            output[0][i] = mixL;
            output[1][i] = mixR;
        }

        return true; // continue processing
    }

    processCommands() {
        const bytesRead = this.#commands.read(this.#buffer);

        if (bytesRead > 0) {
            console.log("Processing commands", this.#buffer.subarray(0, bytesRead));

            let off = 0
            while (off < bytesRead) {
                const command = this.#bufferView.getUint8(off);
                const cmd = commandPool.get();

                cmd.command = command;
                cmd.ch = this.#bufferView.getUint8(off + 1);

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
                    commandPool.put(cmd);
                }
            }

            // TODO use manual merge sort: all new commands must be sorted first then merge two arrays

            for (let i = 0; i < chanLen; i++) {
                this.#commandsByTime[i].sort((a, b) => a.time - b.time);
            }
        }
    }

    runCommands() {
        for (let i = 0; i < chanLen; i++) {
            const selectedChan = this.#channels[i];
            let processed = 0;

            for (const cmd of this.#commandsByTime[i]) {
                if (cmd.time > this.#preciseTime) {
                    break;
                }

                switch (cmd.command) {
                    case cmdSetSample:
                        switch (true) {
                            case cmd.id === 0:
                                selectedChan.active = false;
                                selectedChan.sampleData = null;
                                break;
                            case this.samplesByID[cmd.id] === null:
                                console.log("[piaudio] SetSample failed: Sample not found, id:", cmd.id)
                                selectedChan.active = false;
                                selectedChan.sampleData = null;
                                break;
                            default:
                                selectedChan.active = true;
                                const sample = this.samplesByID[cmd.id];
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

                commandPool.put(cmd);

                processed++;
            }

            if (processed > 0) {
                this.#commandsByTime[i].copyWithin(0, processed);
                this.#commandsByTime[i].length -= processed;
                console.log("Processed commands", processed);
            }
        }
    }
}

registerProcessor("pi-audio-processor", PiAudioProcessor);