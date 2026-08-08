package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type Post struct {
	Title       string
	Slug        string
	Date        string
	Description string
	Summary     string
	Content     template.HTML
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")

	if slug != "" {
		filePath := filepath.Join("content", slug+".md")
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		markdownContent := contentBytes
		if bytes.HasPrefix(contentBytes, []byte("---")) {
			parts := bytes.SplitN(contentBytes, []byte("---"), 3)
			if len(parts) >= 3 {
				markdownContent = parts[2]
			}
		}

		var buf strings.Builder
		if err := goldmark.Convert(markdownContent, &buf); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Format title nicely from slug as fallback
		formattedTitle := strings.ReplaceAll(slug, "-", " ")
		words := strings.Fields(formattedTitle)
		for i, w := range words {
			if len(w) > 0 {
				words[i] = strings.ToUpper(w[:1]) + w[1:]
			}
		}
		titleStr := strings.Join(words, " ")

		data := struct {
			ActivePost *Post
			Posts      []Post
		}{
			ActivePost: &Post{
				Title:   titleStr,
				Content: template.HTML(buf.String()),
			},
		}

		tmpl, err := template.ParseFiles("backend/templates/blog.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	// List all posts and extract frontmatter meta
	files, err := os.ReadDir("content")
	var posts []Post
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
				slugName := strings.TrimSuffix(file.Name(), ".md")

				// Read file to pull frontmatter metadata
				fileBytes, err := os.ReadFile(filepath.Join("content", file.Name()))
				title := strings.ReplaceAll(slugName, "-", " ")
				date := ""
				description := ""

				if err == nil && bytes.HasPrefix(fileBytes, []byte("---")) {
					parts := bytes.SplitN(fileBytes, []byte("---"), 3)
					if len(parts) >= 3 {
						frontmatter := string(parts[1])
						lines := strings.Split(frontmatter, "\n")
						for _, line := range lines {
							if strings.HasPrefix(line, "title:") {
								title = strings.Trim(strings.TrimPrefix(line, "title:"), " \"'")
							} else if strings.HasPrefix(line, "date:") {
								date = strings.Trim(strings.TrimPrefix(line, "date:"), " \"'")
							} else if strings.HasPrefix(line, "description:") {
								description = strings.Trim(strings.TrimPrefix(line, "description:"), " \"'")
							}
						}
					}
				}

				posts = append(posts, Post{
					Title:       title,
					Slug:        slugName,
					Date:        date,
					Description: description,
				})
			}
		}
	}

	data := struct {
		ActivePost *Post
		Posts      []Post
	}{
		ActivePost: nil,
		Posts:      posts,
	}

	tmpl, err := template.ParseFiles("backend/templates/blog.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
