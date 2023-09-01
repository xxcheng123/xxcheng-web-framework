package orm

import (
	"reflect"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

func parseModel(entity any) (*model, error) {

	typ := reflect.TypeOf(entity)
	tableName := typ.Name()
	numField := typ.NumField()
	fields := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		fields[fdType.Name] = &field{
			colName: underscoreName(fdType.Name),
		}
	}
	m := &model{
		tableName: underscoreName(tableName),
		fields:    fields,
	}
	return m, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
