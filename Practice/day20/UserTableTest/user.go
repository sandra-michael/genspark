package UserTableTest

// Create a struct user,
//     Fields: Name, Email, Pay, Bonus

type User struct {
	Name  string
	Email string
	Pay   float64
	Bonus float64
}

//Create a method that calculates total salary for a month for a user

func (u User) CalculateTotalSalary(month string) float64 {

	switch month {
	case "December":
		return u.Pay + u.Bonus
	default:
		return u.Pay
	}

}
