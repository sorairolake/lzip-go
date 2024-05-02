// SPDX-FileCopyrightText: 2024 Shun Sakai
//
// SPDX-License-Identifier: Apache-2.0 OR MIT

package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmdtest"
)

var testFiles = []string{
	"foo.txt",
	"fox.lz",
	"fox_bm.lz",
	"fox_das46.lz",
	"fox_de20.lz",
	"fox_s11.lz",
	"fox_v2.lz",
}

func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func TestCLI(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := exec.Command("go", "build").Run(); err != nil {
		t.Fatal(err)
	}

	defer os.Remove("glzip")

	ts, err := cmdtest.Read("testdata")
	if err != nil {
		t.Fatal(err)
	}

	ts.Setup = func(rootDir string) error {
		for _, testFile := range testFiles {
			err := copyFile(filepath.Join(wd, "testdata", testFile), filepath.Join(rootDir, testFile))
			if err != nil {
				return err
			}
		}

		return nil
	}
	ts.Commands["glzip"] = cmdtest.Program("glzip")
	ts.Run(t, false)
}
