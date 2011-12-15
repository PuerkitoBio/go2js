GoScript
========

Compiles transforming Go into JavaScript so you can continue using a clean and
concise sintaxis.  
In fact, it is used a subset of Go since JavaScript has not native way to
represent some types neither Go's statements, although some of them could be
emulated (but that is not my goal).

Advantages:

+ Using one only language for all development. A great advantage for a company.

+ It allows many type errors to be caught early in the development cycle, due
to static typing.

+ The mathematical expressions are calculated at the translation stage. (ToDo)

+ The lines numbers in the unminified generated JavaScript match up with the
lines numbers in the original source file.

+ Generate minimized JavaScript. (ToDo)

Go sintaxis not supported:

+ Complex numbers, and integers of 64 bits
+ Function type, interface type excepting the empty interface
+ Channels, goroutines (ToDo)
+ Import declarations from core library (ToDo)

Status:

	const				[OK]
	itoa				[OK]
	blank identifier	[OK]
	var					[OK]
	array				[OK]
	slice				[OK]
	ellipsis			[OK]
	map					[OK]
	empty interface		[OK]
	check channel		[OK]
	struct				[OK]
	pointer				[OK]

**Note:** JavaScript can not actually do meaningful integer arithmetic on anything
bigger than 2^53. Also bitwise logical operations only have defined results (per
the spec) up to 32 bits.  
By this reason, the integers of 64 bits are unsupported.


## Installation

	goinstall << DOWNLOAD URL >>


## Configuration

Nothing.


## Operating instructions

<< INSTRUCTIONS TO RUN THE PROGRAM >>


## Copyright and licensing

*Copyright 2011  The "GoScript" Authors*. See file AUTHORS and CONTRIBUTORS.  
Unless otherwise noted, the source files are distributed under the
*GNU General Public License, version 3 or later* found in the LICENSE file.


* * *
*Generated by [GoWizard](https://github.com/kless/GoWizard)*

