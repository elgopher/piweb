// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

var api; // all Go functions are available here

(function startGameLoop() {
    api.init()

    requestAnimationFrame(function tick() {
        api.update()
        api.draw()

        requestAnimationFrame(tick)
    })
})()