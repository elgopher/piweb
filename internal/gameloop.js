// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

var api; // all Go functions are available here

// Ticker runs a callback on every tick based on pi.TPS
class Ticker {
    tickStartTime = 0

    constructor(tickCallback) {
        this.tickCallback = tickCallback
        this.frame = this.frame.bind(this) // needed because executed from RAF
    }

    start() {
        this.tickStartTime = performance.now()
        requestAnimationFrame(this.frame)
    }

    frame() {
        const tickDuration = 1000 / api.tps
        const now = performance.now()

        const elapsed = now - this.tickStartTime;

        let ticks = Math.round(elapsed / tickDuration)
        if (ticks > 60) {
            console.warn("Too many ticks missed: %d. Dropping %d ticks", ticks, ticks - 60)
            ticks = 60
        }
        for (let i = 0; i < ticks; i++) {
            this.tickCallback()
            this.tickStartTime = this.tickStartTime + tickDuration
        }

        requestAnimationFrame(this.frame)
    }
}


(function startGameLoop() {
    api.init()

    const ticker = new Ticker(api.tick)
    ticker.start()

})()