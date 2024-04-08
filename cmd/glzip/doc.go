// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

// Glzip reduces the size of files. It reads and writes of lzip format
// compressed files.
//
// Usage:
//
//	glzip [OPTIONS] <FILE>...
//
// Arguments:
//
//	<FILE>
//		Input file.
//
// Options:
//
//	-stdout
//		Write to standard output, keep input files.
//	-decompress
//		Decompress data.
//	-keep
//		Keep input files.
//	-dictionary-size <BYTE>
//		Set dictionary size in bytes.
//	-version
//		Print version number.
package main
