// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

export class SharedRing {
    static HEADER_BYTES = 12; // Int32[3]: [capacity, head, tail]

    constructor(sab) {
        if (!(sab instanceof SharedArrayBuffer)) throw new TypeError("SharedArrayBuffer is required");
        this.ctrl = new Int32Array(sab, 0, 3);           // [cap, head, tail]
        this.cap = sab.byteLength - SharedRing.HEADER_BYTES;
        if (this.cap <= 0) throw new RangeError("SharedArrayBuffer too small");

        // header initialization (idempotent)
        if (Atomics.load(this.ctrl, 0) !== this.cap) {
            Atomics.store(this.ctrl, 0, this.cap);       // capacity
            Atomics.store(this.ctrl, 1, 0);        // head
            Atomics.store(this.ctrl, 2, 0);        // tail
        }
        this.data = new Uint8Array(sab, SharedRing.HEADER_BYTES, this.cap);
    }

    // how much space is left in the buffer
    freeSpace() {
        const H = Atomics.load(this.ctrl, 1) >>> 0;
        const T = Atomics.load(this.ctrl, 2) >>> 0;
        const used = (H - T) >>> 0;
        return this.cap - used;
    }

    // Writes as many bytes as possible from src; returns number of bytes written.
    write(src) {
        const a = toU8(src);
        if (a.length === 0) return 0;

        const H = Atomics.load(this.ctrl, 1) >>> 0;
        const T = Atomics.load(this.ctrl, 2) >>> 0;
        const used = (H - T) >>> 0;
        const free = this.cap - used;
        if (free === 0) return 0;

        const n = Math.min(free, a.length);
        const w = H % this.cap;

        const first = Math.min(n, this.cap - w);
        this.data.set(a.subarray(0, first), w);
        if (n > first) this.data.set(a.subarray(first, n), 0);

        Atomics.store(this.ctrl, 1, (H + n) >>> 0); // publish head
        return n;
    }

    // Reads into dst; returns number of bytes read.
    read(dst) {
        const b = toU8(dst);
        if (b.length === 0) return 0;

        const H = Atomics.load(this.ctrl, 1) >>> 0;
        const T = Atomics.load(this.ctrl, 2) >>> 0;
        const avail = (H - T) >>> 0;
        if (avail === 0) return 0;

        const n = Math.min(avail, b.length);
        const r = T % this.cap;

        const first = Math.min(n, this.cap - r);
        b.set(this.data.subarray(r, r + first), 0);
        if (n > first) b.set(this.data.subarray(0, n - first), first);

        Atomics.store(this.ctrl, 2, (T + n) >>> 0); // publish tail
        return n;
    }
}

function toU8(x) {
    if (x instanceof Uint8Array) return x;
    if (ArrayBuffer.isView(x)) return new Uint8Array(x.buffer, x.byteOffset, x.byteLength);
    if (x instanceof ArrayBuffer) return new Uint8Array(x);
    throw new TypeError("Expected Uint8Array/ArrayBufferView/ArrayBuffer");
}
