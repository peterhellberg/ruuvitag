package ruuvitag

import "testing"

func TestIsRAWv2(t *testing.T) {
	for _, tt := range []struct {
		data []byte
		want bool
	}{
		{
			[]byte{
				0x99, 0x04, 0x05, 0x11, 0x14,
				0x30, 0x04, 0xc1, 0x43, 0xff,
				0xe0, 0xff, 0xdc, 0x03, 0xf0,
				0xa8, 0x76, 0xdc, 0x4a, 0x1c,
				0xfb, 0xa6, 0x88, 0xc5, 0x75,
				0x7f,
			},
			true,
		},
		{
			[]byte{
				0x99, 0x04, 0x03, 0x4d,
				0x17, 0x01, 0xc1, 0x87,
				0x00, 0x08, 0xff, 0xd5,
				0x04, 0x1a, 0x0c, 0x1f,
			},
			false,
		},
	} {
		if got, want := IsRAWv2(tt.data), tt.want; got != want {
			t.Fatalf("IsRAWv2(%#v) = %v, want %v", tt.data, got, want)
		}
	}
}

func TestParseRAWv2(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, tt := range []struct {
			data []byte
			err  error
			want RAWv2
		}{
			{ // Valid data
				[]byte{0x99, 0x04,
					0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3,
					0x7C, 0x00, 0x04, 0xFF, 0xFC, 0x04,
					0x0C, 0xAC, 0x36, 0x42, 0x00, 0xCD,
					0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
				},
				nil,
				RAWv2{
					DataFormat:  5,
					Temperature: 24.3,
					Humidity:    53.49,
					Pressure:    100044,
					Acceleration: Acceleration{
						X: 4,
						Y: -4,
						Z: 1036,
					},
					Battery:  2977,
					TXPower:  4,
					Movement: 66,
					Sequence: 205,
					MAC:      [6]byte{0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F},
				},
			},
			{ // Maximum values
				[]byte{0x99, 0x04,
					0x05, 0x7F, 0xFF, 0xFF, 0xFE, 0xFF,
					0xFE, 0x7F, 0xFF, 0x7F, 0xFF, 0x7F,
					0xFF, 0xFF, 0xDE, 0xFE, 0xFF, 0xFE,
					0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
				},
				nil,
				RAWv2{
					DataFormat:  5,
					Temperature: 163.835,
					Humidity:    163.835,
					Pressure:    115534,
					Acceleration: Acceleration{
						X: 32767,
						Y: 32767,
						Z: 32767,
					},
					Battery:  3646,
					TXPower:  20,
					Movement: 254,
					Sequence: 65534,
					MAC:      [6]byte{0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F},
				},
			},
			{ // Minimum values
				[]byte{0x99, 0x04,
					0x05, 0x80, 0x01, 0x00, 0x00, 0x00,
					0x00, 0x80, 0x01, 0x80, 0x01, 0x80,
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
				},
				nil,
				RAWv2{
					DataFormat:  5,
					Temperature: -163.835,
					Humidity:    0,
					Pressure:    50000,
					Acceleration: Acceleration{
						X: -32767,
						Y: -32767,
						Z: -32767,
					},
					Battery:  1600,
					TXPower:  -40,
					Movement: 0,
					Sequence: 0,
					MAC:      [6]byte{0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F},
				},
			},
			{ // Invalid values
				[]byte{0x99, 0x04,
					0x05, 0x80, 0x00, 0xFF, 0xFF, 0xFF,
					0xFF, 0x80, 0x00, 0x80, 0x00, 0x80,
					0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
				ErrInvalidValues,
				RAWv2{
					DataFormat:  5,
					Temperature: -163.84,
					Humidity:    163.8375,
					Pressure:    115535,
					Acceleration: Acceleration{
						X: -32768,
						Y: -32768,
						Z: -32768,
					},
					Battery:  3647,
					TXPower:  22,
					Movement: 255,
					Sequence: 65535,
					MAC:      [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
				},
			},
		} {
			raw, err := ParseRAWv2(tt.data)
			if err != tt.err {
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

			if _, err := ParseRAWv2(data); err != ErrDataLength {
				t.Fatalf("expected error ErrDataLength")
			}
		})

		t.Run("unknown Manufacturer ID", func(t *testing.T) {
			data := []byte{0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}

			if _, err := ParseRAWv2(data); err != ErrManufacturerID {
				t.Fatalf("expected ManufacturerID, got %q", err)
			}
		})
	})
}
