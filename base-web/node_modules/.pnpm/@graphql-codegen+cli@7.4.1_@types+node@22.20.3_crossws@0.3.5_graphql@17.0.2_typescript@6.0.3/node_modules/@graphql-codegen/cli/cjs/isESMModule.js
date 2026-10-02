"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isESMModule = void 0;
/**
 * Poor mans ESM detection.
 * Looking at this and you have a better method?
 * Send a PR.
 */
exports.isESMModule = (typeof __dirname === 'string') === false;
