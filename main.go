// Calculate file entropy.
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

	"go.foxforensics.dev/go-mmap"

	"go.foxforensics.dev/entropy/entropy"
)

var Usage = `© 2026 Fox Forensics. Licensed under MIT License.
Usage: entropy PATH

Report bugs at: foxforensics.dev/issues`

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" {
		_, _ = fmt.Fprintln(os.Stderr, Usage)
		os.Exit(2)
	}

	err := filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		path, err = filepath.Abs(path)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return nil
		}

		f, err := os.Open(path)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return nil
		}

		defer func() { _ = f.Close() }()

		m, err := mmap.Map(f, mmap.RDONLY, 0)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			return nil
		}

		defer func() { _ = m.Unmap() }()

		e := entropy.Calculate(m)

		_, _ = fmt.Printf("%0.10f  %s\n", e, path)
		return nil
	})

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
