package emv

import (
	"log"
	"encoding/hex"
)

var (
	h = hex.EncodeToString
)

type PinblockHandler interface {
	wrap(pan []byte, pin []byte) ([]byte, error)
	unwrap(pan []byte, pinblock []byte) ([]byte, error)
}

type pinblockIso4 struct {
	key []byte
}

func (pb pinblockIso4) wrap(pan []byte, pin []byte) (pinblock []byte, err error) {
	log.Printf("entering wrap with pan = %s pin = %s", h(pan), h(pin))
	defer log.Printf("exiting wrap with pinblock = %s", h(pinblock))
	return
}

func (pb pinblockIso4) unwrap(pan []byte, pinblock []byte) (pin []byte, err error) {
	log.Printf("entering unwrap with pan = %s pinblock = %s", h(pan), h(pinblock))
	defer log.Printf("exiting unwrap with pin = %s", h(pin))
	return
}

func NewIso4(key []byte) pinblockIso4 {
	return pinblockIso4{key: key}
}

func must(s []byte, err error) (out []byte) {
	if err == nil {
		out = s
	}
	return
}
