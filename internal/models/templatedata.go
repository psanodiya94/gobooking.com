package models

import "github.com/psanodiya94/gobooking.com/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	FloatMap  map[string]float32
	IntMap    map[string]int
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
