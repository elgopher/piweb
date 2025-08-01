# How to Run

## Live Testing (WebAssembly)

To run your WebAssembly game in the browser with live reloading:

1. Install `wasmserve`:

   ```bash
   go install github.com/hajimehoshi/wasmserve@latest
   ```

2. Start the server in this directory:

   ```bash
   cd _example
   wasmserve
   ```

3. Open your browser and go to:

   ```
   http://localhost:8080
   ```
