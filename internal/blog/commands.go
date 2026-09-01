// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package blog

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func parseForm(r *http.Request) (Form, error) {
	if err := r.ParseForm(); err != nil {
		return Form{}, err
	}

	return Form{
		Title:   r.FormValue("title"),
		Slug:    r.FormValue("slug"),
		Excerpt: r.FormValue("excerpt"),
		Content: r.FormValue("content"),
		Status:  r.FormValue("status"),
	}, nil
}

func (h *AdminHandler) createPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	form, err := parseForm(r)
	if err != nil {
		http.Error(
			w,
			"invalid form submission",
			http.StatusBadRequest,
		)
		return
	}

	formErrors := ValidateForm(form)

	if formErrors.Any() {
		h.renderNewPost(
			w,
			r,
			form,
			formErrors,
			http.StatusUnprocessableEntity,
		)
		return
	}

	post := form.ToPost()

	if post.Status == "published" {
		now := time.Now().UTC()
		post.PublishedAt = &now
	}

	id, err := h.repository.Create(
		r.Context(),
		post,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create blog post",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		blogAdminBasePath+
			"/"+
			strconv.FormatInt(id, 10)+
			"/edit",
		http.StatusSeeOther,
	)
}

func (h *AdminHandler) updatePost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	form, err := parseForm(r)
	if err != nil {
		http.Error(
			w,
			"invalid form submission",
			http.StatusBadRequest,
		)
		return
	}

	formErrors := ValidateForm(form)

	if formErrors.Any() {
		h.renderEditPost(
			w,
			r,
			id,
			form,
			formErrors,
			http.StatusUnprocessableEntity,
		)
		return
	}

	existingPost, err := h.repository.GetByID(
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

	post := form.ToPostWithID(id)

	switch post.Status {
	case "published":
		if existingPost.PublishedAt != nil {
			post.PublishedAt = existingPost.PublishedAt
		} else {
			now := time.Now().UTC()
			post.PublishedAt = &now
		}

	case "draft":
		post.PublishedAt = nil
	}

	err = h.repository.Update(
		r.Context(),
		post,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to update blog post",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		blogAdminBasePath+
			"/"+
			strconv.FormatInt(id, 10)+
			"/edit",
		http.StatusSeeOther,
	)
}

func (h *AdminHandler) deletePost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	err := h.repository.Delete(
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
			"failed to delete blog post",
			http.StatusInternalServerError,
		)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set(
			"HX-Redirect",
			blogAdminBasePath,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(
		w,
		r,
		blogAdminBasePath,
		http.StatusSeeOther,
	)
}

func (h *AdminHandler) publishPost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	err := h.repository.Publish(
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
			"failed to publish blog post",
			http.StatusInternalServerError,
		)
		return
	}

	redirectToEdit(
		w,
		r,
		id,
	)
}

func (h *AdminHandler) unpublishPost(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	err := h.repository.Unpublish(
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
			"failed to unpublish blog post",
			http.StatusInternalServerError,
		)
		return
	}

	redirectToEdit(
		w,
		r,
		id,
	)
}

func redirectToEdit(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	location := blogAdminBasePath +
		"/" +
		strconv.FormatInt(id, 10) +
		"/edit"

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set(
			"HX-Redirect",
			location,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(
		w,
		r,
		location,
		http.StatusSeeOther,
	)
}
