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

func TestReader(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/test.txt.lz")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := lzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), data) {
		t.Error("unexpected mismatch between uncompressed data and test data")
	}
}

func TestReaderInvalidMagic(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/fox_bm.lz")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := lzip.NewReader(file); !errors.Is(err, lzip.ErrInvalidMagic) {
		t.Error("unexpected error type")
	}
}

func TestReaderUnknownVersion(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/fox_v2.lz")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lzip.NewReader(file)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var unknownVersionError *lzip.UnknownVersionError
	if !errors.As(err, &unknownVersionError) {
		t.Fatal("unexpected error type")
	}

	const expected = 2

	if v := unknownVersionError.Version; v != expected {
		t.Errorf("expected unrecognized version number `%v`, got `%v`", expected, v)
	}
}

func TestReaderDictSizeTooSmall(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/fox_s11.lz")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lzip.NewReader(file)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var dictSizeTooSmallError *lzip.DictSizeTooSmallError
	if !errors.As(err, &dictSizeTooSmallError) {
		t.Fatal("unexpected error type")
	}

	const expected = 1 << 11

	if size := dictSizeTooSmallError.DictSize; size != expected {
		t.Errorf("expected too small dictionary size `%v`, got `%v`", expected, size)
	}
}
