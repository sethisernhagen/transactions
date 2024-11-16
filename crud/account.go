package crud

import (
	"net/http"
	"strconv"
	"transactions/resources"
	"transactions/stores"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type AccountCrud struct {
	accountStore stores.AccountStore
}

func NewAccountCrud(accountStore stores.AccountStore) *AccountCrud {
	return &AccountCrud{
		accountStore: accountStore,
	}
}

func (ac AccountCrud) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/", ac.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ac.Get)
	})

	return r
}

func (ac AccountCrud) Get(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		render.Render(w, r, resources.ErrInvalidRequest(err))
		return
	}
	account, err := ac.accountStore.Retrieve(accountID)
	if err != nil {
		render.Render(w, r, resources.ErrNotFound(err))
		return
	}

	if err := render.Render(w, r, resources.NewAccountResponse(account)); err != nil {
		render.Render(w, r, resources.ErrRender(err))
		return
	}
}

func (ac AccountCrud) Create(w http.ResponseWriter, r *http.Request) {
	data := &resources.NewAccountRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, resources.ErrInvalidRequest(err))
		return
	}
	account := data.Account
	result, err := ac.accountStore.Create(account)
	if err != nil {
		render.Render(w, r, resources.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, resources.NewAccountResponse(result))
}
