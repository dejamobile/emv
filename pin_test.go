package emv

import (
	"testing"
	"encoding/hex"
	"log"
	"reflect"
	"time"
)

var (
	pinblockIso4Handler PinblockHandler = NewIso4(must(hex.DecodeString("1fe206a02909246b3d3d3cce8d72975a")))
	pan                                 = "2950000217619"
	pin                                 = "12345"
)

func TestWrapPinblockIso4(t *testing.T) {
	start := time.Now()
	pinblock, err := pinblockIso4Handler.wrap(pan, pin)
	if err != nil {
		t.Fatalf("test failed, cause = %s\n", err.Error())
	}
	log.Printf("pinblock = %s\n", h(pinblock))
	unwrappedPin, err := pinblockIso4Handler.unwrap(pan, pinblock)
	if err != nil {
		t.Fatalf("test failed, cause = %s\n", err.Error())
	}
	log.Printf("unwrapped pin = %s\n", unwrappedPin)

	if !reflect.DeepEqual(pin, unwrappedPin) {
		t.Fatalf("pin and unwrapped pin should be equal")
	}
	log.Printf("TestWrapPinblockIso4 took %s\n", time.Since(start))
}
