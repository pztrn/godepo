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

// Package config responsible for all works with configuration, like
// reading configuration from disk into memory or write configuration
// from memory to disk.
package config

import (
	// stdlib
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	// local
	"github.com/pztrn/godepo/config/struct"

	// other
	"gopkg.in/yaml.v2"
)

type Config struct {
	// Configuration file path.
	configFilePath string
	// Parsed configuration.
	Config configstruct.ConfigStruct
	// Project path.
	ProjectPath string
}

// Initialize initializes configuration storage. It will also add
// neccessary CLI flags.
func (c *Config) Initialize() {
	flag.StringVar(&c.ProjectPath, "projectpath", ".", "Path to project where godepo.yaml file is resided.")
}

// ReadConfig reads configuration into memory.
func (c *Config) ReadConfig() {
	// Before actual configuration reading we should expand "~" if
	// present.
	if strings.Contains(c.ProjectPath, "~") {
		c.resolveHomeDir()
	}

	// Compose configuration file path.
	c.configFilePath = filepath.Join(c.ProjectPath, "godepo.yaml")

	// Check if file exists at all.
	c.checkFileExisting()

	// Read configuration.
	log.Print("Reading configuration from '" + c.configFilePath + "'...")
	c.readFile()
}

// Check if configuration file exists.
func (c *Config) checkFileExisting() {
	if _, err := os.Stat(c.ProjectPath); os.IsNotExist(err) {
		log.Fatalf("Failed to load configuration file: %s", err.Error())
	}
}

func (c *Config) readFile() {
	filedata, err := ioutil.ReadFile(c.configFilePath)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err.Error())
	}
	c.Config = configstruct.ConfigStruct{}
	err1 := yaml.Unmarshal(filedata, &c.Config)
	if err1 != nil {
		log.Fatalf("Failed to parse configuration file from YAML into struct: %s", err1.Error())
	}
}

// Resolves home directory. ATM only NIX systems are supported.
func (c *Config) resolveHomeDir() {
	// ToDo: support Windows.
	if runtime.GOOS == "windows" {
		log.Fatal("GoDepo currently isn't working on Windows. Feel free to submit patches at https://github.com/pztrn/godepo")
	}

	curUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to obtain current user data: %s", err.Error())
	}
	if curUser.HomeDir == "" {
		log.Fatal("You've used tilde ('~') in configuration file path, but current user have no home directory defined. Cannot continue.")
	}

	c.ProjectPath = strings.Replace(c.ProjectPath, "~", curUser.HomeDir, 1)
}
