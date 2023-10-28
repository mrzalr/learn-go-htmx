package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

func main() {
	http.HandleFunc("/", numberFormatterHandler)
	http.HandleFunc("/num-input", handleNumInput)
	http.HandleFunc("/meta-num-input", handleMetaNumInput)
	http.HandleFunc("/meta-num-output", handleMetaNumOutput)

	http.ListenAndServe(":8080", nil)
}

func execTemplate(w http.ResponseWriter, path string, data any) {
	templ, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func numberFormatterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		execTemplate(w, path.Join("template", "index.gohtml"), map[string]any{"numresult": "-"})
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func handleNumInput(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data := r.FormValue("input-num")
		result := parseNumber(data)

		w.Header().Set("Content-Type", "text/html")
		updatedContent := []byte(result)
		w.Write(updatedContent)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func parseNumber(input string) string {
	numbers := strings.Split(input, "\n")
	res := []string{}
	for _, number := range numbers {
		if len(number) == 0 {
			continue
		}

		_, err := strconv.ParseInt(number, 10, 64)

		if err != nil {
			continue
		}

		res = append(res, fmt.Sprintf("62%s", number))
	}
	return strings.Join(res, "\n")
}

func handleMetaNumInput(w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.FormValue("input-num"), "\n")
	count := len(data)
	invalid := 0

	for _, number := range data {
		if len(number) == 0 {

			invalid++
			continue
		}

		_, err := strconv.ParseInt(number, 10, 64)

		if err != nil {
			invalid++
			continue
		}
	}

	templ := `
	<p class="font-bold text-gray-500">count: %d</p>
	<p class="font-bold text-gray-500">
	invalid: %d
	</p>`

	w.Header().Set("Content-Type", "text/html")
	updatedContent := []byte(fmt.Sprintf(templ, count, invalid))
	w.Write(updatedContent)
}

func handleMetaNumOutput(w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.FormValue("output-num"), "\n")
	count := len(data)

	templ := `<p class="font-bold text-gray-500">count: %d</p>`

	w.Header().Set("Content-Type", "text/html")
	updatedContent := []byte(fmt.Sprintf(templ, count))
	w.Write(updatedContent)
}
