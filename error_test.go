// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"math"
	"testing"

	"github.com/sorairolake/lzip-go"
)

func TestErrInvalidMagic(t *testing.T) {
	t.Parallel()

	err := lzip.ErrInvalidMagic
	expected := "lzip: invalid magic number"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestUnsupportedVersionError(t *testing.T) {
	t.Parallel()

	err := lzip.UnsupportedVersionError{0}
	expected := "lzip: unsupported version number"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if v := err.Version; v != 0 {
		t.Errorf("expected unsupported version number `%v`, got `%v`", 0, v)
	}
}

func TestUnknownVersionError(t *testing.T) {
	t.Parallel()

	err := lzip.UnknownVersionError{math.MaxUint8}
	expected := "lzip: unknown version number"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if v := err.Version; v != math.MaxUint8 {
		t.Errorf("expected unknown version number `%v`, got `%v`", math.MaxUint8, v)
	}
}

func TestDictSizeTooSmallError(t *testing.T) {
	t.Parallel()

	err := lzip.DictSizeTooSmallError{lzip.MinDictSize - 1}
	expected := "lzip: dictionary size is too small"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if size := err.DictSize; size != (lzip.MinDictSize - 1) {
		t.Errorf("expected too small dictionary size `%v`, got `%v`", lzip.MinDictSize-1, size)
	}
}

func TestDictSizeTooLargeError(t *testing.T) {
	t.Parallel()

	err := lzip.DictSizeTooLargeError{lzip.MaxDictSize + 1}
	expected := "lzip: dictionary size is too large"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if size := err.DictSize; size != (lzip.MaxDictSize + 1) {
		t.Errorf("expected too large dictionary size `%v`, got `%v`", lzip.MaxDictSize+1, size)
	}
}

func TestInvalidCRCError(t *testing.T) {
	t.Parallel()

	err := lzip.InvalidCRCError{0}
	expected := "lzip: CRC mismatch"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if crc := err.CRC; crc != 0 {
		t.Errorf("expected invalid CRC `%v`, got `%v`", 0, crc)
	}
}

func TestInvalidDataSizeError(t *testing.T) {
	t.Parallel()

	err := lzip.InvalidDataSizeError{0}
	expected := "lzip: data size mismatch"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if size := err.DataSize; size != 0 {
		t.Errorf("expected invalid data size `%v`, got `%v`", 0, size)
	}
}

func TestInvalidMemberSizeError(t *testing.T) {
	t.Parallel()

	err := lzip.InvalidMemberSizeError{0}
	expected := "lzip: member size mismatch"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if size := err.MemberSize; size != 0 {
		t.Errorf("expected invalid member size `%v`, got `%v`", 0, size)
	}
}

func TestMemberSizeTooLargeError(t *testing.T) {
	t.Parallel()

	err := lzip.MemberSizeTooLargeError{lzip.MaxMemberSize + 1}
	expected := "lzip: member size is too large"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}

	if size := err.MemberSize; size != (lzip.MaxMemberSize + 1) {
		t.Errorf("expected too large member size `%v`, got `%v`", lzip.MaxMemberSize+1, size)
	}
}
