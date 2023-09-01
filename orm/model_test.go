package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestTable struct {
	ID        int64
	Username  string
	UserGrade int
	Age       int
}

func Test_parseModel(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantErr error
		want    *model
	}{
		{
			name: "common",
			entity: TestTable{
				ID:        1,
				Username:  "xxcheng",
				UserGrade: 10,
				Age:       18,
			},
			wantErr: nil,
			want: &model{
				tableName: "test_table",
				fields: map[string]*field{
					"ID": {
						colName: "i_d",
					},
					"Username": {
						colName: "username",
					},
					"UserGrade": {
						colName: "user_grade",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseModel(tc.entity)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Equal(t, tc.want.tableName, got.tableName)
			assert.Equal(t, len(tc.want.fields), len(got.fields))
			for key, f := range tc.want.fields {
				assert.Equal(t, f, got.fields[key])
			}
		})
	}
}
