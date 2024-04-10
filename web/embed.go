// Copyright 2022 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/log"
)

//go:embed **/dist/*
var files embed.FS

func walkNames(fsys fs.FS) []string {
	var names []string
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			names = append(names, path)
		}
		return nil
	}
	if err := fs.WalkDir(fsys, ".", walkDirFunc); err != nil {
		panic("assetNames failure: " + err.Error())
	}
	return names
}

func Init(engine *gin.Engine) {
	for _, name := range []string{"formatify/dist"} {
		sub, _ := fs.Sub(files, name)
		log.Zap.Debugln("web/"+name+":", walkNames(sub))
		engine.StaticFS("/"+name, http.FS(sub))
	}
}
