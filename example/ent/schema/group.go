package schema

import (
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
		field.Time("created_at").Annotations(entproto.Field(entproto.WithFieldImmutable(true), entproto.WithFieldOrderable(true))),
		field.Time("updated_at").Annotations(entproto.Field(entproto.WithFieldImmutable(true), entproto.WithFieldOrderable(true))),
	}
}

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group_users", User.Type).
			Ref("user_groups"),
		// Through("groups", UserGroup.Type), // 多对多关联也可以手动指定关联表
	}
}

func (Group) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.API(entproto.WithAPIPattern("/v1")),
	}
}
