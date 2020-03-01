package ruuvitag

import "bytes"

var manufacturerIDRAWv2 = []byte{0x99, 0x4, 0x5}

// RAWv2 Data Format 5
type RAWv2 struct {
	DataFormat   uint8        // 5
	Temperature  float64      // in 0.005 degrees
	Humidity     float64      // in 0.0025% (0-163.83% range, though realistically 0-100%)
	Pressure     uint32       // in 1 Pa units, with offset of -50 000 Pa
	Acceleration Acceleration // -32000 to 32000 (mG), however the sensor on RuuviTag supports only 16 G max
	Battery      uint16       // in mV. first 11 bits is the battery voltage
	TXPower      int8         // last 5 bits unsigned are the TX power (in 2dBm steps)
	Movement     uint8        // incremented by motion detection interrupts from accelerometer
	Sequence     uint16       // each time a measurement is taken, this is incremented by one, used for measurement de-duplication
	MAC          [6]byte      // 48bit MAC address
}

// IsRAWv2 tests wether the byte slice data begins with the prefix 0x99 0x4 0x5
func IsRAWv2(data []byte) bool {
	return bytes.HasPrefix(data, manufacturerIDRAWv2)
}

func ParseRAWv2(data []byte) (RAWv2, error) {
	if len(data) != 26 {
		return RAWv2{}, ErrDataLength
	}

	if !bytes.HasPrefix(data, manufacturerID) {
		return RAWv2{}, ErrManufacturerID
	}

	var d dataFormat5

	r := bytes.NewReader(data[2:26])

	if err := readBigEndianBinary(r, &d); err != nil {
		return RAWv2{}, err
	}

	return d.toRAWv2()
}

type dataFormat5 struct {
	DataFormat  uint8
	Temperature int16
	Humidity    uint16
	Pressure    uint16
	X, Y, Z     int16
	PowerInfo   uint16
	Movement    uint8
	Sequence    uint16
	MAC         [6]byte
}

func (d dataFormat5) toRAWv2() (RAWv2, error) {
	return RAWv2{
		DataFormat:   d.dataFormat(),
		Temperature:  d.temperature(),
		Humidity:     d.humidity(),
		Pressure:     d.pressure(),
		Acceleration: d.acceleration(),
		Battery:      d.battery(),
		TXPower:      d.txPower(),
		Movement:     d.movement(),
		Sequence:     d.sequence(),
		MAC:          d.mac(),
	}, d.validate()
}

func (d dataFormat5) validate() error {
	switch {
	case
		d.Humidity == 65535,
		d.Pressure == 65535,
		d.Movement == 255,
		d.Sequence == 65535,
		d.PowerInfo>>5 == 2047,
		d.PowerInfo&0x1F == 31:
		return ErrInvalidValues
	}

	return nil
}

func (d dataFormat5) dataFormat() uint8 {
	return d.DataFormat
}

func (d dataFormat5) temperature() float64 {
	return float64(d.Temperature) * 0.005
}

func (d dataFormat5) humidity() float64 {
	return float64(d.Humidity) * 0.0025
}

func (d dataFormat5) pressure() uint32 {
	return uint32(d.Pressure) + 50000
}

func (d dataFormat5) acceleration() Acceleration {
	return Acceleration{d.X, d.Y, d.Z}
}

func (d dataFormat5) battery() uint16 {
	return d.PowerInfo>>5 + 1600
}

func (d dataFormat5) txPower() int8 {
	return int8((d.PowerInfo & 0x1F * 2) - 40)
}

func (d dataFormat5) movement() uint8 {
	return d.Movement
}

func (d dataFormat5) sequence() uint16 {
	return d.Sequence
}

func (d dataFormat5) mac() [6]byte {
	return d.MAC
}
