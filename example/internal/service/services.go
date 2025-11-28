package service

import ent "github.com/syralon/entc-gen-go/example/ent"

type Services struct {
	GroupService *GroupService
	UserService  *UserService
}

func NewServices(client *ent.Client) *Services {
	return &Services{
		NewGroupService(client),
		NewUserService(client),
	}
}
