//// [tests/cases/conformance/types/members/typesWithProtectedConstructor.ts] ////

//// [typesWithProtectedConstructor.ts]
class C {
    protected constructor() { }
}

var c = new C(); // error C is protected
var r: () => void = c.constructor;

class C2 {
    protected constructor(x: number);
    protected constructor(x: any) { }
}

var c2 = new C2(); // error C2 is protected
var r2: (x: number) => void = c2.constructor;

//// [typesWithProtectedConstructor.js]
class C {
    constructor() { }
}
var c = new C(); // error C is protected
var r = c.constructor;
class C2 {
    constructor(x) { }
}
var c2 = new C2(); // error C2 is protected
var r2 = c2.constructor;


//// [typesWithProtectedConstructor.d.ts]
declare class C {
    protected constructor();
}
declare var c: any; // error C is protected
declare var r: () => void;
declare class C2 {
    protected constructor(x: number);
}
declare var c2: any; // error C2 is protected
declare var r2: (x: number) => void;
