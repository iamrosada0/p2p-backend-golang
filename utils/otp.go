package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP gera um código OTP de 6 dígitos
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(999999))
}
