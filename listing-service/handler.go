package main

import (
	"encoding/json"
	"github.com/adewaleafolabi/listing/db"
	"github.com/adewaleafolabi/listing/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/xid"
	"net/http"
	"strconv"
	"database/sql"
)

func PropertyRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{id}", GetProperty)
	router.Post("/", CreateProperty)
	router.Get("/", ListProperties)
	return router
}

func GetProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	property, err := db.FindProperty(ctx, id)

	if err != nil {
		response := make(map[string]string)
		code:= http.StatusInternalServerError
		response["message"] = err.Error()

		switch err {
		case sql.ErrNoRows:
			response["message"] = "Not found"
			code = http.StatusNotFound
		default:
			response["message"] = err.Error()
		}

		render.Status(r, code)
		render.JSON(w, r, response)

	} else {
		render.JSON(w, r, property)
	}

}

func ListProperties(w http.ResponseWriter, r *http.Request) {
	lastId := chi.URLParam(r, "lastId")
	limit, _ := strconv.ParseInt(chi.URLParam(r, "limit"), 0, 64)

	if limit <= 0 {
		limit = 100
	}
	props, err := db.ListProperties(r.Context(), limit, lastId)

	if err != nil {
		response := make(map[string]string)
		response["message"] = err.Error()
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
	} else {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, props)
	}
}

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	data := model.Property{}
	err := json.NewDecoder(r.Body).Decode(&data)
	data.ID = xid.New().String()
	response := make(map[string]string)

	if err != nil {

		response["message"] = err.Error()
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response)
		return
	}

	if err = db.SaveProperty(r.Context(), data); err != nil {
		response["message"] = err.Error()
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, data)

}
