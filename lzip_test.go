// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"math"
	"testing"

	"github.com/sorairolake/lzip-go"
)

func TestMagicSize(t *testing.T) {
	t.Parallel()

	if size := lzip.MagicSize; size != 4 {
		t.Errorf("expected magic number size `%v`, got `%v`", 4, size)
	}
}

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

	if magic := lzip.Magic; magic != "LZIP" {
		t.Errorf("expected magic number `%v`, got `%v`", "LZIP", magic)
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

	if size := lzip.MinDictSize; size != (1 << 12) {
		t.Errorf("expected minimum dictionary size `%v`, got `%v`", (1 << 12), size)
	}
}

func TestMaxDictSize(t *testing.T) {
	t.Parallel()

	if size := lzip.MaxDictSize; size != (1 << 29) {
		t.Errorf("expected maximum dictionary size `%v`, got `%v`", (1 << 29), size)
	}
}

func TestDefaultDictSize(t *testing.T) {
	t.Parallel()

	if size := uint32(lzip.DefaultDictSize); size != (8 * uint32(math.Pow(2.0, 20.0))) {
		t.Errorf("expected default dictionary size `%v`, got `%v`", (8 * uint32(math.Pow(2.0, 20.0))), size)
	}
}

func TestMaxMemberSize(t *testing.T) {
	t.Parallel()

	if size := uint64(lzip.MaxMemberSize); size != (2 * uint64(math.Pow(2.0, 50.0))) {
		t.Errorf("expected maximum member size `%v`, got `%v`", (2 * uint64(math.Pow(2.0, 50.0))), size)
	}
}
