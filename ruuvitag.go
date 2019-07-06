package ruuvitag

import (
	"bytes"
	"encoding/binary"
)

// ManufacturerID sent by RuuviTags
const ManufacturerID = 0x9904

// RAWv1 Data Format 3
type RAWv1 struct {
	DataFormat   uint8   // 3
	Humidity     float64 // 0.0 % to 100 % in 0.5 % increments.
	Temperature  float64 // -127.99 °C to +127.99 °C in 0.01 °C increments.
	Pressure     uint32  // 50000 Pa to 115536 Pa in 1 Pa increments.
	Battery      uint16  // 0 mV to 65536 mV in 1 mV increments, practically 1800 ... 3600 mV.
	Acceleration         // -32000 to 32000 (mG), however the sensor on RuuviTag supports only 16 G max
}

// Acceleration vector
type Acceleration struct {
	X int16
	Y int16
	Z int16
}

// DecodeRAWv1 from "Manufacturer Specific Data"-field.
// Manufacturer ID is 0x0499. The actual data payload is:
//
// 		Offset 	Allowed values 	Description
// 		0 				3 							Data format definition (3 = current sensor readings)
// 		1 				0…200 					Humidity (one lsb is 0.5%, e.g. 128 is 64%)
// 		2 				-127…127 				Temperature (MSB is sign, next 7 bits are decimal value)
// 		3 				0…99 						Temperature (fraction, 1/100.)
// 		4 - 5 		0…65535 				Pressure (Most Significant Byte first, value - 50kPa)
// 		6-7 			-32767…32767 		Acceleration-X (Most Significant Byte first)
// 		8 - 9 		-32767…32767 		Acceleration-Y (Most Significant Byte first)
// 		10 - 11 	-32767…32767 		Acceleration-Z (Most Significant Byte first)
// 		12 - 13 	0…65535 				Battery voltage (millivolts). MSB First
//
// 	https://github.com/ruuvi/ruuvi-sensor-protocols#data-format-3-protocol-specification-rawv1
//
func DecodeRAWv1(data []byte) (RAWv1, error) {
	if len(data) != 16 {
		return RAWv1{}, ErrDataLength
	}

	var p struct {
		ManufacturerID      uint16
		DataFormat          uint8
		Humidity            uint8
		Temperature         uint8
		TemperatureFraction uint8
		Pressure            uint16
		X, Y, Z             int16
		BatteryVoltageMv    uint16
	}

	r := bytes.NewReader(data)

	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return RAWv1{}, err
	}

	if p.ManufacturerID != ManufacturerID {
		return RAWv1{}, ErrManufacturerID
	}

	return RAWv1{
		DataFormat:   p.DataFormat,
		Humidity:     calculateHumidity(p.Humidity),
		Temperature:  calculateTemperature(p.Temperature, p.TemperatureFraction),
		Pressure:     calculatePressure(p.Pressure),
		Acceleration: Acceleration{p.X, p.Y, p.Z},
		Battery:      p.BatteryVoltageMv,
	}, nil
}

func calculateHumidity(h uint8) float64 {
	return float64(h) * 0.5
}

func calculateTemperature(d, f uint8) float64 {
	t := float64(d&127) + float64(f)*0.01

	if d&128 == 0 {
		return t
	}

	return -t
}

func calculatePressure(p uint16) uint32 {
	return uint32(p) + 50000
}

// Error type
type Error string

func (e Error) Error() string {
	return string(e)
}

// Error constants
const (
	ErrDataLength     = Error("unexpected Manufacturer data length")
	ErrManufacturerID = Error("unknown Manufacturer ID")
)
