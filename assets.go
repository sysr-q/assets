package assets

import (
	"io/ioutil"
	"os"
)

// Read is a superficial wrapper on top of io/ioutil.ReadAll, which simply calls
// os.Open(file) first. It exists mainly to be rewritten by the magic Assets
// preprocessor.
func Read(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// MustRead is identical to Read, however if an error is returned, MustRead
// will panic.
func MustRead(file string) []byte {
	b, err := Read(file)
	if err != nil {
		panic(err)
	}
	return b
}
