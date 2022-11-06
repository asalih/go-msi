package msi

import (
	"io"

	"github.com/asalih/go-mscfb"
)

type Table struct {
	Name           string
	Columns        []*Column
	LongStringRefs bool
}

func NewTable(name string, columns []*Column, longStringRefs bool) *Table {
	return &Table{
		Name:           name,
		Columns:        columns,
		LongStringRefs: longStringRefs,
	}
}

func (t *Table) StreamName() string {
	return NameEncode(t.Name, true)
}

func (t *Table) ReadRows(stream *mscfb.Stream) ([][]*ValueRef, error) {
	dataLength, err := stream.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	_, err = stream.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var rowSize uint64
	for _, column := range t.Columns {
		rowSize += column.ColumnType.Width(t.LongStringRefs)
	}

	var numRows int64
	if rowSize > 0 {
		numRows = int64(dataLength) / int64(rowSize)
	}

	rows := make([][]*ValueRef, numRows)
	for i := int64(0); i < numRows; i++ {
		rows[i] = make([]*ValueRef, 0)
	}

	for _, column := range t.Columns {
		for i := int64(0); i < numRows; i++ {

			value, err := column.ColumnType.ReadValue(stream, t.LongStringRefs)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			rows[i] = append(rows[i], value)
		}
	}

	return rows, nil
}
