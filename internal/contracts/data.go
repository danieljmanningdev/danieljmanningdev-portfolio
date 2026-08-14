package contracts

import (
	"fmt"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type contractListItem struct {
	models.Contract
	DisplayValue string
}

type contractsPageData struct {
	Title     string
	Contracts []contractListItem
}

type contractPageData struct {
	Title        string
	Contract     models.Contract
	DisplayValue string
}

type contractFormPageData struct {
	Title    string
	Form     contractForm
	Errors   contractFormErrors
	Clients  []models.Client
	Projects []models.Project
}

func formatContractValue(
	valueCents *int64,
) string {
	if valueCents == nil {
		return ""
	}

	return fmt.Sprintf(
		"£%.2f",
		float64(*valueCents)/100,
	)
}
