package config

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"portfolio-website/internals/model"
	"portfolio-website/internals/templates"

	"github.com/resend/resend-go/v3"
)

type Client struct {
	client *resend.Client
	From   string
}

func NewClient() *Client {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("FROM")
	if from == "" {
		from = "Your Blog <newsletter@yourdomain.com>"
	}
	return &Client{
		client: resend.NewClient(apiKey),
		From:   from,
	}
}

func (c *Client) SendBlogNotification(ctx context.Context, to string, data model.Blog) error {
	siteURL := os.Getenv("BASE_URL")
	if siteURL == "" {
		siteURL = "https://yourdomain.com"
	}

	t := template.Must(template.New("newsletter").Parse(templates.EmailNewsletter))
	var buf bytes.Buffer

	templateData := map[string]interface{}{
		"Blog":    data,
		"SiteURL": siteURL,
	}

	if err := t.Execute(&buf, templateData); err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    c.From,
		To:      []string{to},
		Subject: data.Blog_Title,
		Html:    buf.String(),
	}

	_, err := c.client.Emails.Send(params)
	return err
}

func (c *Client) SendProjectNotification(ctx context.Context, to string, data model.Project) error {
	siteURL := os.Getenv("BASE_URL")
	if siteURL == "" {
		siteURL = "https://yourdomain.com"
	}

	t := template.Must(template.New("project_newsletter").Parse(templates.EmailProject))
	var buf bytes.Buffer

	templateData := map[string]interface{}{
		"Project": data,
		"SiteURL": siteURL,
	}

	if err := t.Execute(&buf, templateData); err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    c.From,
		To:      []string{to},
		Subject: "New Project: " + data.Project_Name,
		Html:    buf.String(),
	}

	_, err := c.client.Emails.Send(params)
	return err
}

func (c *Client) SendCaseStudyNotification(ctx context.Context, to string, data model.Case_Study) error {
	siteURL := os.Getenv("BASE_URL")
	if siteURL == "" {
		siteURL = "https://yourdomain.com"
	}

	t := template.Must(template.New("case_study_newsletter").Parse(templates.EmailCaseStudy))
	var buf bytes.Buffer

	templateData := map[string]interface{}{
		"CaseStudy": data,
		"SiteURL":   siteURL,
	}

	if err := t.Execute(&buf, templateData); err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    c.From,
		To:      []string{to},
		Subject: "New Case Study: " + data.Case_Study_Title,
		Html:    buf.String(),
	}

	_, err := c.client.Emails.Send(params)
	return err
}

