// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "/wails/runtime.js";

/**
 * Fatal is equivalent to l.Print() followed by a call to [os.Exit](1).
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Fatal(...v) {
    return Call.ByID(960027757, v);
}

/**
 * Fatalf is equivalent to l.Printf() followed by a call to [os.Exit](1).
 * @param {string} format
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Fatalf(format, ...v) {
    return Call.ByID(528639313, format, v);
}

/**
 * Fatalln is equivalent to l.Println() followed by a call to [os.Exit](1).
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Fatalln(...v) {
    return Call.ByID(3643628871, v);
}

/**
 * Flags returns the output flags for the logger.
 * The flag bits are [Ldate], [Ltime], and so on.
 * @returns {Promise<number>}
 */
export function Flags() {
    return Call.ByID(2740549504);
}

/**
 * Output writes the output for a logging event. The string s contains
 * the text to print after the prefix specified by the flags of the
 * Logger. A newline is appended if the last character of s is not
 * already a newline. Calldepth is used to recover the PC and is
 * provided for generality, although at the moment on all pre-defined
 * paths it will be 2.
 * @param {number} calldepth
 * @param {string} s
 * @returns {Promise<void>}
 */
export function Output(calldepth, s) {
    return Call.ByID(2365224640, calldepth, s);
}

/**
 * Panic is equivalent to l.Print() followed by a call to panic().
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Panic(...v) {
    return Call.ByID(3360441208, v);
}

/**
 * Panicf is equivalent to l.Printf() followed by a call to panic().
 * @param {string} format
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Panicf(format, ...v) {
    return Call.ByID(1846388794, format, v);
}

/**
 * Panicln is equivalent to l.Println() followed by a call to panic().
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Panicln(...v) {
    return Call.ByID(2472961110, v);
}

/**
 * Prefix returns the output prefix for the logger.
 * @returns {Promise<string>}
 */
export function Prefix() {
    return Call.ByID(4101446653);
}

/**
 * Print calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Print].
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Print(...v) {
    return Call.ByID(239783172, v);
}

/**
 * Printf calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Printf].
 * @param {string} format
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Printf(format, ...v) {
    return Call.ByID(3787542854, format, v);
}

/**
 * Println calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Println].
 * @param {any[]} v
 * @returns {Promise<void>}
 */
export function Println(...v) {
    return Call.ByID(2875326178, v);
}

/**
 * SetFlags sets the output flags for the logger.
 * The flag bits are [Ldate], [Ltime], and so on.
 * @param {number} flag
 * @returns {Promise<void>}
 */
export function SetFlags(flag) {
    return Call.ByID(3820942102, flag);
}

/**
 * SetPrefix sets the output prefix for the logger.
 * @param {string} prefix
 * @returns {Promise<void>}
 */
export function SetPrefix(prefix) {
    return Call.ByID(3047817107, prefix);
}
