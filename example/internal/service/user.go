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

func UserToProto(data *ent.User) *pb.User {
	return &pb.User{
		Name:      data.Name,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
		GroupId:   data.GroupID,
		Status:    data.Status,
	}
}

func UserFromProto(data *pb.User) *ent.User {
	return &ent.User{
		Name:      data.Name,
		CreatedAt: data.CreatedAt.AsTime(),
		UpdatedAt: data.UpdatedAt.AsTime(),
		GroupID:   data.GroupId,
		Status:    data.Status,
	}
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	client *ent.UserClient
}

func NewUserService(client *ent.Client) *UserService {
	return &UserService{
		client: client.User,
	}
}

func (s *UserService) Get(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	data, err := s.client.Get(ctx, int(request.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Data: UserToProto(data),
	}, nil
}

func (s *UserService) List(ctx context.Context, request *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	conditions := entproto.Selectors[predicate.User](request.Options.Name.Selector(user.FieldName), request.Options.CreatedAt.Selector(user.FieldCreatedAt), request.Options.UpdatedAt.Selector(user.FieldUpdatedAt), request.Options.GroupId.Selector(user.FieldGroupID), request.Options.Status.Selector(user.FieldStatus))
	query := s.client.Query()
	query = query.Where(conditions...)

	if e := request.Options.UserGroups; e != nil {
		query.WithUserGroups(func(eq *ent.GroupQuery) {
			eq.Where(entproto.Selectors[predicate.Group](e.Name.Selector(group.FieldName), e.CreatedAt.Selector(group.FieldCreatedAt), e.UpdatedAt.Selector(group.FieldUpdatedAt))...)
		})
	}

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
