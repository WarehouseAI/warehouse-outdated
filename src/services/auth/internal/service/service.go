package service

import (
	"context"
	"warehouse/gen"

	r "warehouse/src/internal/database/redisdb"
	gw "warehouse/src/services/auth/internal/gateway"
	"warehouse/src/services/auth/internal/service/login"
	"warehouse/src/services/auth/internal/service/logout"
	"warehouse/src/services/auth/internal/service/register"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/sirupsen/logrus"
)

type AuthServiceProvider struct {
	ctx             context.Context
	sessionDatabase *r.RedisDatabase
	userGateway     *gw.UserGrpcConnection
	logger          *logrus.Logger
}

func NewAuthService(sessionDatabase *r.RedisDatabase, logger *logrus.Logger) *AuthServiceProvider {
	userGateway := gw.NewUserGrpcConnection("user-service:8001")
	ctx := context.Background()

	return &AuthServiceProvider{
		ctx:             ctx,
		sessionDatabase: sessionDatabase,
		userGateway:     userGateway,
		logger:          logger,
	}
}

func (pvd *AuthServiceProvider) Login(userCreds *m.LoginRequest) (*r.Session, error) {
	return login.Login(userCreds, pvd.userGateway, pvd.sessionDatabase, pvd.logger, pvd.ctx)
}

func (pvd *AuthServiceProvider) Logout(sessionId string) error {
	return logout.Logout(sessionId, pvd.sessionDatabase, pvd.logger, pvd.ctx)
}

func (pvd *AuthServiceProvider) Register(userInfo *gen.CreateUserMsg) (*m.RegisterResponse, error) {
	return register.Register(userInfo, pvd.userGateway, pvd.logger, pvd.ctx)
}
