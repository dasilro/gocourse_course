package course

import (
	"context"
	"log"
	"time"

	"github.com/dasilro/gocourse_domain/domain"
)

type (
	Filters struct {
		Name string
	}

	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset int, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{log: log, repo: repo}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {
	startDateParsed, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return nil, err
	}

	endDateParsed, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return nil, err
	}

	course := domain.Course{Name: name, StartDate: startDateParsed, EndDate: endDateParsed}
	if err := s.repo.Create(ctx, &course); err != nil {
		s.log.Println(err.Error())
		return nil, err
	}
	return &course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset int, limit int) ([]domain.Course, error) {
	return s.repo.GetAll(ctx, filters, offset, limit)
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	return s.repo.Get(ctx, id)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s service) Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error {
	var startDateParsed *time.Time
	if startDate != nil {
		parsed, err := time.Parse(time.DateOnly, *startDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		startDateParsed = &parsed
	}

	var endDateParsed *time.Time
	if endDate != nil {
		parsed, err := time.Parse(time.DateOnly, *endDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		endDateParsed = &parsed
	}

	return s.repo.Update(ctx, id, name, startDateParsed, endDateParsed)
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
