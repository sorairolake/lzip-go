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
	"slices"

	"github.com/ulikunitz/xz/lzma"
)

// Reader is an [io.Reader] that can be read to retrieve uncompressed data from
// a lzip-format compressed file. A file may contain multiple members, which are
// decoded sequentially.
type Reader struct {
	r            io.Reader
	decompressor *lzma.Reader
	trailer
	members [][]byte
	current int
}

// NewReader creates a new [Reader] reading the given reader.
func NewReader(r io.Reader) (*Reader, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, io.EOF
	}
	if len(data) < headerSize {
		return nil, ErrInvalidMagic
	}

	z := new(Reader)
	z.members = splitMembers(data)
	if err := z.nextMember(); err != nil {
		return nil, err
	}

	return z, nil
}

// splitMembers locates the members in data and returns their bytes. Member
// boundaries are found by chaining the member size stored in each trailer,
// starting from the tail of data. A chain that does not close at the beginning
// of data means data is a single possibly corrupt member, which is reported as
// one member spanning all of data so that Read can surface the corruption.
func splitMembers(data []byte) [][]byte {
	var ends []int
	end := len(data)
	for end > 0 {
		if end < trailerSize {
			break
		}

		trailer := data[end-trailerSize : end]
		memberSize := binary.LittleEndian.Uint64(trailer[12:trailerSize])
		start := end - int(memberSize)
		if start < 0 || start+headerSize > end {
			break
		}
		if !slices.Equal(data[start:start+magicSize], []byte(magic)) {
			break
		}

		ends = append(ends, end)
		end = start
	}
	if end != 0 {
		return [][]byte{data}
	}

	members := make([][]byte, len(ends))
	for i, end := range ends {
		start := 0
		if i < len(ends)-1 {
			start = ends[i+1]
		}
		members[len(ends)-1-i] = data[start:end]
	}
	return members
}

// nextMember sets up the decoder for the next member. It returns [io.EOF] when
// all members have been decoded.
func (z *Reader) nextMember() error {
	if z.current >= len(z.members) {
		return io.EOF
	}
	member := z.members[z.current]
	z.current++

	dictSize, err := parseHeader(member)
	if err != nil {
		return err
	}

	if len(member) < headerSize+trailerSize {
		return &InvalidMemberSizeError{uint64(len(member))}
	}

	trailer := member[len(member)-trailerSize:]
	var lzmaHeader [lzma.HeaderLen]byte

	lzmaHeader[0] = lzma.Properties{LC: 3, LP: 0, PB: 2}.Code()
	binary.LittleEndian.PutUint32(lzmaHeader[1:5], dictSize)
	copy(lzmaHeader[5:], trailer[4:12])

	z.memberSize = uint64(len(member))
	if memberSize := z.memberSize; memberSize > MaxMemberSize {
		return &MemberSizeTooLargeError{memberSize}
	}

	body := slices.Concat(lzmaHeader[:], member[headerSize:])
	z.r = bytes.NewReader(body)

	if z.decompressor, err = lzma.NewReader(z.r); err != nil {
		return err
	}

	z.crc = 0
	z.dataSize = 0
	return nil
}

// parseHeader validates the header of a member and returns its dictionary
// size.
func parseHeader(member []byte) (uint32, error) {
	if !slices.Equal(member[:magicSize], []byte(magic)) {
		return 0, ErrInvalidMagic
	}

	switch v := member[4]; v {
	case 0:
		return 0, &UnsupportedVersionError{v}
	case 1:
	default:
		return 0, &UnknownVersionError{v}
	}

	dictSize := uint32(1 << (member[5] & 0x1F))
	dictSize -= (dictSize / 16) * uint32((member[5]>>5)&0x07)

	switch {
	case dictSize < MinDictSize:
		return 0, &DictSizeTooSmallError{dictSize}
	case dictSize > MaxDictSize:
		return 0, &DictSizeTooLargeError{dictSize}
	}

	return dictSize, nil
}

// Read reads uncompressed data from the stream.
func (z *Reader) Read(p []byte) (n int, err error) {
	for {
		n, err = z.decompressor.Read(p)

		if n > 0 {
			z.crc = crc32.Update(z.crc, crc32.IEEETable, p[:n])
			z.dataSize += uint64(n)

			return n, nil
		}

		if err == nil {
			continue
		}

		if errors.Is(err, io.EOF) {
			var trailer [trailerSize]byte
			if _, err := io.ReadFull(z.r, trailer[:]); err != nil {
				return n, err
			}

			crc := binary.LittleEndian.Uint32(trailer[:4])
			if crc != z.crc {
				return n, &InvalidCRCError{crc}
			}

			dataSize := binary.LittleEndian.Uint64(trailer[4:12])
			if dataSize != z.dataSize {
				return n, &InvalidDataSizeError{dataSize}
			}

			memberSize := binary.LittleEndian.Uint64(trailer[12:])
			if memberSize != z.memberSize {
				return n, &InvalidMemberSizeError{memberSize}
			}

			if err := z.nextMember(); err != nil {
				return n, err
			}
			continue
		}

		return n, err
	}
}
