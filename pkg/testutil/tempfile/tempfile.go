package tempfile

import (
	"io/ioutil"
	"os"

	"github.com/docker/docker/pkg/testutil/assert"
)

// TempFile is a temporary file that can be used with unit tests. TempFile
// reduces the boilerplate setup required in each test case by handling
// setup errors.
type TempFile struct {
	File *os.File
}

// NewTempFile returns a new temp file with contents
func NewTempFile(t assert.TestingT, prefix string, content string) *TempFile {
	file, err := ioutil.TempFile("", prefix+"-")
	assert.NilError(t, err)

	_, err = file.Write([]byte(content))
	assert.NilError(t, err)
	file.Close()
	return &TempFile{File: file}
}

// Name returns the filename
func (f *TempFile) Name() string {
	return f.File.Name()
}

// Remove removes the file
func (f *TempFile) Remove() error {
	return os.Remove(f.Name())
}

// TempDir is a temporary directory that can be used with unit tests, and
// removed at the end of the test case.
type TempDir struct {
	path string
}

// NewTempDir returns a new temp directory for use with unit tests
func NewTempDir(t assert.TestingT, prefix string) *TempDir {
	name, err := ioutil.TempDir("", prefix)
	assert.NilError(t, err)

	return &TempDir{path: name}
}

// Path returns the path to the directory
func (t *TempDir) Path() string {
	return t.path
}

// Remove the temporary directory
func (t *TempDir) Remove() error {
	return os.RemoveAll(t.path)
}
