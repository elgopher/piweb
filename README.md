# piweb

**Experimental Pi backend for web browsers.**
Powered by [TinyGo](https://tinygo.org/) and the [Audio Worklet API](https://developer.mozilla.org/en-US/docs/Web/API/AudioWorklet).

---

## 🎯 Project Goals

The aim of this project is to create a new backend for [Pi](https://github.com/elgopher/pi) games that runs efficiently in modern web browsers and offers significant improvements over the standard [piebiten](https://github.com/elgopher/pi/piebiten) backend:

### 🔊 Better Audio

* **Minimal latency** — as low as **20 ms** (compared to 60 ms in piebiten)
* **Glitch-free playback** — thanks to audio processing in a high-priority, separate audio thread (via Audio Worklet)

### 📦 Smaller Binary Size

* At least **2× smaller** `.wasm` output
* **No third-party dependencies**

### ⚡ Higher Performance

* Games will run at **significantly higher frame rates**
* Lower CPU usage — better performance on mobile and low-end devices

---

## 🧪 How It Works

These improvements are possible thanks to:

* **Audio Worklet API** – the browser-native API for real-time, low-latency audio processing
* **TinyGo compiler** – a lightweight Go compiler that generates **small and fast WebAssembly binaries**

---

## 🚧 Status

> This is an experimental project.
> Not all Pi features are supported yet.
