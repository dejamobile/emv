package emv

import (
	"log"
	"encoding/hex"
	"crypto/rand"
	"crypto/aes"
	"errors"
)

var (
	h = hex.EncodeToString
)

const (
	OFFSET_RANDOM_PIN = 8
	OFFSET_PIN        = 1
	BLOCK_FORMAT      = 4

	PIN_BLOCK_SIZE     = 16
	STARTED_OFFSET_PAN = 1

	FILL_DIGIT       = 10
	PAN_PADDING_INT  = 0
	PIN_FILLING_BYTE = 0xAA

	NIBBLE_SIZE  = 4
	MASK_FOR_PIN = 0x0F
	ASCII_OFFSET = 48
)

type PinblockHandler interface {
	Wrap(pan string, pin string) ([]byte, error)
	Unwrap(pan string, pinblock []byte) (pin string, err error)
}

type pinblockIso4 struct {
	key []byte
}

func (pb pinblockIso4) Wrap(pan string, pin string) (pinblock []byte, err error) {
	if len(pan) > 16 || len(pan) < 12 {
		return nil, errors.New("incorrect pan length (must be between 12 and 16")
	}
	if len(pin) > 8 || len(pin) < 4 {
		return nil, errors.New("incorrect pin length (must be between 4 and 8")
	}
	formattedPin := formatPINToPlainTextPinField(pin, BLOCK_FORMAT)
	log.Printf("pinblock : ", h(formattedPin))
	formattedPan := formatPANToPlainTextPANField(pan)
	interBlockA := encryptAes128Ecb(formattedPin, pb.key)
	interBlockB := xor(interBlockA, formattedPan)
	pinblock = encryptAes128Ecb(interBlockB, pb.key)
	return pinblock, nil
}

func (pb pinblockIso4) Unwrap(pan string, pinblock []byte) (pin string, err error) {
	if len(pan) > 16 || len(pan) < 12 {
		return "", errors.New("incorrect pan length (must be between 12 and 16")
	}
	interBlockB := decryptAes128Ecb(pinblock, pb.key)
	formattedPan := formatPANToPlainTextPANField(pan)
	interBlockA := xor(formattedPan, interBlockB)
	decryptedPinBlock := decryptAes128Ecb(interBlockA, pb.key)
	pin = retrievePinFromPinBlock(decryptedPinBlock)
	return pin, nil
}

func NewIso4(key []byte) pinblockIso4 {
	return pinblockIso4{key}
}

func formatPINToPlainTextPinField(pin string, blockFormat int) []byte {
	pinblock := make([]byte, PIN_BLOCK_SIZE)
	length := len(pin)
	pinblock[0] = byte(blockFormat<<NIBBLE_SIZE + length)
	for i := 0; i < length-1; i = i + 2 {
		pinblock[i/2+1] = (pin[i]-ASCII_OFFSET)<<NIBBLE_SIZE + (pin[i+1] - ASCII_OFFSET)
	}
	if length%2 != 0 {
		pinblock[length/2+1] = (pin[length-1]-ASCII_OFFSET)<<NIBBLE_SIZE + FILL_DIGIT
	}
	index := OFFSET_PIN + length/2 + length%2
	for i := index; i < OFFSET_RANDOM_PIN; i++ {
		pinblock[i] = PIN_FILLING_BYTE
	}
	rand.Read(pinblock[OFFSET_RANDOM_PIN:])
	return pinblock
}

func formatPANToPlainTextPANField(pan string) [] byte {
	formattedPan := make([]byte, PIN_BLOCK_SIZE)
	length := len(pan)
	panLength := length - 12
	formattedPan[0] = byte((panLength)<<NIBBLE_SIZE) + (pan[0] - ASCII_OFFSET)
	for i := STARTED_OFFSET_PAN; i < length-1; i = i + 2 {
		formattedPan[i/2+STARTED_OFFSET_PAN] = (pan[i]-48)<<NIBBLE_SIZE + (pan[i+1] - ASCII_OFFSET)
	}
	if (length-STARTED_OFFSET_PAN)%2 != 0 {
		formattedPan[length/2+STARTED_OFFSET_PAN] = (pan[length-1]-48)<<NIBBLE_SIZE + PAN_PADDING_INT
	}
	return formattedPan
}

func decryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	return doCipher(data, key, cipher.Decrypt)
}

func encryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	return doCipher(data, key, cipher.Encrypt)
}

func doCipher(in, key []byte, f cipherFunc) []byte {
	out := make([]byte, len(in))
	size := 16

	for bs, be := 0, size; bs < len(in); bs, be = bs+size, be+size {
		f(out[bs:be], in[bs:be])
	}
	return out
}

func retrievePinFromPinBlock(pinBlock []byte) string {
	length := int(pinBlock[0] & MASK_FOR_PIN)
	byteArray := pinBlock[OFFSET_PIN:OFFSET_PIN+length/2+length%2]
	pin := h(byteArray)
	if length%2 != 0 {
		pin = pin[:length]
	}
	return pin
}

type cipherFunc func(dst, src []byte)

func xor(a, b []byte) (out []byte) {
	out = make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		out[i] = a[i] ^ b[i]
	}
	return
}
