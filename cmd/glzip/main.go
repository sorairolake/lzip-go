// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sorairolake/lzip-go"
)

const version = "0.2.0"

type options struct {
	version        bool
	stdout         bool
	decompress     bool
	keep           bool
	dictionarySize uint
}

func compress(file string, output *os.File, opt options) error {
	if !opt.stdout {
		out, err := os.Create(file + ".lz")
		if err != nil {
			return err
		}

		output = out
	}

	defer output.Close()

	writerOpt := &lzip.WriterOptions{DictSize: uint32(opt.dictionarySize)}

	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	writer, err := lzip.NewWriterOptions(bufWriter, writerOpt)
	if err != nil {
		return err
	}
	defer writer.Close()

	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer input.Close()

	if _, err := io.Copy(writer, input); err != nil {
		return err
	}

	if !opt.keep {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func uncompress(file string, output *os.File, opt options) error {
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer input.Close()

	bufReader := bufio.NewReader(input)

	reader, err := lzip.NewReader(bufReader)
	if err != nil {
		return err
	}

	if !opt.stdout {
		out, err := os.Create(strings.TrimSuffix(file, ".lz"))
		if err != nil {
			return err
		}

		output = out
	}

	defer output.Close()

	if _, err := io.Copy(output, reader); err != nil {
		return err
	}

	if !opt.keep {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	opt := new(options)

	flag.BoolVar(&opt.version, "version", false, "Print version number")
	flag.BoolVar(&opt.stdout, "stdout", false, "Write to standard output, keep input files")
	flag.BoolVar(&opt.decompress, "decompress", false, "Decompress data")
	flag.BoolVar(&opt.keep, "keep", false, "Keep input files")
	flag.UintVar(&opt.dictionarySize, "dictionary-size", lzip.DefaultDictSize, "Set dictionary size in bytes")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: glzip [OPTIONS] <FILE>...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	if opt.version {
		fmt.Printf("glzip %v\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	output := os.Stdout

	if opt.stdout {
		opt.keep = true
	}

	if !opt.decompress {
		for _, file := range args {
			if err := compress(file, output, *opt); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		for _, file := range args {
			if err := uncompress(file, output, *opt); err != nil {
				log.Fatal(err)
			}
		}
	}
}
