// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

/**
 * @callback factory
 * @return {Object} new object
 */

// Pool of objects. Used to reduce the number of allocations (which is crucial in hot paths).
export class Pool {
    #pool = [];
    /** @type factory */
    #factory;

    /** @param {factory} f */
    constructor(f) {
        this.#factory = f;
    }

    get() {
        return this.#pool.length ? this.#pool.pop() : this.#factory();
    }

    put(obj) {
        this.#pool.push(obj);
    }
}
