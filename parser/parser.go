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

package parser

import (
	// stdlib
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	// local
	"github.com/pztrn/godepo/config/struct"
)

// Parser responsible for all parsing actions.
type Parser struct {
	files []string
}

// Handler for filepath.Walk call.
// If there is something that we have to ALWAYS ignore - it should be
// here.
func (p *Parser) filepathWalkHandler(path string, info os.FileInfo, err error) error {
	// If something should be skipped - it should be right there.
	// Skip vendor directory and all it's contents.
	if strings.Contains(path, "/vendor/") {
		return nil
	}

	// Check that we have file here, not a directory.
	// ToDo: should we also work with symlinks, or such files
	// should be ignored?
	if info.IsDir() {
		return nil
	}

	// Do not append files that didn't end on ".go".
	if !strings.HasSuffix(path, ".go") {
		return nil
	}

	p.files = append(p.files, path)

	return nil
}

// Parse parses dependencies, gets only unique of them, reformats-preformats
// and returns a slice with unique and ready-to-install dependencies.
// Example return:
//
//     []configstruct.ConfigPackage{
//	       configstruct.ConfigPackage{
//	           ImportPath: "github.com/pztrn/flagger",
//             SourcePath: "https://github.com/pztrn/flagger",
//             VCSName: "git"
//         }
//     }
func (p *Parser) Parse() []*configstruct.ConfigPackage {
	//var packagesData []*configstruct.ConfigPackage

	// Get files list in project directory.
	p.files = []string{}
	filepath.Walk(cfg.ProjectPath, p.filepathWalkHandler)
	if debug {
		log.Printf("[DEBUG] Got %d files", len(p.files))
	}

	_ = p.parseFilesForPackages()

	return nil
}

// Parses passed files, composes a list of packages and returns.
func (p *Parser) parseFilesForPackages() []string {
	if debug {
		log.Print("Parsing project files for unique packages...")
	}

	var packages []string

	for i := range p.files {
		if debug {
			log.Print("[DEBUG] Reading file: '" + p.files[i] + "'...")
		}

		fileDataAsBytes, err := ioutil.ReadFile(p.files[i])
		if err != nil {
			log.Fatalf("Failed to read file '"+p.files[i]+"': %s", err.Error())
		}

		// Get string and a slice of lines.
		fileData := strings.Split(string(fileDataAsBytes), "\n")
		if debug {
			log.Printf("[DEBUG]\tFile contains %d lines", len(fileData))
		}

		var importsStart = false
		var importsEnd = false
		var multilineComment = false

		for i := range fileData {
			var packageName = ""
			// Skip multiline comments.
			if strings.Contains(fileData[i], "/*") {
				multilineComment = true
				continue
			}

			if strings.Contains(fileData[i], "*/") && multilineComment {
				multilineComment = false
				continue
			}

			if !strings.Contains(fileData[i], "*/") && multilineComment {
				continue
			}

			// Working with "import()".
			if strings.Contains(fileData[i], "import (") {
				importsStart = true
				continue
			}

			if importsStart {
				// Stop searching for packages when imports was parsed.
				if strings.Contains(fileData[i], ")") {
					if debug {
						log.Print("[DEBUG]\tImports ended.")
						importsStart = false
						importsEnd = true
						continue
					}
				}

				// Skip comments.
				if strings.Contains(fileData[i], "//") {
					continue
				}

				// Skip lines that didn't start with ".
				if !strings.Contains(fileData[i], "\"") {
					continue
				}

				packageName = strings.Trim(fileData[i], " ")
			}

			// Working with one-line imports.
			if strings.Contains(fileData[i], "import \"") {
				packageName = strings.Split(fileData[i], " ")[1]
			}

			// Do some magic.
			if packageName != "" {
				// Clearing from spaces.
				packageName = strings.Trim(fileData[i], " ")
				// Clearing from tabs.
				packageName = strings.Trim(packageName, "\t")
				// Clearing from doublequotes.
				packageName = strings.Trim(packageName, "\"")
				// Clearing from "_"s.
				packageName = strings.TrimLeft(packageName, "_ \"")

				// If we have no dots in first part of package name,
				// then we consider this package as builtin and skip.
				packageDomain := strings.Split(packageName, "/")[0]
				if !strings.Contains(packageDomain, ".") {
					continue
				}

				if debug {
					log.Printf("[DEBUG] \t\tFound import: %s", packageName)
				}

				// Check for uniquiness.
				var packageUnique = true
				for ii := range packages {
					if packages[ii] == packageName {
						packageUnique = false
					}
				}
				if packageUnique {
					packages = append(packages, packageName)
				}
			}

			if importsEnd {
				break
			}
		}
	}

	if debug {
		log.Printf("[DEBUG] Found %d unique packages", len(packages))
		log.Print(packages)
	}

	return packages
}
