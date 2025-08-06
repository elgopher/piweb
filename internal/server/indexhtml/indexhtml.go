// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package indexhtml

import "bytes"

var scripts = []byte(`
<script src="wasm_exec.js"></script>
<script src="bytebuffer.js"></script>
<script src="gamepad.js"></script>
<script src="gameloop.js"></script>
<script src="canvas.js"></script>

<script>
    window.addEventListener("load", async () => {
        const resp = await fetch("main.wasm");
        if (!resp.ok) {
            const errorElement = document.createElement('pre');
            errorElement.innerText = await resp.text();
            document.body.appendChild(errorElement);
            return;
        }

        const go = new Go();
        const result = await WebAssembly.instantiateStreaming(resp, go.importObject);
        go.run(result.instance);

        startGameLoop();
    });
</script>
`)

func PutScripts(content []byte) []byte {
	return bytes.Replace(content, []byte("$$SCRIPTS$$"), scripts, 1)
}
