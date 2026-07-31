package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"pamagi/bootstrap"
	"pamagi/domain"
)

type gmailService struct {
	env *bootstrap.Env
}

func NewGmailService(env *bootstrap.Env) domain.EmailService {
	return &gmailService{env: env}
}

func (g *gmailService) SendOTP(ctx context.Context, toEmail string, otpCode string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := g.env.GmailSenderEmail
	appPassword := g.env.GmailAppPassword

	auth := smtp.PlainAuth("", senderEmail, appPassword, smtpHost)

	// Format Header agar nama pengirim berubah
	header := fmt.Sprintf("From: \"LinguaFlip Support\" <%s>\r\n", senderEmail)
	header += fmt.Sprintf("To: %s\r\n", toEmail)
	header += "Subject: [LinguaFlip] Kode OTP Pemulihan Password\r\n"
	header += "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	// Template HTML yang lebih bagus + Logo LinguaFlip
	// Catatan: Ganti 'https://mydm.cloud' dengan domain aslimu nanti saat production.
	// Saat testing di localhost, Gmail mungkin memblokir gambar lokal.
	logoURL := "https://mydm.cloud/assets/logo.png" 
	
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
				<div style="text-align: center; margin-bottom: 20px;">
					<img src="%s" alt="LinguaFlip Logo" style="max-width: 120px;">
				</div>
				<h2 style="color: #333333; text-align: center;">Pemulihan Password</h2>
				<p style="color: #555555; font-size: 16px;">Halo,</p>
				<p style="color: #555555; font-size: 16px;">Kami menerima permintaan untuk mereset password akun LinguaFlip kamu. Berikut adalah kode OTP kamu:</p>
				<div style="text-align: center; margin: 30px 0;">
					<span style="font-size: 32px; font-weight: bold; color: #1abc9c; letter-spacing: 5px; background-color: #f9f9f9; padding: 10px 20px; border-radius: 4px; border: 1px dashed #1abc9c;">%s</span>
				</div>
				<p style="color: #555555; font-size: 14px;">Kode ini hanya berlaku selama <strong>5 menit</strong>. Jangan berikan kode ini kepada siapa pun.</p>
				<hr style="border: none; border-top: 1px solid #eeeeee; margin: 30px 0;">
				<p style="color: #999999; font-size: 12px; text-align: center;">Ini adalah pesan otomatis dari tim developer LinguaFlip. Mohon tidak membalas email ini.</p>
			</div>
		</body>
		</html>`, logoURL, otpCode)
	
	message := []byte(header + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("gagal mengirim email via Gmail: %v", err)
	}
	return nil
}