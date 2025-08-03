// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

var api; // all Go functions are available here

// Ticker runs a callback on every tick based on pi.TPS
class Ticker {
    #tickStartTime = 0;

    constructor(tickCallback) {
        this.tickCallback = tickCallback;
    }

    start() {
        this.#tickStartTime = performance.now();
        requestAnimationFrame(this.#frame);
    }

    #frame = () => {
        const tickDuration = 1000 / api.tps;
        const now = performance.now();

        const elapsed = now - this.#tickStartTime;

        let ticks = Math.round(elapsed / tickDuration);
        let actualTicks = ticks;
        if (ticks > api.tps) {
            console.warn("Too many ticks missed: %d. Dropping %d ticks", ticks, ticks - api.tps);
            actualTicks = api.tps;
        }
        if (actualTicks > 0) {
            this.tickCallback(actualTicks);
            this.#tickStartTime = this.#tickStartTime + tickDuration * ticks;
        }

        requestAnimationFrame(this.#frame);
    }
}


(function startGameLoop() {
    api.init();

    const ticker = new Ticker(api.tick);
    ticker.start();
})()