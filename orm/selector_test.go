package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no from",
			builder: &Selector[User]{},
			wantErr: nil,
			wantQuery: &Query{
				SQL:  "SELECT * FROM `user`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: (&Selector[User]{}).From("user"),
			wantErr: nil,
			wantQuery: &Query{
				SQL:  "SELECT * FROM `user`;",
				Args: nil,
			},
		},
		{
			name:    "from db.table",
			builder: (&Selector[User]{}).From("user.user"),
			wantErr: nil,
			wantQuery: &Query{
				SQL:  "SELECT * FROM `user`.`user`;",
				Args: nil,
			},
		},
		{
			name: "where id = 1",
			builder: (&Selector[User]{}).From("user").
				Where(C("id").EQ(1)),
			wantErr: nil,
			wantQuery: &Query{
				SQL: "SELECT * FROM `user` WHERE (`id` = ?);",
				Args: []any{
					1,
				},
			},
		},
		{
			name: "where id = 1 or id =2",
			builder: (&Selector[User]{}).From("user").
				Where(C("id").EQ(1).OR(C("id").EQ(2))),
			wantErr: nil,
			wantQuery: &Query{
				SQL: "SELECT * FROM `user` WHERE ((`id` = ?) OR (`id` = ?));",
				Args: []any{
					1, 2,
				},
			},
		},

		{
			name: "not username=xxcheng",
			builder: (&Selector[User]{}).From("user").
				Where(NOT(C("username").EQ("xxcheng"))),
			wantErr: nil,
			wantQuery: &Query{
				SQL: "SELECT * FROM `user` WHERE (NOT (`username` = ?));",
				Args: []any{
					"xxcheng",
				},
			},
		},
		{
			name: "where id = 1 or id =2 and not username=xxcheng",
			builder: (&Selector[User]{}).From("user").
				Where(C("id").EQ(1).OR(C("id").EQ(2)).AND(NOT(C("username").EQ("xxcheng")))),
			wantErr: nil,
			wantQuery: &Query{
				SQL: "SELECT * FROM `user` WHERE (((`id` = ?) OR (`id` = ?)) AND (NOT (`username` = ?)));",
				Args: []any{
					1, 2, "xxcheng",
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			q, err := testCase.builder.Build()
			assert.Equal(t, nil, err)
			assert.Equal(t, testCase.wantQuery, q)
		})
	}
}

type User struct {
	ID       int
	Username string
	Grade    uint8
}
