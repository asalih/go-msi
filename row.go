package msi

type Rows struct {
	StringPool   *StringPool
	Table        *Table
	Rows         [][]*ValueRef
	NextRowIndex int
}

type Row struct {
	Table  *Table
	Values []Value
}

func NewRows(stringPool *StringPool, table *Table, rows [][]*ValueRef) *Rows {
	return &Rows{
		StringPool: stringPool,
		Table:      table,
		Rows:       rows,
	}
}

func (r *Rows) Next() *Row {
	if r.NextRowIndex >= len(r.Rows) {
		return nil
	}

	rr := r.Rows[r.NextRowIndex]
	values := make([]Value, 0, len(rr))
	for _, v := range rr {
		values = append(values, v.ToValue(r.StringPool))
	}

	row := &Row{
		Table:  r.Table,
		Values: values,
	}

	r.NextRowIndex++

	return row
}

func NewRow(t *Table, values []Value) *Row {
	return &Row{
		Table:  t,
		Values: values,
	}
}
