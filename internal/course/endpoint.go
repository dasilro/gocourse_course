package course

import (
	"context"
	"errors"

	"github.com/dasilro/go_lib_response/response"
	"github.com/dasilro/gocourse_meta/meta"
)

type (
	//Controller func(w http.ResponseWriter, r *http.Request)
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	UpdateReq struct {
		ID        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"err,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service, paginatorLimitDefault string) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, paginatorLimitDefault),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.Name == "" {

			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {

			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate == "" {

			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		course, err := s.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("success", course, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		course, err := s.Get(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil), nil
	}
}

func makeGetAllEndpoint(s Service, paginatorLimitDefault string) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)
		filters := Filters{
			Name: req.Name,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, paginatorLimitDefault)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", courses, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil, nil), nil
	}
}
