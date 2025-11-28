package entproto

import (
	"errors"
	"path"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
	googleapi "google.golang.org/genproto/googleapis/api/annotations"
)

var ErrUnknownMethod = errors.New("unknown method")

//go:generate stringer -type APIMethod
type APIMethod uint8

const (
	GET APIMethod = 1 << iota
	LIST
	CREATE
	UPDATE
	DELETE
	ALL = 1<<iota - 1

	ReadOnly = GET | LIST
)

func (a APIMethod) Methods() []APIMethod {
	return bitTwiddling(a)
}

func (a APIMethod) Name() string {
	return strcase.ToCamel(a.String())
}

func (a APIMethod) Rule(prefix string) (*googleapi.HttpRule, error) {
	rule := &googleapi.HttpRule{}
	switch a {
	case GET:
		rule.Pattern = &googleapi.HttpRule_Get{Get: path.Join(prefix, "{id}")}
	case LIST:
		rule.Pattern = &googleapi.HttpRule_Get{Get: prefix}
	case CREATE:
		rule.Pattern = &googleapi.HttpRule_Post{Post: prefix}
		rule.Body = "*"
	case UPDATE:
		rule.Pattern = &googleapi.HttpRule_Put{Put: path.Join(prefix, "{id}")}
		rule.Body = "*"
	case DELETE:
		rule.Pattern = &googleapi.HttpRule_Delete{Delete: path.Join(prefix, "{id}")}
	default:
		return nil, ErrUnknownMethod
	}
	return rule, nil
}

type APIOptions struct {
	Pattern     string
	Method      APIMethod
	DisableEdge bool
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

func WithAPIDisableEdge(disable bool) func(a *apiAnnotation) {
	return func(a *apiAnnotation) {
		a.DisableEdge = disable
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
	Method: ALL,
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
