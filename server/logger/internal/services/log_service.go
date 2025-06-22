package services

import (
	"context"
	"errors"
	"logger/types"
	"time"
)

type LogService struct {
	repo types.LogRepositoryInterface
}

func NewLogService(repo types.LogRepositoryInterface) types.LogServiceInterface {
	return &LogService{
		repo: repo,
	}
}

func (s *LogService) GetAllLogs(ctx context.Context) ([]types.Log, error) {
	return s.repo.FindAll(ctx)
}

func (s *LogService) GetLogByID(ctx context.Context, id string) (*types.Log, error) {
	if id == "" {
		return nil, errors.New("log ID is required")
	}
	
	log, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if log == nil {
		return nil, errors.New("log not found")
	}
	
	return log, nil
}

func (s *LogService) CreateLog(ctx context.Context, req types.CreateLogRequest) (*types.Log, error) {
	if req.Name == "" {
		return nil, errors.New("name field is required")
	}
	if req.Data == nil {
		return nil, errors.New("data field is required")
	}

	log := &types.Log{
		Name:      req.Name,
		Data:      req.Data,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return nil, err
	}

	return log, nil
}

func (s *LogService) UpdateLog(ctx context.Context, id string, req types.UpdateLogRequest) (*types.Log, error) {
	if id == "" {
		return nil, errors.New("log ID is required")
	}

	// Check if log exists
	existingLog, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingLog == nil {
		return nil, errors.New("log not found")
	}

	// Update fields if provided
	updatedLog := &types.Log{
		ID:        existingLog.ID,
		Name:      existingLog.Name,
		Data:      existingLog.Data,
		CreatedAt: existingLog.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}

	if req.Name != "" {
		updatedLog.Name = req.Name
	}
	if req.Data != nil {
		updatedLog.Data = req.Data
	}

	if err := s.repo.Update(ctx, id, updatedLog); err != nil {
		return nil, err
	}

	return updatedLog, nil
}

func (s *LogService) DeleteLog(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("log ID is required")
	}

	// Check if log exists
	existingLog, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingLog == nil {
		return errors.New("log not found")
	}

	return s.repo.Delete(ctx, id)
}

func (s *LogService) DropAllLogs(ctx context.Context) error {
	return s.repo.DropCollection(ctx)
}

func (s *LogService) GetLogStats(ctx context.Context) (*types.LogStats, error) {
	return s.repo.GetStats(ctx)
}