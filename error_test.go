// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
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

func TestErrUnsupportedVersion(t *testing.T) {
	t.Parallel()

	err := lzip.ErrUnsupportedVersion
	expected := "lzip: unsupported version number"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrUnknownVersion(t *testing.T) {
	t.Parallel()

	err := lzip.ErrUnknownVersion
	expected := "lzip: unknown version number"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrDictSizeTooSmall(t *testing.T) {
	t.Parallel()

	err := lzip.ErrDictSizeTooSmall
	expected := "lzip: dictionary size is too small"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrDictSizeTooLarge(t *testing.T) {
	t.Parallel()

	err := lzip.ErrDictSizeTooLarge
	expected := "lzip: dictionary size is too large"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrInvalidCRC(t *testing.T) {
	t.Parallel()

	err := lzip.ErrInvalidCRC
	expected := "lzip: CRC mismatch"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrInvalidDataSize(t *testing.T) {
	t.Parallel()

	err := lzip.ErrInvalidDataSize
	expected := "lzip: data size mismatch"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrInvalidMemberSize(t *testing.T) {
	t.Parallel()

	err := lzip.ErrInvalidMemberSize
	expected := "lzip: member size mismatch"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}

func TestErrMemberSizeTooLarge(t *testing.T) {
	t.Parallel()

	err := lzip.ErrMemberSizeTooLarge
	expected := "lzip: member size is too large"

	if err.Error() != expected {
		t.Error("unexpected error message")
	}
}
