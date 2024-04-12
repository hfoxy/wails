// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "@wailsio/runtime";

/**
 * Fatal is equivalent to l.Print() followed by a call to [os.Exit](1).
 */
export function Fatal(...v: any[]): Promise<void> {
    return Call.ByID(960027757, v);
}

/**
 * Fatalf is equivalent to l.Printf() followed by a call to [os.Exit](1).
 */
export function Fatalf(format: string, ...v: any[]): Promise<void> {
    return Call.ByID(528639313, format, v);
}

/**
 * Fatalln is equivalent to l.Println() followed by a call to [os.Exit](1).
 */
export function Fatalln(...v: any[]): Promise<void> {
    return Call.ByID(3643628871, v);
}

/**
 * Flags returns the output flags for the logger.
 * The flag bits are [Ldate], [Ltime], and so on.
 */
export function Flags(): Promise<number> {
    return Call.ByID(2740549504);
}

/**
 * Output writes the output for a logging event. The string s contains
 * the text to print after the prefix specified by the flags of the
 * Logger. A newline is appended if the last character of s is not
 * already a newline. Calldepth is used to recover the PC and is
 * provided for generality, although at the moment on all pre-defined
 * paths it will be 2.
 */
export function Output(calldepth: number, s: string): Promise<void> {
    return Call.ByID(2365224640, calldepth, s);
}

/**
 * Panic is equivalent to l.Print() followed by a call to panic().
 */
export function Panic(...v: any[]): Promise<void> {
    return Call.ByID(3360441208, v);
}

/**
 * Panicf is equivalent to l.Printf() followed by a call to panic().
 */
export function Panicf(format: string, ...v: any[]): Promise<void> {
    return Call.ByID(1846388794, format, v);
}

/**
 * Panicln is equivalent to l.Println() followed by a call to panic().
 */
export function Panicln(...v: any[]): Promise<void> {
    return Call.ByID(2472961110, v);
}

/**
 * Prefix returns the output prefix for the logger.
 */
export function Prefix(): Promise<string> {
    return Call.ByID(4101446653);
}

/**
 * Print calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Print].
 */
export function Print(...v: any[]): Promise<void> {
    return Call.ByID(239783172, v);
}

/**
 * Printf calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Printf].
 */
export function Printf(format: string, ...v: any[]): Promise<void> {
    return Call.ByID(3787542854, format, v);
}

/**
 * Println calls l.Output to print to the logger.
 * Arguments are handled in the manner of [fmt.Println].
 */
export function Println(...v: any[]): Promise<void> {
    return Call.ByID(2875326178, v);
}

/**
 * SetFlags sets the output flags for the logger.
 * The flag bits are [Ldate], [Ltime], and so on.
 */
export function SetFlags(flag: number): Promise<void> {
    return Call.ByID(3820942102, flag);
}

/**
 * SetPrefix sets the output prefix for the logger.
 */
export function SetPrefix(prefix: string): Promise<void> {
    return Call.ByID(3047817107, prefix);
}

export function output(pc: number, calldepth: number, appendOutput: any): Promise<void> {
    return Call.ByID(3540906208, pc, calldepth, appendOutput);
}
