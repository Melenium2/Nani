package file_test

import (
	"Nani/internal/app/file"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

func removeFile(path string) error {
	return os.Remove(path)
}

func TestWriteLines_ShouldCreateFileAndWriteLinesToThemThenDelete_NoErrors(t *testing.T) {
	filename := "randomfile.txt"
	reader := file.New(filename)
	reader.WriteLines("privet", "poka", "hello")
	f, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.Greater(t, len(f), 0)

	assert.NoError(t, removeFile(filename))
}

func TestRead_ShoudReadFromFileFirst5Lines_NoError(t *testing.T) {
	filename := "tmp.txt"
	reader := file.New(filename)
	err := ioutil.WriteFile(filename, []byte("First\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst\n"), 0644)
	assert.NoError(t, err)
	r, err := reader.Read(5)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(r))

	assert.NoError(t, removeFile(filename))
}

func TestReadAll_ShoudReturnAllContentFromFile_NoError(t *testing.T) {
	filename := "tmp.txt"
	content := []byte("First\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst")
	reader := file.New(filename)
	err := ioutil.WriteFile(filename, content, 0644)
	assert.NoError(t, err)
	r, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, len(content), len(r))

	assert.NoError(t, removeFile(filename))
}

func TestReadAllSlice_ShouldReturnAllContentFromFileLineByLineInsideSlice_NoError(t *testing.T) {
	filename := "tmp.txt"
	content := []byte("First\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst")
	reader := file.New(filename)
	err := ioutil.WriteFile(filename, content, 0644)
	assert.NoError(t, err)
	r, err := reader.ReadAllSlice("\n")
	assert.NoError(t, err)
	assert.Equal(t, 8, len(r))

	assert.NoError(t, removeFile(filename))
}

func TestWriteLines_OneGShouldWriteToFileAfterAnotherG_NoErrorAndSequentialContent(t *testing.T) {
	filename := "tmp.txt"
	content := []byte("First\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst\nFirst")
	reader := file.New(filename)
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		t.Log("reading wait")
		time.Sleep(100)
		l, err := reader.Read(5)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(l))
		wait.Done()
		t.Log("Reading g comp")
	}()

	go func() {
		t.Log("writing g wait")
		var err error
		for i := 0; i < 100; i++ {
			err = reader.WriteLines(string(content))
		}
		assert.NoError(t, err)
		wait.Done()
		t.Log("writing g comp")
	}()

	t.Log("wait")
	var err error
	for i := 0; i < 100; i++ {
		err = reader.WriteLines(string(content))
	}
	assert.NoError(t, err)
	t.Log("complete")

	wait.Wait()

	assert.NoError(t, removeFile(filename))
}

func TestRead_ShouldReturnErrorCozFileNotFound_Error(t *testing.T) {
	reader := file.New("file.txt")
	r, err := reader.Read(100)
	assert.Error(t, err)
	assert.Len(t, r, 0)
}

func TestWriteLines_ShouldCreateNewFileThenWriteLines_NoError(t *testing.T) {
	filename := "newfile.txt"
	reader := file.New(filename)
	err := reader.WriteLines("Hi", "hello")
	assert.NoError(t, err)
	b, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, b)

	assert.NoError(t, removeFile(filename))
}
