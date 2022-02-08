package dilerr

import "fmt"

// creationError for messaging about error in creation process
type creationError struct {
	text string
}

func NewCreationError(text string) error {
	return &creationError{text}
}

func (ce *creationError) Error() string {
	return ce.text
}

func GetAlreadyExistError(alias string) error {
	return NewCreationError(
		fmt.Sprintf(
			"container with alias '%s' already exists", alias,
		),
	)
}

// typeError for messaging about error in type checking
type typeError struct {
	text string
}

func NewTypeError(text string) error {
	return &typeError{text}
}

func (te *typeError) Error() string {
	return te.text
}

// TypeError for messaging about error in type checking
type getError struct {
	text string
}

func NewGetError(text string) error {
	return &getError{text}
}

func (ge *getError) Error() string {
	return ge.text
}

type threadError struct {
	text string
}

func NewThreadError(text string) error {
	return &threadError{text}
}

func (te *threadError) Error() string {
	return te.text
}
