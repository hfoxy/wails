// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Call} from "@wailsio/runtime";

import {$anon1, Person} from "./models";

/**
 * Greet does XYZ
 * It has a multiline doc comment
 * The comment has even some *\/ traps!!
 */
export function Greet(str: string, people: Person[], $2: $anon1, assoc: { [_: string]: boolean | null }, $4: (number | null)[], ...other: string[]): Promise<[Person, any, number[]]> {
    return Call.ByName("main.GreetService.Greet", str, people, $2, assoc, $4, other);
}
