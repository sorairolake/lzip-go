// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip_test

import (
	"testing"

	"github.com/sorairolake/lzip-go"
)

func TestHeaderSize(t *testing.T) {
	t.Parallel()

	const expected = 6

	if size := lzip.HeaderSize; size != expected {
		t.Errorf("expected header size `%v`, got `%v`", expected, size)
	}
}

func TestTrailerSize(t *testing.T) {
	t.Parallel()

	const expected = 20

	if size := lzip.TrailerSize; size != expected {
		t.Errorf("expected trailer size `%v`, got `%v`", expected, size)
	}
}

func TestMagic(t *testing.T) {
	t.Parallel()

	const expected = "LZIP"

	if magic := lzip.Magic; magic != expected {
		t.Errorf("expected magic number `%v`, got `%v`", expected, magic)
	}
}

func TestMagicSize(t *testing.T) {
	t.Parallel()

	const expected = 4

	if size := lzip.MagicSize; size != expected {
		t.Errorf("expected magic number size `%v`, got `%v`", expected, size)
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()

	const (
		expectedV0 = iota
		expectedV1
	)

	if v0 := lzip.Version0; v0 != expectedV0 {
		t.Errorf("expected version number `%v`, got `%v`", expectedV0, v0)
	}

	if v1 := lzip.Version1; v1 != expectedV1 {
		t.Errorf("expected version number `%v`, got `%v`", expectedV1, v1)
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

	const expected = 8 * (1 << 20)

	if size := lzip.DefaultDictSize; size != expected {
		t.Errorf("expected default dictionary size `%v`, got `%v`", expected, size)
	}
}

func TestMaxMemberSize(t *testing.T) {
	t.Parallel()

	const expected = 2 * (1 << 50)

	if size := lzip.MaxMemberSize; size != expected {
		t.Errorf("expected maximum member size `%v`, got `%v`", expected, size)
	}
}
