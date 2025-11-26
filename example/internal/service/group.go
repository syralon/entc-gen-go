package service

import (
	"context"
	ent "github.com/syralon/entc-gen-go/example/ent"
	group "github.com/syralon/entc-gen-go/example/ent/group"
	predicate "github.com/syralon/entc-gen-go/example/ent/predicate"
	user "github.com/syralon/entc-gen-go/example/ent/user"
	pb "github.com/syralon/entc-gen-go/example/proto/example"
	entproto "github.com/syralon/entc-gen-go/proto/syralon/entproto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func GroupToProto(data *ent.Group) *pb.Group {
	return &pb.Group{
		Name:      data.Name,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}
}

func GroupFromProto(data *pb.Group) *ent.Group {
	return &ent.Group{
		Name:      data.Name,
		CreatedAt: data.CreatedAt.AsTime(),
		UpdatedAt: data.UpdatedAt.AsTime(),
	}
}

type GroupService struct {
	pb.UnimplementedGroupServiceServer
	client *ent.GroupClient
}

func NewGroupService(client *ent.Client) *GroupService {
	return &GroupService{
		client: client.Group,
	}
}

func (s *GroupService) Get(ctx context.Context, request *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	data, err := s.client.Get(ctx, int(request.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.GetGroupResponse{
		Data: GroupToProto(data),
	}, nil
}

func (s *GroupService) List(ctx context.Context, request *pb.ListGroupRequest) (*pb.ListGroupResponse, error) {
	conditions := entproto.Selectors[predicate.Group](request.Options.Name.Selector(group.FieldName), request.Options.CreatedAt.Selector(group.FieldCreatedAt), request.Options.UpdatedAt.Selector(group.FieldUpdatedAt))
	query := s.client.Query()
	query = query.Where(conditions...)

	if e := request.Options.GroupUsers; e != nil {
		query.WithGroupUsers(func(eq *ent.UserQuery) {
			eq.Where(entproto.Selectors[predicate.User](e.Name.Selector(user.FieldName), e.CreatedAt.Selector(user.FieldCreatedAt), e.UpdatedAt.Selector(user.FieldUpdatedAt), e.GroupId.Selector(user.FieldGroupID), e.Status.Selector(user.FieldStatus))...)
		})
	}

	if paginator := request.GetPaginator(); paginator != nil {
		switch page := paginator.GetPaginator().(type) {
		case *entproto.Paginator_Classical:
			query = query.Order(page.Classical.OrderSelector()).Offset(int(page.Classical.GetLimit() * (page.Classical.GetPage() - 1))).Limit(int(page.Classical.GetLimit()))
		case *entproto.Paginator_Infinite:
			query = query.Order(group.ByID()).Limit(int(page.Infinite.GetLimit()))
			if sequence := page.Infinite.GetSequence(); sequence > 0 {
				query = query.Where(group.IDLT(int(page.Infinite.GetSequence())))
			}
		}
	}

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListGroupResponse{
		Data: Trans(data, GroupToProto),
	}, nil
}

func (s *GroupService) ListGroupUsers(ctx context.Context, request *pb.ListGroupGroupUsersRequest) (*pb.ListUserResponse, error) {
	query := s.client.Query().Where(group.ID(int(request.GroupId))).QueryGroupUsers().Where(entproto.Selectors[predicate.User](
		request.Options.Name.Selector(user.FieldName),
		request.Options.CreatedAt.Selector(user.FieldCreatedAt),
		request.Options.UpdatedAt.Selector(user.FieldUpdatedAt),
		request.Options.GroupId.Selector(user.FieldGroupID),
		request.Options.Status.Selector(user.FieldStatus))...)

	if paginator := request.GetPaginator(); paginator != nil {
		switch page := paginator.GetPaginator().(type) {
		case *entproto.Paginator_Classical:
			query = query.Order(page.Classical.OrderSelector()).Offset(int(page.Classical.GetLimit() * (page.Classical.GetPage() - 1))).Limit(int(page.Classical.GetLimit()))
		case *entproto.Paginator_Infinite:
			query = query.Order(user.ByID()).Limit(int(page.Infinite.GetLimit()))
			if sequence := page.Infinite.GetSequence(); sequence > 0 {
				query = query.Where(user.IDLT(int(page.Infinite.GetSequence())))
			}
		}
	}
	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.ListUserResponse{
		Data: Trans(data, UserToProto),
	}, nil
}
