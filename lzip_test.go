// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"slices"
	"testing"

	"github.com/sorairolake/lzip-go"
)

func TestHeaderSize(t *testing.T) {
	t.Parallel()

	if size := lzip.HeaderSize; size != 6 {
		t.Errorf("expected header size `%v`, got `%v`", 6, size)
	}
}

func TestTrailerSize(t *testing.T) {
	t.Parallel()

	if size := lzip.TrailerSize; size != 20 {
		t.Errorf("expected trailer size `%v`, got `%v`", 20, size)
	}
}

func TestMagic(t *testing.T) {
	t.Parallel()

	expected := [lzip.MagicSize]byte{0x4c, 0x5a, 0x49, 0x50}
	if !slices.Equal([]byte(lzip.Magic), expected[:]) {
		t.Error("unexpected magic number")
	}
}

func TestMagicSize(t *testing.T) {
	t.Parallel()

	if size := lzip.MagicSize; size != 4 {
		t.Errorf("expected magic number size `%v`, got `%v`", 4, size)
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()

	if v0 := lzip.Version0; v0 != 0 {
		t.Errorf("expected version number `%v`, got `%v`", 0, v0)
	}

	if v1 := lzip.Version1; v1 != 1 {
		t.Errorf("expected version number `%v`, got `%v`", 1, v1)
	}
}

func TestMinDictSize(t *testing.T) {
	t.Parallel()

	const expected = 1 << 12

	if size := lzip.MinDictSize; size != expected {
		t.Errorf("expected minimum dictionary size `%v`, got `%v`", expected, size)
	}
}

func TestMaxDictSize(t *testing.T) {
	t.Parallel()

	const expected = 1 << 29

	if size := lzip.MaxDictSize; size != expected {
		t.Errorf("expected maximum dictionary size `%v`, got `%v`", expected, size)
	}
}

func TestDefaultDictSize(t *testing.T) {
	t.Parallel()

	const expected = 1 << 23

	if size := lzip.DefaultDictSize; size != expected {
		t.Errorf("expected default dictionary size `%v`, got `%v`", expected, size)
	}
}

func TestMaxMemberSize(t *testing.T) {
	t.Parallel()

	const expected = 1 << 51

	if size := lzip.MaxMemberSize; size != expected {
		t.Errorf("expected maximum member size `%v`, got `%v`", expected, size)
	}
}
