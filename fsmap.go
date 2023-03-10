package fsmap

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func mkdirTemp(dir string) (string, error) {
	try := 0
	for {
		name := filepath.Join(dir, strconv.FormatUint(rand.Uint64(), 10))
		err := os.Mkdir(name, os.ModePerm)
		if err == nil {
			return name, nil
		}
		if os.IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return "", &os.PathError{Op: "fsmap.mkdirtemp", Path: dir, Err: os.ErrExist}
		}
		if os.IsNotExist(err) {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return "", err
			}
		}
		return "", err
	}
}

func hashAlgorithm(f func() hash.Hash) func([]byte) string {
	return func(b []byte) string {
		h := f()
		h.Write(b)
		return hex.EncodeToString(h.Sum(nil))
	}
}

var algorithms = map[string]func([]byte) string{
	"SHA-512": hashAlgorithm(sha512.New),
}

type Fsmap struct {
	prefix string
	algo   func([]byte) string
}

func New(prefix string, algorithm string) (f *Fsmap, err error) {
	f = &Fsmap{prefix: prefix}

	var ok bool
	if f.algo, ok = algorithms[algorithm]; !ok {
		err = errors.New("fsmap: invalid algorithm")
	}

	return
}

func (f Fsmap) Lookup(key []byte, create bool) (path string, err error) {
	base := filepath.Join(f.prefix, f.algo(key))

	var fi os.FileInfo
	if fi, err = os.Stat(base); err != nil && (!create || !os.IsNotExist(err)) {
		return
	}

	if os.IsNotExist(err) {
		if err = os.Mkdir(base, os.ModePerm); err != nil {
			return
		}
	} else if !fi.IsDir() {
		err = os.ErrInvalid
		return
	}

	var entries []os.DirEntry
	if entries, err = os.ReadDir(base); err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			keypath := filepath.Join(base, entry.Name(), "key")
			var keydata []byte
			if keydata, err = ioutil.ReadFile(keypath); err != nil {
				return
			}

			if bytes.Equal(key, keydata) {
				path = filepath.Join(base, entry.Name())
				return
			}
		}
	}

	if !create {
		err = os.ErrNotExist
		return
	}

	if path, err = mkdirTemp(base); err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Join(path, "key"), key, 0644)
	return
}
