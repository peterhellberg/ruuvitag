package ruuvitag

import "testing"

func TestDecodeRAWv1(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, tt := range []struct {
			data []byte
			want RAWv1
		}{
			{
				[]byte{
					0x99, 0x04, 0x03, 0x4d,
					0x17, 0x01, 0xc1, 0x87,
					0x00, 0x08, 0xff, 0xd5,
					0x04, 0x1a, 0x0c, 0x1f,
				},
				RAWv1{
					DataFormat:  3,
					Humidity:    38.5,
					Temperature: 23.01,
					Pressure:    99543,
					Battery:     3103,
					Acceleration: Acceleration{
						X: 8,
						Y: -43,
						Z: 1050,
					},
				},
			},
			{
				[]byte{
					0x99, 0x04, 0x03, 0x60,
					0x0f, 0x06, 0xc1, 0xc9,
					0xff, 0xf4, 0xff, 0xd0,
					0x04, 0x1c, 0x0c, 0x13,
				},
				RAWv1{
					DataFormat:  3,
					Humidity:    48,
					Temperature: 15.06,
					Pressure:    99609,
					Battery:     3091,
					Acceleration: Acceleration{
						X: -12,
						Y: -48,
						Z: 1052,
					},
				},
			},
			{
				[]byte{
					0x99, 0x04, 0x03, 0x49,
					0x82, 0x01, 0xc1, 0x82,
					0xff, 0xf9, 0xff, 0xd4,
					0x04, 0x24, 0x0c, 0x13,
				},
				RAWv1{
					DataFormat:  3,
					Humidity:    36.5,
					Temperature: -2.01,
					Pressure:    99538,
					Battery:     3091,
					Acceleration: Acceleration{
						X: -7,
						Y: -44,
						Z: 1060,
					},
				},
			},
		} {
			raw, err := DecodeRAWv1(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := raw, tt.want; got != want {
				t.Fatalf("raw = %+v, want %+v", got, want)
			}
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		t.Run("unexpected Manufacturer data length", func(t *testing.T) {
			data := []byte{}

			if _, err := DecodeRAWv1(data); err != ErrDataLength {
				t.Fatalf("expected error ErrDataLength")
			}
		})

		t.Run("unknown Manufacturer ID", func(t *testing.T) {
			data := []byte{
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0,
			}

			if _, err := DecodeRAWv1(data); err != ErrManufacturerID {
				t.Fatalf("expected ManufacturerID")
			}
		})
	})
}

func TestCalculateHumidity(t *testing.T) {
	for _, tt := range []struct {
		h    uint8
		want float64
	}{
		{0x00, 0},
		{0x80, 64},
		{0xC8, 100},
	} {
		if got := calculateHumidity(tt.h); got != tt.want {
			t.Fatalf("calculateHumidity(0x%X) = %f, want %f", tt.h, got, tt.want)
		}
	}
}

func TestCalculateTemperature(t *testing.T) {
	for _, tt := range []struct {
		d    uint8
		f    uint8
		want float64
	}{
		{0x00, 0x00, 0.0},
		{0x81, 0x45, -1.69},
		{0x01, 0x45, 1.69},
	} {
		if got := calculateTemperature(tt.d, tt.f); got != tt.want {
			t.Fatalf("calculateTemperature(0x%X, 0x%X) = %f, want %f", tt.d, tt.f, got, tt.want)
		}
	}
}

func TestCalculatePressure(t *testing.T) {
	for _, tt := range []struct {
		b    []byte
		want uint32
	}{
		{nil, 0},
		{[]byte{0x00, 0x00}, 50000},
		{[]byte{0xC8, 0x7D}, 101325}, // (average sea-level pressure)
		{[]byte{0xFF, 0xFF}, 115535},
	} {
		if got := calculatePressure(tt.b); got != tt.want {
			t.Fatalf("calculatePressure(%X) = %d, want %d", tt.b, got, tt.want)
		}
	}
}

func TestCalculateAcceleration(t *testing.T) {
	for _, tt := range []struct {
		b    []byte
		want Acceleration
	}{
		{nil, Acceleration{}},
		{[]byte{0xFC, 0x18, 0x03, 0xE8, 0x0, 0x0}, Acceleration{-1000, 1000, 0}},
	} {
		if got := calculateAcceleration(tt.b); got != tt.want {
			t.Fatalf("calculateAcceleration(%#v) = %+v, want %+v", tt.b, got, tt.want)
		}
	}
}

func TestCaluclateAxis(t *testing.T) {
	for _, tt := range []struct {
		b    []byte
		want int16
	}{
		{nil, 0},
		{[]byte{0xFC, 0x18}, -1000},
		{[]byte{0x03, 0xE8}, 1000},
	} {
		if got := calculateAxis(tt.b); got != tt.want {
			t.Fatalf("calculateAxis(%#v) = %d, want %d", tt.b, got, tt.want)
		}
	}
}

func TestError(t *testing.T) {
	var err error = Error("test error")

	if got, want := err.Error(), "test error"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
}
