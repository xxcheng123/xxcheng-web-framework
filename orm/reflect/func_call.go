package reflect

import "reflect"

/*
具体的一个值->Typeof Type
具体的一个值->Typeof Value

Type
Value
	Type ->Type
Method
	Func ->Value

Type{
	Field
	Method
}

**/

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	//传入的类型的上面有哪些方法
	typ := reflect.TypeOf(entity)
	//获取每个方法
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		//方法 Method
		method := typ.Method(i)
		//方法的值 Value
		fn := method.Func
		// Type
		fn.Type()
		numIn := fn.Type().NumIn()
		input := make([]reflect.Type, 0, numIn)
		inputValues := make([]reflect.Value, 0, numIn)
		//第一个是对象本身
		input = append(input, typ)
		inputValues = append(inputValues, reflect.ValueOf(entity))
		for j := 1; j < numIn; j++ {
			input = append(input, fn.Type().In(j))
			inputValues = append(inputValues, reflect.Zero(fn.Type().In(i)))
		}
		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			output = append(output, fn.Type().Out(j))
		}
		resValues := fn.Call(inputValues)
		result := make([]any, 0, len(resValues))
		for _, v := range resValues {
			result = append(result, v.Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
	}
	return res, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
