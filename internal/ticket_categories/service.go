package ticket_categories

import (
	"context"

	"github.com/yoyojr404/pkg/apperror"
)

type Service interface {
	GetAll(ctx context.Context, page, limit int) (*TicketCategoryResult, error)
	GetByID(ctx context.Context, id int64) (*CategoryResponse, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) GetAll(ctx context.Context, page, limit int) (*TicketCategoryResult, error) {
	offset := (page - 1) * limit

	entities, total, err := s.repo.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, apperror.NewInternal("Failed to retrieve ticket categories", err)
	}

	return &TicketCategoryResult{
		Entities:     entities,
		CurrentPage:  page,
		TotalPerPage: limit,
		Total:        int64(total),
	}, nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id int64) (*CategoryResponse, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal("Failed to fetch ticket category", err)
	}

	if entity == nil {
		return nil, apperror.NewNotFound("Ticket category not found", ErrCategoryNotFound)
	}

	res := ToResponse(entity)
	return &res, nil
}
