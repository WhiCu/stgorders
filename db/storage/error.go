package storage

import (
	"errors"
	"fmt"
)

var (
	ErrPreloadCache = errors.New("could not preload cache")
)

func WrapErrPreloadCache(err error) error {
	return Wrap(ErrPreloadCache, err)
}

func Wrap(main error, sub error) error {
	return fmt.Errorf("%v: %v", main, sub)
}
