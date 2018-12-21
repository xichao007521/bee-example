package services

import "do-global.com/bee-example/models"

type UserService struct {
	Service
}

func (*UserService) Login(username string, password string) *models.User {
	// TODO do something
	return &models.User{
		Id:   1,
		Name: "u1",
	}
}
