package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tales-of-dread/internal/domain"
)

func TestPart2HTTPWorkflow(t *testing.T) {
	s := setupTestServer(t)
	mux := s.Routes()

	// 1. Manually set game state to PhasePart2Twist
	s.game.Phase = domain.PhasePart2Twist
	s.game.Part2 = domain.NewPart2State(nil)

	// 2. POST /api/part2/start
	req := httptest.NewRequest(http.MethodPost, "/api/part2/start", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/part2/start expected 200, got %d", w.Code)
	}

	var state domain.GameState
	if err := json.Unmarshal(w.Body.Bytes(), &state); err != nil {
		t.Fatalf("failed decoding state: %v", err)
	}
	if state.Phase != domain.PhasePart2Battle {
		t.Errorf("expected phase %s, got %s", domain.PhasePart2Battle, state.Phase)
	}
	if state.Part2 == nil {
		t.Fatal("expected Part2 state to be non-nil")
	}

	// 3. POST /api/part2/roll
	rollBody := []byte(`{"locked": [true, false, false, false, false]}`)
	req = httptest.NewRequest(http.MethodPost, "/api/part2/roll", bytes.NewBuffer(rollBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/part2/roll expected 200, got %d", w.Code)
	}

	// 4. POST /api/part2/relic (eye_of_truth)
	relicBody := []byte(`{"relic": "eye_of_truth"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/part2/relic", bytes.NewBuffer(relicBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/part2/relic expected 200, got %d", w.Code)
	}

	// 5. POST /api/part2/call-bluff
	bluffBody := []byte(`{"call_bluff": true}`)
	req = httptest.NewRequest(http.MethodPost, "/api/part2/call-bluff", bytes.NewBuffer(bluffBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/part2/call-bluff expected 200, got %d", w.Code)
	}

	// 6. Test ending endpoint
	s.game.Phase = domain.PhasePart2Ending
	s.game.Part2.WitchHP = 0
	endingBody := []byte(`{"choice": "fire"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/part2/ending", bytes.NewBuffer(endingBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/part2/ending expected 200, got %d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &state); err != nil {
		t.Fatalf("failed decoding state: %v", err)
	}
	if state.Part2.ChosenEnding != "fire" {
		t.Errorf("expected chosen ending 'fire', got '%s'", state.Part2.ChosenEnding)
	}
}
