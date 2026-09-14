// Calculate file entropy and average.
//
// Usage:
//
//	entropy path
//
// The arguments are:
//
//	path
//	    File or folder to calculate (required).
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.foxforensics.eu/entropy/entropy"
	"go.foxforensics.eu/go-mmap"
)

var Usage = `© 2026 Fox Forensics. Licensed under MIT License.
Usage: entropy PATH

Report bugs at: foxforensics.eu/issues`

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" {
		_, _ = fmt.Fprintln(os.Stderr, Usage)
		os.Exit(2)
	}

	if err := filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		path, err = filepath.Abs(path)

		if err != nil {
			return err
		}

		f, err := os.Open(path)

		if err != nil {
			return err
		}

		defer func() {
			_ = f.Close()
		}()

		m, err := mmap.Map(f, mmap.RDONLY, 0)

		if err != nil {
			return err
		}

		defer func() {
			_ = m.Unmap()
		}()

		avg, ent := entropy.Calculate(m)

		_, _ = fmt.Printf("%0.10f %3d  %s\n", ent, avg, path)

		return nil
	}); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
