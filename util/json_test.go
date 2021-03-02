package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertObjectToJSON(t *testing.T) {
	type T struct {
		Field string
	}
	obj := new(T)
	obj.Field = "A"
	obj2 := new(T)
	err := FromJSON(ToJSON(obj), obj2)

	assert.Nil(t, err, "Deve encriptar o texto")
	assert.Equal(t, obj2.Field, obj.Field, "Deve-se converter um objeto para JSON e vice-versa")
}
