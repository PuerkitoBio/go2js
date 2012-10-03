/* Generated by GoScript (github.com/kless/goscript) */









var PASS = true;


function person(name, age) {
	this.name=name;
	this.age=age
}



function getOlder(people) {
	if (people.length === 0) {
		return [new person(), false];
	}

	var older = people[0];

	var value; for (var _ in people) { value = people[_];
		if (value.age > older.age) {
			older = value;
		}
	}
	return [older, true];
}

function main() {
	var pass = true;

	document.write("<br><br>== More functions<br>");

	
	var ok = false;
	var older = new person("", 0);



	var paul = new person("Paul", 23);
	var jim = new person("Jim", 24);
	var sam = new person("Sam", 84);
	var rob = new person("Rob", 54);
	var karl = new person("Karl", 19);

	var tests = g.Slice("", ["Jim", "Sam", "Sam", "Karl"]);

	older = getOlder(paul, jim)[0];
	if (JSON.stringify(older.name) !== JSON.stringify(tests.f[0])) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: (getOlder paul,jim) => got " + older.name + ", want " + tests.f[0] + "<br>");

		pass = false, PASS = false;
	}

	older = getOlder(paul, jim, sam)[0];
	if (JSON.stringify(older.name) !== JSON.stringify(tests.f[1])) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: (getOlder paul,jim,sam) => got " + older.name + ", want " + tests.f[1] + "<br>");

		pass = false, PASS = false;
	}

	older = getOlder(paul, jim, sam, rob)[0];
	if (JSON.stringify(older.name) !== JSON.stringify(tests.f[2])) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: (getOlder paul,jim,sam,rob) => got " + older.name + ", want " + tests.f[2] + "<br>");

		pass = false, PASS = false;
	}

	older = getOlder(karl)[0];
	if (JSON.stringify(older.name) !== JSON.stringify(tests.f[3])) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: (getOlder karl) => got " + older.name + ", want " + tests.f[3] + "<br>");

		pass = false, PASS = false;
	}


	_ = getOlder(), older = _[0], ok = _[1];
	if (ok) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;FAIL: (getOlder) => got " + ok + ", want " + !ok + "<br>");
		pass = false, PASS = false;
	}

	if (pass) {
		document.write("&nbsp;&nbsp;&nbsp;&nbsp;pass<br>");
	}

	if (PASS) {
		document.write("<br>PASS<br>");
	} else {
		document.write("<br>FAIL<br>");
		alert("Fail: More functions");
	}
} main();
