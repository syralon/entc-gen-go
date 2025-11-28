module github.com/syralon/entc-gen-go

go 1.24.6

toolchain go1.24.10

replace (
	github.com/syralon/entc-gen-go/pkg/annotations => ./pkg/annotations
	github.com/syralon/entc-gen-go/proto => ./proto
)

require (
	entgo.io/ent v0.14.5
	github.com/dave/jennifer v1.7.1
	github.com/go-kratos/kratos/v2 v2.9.1
	github.com/go-openapi/inflect v0.19.0
	github.com/golang/protobuf v1.5.4
	github.com/google/gnostic v0.7.1
	github.com/google/wire v0.7.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/iancoleman/strcase v0.3.0
	github.com/jhump/protoreflect/v2 v2.0.0-beta.2
	github.com/syralon/entc-gen-go/pkg/annotations v0.0.0-00010101000000-000000000000
	github.com/syralon/entc-gen-go/proto v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/sdk v1.37.0
	google.golang.org/genproto/googleapis/api v0.0.0-20251111163417-95abcf5c77ba
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.10
)

require (
	ariga.io/atlas v0.32.1-0.20250325101103-175b25e1c1b9 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/google/gnostic-models v0.7.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/hashicorp/hcl/v2 v2.18.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251103181224-f26f9409b101 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
