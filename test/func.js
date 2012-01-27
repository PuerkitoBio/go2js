/* Generated by GoScript <github.com/kless/GoScript> */






function singleLine() { console.log("Hello world!"); }

var x = 10;

(function() {
	x = 13;
}());

function simpleFunc() {

	var max = function(a, b) {
		if (a > b) {
			return a;
		}
		return b;
	};

	var x = 3;
	var y = 4;
	var z = 5;

	var max_xy = max(x, y);
	var max_xz = max(x, z);

	alert("max(" + x + ", " + y + ") = " + max_xy + "\n");
	alert("max(" + x + ", " + z + ") = " + max_xz + "\n");
	alert("max(" + y + ", " + z + ") = " + max(y, z) + "\n");
}

function twoOuputValues() {

	var SumAndProduct = function(A, B) {
		return [A + B, A * B];
	};

	var x = 3;
	var y = 4;
	var _ = SumAndProduct(x, y), xPLUSy = _[0], xTIMESy = _[1];

	alert("" + x + " + " + y + " = " + xPLUSy + "\n");
	alert("" + x + " * " + y + " = " + xTIMESy + "\n");
}

function resultVariable() {


	var MySqrt = function(f) { var s = 0, ok = false;
		if (f > 0) {
			s = Math.sqrt(f), ok = true;
		}
		return [s, ok];
	};

	for (var i = -2.0; i <= 10; i++) {
		var _ = MySqrt(i), sqroot = _[0], ok = _[1];
		if (ok) {
			alert("The square root of " + i + " is " + sqroot + "\n");
		} else {
			alert("Sorry, no square root for " + i + "\n");
		}
	}
}

function testReturn_1() {
	var MySqrt = function(f) { var squareroot = 0, ok = false;
		if (f > 0) {
			squareroot = Math.sqrt(f), ok = true;
		}
		return [squareroot, ok];
	};

	var check = MySqrt(5)[1];
	alert(check + "\n");
}

function testReturn_2(n) { var ok = false;
	if (n > 0) {
		ok = true;
	}
	return ok;
}

function testPanic() {
	throw new Error("unreachable");
	throw new Error("not implemented: " + "foo" + "");
}

function main() {
	console.log("\n== singleLine()\n\n");
	singleLine();

	console.log("\n== simpleFunc()\n\n");
	simpleFunc();

	console.log("\n== twoOuputValues()\n\n");
	twoOuputValues();

	console.log("\n== resultVariable()\n\n");
	resultVariable();

	console.log("\n== testReturn_1()\n\n");
	testReturn_1();

	console.log("\n== testReturn_2(-1)\n\n");
	console.log(testReturn_2(-1) + "\n");

	console.log("\n== testPanic()\n\n");
	testPanic();
}
