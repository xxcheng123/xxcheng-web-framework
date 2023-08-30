package reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("空的")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持的类型")
	}

	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)
		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}
	return res, nil
}

func SetField(entity any, field string, newValue any) error {
	fieldVal := reflect.ValueOf(entity)
	/**

	为什么传入结构体，然后使用这种方法不行？
	因为传进来的 entity 类型是一个 any，就是 interface{}
	而 NumField、FieldByName 这些方法只能使用在 struct 结构体上面。
	但是为什么传入一个指针就可以了呢
	if fieldVal.Kind() == reflect.Struct {
		fieldVal = reflect.ValueOf(&entity)
		fieldTyp := reflect.TypeOf(&entity)
		for fieldTyp.Kind() == reflect.Pointer {
			fieldVal = fieldVal.Elem()
			fieldTyp = fieldTyp.Elem()
		}
		numField := fieldTyp.NumField()
		for i := 0; i < numField; i++ {
			if fieldTyp.Field(i).Name == field && fieldVal.Field(i).CanSet() {
				fieldVal.Field(i).Set(reflect.ValueOf(newValue))
				return nil
			}
			return errors.New("不允许更改")
		}
	}
	*/
	for fieldVal.Kind() == reflect.Pointer {
		fieldVal = fieldVal.Elem()
	}
	var val reflect.Value
	if field == "" {
		val = fieldVal
	} else {
		val = fieldVal.FieldByName(field)
	}
	if val.CanSet() {
		newVal := reflect.ValueOf(newValue)
		val.Set(newVal)
		return nil
	}
	return errors.New("不支持更改")
}
