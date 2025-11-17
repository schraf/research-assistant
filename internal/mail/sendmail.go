package mail

import (
	"log/slog"
	"net/smtp"
	"os"
)

func SendEmail(subject string, body string) error {
	smtpServer := os.Getenv("MAIL_SMTP_SERVER")
	smtpPort := os.Getenv("MAIL_SMTP_PORT")
	senderEmail := os.Getenv("MAIL_SENDER_EMAIL")
	password := os.Getenv("MAIL_SENDER_PASSWORD")
	recipientEmail := os.Getenv("MAIL_RECIPIENT_EMAIL")

	slog.Info("sending_email",
		slog.String("server", smtpServer),
		slog.String("port", smtpPort),
		slog.String("recipient", recipientEmail),
		slog.String("subject", subject),
	)

	// --- Email Content ---
	// The message body must be a byte array.
	// Note the headers (To, Subject, MIME-Version, Content-Type) are crucial for
	// the email client to render the message correctly.
	message := []byte(
		"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" + // This blank line separates headers from the body
			body + "\n",
	)

	auth := smtp.PlainAuth("", senderEmail, password, smtpServer)

	err := smtp.SendMail(
		smtpServer+":"+smtpPort,
		auth,
		senderEmail,
		[]string{recipientEmail},
		message,
	)

	if err != nil {
		return err
	}

	slog.Info("email_sent")

	return nil
}
