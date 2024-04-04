// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip

import "errors"

var (
	// ErrInvalidMagic represents an error due to the magic number was
	// invalid.
	ErrInvalidMagic = errors.New("lzip: invalid magic number")

	// ErrUnsupportedVersion represents an error due to the version number
	// stored in the header indicated the lzip format which is not
	// supported by this package.
	ErrUnsupportedVersion = errors.New("lzip: unsupported version number")

	// ErrUnknownVersion represents an error due to the version number
	// stored in the header was not recognized by this package.
	ErrUnknownVersion = errors.New("lzip: unknown version number")

	// ErrDictSizeTooSmall represents an error due to the dictionary size
	// was smaller than 4 KiB.
	ErrDictSizeTooSmall = errors.New("lzip: dictionary size is too small")

	// ErrDictSizeTooLarge represents an error due to the dictionary size
	// was larger than 512 MiB.
	ErrDictSizeTooLarge = errors.New("lzip: dictionary size is too large")

	// ErrInvalidCRC represents an error due to a CRC of the original
	// uncompressed data mismatched.
	ErrInvalidCRC = errors.New("lzip: CRC mismatch")

	// ErrInvalidDataSize represents an error due to the size of the
	// original uncompressed data stored in the trailer and the actual size
	// of it mismatched.
	ErrInvalidDataSize = errors.New("lzip: data size mismatch")

	// ErrInvalidMemberSize represents an error due to the total size of
	// the member stored in the trailer and the actual total size of it
	// mismatched.
	ErrInvalidMemberSize = errors.New("lzip: member size mismatch")

	// ErrDictSizeTooLarge represents an error due to the member size was
	// larger than 2 PiB.
	ErrMemberSizeTooLarge = errors.New("lzip: member size is too large")
)
