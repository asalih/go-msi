package msi

type ColumnBuilder struct {
	Name          string
	IsLocalizable bool
	IsNullable    bool
	IsPrimarykey  bool
	ValueRange    valueRange
	ForeignKey    foreignKey
	Category      Category
	EnumValues    []string
}

func NewColumnBuilder(name string) *ColumnBuilder {
	return &ColumnBuilder{
		Name:       name,
		EnumValues: make([]string, 0),
	}
}

// Makes the column be localizable.
func (b *ColumnBuilder) SetLocalizable() *ColumnBuilder {
	b.IsLocalizable = true
	return b
}

// Makes the column allow null values.
func (b *ColumnBuilder) SetNullable() *ColumnBuilder {
	b.IsNullable = true
	return b
}

// Makes the column be a primary key column.
func (b *ColumnBuilder) SetPrimaryKey() *ColumnBuilder {
	b.IsPrimarykey = true
	return b
}

// Makes the column only permit values in the given range.
func (b *ColumnBuilder) SetRange(min, max int32) *ColumnBuilder {
	b.ValueRange = valueRange{Min: min, Max: max}
	return b
}

// Makes the column refer to a key column in another table.
func (b *ColumnBuilder) SetForeignKey(tableName string, colIndex int32) *ColumnBuilder {
	b.ForeignKey = foreignKey{
		TableName:   tableName,
		ColumnIndex: colIndex,
	}
	return b
}

// For string columns, makes the column use the specified data format.
func (b *ColumnBuilder) SetCategory(category Category) *ColumnBuilder {
	b.Category = category
	return b
}

// Makes the column only permit the given values.
func (b *ColumnBuilder) SetEnumValues(values ...string) *ColumnBuilder {
	b.EnumValues = values
	return b
}

// Makes the column only permit the given values.
func (b *ColumnBuilder) Int16() *Column {
	return b.withType(ColumnTypeInt16)
}

// Makes the column only permit the given values.
func (b *ColumnBuilder) Int32() *Column {
	return b.withType(ColumnTypeInt32)
}

// Builds a column that stores a string.
func (b *ColumnBuilder) String(maxLen int) *Column {
	return b.withTypeSize(ColumnTypeStr, maxLen)
}

// Builds a column that stores an identifier string.
func (b *ColumnBuilder) IDString(maxLen int) *Column {
	return b.SetCategory(CategoryIdentifier).String(maxLen)
}

// Builds a column that stores a text string.
func (b *ColumnBuilder) TextString(maxLen int) *Column {
	return b.SetCategory(CategoryText).String(maxLen)
}

// Builds a column that stores a formatted string.
func (b *ColumnBuilder) FormattedString(maxLen int) *Column {
	return b.SetCategory(CategoryFormatted).String(maxLen)
}

// Builds a column that refers to a binary data stream.
func (b *ColumnBuilder) Binary() *Column {
	return b.SetCategory(CategoryBinary).String(0)
}

// Makes the column only permit the given values.
func (b *ColumnBuilder) withType(colType ColumnType) *Column {
	return b.withTypeSize(colType, 0)
}

// Makes the column only permit the given values.
func (b *ColumnBuilder) withTypeSize(colType ColumnType, size int) *Column {
	return &Column{
		Name:             b.Name,
		ColumnType:       colType,
		ColumnStringSize: size,
		IsLocalizable:    b.IsLocalizable,
		IsNullable:       b.IsNullable,
		IsPrimarykey:     b.IsPrimarykey,
		ValueRange:       b.ValueRange,
		ForeignKey:       b.ForeignKey,
		Category:         b.Category,
		EnumValues:       b.EnumValues,
	}
}

func (b *ColumnBuilder) withBitFields(typeBits int32) (*Column, error) {
	isNullable := typeBits&COL_NULLABLE_BIT != 0

	colType, size, err := ColumnTypeFromBitField(typeBits)
	if err != nil {
		return nil, err
	}

	return &Column{
		Name:             b.Name,
		ColumnType:       colType,
		ColumnStringSize: size,
		IsLocalizable:    typeBits&COL_LOCALIZABLE_BIT != 0,
		IsNullable:       isNullable || b.IsNullable,
		IsPrimarykey:     typeBits&COL_PRIMARY_KEY_BIT != 0,
		ValueRange:       b.ValueRange,
		ForeignKey:       b.ForeignKey,
		Category:         b.Category,
		EnumValues:       b.EnumValues,
	}, nil
}
