[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

# Z80 Processor in Go

This is an implementation of the Mostek / Zilog Z80 processor in Go

If you find this project useful, you may want to [__Buy me a
Coffee!__ :coffee:](https://www.buymeacoffee.com/codesqueak) Thanks :thumbsup:

## How to use

import (
"github.com/codesqueak/z80"
)

## Undocumented instruction

The code attempts to faithfully reproduce the numerous undocumented instructions in the Z80. I have tested against a
real device but if you find any issues, let me know.

## How to make a machine

To make a machine you need three components, the CPU, Memory and I/O. To see a simple example, look at the test in
core_instructions_test.go.  