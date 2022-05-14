package dilema_test

import (
	"testing"

	"github.com/kjushka/dilema"
	"github.com/kjushka/dilema/test_data/service"

	"github.com/stretchr/testify/assert"
)

func TestRegisterSingleTone(t *testing.T) {
	di := dilema.Init()
	
	constructor := service.NewSomeActionByWithParams
	a, b := 1, 2
	strA := "1"
	strB := "2"

	err := di.RegisterSingletone("some_action", constructor, strA, b)
	assert.NotNil(t, err)
	err = di.RegisterSingletone("some_action", constructor, a, strB)
	assert.NotNil(t, err)
	err = di.RegisterSingletone("some_action", constructor, a, b)
	assert.Nil(t, err)
	err = di.RegisterSingletone("some_action", constructor, a, b)
	assert.NotNil(t, err)
}
