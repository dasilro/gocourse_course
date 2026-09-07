package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("Name is required")
var ErrStartDateRequired = errors.New("Start date is required")
var ErrStartDateInvalid = errors.New("Start date is invalid")
var ErrEndDateRequired = errors.New("End date is required")
var ErrEndDateInvalid = errors.New("End date is invalid")
var ErrEndLesserStart = errors.New("End date must be bigger or equal than start date")

type ErrNotFound struct {
	CourseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Course %s doesn't exist", e.CourseID)
}
