package msi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type PropertyValue struct {
	Empty    bool
	Null     bool
	I1       int8
	I2       int16
	I4       int32
	LpStr    string
	FileTime int64
}

func ReadPropValue(rdr io.ReadSeeker, codePage CodePage) (*PropertyValue, error) {
	var typeNumber uint32
	err := binary.Read(rdr, binary.LittleEndian, &typeNumber)
	if err != nil {
		return nil, err
	}

	switch typeNumber {
	case 0:
		return &PropertyValue{Empty: true}, nil
	case 1:
		return &PropertyValue{Null: true}, nil
	case 2:
		var value int16
		err = binary.Read(rdr, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return &PropertyValue{I2: value}, nil
	case 3:
		var value int32
		err = binary.Read(rdr, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return &PropertyValue{I4: value}, nil
	case 16:
		var value int8
		err = binary.Read(rdr, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return &PropertyValue{I1: value}, nil
	case 30:
		var length uint32
		err = binary.Read(rdr, binary.LittleEndian, &length)
		if err != nil {
			return nil, err
		}
		if length != 0 {
			length = length - 1
		}
		var value = make([]byte, 0, length)
		for i := 0; i < int(length); i++ {
			var b uint8
			err = binary.Read(rdr, binary.LittleEndian, &b)
			if err != nil {
				return nil, err
			}
			value = append(value, b)
		}
		var term uint8
		err = binary.Read(rdr, binary.LittleEndian, &term)
		if err != nil {
			return nil, err
		}
		if term != 0 {
			return nil, fmt.Errorf("invalid string terminator: %v", term)
		}

		str, err := codePage.Decode(value)
		if err != nil {
			return nil, err
		}

		return &PropertyValue{LpStr: str}, nil
	case 64:
		var value int64
		err = binary.Read(rdr, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		return &PropertyValue{FileTime: value}, nil
	default:
		return nil, fmt.Errorf("invalid property type: %v", typeNumber)
	}
}

func (p *PropertyValue) MinimumVersion() PropertyFormatVersion {
	if p.I1 == int8(PropertyFormatVersion1) {
		return PropertyFormatVersion1
	}
	return PropertyFormatVersion0
}
