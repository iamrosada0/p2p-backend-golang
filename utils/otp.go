package utils

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

// GenerateOTP gera um código OTP de 6 dígitos
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(999999))
}

// SendOTPByEmail envia um código OTP por email
func SendOTPByEmail(email, otp string) error {
	// Configuração do cliente de e-mail (utilizando go-gomail)
	d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-email-password")

	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "OTP Verification")
	m.SetBody("text/plain", fmt.Sprintf("Seu código OTP é: %s", otp))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
