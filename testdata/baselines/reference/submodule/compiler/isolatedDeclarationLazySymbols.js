//// [tests/cases/compiler/isolatedDeclarationLazySymbols.ts] ////

//// [isolatedDeclarationLazySymbols.ts]
export function foo() {

}

const o = {
    ["prop.inner"]: "a",
    prop: {
        inner: "b",
    }
} as const

foo[o["prop.inner"]] ="A";
foo[o.prop.inner] = "B";

export class Foo {
    [o["prop.inner"]] ="A"
    [o.prop.inner] = "B"
}

export let oo = {
    [o['prop.inner']]:"A",
    [o.prop.inner]: "B",
}

//// [isolatedDeclarationLazySymbols.js]
export function foo() {
}
const o = {
    ["prop.inner"]: "a",
    prop: {
        inner: "b",
    }
};
foo[o["prop.inner"]] = "A";
foo[o.prop.inner] = "B";
export class Foo {
    [o["prop.inner"]] = "A"[o.prop.inner] = "B";
}
export let oo = {
    [o['prop.inner']]: "A",
    [o.prop.inner]: "B",
};


//// [isolatedDeclarationLazySymbols.d.ts]
export declare function foo(): void;
export declare class Foo {
}
export declare let oo: {
    a: string;
    b: string;
};
