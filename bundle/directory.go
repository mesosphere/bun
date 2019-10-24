package bundle

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

// DirType represent different types of the hosts.
type DirType string

const (
	// DTRoot is a bundle root directory
	DTRoot DirType = "root"
	// DTMaster directory
	DTMaster = "master"
	// DTAgent directory
	DTAgent = "agent"
	// DTPublicAgent directory
	DTPublicAgent = "public agent"
)

type directory struct {
	Type DirType
	Path string
}

type bulkCloser []io.Closer

func (bc bulkCloser) Close() error {
	var e []string
	for _, c := range bc {
		if err := c.Close(); err != nil {
			e = append(e, err.Error())
		}
	}
	if len(e) > 0 {
		return errors.New(strings.Join(e, "\n"))
	}
	return nil
}

type namer string

func (name namer) Name() string {
	return string(name)
}

// File is a safe way to access bundle files.
type File interface {
	io.ReadCloser
	Name() string
}

// OpenFile opens the files of the typeName file type.
// If the file is not found, it tries to open it from a correspondent .gzip archive.
// If the .gzip archive is not found as well then returns an error.
// Caller is responsible for closing the file.
func (d directory) OpenFile(typeName FileTypeName) (File, error) {
	fileType := GetFileType(typeName)
	ok := false
	for _, dirType := range fileType.DirTypes {
		if dirType == d.Type {
			ok = true
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("%v files do not belong to %v hosts",
			fileType.Name,
			d.Type)
	}
	var notFound []string
	for _, localPath := range fileType.Paths {
		filePath := path.Join(d.Path, localPath)
		file, err := os.Open(filePath)
		if err == nil {
			return file, nil // found
		}
		if !os.IsNotExist(err) {
			return nil, err // error
		}
		// not found
		notFound = append(notFound, filePath)
		// try to open correspondent .gz file
		filePath += ".gz"
		file, err = os.Open(filePath)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err // error
			}
			notFound = append(notFound, filePath)
			continue // not found
		}
		// found
		r, err := gzip.NewReader(file)
		if err != nil {
			return nil, err // error
		}
		return struct {
			io.Reader
			io.Closer
			namer
		}{io.Reader(r), bulkCloser{r, file}, namer(filePath)}, nil
	}
	return nil, fmt.Errorf("none of the following files are found:\n%v",
		strings.Join(notFound, "\n"))
}

// ReadJSON reads JSON-encoded data from the bundle file and stores the result in
// the value pointed to by v.
func (d directory) ReadJSON(typeName FileTypeName, v interface{}) error {
	fileType := GetFileType(typeName)
	if fileType.ContentType != CTJson {
		panic(fmt.Sprintf("Content of the %v file is not JSON", typeName))
	}
	file, err := d.OpenFile(typeName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("bun.directory.ReadJSON: Cannot close file: %v", err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

type SearchRequest struct {
	PatternS string
	PatternR *regexp.Regexp
	F        func(line string, n int) bool
}

func (d directory) ScanLines(t FileTypeName, f func(n int, line string) bool) error {
	file, err := d.OpenFile(t)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("bun.directory.FindLine: Cannot close file %v with error: %v",
				file.Name(), err)
			if strings.HasSuffix(file.Name(), ".gz") {
				log.Printf("The .gz file might be corrupted. Try to fix it with"+
					" the gzrecover command and run the check again:\n"+
					"1) brew install gzrt\n"+
					"2) gzrecover -o %v %v", strings.TrimSuffix(file.Name(), ".gz"),
					file.Name())
			}
		}
	}()
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		if f(i, line) {
			return nil
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	// Not found
	return nil
}
