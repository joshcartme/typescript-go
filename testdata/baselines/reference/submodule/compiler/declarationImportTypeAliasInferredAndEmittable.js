//// [tests/cases/compiler/declarationImportTypeAliasInferredAndEmittable.ts] ////

//// [foo.ts]
class Conn {
    constructor() { }
    item = 3;
    method() { }
}

export = Conn;
//// [usage.ts]
type Conn = import("./foo");
declare var x: Conn;

export class Wrap {
    connItem: number;
    constructor(c = x) {
        this.connItem = c.item;
    }
}


//// [foo.js]
"use strict";
class Conn {
    constructor() { }
    item = 3;
    method() { }
}
module.exports = Conn;
//// [usage.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Wrap = void 0;
class Wrap {
    connItem;
    constructor(c = x) {
        this.connItem = c.item;
    }
}
exports.Wrap = Wrap;


//// [foo.d.ts]
declare class Conn {
    constructor();
    item: number;
    method(): void;
}
export = Conn;
//// [usage.d.ts]
export declare class Wrap {
    connItem: number;
    constructor(c?: import("./foo"));
}
