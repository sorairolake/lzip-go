// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package main

import (
	"bufio"
	"io"
	"os"

	"github.com/sorairolake/lzip-go"
)

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
