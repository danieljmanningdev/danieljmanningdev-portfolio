package contracts

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

const contractsBasePath = "/dashboard/contracts"

type ContractsHandler struct {
	contractRepository *repository.ContractRepository
	clientRepository   *repository.ClientRepository
	projectRepository  *repository.ProjectRepository

	contractsTemplates    *template.Template
	newContractTemplates  *template.Template
	contractTemplates     *template.Template
	editContractTemplates *template.Template
}

func NewContractsHandler(
	db *sql.DB,
	templateDir string,
) (*ContractsHandler, error) {
	contractRepo := repository.NewContractRepository(db)
	clientRepo := repository.NewClientRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	contractsTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/contracts/index.html",
	)
	if err != nil {
		return nil, err
	}

	newContractTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/contracts/new.html",
	)
	if err != nil {
		return nil, err
	}

	contractTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/contracts/show.html",
	)
	if err != nil {
		return nil, err
	}

	editContractTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/contracts/edit.html",
	)
	if err != nil {
		return nil, err
	}

	return &ContractsHandler{
		contractRepository:    contractRepo,
		clientRepository:      clientRepo,
		projectRepository:     projectRepo,
		contractsTemplates:    contractsTemplates,
		newContractTemplates:  newContractTemplates,
		contractTemplates:     contractTemplates,
		editContractTemplates: editContractTemplates,
	}, nil
}

func (h *ContractsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.handleContractGET(w, r)

	case http.MethodPost:
		h.handleContractPOST(w, r)

	default:
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ContractsHandler) handleContractGET(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		contractsBasePath,
	)

	if path == "" || path == "/" {
		h.listContracts(w, r)
		return
	}

	if path == "/new" {
		h.renderNewContract(
			w,
			r,
			contractForm{
				Status: contractStatusDraft,
			},
			contractFormErrors{},
			http.StatusOK,
		)
		return
	}

	if strings.HasSuffix(path, "/edit") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/edit",
		)

		id, err := strconv.ParseInt(
			idString,
			10,
			64,
		)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		h.editContract(w, r, id)
		return
	}

	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.showContract(w, r, id)
}

func (h *ContractsHandler) handleContractPOST(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		contractsBasePath,
	)

	if path == "/new" {
		h.createContract(w, r)
		return
	}

	if strings.HasSuffix(path, "/cancel") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/cancel",
		)

		id, err := strconv.ParseInt(
			idString,
			10,
			64,
		)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		h.cancelContract(w, r, id)
		return
	}

	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.updateContract(w, r, id)
}
