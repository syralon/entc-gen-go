package entproto

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"errors"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
	"path"
	"strings"
)

var ErrUnknownMethod = errors.New("unknown method")

//go:generate stringer -type APIMethod
type APIMethod uint8

const (
	APIGet APIMethod = 1 << iota
	APIList
	APICreate
	APIUpdate
	APIDelete
	APIAll = 1<<iota - 1
)

func (a APIMethod) Methods() []APIMethod {
	return bitTwiddling(a)
}

func (a APIMethod) Name() string {
	return strings.TrimPrefix(a.String(), "API")
}

func (a APIMethod) Rule(prefix string) (*googleapi.HttpRule, error) {
	rule := &googleapi.HttpRule{}
	switch a {
	case APIGet:
		rule.Pattern = &googleapi.HttpRule_Get{Get: path.Join(prefix, "{id}")}
	case APIList:
		rule.Pattern = &googleapi.HttpRule_Get{Get: prefix}
	case APICreate:
		rule.Pattern = &googleapi.HttpRule_Post{Post: prefix}
		rule.Body = "*"
	case APIUpdate:
		rule.Pattern = &googleapi.HttpRule_Post{Post: path.Join(prefix, "{id}")}
		rule.Body = "*"
	case APIDelete:
		rule.Pattern = &googleapi.HttpRule_Delete{Delete: path.Join(prefix, "{id}")}
	default:
		return nil, ErrUnknownMethod
	}
	return rule, nil
}

type APIOptions struct {
	Pattern string
	Method  APIMethod
}

type apiAnnotation struct {
	APIOptions
}

func (a *apiAnnotation) Name() string { return apiAnnotationName }

type APIOption func(a *apiAnnotation)

func WithAPIPattern(pattern string) func(a *apiAnnotation) {
	return func(a *apiAnnotation) {
		a.Pattern = pattern
	}
}

func WithAPIMethods(methods ...APIMethod) func(a *apiAnnotation) {
	return func(a *apiAnnotation) {
		a.Method = 0
		for _, method := range methods {
			a.Method |= method
		}
	}
}

func API(opts ...APIOption) entc.Annotation {
	a := &apiAnnotation{
		APIOptions: defaultAPIOptions,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

var defaultAPIOptions = APIOptions{
	Method: APIAll,
}

func GetAPIOptions(annotations gen.Annotations) (APIOptions, error) {
	s := &apiAnnotation{}
	err := Get(annotations, apiAnnotationName, s)
	if errors.Is(err, ErrAnnotationNotFound) {
		return defaultAPIOptions, nil
	}
	if err != nil {
		return APIOptions{}, err
	}
	return s.APIOptions, nil
}
