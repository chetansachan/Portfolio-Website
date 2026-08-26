package templates

import _ "embed"

//go:embed email_newsletter_template.html
var EmailNewsletter string

//go:embed email_project_template.html
var EmailProject string
