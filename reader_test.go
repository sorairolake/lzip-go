// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"slices"
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

func TestReaderFromEmpty(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/em.lz")
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

	data := []byte{}

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

func TestReaderNonZeroFirstByte(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/fox_nz.lz")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := lzip.NewReader(file); err == nil {
		t.Error("unexpected success")
	}
}

func TestReadInvalidCRC(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/fox_bcrc.lz")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := lzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, reader)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var invalidCRCError *lzip.InvalidCRCError
	if !errors.As(err, &invalidCRCError) {
		t.Fatal("unexpected error type")
	}

	const expected = 0xEB50_CC6B

	if crc := invalidCRCError.CRC; crc != expected {
		t.Errorf("expected CRC `%v`, got `%v`", expected, crc)
	}
}

func TestReadInvalidMemberSize(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/fox_mes81.lz")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := lzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, reader)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var invalidMemberSizeError *lzip.InvalidMemberSizeError
	if !errors.As(err, &invalidMemberSizeError) {
		t.Fatal("unexpected error type")
	}

	const expected = 81

	if memberSize := invalidMemberSizeError.MemberSize; memberSize != expected {
		t.Errorf("expected member size `%v`, got `%v`", expected, memberSize)
	}
}

func concatFiles(t *testing.T, names ...string) []byte {
	t.Helper()

	var buf bytes.Buffer

	for _, name := range names {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		buf.Write(data)
	}

	return buf.Bytes()
}

func TestReaderMultiMember(t *testing.T) {
	t.Parallel()

	compressed := concatFiles(t, "testdata/test.txt.lz", "testdata/em.lz", "testdata/test.txt.lz")

	reader, err := lzip.NewReader(bytes.NewReader(compressed))
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

	if !bytes.Equal(buf.Bytes(), slices.Concat(data, data)) {
		t.Error("unexpected mismatch between uncompressed data and test data")
	}
}

func TestReadMultiMemberInvalidCRC(t *testing.T) {
	t.Parallel()

	compressed := concatFiles(t, "testdata/test.txt.lz", "testdata/fox_bcrc.lz")

	reader, err := lzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, reader)
	if err == nil {
		t.Fatal("unexpected success")
	}

	var invalidCRCError *lzip.InvalidCRCError
	if !errors.As(err, &invalidCRCError) {
		t.Fatalf("unexpected error type: %v", err)
	}

	const expected = 0xEB50_CC6B

	if crc := invalidCRCError.CRC; crc != expected {
		t.Errorf("expected CRC `%v`, got `%v`", expected, crc)
	}
}

func TestReadTrailingData(t *testing.T) {
	t.Parallel()

	compressed := slices.Concat(concatFiles(t, "testdata/test.txt.lz"), []byte("trailing data"))

	reader, err := lzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); !errors.Is(err, lzip.ErrInvalidMagic) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReadMissingTrailer(t *testing.T) {
	t.Parallel()

	compressed := concatFiles(t, "testdata/test.txt.lz")
	compressed = compressed[:len(compressed)-lzip.TrailerSize]

	reader, err := lzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReadAfterEOF(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/test.txt.lz")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := lzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Fatal(err)
	}

	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, io.EOF) {
		t.Errorf("unexpected error: %v", err)
	}
}
