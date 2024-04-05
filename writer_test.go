// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/sorairolake/lzip-go"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}

	text := string(data)

	var buf bytes.Buffer
	writer := lzip.NewWriter(&buf)

	n, err := io.WriteString(writer, text)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(text) {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader, err := lzip.NewReader(&buf)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if _, err := io.Copy(&out, reader); err != nil {
		t.Fatal(err)
	}

	if out.String() != text {
		t.Error("unexpected mismatch between uncompressed data and test data")
	}
}

func TestWriterOptions(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}

	text := string(data)

	opt := &lzip.WriterOptions{4 * 1024 * 1024}

	var buf bytes.Buffer

	writer, err := lzip.NewWriterOptions(&buf, opt)
	if err != nil {
		t.Fatal(err)
	}

	n, err := io.WriteString(writer, text)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(text) {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader, err := lzip.NewReader(&buf)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if _, err := io.Copy(&out, reader); err != nil {
		t.Fatal(err)
	}

	if out.String() != text {
		t.Error("unexpected mismatch between uncompressed data and test data")
	}
}

func TestWriterOptionsDictSizeTooSmall(t *testing.T) {
	t.Parallel()

	opt := &lzip.WriterOptions{(4 * 1024) - 1}

	if _, err := lzip.NewWriterOptions(io.Discard, opt); !errors.Is(err, lzip.ErrDictSizeTooSmall) {
		t.Error("unexpected error type")
	}
}

func TestWriterOptionsDictSizeTooLarge(t *testing.T) {
	t.Parallel()

	opt := &lzip.WriterOptions{(512 * 1024 * 1024) + 1}

	if _, err := lzip.NewWriterOptions(io.Discard, opt); !errors.Is(err, lzip.ErrDictSizeTooLarge) {
		t.Error("unexpected error type")
	}
}

func TestVerifyWriterOptions(t *testing.T) {
	t.Parallel()

	opt := &lzip.WriterOptions{4 * 1024}
	if err := opt.Verify(); err != nil {
		t.Fatal(err)
	}

	opt = &lzip.WriterOptions{512 * 1024 * 1024}
	if err := opt.Verify(); err != nil {
		t.Fatal(err)
	}

	opt = &lzip.WriterOptions{(4 * 1024) - 1}
	if err := opt.Verify(); !errors.Is(err, lzip.ErrDictSizeTooSmall) {
		t.Error("unexpected error type")
	}

	opt = &lzip.WriterOptions{(512 * 1024 * 1024) + 1}
	if err := opt.Verify(); !errors.Is(err, lzip.ErrDictSizeTooLarge) {
		t.Error("unexpected error type")
	}
}
