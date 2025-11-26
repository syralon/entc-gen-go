package openapi

import (
	"encoding/json"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	openapiv3 "github.com/google/gnostic/openapiv3"
)

const (
	schemaAnnotationName = "openapi_annotation"
)

type schema struct {
	Schema *openapiv3.Schema `json:"schema"`
}

func (s *schema) Name() string {
	return schemaAnnotationName
}

func Schema(s *openapiv3.Schema) entc.Annotation {
	return &schema{Schema: s}
}

func GetSchema(annotations gen.Annotations) (*openapiv3.Schema, error) {
	ann, ok := annotations[schemaAnnotationName]
	if !ok {
		return nil, nil
	}
	// TODO skip json unmarshal
	data, err := json.Marshal(ann)
	if err != nil {
		return nil, err
	}
	s := &schema{}
	err = json.Unmarshal(data, s)
	if err != nil {
		return nil, err
	}
	return s.Schema, nil
}
