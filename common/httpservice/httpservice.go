package httpservice

import (
	"context"

	"gorm.io/gorm"

	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/config"
)

type Service struct {
	cfg         config.KVStore
	postgreConf *gorm.DB
}

func NewService(
	cfg config.KVStore,
	postgreConf *gorm.DB,
) *Service {
	return &Service{
		cfg:         cfg,
		postgreConf: postgreConf,
	}
}

func (s *Service) GetPostgreConf() *gorm.DB {
	return s.postgreConf
}

func (s *Service) GetServiceHealth(_ context.Context) error {
	// do health check logic here
	return nil
}
