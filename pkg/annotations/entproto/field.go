package entproto

import (
	"entgo.io/ent/schema/field"
	"errors"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type Filter uint16

const (
	FilterEQ = 1 << iota
	FilterNE
	FilterGT
	FilterGTE
	FilterLT
	FilterLTE
	FilterBETWEEN
	FilterIN
	FilterAll = 1<<iota - 1
)

func (f Filter) Filters() []Filter {
	return bitTwiddling(f)
}

type FieldOptions struct {
	Name         string
	Visible      bool
	Filterable   bool
	Immutable    bool
	Settable     bool
	Filter       Filter
	Orderable    bool
	Type         field.Type
	TypeRepeated bool
}

type fieldAnnotation struct {
	FieldOptions
}

func (a *fieldAnnotation) Name() string { return fieldAnnotationName }

type FieldOption func(*fieldAnnotation)

func WithFieldName(name string) FieldOption {
	return func(a *fieldAnnotation) {
		a.FieldOptions.Name = name
	}
}

func WithFieldImmutable(immutable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Immutable = immutable
	}
}

func WithFieldSettable(settable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Settable = settable
	}
}

func WithFieldFilterable(filterable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Filterable = filterable
	}
}

func WithFieldVisible(visible bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Visible = visible
	}
}

func WithFieldOrderable(orderable bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Orderable = orderable
	}
}

func WithFieldFilter(filters ...Filter) FieldOption {
	return func(a *fieldAnnotation) {
		a.Filter = 0
		for _, f := range filters {
			a.Filter |= f
		}
	}
}

func WithFieldType(fieldType field.Type, repeated ...bool) FieldOption {
	return func(a *fieldAnnotation) {
		a.Type = fieldType
		a.TypeRepeated = len(repeated) > 0 && repeated[0]
	}
}

func Field(opts ...FieldOption) entc.Annotation {
	a := &fieldAnnotation{
		FieldOptions: defaultFieldOption,
	}
	for _, option := range opts {
		option(a)
	}
	return a
}

var defaultFieldOption = FieldOptions{
	Visible:    true,
	Filterable: true,
	Filter:     FilterAll,
}

func GetFieldOptions(annotations gen.Annotations) (FieldOptions, error) {
	s := &fieldAnnotation{}
	err := Get(annotations, fieldAnnotationName, s)
	if errors.Is(err, ErrAnnotationNotFound) {
		return defaultFieldOption, nil
	}
	if err != nil {
		return FieldOptions{}, err
	}
	return s.FieldOptions, nil
}
