package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// FmtPretty converts any interface to JSON with indentation, for use in logging where better readability is required. Errors are ignored.
func FmtPretty(v interface{}) string {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// Ignore error
	}
	return string(jsonData)
}

// FmtGigabytes converts bytes to a printable string of gigabytes, rounded to the closest integer.
func FmtGigabytes(bytes uint64) string {
	return fmt.Sprintf("%.0fGB", float64(bytes)/1024/1024/1024)
}

func StringToBytes(sizeString string) (uint64, error) {
	var sizeBytes uint64
	var scaling uint64 = 1
	var err error

	if strings.HasSuffix(sizeString, "G") {
		sizeString = strings.TrimSuffix(sizeString, "G")
		scaling = 1024 * 1024 * 1024
	} else if strings.HasSuffix(sizeString, "M") {
		sizeString = strings.TrimSuffix(sizeString, "M")
		scaling = 1024 * 1024
	}

	sizeBytes, err = strconv.ParseUint(sizeString, 10, 64)
	if err != nil {
		return 0, err
	}
	sizeBytes = sizeBytes * scaling

	return sizeBytes, nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherwise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
// Source: https://stackoverflow.com/a/21067803
// More discussion: https://github.com/golang/go/issues/56172
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
