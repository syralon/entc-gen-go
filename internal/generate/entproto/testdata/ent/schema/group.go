package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
)

// Group holds the schema definition for the Group entity.
type Group struct {
	ent.Schema
}

// Fields of the Group.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("status").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Immutable(),
	}
}

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("group"),
	}
}

func (Group) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.API(entproto.WithAPIPattern("/v1")),
	}
}
