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
		"MSecTime": func(t time.Time) string {
			timestamp := t.Format("2006/1/2 15:04:05.000")
			unixTimestamp := t.Unix()
			return fmt.Sprintf("%s (%d)", timestamp, unixTimestamp)
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
