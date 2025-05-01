package handlers

import (
	"LoadBalancer/internal/models"
	"LoadBalancer/internal/repository"
	"encoding/json"
	"net/http"
	"strings"
)

type ClientHandler struct {
	Repo repository.ClientRepository
}

// POST /clients
func (h *ClientHandler) CreateOrUpdateClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if err := h.Repo.UpsertClient(r.Context(), &client); err != nil {
		http.Error(w, `{"error": "DB error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Client saved"})
}

// GET /clients
func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	clients, err := h.Repo.ListClients(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Failed to list clients"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// GET /clients/{id}
func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractID(r.URL.Path, "/clients/")
	client, err := h.Repo.GetClient(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "Client not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// DELETE /clients/{id}
func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractID(r.URL.Path, "/clients/")
	if err := h.Repo.DeleteClient(r.Context(), id); err != nil {
		http.Error(w, `{"error": "Failed to delete"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// 🔧 Вспомогательная функция для извлечения ID из пути
func extractID(path string, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}
