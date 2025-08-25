export {}; // Make this file a module so we can safely augment globals.

declare global {
    /** Global AudioWorklet timeline time in seconds (renderer time). */
    const currentTime: number;

    /** Global frame counter (number of rendered sample frames). */
    const currentFrame: number;

    /** Global sample rate in Hz. */
    const sampleRate: number;

    interface AudioParamDescriptor {
        name: string;
        automationRate?: 'a-rate' | 'k-rate';
        minValue?: number;
        maxValue?: number;
        defaultValue?: number;
    }

    interface AudioWorkletNodeOptions {
        numberOfInputs?: number;
        numberOfOutputs?: number;
        outputChannelCount?: number[];
        parameterData?: Record<string, number>;
        processorOptions?: any;
    }

    interface AudioWorkletProcessor {
        /** Bidirectional message port to the main thread. */
        readonly port: MessagePort;

        /**
         * Render callback. Return true to continue processing, false to terminate.
         * - inputs/outputs: [node][channel] -> Float32Array (128 frames per quantum)
         * - parameters: map from parameter name to a Float32Array of values for this quantum
         */
        process(
            inputs: ReadonlyArray<ReadonlyArray<Float32Array>>,
            outputs: ReadonlyArray<ReadonlyArray<Float32Array>>,
            parameters: Record<string, Float32Array>
        ): boolean;
    }

    interface AudioWorkletProcessorConstructor {
        new (options?: AudioWorkletNodeOptions): AudioWorkletProcessor;
        /** Optional static parameter descriptors. */
        readonly parameterDescriptors?: ReadonlyArray<AudioParamDescriptor>;
    }

    /** Base class provided in the AudioWorkletGlobalScope. */
    const AudioWorkletProcessor: {
        prototype: AudioWorkletProcessor;
    };

    /**
     * Register a processor under a name used by AudioWorkletNode on the main thread.
     */
    function registerProcessor(
        name: string,
        processorCtor: AudioWorkletProcessorConstructor
    ): void;
}