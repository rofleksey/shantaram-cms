package util

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidCredentials = errors.New("неверный логин или пароль")

func CopyFile(src, dst string) (err error) {
	r, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}

	defer func() {
		if err1 := w.Close(); err == nil {
			err = err1
		}
	}()

	if _, err = w.ReadFrom(r); err != nil {
		return fmt.Errorf("could not read from source: %w", err)
	}

	return err
}
