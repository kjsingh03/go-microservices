package types

import (
	"context"
)

type LogServiceInterface interface {
	GetAllLogs(ctx context.Context) ([]Log, error)
	GetLogByID(ctx context.Context, id string) (*Log, error)
	CreateLog(ctx context.Context, req CreateLogRequest) (*Log, error)
	UpdateLog(ctx context.Context, id string, req UpdateLogRequest) (*Log, error)
	DeleteLog(ctx context.Context, id string) error
	DropAllLogs(ctx context.Context) error
	GetLogStats(ctx context.Context) (*LogStats, error)
}

type LogRepositoryInterface interface {
	FindAll(ctx context.Context) ([]Log, error)
	FindByID(ctx context.Context, id string) (*Log, error)
	Create(ctx context.Context, log *Log) error
	Update(ctx context.Context, id string, log *Log) error
	Delete(ctx context.Context, id string) error
	DropCollection(ctx context.Context) error
	GetStats(ctx context.Context) (*LogStats, error)
}