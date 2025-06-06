//// [tests/cases/conformance/es6/destructuring/iterableArrayPattern10.ts] ////

//// [iterableArrayPattern10.ts]
class Bar { x }
class Foo extends Bar { y }
class FooIterator {
    next() {
        return {
            value: new Foo,
            done: false
        };
    }

    [Symbol.iterator]() {
        return this;
    }
}

function fun([a, b]) { }
fun(new FooIterator);

//// [iterableArrayPattern10.js]
class Bar {
    x;
}
class Foo extends Bar {
    y;
}
class FooIterator {
    next() {
        return {
            value: new Foo,
            done: false
        };
    }
    [Symbol.iterator]() {
        return this;
    }
}
function fun([a, b]) { }
fun(new FooIterator);
