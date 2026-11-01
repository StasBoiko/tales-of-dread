package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"tales-of-dread/internal/domain"
	"tales-of-dread/web"
)

func setupTestServer(t *testing.T) *Server {
	t.Helper()
	g := domain.NewGame()
	httpFS, rawFS, err := web.SubFS()
	if err != nil {
		t.Fatalf("failed loading web FS: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewServer(g, logger, httpFS, rawFS)
}

func TestHealthz(t *testing.T) {
	s := setupTestServer(t)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "Valid GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid DELETE request",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/healthz", nil)
			w := httptest.NewRecorder()

			s.Routes().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.method == http.MethodGet {
				var resp StatusResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed decoding healthz response: %v", err)
				}
				if resp.Status != "ok" {
					t.Errorf("expected status 'ok', got '%s'", resp.Status)
				}
			}
		})
	}
}

func TestAPIWorkflow(t *testing.T) {
	s := setupTestServer(t)
	mux := s.Routes()

	// 1. GET /api/state initially
	req := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/state expected 200, got %d", w.Code)
	}

	// 2. POST /api/start
	req = httptest.NewRequest(http.MethodPost, "/api/start", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/start expected 200, got %d", w.Code)
	}

	var state domain.GameState
	if err := json.Unmarshal(w.Body.Bytes(), &state); err != nil {
		t.Fatalf("failed decoding start response: %v", err)
	}
	if state.Phase != domain.PhasePlayerTurn {
		t.Errorf("expected phase %s, got %s", domain.PhasePlayerTurn, state.Phase)
	}

	// 3. POST /api/hit
	req = httptest.NewRequest(http.MethodPost, "/api/hit", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/hit expected 200, got %d", w.Code)
	}

	// 4. POST /api/stand (if still player turn)
	if s.game.Phase == domain.PhasePlayerTurn {
		req = httptest.NewRequest(http.MethodPost, "/api/stand", nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("POST /api/stand expected 200, got %d", w.Code)
		}
	}

	// 5. POST /api/reset
	req = httptest.NewRequest(http.MethodPost, "/api/reset", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/reset expected 200, got %d", w.Code)
	}
	var resetState domain.GameState
	if err := json.Unmarshal(w.Body.Bytes(), &resetState); err != nil {
		t.Fatalf("failed decoding reset state: %v", err)
	}
	if resetState.Phase != domain.PhaseTitle {
		t.Errorf("expected phase 'title' after reset, got %s", resetState.Phase)
	}
}

func TestSelectArtifactValidation(t *testing.T) {
	s := setupTestServer(t)
	mux := s.Routes()

	// Empty body
	req := httptest.NewRequest(http.MethodPost, "/api/select-artifact", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty artifactId, got %d", w.Code)
	}
}

func TestCampfireEndpoints(t *testing.T) {
	s := setupTestServer(t)
	mux := s.Routes()

	// Switch game to campfire phase
	s.game.Phase = domain.PhaseCampfire
	s.game.PlayerCandles = 2
	s.game.MaxCandles = 5
	s.game.CampfireHealed = false

	// POST /api/campfire-heal
	req := httptest.NewRequest(http.MethodPost, "/api/campfire-heal", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/campfire-heal expected 200, got %d", w.Code)
	}

	var state domain.GameState
	if err := json.Unmarshal(w.Body.Bytes(), &state); err != nil {
		t.Fatalf("failed decoding state: %v", err)
	}
	if !state.CampfireHealed {
		t.Error("expected CampfireHealed to be true in state snapshot")
	}

	// POST /api/leave-campfire
	req = httptest.NewRequest(http.MethodPost, "/api/leave-campfire", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/leave-campfire expected 200, got %d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &state); err != nil {
		t.Fatalf("failed decoding state: %v", err)
	}
	if state.Phase != domain.PhaseArtifactSelect {
		t.Errorf("expected phase %s, got %s", domain.PhaseArtifactSelect, state.Phase)
	}
}
