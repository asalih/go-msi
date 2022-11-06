package msi

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/asalih/go-mscfb"
)

const (
	LONG_STRING_REFS_BIT uint32 = 0x8000_0000
	MAX_STRING_REF       int32  = 0xff_ffff
)

type StringPoolBuilder struct {
	CodePage           CodePage
	LongStringRefs     bool
	LengthAndRefCounts []stringPoolLRC
}

type StringPool struct {
	CodePage       CodePage
	Strings        []poolStrings
	LongStringRefs bool
	IsModified     bool
}

type StringRef struct {
	Num int32
}

func (s *StringRef) Index() int64 {
	return int64(s.Num - 1)
}

func (s *StringRef) Read(reader *mscfb.Stream, longStringRefs bool) (StringRef, error) {
	var numRef uint16
	err := binary.Read(reader, binary.LittleEndian, &numRef)
	if err != nil {
		return StringRef{}, err
	}

	num := int32(numRef)
	if longStringRefs {
		var lfr uint8
		err = binary.Read(reader, binary.LittleEndian, &lfr)
		if err != nil {
			return StringRef{}, err
		}

		num |= int32(lfr) << 16
	}

	return StringRef{Num: num}, nil
}

type stringPoolLRC struct {
	Length    uint32
	RefCounts uint16
}

type poolStrings struct {
	Value    string
	RefCount uint16
}

func (pool *StringPoolBuilder) ReadFromPool(stream *mscfb.Stream) error {
	if pool.LengthAndRefCounts == nil {
		pool.LengthAndRefCounts = make([]stringPoolLRC, 0)
	}

	var codepage uint32
	err := binary.Read(stream, binary.LittleEndian, &codepage)
	if err != nil {
		return nil
	}

	lsr := (codepage & LONG_STRING_REFS_BIT) != 0
	codepage = (codepage & ^LONG_STRING_REFS_BIT)
	codePageID := CodePageFromID(int(codepage))
	if codePageID == -1 {
		return fmt.Errorf("invalid codepage: %v", codePageID)
	}

	for {
		var len uint16
		var refCount uint16

		err = binary.Read(stream, binary.LittleEndian, &len)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = binary.Read(stream, binary.LittleEndian, &refCount)
		if err != nil {
			return err
		}

		if len == 0 && refCount > 0 {
			var lenW uint16
			err = binary.Read(stream, binary.LittleEndian, &lenW)
			if err != nil {
				return err
			}

			var refCountW uint16
			err = binary.Read(stream, binary.LittleEndian, &refCountW)
			if err != nil {
				return err
			}

			splrc := stringPoolLRC{
				Length:    (uint32(refCount) << 16) | uint32(lenW),
				RefCounts: refCountW,
			}
			pool.LengthAndRefCounts = append(pool.LengthAndRefCounts, splrc)
			continue
		}

		splrc := stringPoolLRC{
			Length:    uint32(len),
			RefCounts: refCount,
		}
		pool.LengthAndRefCounts = append(pool.LengthAndRefCounts, splrc)
	}

	pool.CodePage = codePageID
	pool.LongStringRefs = lsr

	return nil
}

func (pool *StringPoolBuilder) BuildFromData(stream *mscfb.Stream) (*StringPool, error) {
	strings := make([]poolStrings, 0)

	for _, ref := range pool.LengthAndRefCounts {
		cpd, err := pool.readExact(stream, ref)
		if err != nil {
			return nil, err
		}

		ps := poolStrings{
			Value:    cpd,
			RefCount: ref.RefCounts,
		}
		strings = append(strings, ps)
	}

	return &StringPool{
		CodePage:       pool.CodePage,
		Strings:        strings,
		LongStringRefs: pool.LongStringRefs,
		IsModified:     false,
	}, nil
}

func (pool *StringPoolBuilder) readExact(stream *mscfb.Stream, lrc stringPoolLRC) (string, error) {
	buf := make([]byte, lrc.Length)
	var total int
	var err error

	for total < int(lrc.Length) {
		num, err := stream.Read(buf[total:])
		if err != nil {
			return "", err
		}
		total += num
	}

	cpd, err := pool.CodePage.Decode(buf)
	if err != nil {
		return "", err
	}

	return cpd, nil
}

func (s *StringPool) Get(ref StringRef) string {
	index := ref.Index()
	if index < int64(len(s.Strings)) {
		return s.Strings[index].Value
	}

	return ""
}
