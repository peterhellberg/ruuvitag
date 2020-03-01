package ruuvitag

import "bytes"

var manufacturerIDRAWv1 = []byte{0x99, 0x4, 0x3}

// RAWv1 Data Format 3
type RAWv1 struct {
	DataFormat   uint8        // 3
	Humidity     float64      // 0.0 % to 100 % in 0.5 % increments.
	Temperature  float64      // -127.99 °C to +127.99 °C in 0.01 °C increments.
	Pressure     uint32       // 50000 Pa to 115536 Pa in 1 Pa increments.
	Acceleration Acceleration // -32000 to 32000 (mG), however the sensor on RuuviTag supports only 16 G max
	Battery      uint16       // 0 mV to 65536 mV in 1 mV increments, practically 1800 ... 3600 mV.
}

// IsRAWv1 tests wether the byte slice data begins with the prefix 0x99 0x4 0x3
func IsRAWv1(data []byte) bool {
	return bytes.HasPrefix(data, manufacturerIDRAWv1)
}

// ParseRAWv1 from Bluetooth "Manufacturer Specific Data"-field.
//
// Data Format 3 Protocol Specification (RAWv1) can be found at:
//
// https://github.com/ruuvi/ruuvi-sensor-protocols#data-format-3-protocol-specification-rawv1
//
func ParseRAWv1(data []byte) (RAWv1, error) {
	if len(data) != 16 {
		return RAWv1{}, ErrDataLength
	}

	if !bytes.HasPrefix(data, manufacturerID) {
		return RAWv1{}, ErrManufacturerID
	}

	var d dataFormat3

	r := bytes.NewReader(data[2:16])

	if err := readBigEndianBinary(r, &d); err != nil {
		return RAWv1{}, err
	}

	return d.toRAWv1()
}

type dataFormat3 struct {
	DataFormat          uint8
	Humidity            uint8
	Temperature         uint8
	TemperatureFraction uint8
	Pressure            uint16
	X, Y, Z             int16
	Battery             uint16
}

func (d dataFormat3) toRAWv1() (RAWv1, error) {
	return RAWv1{
		DataFormat:   d.dataFormat(),
		Humidity:     d.humidity(),
		Temperature:  d.temperature(),
		Pressure:     d.pressure(),
		Acceleration: d.acceleration(),
		Battery:      d.battery(),
	}, d.validate()
}

func (d dataFormat3) validate() error {
	return nil
}

func (d dataFormat3) dataFormat() uint8 {
	return d.DataFormat
}

func (d dataFormat3) humidity() float64 {
	return float64(d.Humidity) * 0.5
}

func (d dataFormat3) temperature() float64 {
	t := float64(d.Temperature&127) + float64(d.TemperatureFraction)*0.01

	if d.Temperature&128 == 0 {
		return t
	}

	return -t
}

func (d dataFormat3) pressure() uint32 {
	return uint32(d.Pressure) + 50000
}

func (d dataFormat3) acceleration() Acceleration {
	return Acceleration{d.X, d.Y, d.Z}
}

func (d dataFormat3) battery() uint16 {
	return d.Battery
}
