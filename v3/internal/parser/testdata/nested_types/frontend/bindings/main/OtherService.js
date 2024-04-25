// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "@wailsio/runtime";

import {App, Field} from "./models.js";

/**
 * @param {App} app
 * @returns {Promise<Field[][]>}
 */
export function GetBoard(app) {
    return Call.ByID(1871244943, app);
}

/**
 * @param {App} app
 * @returns {Promise<{ [_: string]: { [_: string]: string } }>}
 */
export function GetLocale(app) {
    return Call.ByID(4038436675, app);
}

/**
 * @param {{ [_: string]: { [_: string]: string } }} locale
 * @param {Field[][]} board
 * @returns {Promise<App>}
 */
export function NewApp(locale, board) {
    return Call.ByID(3890694354, locale, board);
}
