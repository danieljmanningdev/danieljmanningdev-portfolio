package blog

import (
	"strings"
	"testing"
)

func TestJournalFieldLimits(t *testing.T) {
	valid := Form{Title: "Title", Slug: "a-post", Excerpt: "An excerpt", Content: "Some Markdown", Status: "draft"}
	if errors := ValidateForm(valid); errors.Any() {
		t.Fatalf("valid form: %+v", errors)
	}
	for _, mutate := range []func(*Form){
		func(f *Form) { f.Title = strings.Repeat("a", 201) },
		func(f *Form) { f.Slug = strings.Repeat("a", 201) },
		func(f *Form) { f.Excerpt = strings.Repeat("a", 501) },
		func(f *Form) { f.Content = strings.Repeat("a", (1<<20)+1) },
		func(f *Form) { f.Content = string([]byte{0xff}) },
	} {
		form := valid
		mutate(&form)
		if !ValidateForm(form).Any() {
			t.Fatal("invalid content accepted")
		}
	}
}

func TestJournalMarkdownDoesNotEnableRawHTMLOrScriptURLs(t *testing.T) {
	html, err := RenderMarkdown("<script>alert(1)</script>\n\n[link](javascript:alert(1))\n\n<img src=x onerror=alert(1)>")
	if err != nil {
		t.Fatal(err)
	}
	for _, unsafe := range []string{"<script", "onerror=", "href=\"javascript:"} {
		if strings.Contains(string(html), unsafe) {
			t.Fatalf("unsafe Markdown output: %s", html)
		}
	}
}
