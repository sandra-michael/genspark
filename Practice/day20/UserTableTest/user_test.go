package UserTableTest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateTotalSalary(t *testing.T) {
	inputUser := User{Name: "testUser", Email: "test@email", Pay: 20000.0, Bonus: 300.0}
	tt := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name: "test normal month",
			//inputUser:  User{Name: "testUser", Email: "test@email", Pay: 20000.0, Bonus: 300.0},
			//inputMonth: "january",
			input: "jan",
			want:  20000.0,
		},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			got := inputUser.CalculateTotalSalary(tr.input)

			// require would fail the current test
			require.Equal(t, tr.want, got)
		})
	}

}

func TestCalculateTotalSalaryV2(t *testing.T) {
	//inputUser := User{Name: "testUser", Email: "test@email", Pay: 20000.0, Bonus: 300.0}
	//can add both the inputs in an args struct
	type args struct {
		User  User
		Month string
	}
	tt := []struct {
		name  string
		input args
		want  float64
	}{
		{
			name:  "test normal month",
			input: args{User: User{Name: "testUser", Email: "test@email", Pay: 20000.0, Bonus: 300.0}, Month: "january"},
			// inputUser:  User{Name: "testUser", Email: "test@email", Pay: 20000.0, Bonus: 300.0},
			// inputMonth: "january",
			want: 20000.0,
		},
		{
			name:  "test bonus month december ",
			input: args{User: User{Name: "testUser", Email: "test@email", Pay: 20000.0, Bonus: 3000.0}, Month: "December"},

			// inputUser:  User{Name: "testUser2", Email: "test2@email", Pay: 20000.0, Bonus: 3000.0},
			// inputMonth: "December",
			want: 23000.0,
		},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			got := tr.input.User.CalculateTotalSalary(tr.input.Month)

			// require would fail the current test
			require.Equal(t, tr.want, got)
		})
	}

}
