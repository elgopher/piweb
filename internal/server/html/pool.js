// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// pool.js
export function createPool(factory) {
    const pool = [];
    return {
        get() { return pool.length ? pool.pop() : factory(); },
        put(obj) { pool.push(obj); }
    };
}