package render

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

type Renderer struct {
	tmpls map[string]*template.Template
}

func New(files fs.FS) (*Renderer, error) {
	tmpls := make(map[string]*template.Template)
	pages := []string{"dashboard", "profile", "notifications", "course-catalog", "course-page"}

	for _, page := range pages {
		tmpl, err := template.ParseFS(files, "templates/layouts/base.html", fmt.Sprintf("templates/%s.html", page))
		if err != nil {
			return nil, fmt.Errorf("error parsing %s template: %w", page, err)
		}
		tmpls[page] = tmpl
	}

	return &Renderer{tmpls: tmpls}, nil
}

func (r *Renderer) Render(w io.Writer, name string, data any) error {
	tmpl, ok := r.tmpls[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}
