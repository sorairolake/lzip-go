// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package main

import (
	"flag"
	"fmt"

	"github.com/sorairolake/lzip-go"
)

const version = "0.3.5"

type options struct {
	version        bool
	stdout         bool
	decompress     bool
	keep           bool
	dictionarySize uint
}

var opt options

func init() {
	flag.BoolVar(&opt.version, "version", false, "Print version number")
	flag.BoolVar(&opt.stdout, "stdout", false, "Write to standard output, keep input files")
	flag.BoolVar(&opt.decompress, "decompress", false, "Decompress data")
	flag.BoolVar(&opt.keep, "keep", false, "Keep input files")
	flag.UintVar(&opt.dictionarySize, "dictionary-size", lzip.DefaultDictSize, "Set dictionary size in bytes")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: glzip [OPTIONS] <FILE>...\n")
		flag.PrintDefaults()
	}
}
