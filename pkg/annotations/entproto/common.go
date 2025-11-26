package entproto

import (
	"encoding/json"
	"entgo.io/ent/entc/gen"
	"errors"
)

const (
	fieldAnnotationName = "entproto_field_annotation"
	apiAnnotationName   = "entproto_api_annotation"
)

var (
	ErrAnnotationNotFound = errors.New("annotation not found")
)

func Get(annotations gen.Annotations, name string, v any) error {
	ann, ok := annotations[name]
	if !ok {
		return ErrAnnotationNotFound
	}
	data, err := json.Marshal(ann)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func bitTwiddling[T ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint](a T) []T {
	if a == 0 {
		return nil
	}
	vals := make([]T, 0)
	for v := a; v != 0; v &= v - 1 {
		lowest := v & -v
		vals = append(vals, lowest)
	}
	return vals
}
