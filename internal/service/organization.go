package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"QuickSlot/pkg/database/redis"
)

type OrganizationService struct {
	repo        repository.OrganizationRepository
	redisClient *redis.RedisClient
}

func NewOrganizationService(r repository.OrganizationRepository, rc *redis.RedisClient) *OrganizationService {
	return &OrganizationService{repo: r, redisClient: rc}
}

func (s *OrganizationService) invalidateCache(ctx context.Context, id int64) {
	if s.redisClient != nil && s.redisClient.Client != nil {
		s.redisClient.Client.Del(ctx, "cache:org:all")
		if id > 0 {
			s.redisClient.Client.Del(ctx, fmt.Sprintf("cache:org:%d", id))
		}
	}
}

func (s *OrganizationService) Create(ctx context.Context, name string, ownerID int64) (int64, error) {
	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, repository.ErrConflict
	}

	org := &model.Organization{
		Name:    name,
		OwnerID: ownerID,
	}
	id, err := s.repo.CreateOrg(ctx, org)
	if err == nil {
		s.invalidateCache(ctx, id)
	}
	return id, err
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]model.Organization, error) {
	if s.redisClient != nil && s.redisClient.Client != nil {
		cached, err := s.redisClient.Client.Get(ctx, "cache:org:all").Result()
		if err == nil {
			var orgs []model.Organization
			if json.Unmarshal([]byte(cached), &orgs) == nil {
				return orgs, nil
			}
		}
	}

	orgs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if s.redisClient != nil && s.redisClient.Client != nil {
		if data, err := json.Marshal(orgs); err == nil {
			s.redisClient.Client.Set(ctx, "cache:org:all", data, 10*time.Minute)
		}
	}

	return orgs, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id int64) (*model.Organization, error) {
	key := fmt.Sprintf("cache:org:%d", id)
	if s.redisClient != nil && s.redisClient.Client != nil {
		cached, err := s.redisClient.Client.Get(ctx, key).Result()
		if err == nil {
			var org model.Organization
			if json.Unmarshal([]byte(cached), &org) == nil {
				return &org, nil
			}
		}
	}

	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.redisClient != nil && s.redisClient.Client != nil {
		if data, err := json.Marshal(org); err == nil {
			s.redisClient.Client.Set(ctx, key, data, 10*time.Minute)
		}
	}

	return org, nil
}

func (s *OrganizationService) Update(ctx context.Context, id int64, update *model.OrganizationUpdate) error {
	err := s.repo.Update(ctx, id, update)
	if err == nil {
		s.invalidateCache(ctx, id)
	}
	return err
}

func (s *OrganizationService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.invalidateCache(ctx, id)
	}
	return err
}
