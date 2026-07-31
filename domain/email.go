package domain

import "context"

// Kontrak untuk layanan pengirim email
type EmailService interface {
	SendOTP(ctx context.Context, toEmail string, otpCode string) error
}