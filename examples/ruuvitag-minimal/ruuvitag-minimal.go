package main

import (
	"fmt"

	"github.com/peterhellberg/ruuvitag"
)

func main() {
	// Manufacturer data in RAWv1 format
	data := []byte{
		0x99, 0x04, 0x03, 0x49,
		0x82, 0x01, 0xc1, 0x82,
		0xff, 0xf9, 0xff, 0xd4,
		0x04, 0x24, 0x0c, 0x13,
	}

	if raw, err := ruuvitag.ParseRAWv1(data); err == nil {
		fmt.Printf("%+v\n", raw)
	}
}
