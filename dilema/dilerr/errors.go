package dilerr

// CreationError for messaging about error in creation process
type CreationError struct {
	text string
}

func NewCreationError(text string) *CreationError {
	return &CreationError{text}
}

func (pe *CreationError) Error() string {
	return pe.text
}

// TypeError for messaging about error in type checking
type TypeError struct {
	text string
}

func NewTypeError(text string) *TypeError {
	return &TypeError{text}
}

func (te *TypeError) Error() string {
	return te.text
}
