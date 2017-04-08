package widgets

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

// NewForm new form
func NewForm(req *http.Request, lang, action, next, title string, fields ...interface{}) Form {
	if next == "" {
		next = action
	}
	return Form{
		"id":             uuid.New().String(),
		"lang":           lang,
		"method":         http.MethodPost,
		"action":         action,
		"next":           next,
		"title":          title,
		"fields":         fields,
		csrf.TemplateTag: csrf.TemplateField(req),
	}
}

// Form form
type Form map[string]interface{}

// Method set method
func (p Form) Method(m string) {
	p["method"] = m
}

// Field field
type Field map[string]interface{}

// Readonly set readonly
func (p Field) Readonly() {
	p["readonly"] = true
}

// Helper set helper message
func (p Field) Helper(h string) {
	p["helper"] = h
}

// NewTextarea new textarea field
func NewTextarea(id, label, value string, row int) Field {
	return Field{
		"id":    id,
		"type":  "textarea",
		"label": label,
		"value": value,
		"row":   3,
	}
}

// NewTextField new text field
func NewTextField(id, label, value string) Field {
	return Field{
		"id":    id,
		"type":  "text",
		"label": label,
		"value": value,
	}
}

// NewHiddenField new hidden field
func NewHiddenField(id string, value interface{}) Field {
	return Field{
		"id":    id,
		"type":  "hidden",
		"value": value,
	}
}

// NewEmailField new email field
func NewEmailField(id, label, value string) Field {
	return Field{
		"id":    id,
		"type":  "email",
		"label": label,
		"value": value,
	}
}

// NewPasswordField new password field
func NewPasswordField(id, label, helper string) Field {
	return Field{
		"id":     id,
		"type":   "password",
		"label":  label,
		"helper": helper,
	}
}
