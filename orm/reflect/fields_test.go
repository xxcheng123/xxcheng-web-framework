package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFields(t *testing.T) {
	type User struct {
		Name string
		age  int
	}
	tests := []struct {
		name    string
		entity  any
		want    map[string]any
		wantErr error
	}{
		{
			name: "common",
			entity: User{
				Name: "xxcheng",
				age:  18,
			},
			wantErr: nil,
			want: map[string]any{
				"Name": "xxcheng",
				"age":  0,
			},
		},
		{
			name: "ptr",
			entity: &User{
				Name: "xxcheng",
				age:  18,
			},
			wantErr: nil,
			want: map[string]any{
				"Name": "xxcheng",
				"age":  0,
			},
		},
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("空的"),
		},
		{
			name:    "nil ptr",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持零值"),
		},
		{
			name:    "common type",
			entity:  19,
			wantErr: errors.New("不支持的类型"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IterateFields(tt.entity)
			if err == nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSetField(t *testing.T) {
	type User struct {
		Name string
		age  int
	}
	tests := []struct {
		name       string
		entity     any
		wantErr    error
		wantField  string
		wantValue  any
		wantResult any
	}{
		{
			name: "common",
			entity: User{
				Name: "xxcheng",
			},
			wantErr:   errors.New("不支持更改"),
			wantField: "Name",
			wantValue: "jpc",
		},
		{
			name: "common",
			entity: &User{
				Name: "xxcheng",
			},
			wantErr:   nil,
			wantField: "Name",
			wantValue: "jpc",
			wantResult: &User{
				Name: "jpc",
			},
		},
		{
			name:      "common num",
			entity:    88,
			wantErr:   errors.New("不支持更改"),
			wantValue: 99,
		},
		{
			name: "ptr num",
			entity: func() *int {
				num := 88
				return &num
			}(),
			wantErr:   nil,
			wantValue: 99,
			wantResult: func() *int {
				num := 99
				return &num
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetField(tt.entity, tt.wantField, tt.wantValue)
			if err != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.wantResult, tt.entity)
		})
	}
}
