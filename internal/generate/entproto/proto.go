package entproto

type BuildOption func(*builder)

func (fn BuildOption) applyService(sb *ServiceBuilder) {
	fn(&sb.builder)
}

func (fn BuildOption) applyEnt(eb *EntBuilder) {
	fn(&eb.builder)
}

func WithProtoPackage(pkg string) BuildOption {
	return func(b *builder) {
		b.protoPackage = pkg
	}
}
func WithGoPackage(pkg string) BuildOption {
	return func(b *builder) {
		b.goPackage = pkg
	}
}

type builder struct {
	protoPackage string
	goPackage    string
}
