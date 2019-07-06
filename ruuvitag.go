package ruuvitag

import (
	"bytes"
	"encoding/binary"
)

// ManufacturerID sent by RuuviTags
var ManufacturerID = []byte{0x99, 0x4}

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
	// Verify the length of the manufacturer data
	if len(data) != 16 {
		return RAWv1{}, ErrDataLength
	}

	// Verify expected Manufacturer ID
	if !bytes.HasPrefix(data, ManufacturerID) {
		return RAWv1{}, ErrManufacturerID
	}

	// Payload in the manufacturer data
	payload := data[2:16]

	return RAWv1{
		DataFormat:   payload[0],
		Humidity:     calculateHumidity(payload[1]),
		Temperature:  calculateTemperature(payload[2], payload[3]),
		Pressure:     calculatePressure(payload[4:6]),
		Acceleration: calculateAcceleration(payload[6:12]),
		Battery:      calculateBattery(payload[12:14]),
	}, nil
}

func calculateHumidity(h byte) float64 {
	return float64(h) * 0.5
}

func calculateTemperature(d, f uint8) float64 {
	t := float64(d&127) + float64(f)*0.01

	if d&128 == 0 {
		return t
	}

	return -t
}

func calculatePressure(b []byte) uint32 {
	if b == nil {
		return 0
	}

	return 50000 + uint32(binary.BigEndian.Uint16(b))
}

func calculateAcceleration(b []byte) Acceleration {
	if len(b) != 6 {
		return Acceleration{}
	}

	return Acceleration{
		X: calculateAxis(b[0:2]),
		Y: calculateAxis(b[2:4]),
		Z: calculateAxis(b[4:6]),
	}
}

func calculateAxis(b []byte) int16 {
	if b == nil {
		return 0
	}

	return int16(binary.BigEndian.Uint16(b))
}

func calculateBattery(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
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
