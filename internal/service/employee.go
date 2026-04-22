package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"QuickSlot/pkg/database/redis"
)

type EmployeeService struct {
	empRepo     repository.EmployeeRepository
	orgRepo     repository.OrganizationRepository
	redisClient *redis.RedisClient
}

func NewEmployeeService(e repository.EmployeeRepository, o repository.OrganizationRepository, rc *redis.RedisClient) *EmployeeService {
	return &EmployeeService{empRepo: e, orgRepo: o, redisClient: rc}
}

func (s *EmployeeService) invalidateCache(ctx context.Context, id int64, orgID int64) {
	if s.redisClient != nil && s.redisClient.Client != nil {
		if id > 0 {
			s.redisClient.Client.Del(ctx, fmt.Sprintf("cache:emp:%d", id))
		}
		if orgID > 0 {
			s.redisClient.Client.Del(ctx, fmt.Sprintf("cache:emp:org:%d", orgID))
		}
	}
}

func (s *EmployeeService) Create(ctx context.Context, emp *model.Employee) (int64, error) {
	exists, err := s.empRepo.ExistsByEmailOrPhone(ctx, emp.Email, emp.Phone)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, repository.ErrConflict
	}

	_, err = s.orgRepo.GetByID(ctx, emp.OrganizationID)
	if err != nil {
		if errors.Is(err, repository.ErrOrgNotFound) {
			return 0, errors.New("organization does not exist or is deleted")
		}
		return 0, err
	}

	id, err := s.empRepo.CreateEmployee(ctx, emp)
	if err == nil {
		s.invalidateCache(ctx, id, emp.OrganizationID)
	}
	return id, err
}

func (s *EmployeeService) GetByOrganization(ctx context.Context, orgID int64) ([]model.Employee, error) {
	key := fmt.Sprintf("cache:emp:org:%d", orgID)
	if s.redisClient != nil && s.redisClient.Client != nil {
		cached, err := s.redisClient.Client.Get(ctx, key).Result()
		if err == nil {
			var employees []model.Employee
			if json.Unmarshal([]byte(cached), &employees) == nil {
				return employees, nil
			}
		}
	}

	employees, err := s.empRepo.GetByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if s.redisClient != nil && s.redisClient.Client != nil {
		if data, err := json.Marshal(employees); err == nil {
			s.redisClient.Client.Set(ctx, key, data, 10*time.Minute)
		}
	}

	return employees, nil
}

func (s *EmployeeService) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	key := fmt.Sprintf("cache:emp:%d", id)
	if s.redisClient != nil && s.redisClient.Client != nil {
		cached, err := s.redisClient.Client.Get(ctx, key).Result()
		if err == nil {
			var emp model.Employee
			if json.Unmarshal([]byte(cached), &emp) == nil {
				return &emp, nil
			}
		}
	}

	emp, err := s.empRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.redisClient != nil && s.redisClient.Client != nil {
		if data, err := json.Marshal(emp); err == nil {
			s.redisClient.Client.Set(ctx, key, data, 10*time.Minute)
		}
	}

	return emp, nil
}

func (s *EmployeeService) Update(ctx context.Context, id int64, update *model.EmployeeUpdate) error {
	if update.OrganizationID != nil {
		_, err := s.orgRepo.GetByID(ctx, *update.OrganizationID)
		if err != nil {
			if errors.Is(err, repository.ErrOrgNotFound) {
				return errors.New("organization does not exist or is deleted")
			}
			return err
		}
	}
	err := s.empRepo.Update(ctx, id, update)
	
	if err == nil {
		var orgID int64
		if update.OrganizationID != nil {
			orgID = *update.OrganizationID
		} else {
			if e, eErr := s.empRepo.GetByID(ctx, id); eErr == nil {
				orgID = e.OrganizationID
			}
		}
		s.invalidateCache(ctx, id, orgID)
	}

	return err
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	var orgID int64
	if e, err := s.empRepo.GetByID(ctx, id); err == nil {
		orgID = e.OrganizationID
	}

	err := s.empRepo.Delete(ctx, id)
	if err == nil {
		s.invalidateCache(ctx, id, orgID)
	}
	return err
}
