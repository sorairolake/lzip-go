// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/sorairolake/lzip-go"
)

func uncompress(file string, output *os.File, opt options) (err error) {
	input, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func() {
		if e := input.Close(); e != nil {
			err = e
		}
	}()

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

	defer func() {
		if e := output.Close(); e != nil {
			err = e
		}
	}()

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
