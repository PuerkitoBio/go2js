/* Generated by GoScript <github.com/kless/GoScript> */



function initialValue() {
	var s1 = g.NilSlice();
	var s2 = g.NewSlice([]);
	var s3 = g.NewSlice([2]);
	var s4 = g.NewSlice([2, 4]);
	var s5 = g.MakeSlice(0, 0);
	var s6 = g.MakeSlice(0, 5);
	var s7 = g.MakeSlice(0, 5, 10);

	var msg = "nil";
	if (s1.isNil && !s2.isNil && !s3.isNil && !s4.isNil && !s5.isNil && !s6.isNil && !s7.isNil) {

		document.write("[OK] " + msg + "<br>");
	} else {
		document.write("[Error] " + msg + "<br>");
	}

	msg = "length";
	if (s1.len === 0 && s2.len === 0 && s3.len === 1 && s4.len === 2 && s5.len === 0 && s6.len === 5 && s7.len === 5) {

		document.write("[OK] " + msg + "<br>");
	} else {
		document.write("[Error] " + msg + "<br>");
	}

	msg = "capacity";
	if (s1.cap === 0 && s2.cap === 0 && s3.cap === 1 && s4.cap === 2 && s5.cap === 0 && s6.cap === 5 && s7.cap === 10) {

		document.write("[OK] " + msg + "<br>");
	} else {
		document.write("[Error] " + msg + "<br>");
	}
}

function shortHand() {

	var array = g.MakeArray([10], 0, ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j']);

	var a_slice = g.NilSlice(), b_slice = g.NilSlice();

	var msg = "slicing";

	a_slice.set(array, 4, 8);

	if (a_slice.toString() === "efgh" && a_slice.len === 4 && a_slice.cap === 6) {
		document.write("[OK] " + msg + "<br>");
	} else {
		document.write("[Error] " + msg + "<br>");
	}


	a_slice.set(array, 6, 7);

	if (a_slice.toString() === "g") {
		document.write("[OK]<br>");
	} else {
		document.write("[Error]<br>");
	}


	msg = "shorthand";

	a_slice.set(array, 0, 3);

	if (a_slice.toString() === "abc" && a_slice.len === 3 && a_slice.cap === 10) {
		document.write("[OK] " + msg + "<br>");
	} else {
		document.write("[Error] " + msg + "<br>");
	}


	a_slice.set(array, 5);

	if (a_slice.toString() === "fghij") {
		document.write("[OK]<br>");
	} else {
		document.write("[Error]<br>");
	}


	a_slice.set(array, 0);

	if (a_slice.toString() === "abcdefghij") {
		document.write("[OK]<br>");
	} else {
		document.write("[Error]<br>");
	}


	msg = "slice of a slice";

	a_slice.set(array, 3, 7);

	if (a_slice.toString() === "defg" && a_slice.len === 4 && a_slice.cap === 7) {
		document.write("[OK] " + msg + "<br>");
	} else {
		document.write("[Error] " + msg + "<br>");
	}


	b_slice.set(a_slice, 1, 3);

	if (b_slice.toString() === "ef" && b_slice.len === 2 && b_slice.cap === 6) {
		document.write("[OK]<br>");
	} else {
		document.write("[Error]<br>");
	}


	b_slice.set(a_slice, 0, 3);

	if (b_slice.toString() === "def") {
		document.write("[OK]<br>");
	} else {
		document.write("[Error]<br>");
	}


	b_slice.set(a_slice, 0);

	if (b_slice.toString() === "defg") {
		document.write("[OK]<br>");
	} else {
		document.write("[Error]<br>");
	}

}




function Max(slice) {
	var max = slice[0];
	for (var index = 1; index < slice.length; index++) {
		if (slice[index] > max) {
			max = slice[index];
		}
	}
	return max;
}

function useFunc() {

	var A1 = g.MakeArray([10], 0, [1, 2, 3, 4, 5, 6, 7, 8, 9]);
	var A2 = g.MakeArray([4], 0, [1, 2, 3, 4]);
	var A3 = g.MakeArray([1], 0, [1]);


	var slice = g.NilSlice();

	slice.set(A1, 0);

	if (Max(slice.f) === 9) {
		document.write("[OK] A1<br>");
	} else {
		document.write("[Error] A1<br>");
	}


	slice.set(A2, 0);

	if (Max(slice.f) === 4) {
		document.write("[OK] A2<br>");
	} else {
		document.write("[Error] A2<br>");
	}


	slice.set(A3, 0);

	if (Max(slice.f) === 1) {
		document.write("[OK] A3<br>");
	} else {
		document.write("[Error] A3<br>");
	}

}



function PrintByteSlice(name, slice) {
	var s = "" + name + " is : [";
	for (var index = 0; index < slice.len - 1; index++) {
		s += "" + slice.f[index] + ",";
	}
	s += "" + slice.f[slice.len - 1] + "]";

	document.write(s + "<br>");
	return s;
}

function reference() {

	var A = g.MakeArray([10], 0, ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j']);


	var slice1 = g.NewSlice(A, 3, 7);
	var slice2 = g.NewSlice(A, 5);
	var slice3 = g.NewSlice(slice1, 0, 2);


	document.write("=== First content of A and the slices<br>");
	PrintByteSlice("A", g.NewSlice(A, 0));
	PrintByteSlice("slice1", slice1);
	PrintByteSlice("slice2", slice2);
	PrintByteSlice("slice3", slice3);


	A[4] = 'E';
	document.write("<br>=== Content of A and the slices, after changing 'e' to 'E' in array A<br>");
	PrintByteSlice("A", g.NewSlice(A, 0));
	PrintByteSlice("slice1", slice1);
	PrintByteSlice("slice2", slice2);
	PrintByteSlice("slice3", slice3);


	slice2[1] = 'G';
	document.write("<br>=== Content of A and the slices, after changing 'g' to 'G' in slice2<br>");
	PrintByteSlice("A", g.NewSlice(A, 0));
	PrintByteSlice("slice1", slice1);
	PrintByteSlice("slice2", slice2);
	PrintByteSlice("slice3", slice3);
}



function resize() {
	var slice = g.NilSlice();


	slice = g.MakeSlice(0, 4, 5);

	if (slice.len === 4 && slice.cap === 5 && slice.f[0] === 0 && slice.f[1] === 0 && slice.f[2] === 0 && slice.f[3] === 0) {

		document.write("[OK] allocation<br>");
	} else {
		document.write("[Error] allocation<br>");
	}



	slice.f[1] = 2, slice.f[3] = 3;

	if (slice.f[0] === 0 && slice.f[1] === 2 && slice.f[2] === 0 && slice.f[3] === 3) {
		document.write("[OK] change<br>");
	} else {
		document.write("[Error] change<br>");
	}



	slice = g.MakeSlice(0, 2);

	if (slice.len === 2 && slice.cap === 2 && slice.f[0] === 0 && slice.f[1] === 0) {
		document.write("[OK] resize<br>");
	} else {
		document.write("[Error] resize<br>");
	}

}



function main() {
	document.write("<br>== initialValue<br>");
	initialValue();
	document.write("<br>== shortHand<br>");
	shortHand();
	document.write("<br>== useFunc<br>");
	useFunc();
	document.write("<br>== reference<br>");
	reference();
	document.write("<br>== resize<br>");
	resize();
} main();
