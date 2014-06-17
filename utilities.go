package main

import "path/filepath"

import "github.com/flosch/pongo"

func Render(name string, context *pongo.Context) string {
    template := pongo.Must(pongo.FromFile(filepath.Join(TEMPLATE_PATH, name), nil))
    template.SetDebug(true)
    out, err := template.Execute(context)
    if err != nil { panic(err) }
    return *out
}
