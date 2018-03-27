// GoDepo - advanced dependencies management tool for Golang.
//
// Copyright (c) 2018, Stanislav N. aka pztrn.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject
// to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	// stdlib
	"flag"
	"log"

	// local
	"github.com/pztrn/godepo/config"
	"github.com/pztrn/godepo/parser"
)

var (
	// Flags that controls behaviour.
	// Should we ensure that dependencies are exactly the same as in
	// godepo.lock file?
	ensure = false

	// Global debug? Will output much logs!
	debug = false

	// Instances of structs.
	prsr *parser.Parser
)

func main() {
	log.Print("Starting godepo...")

	flag.BoolVar(&debug, "debug", false, "Activate debug mode. Will print many more log lines")
	flag.BoolVar(&ensure, "ensure", false, "Ensure that vendored packages are same as specified in godepo.lock file")

	cfg := config.New()
	cfg.Initialize()

	prsr = parser.New(cfg)

	flag.Parse()

	config.SetDebug(debug)
	parser.SetDebug(debug)

	cfg.ReadConfig()

	if debug {
		log.Print("Debug mode activated!")
	}

	if !ensure {
		flag.PrintDefaults()
	}

	if ensure {
		ensuredeps()
	}
}

func ensuredeps() {
	log.Print("Ensuring that dependencies are at the very save revision/version/branch as specified in godepo.yaml and godepo.lock...")
	prsr.Parse()
}
