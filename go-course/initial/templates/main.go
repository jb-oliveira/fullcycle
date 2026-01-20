package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func MyUpper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	// curso := Curso{"GO", 40}
	// tmp := template.New("CursoTemplate")
	// tmp, err := tmp.Parse("Curso: {{.Nome}} - Carga Horaria {{.CargaHoraria}}")
	// if err != nil {
	// 	panic(err)
	// }
	// err = tmp.Execute(os.Stdout, curso)
	// if err != nil {
	// 	panic(err)
	// }

	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	t := template.Must(template.New("content.html").ParseFiles(templates...))
	// 	err := t.Execute(w, Cursos{
	// 		{"GO", 40},
	// 		{"Java", 40},
	// 		{"Python", 50},
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// })
	// http.ListenAndServe(":8080", nil)

	t := template.New("content.html")
	t.Funcs(template.FuncMap{"MyUpper": MyUpper})
	t = template.Must(t.ParseFiles(templates...))
	err := t.Execute(os.Stdout, Cursos{
		{"GO", 40},
		{"Java", 40},
		{"Python", 50},
	})
	if err != nil {
		log.Fatal(err)
	}
}
