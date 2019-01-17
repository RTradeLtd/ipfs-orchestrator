// Code generated by fileb0x at "2019-01-16 18:56:44.366531 -0800 PST m=+0.020538983" from config file "b0x.yml" DO NOT EDIT.
// modification hash(436f5d4a341404d73dded968640eeefb.e5979db15ff7a7144261cbf60c4e3094)

package internal

import (
	"bytes"

	"context"
	"io"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/webdav"
)

var (
	// CTX is a context for webdav vfs
	CTX = context.Background()

	// FS is a virtual memory file system
	FS = webdav.NewMemFS()

	// Handler is used to server files through a http handler
	Handler *webdav.Handler

	// HTTP is the http file system
	HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct {
	// Prefix allows to limit the path of all requests. F.e. a prefix "css" would allow only calls to /css/*
	Prefix string
}

// FileIpfsInternalIpfsStartSh is "ipfs/internal/ipfs_start.sh"
var FileIpfsInternalIpfsStartSh = []byte("\x23\x21\x2f\x62\x69\x6e\x2f\x73\x68\x0a\x0a\x23\x20\x4d\x6f\x64\x69\x66\x69\x65\x64\x20\x49\x50\x46\x53\x20\x6e\x6f\x64\x65\x20\x69\x6e\x69\x74\x69\x61\x6c\x69\x7a\x61\x74\x69\x6f\x6e\x20\x73\x63\x72\x69\x70\x74\x2e\x0a\x23\x20\x4d\x6f\x75\x6e\x74\x20\x74\x6f\x20\x2f\x75\x73\x72\x2f\x6c\x6f\x63\x61\x6c\x2f\x62\x69\x6e\x2f\x73\x74\x61\x72\x74\x5f\x69\x70\x66\x73\x0a\x23\x20\x53\x6f\x75\x72\x63\x65\x3a\x20\x68\x74\x74\x70\x73\x3a\x2f\x2f\x67\x69\x74\x68\x75\x62\x2e\x63\x6f\x6d\x2f\x69\x70\x66\x73\x2f\x67\x6f\x2d\x69\x70\x66\x73\x2f\x62\x6c\x6f\x62\x2f\x24\x7b\x49\x50\x46\x53\x5f\x56\x45\x52\x53\x49\x4f\x4e\x7d\x2f\x62\x69\x6e\x2f\x63\x6f\x6e\x74\x61\x69\x6e\x65\x72\x5f\x64\x61\x65\x6d\x6f\x6e\x0a\x0a\x73\x65\x74\x20\x2d\x65\x0a\x0a\x23\x20\x61\x72\x67\x75\x6d\x65\x6e\x74\x73\x20\x70\x72\x6f\x76\x69\x64\x65\x64\x20\x74\x68\x72\x6f\x75\x67\x68\x20\x73\x74\x72\x69\x6e\x67\x20\x74\x65\x6d\x70\x6c\x61\x74\x65\x73\x0a\x44\x49\x53\x4b\x5f\x4d\x41\x58\x3d\x25\x64\x47\x42\x0a\x0a\x23\x20\x73\x65\x74\x20\x76\x61\x72\x69\x61\x62\x6c\x65\x73\x0a\x75\x73\x65\x72\x3d\x69\x70\x66\x73\x0a\x72\x65\x70\x6f\x3d\x22\x24\x49\x50\x46\x53\x5f\x50\x41\x54\x48\x22\x0a\x0a\x23\x20\x73\x65\x74\x20\x75\x73\x65\x72\x0a\x69\x66\x20\x5b\x20\x22\x24\x28\x69\x64\x20\x2d\x75\x29\x22\x20\x2d\x65\x71\x20\x30\x20\x5d\x3b\x20\x74\x68\x65\x6e\x0a\x20\x20\x65\x63\x68\x6f\x20\x22\x63\x68\x61\x6e\x67\x69\x6e\x67\x20\x75\x73\x65\x72\x20\x74\x6f\x20\x24\x75\x73\x65\x72\x22\x0a\x20\x20\x23\x20\x65\x6e\x73\x75\x72\x65\x20\x66\x6f\x6c\x64\x65\x72\x20\x69\x73\x20\x77\x72\x69\x74\x61\x62\x6c\x65\x0a\x20\x20\x73\x75\x2d\x65\x78\x65\x63\x20\x22\x24\x75\x73\x65\x72\x22\x20\x74\x65\x73\x74\x20\x2d\x77\x20\x22\x24\x72\x65\x70\x6f\x22\x20\x7c\x7c\x20\x63\x68\x6f\x77\x6e\x20\x2d\x52\x20\x2d\x2d\x20\x22\x24\x75\x73\x65\x72\x22\x20\x22\x24\x72\x65\x70\x6f\x22\x0a\x20\x20\x23\x20\x72\x65\x73\x74\x61\x72\x74\x20\x73\x63\x72\x69\x70\x74\x20\x77\x69\x74\x68\x20\x6e\x65\x77\x20\x70\x72\x69\x76\x69\x6c\x65\x67\x65\x73\x0a\x20\x20\x65\x78\x65\x63\x20\x73\x75\x2d\x65\x78\x65\x63\x20\x22\x24\x75\x73\x65\x72\x22\x20\x22\x24\x30\x22\x20\x22\x24\x40\x22\x0a\x66\x69\x0a\x0a\x23\x20\x63\x68\x65\x63\x6b\x20\x65\x78\x65\x63\x2c\x20\x72\x65\x70\x6f\x72\x74\x20\x76\x65\x72\x73\x69\x6f\x6e\x0a\x69\x70\x66\x73\x20\x76\x65\x72\x73\x69\x6f\x6e\x0a\x0a\x23\x20\x63\x68\x65\x63\x6b\x20\x66\x6f\x72\x20\x65\x78\x69\x73\x74\x69\x6e\x67\x20\x72\x65\x70\x6f\x20\x2d\x20\x6f\x74\x68\x65\x72\x77\x69\x73\x65\x20\x69\x6e\x69\x74\x20\x6e\x65\x77\x20\x6f\x6e\x65\x0a\x69\x66\x20\x5b\x20\x2d\x65\x20\x22\x24\x72\x65\x70\x6f\x2f\x63\x6f\x6e\x66\x69\x67\x22\x20\x5d\x3b\x20\x74\x68\x65\x6e\x0a\x20\x20\x65\x63\x68\x6f\x20\x22\x66\x6f\x75\x6e\x64\x20\x49\x50\x46\x53\x20\x66\x73\x2d\x72\x65\x70\x6f\x20\x61\x74\x20\x24\x72\x65\x70\x6f\x22\x0a\x65\x6c\x73\x65\x0a\x20\x20\x69\x70\x66\x73\x20\x69\x6e\x69\x74\x0a\x20\x20\x69\x70\x66\x73\x20\x63\x6f\x6e\x66\x69\x67\x20\x41\x64\x64\x72\x65\x73\x73\x65\x73\x2e\x41\x50\x49\x20\x2f\x69\x70\x34\x2f\x30\x2e\x30\x2e\x30\x2e\x30\x2f\x74\x63\x70\x2f\x35\x30\x30\x31\x0a\x20\x20\x69\x70\x66\x73\x20\x63\x6f\x6e\x66\x69\x67\x20\x41\x64\x64\x72\x65\x73\x73\x65\x73\x2e\x47\x61\x74\x65\x77\x61\x79\x20\x2f\x69\x70\x34\x2f\x30\x2e\x30\x2e\x30\x2e\x30\x2f\x74\x63\x70\x2f\x38\x30\x38\x30\x0a\x66\x69\x0a\x0a\x23\x20\x73\x65\x74\x20\x64\x61\x74\x61\x73\x74\x6f\x72\x65\x20\x71\x75\x6f\x74\x61\x0a\x69\x70\x66\x73\x20\x63\x6f\x6e\x66\x69\x67\x20\x44\x61\x74\x61\x73\x74\x6f\x72\x65\x2e\x53\x74\x6f\x72\x61\x67\x65\x4d\x61\x78\x20\x24\x44\x49\x53\x4b\x5f\x4d\x41\x58\x0a\x0a\x23\x20\x72\x65\x6c\x65\x61\x73\x65\x20\x6c\x6f\x63\x6b\x73\x0a\x69\x70\x66\x73\x20\x72\x65\x70\x6f\x20\x66\x73\x63\x6b\x0a\x0a\x23\x20\x69\x66\x20\x74\x68\x65\x20\x66\x69\x72\x73\x74\x20\x61\x72\x67\x75\x6d\x65\x6e\x74\x20\x69\x73\x20\x64\x61\x65\x6d\x6f\x6e\x0a\x69\x66\x20\x5b\x20\x22\x24\x31\x22\x20\x3d\x20\x22\x64\x61\x65\x6d\x6f\x6e\x22\x20\x5d\x3b\x20\x74\x68\x65\x6e\x0a\x20\x20\x23\x20\x66\x69\x6c\x74\x65\x72\x20\x74\x68\x65\x20\x66\x69\x72\x73\x74\x20\x61\x72\x67\x75\x6d\x65\x6e\x74\x20\x75\x6e\x74\x69\x6c\x0a\x20\x20\x23\x20\x68\x74\x74\x70\x73\x3a\x2f\x2f\x67\x69\x74\x68\x75\x62\x2e\x63\x6f\x6d\x2f\x69\x70\x66\x73\x2f\x67\x6f\x2d\x69\x70\x66\x73\x2f\x70\x75\x6c\x6c\x2f\x33\x35\x37\x33\x0a\x20\x20\x23\x20\x68\x61\x73\x20\x62\x65\x65\x6e\x20\x72\x65\x73\x6f\x6c\x76\x65\x64\x0a\x20\x20\x73\x68\x69\x66\x74\x0a\x65\x6c\x73\x65\x0a\x20\x20\x23\x20\x70\x72\x69\x6e\x74\x20\x64\x65\x70\x72\x65\x63\x61\x74\x69\x6f\x6e\x20\x77\x61\x72\x6e\x69\x6e\x67\x0a\x20\x20\x23\x20\x67\x6f\x2d\x69\x70\x66\x73\x20\x75\x73\x65\x64\x20\x74\x6f\x20\x68\x61\x72\x64\x63\x6f\x64\x65\x20\x22\x69\x70\x66\x73\x20\x64\x61\x65\x6d\x6f\x6e\x22\x20\x69\x6e\x20\x69\x74\x27\x73\x20\x65\x6e\x74\x72\x79\x70\x6f\x69\x6e\x74\x0a\x20\x20\x23\x20\x74\x68\x69\x73\x20\x77\x6f\x72\x6b\x61\x72\x6f\x75\x6e\x64\x20\x73\x75\x70\x70\x6f\x72\x74\x73\x20\x74\x68\x65\x20\x6e\x65\x77\x20\x73\x79\x6e\x74\x61\x78\x20\x73\x6f\x20\x70\x65\x6f\x70\x6c\x65\x20\x73\x74\x61\x72\x74\x20\x73\x65\x74\x74\x69\x6e\x67\x20\x64\x61\x65\x6d\x6f\x6e\x20\x65\x78\x70\x6c\x69\x63\x69\x74\x6c\x79\x0a\x20\x20\x23\x20\x77\x68\x65\x6e\x20\x6f\x76\x65\x72\x77\x72\x69\x74\x69\x6e\x67\x20\x43\x4d\x44\x0a\x20\x20\x65\x63\x68\x6f\x20\x22\x44\x45\x50\x52\x45\x43\x41\x54\x45\x44\x3a\x20\x61\x72\x67\x75\x6d\x65\x6e\x74\x73\x20\x68\x61\x76\x65\x20\x62\x65\x65\x6e\x20\x73\x65\x74\x20\x62\x75\x74\x20\x74\x68\x65\x20\x66\x69\x72\x73\x74\x20\x61\x72\x67\x75\x6d\x65\x6e\x74\x20\x69\x73\x6e\x27\x74\x20\x27\x64\x61\x65\x6d\x6f\x6e\x27\x22\x20\x3e\x26\x32\x0a\x66\x69\x0a\x0a\x65\x78\x65\x63\x20\x69\x70\x66\x73\x20\x64\x61\x65\x6d\x6f\x6e\x20\x22\x24\x40\x22\x0a")

func init() {
	err := CTX.Err()
	if err != nil {
		panic(err)
	}

	err = FS.Mkdir(CTX, "ipfs/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	err = FS.Mkdir(CTX, "ipfs/internal/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	var f webdav.File

	f, err = FS.OpenFile(CTX, "ipfs/internal/ipfs_start.sh", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FileIpfsInternalIpfsStartSh)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	Handler = &webdav.Handler{
		FileSystem: FS,
		LockSystem: webdav.NewMemLS(),
	}

}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {
	path = hfs.Prefix + path

	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(f)
	return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// WalkDirs looks for files in the given dir and returns a list of files in it
// usage for all files in the b0x: WalkDirs("", false)
func WalkDirs(name string, includeDirsInList bool, files ...string) ([]string, error) {
	f, err := FS.OpenFile(CTX, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	fileInfos, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	for _, info := range fileInfos {
		filename := path.Join(name, info.Name())

		if includeDirsInList || !info.IsDir() {
			files = append(files, filename)
		}

		if info.IsDir() {
			files, err = WalkDirs(filename, includeDirsInList, files...)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}
