package adaptor

import "project/internal/action"

type secondService struct {
	someField int
}

func NewSecondService(sf int) action.SecondService {
	return &secondService{sf}
}

func (ss *secondService) MethodThree() string {
	return ""
}
