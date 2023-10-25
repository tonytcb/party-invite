package domain

import "fmt"

type ErrInvalidArgument struct {
	OriginalErr string
	Description string
}

func NewErrInvalidArgument(originalErr string, description string) *ErrInvalidArgument {
	return &ErrInvalidArgument{OriginalErr: originalErr, Description: description}
}

func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("%s: %s", e.Description, e.OriginalErr)
}
