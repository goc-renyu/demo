package fun

import (
	"demo/errr"
	"os"
)

func OpenFile1() error {
	_, err := os.Open("filename.ext")
	if err != nil {
		return err
	}
	return nil
}

func OpenFile2() error {
	return errr.NewError()
}
