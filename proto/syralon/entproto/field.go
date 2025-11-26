package entproto

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

func Selectors[T ~func(*sql.Selector)](data ...[]func(*sql.Selector)) []T {
	t := make([]T, 0, len(data))
	for _, selectors := range data {
		for _, sel := range selectors {
			t = append(t, T(sel))
		}
	}
	return t
}

func (f *Int32Field) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *Int32Field_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *Int32Field_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *Int32Field_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *Int32Field_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *Int32Field_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *Int32Field_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *Int32Field_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *Int32Field_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *Int32Field_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	default:
		return nil
	}
}

func (f *Uint32Field) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *Uint32Field_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *Uint32Field_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *Uint32Field_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *Uint32Field_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *Uint32Field_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *Uint32Field_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *Uint32Field_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *Uint32Field_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *Uint32Field_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	default:
		return nil
	}
}

func (f *Int64Field) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *Int64Field_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *Int64Field_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *Int64Field_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *Int64Field_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *Int64Field_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *Int64Field_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *Int64Field_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *Int64Field_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *Int64Field_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	default:
		return nil
	}
}

func (f *Uint64Field) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *Uint64Field_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *Uint64Field_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *Uint64Field_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *Uint64Field_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *Uint64Field_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *Uint64Field_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *Uint64Field_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *Uint64Field_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *Uint64Field_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	default:
		return nil
	}
}

func (f *FloatField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *FloatField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *FloatField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *FloatField_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *FloatField_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *FloatField_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *FloatField_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *FloatField_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *FloatField_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *FloatField_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	default:
		return nil
	}
}

func (f *DoubleField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *DoubleField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *DoubleField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *DoubleField_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *DoubleField_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *DoubleField_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *DoubleField_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *DoubleField_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *DoubleField_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *DoubleField_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	default:
		return nil
	}
}

func (f *StringField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *StringField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *StringField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	case *StringField_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt)}
	case *StringField_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte)}
	case *StringField_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt)}
	case *StringField_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte)}
	case *StringField_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){sql.FieldGTE(column, v.Between.List[0]), sql.FieldLT(column, v.Between.List[1])}
	case *StringField_In:
		return []func(*sql.Selector){sql.FieldIn(column, v.In.List...)}
	case *StringField_NotIn:
		return []func(*sql.Selector){sql.FieldNotIn(column, v.NotIn.List...)}
	case *StringField_EqFold:
		return []func(*sql.Selector){sql.FieldEqualFold(column, v.EqFold)}
	case *StringField_Contains:
		return []func(*sql.Selector){sql.FieldContains(column, v.Contains)}
	case *StringField_HasPrefix:
		return []func(*sql.Selector){sql.FieldHasPrefix(column, v.HasPrefix)}
	case *StringField_HasSuffix:
		return []func(*sql.Selector){sql.FieldHasSuffix(column, v.HasSuffix)}
	default:
		return nil
	}
}

func (f *BoolField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *BoolField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *BoolField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	default:
		return nil
	}
}

func (f *DurationField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	var precision time.Duration
	switch f.GetPrecision() {
	case PRECISION_MILLISECOND:
		precision = time.Millisecond
	case PRECISION_SECOND:
		precision = time.Second
	case PRECISION_MINUTE:
		precision = time.Minute
	case PRECISION_HOUR:
		precision = time.Hour
	case PRECISION_DAY:
		precision = 24 * time.Hour
	default:
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *DurationField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq.AsDuration()/precision)}
	case *DurationField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne.AsDuration()/precision)}
	case *DurationField_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt.AsDuration()/precision)}
	case *DurationField_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte.AsDuration()/precision)}
	case *DurationField_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt.AsDuration()/precision)}
	case *DurationField_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte.AsDuration()/precision)}
	case *DurationField_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){
			sql.FieldGTE(column, v.Between.List[0].AsDuration()/precision),
			sql.FieldLT(column, v.Between.List[1].AsDuration()/precision),
		}
	default:
		return nil
	}
}

func (f *TimestampField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *TimestampField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq.AsTime())}
	case *TimestampField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne.AsTime())}
	case *TimestampField_Gt:
		return []func(*sql.Selector){sql.FieldGT(column, v.Gt.AsTime())}
	case *TimestampField_Gte:
		return []func(*sql.Selector){sql.FieldGTE(column, v.Gte.AsTime())}
	case *TimestampField_Lt:
		return []func(*sql.Selector){sql.FieldLT(column, v.Lt.AsTime())}
	case *TimestampField_Lte:
		return []func(*sql.Selector){sql.FieldLTE(column, v.Lte.AsTime())}
	case *TimestampField_Between:
		if len(v.Between.List) != 2 {
			return nil
		}
		return []func(*sql.Selector){
			sql.FieldGTE(column, v.Between.List[0].AsTime()),
			sql.FieldLT(column, v.Between.List[1].AsTime()),
		}
	default:
		return nil
	}
}

func (f *BytesField) Selector(column string) []func(*sql.Selector) {
	if f == nil {
		return nil
	}
	switch v := f.GetOperator().(type) {
	case *BytesField_Eq:
		return []func(*sql.Selector){sql.FieldEQ(column, v.Eq)}
	case *BytesField_Ne:
		return []func(*sql.Selector){sql.FieldNEQ(column, v.Ne)}
	default:
		return nil
	}
}
