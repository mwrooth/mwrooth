package service

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mwrooth/gateway/pkg/api/general"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Service struct {
	*general.Resources
	*general.ServiceGeneral
}

func NewLogger(logLevel string) (*zap.SugaredLogger, error) {
	var cfg zap.Config

	switch os.Getenv("ENVIRONMENT") {
	case "PRODUCTION":
		cfg = zap.NewProductionConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.Level.SetLevel(zapcore.DebugLevel)

	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(log)

	return log.Sugar(), err
}

func New() *Service {
	// Setup logger
	logger, err := NewLogger(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Errorf("failed to init zap logger: %+v", err))
	}

	s := &Service{
		Resources: &general.Resources{
			Log: logger,
			GrpcClients: map[string]*general.GrpcClient{
				"MWROOTH":   {},
				"DETECTION": {},
			},
		},
	}

	s.Engine = s.NewEngine()
	return s
}

func (s *Service) GetGrpcClients() map[string]*general.GrpcClient {
	return s.GrpcClients
}

func (s *Service) GetEngine() *gin.Engine {
	return s.Engine
}

func (s *Service) SetBaseConfig(cfg *general.BaseConfig) {
	s.BaseCfg = cfg
}

func (s *Service) GetLog() *zap.SugaredLogger {
	return s.Log
}
