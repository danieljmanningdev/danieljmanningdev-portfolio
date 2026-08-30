// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package clients

import (
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type clientsPageData struct {
	Title   string
	Clients []models.Client
}

type clientPageData struct {
	Title  string
	Client models.Client
}

type clientFormPageData struct {
	Title  string
	Form   clientForm
	Errors clientFormErrors
}
