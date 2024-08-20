package handlers

import (
	"encoding/json"
	"net/http"

	"tier3-app/models"
	"tier3-app/services"
)

type QueueHandler struct {
	Service *services.QueueService
}

func NewQueueHandler(service *services.QueueService) *QueueHandler {
	return &QueueHandler{Service: service}
}

func (h *QueueHandler) GetQueue(w http.ResponseWriter, r *http.Request) {
	queue, err := h.Service.GetQueue()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if queue == nil {
		queue = []models.QueueItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queue)
}

func (h *QueueHandler) AddToQueue(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.Service.AddToQueue(requestData.Name, requestData.Email, requestData.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
