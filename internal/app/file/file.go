package file

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

// Read all content from given file and return string content
// @path: string (path to file)
// @return string (content from file), @error Error
func ReadAll(path string) (string, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(d), nil
}

type LocalFileWorker interface {
	ChangePath(path string)
	Read(lines ...int) ([]string, error)
	ReadAll() ([]byte, error)
	ReadAllSlice(separator ...string) ([]string, error)
	WriteLines(lines ...string) error
}

// Struct to help you work with files
type FileReader struct {
	path  string
	mutex sync.RWMutex
	debug bool
}

// Change path to file
// @Path: string (new file path)
func (f *FileReader) ChangePath(path string) {
	f.path = path
}

// Read 0 or N lines from file and return slice
// @Lines: slice (may contain the number of lines to be parsed)
// @return []string (slice of strings with lines from file) @error Error
func (f *FileReader) Read(lines ...int) ([]string, error) {
	var scanLines = 0
	if len(lines) > 0 {
		scanLines = lines[0]
	}

	if f.debug {
		log.Print("Locked")
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanLines)
	l := make([]string, 0)
	var i = 0
	for i < scanLines && sc.Scan() {
		l = append(l, sc.Text())
		i++
	}
	file.Close()

	if f.debug {
		log.Print("Unlocked")
	}

	return l, nil
}

// Read all bytes from file and return slice of bytes
// @return []byte (slice with file content) @error Error
func (f *FileReader) ReadAll() ([]byte, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return ioutil.ReadFile(f.path)
}

// Read all bytes from file and return slice of string
// @separator: []string (may contain separator for strings)
// @return []string (lines from file) @error Error
func (f *FileReader) ReadAllSlice(separator ...string) ([]string, error) {
	content, err := f.ReadAll()
	if err != nil {
		return nil, err
	}

	var sep = ""
	if len(separator) > 0 {
		sep = separator[0]
	}

	return strings.Split(string(content), sep), nil
}

// Write lines to file
// @lines: []string (should contain lines for writing)
// @return Error
func (f *FileReader) WriteLines(lines ...string) error {
	if f.debug {
		log.Print("Locked")
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}
	file.Close()

	if f.debug {
		log.Print("Unlocked")
	}
	return nil
}

// Create new instance of FileReader
// @path: string (path to file)
// @debug []bool (may contain debug flag for debugging)
// @return *FileReader
func New(path string, debug ...bool) *FileReader {
	var d = false
	if len(debug) > 0 {
		d = debug[0]
	}

	return &FileReader{
		path:  path,
		debug: d,
	}
}
