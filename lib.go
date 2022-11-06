package msi

import (
	"fmt"
	"io"
	"strings"

	"github.com/asalih/go-mscfb"
)

type columnMapValue struct {
	Index int
	Name  string
	Type  int
}
type columnMap map[string][]columnMapValue

type tableColumnKey struct {
	Table  string
	Column string
}

type MSIPackage struct {
	CompoundFile *mscfb.CompoundFile

	PackageType PackageType
	SummaryInfo *SummaryInfo
	StringPool  *StringPool
	Tables      map[string]*Table
}

func Open(rdr io.ReadSeeker) (*MSIPackage, error) {
	msiReader, err := mscfb.Open(rdr, mscfb.ValidationPermissive)
	if err != nil {
		return nil, err
	}

	rootEntry := msiReader.RootEntry()
	packageType := PackageTypeFromCLSID(rootEntry.CLSID)

	summaryStream, err := msiReader.OpenStream(SUMMARY_INFO_STREAM_NAME)
	if err != nil {
		return nil, err
	}

	summaryInfo := SummaryInfo{}
	_, err = summaryInfo.ReadSummaryInfo(summaryStream)
	if err != nil {
		return nil, err
	}

	stringTableStreamName := NameEncode(STRING_POOL_TABLE_NAME, true)
	stringTableStream, err := msiReader.OpenStream(stringTableStreamName)
	if err != nil {
		return nil, err
	}

	poolBuilder := StringPoolBuilder{}
	err = poolBuilder.ReadFromPool(stringTableStream)
	if err != nil {
		return nil, err
	}

	stringDataStreamName := NameEncode(STRING_DATA_TABLE_NAME, true)
	stringDataStream, err := msiReader.OpenStream(stringDataStreamName)
	if err != nil {
		return nil, err
	}

	stringPool, err := poolBuilder.BuildFromData(stringDataStream)
	if err != nil {
		return nil, err
	}

	// Read _Tables
	allTables := make(map[string]*Table)

	tablesTable := makeTablesTable(stringPool.LongStringRefs)
	tablesStreamName := tablesTable.StreamName()

	isTablesTableExist, err := msiReader.Exists(tablesStreamName)
	if err != nil {
		return nil, err
	}

	tableNames := make(map[string]struct{})
	if isTablesTableExist {
		tablesStream, err := msiReader.OpenStream(tablesStreamName)
		if err != nil {
			return nil, err
		}

		tr, err := tablesTable.ReadRows(tablesStream)
		if err != nil {
			return nil, err
		}

		rows := NewRows(stringPool, tablesTable, tr)
		for {
			row := rows.Next()
			if row == nil {
				break
			}

			name := row.Values[0].(string)
			if _, ok := tableNames[name]; ok {
				return nil, fmt.Errorf("duplicate table name: %s", name)
			}
			tableNames[name] = struct{}{}
		}
	}

	allTables[tablesTable.Name] = tablesTable

	// Read _Columns
	columnsTable := makeColumnsTable(stringPool.LongStringRefs)
	columnsTableStreamName := columnsTable.StreamName()

	isColumnsTableExist, err := msiReader.Exists(columnsTableStreamName)
	if err != nil {
		return nil, err
	}

	columnsMap := make(columnMap)
	for tableName := range tableNames {
		columnsMap[tableName] = make([]columnMapValue, 0)
	}

	if isColumnsTableExist {
		columnsStream, err := msiReader.OpenStream(columnsTableStreamName)
		if err != nil {
			return nil, err
		}

		cr, err := columnsTable.ReadRows(columnsStream)
		if err != nil {
			return nil, err
		}

		rows := NewRows(stringPool, columnsTable, cr)
		for {
			row := rows.Next()
			if row == nil {
				break
			}

			tableName := row.Values[0].(string)
			columnIndex := row.Values[1].(int)
			columnName := row.Values[2].(string)
			columnType := row.Values[3].(int)

			if _, ok := columnsMap[tableName]; !ok {
				return nil, fmt.Errorf("invalid table name: %s", tableName)
			}

			for _, v := range columnsMap[tableName] {
				if v.Index == columnIndex {
					return nil, fmt.Errorf("duplicate column index: %s.%d", tableName, columnIndex)
				}
			}

			cMapValue := columnMapValue{
				Index: columnIndex,
				Name:  columnName,
				Type:  columnType,
			}
			columnsMap[tableName] = append(columnsMap[tableName], cMapValue)
		}
	}

	allTables[columnsTable.Name] = columnsTable

	//Read _Validation

	validationMap := make(map[tableColumnKey][]*ValueRef)
	validationTable := makeValidationTable(stringPool.LongStringRefs)
	validationTableStreamName := validationTable.StreamName()

	isValidationTableExist, err := msiReader.Exists(validationTableStreamName)
	if err != nil {
		return nil, err
	}

	if isValidationTableExist {
		validationStream, err := msiReader.OpenStream(validationTableStreamName)
		if err != nil {
			return nil, err
		}
		rows, err := validationTable.ReadRows(validationStream)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			tableName := row[0].ToValue(stringPool).(string)
			columnName := row[1].ToValue(stringPool).(string)
			key := tableColumnKey{
				Table:  tableName,
				Column: columnName,
			}

			if _, ok := validationMap[key]; ok {
				return nil, fmt.Errorf("duplicate validation for table %s column %s", tableName, columnName)
			}

			validationMap[key] = row
		}
	}

	// Construct Table objects from column/validation data:
	for tableName, columnSpecs := range columnsMap {
		if len(columnSpecs) == 0 {
			return nil, fmt.Errorf("no columns for table %s", tableName)
		}

		len := len(columnSpecs)
		idx := 0
		lastIdx := len - 1

		if columnSpecs[idx].Index != 1 ||
			columnSpecs[lastIdx].Index != len {
			return nil, fmt.Errorf("table %s does not have a complete set of columns", tableName)
		}

		columns := make([]*Column, 0, len)
		for _, columnSpec := range columnSpecs {
			builder := NewColumnBuilder(columnSpec.Name)
			key := tableColumnKey{
				Table:  tableName,
				Column: columnSpec.Name,
			}

			if valueRefs, ok := validationMap[key]; ok {
				isNullable := valueRefs[2].ToValue(stringPool).(string) == "Y"
				if isNullable {
					builder.SetNullable()
				}

				minValue := valueRefs[3].ToValue(stringPool)
				maxValue := valueRefs[4].ToValue(stringPool)
				if minValue != nil && maxValue != nil {
					mi := int32(minValue.(int))
					ma := int32(maxValue.(int))
					builder.SetRange(mi, ma)
				}

				keyTable := valueRefs[5].ToValue(stringPool)
				keyColumn := valueRefs[6].ToValue(stringPool)
				if keyTable != nil && keyColumn != nil {
					kc := int32(keyColumn.(int))
					builder.SetForeignKey(keyTable.(string), kc)
				}

				categoryValue := valueRefs[7].ToValue(stringPool)
				if categoryValue != nil {
					c := CategoryFromString(categoryValue.(string))
					if c == -1 {
						return nil, fmt.Errorf("invalid category value: %s", categoryValue)
					}

					builder.SetCategory(c)
				}

				enumValues := valueRefs[8].ToValue(stringPool)
				if enumValues != nil {
					v := strings.Split(enumValues.(string), ";")
					builder.SetEnumValues(v...)
				}
			}

			col, err := builder.withBitFields(int32(columnSpec.Type))
			if err != nil {
				return nil, err
			}

			columns = append(columns, col)
		}

		table := NewTable(tableName, columns, stringPool.LongStringRefs)
		allTables[tableName] = table
	}

	return &MSIPackage{
		CompoundFile: msiReader,
		PackageType:  packageType,
		SummaryInfo:  &summaryInfo,
		StringPool:   stringPool,
		Tables:       allTables,
	}, nil
}

