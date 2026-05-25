package handler

import (
	"encoding/json"
	"insider-league/domain"
	"net/http"
	"strconv"
)

type ApiHandler struct {
	sim domain.Simulation
}

func NewApiHandler(sim domain.Simulation) *ApiHandler {
	return &ApiHandler{sim: sim}
}

func (h *ApiHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	standings, _ := h.sim.GetStandings()
	json.NewEncoder(w).Encode(standings)
}

func (h *ApiHandler) SimulateNextWeek(w http.ResponseWriter, r *http.Request) {
	h.sim.PlayNextWeek()
	h.GetStandings(w, r)
}

func (h *ApiHandler) PlayAll(w http.ResponseWriter, r *http.Request) {
	h.sim.PlayAll()
	h.GetStandings(w, r)
}

func (h *ApiHandler) EditMatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MatchID   int `json:"match_id"`
		HomeScore int `json:"home_score"`
		AwayScore int `json:"away_score"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	h.sim.EditMatch(req.MatchID, req.HomeScore, req.AwayScore)
	h.GetStandings(w, r)
}

func (h *ApiHandler) GetMatchesByWeek(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	weekStr := r.URL.Query().Get("week")
	if weekStr == "" {
		http.Error(w, "week parametresi zorunludur", http.StatusBadRequest)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 6 {
		http.Error(w, "Geçersiz hafta numarası. 1 ile 6 arasında olmalıdır.", http.StatusBadRequest)
		return
	}

	matches, err := h.sim.GetMatchesByWeek(week)
	if err != nil {
		http.Error(w, "Maçlar getirilirken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if matches == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bu haftaya ait oynanmış bir maç bulunmamaktadır. Lütfen simülasyonu ilerletin.",
		})
		return
	}

	json.NewEncoder(w).Encode(matches)
}
