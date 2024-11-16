package crud

import (
	"net/http"
	"strconv"

	"transactions/resources"
	"transactions/stores"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type TransactionCrud struct {
	transactionStore stores.TransactionStore
}

func NewTransactionCrud(transactionStore stores.TransactionStore) *TransactionCrud {
	return &TransactionCrud{
		transactionStore: transactionStore,
	}
}

func (tc *TransactionCrud) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/", tc.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", tc.Get)
	})

	return r
}

func (tc *TransactionCrud) Get(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := chi.URLParam(r, "id")
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		render.Render(w, r, resources.ErrInvalidRequest(err))
		return
	}
	transaction, err := tc.transactionStore.Retrieve(transactionID)
	if err != nil {
		render.Render(w, r, resources.ErrNotFound(err))
		return
	}

	if err := render.Render(w, r, resources.NewTransactionResponse(transaction)); err != nil {
		render.Render(w, r, resources.ErrRender(err))
		return
	}
}

func (c *TransactionCrud) Create(w http.ResponseWriter, r *http.Request) {
	data := &resources.NewTransactionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, resources.ErrInvalidRequest(err))
		return
	}
	transaction := data.Transaction
	result, err := c.transactionStore.Create(transaction)
	if err != nil {
		render.Render(w, r, resources.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, resources.NewTransactionResponse(result))
}
