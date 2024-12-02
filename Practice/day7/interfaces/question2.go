package main

import "fmt"

// q2. Define an interface named Notifier with two methods:
//     SendEmail(): A method that takes a string argument representing the email content and returns a string indicating the email has been sent.
//     SendSMS(): A method that takes a string argument representing the SMS content and returns a string indicating the SMS has been sent.
//     Implement this interface in a struct called NotificationService.

type Notifier interface {
	SendEmail(str string) string
	SendSMS(str string) string
}

type NotificationService struct {
	//this message can be path of the method body
	ipMessage string
}

func (ns NotificationService) SendEmail(str string) string {
	fmt.Println("An Email notification is sent for message : ", ns.ipMessage)
	opMessage := "Email sent"
	return opMessage
}

func (ns NotificationService) SendSMS(str string) string {
	fmt.Println("An SMS notification is sent for message : ", ns.ipMessage)
	opMessage := "SMS sent"
	return opMessage
}

func main() {
	//ns := NotificationService{"some content"}
	var ns Notifier = NotificationService{"some content"}
	ns.SendEmail("sdfd ")
	ns.SendSMS(" sdfd  dsf ")

}
