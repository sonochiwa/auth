package services

import "go.uber.org/zap"

type UserService struct {
	log *zap.Logger
}

func NewUserService(logger *zap.Logger) *UserService {
	return &UserService{
		log: logger,
	}
}
