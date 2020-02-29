package main

import (
	"fmt"

	"github.com/peterhellberg/ruuvitag"
	"github.com/raff/goble"
)

func main() {
	var quit chan bool

	ble := goble.New()

	ble.On("stateChange", func(ev goble.Event) bool {
		switch ev.State {
		case "poweredOn":
			ble.StartScanning(nil, true)
			return false
		default:
			ble.StopScanning()
			quit <- true

			return true
		}
	})

	ble.On("discover", func(ev goble.Event) bool {
		p := ev.Peripheral
		a := p.Advertisement
		d := a.ManufacturerData

		switch {
		case ruuvitag.IsRAWv2(d):
			if raw, err := ruuvitag.ParseRAWv2(d); err == nil {
				fmt.Printf("[%s] RSSI: %d: %+v\n", p.Uuid, p.Rssi, raw)
			}
		case ruuvitag.IsRAWv1(d):
			if raw, err := ruuvitag.ParseRAWv1(d); err == nil {
				fmt.Printf("[%s] RSSI: %d: %+v\n", p.Uuid, p.Rssi, raw)
			}
		}

		return false
	})

	ble.Init()

	<-quit
}
