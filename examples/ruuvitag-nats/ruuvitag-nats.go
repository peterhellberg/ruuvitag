package main

import (
	"context"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/nats-io/nats.go"
	"github.com/peterhellberg/ruuvitag"
)

func main() {
	d, err := dev.DefaultDevice()
	if err != nil {
		panic(err)
	}
	ble.SetDefaultDevice(d)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	ctx := ble.WithSigHandler(context.WithCancel(context.Background()))

	ble.Scan(ctx, true, func(a ble.Advertisement) {
		raw, err := ruuvitag.ParseRAWv1(a.ManufacturerData())
		if err == nil {
			c.Publish("ruuvitag."+a.Addr().String(), struct {
				ruuvitag.RAWv1
				RSSI int
				Time time.Time
			}{raw, a.RSSI(), time.Now()})
		}
	}, func(a ble.Advertisement) bool {
		return ruuvitag.IsRAWv1(a.ManufacturerData())
	})

	c.Close()
}
