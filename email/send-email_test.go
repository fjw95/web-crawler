package email

import (
	"os"
	"strconv"
	"testing"

	"golang.org/x/net/context"
)

func TestSendEmail(t *testing.T) {
	port, _ := strconv.Atoi(os.Getenv("TEST_SMTP_PORT"))
	c := WithSmtp(context.Background(),
		&SmtpConfig{
			Server:   os.Getenv("TEST_SMTP_ADDRESS"),
			Port:     port,
			User:     os.Getenv("TEST_SMTP_USER"),
			Password: os.Getenv("TEST_SMTP_PASSWORD"),
		},
	)

	message := NewHtmlMessage("helo test", "hello everybody!!")
	message.To = []string{"fail@gmail.com"}

	c = WithMessage(c, &message)

	if err := Send(c); err != nil {
		t.Error(err)
	}
}
