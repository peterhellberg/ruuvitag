# ruuvitag

<img src="https://ruuvi.com/assets/images/ruuvitag.jpg" align="right" width="280">

[![Build Status](https://travis-ci.org/peterhellberg/ruuvitag.svg?branch=master)](https://travis-ci.org/peterhellberg/ruuvitag)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/ruuvitag)](https://goreportcard.com/report/github.com/peterhellberg/ruuvitag)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/ruuvitag)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/ruuvitag#license-mit)

This is a Go package for decoding [RuuviTag](https://ruuvi.com/ruuvitag-specs/) sensor data.

(Currently only supports the [RAWv1](https://github.com/ruuvi/ruuvi-sensor-protocols#data-format-3-protocol-specification-rawv1) data format)

## Requirements

 - <https://ruuvi.com/ruuvitag-specs/> - RuuviTag Environmental Sensor

## Installation

    go get -u github.com/peterhellberg/ruuvitag

## Usage

### Minimal example

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

	if raw, err := ruuvitag.DecodeRAWv1(data); err == nil {
		fmt.Printf("%+v\n", raw)
	}
}
```

```sh
{DataFormat:3 Humidity:36.5 Temperature:-2.01 Pressure:99538 Battery:3091 Acceleration:{X:-7 Y:-44 Z:1060}}
```

### Using `ruuvitag` with [github.com/go-ble/ble](https://github.com/go-ble/ble)

```go
package main

import (
	"bytes"
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
	raw, err := ruuvitag.DecodeRAWv1(a.ManufacturerData())
	if err == nil {
		fmt.Printf("[%s] RSSI: %3d: %+v\n", a.Addr(), a.RSSI(), raw)
	}
}

func filter(a ble.Advertisement) bool {
	return bytes.HasPrefix(a.ManufacturerData(), []byte{0x99, 0x4, 0x3})
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

### Publishing to `NATS`

```go
package main

import (
	"bytes"
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
		raw, err := ruuvitag.DecodeRAWv1(a.ManufacturerData())
		if err == nil {
			c.Publish("ruuvitag."+a.Addr().String(), struct {
				ruuvitag.RAWv1
				RSSI int
				Time time.Time
			}{raw, a.RSSI(), time.Now()})
		}
	}, func(a ble.Advertisement) bool {
		return bytes.HasPrefix(a.ManufacturerData(), []byte{0x99, 0x4, 0x3})
	})

	c.Close()
}
```

<img src="https://data.gopher.se/gopher/viking-gopher.svg" align="right" width="30%" height="300">

## License (MIT)

Copyright (c) 2019 [Peter Hellberg](https://c7.se)

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
