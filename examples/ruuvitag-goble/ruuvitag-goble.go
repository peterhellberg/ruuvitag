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

		if ruuvitag.IsRAWv1(a.ManufacturerData) {
			if raw, err := ruuvitag.ParseRAWv1(a.ManufacturerData); err == nil {
				fmt.Printf("[%s] RSSI: %d: %+v\n", p.Uuid, p.Rssi, raw)
			}
		}

		return false
	})

	ble.Init()

	<-quit
}
