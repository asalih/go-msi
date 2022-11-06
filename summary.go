package msi

import (
	"bytes"
	"fmt"

	"github.com/asalih/go-mscfb"
)

type SummaryInfo struct {
	Properties *PropertySet
}

const defaultOsVersion = 10

var fmtIdSummaryInfo = []byte("\xe0\x85\x9f\xf2\xf9\x4f\x68\x10\xab\x91\x08\x00\x2b\x27\xb3\xd9")

func NewSummary() *SummaryInfo {
	propset := NewPropertySet(Win32, defaultOsVersion, fmtIdSummaryInfo)
	propset.CodePage = CodePageDefault()

	return &SummaryInfo{
		Properties: propset,
	}
}

func (s *SummaryInfo) ReadSummaryInfo(reader *mscfb.Stream) (*SummaryInfo, error) {
	propertySet, err := ReadPropertySet(reader)
	if err != nil {
		return nil, err
	}

	if propertySet.FmtID != nil && !bytes.Equal(propertySet.FmtID, fmtIdSummaryInfo) {
		return nil, fmt.Errorf("invalid property set format id")
	}

	s.Properties = propertySet

	return s, nil
}
