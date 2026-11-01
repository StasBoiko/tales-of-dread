package handler

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"unicode/utf8"

	"tales-of-dread/internal/domain"
)

// Server coordinates HTTP routing, API requests, and web asset delivery.
type Server struct {
	game   *domain.Game
	logger *slog.Logger
	httpFS http.FileSystem
	rawFS  fs.FS
}

// SelectArtifactRequest represents the payload for choosing a dark relic.
type SelectArtifactRequest struct {
	ArtifactID string `json:"artifactId"`
}

// ErrorResponse represents a structured error returned to the frontend.
type ErrorResponse struct {
	Error string `json:"error"`
}

// StatusResponse represents the liveness probe response.
type StatusResponse struct {
	Status string `json:"status"`
}

// Part2RollRequest specifies which dice are held before re-rolling.
type Part2RollRequest struct {
	Locked []bool `json:"locked"`
}

// Part2CallBluffRequest specifies whether the player calls bluff or accepts showdown.
type Part2CallBluffRequest struct {
	CallBluff bool `json:"call_bluff"`
}

// Part2RelicRequest specifies mutated relic parameters.
type Part2RelicRequest struct {
	Relic      string `json:"relic"`
	DieIndex   int    `json:"die_index"`
	TargetFace int    `json:"target_face"`
}

// Part2EndingRequest specifies the moral climax choice.
type Part2EndingRequest struct {
	Choice string `json:"choice"`
}

// NewServer creates a configured HTTP server handler.
func NewServer(game *domain.Game, logger *slog.Logger, httpFS http.FileSystem, rawFS fs.FS) *Server {
	return &Server{
		game:   game,
		logger: logger,
		httpFS: httpFS,
		rawFS:  rawFS,
	}
}

// StartGameRequest represents optional configuration when launching a run.
type StartGameRequest struct {
	Nightmare bool `json:"nightmare"`
}

// Routes constructs the HTTP multiplexer with all API and static file routes.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", s.HandleHealthz)
	mux.HandleFunc("/api/start", s.HandleStart)
	mux.HandleFunc("/api/nightmare/start", s.HandleNightmareStart)
	mux.HandleFunc("/api/reset", s.HandleReset)
	mux.HandleFunc("/api/hit", s.HandleHit)
	mux.HandleFunc("/api/stand", s.HandleStand)
	mux.HandleFunc("/api/next-round", s.HandleNextRound)
	mux.HandleFunc("/api/campfire-heal", s.HandleCampfireHeal)
	mux.HandleFunc("/api/leave-campfire", s.HandleLeaveCampfire)
	mux.HandleFunc("/api/select-artifact", s.HandleSelectArtifact)
	mux.HandleFunc("/api/state", s.HandleState)

	// Part II: Witch's Eye Bone Ritual routes
	mux.HandleFunc("/api/part2/start", s.HandlePart2Start)
	mux.HandleFunc("/api/part2/roll", s.HandlePart2Roll)
	mux.HandleFunc("/api/part2/call-bluff", s.HandlePart2CallBluff)
	mux.HandleFunc("/api/part2/next-round", s.HandlePart2NextRound)
	mux.HandleFunc("/api/part2/relic", s.HandlePart2Relic)
	mux.HandleFunc("/api/part2/ending", s.HandlePart2Ending)
	mux.HandleFunc("/api/dev/trigger-part2", s.HandleDevTriggerPart2)
	mux.HandleFunc("/api/nightmare/trigger-part2", s.HandleNightmareTriggerPart2)
	mux.HandleFunc("/api/dev/trigger-ending", s.HandleDevTriggerEnding)

	// Fallback handler for embedded SPA frontend and public assets
	mux.HandleFunc("/", s.HandleStaticOrSPA)

	return mux
}

// HandleHealthz responds to service health and liveness probes.
func (s *Server) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.respondJSON(w, http.StatusOK, StatusResponse{Status: "ok"})
}

// HandleStart initializes or resets the duel run.
func (s *Server) HandleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	nightmare := r.URL.Query().Get("nightmare") == "true"
	if r.Body != nil && r.ContentLength > 0 {
		var req StartGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Nightmare {
			nightmare = true
		}
	}

	if err := s.game.StartGameMode(nightmare); err != nil {
		s.logger.Error("failed starting game", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to start duel")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleNightmareStart starts a run directly in Nightmare mode.
func (s *Server) HandleNightmareStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.StartNightmare(); err != nil {
		s.logger.Error("failed starting nightmare run", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to start nightmare run")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleReset resets the game run back to the initial title and prologue phase.
func (s *Server) HandleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	s.game.ResetToTitle()
	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleHit draws a single card for the player.
func (s *Server) HandleHit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.Hit(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed processing hit", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to draw card")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleStand concludes the player's turn and triggers boss AI resolution.
func (s *Server) HandleStand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.Stand(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) || errors.Is(err, domain.ErrInvalidPlayerHand) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed processing stand", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to stand")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleNextRound advances to a new hand against the active boss.
func (s *Server) HandleNextRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.NextRound(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed advancing to next round", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to advance round")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleCampfireHeal restores player health candles at the respite campfire.
func (s *Server) HandleCampfireHeal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.CampfireHeal(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleLeaveCampfire concludes the campfire rest and advances to dark relic selection.
func (s *Server) HandleLeaveCampfire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.LeaveCampfire(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed leaving campfire", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to leave campfire")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleSelectArtifact attaches a dark relic to the player and unlocks the next boss encounter.
func (s *Server) HandleSelectArtifact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req SelectArtifactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "malformed request payload")
		return
	}

	if utf8.RuneCountInString(strings.TrimSpace(req.ArtifactID)) == 0 {
		s.respondError(w, http.StatusBadRequest, "artifactId is required")
		return
	}

	if err := s.game.SelectArtifact(req.ArtifactID); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) || errors.Is(err, domain.ErrArtifactNotFound) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed selecting artifact", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed selecting artifact")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleState returns the current snapshot of duel state.
