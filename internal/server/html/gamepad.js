// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

const gamepadEventTypeUp = 1;
const gamepadEventTypeDown = 2;

const gamepadButtons = [0, 1, 2, 3, 12, 13, 14, 15];

class Gamepad {
    events = new ByteBuffer();
    #state = [];
    #reusedState = new GamepadState();

    constructor() {
        for (let i = 0; i < 4; i++) {
            this.#state[i] = new GamepadState();
        }
    }

    onAnimationFrame() {
        const snapshot = navigator.getGamepads() // This cannot be avoided. Allocates a lot!
        for (let i = 0; i < 4; i++) {
            this.#reusedState.setState(snapshot[i]);
            this.#state[i].publishEvents(this.#reusedState, i, this.events);
        }
    }
}

class GamepadState {
    pressed = {};

    constructor() {
        for (const btn of gamepadButtons) {
            this.pressed[btn] = false;
        }
    }

    publishEvents(newState, idx, eventsByteBuffer) {
        for (const btn of gamepadButtons) {
            const newBtnState = newState.pressed[btn];
            const oldBtnState = this.pressed[btn];
            if (newBtnState !== oldBtnState) {
                const eventType = !newBtnState ? gamepadEventTypeUp : gamepadEventTypeDown;
                eventsByteBuffer.write(eventType);
                eventsByteBuffer.write(idx);
                eventsByteBuffer.write(btn);
                this.pressed[btn] = newBtnState;
            }
        }
    }

    setState(snapshot) {
        if (snapshot != null) {
            for (const btn of gamepadButtons) {
                this.pressed[btn] = snapshot.buttons[btn].pressed;
            }

            const axes = snapshot.axes;
            
            const verticalAxis = axes[1];
            if (verticalAxis > 0.5) {
                this.pressed[13] = true;
            }
            if (verticalAxis < -0.5) {
                this.pressed[12] = true;
            }

            const horizontalAxis = axes[0];
            if (horizontalAxis > 0.5) {
                this.pressed[15] = true;
            }
            if (horizontalAxis < -0.5) {
                this.pressed[14] = true;
            }
        } else {
            for (const btn of gamepadButtons) {
                this.pressed[btn] = false;
            }
        }
    }
}

var gamepad = new Gamepad();