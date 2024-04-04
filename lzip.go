// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

// Package lzip implements reading and writing of lzip format compressed files.
// The package supports version 1 of the specification.
//
// See the following for the specification:
//
//   - https://www.nongnu.org/lzip/manual/lzip_manual.html#File-format
//   - https://datatracker.ietf.org/doc/html/draft-diaz-lzip-09#section-2
package lzip

import (
	"math/bits"

	"github.com/ulikunitz/xz/lzma"
)

const (
	magicSize   = 4
	headerSize  = 6
	trailerSize = 20
)

const magic = "LZIP"

type version uint8

const (
	version0 version = iota
	version1
)

const (
	minDictSize     = lzma.MinDictCap
	maxDictSize     = 1 << 29
	defaultDictSize = 0x800000
)

const maxMemberSize = 0x8000000000000

type header struct {
	magic [magicSize]byte
	version
	dictSize uint8
}

func newHeader(dictSize int) *header {
	ds := bits.Len(uint(dictSize - 1))

	if dictSize > minDictSize {
		base := 1 << dictSize
		frac := base / 16

		for i := 7; i >= 1; i-- {
			if (base - (i * frac)) >= ds {
				ds |= i << 5
			}
		}
	}

	z := &header{[magicSize]byte([]byte(magic)), version1, uint8(ds)}

	return z
}

type trailer struct {
	crc        uint32
	dataSize   uint64
	memberSize uint64
}
