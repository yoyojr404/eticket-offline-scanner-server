package tickets

import (
	"context"

	"github.com/yoyojr404/pkg/apperror"
)

type Service interface {
	GetAll(ctx context.Context, page, perPage int) (*TicketResult, error)
	GetByID(ctx context.Context, id int64) (*Ticket, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (s *serviceImpl) GetAll(ctx context.Context, page, perPage int) (*TicketResult, error) {
	offset := (page - 1) * perPage
	entities, total, err := s.repo.FindAll(ctx, offset, perPage)
	if err != nil {
		return nil, apperror.NewInternal("Failed to retrieve tickets", err)
	}

	return &TicketResult{
		Entities:     entities,
		CurrentPage:  page,
		TotalPerPage: perPage,
		Total:        int64(total),
	}, nil
}

func (s *serviceImpl) GetByID(ctx context.Context, id int64) (*Ticket, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal("Failed to fetch ticket", err)
	}

	if entity == nil {
		return nil, apperror.NewNotFound("Ticket not found", ErrTicketNotFound)
	}

	return entity, nil
}
