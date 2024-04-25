// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export class Aliases {
    "Name": string;
    "Numbers": number[];
    "Settings": { [_: string]: string };
    "Person": Human;

    /** Creates a new Aliases instance. */
    constructor(source: Partial<Aliases> = {}) {
        if (!("Name" in source)) {
            this["Name"] = "";
        }
        if (!("Numbers" in source)) {
            this["Numbers"] = [];
        }
        if (!("Settings" in source)) {
            this["Settings"] = {};
        }
        if (!("Person" in source)) {
            this["Person"] = (new Human());
        }

        Object.assign(this, source);
    }

    /** Creates a new Aliases instance from a string or object. */
    static createFrom(source: string | object = {}): Aliases {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Aliases(parsedSource as Partial<Aliases>);
    }
}

/**
 * Human is a copy of Person
 * it is only visible without GODEBUG=gotypesalias=1
 */
export class Human {
    "name": string;

    /** Creates a new Human instance. */
    constructor(source: Partial<Human> = {}) {
        if (!("name" in source)) {
            this["name"] = "";
        }

        Object.assign(this, source);
    }

    /** Creates a new Human instance from a string or object. */
    static createFrom(source: string | object = {}): Human {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Human(parsedSource as Partial<Human>);
    }
}
