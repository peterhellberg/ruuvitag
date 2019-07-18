package main

import (
	"context"
	"fmt"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/peterhellberg/ruuvitag"
)

func setup(ctx context.Context) context.Context {
	d, err := dev.DefaultDevice()
	if err != nil {
		panic(err)
	}
	ble.SetDefaultDevice(d)

	return ble.WithSigHandler(context.WithCancel(ctx))
}

func main() {
	ctx := setup(context.Background())

	ble.Scan(ctx, true, handler, filter)
}

func handler(a ble.Advertisement) {
	raw, err := ruuvitag.ParseRAWv1(a.ManufacturerData())
	if err == nil {
		fmt.Printf("[%s] RSSI: %3d: %+v\n", a.Addr(), a.RSSI(), raw)
	}
}

func filter(a ble.Advertisement) bool {
	return ruuvitag.IsRAWv1(a.ManufacturerData())
}
