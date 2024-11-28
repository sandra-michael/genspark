package main

/* q2. Create a program to store and print a person's and their project's details. Declare and initialize variables for the following details,
   Project name (string)
   Code lines written (uint8)
   Bugs found (int)
   Is the project complete? (bool)
   Average lines of code written per hour (float64)
   Team lead name (string)
   Project deadline in days (int)
   Additionally, demonstrate a uint overflow by initializing the largest possible value for uint and then adding 1 to it */

// q4. Print default values and Type names of variables from question 2 using printf
// Quick Tip, Use %v if not sure about what verb should be used,
// but don't use it in this question :)
// but generally using %v should be fine

import "fmt"

type PersonProjDetail struct {
	Name                      string
	ProjectName               string
	CodeLines                 uint8
	BugsFound                 int
	IsComplete                bool
	AverageLinesOfCodePerHour float64
	TeamLeadName              string
	ProjectDeadlineInDays     int
}

func PrintPerson(p PersonProjDetail) {
	fmt.Printf("Name : %s , Type : %T \n ", p.Name, p.Name)
	fmt.Println(p.ProjectName)
	fmt.Printf("CodeLines : %v , Type : %T", p.CodeLines, p.CodeLines)
	fmt.Println(p.BugsFound)
	fmt.Println(p.IsComplete)
	fmt.Println(p.AverageLinesOfCodePerHour)
	fmt.Println(p.TeamLeadName)
	fmt.Println(p.ProjectDeadlineInDays)
}

func main() {
	personDet := PersonProjDetail{
		Name:                      "Sandra ",
		ProjectName:               "Project1",
		CodeLines:                 100, //.\PersonProject.go:41:30: cannot use 1000 (untyped int constant) as uint8 value in struct literal (overflows)
		BugsFound:                 3,
		IsComplete:                true,
		AverageLinesOfCodePerHour: 3.2,
		TeamLeadName:              "tl",
		ProjectDeadlineInDays:     3,
	}

	PrintPerson(personDet)

}
