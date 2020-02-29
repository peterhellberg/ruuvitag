package ruuvitag

import "testing"

func TestIsRAWv1(t *testing.T) {
	for _, tt := range []struct {
		data []byte
		want bool
	}{
		{
			[]byte{
				0x99, 0x04, 0x03, 0x4d,
				0x17, 0x01, 0xc1, 0x87,
				0x00, 0x08, 0xff, 0xd5,
				0x04, 0x1a, 0x0c, 0x1f,
			},
			true,
		},
		{
			[]byte{
				0x88, 0x03, 0x06, 0x4d,
				0x17, 0x01, 0xc1, 0x87,
				0x00, 0x08, 0xff, 0xd5,
				0x04, 0x1a, 0x0c, 0x1f,
			},
			false,
		},
	} {
		if got, want := IsRAWv1(tt.data), tt.want; got != want {
			t.Fatalf("IsRAWv1(%#v) = %v, want %v", tt.data, got, want)
		}
	}
}

func TestParseRAWv1(t *testing.T) {
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
			raw, err := ParseRAWv1(tt.data)
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

			if _, err := ParseRAWv1(data); err != ErrDataLength {
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

			if _, err := ParseRAWv1(data); err != ErrManufacturerID {
				t.Fatalf("expected ManufacturerID, got %q", err)
			}
		})
	})
}

func TestDataFormat3(t *testing.T) {
	t.Run("humidity", func(t *testing.T) {
		for _, tt := range []struct {
			h    uint8
			want float64
		}{
			{0x00, 0},
			{0x80, 64},
			{0xC8, 100},
		} {
			d := dataFormat3{Humidity: tt.h}

			if got := d.humidity(); got != tt.want {
				t.Fatalf("d.humidity() = %f, want %f", got, tt.want)
			}
		}
	})

	t.Run("temperature", func(t *testing.T) {
		for _, tt := range []struct {
			t    uint8
			f    uint8
			want float64
		}{
			{0x00, 0x00, 0.0},
			{0x81, 0x45, -1.69},
			{0x01, 0x45, 1.69},
		} {
			d := dataFormat3{Temperature: tt.t, TemperatureFraction: tt.f}

			if got := d.temperature(); got != tt.want {
				t.Fatalf("d.temperature() = %f, want %f", got, tt.want)
			}
		}
	})

	t.Run("pressure", func(t *testing.T) {
		for _, tt := range []struct {
			p    uint16
			want uint32
		}{
			{0x0000, 50000},
			{0xC87D, 101325}, // (average sea-level pressure)
			{0xFFFF, 115535},
		} {
			d := dataFormat3{Pressure: tt.p}

			if got := d.pressure(); got != tt.want {
				t.Fatalf("d.pressure() = %d, want %d", got, tt.want)
			}
		}
	})
}
