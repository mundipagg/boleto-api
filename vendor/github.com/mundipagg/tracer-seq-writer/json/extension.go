package json

import (
	"github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/mundipagg/tracer-seq-writer/json/encoder"
	"reflect"
)

type CaseStrategyExtension struct {
	jsoniter.DummyExtension
	Strategy func(string) string
}

func (cs *CaseStrategyExtension) CreateMapKeyEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	if typ.Kind() == reflect.String {
		return &encoder.Map{
			Strategy: cs.Strategy,
		}
	}
	return nil
}

func (cs *CaseStrategyExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	ty := typ.Type1()
	if ty.Kind() == reflect.Struct {
		return &encoder.Struct{
			Type:     ty,
			Strategy: cs.Strategy,
		}
	}
	return nil
}
