package consts

import (
	"errors"
)

var (
	ErrNoData = errors.New("data not found")
	ErrCrypto = errors.New("crypto error")
)
