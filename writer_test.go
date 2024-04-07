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

	opt := &lzip.WriterOptions{1 << 22}

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

	const expected = (1 << 12) - 1

	opt := &lzip.WriterOptions{expected}

	_, err := lzip.NewWriterOptions(io.Discard, opt)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var dictSizeTooSmallError *lzip.DictSizeTooSmallError
	if !errors.As(err, &dictSizeTooSmallError) {
		t.Fatal("unexpected error type")
	}

	if size := dictSizeTooSmallError.DictSize; size != expected {
		t.Errorf("expected too small dictionary size `%v`, got `%v`", expected, size)
	}
}

func TestWriterOptionsDictSizeTooLarge(t *testing.T) {
	t.Parallel()

	const expected = (1 << 29) + 1

	opt := &lzip.WriterOptions{expected}

	_, err := lzip.NewWriterOptions(io.Discard, opt)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var dictSizeTooLargeError *lzip.DictSizeTooLargeError
	if !errors.As(err, &dictSizeTooLargeError) {
		t.Fatal("unexpected error type")
	}

	if size := dictSizeTooLargeError.DictSize; size != expected {
		t.Errorf("expected too large dictionary size `%v`, got `%v`", expected, size)
	}
}

func TestVerifyWriterOptions(t *testing.T) {
	t.Parallel()

	opt := &lzip.WriterOptions{1 << 12}
	if err := opt.Verify(); err != nil {
		t.Fatal(err)
	}

	opt = &lzip.WriterOptions{1 << 29}
	if err := opt.Verify(); err != nil {
		t.Fatal(err)
	}

	var expected uint32 = (1 << 12) - 1
	opt = &lzip.WriterOptions{expected}

	err := opt.Verify()
	if err == nil {
		t.Fatal("unexpected success")
	}

	var dictSizeTooSmallError *lzip.DictSizeTooSmallError
	if !errors.As(err, &dictSizeTooSmallError) {
		t.Fatal("unexpected error type")
	}

	if size := dictSizeTooSmallError.DictSize; size != expected {
		t.Errorf("expected too small dictionary size `%v`, got `%v`", expected, size)
	}

	expected = (1 << 29) + 1
	opt = &lzip.WriterOptions{expected}

	err = opt.Verify()
	if err == nil {
		t.Fatal("unexpected success")
	}

	var dictSizeTooLargeError *lzip.DictSizeTooLargeError
	if !errors.As(err, &dictSizeTooLargeError) {
		t.Fatal("unexpected error type")
	}

	if size := dictSizeTooLargeError.DictSize; size != expected {
		t.Errorf("expected too large dictionary size `%v`, got `%v`", expected, size)
	}
}
