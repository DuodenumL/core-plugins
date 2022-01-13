package command

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func call(plugin interface{}, method string, args ...interface{}) ([]interface{}, error) {
	pluginV := reflect.ValueOf(plugin)
	funcV := pluginV.MethodByName(method)
	if funcV.IsZero() {
		return nil, errors.Wrap(ErrNoSuchMethod, fmt.Sprintf("method name: %s", method))
	}
	in := []reflect.Value{}
	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}
	out := funcV.Call(in)
	res := []interface{}{}
	for _, elem := range out {
		res = append(res, elem.Interface())
	}
	return res, nil
}

func printResult(result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		logrus.Errorf("[printResult] failed to marshal result %+v, err: %v", result, err)
		fmt.Println("invalid json")
	}

	fmt.Println(string(body))
}

func newPtrOfType(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return reflect.New(reflect.Indirect(reflect.ValueOf(v)).Type()).Interface()
}

func makeSliceOfType(v interface{}, cap int) reflect.Value {
	return reflect.MakeSlice(reflect.TypeOf(v), 0, cap)
}

func makeMapOfType(v interface{}) reflect.Value {
	return reflect.MakeMap(reflect.TypeOf(v))
}
