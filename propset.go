package msi

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/asalih/go-mscfb"
)

type OperatingSystem int

const (
	Win16 OperatingSystem = iota
	Mac
	Win32
)

const (
	BYTE_ORDER_MARK   uint16 = 0xfffe
	PROPERTY_CODEPAGE uint32 = 1
)

type PropertyFormatVersion uint16

const (
	PropertyFormatVersion0 PropertyFormatVersion = 0
	PropertyFormatVersion1 PropertyFormatVersion = 1
)

func (p PropertyFormatVersion) VersionNumber() int {
	switch p {
	case PropertyFormatVersion0:
		return 0
	case PropertyFormatVersion1:
		return 1
	default:
		return -1
	}
}

type PropertySet struct {
	OS        OperatingSystem
	OSVersion uint16
	CLSID     []byte
	FmtID     []byte
	CodePage  CodePage
	//todo: must be binary tree
	Properties map[uint32]*PropertyValue
}

func NewPropertySet(os OperatingSystem, osVersion uint16, fmtId []byte) *PropertySet {
	return &PropertySet{
		OS:         os,
		OSVersion:  osVersion,
		FmtID:      fmtId,
		CLSID:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		CodePage:   CodePageDefault(),
		Properties: make(map[uint32]*PropertyValue),
	}
}

func ReadPropertySet(reader *mscfb.Stream) (*PropertySet, error) {
	var byteOrder uint16
	err := binary.Read(reader, binary.LittleEndian, &byteOrder)
	if err != nil {
		return nil, err
	}

	if byteOrder != BYTE_ORDER_MARK {
		return nil, fmt.Errorf("invalid byte order mark")
	}

	var propertyFormatVersion uint16
	err = binary.Read(reader, binary.LittleEndian, &propertyFormatVersion)
	if err != nil {
		return nil, err
	}

	if propertyFormatVersion != uint16(PropertyFormatVersion1) &&
		propertyFormatVersion != uint16(PropertyFormatVersion0) {
		return nil, fmt.Errorf("invalid property format version")
	}

	var osVersion uint16
	err = binary.Read(reader, binary.LittleEndian, &osVersion)
	if err != nil {
		return nil, err
	}

	var os uint16
	err = binary.Read(reader, binary.LittleEndian, &os)
	if err != nil {
		return nil, err
	}

	switch os {
	case uint16(Win16):
	case uint16(Mac):
	case uint16(Win32):
		break
	default:
		return nil, fmt.Errorf("invalid operating system")
	}

	var clsid [16]byte
	err = binary.Read(reader, binary.LittleEndian, &clsid)
	if err != nil {
		return nil, err
	}

	var reserved uint32
	err = binary.Read(reader, binary.LittleEndian, &reserved)
	if err != nil {
		return nil, err
	}
	if reserved < 1 {
		return nil, fmt.Errorf("invalid reserved value: %v", reserved)
	}

	//section header
	var fmtId [16]byte
	err = binary.Read(reader, binary.LittleEndian, &fmtId)
	if err != nil {
		return nil, err
	}

	var sectionOffset uint32
	err = binary.Read(reader, binary.LittleEndian, &sectionOffset)
	if err != nil {
		return nil, err
	}

	//section
	_, err = reader.Seek(int64(sectionOffset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	var sectionSize uint32
	err = binary.Read(reader, binary.LittleEndian, &sectionSize)
	if err != nil {
		return nil, err
	}

	var propertyCount uint32
	err = binary.Read(reader, binary.LittleEndian, &propertyCount)
	if err != nil {
		return nil, err
	}

	var propertyOffset = make(map[uint32]uint32)

	for i := 0; i < int(propertyCount); i++ {
		var name uint32
		var offset uint32
		err = binary.Read(reader, binary.LittleEndian, &name)
		if err != nil {
			return nil, err
		}

		err = binary.Read(reader, binary.LittleEndian, &offset)
		if err != nil {
			return nil, err
		}

		if _, ok := propertyOffset[name]; ok {
			return nil, fmt.Errorf("duplicate property name: %v", name)
		}

		propertyOffset[name] = offset
	}

	var codePageRead CodePage
	offset, ok := propertyOffset[PROPERTY_CODEPAGE]
	if ok {
		_, err = reader.Seek(int64(sectionOffset)+int64(offset), io.SeekStart)
		if err != nil {
			return nil, err
		}

		propVal, err := ReadPropValue(reader, CodePageDefault())
		if err != nil {
			return nil, err
		}

		cp := CodePageFromID(int(propVal.I2))
		if cp == -1 {
			return nil, fmt.Errorf("invalid code page: %v", propVal.I2)
		}

	} else {
		codePageRead = CodePageDefault()
	}

	propertyValues := make(map[uint32]*PropertyValue)
	for name, offset := range propertyOffset {
		_, err = reader.Seek(int64(sectionOffset)+int64(offset), io.SeekStart)
		if err != nil {
			return nil, err
		}

		propVal, err := ReadPropValue(reader, codePageRead)
		if err != nil {
			return nil, err
		}

		if propVal.MinimumVersion() > PropertyFormatVersion(propertyFormatVersion) {
			return nil, fmt.Errorf("invalid property format version: %v", propVal.MinimumVersion())
		}

		propertyValues[name] = propVal
	}

	return &PropertySet{
		OS:         OperatingSystem(os),
		OSVersion:  osVersion,
		CLSID:      clsid[:],
		FmtID:      fmtId[:],
		CodePage:   codePageRead,
		Properties: propertyValues,
	}, nil
}
