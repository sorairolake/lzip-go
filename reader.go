// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"reflect"
	"slices"

	"github.com/ulikunitz/xz/lzma"
)

// Reader is an [io.Reader] that can be read to retrieve uncompressed data from
// a lzip-format compressed file.
type Reader struct {
	r            io.Reader
	decompressor *lzma.Reader
	trailer
}

// NewReader creates a new [Reader] reading the given reader.
func NewReader(r io.Reader) (*Reader, error) {
	z := new(Reader)

	var header [headerSize]byte
	if _, err := r.Read(header[:]); err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(header[:magicSize], []byte(magic)) {
		return nil, ErrInvalidMagic
	}

	switch v := header[4]; v {
	case 0:
		return nil, ErrUnsupportedVersion
	case 1:
	default:
		return nil, ErrUnknownVersion
	}

	dictSize := 1 << (int(header[5] & 0x1f))
	dictSize -= (dictSize / 16) * int((header[5]>>5)&0x07)

	switch {
	case dictSize < minDictSize:
		return nil, ErrDictSizeTooSmall
	case dictSize > maxDictSize:
		return nil, ErrDictSizeTooLarge
	}

	rb, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var lzmaHeader [lzma.HeaderLen]byte
	lzmaHeader[0] = lzma.Properties{LC: 3, LP: 0, PB: 2}.Code()
	binary.LittleEndian.PutUint32(lzmaHeader[1:5], uint32(dictSize))
	copy(lzmaHeader[5:], rb[len(rb)-16:len(rb)-8])

	z.trailer.memberSize = uint64(headerSize + len(rb))
	if z.trailer.memberSize > maxMemberSize {
		return nil, ErrMemberSizeTooLarge
	}

	rb = slices.Concat(lzmaHeader[:], rb)

	r = bytes.NewReader(rb)

	z.decompressor, err = lzma.NewReader(r)
	if err != nil {
		return nil, err
	}

	z.r = r

	return z, nil
}

// Read reads uncompressed data from the stream.
func (z *Reader) Read(p []byte) (n int, err error) {
	for n == 0 {
		n, err = z.decompressor.Read(p)
		if err != nil {
			return n, err
		}

		z.trailer.crc = crc32.Update(z.trailer.crc, crc32.IEEETable, p[:n])
		z.trailer.dataSize += uint64(n)

		if !errors.Is(err, io.EOF) {
			return n, err
		}

		var trailer [trailerSize]byte
		if _, err := io.ReadFull(z.r, trailer[:]); err != nil {
			return n, err
		}

		crc := binary.LittleEndian.Uint32(trailer[:4])
		if crc != z.trailer.crc {
			return n, ErrInvalidCRC
		}

		dataSize := binary.LittleEndian.Uint64(trailer[4:12])
		if dataSize != z.trailer.dataSize {
			return n, ErrInvalidDataSize
		}

		memberSize := binary.LittleEndian.Uint64(trailer[12:])
		if memberSize != z.trailer.memberSize {
			return n, ErrInvalidMemberSize
		}
	}

	return n, nil
}
