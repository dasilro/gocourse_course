package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("Name is required")
var ErrStartDateRequired = errors.New("Start date is required")
var ErrEndDateRequired = errors.New("End date is required")

type ErrNotFound struct {
	CourseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Course %s doesn't exist", e.CourseID)
}
