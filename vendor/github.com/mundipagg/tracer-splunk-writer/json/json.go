package json

import (
	"github.com/json-iterator/go"
)

func New() jsoniter.API {
	return NewWithCaseStrategy(func(s string) string {
		return s
	})
}

func NewWithCaseStrategy(strategy func(string) string) jsoniter.API {
	json := jsoniter.Config{
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       false,
		ObjectFieldMustBeSimpleString: true,
	}.Froze()
	json.RegisterExtension(&CaseStrategyExtension{
		Strategy: strategy,
	})
	return json
}
