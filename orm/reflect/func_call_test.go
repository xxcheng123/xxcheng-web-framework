package reflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestUser struct {
	Name string
	age  int
}

func NewTestUser(name string, age int) TestUser {
	return TestUser{
		Name: name,
		age:  age,
	}
}
func NewTestUserPtr(name string, age int) *TestUser {
	return &TestUser{
		Name: name,
		age:  age,
	}
}
func (t TestUser) GetName() string {
	return t.Name
}
func (t *TestUser) SetName(name string) {
	t.Name = name
}
func (t TestUser) getAge() int {
	return t.age
}

func TestIterateFunc(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		want    map[string]FuncInfo
		wantErr error
	}{
		{
			name:    "common",
			entity:  NewTestUser("xxcheng", 18),
			wantErr: nil,
			want: map[string]FuncInfo{
				"GetName": {
					Name: "GetName",
					InputTypes: []reflect.Type{
						reflect.TypeOf(TestUser{}),
					},
					OutputTypes: []reflect.Type{
						reflect.TypeOf(""),
					},
					Result: []any{
						"xxcheng",
					},
				},
			},
		},
		{
			name:    "ptr",
			entity:  NewTestUserPtr("xxcheng", 18),
			wantErr: nil,
			want: map[string]FuncInfo{
				"GetName": {
					Name: "GetName",
					InputTypes: []reflect.Type{
						reflect.TypeOf(&TestUser{}),
					},
					OutputTypes: []reflect.Type{
						reflect.TypeOf(""),
					},
					Result: []any{
						"xxcheng",
					},
				},
				"SetName": {
					Name: "SetName",
					InputTypes: []reflect.Type{
						reflect.TypeOf(&TestUser{}),
						reflect.TypeOf(""),
					},
					OutputTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFunc(tt.entity)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
