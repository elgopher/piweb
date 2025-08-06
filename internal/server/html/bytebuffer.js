// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

class ByteBuffer {
    buf = new Uint8Array(1);
    #len = 0;

    write(b) {
        if (this.#len === this.buf.byteLength) {
            const resizedBuf = new Uint8Array(this.buf.byteLength * 2);
            resizedBuf.set(this.buf);
            this.buf = resizedBuf;
        }

        this.buf[this.#len] = b;
        this.#len++;
    }

    clear() {
        this.#len = 0;
    }

    length() {
        return this.#len;
    }
}
