package main

import (
	"log"
	"net/http"

	"github.com/segmentio/ksuid"
	"github.com/tinrab/meower/db"
	"github.com/tinrab/meower/event"
	"github.com/tinrab/meower/schema"
	"github.com/tinrab/meower/util"
)

type CreateMeowResponse struct {
	ID string `json:"id"`
}

func createMeowHandler(w http.ResponseWriter, r *http.Request) {
	// Read parameters
	body := r.FormValue("body")
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	// Create meow
	meow := schema.Meow{
		ID:   ksuid.New().String(),
		Body: body,
	}
	if err := db.InsertMeow(meow); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create meow")
		return
	}

	// Publish event
	if err := event.PublishMeowCreated(meow); err != nil {
		log.Println(err)
	}

	// Return new meow
	util.ResponseOk(w, CreateMeowResponse{ID: meow.ID})
}
