package blog

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	authservice "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
)

func (h *AdminHandler) listPosts(
	w http.ResponseWriter,
	r *http.Request,
) {
	posts, err := h.repository.ListAll(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to load blog posts",
			http.StatusInternalServerError,
		)
		return
	}

	h.render(
		w,
		h.indexTemplates,
		ListPageData{
			Title: "Blog — Daniel J. Manning",
			Posts: posts,
		},
	)
}

func (h *AdminHandler) editPost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	post, err := h.repository.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to load blog post",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderEditPost(
		w,
		r,
		id,
		FormFromPost(post),
		FormErrors{},
		http.StatusOK,
	)
}

func (h *AdminHandler) renderNewPost(
	w http.ResponseWriter,
	r *http.Request,
	form Form,
	errors FormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = "draft"
	}

	h.renderStatus(
		w,
		h.newTemplates,
		FormPageData{
			Title:     "New Blog Post — Daniel J. Manning",
			Form:      form,
			Errors:    errors,
			CSRFToken: csrfTokenFromRequest(r),
		},
		status,
	)
}

func (h *AdminHandler) renderEditPost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
	form Form,
	errors FormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = "draft"
	}

	h.renderStatus(
		w,
		h.editTemplates,
		EditPageData{
			Title:     "Edit Blog Post — Daniel J. Manning",
			PostID:    id,
			Form:      form,
			Errors:    errors,
			CSRFToken: csrfTokenFromRequest(r),
		},
		status,
	)
}

func (h *AdminHandler) render(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
) {
	h.renderStatus(
		w,
		tmpl,
		data,
		http.StatusOK,
	)
}

func (h *AdminHandler) renderStatus(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
	status int,
) {
	var body bytes.Buffer

	if err := tmpl.ExecuteTemplate(
		&body,
		"base",
		data,
	); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(status)

	_, _ = body.WriteTo(w)
}

func csrfTokenFromRequest(
	r *http.Request,
) string {
	token, _ := authservice.AdminCSRFTokenFromContext(
		r.Context(),
	)

	return token
}
