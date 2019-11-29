package seq

import (
	"fmt"
	"reflect"
)

type Entry map[string]interface{}

func (log Entry) Add(name string, value interface{}) Entry {
	log[name] = value
	return log
}

func NewEntry(p ...interface{}) Entry {
	if len(p) == 1 {
		c, ok := p[0].([]interface{})
		if ok {
			p = c
		}
	}
	normalized := Entry{}
	for _, item := range p {
		if item == nil {
			continue
		}
		itemType := reflect.TypeOf(item)
		v := reflect.ValueOf(item)
		inner := Entry{}
		switch itemType.Kind() {
		case reflect.Map:
			for _, key := range v.MapKeys() {
				inner[fmt.Sprint(key.Interface())] = v.MapIndex(key).Interface()
			}
		case reflect.Ptr, reflect.Interface:
			if !v.IsNil() {
				inner = NewEntry(v.Elem().Interface())
			}
		default:
			inner[itemType.Name()] = item
		}
		normalized = Merge(normalized, inner)
	}
	return normalized
}

func Merge(np Entry, other Entry) Entry {
	collisions := map[string]int{}
	r := Entry{}
	for key, value := range np {
		r[key] = value
	}
	for key, value := range other {
		if _, ok := r[key]; ok {
			var index int
			if index, ok = collisions[key]; !ok {
				index = 0
			}
			index += 1
			collisions[key] = index
			key = fmt.Sprintf("%s%d", key, index)
		}
		r[key] = value
	}
	return r
}
