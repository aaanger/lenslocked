package controllers

import (
	"html/template"
	"net/http"
)

type Static struct {
	Template Template
}

func (static Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	static.Template.Execute(w, nil)
}

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{Question: "Is there a free version?",
			Answer: "Yes",
		},
		{Question: "What are your support hours?",
			Answer: "We have support staff answering emails 24/7",
		},
		{Question: "How do I contact support?",
			Answer: `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
