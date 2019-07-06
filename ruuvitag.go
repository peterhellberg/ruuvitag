package ruuvitag

import (
	"bytes"
	"encoding/binary"
)

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

// DecodeRAWv1 from Bluetooth "Manufacturer Specific Data"-field.
//
// Data Format 3 Protocol Specification (RAWv1) can be found at:
//
// https://github.com/ruuvi/ruuvi-sensor-protocols#data-format-3-protocol-specification-rawv1
//
func DecodeRAWv1(data []byte) (RAWv1, error) {
	if len(data) != 16 {
		return RAWv1{}, ErrDataLength
	}

	if !bytes.HasPrefix(data, []byte{0x99, 0x04}) {
		return RAWv1{}, ErrManufacturerID
	}

	var p struct {
		DataFormat          uint8
		Humidity            uint8
		Temperature         uint8
		TemperatureFraction uint8
		Pressure            uint16
		X, Y, Z             int16
		BatteryVoltageMv    uint16
	}

	r := bytes.NewReader(data[2:16])

	if err := binary.Read(r, binary.BigEndian, &p); err != nil {
		return RAWv1{}, err
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