func (p *MSIPackage) Streams() *Streams {
	return NewStreams(p.CompoundFile.Directory.RootStorageEntries())
}

func (p *MSIPackage) ReadStream(streamName string) (io.ReadSeeker, error) {
	if !NameIsValid(streamName, false) {
		return nil, fmt.Errorf("invalid stream name: %s", streamName)
	}

	encoded := NameEncode(streamName, false)
	isStream, err := p.CompoundFile.IsStream(encoded)
	if err != nil {
		return nil, err
	}
	if !isStream {
		return nil, fmt.Errorf("stream %s does not exist", streamName)
	}

	return p.CompoundFile.OpenStream(encoded)
}

func makeTablesTable(longStringRefs bool) *Table {
	col := NewColumnBuilder("Name").SetPrimaryKey().String(64)

	return NewTable(TABLES_TABLE_NAME, []*Column{col}, longStringRefs)
}

func makeColumnsTable(longStringRefs bool) *Table {
	cols := []*Column{
		NewColumnBuilder("Table").SetPrimaryKey().String(64),
		NewColumnBuilder("Number").SetPrimaryKey().Int16(),
		NewColumnBuilder("Name").String(64),
		NewColumnBuilder("Type").Int16(),
	}

	return NewTable(COLUMNS_TABLE_NAME, cols, longStringRefs)
}

func makeValidationTable(longStringRefs bool) *Table {
	var min int32 = -0x7fff_ffff
	var max int32 = 0x7fff_ffff

	categoriesAsStrings := make([]string, len(AllCategories))
	for i, category := range AllCategories {
		categoriesAsStrings[i] = category.String()
	}

	cols := []*Column{
		NewColumnBuilder("Table").SetPrimaryKey().IDString(32),
		NewColumnBuilder("Column").SetPrimaryKey().IDString(32),
		NewColumnBuilder("Nullable").SetEnumValues("Y", "N").String(4),
		NewColumnBuilder("MinValue").SetNullable().SetRange(min, max).Int32(),
		NewColumnBuilder("MaxValue").SetNullable().SetRange(min, max).Int32(),
		NewColumnBuilder("KeyTable").SetNullable().IDString(255),
		NewColumnBuilder("KeyColumn").SetNullable().SetRange(1, 32).Int16(),
		NewColumnBuilder("Category").SetNullable().SetEnumValues(categoriesAsStrings...).String(32),
		NewColumnBuilder("Set").SetNullable().TextString(255),
		NewColumnBuilder("Description").SetNullable().TextString(255),
	}

	return NewTable(VALIDATION_TABLE_NAME, cols, longStringRefs)
}
