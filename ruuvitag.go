package ruuvitag

import (
	"encoding/binary"
	"io"
)

// Prefix for manufacturer data handled by this package
var manufacturerID = []byte{0x99, 0x4}

// Error type
type Error string

// Error value
func (e Error) Error() string {
	return string(e)
}

// Error constants
const (
	ErrDataLength     = Error("unexpected Manufacturer data length")
	ErrInvalidValues  = Error("invalid values")
	ErrManufacturerID = Error("unknown Manufacturer ID")
)

// Acceleration vector in mG
type Acceleration struct {
	X int16
	Y int16
	Z int16
}

func readBigEndianBinary(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}
