# piweb

**Experimental Pi backend for web browsers.**

---

## 🎯 Project Goals

The aim of this project is to create a new backend for [Pi](https://github.com/elgopher/pi) games that runs efficiently in modern web browsers and offers significant improvements over the standard [piebiten](https://github.com/elgopher/pi/piebiten) backend:

### 🔊 Better Audio

* **Minimal latency** — as low as **20 ms** (compared to 60 ms in piebiten)
* **Glitch-free playback** — thanks to audio processing in a high-priority, separate audio thread

### 📦 Smaller Binary Size

* At least **2× smaller** `.wasm` output

### ⚡ Higher Performance

* Significantly lower number of memory allocations, therefore less CPU time spent on garbage collection
* Games will run at **higher frame rates**
* Lower CPU usage — better performance on mobile and low-end devices

---

## 🧪 How It Works

These improvements are possible thanks to:

* **Audio Worklet API** – the browser-native API for real-time, low-latency audio processing
* Writing the code directly in **JavaScript**
* Reducing the number of dependencies

---

## 🚧 Status

> This is an experimental project.
> Most features are not ready.
