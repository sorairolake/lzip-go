// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/sorairolake/lzip-go"
)

func Example() {
	const text = "The quick brown fox jumps over the lazy dog."

	opt := &lzip.WriterOptions{4 * 1024 * 1024}

	var buf bytes.Buffer

	writer, err := lzip.NewWriterOptions(&buf, opt)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.WriteString(writer, text); err != nil {
		log.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	reader, err := lzip.NewReader(&buf)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		log.Fatal(err)
	}

	// Output:
	// The quick brown fox jumps over the lazy dog.
}

func ExampleReader() {
	file, err := os.Open("testdata/fox.lz")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader, err := lzip.NewReader(bufio.NewReader(file))
	if err != nil {
		log.Print(err)

		return
	}

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		log.Print(err)

		return
	}

	// Output:
	// The quick brown fox jumps over the lazy dog.
}

func ExampleWriter() {
	const text = "The quick brown fox jumps over the lazy dog."

	var buf bytes.Buffer
	writer := lzip.NewWriter(&buf)

	if _, err := io.WriteString(writer, text); err != nil {
		log.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	// Output:
}
