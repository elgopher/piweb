# How to Run

## Live Testing (WebAssembly)

To run your WebAssembly game in the browser with live reloading:

1. Start the web server in this directory:

   ```bash
   cd _example
   go run .
   ```

2. Open your browser and go to:

   ```
   http://localhost:8080
   ```

3. **Refresh** the browser to recompile the code and restart the game

# How to change the template

1. You can modify all files like index.html, main.css etc.
2. Just put the altered file into _example directory. You can find all default static files in [html directory](../internal/server/html)
3. Of course, you can add new files as well.