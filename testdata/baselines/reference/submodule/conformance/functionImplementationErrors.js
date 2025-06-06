//// [tests/cases/conformance/functions/functionImplementationErrors.ts] ////

//// [functionImplementationErrors.ts]
// FunctionExpression with no return type annotation with multiple return statements with unrelated types
var f1 = function () {
    return '';
    return 3;
};
var f2 = function x() {
    return '';
    return 3;
};
var f3 = () => {
    return '';
    return 3;
};

// FunctionExpression with no return type annotation with return branch of number[] and other of string[]
var f4 = function () {
    if (true) {
        return [''];
    } else {
        return [1];
    }
}

// Function implemetnation with non -void return type annotation with no return
function f5(): number {
}

var m;
// Function signature with parameter initializer referencing in scope local variable
function f6(n = m) {
    var m = 4;
}

// Function signature with initializer referencing other parameter to the right
function f7(n = m, m?) {
}

// FunctionExpression with non -void return type annotation with a throw, no return, and other code
// Should be error but isn't
undefined === function (): number {
    throw undefined;
    var x = 4;
};

class Base { private x; }
class AnotherClass { private y; }
class Derived1 extends Base { private m; }
class Derived2 extends Base { private n; }
function f8() {
    return new Derived1();
    return new Derived2();    
}
var f9 = function () {
    return new Derived1();
    return new Derived2();
};
var f10 = () => {
    return new Derived1();
    return new Derived2();
};
function f11() {
    return new Base();
    return new AnotherClass();
}
var f12 = function () {
    return new Base();
    return new AnotherClass();
};
var f13 = () => {
    return new Base();
    return new AnotherClass();
};


//// [functionImplementationErrors.js]
// FunctionExpression with no return type annotation with multiple return statements with unrelated types
var f1 = function () {
    return '';
    return 3;
};
var f2 = function x() {
    return '';
    return 3;
};
var f3 = () => {
    return '';
    return 3;
};
// FunctionExpression with no return type annotation with return branch of number[] and other of string[]
var f4 = function () {
    if (true) {
        return [''];
    }
    else {
        return [1];
    }
};
// Function implemetnation with non -void return type annotation with no return
function f5() {
}
var m;
// Function signature with parameter initializer referencing in scope local variable
function f6(n = m) {
    var m = 4;
}
// Function signature with initializer referencing other parameter to the right
function f7(n = m, m) {
}
// FunctionExpression with non -void return type annotation with a throw, no return, and other code
// Should be error but isn't
undefined === function () {
    throw undefined;
    var x = 4;
};
class Base {
    x;
}
class AnotherClass {
    y;
}
class Derived1 extends Base {
    m;
}
class Derived2 extends Base {
    n;
}
function f8() {
    return new Derived1();
    return new Derived2();
}
var f9 = function () {
    return new Derived1();
    return new Derived2();
};
var f10 = () => {
    return new Derived1();
    return new Derived2();
};
function f11() {
    return new Base();
    return new AnotherClass();
}
var f12 = function () {
    return new Base();
    return new AnotherClass();
};
var f13 = () => {
    return new Base();
    return new AnotherClass();
};
