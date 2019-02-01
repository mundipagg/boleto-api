package tmpl

import (
	"github.com/mundipagg/boleto-api/util"
	. "github.com/PMoneda/flow"
)

//Using flow to mapping data in JSON from a contract to other
func TransformFromJSON(content string, from string, to string, obj interface{}) (interface{}) {
	var data Transform = Transform(content)
	var input Transform = Transform(from)
	var output Transform = Transform(to)

	success, err := data.TransformFromJSON(input, output, GetFuncMaps())

	if err != nil{
		return err
	}

	if (obj == nil){
		return success
	}

	result := util.ParseJSON(success, obj)
	return result
}

//Using flow to mapping data in XML from a contract to other
func TransformFromXML(content string, from string, to string, obj interface{}) (interface{}) {
	var data Transform = Transform(content)
	var input Transform = Transform(from)
	var output Transform = Transform(to)

	success, err := data.TransformFromXML(input, output, GetFuncMaps())

	if err != nil{
		return err
	}

	if (obj == nil){
		return success
	}

	result := util.ParseJSON(success, obj)
	return result
}