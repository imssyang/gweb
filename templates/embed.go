// Copyright 2020 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/log"
)

var (
	//go:embed *.tmpl **/*
	files embed.FS

	funcMap = template.FuncMap{
		"Year": func() int {
			return time.Now().Year()
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"FormatAsDate": func(t time.Time) string {
			year, month, day := t.Date()
			return fmt.Sprintf("%d/%02d/%02d", year, month, day)
		},
		"NoEscape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
)

func Init(engine *gin.Engine) {
	templ := template.Must(
		template.New("").Delims("{{", "}}").Funcs(funcMap).ParseFS(
			files,
			"*.tmpl", "**/*.tmpl",
		),
	)

	for _, templ := range templ.Templates() {
		log.Zap.Debugln("template: ", templ.Name())
	}

	engine.SetFuncMap(funcMap)
	engine.SetHTMLTemplate(templ)
}
