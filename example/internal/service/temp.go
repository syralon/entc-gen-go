package service

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/syralon/entc-gen-go/example/ent"
	"github.com/syralon/entc-gen-go/example/ent/group"
	"github.com/syralon/entc-gen-go/example/ent/predicate"
	"github.com/syralon/entc-gen-go/example/ent/user"
	pb "github.com/syralon/entc-gen-go/example/proto/example"
	"github.com/syralon/entc-gen-go/proto/syralon/entproto"
)

func (s *GroupService) list(ctx context.Context, req *pb.ListGroupRequest) (*pb.ListGroupResponse, error) {
	conditions := entproto.Selectors[predicate.Group](
		req.GetOptions().GetId().Selector(group.FieldID),
		req.GetOptions().GetName().Selector(group.FieldName),
		req.GetOptions().GetCreatedAt().Selector(group.FieldCreatedAt),
		req.GetOptions().GetUpdatedAt().Selector(group.FieldUpdatedAt),
	)
	query := s.client.Query()
	query = query.Where(conditions...)

	if req.GetOptions().GetGroupUsers() != nil {
		u := req.GetOptions().GetGroupUsers()
		query.WithGroupUsers(func(uq *ent.UserQuery) {
			uq.Where(entproto.Selectors[predicate.User](
				u.GetId().Selector(user.FieldID),
				u.GetName().Selector(user.FieldName),
				u.GetCreatedAt().Selector(user.FieldCreatedAt),
				u.GetUpdatedAt().Selector(user.FieldUpdatedAt),
				u.GetGroupId().Selector(user.FieldGroupID),
				u.GetStatus().Selector(user.FieldStatus),
			)...)
		})
	}
	if req.GetPaginator() != nil {
		switch page := req.GetPaginator().GetPaginator().(type) {
		case *entproto.Paginator_Classical:
			query = query.
				Order(page.Classical.OrderSelector()).
				Offset(int(page.Classical.GetLimit() * (page.Classical.GetPage() - 1))).
				Limit(int(page.Classical.GetLimit()))
		case *entproto.Paginator_Infinite:
			query = query.
				Order(group.ByID(sql.OrderDesc())).
				Limit(int(page.Infinite.GetLimit()))
			if page.Infinite.GetSequence() > 0 {
				query = query.Where(group.IDLT(int(page.Infinite.GetSequence())))
			}
		default:
		}
	}
	groups, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListGroupResponse{Data: Trans(groups, GroupToProto)}, nil
}
