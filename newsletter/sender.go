package newsletter

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"portfolio-website/internals/model"

	"github.com/resend/resend-go/v3"
)

type Client struct {
	client *resend.Client
	From   string
}

func NewClient() *Client {
	apiKey := os.Getenv("RESEND_API_KEY")
	return &Client{
		client: resend.NewClient(apiKey),
		From:   "Your Blog <newsletter@yourdomain.com>",
	}
}

func (c *Client) SendBlogNotification(ctx context.Context, to string, data model.Blog) error {
	tmplt := `<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <title>Welcome</title>
    <style>
        body {
            font-family: sans-serif;
            background-color: #f6f9fc;
            padding: 20px;
        }

        .container {
            background-color: #ffffff;
            padding: 40px;
            border-radius: 8px;
            max-width: 600px;
            margin: 0 auto;
        }

        .button {
            background-color: #000000;
            color: #ffffff;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 4px;
            display: inline-block;
        }
    </style>
</head>
</html>,`

	t := template.Must(template.New("newsletter").Parse(tmplt))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From: c.From,
		To: []string{to},
		Subject: data.Blog_Title,
		Html: buf.String(),
	}

	_, err := c.client.Emails.Send(params)
	return err
}
