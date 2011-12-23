/* Generated by GoJscript <github.com/kless/GoJscript> */

function testIf() {
	var x = 5;


	if (x > 10) {
		alert("x is greater than 10\n");
	} else {
		alert("x is less than 10\n");
	}


	var x = 12; if (x > 10) {
		alert("x is greater than 10\n");
	} else {
		alert("x is less than 10\n");
	}


	var i = 7;

	if (i === 3) {
		alert("i is equal to 3\n");
	} else if (i < 3) {
		alert("i is less than 3\n");
	} else {
		alert("i is greater than 3\n");
	}
}

function testSwitch() {
	var i = 10;


	switch (i) {
	case 1:
		alert("i is equal to 1\n"); break;
	case 2: case 3: case 4:
		alert("i is equal to 2, 3 or 4\n"); break;
	case 10:
		alert("i is equal to 10\n"); break;
	default:
		alert("All I know is that i is an integer\n");
	}


	switch (1) {
	case i < 10:
		alert("i is less than 10\n"); break;
	case i > 10: case i < 0:
		alert("i is either bigger than 10 or less than 0\n"); break;
	case i === 10:
		alert("i is equal to 10\n"); break;
	default:
		alert("This won't be printed anyway\n");
	}


	var i = 6; switch (1) {
	case 4:
		alert("was <= 4\n");
		
	case 5:
		alert("was <= 5\n");
		
	case 6:
		alert("was <= 6\n");
		
	case 7:
		alert("was <= 7\n");
		
	case 8:
		alert("was <= 8\n"); break;

	default:
		alert("default case\n");
	}


	switch (i) {
	default: break;
	case 1: case 3: case 5: case 7: case 9:
		return "odd";
	case 2: case 4: case 6: case 8:
		return "even";
	}

	return "";
}

function testFor() {
	var sum = 0;


	for (var i = 0; i < 10; i++) {
		sum += i;
	}
	alert("sum is equal to \n");


	var sum = 1;
	for (; sum < 1000;) {
		sum += sum;
	}
	alert("sum is equal to \n");


	var sum = 1;
	for (; sum < 1000;) {
		sum += sum;
	}
	alert("sum is equal to \n");


	for (;;) {
		alert("I loop for ever!\n");
	}


	for (var i = 10; i > 0; i--) {
		if (i < 5) {
			break;
		}
		alert("i\n");
	}


	for (var i = 10; i > 0; i--) {
		if (i === 5) {
			continue;
		}
		alert("i\n");
	}
}

function testRange() {
	var s = [2,3,5];




}
