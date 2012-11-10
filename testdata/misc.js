









var PASS = true;



function argArray(arr) {
	var pass = true;

	if (arr.len() === 3 && arr.cap() === 3 && arr.v[0] === 1 && arr.v[1] === 2 && arr.v[2] === 3) {
	} else {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: argArray<br>");
		pass = false, PASS = false;
	}

	if (pass) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;pass<br>");
	}
	return arr;
}

function argEllipsis(arr) {
	var pass = true;

	if (arr.len() === 2 && arr.cap() === 2 && arr.v[0] === 5 && arr.v[1] === 6) {
	} else {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: argEllipsis<br>");
		pass = false, PASS = false;
	}

	if (pass) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;pass<br>");
	}
	return arr;
}

function argSlice(s) {
	var pass = true;

	if (s.len === 2 && s.cap === 2 && s.str() === "89" && s.get()[0] === '8' && s.get()[1] === '9') {
	} else {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: argSlice<br>");
		pass = false, PASS = false;
	}

	if (pass) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;pass<br>");
	}
	return s;
}

function argMap(m) {
	var pass = true;

	if (m.len() === 2 && m.get(1)[0] === "foo" && m.get(2)[0] === "bar") {
	} else {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: argSlice<br>");
		pass = false, PASS = false;
	}

	if (pass) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;pass<br>");
	}
	return m;
}

function main() {
	document.write("<br><br>== Miscellaneous<br><br>");

	document.write("=== RUN argArray<br>");
	var a = g.MkArray([3], 0, [1, 2, 3]);
	a = argArray(a);
	argArray(g.MkArray([3], 0, [1, 2, 3]));

	document.write("=== RUN argEllipsis<br>");
	var ell = g.MkArray([2], 0, [5, 6]);
	ell = argEllipsis(ell);
	argEllipsis(g.MkArray([2], 0, [5, 6]));

	document.write("=== RUN argSlice<br>");
	var s = g.Slice(0, ['8', '9']);
	s = argSlice(s);
	argSlice(g.Slice(0, ['8', '9']));

	document.write("=== RUN argMap<br>");
	var m = g.Map("", {1: "foo", 2: "bar"});
	m = argMap(m);
	argMap(g.Map("", {1: "foo", 2: "bar"}));

	if (PASS) {
		document.write("PASS<br>");
	} else {
		document.write("FAIL<br>");
		alert("Fail: Miscellaneous");
	}
} main();
/* Generated by GoScript (github.com/kless/goscript) */
