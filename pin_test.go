package emv

import (
	"testing"
	"log"
	"reflect"
	"time"
)

var (
	key                                 = make([]byte, 16)
	pinblockIso4Handler PinblockHandler = NewIso4(key)
	pan                                 = "2950000217619"
	pin                                 = "12345"
)

func TestWrapPinblockIso4(t *testing.T) {
	start := time.Now()
	pinblock, err := pinblockIso4Handler.Wrap(pan, pin)
	if err != nil {
		t.Fatalf("test failed, cause = %s\n", err.Error())
	}
	log.Printf("pinblock = %s\n", h(pinblock))
	unwrappedPin, err := pinblockIso4Handler.Unwrap(pan, pinblock)
	if err != nil {
		t.Fatalf("test failed, cause = %s\n", err.Error())
	}
	log.Printf("unwrapped pin = %s\n", unwrappedPin)

	if !reflect.DeepEqual(pin, unwrappedPin) {
		t.Fatalf("pin and unwrapped pin should be equal")
	}
	log.Printf("TestWrapPinblockIso4 took %s\n", time.Since(start))
}
