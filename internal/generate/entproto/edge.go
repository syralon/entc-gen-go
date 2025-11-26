package entproto

type (
	Edge interface {
		SetEdge(ctx Context) error
	}
	Edges []Edge
)

func (e Edges) SetEdge(ctx Context) error {
	for _, v := range e {
		if v == nil {
			continue
		}
		if err := v.SetEdge(ctx); err != nil {
			return err
		}
	}
	return nil
}
