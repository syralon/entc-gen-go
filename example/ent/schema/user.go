package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	openapiv3 "github.com/google/gnostic/openapiv3"
	"github.com/syralon/entc-gen-go/pkg/annotations/entproto"
	"github.com/syralon/entc-gen-go/pkg/annotations/openapi"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("The user name").Annotations(openapi.Schema(&openapiv3.Schema{Description: "The name"})),
		field.Time("created_at").Comment("The created time").Annotations(entproto.Field(entproto.WithFieldImmutable(true), entproto.WithFieldOrderable(true))),
		field.Time("updated_at").Comment("The latest update time").Annotations(entproto.Field(entproto.WithFieldImmutable(true))),
		field.Int64("group_id").Annotations(entproto.Field(entproto.WithFieldOrderable(true), entproto.WithFieldSettable(true))),
		field.Int32("status").Annotations(entproto.Field(entproto.WithFieldSettable(true))),
	}
}

// Edges of the User.
// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_groups", Group.Type),
		// Through("groups", UserGroup.Type), // 多对多关联也可以手动指定关联表

		// edge.To("members", User.Type).
		// 	From("leader").
		// 	Field("leader_id").
		// 	Unique(),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.API(entproto.WithAPIPattern("/v1")),
	}
}