func (s *Server) HandleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandlePart2Start begins Part II: Witch's Eye Dice Duel against the Morovaya Witch.
func (s *Server) HandlePart2Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.StartPart2(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed starting part2", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to start part 2")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandlePart2Roll re-rolls unlocked dice in the Witch's duel.
func (s *Server) HandlePart2Roll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req Part2RollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Locked = nil
	}

	if err := s.game.Part2Reroll(req.Locked); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandlePart2CallBluff resolves the hand by either challenging the Witch or showdown.
func (s *Server) HandlePart2CallBluff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req Part2CallBluffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "malformed request payload")
		return
	}

	if err := s.game.Part2CallBluff(req.CallBluff); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed processing call bluff", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to resolve hand")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandlePart2NextRound prepares the next round of dice rolling against the Witch.
func (s *Server) HandlePart2NextRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.game.Part2NextRound(); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.Error("failed preparing next dice round", "error", err)
		s.respondError(w, http.StatusInternalServerError, "failed to advance round")
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandlePart2Relic applies a mutated relic power during the bone ritual.
func (s *Server) HandlePart2Relic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req Part2RelicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "malformed request payload")
		return
	}

	if utf8.RuneCountInString(strings.TrimSpace(req.Relic)) == 0 {
		s.respondError(w, http.StatusBadRequest, "relic is required")
		return
	}

	if err := s.game.Part2Relic(req.Relic, req.DieIndex, req.TargetFace); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandlePart2Ending locks in the fateful conclusion choice.
func (s *Server) HandlePart2Ending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req Part2EndingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "malformed request payload")
		return
	}

	if utf8.RuneCountInString(strings.TrimSpace(req.Choice)) == 0 {
		s.respondError(w, http.StatusBadRequest, "choice is required")
		return
	}

	if err := s.game.Part2Ending(req.Choice); err != nil {
		if errors.Is(err, domain.ErrInvalidPhase) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleStaticOrSPA delivers embedded frontend assets with SPA fallback to index.html.
func (s *Server) HandleStaticOrSPA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	cleanedPath := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if cleanedPath == "" || cleanedPath == "." {
		cleanedPath = "index.html"
	}

	// Check if file exists in raw embedded FS
	if s.rawFS != nil {
		f, err := s.rawFS.Open(cleanedPath)
		if err == nil {
			_ = f.Close()
			http.FileServer(s.httpFS).ServeHTTP(w, r)
			return
		}
	}

	// Fallback to index.html for SPA routing
	indexFile, err := s.httpFS.Open("/index.html")
	if err != nil {
		s.respondError(w, http.StatusNotFound, "not found")
		return
	}
	defer func() {
		_ = indexFile.Close()
	}()

	stat, err := indexFile.Stat()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "failed reading asset")
		return
	}

	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile)
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed encoding json response", "error", err)
	}
}

func (s *Server) respondError(w http.ResponseWriter, status int, msg string) {
	s.respondJSON(w, status, ErrorResponse{Error: msg})
}

// HandleDevTriggerPart2 fast-forwards the state into Part 2 for verification.
func (s *Server) HandleDevTriggerPart2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	nightmare := r.URL.Query().Get("nightmare") == "true"
	if nightmare {
		s.game.JumpToPart2(true)
	} else {
		s.game.TriggerPart2ForDev()
	}
	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleNightmareTriggerPart2 fast-forwards directly to Nightmare Part 2 battle.
func (s *Server) HandleNightmareTriggerPart2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.game.JumpToPart2(true)
	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}

// HandleDevTriggerEnding fast-forwards directly to the ending choice/narrative.
func (s *Server) HandleDevTriggerEnding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	nightmare := r.URL.Query().Get("nightmare") == "true"
	choice := r.URL.Query().Get("choice")
	s.game.TriggerEndingForDev(choice, nightmare)
	s.respondJSON(w, http.StatusOK, s.game.GetSnapshot())
}
