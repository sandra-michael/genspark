package mail

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(orderId string) {

	// Replace with your Mailtrap credentials
	mailHost := "sandbox.smtp.mailtrap.io" // Mailtrap SMTP server
	mailPort := 587                        // Port (587 is typical for Mailtrap)
	mailUser := "73a60a43601fa6"           // Mailtrap username
	mailPass := "8aebece8aef83c"           // Mailtrap password

	// Email details
	from := "sender@example.com"  // Sender email
	to := "recipient@example.com" // Recipient email
	subject := "Test Order Confirmation"
	body := fmt.Sprint("Hey, \n this is an order confirmation email for order id : ", orderId)

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Configure the SMTP dialer
	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Could not send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
