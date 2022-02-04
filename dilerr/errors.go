package dilerr

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
