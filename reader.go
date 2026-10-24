// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package lzip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"math"
	"slices"

	"github.com/ulikunitz/xz/lzma"
)

// countingReader counts the number of bytes read from the underlying reader.
type countingReader struct {
	r     io.Reader
	count uint64
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	cr.count += uint64(n)

	return n, err
}

// Reader is an [io.Reader] that can be read to retrieve uncompressed data from
// a lzip-format compressed file.
//
// The file may consist of multiple members, in which case the uncompressed
// data of all members is concatenated.
type Reader struct {
	r            *countingReader
	decompressor *lzma.Reader
	trailer
	err error
}

// NewReader creates a new [Reader] reading the given reader.
//
// The reader may read more data than necessary from r.
func NewReader(r io.Reader) (*Reader, error) {
	z := &Reader{r: &countingReader{r: bufio.NewReader(r)}}

	if err := z.nextMember(); err != nil {
		return nil, err
	}

	return z, nil
}

// nextMember reads the header of the next member and sets up the decompressor
// for it. It returns [io.EOF] if there are no more members.
func (z *Reader) nextMember() error {
	z.r.count = 0

	var header [headerSize]byte

	if _, err := io.ReadFull(z.r, header[:]); err != nil {
		return err
	}

	if !slices.Equal(header[:magicSize], []byte(magic)) {
		return ErrInvalidMagic
	}

	switch v := header[4]; v {
	case 0:
		return &UnsupportedVersionError{v}
	case 1:
	default:
		return &UnknownVersionError{v}
	}

	dictSize := uint32(1 << (header[5] & 0x1F))
	dictSize -= (dictSize / 16) * uint32((header[5]>>5)&0x07)

	switch {
	case dictSize < MinDictSize:
		return &DictSizeTooSmallError{dictSize}
	case dictSize > MaxDictSize:
		return &DictSizeTooLargeError{dictSize}
	}

	var lzmaHeader [lzma.HeaderLen]byte

	lzmaHeader[0] = lzma.Properties{LC: 3, LP: 0, PB: 2}.Code()
	binary.LittleEndian.PutUint32(lzmaHeader[1:5], dictSize)
	// The uncompressed size is unknown, so the end of stream marker is used.
	binary.LittleEndian.PutUint64(lzmaHeader[5:], math.MaxUint64)

	decompressor, err := lzma.NewReader(io.MultiReader(bytes.NewReader(lzmaHeader[:]), z.r))
	if err != nil {
		return err
	}

	z.decompressor = decompressor
	z.trailer = trailer{}

	return nil
}

// Read reads uncompressed data from the stream.
func (z *Reader) Read(p []byte) (int, error) {
	if z.err != nil {
		return 0, z.err
	}

	for {
		n, err := z.decompressor.Read(p)

		if n > 0 {
			z.crc = crc32.Update(z.crc, crc32.IEEETable, p[:n])
			z.dataSize += uint64(n)

			return n, nil
		}

		if errors.Is(err, io.EOF) {
			err = z.readTrailer()
			if err == nil {
				err = z.nextMember()
			}
		}

		if err != nil {
			z.err = err

			return 0, err
		}
	}
}

// readTrailer reads the trailer of the current member and verifies it.
func (z *Reader) readTrailer() error {
	var trailer [trailerSize]byte
	if _, err := io.ReadFull(z.r, trailer[:]); err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}

		return err
	}

	crc := binary.LittleEndian.Uint32(trailer[:4])
	if crc != z.crc {
		return &InvalidCRCError{crc}
	}

	dataSize := binary.LittleEndian.Uint64(trailer[4:12])
	if dataSize != z.dataSize {
		return &InvalidDataSizeError{dataSize}
	}

	memberSize := binary.LittleEndian.Uint64(trailer[12:])
	if memberSize != z.r.count {
		return &InvalidMemberSizeError{memberSize}
	}

	return nil
}
