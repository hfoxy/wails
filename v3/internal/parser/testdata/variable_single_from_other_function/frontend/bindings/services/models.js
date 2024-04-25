// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export class Address {
    /**
     * Creates a new Address instance.
     * @param {Partial<Address>} [source = {}] - The source object to create the Address.
     */
    constructor(source = {}) {
        if (!("Street" in source)) {
            /**
             * @member
             * @type {string}
             */
            this["Street"] = "";
        }
        if (!("State" in source)) {
            /**
             * @member
             * @type {string}
             */
            this["State"] = "";
        }
        if (!("Country" in source)) {
            /**
             * @member
             * @type {string}
             */
            this["Country"] = "";
        }

        Object.assign(this, source);
    }

    /**
     * Creates a new Address instance from a string or object.
     * @param {string|object} source - The source data to create a Address instance from.
     * @returns {Address} A new Address instance.
     */
    static createFrom(source) {
        let parsedSource = typeof source === 'string' ? JSON.parse(source) : source;
        return new Address(parsedSource);
    }
};
