# ruuvitag

<img src="https://ruuvi.com/assets/images/ruuvitag.jpg" align="right" width="280">

[![Build Status](https://travis-ci.org/peterhellberg/ruuvitag.svg?branch=master)](https://travis-ci.org/peterhellberg/ruuvitag)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/ruuvitag)](https://goreportcard.com/report/github.com/peterhellberg/ruuvitag)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/ruuvitag)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/ruuvitag#license-mit)

This is a Go package for decoding [RuuviTag](https://ruuvi.com/ruuvitag-specs/) sensor data.

Currently supports the
[RAWv1](https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_03.md) and
[RAWv2](https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md) data formats.

## Requirements

 - <https://ruuvi.com/ruuvitag-specs/> - RuuviTag Environmental Sensor

## Installation

    go get -u github.com/peterhellberg/ruuvitag

## Usage

### Minimal example

[embedmd]:# (examples/ruuvitag-minimal/ruuvitag-minimal.go)
```go
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
```

```sh
{DataFormat:3 Humidity:36.5 Temperature:-2.01 Pressure:99538 Battery:3091 Acceleration:{X:-7 Y:-44 Z:1060}}
```

### Using `ruuvitag` with [github.com/go-ble/ble](https://github.com/go-ble/ble)

[embedmd]:# (examples/ruuvitag-ble/ruuvitag-ble.go)
```go
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
```

```sh
[73cec0e0d7cf46a4b7cd90e847fd3564] RSSI: -80: {DataFormat:3 Humidity:45 Temperature:22.8 Pressure:99522 Battery:3097 Acceleration:{X:10 Y:-46 Z:1044}}
[40c9d52c44974cc296fe2b443aef6850] RSSI: -76: {DataFormat:3 Humidity:42.5 Temperature:23.17 Pressure:99457 Battery:3103 Acceleration:{X:-24 Y:-40 Z:1009}}
[40c9d52c44974cc296fe2b443aef6850] RSSI: -77: {DataFormat:3 Humidity:42.5 Temperature:23.17 Pressure:99457 Battery:3103 Acceleration:{X:-27 Y:-39 Z:1008}}
[73cec0e0d7cf46a4b7cd90e847fd3564] RSSI: -79: {DataFormat:3 Humidity:45 Temperature:22.8 Pressure:99522 Battery:3097 Acceleration:{X:10 Y:-43 Z:1048}}
[73cec0e0d7cf46a4b7cd90e847fd3564] RSSI: -79: {DataFormat:3 Humidity:45 Temperature:22.8 Pressure:99521 Battery:3103 Acceleration:{X:11 Y:-44 Z:1049}}
[73cec0e0d7cf46a4b7cd90e847fd3564] RSSI: -79: {DataFormat:3 Humidity:45 Temperature:22.8 Pressure:99522 Battery:3097 Acceleration:{X:10 Y:-42 Z:1050}}
[40c9d52c44974cc296fe2b443aef6850] RSSI: -77: {DataFormat:3 Humidity:42.5 Temperature:23.17 Pressure:99456 Battery:3103 Acceleration:{X:-30 Y:-42 Z:1012}}
[649a163a6ad4462fa3b7dbedcbe47e25] RSSI: -99: {DataFormat:3 Humidity:64 Temperature:16.26 Pressure:99628 Battery:3103 Acceleration:{X:11 Y:118 Z:1033}}
[73cec0e0d7cf46a4b7cd90e847fd3564] RSSI: -80: {DataFormat:3 Humidity:45 Temperature:22.8 Pressure:99523 Battery:3097 Acceleration:{X:7 Y:-45 Z:1049}}
[649a163a6ad4462fa3b7dbedcbe47e25] RSSI: -99: {DataFormat:3 Humidity:64 Temperature:16.26 Pressure:99628 Battery:3103 Acceleration:{X:12 Y:119 Z:1038}}
[40c9d52c44974cc296fe2b443aef6850] RSSI: -78: {DataFormat:3 Humidity:42.5 Temperature:23.17 Pressure:99457 Battery:3103 Acceleration:{X:-26 Y:-41 Z:1010}}
[40c9d52c44974cc296fe2b443aef6850] RSSI: -77: {DataFormat:3 Humidity:42.5 Temperature:23.17 Pressure:99457 Battery:3103 Acceleration:{X:-27 Y:-40 Z:1013}}
[649a163a6ad4462fa3b7dbedcbe47e25] RSSI: -99: {DataFormat:3 Humidity:64 Temperature:16.26 Pressure:99630 Battery:3109 Acceleration:{X:14 Y:116 Z:1035}}
```

### Using `ruuvitag` with [github.com/raff/goble](https://github.com/raff/goble) (compatible with macOS Mojave)

[embedmd]:# (examples/ruuvitag-goble/ruuvitag-goble.go)
```go
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
```

```sh
[1592c47816f74834a98dbea4e54b9812] RSSI: 27: {DataFormat:3 Humidity:48 Temperature:24.24 Pressure:101009 Acceleration:{X:-701 Y:-463 Z:595} Battery:3091}
[a7bd46a2cac146f580c61b2fda061bad] RSSI: -87: {DataFormat:3 Humidity:45 Temperature:25.48 Pressure:101092 Acceleration:{X:-123 Y:0 Z:1030} Battery:3097}
[1592c47816f74834a98dbea4e54b9812] RSSI: 27: {DataFormat:3 Humidity:48 Temperature:24.24 Pressure:101009 Acceleration:{X:-700 Y:-464 Z:594} Battery:3085}
[42078647994e4f999f556488ad8a21eb] RSSI: 27: {DataFormat:3 Humidity:46 Temperature:23.37 Pressure:100947 Acceleration:{X:-27 Y:-42 Z:1012} Battery:3079}
[1592c47816f74834a98dbea4e54b9812] RSSI: 16: {DataFormat:3 Humidity:48 Temperature:24.24 Pressure:101008 Acceleration:{X:-701 Y:-462 Z:593} Battery:3079}
[a7bd46a2cac146f580c61b2fda061bad] RSSI: -89: {DataFormat:3 Humidity:45.5 Temperature:25.36 Pressure:101091 Acceleration:{X:-126 Y:4 Z:1028} Battery:3103}
[a7bd46a2cac146f580c61b2fda061bad] RSSI: -89: {DataFormat:3 Humidity:45.5 Temperature:25.35 Pressure:101092 Acceleration:{X:-123 Y:2 Z:1023} Battery:3091}
```

### Publishing to `NATS`

[embedmd]:# (examples/ruuvitag-nats/ruuvitag-nats.go)
```go
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
```

<img src="https://data.gopher.se/gopher/viking-gopher.svg" align="right" width="30%" height="300">

## License (MIT)

Copyright (c) 2019-2020 [Peter Hellberg](https://c7.se)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
