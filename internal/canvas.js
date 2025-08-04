// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// piweb uses 2D canvas for gfx rendering.
// Canvas has the size equal to Pi screen.
// piweb scales the canvas by using CSS styles, though.
// piweb uses pixel-perfect integer scaling.

var canvas = newCanvas();
var imageData; // https://developer.mozilla.org/en-US/docs/Web/API/ImageData

function prepareCanvas() {
    canvas = newCanvas();
    resizeCanvas();
    document.body.appendChild(canvas);
    centerCanvasOnTheScreen();
    window.addEventListener("resize", resizeCanvas);
    canvas.addEventListener("contextmenu", (e) => {
        e.preventDefault();
    });
}

function newCanvas() {
    const canvas = document.createElement("canvas");
    const css = canvas.style;
    css.imageRendering = "pixelated";
    return canvas;
}

function resizeCanvas() {
    const scale = screenScale();

    canvas.width = api.screenWidth;
    canvas.height = api.screenHeight;
    canvas.style.width = api.screenWidth / devicePixelRatio * scale + "px";
    canvas.style.height = api.screenHeight / devicePixelRatio * scale + "px";

    const ctx = canvas.getContext("2d");
    imageData = ctx.createImageData(canvas.width, canvas.height);
}

function screenScale() {
    const realw = window.innerWidth * devicePixelRatio;
    const realh = window.innerHeight * devicePixelRatio;
    const sw = realw / api.screenWidth;
    const sh = realh / api.screenHeight;
    let scale = Math.floor(Math.min(sw, sh));
    if (scale === 0) {
        scale = 1;
    }
    return scale;
}

const prevMouseXY = {X: 0, Y: 0}

function addMouseMoveListener(canvas, callback) {
    canvas.addEventListener("mousemove", (e) => {
        // mousemove event is reported many times during a single frame.
        // Filter events before calling Go code to avoid significant
        // memory allocations.
        const m = devicePixelRatio/screenScale()
        const x = Math.floor(e.offsetX * m)
        const y = Math.floor(e.offsetY * m)

        if (prevMouseXY.x === x && prevMouseXY.y === y) {
            return
        }

        callback(x,y)

        prevMouseXY.x = x
        prevMouseXY.y = y
    })
}

function centerCanvasOnTheScreen() {
    document.documentElement.style.height = "100%";
    const body = document.body.style;
    body.margin = "0px";
    body.display = "grid";
    body.placeItems = "center";
    body.height = "100%";
}

function updateCanvas() {
    const ctx = canvas.getContext("2d");
    ctx.putImageData(imageData, 0, 0);
}