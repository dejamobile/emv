package emv

import (
	"testing"
	"encoding/hex"
	"log"
	"reflect"
)

var (
	pinblockIso4Handler PinblockHandler = NewIso4(must(hex.DecodeString("01020304050607080102030405060708")))
	pan                                 = must(hex.DecodeString("00123456789012"))
	pin                                 = must(hex.DecodeString("1234"))
)

func TestWrapPinblockIso4(t *testing.T) {
	pinblock, err := pinblockIso4Handler.wrap(pan, pin)
	if err != nil {
		t.Fatalf("test failed, cause = %s\n", err.Error())
	}
	log.Printf("pinblock = %s\n", h(pinblock))
	unwrappedPin, err := pinblockIso4Handler.unwrap(pan, pinblock)
	if err != nil {
		t.Fatalf("test failed, cause = %s\n", err.Error())
	}
	log.Printf("unwrapped pin = %s\n", h(unwrappedPin))

	if !reflect.DeepEqual(pin, unwrappedPin) {
		t.Fatalf("pin and unwrapped pin should be equal")
	}
}
