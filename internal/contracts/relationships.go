package contracts

import (
	"context"
	"database/sql"
	"errors"
)

func (h *ContractsHandler) validateContractRelationships(
	ctx context.Context,
	form contractForm,
	formErrors contractFormErrors,
) (contractFormErrors, error) {
	if form.ClientID <= 0 {
		return formErrors, nil
	}

	_, err := h.clientRepository.GetByID(
		ctx,
		form.ClientID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			formErrors.ClientID =
				"Selected client does not exist."

			return formErrors, nil
		}

		return formErrors, err
	}

	// A contract may exist without being attached to a project.
	if form.ProjectID <= 0 {
		return formErrors, nil
	}

	project, err := h.projectRepository.GetByID(
		ctx,
		form.ProjectID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			formErrors.ProjectID =
				"Selected project does not exist."

			return formErrors, nil
		}

		return formErrors, err
	}

	if project.ClientID != form.ClientID {
		formErrors.ProjectID =
			"Selected project does not belong to the selected client."
	}

	return formErrors, nil
}
