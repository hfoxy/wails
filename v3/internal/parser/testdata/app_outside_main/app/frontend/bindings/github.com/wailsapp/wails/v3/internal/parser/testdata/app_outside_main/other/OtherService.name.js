// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "@wailsio/runtime";

import * as models from "../../../../../../../../../models/models.js";

/**
 * @param {models.Person} person
 * @returns {Promise<string>}
 */
export function Greet(person) {
    return Call.ByName("github.com/wailsapp/wails/v3/internal/parser/testdata/app_outside_main/other.OtherService.Greet", person);
}
