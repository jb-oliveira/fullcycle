package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

type Course struct {
	Nome         string
	CargaHoraria int
}

type Courses []Course

func MyUpper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	// course := Course{"GO", 40}
	// tmp := template.New("CourseTemplate")
	// tmp, err := tmp.Parse("Course: {{.Nome}} - Carga Horaria {{.CargaHoraria}}")
	// if err != nil {
	// 	panic(err)
	// }
	// err = tmp.Execute(os.Stdout, course)
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
	// 	err := t.Execute(w, Courses{
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
	err := t.Execute(os.Stdout, Courses{
		{"GO", 40},
		{"Java", 40},
		{"Python", 50},
	})
	if err != nil {
		log.Fatal(err)
	}
}
