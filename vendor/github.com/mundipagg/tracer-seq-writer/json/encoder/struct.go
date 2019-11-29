package encoder

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

type Struct struct {
	Type     reflect.Type
	Strategy func(string) string
}

func (changer *Struct) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (changer *Struct) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	beforeBuffer := stream.Buffer()
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "a error occurred while serialization of '%v', error: '%v'", changer.Type.Name(), err)
			stream.SetBuffer(beforeBuffer)
		}
	}()
	v := reflect.NewAt(changer.Type, ptr).Elem()
	switch value := v.Interface().(type) {
	case json.Marshaler:
		valueJ, _ := value.MarshalJSON()
		var valueM interface{}
		_ = json.Unmarshal(valueJ, &valueM)
		stream.WriteVal(valueM)
	case error:
		stream.WriteString(value.Error())
	default:
		stream.WriteObjectStart()
		numFields := v.NumField()
		if numFields > 0 {
			var i int
			for i = 0; i < numFields; i++ {
				fv := v.Field(i)
				ft := changer.Type.Field(i)
				if changer.writeField(ft, fv, stream, false) {
					break
				}
			}
			i++
			for ; i < numFields; i++ {
				fv := v.Field(i)
				ft := changer.Type.Field(i)
				changer.writeField(ft, fv, stream, true)
			}
		}
		stream.WriteObjectEnd()
	}
}

func (changer *Struct) writeField(structField reflect.StructField, value reflect.Value, stream *jsoniter.Stream, needsComma bool) bool {
	if !value.CanInterface() {
		return false
	}

	tag := strings.TrimSpace(structField.Tag.Get("json"))
	if len(tag) == 0 {
		if needsComma {
			stream.WriteMore()
		}
		stream.WriteObjectField(changer.Strategy(structField.Name))
		stream.WriteVal(value.Interface())

	} else {

		pieces := strings.Split(tag, ",")

		if len(pieces) > 1 {
			if pieces[1] == "omitempty" {
				isZero := func() (isZero bool) {
					defer func() {
						if recover() != nil {
							isZero = false
						}
					}()
					return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
				}()
				isNil := func() (isNil bool) {
					defer func() {
						if recover() != nil {
							isNil = false
						}
					}()
					return value.IsNil()
				}()
				if isNil || isZero {
					return false
				}
			}
		}

		if pieces[0] == "-" {
			return false
		}

		if needsComma {
			stream.WriteMore()
		}
		stream.WriteObjectField(changer.Strategy(pieces[0]))
		stream.WriteVal(value.Interface())

	}
	return true
}
