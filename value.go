package msi

type ValueRef struct {
	IsNull bool
	IsInt  bool
	IsStr  bool

	Value interface{}
}

type Value interface{}

func (v *ValueRef) ToValue(pool *StringPool) Value {
	if v.IsNull {
		return nil
	}

	if v.IsInt {
		return (Value)(v.Value)
	}

	if v.IsStr {
		ref, ok := v.Value.(StringRef)
		if !ok {
			return nil
		}

		return Value(pool.Get(ref))
	}

	return nil
}
