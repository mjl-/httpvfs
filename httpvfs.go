/*
Package httpvfs implements functions to make a godoc vfs.Filesystem implement the http.FileSystem interface.

NOTE: It's recommended you use the newer package
https://godoc.org/bitbucket.org/mjl/httpasset instead (serves
"index.html" properly, doesn't return server errors on directory
listings).

Together with bitbucket.org/mjl/asset this makes it easy to serve static files from a zip file embedded in a Go binary.

See https://godoc.org/bitbucket.org/mjl/httpvfs for documentation.
See https://godoc.org/bitbucket.org/mjl/asset on how to use embedded zip files in Go binaries.

Public domain, created by Mechiel Lukkien.
*/
package httpvfs

import (
	"golang.org/x/tools/godoc/vfs"
	"net/http"
	"os"
)

// HttpVfs implements an Open function that returns a http.File. It implements http.FileSystem, for use by http.FileServer, so we can serve static files from a vfs.FileSystem.
type HttpVfs struct {
	vfs vfs.FileSystem
}

// New makes a new HttpVfs out of a godoc vfs.FileSystem.
func New(vfs vfs.FileSystem) *HttpVfs {
	return &HttpVfs{vfs: vfs}
}

func (fs *HttpVfs) Open(name string) (file http.File, err error) {
	f, err := fs.vfs.Open(name)
	if err != nil {
		return
	}
	file = &HttpVfsFile{file: f, path: name, httpVfs: fs}
	return
}

// HttpVfsFile implements http.File
type HttpVfsFile struct {
	file    vfs.ReadSeekCloser
	path    string
	httpVfs *HttpVfs
}

func (f *HttpVfsFile) Close() error {
	return f.file.Close()
}

func (f *HttpVfsFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *HttpVfsFile) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *HttpVfsFile) Stat() (os.FileInfo, error) {
	return f.httpVfs.vfs.Stat(f.path)
}

func (f *HttpVfsFile) Readdir(count int) ([]os.FileInfo, error) {
	return f.httpVfs.vfs.ReadDir(f.path)
}
