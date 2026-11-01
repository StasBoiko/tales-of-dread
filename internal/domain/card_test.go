package domain

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		hand     []*Card
		expected int
	}{
		{
			name: "Standard cards",
			hand: []*Card{
				{BaseValue: 2, IsAce: false},
				{BaseValue: 7, IsAce: false},
			},
			expected: 9,
		},
		{
			name: "Single Ace as 11",
			hand: []*Card{
				{BaseValue: 11, IsAce: true},
				{BaseValue: 9, IsAce: false},
			},
			expected: 20,
		},
		{
			name: "Single Ace reduced to 1 on potential bust",
			hand: []*Card{
				{BaseValue: 11, IsAce: true},
				{BaseValue: 9, IsAce: false},
				{BaseValue: 5, IsAce: false},
			},
			expected: 15,
		},
		{
			name: "Two Aces",
			hand: []*Card{
				{BaseValue: 11, IsAce: true},
				{BaseValue: 11, IsAce: true},
			},
			expected: 12, // 11 + 1
		},
		{
			name: "Three Aces and a Nine",
			hand: []*Card{
				{BaseValue: 11, IsAce: true},
				{BaseValue: 11, IsAce: true},
				{BaseValue: 11, IsAce: true},
				{BaseValue: 9, IsAce: false},
			},
			expected: 12, // 9 + 1 + 1 + 1
		},
		{
			name: "Caged card halved",
			hand: []*Card{
				{BaseValue: 10, IsAce: false, IsCaged: true},
				{BaseValue: 8, IsAce: false},
			},
			expected: 13, // 5 + 8
		},
		{
			name: "Caged low value card minimum 1",
			hand: []*Card{
				{BaseValue: 1, IsAce: false, IsCaged: true},
				{BaseValue: 5, IsAce: false},
			},
			expected: 6, // 1 + 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateScore(tt.hand)
			if score != tt.expected {
				t.Errorf("CalculateScore() = %d, expected %d", score, tt.expected)
			}
		})
	}
}

func TestDeckDrawAndReshuffle(t *testing.T) {
	deck := NewStandardDeck()
	if deck.Remaining() != 52 {
		t.Fatalf("expected 52 cards, got %d", deck.Remaining())
	}

	card, err := deck.Draw()
	if err != nil {
		t.Fatalf("unexpected error drawing card: %v", err)
	}
	if card == nil {
		t.Fatal("expected non-nil card")
	}
	if deck.Remaining() != 51 {
		t.Fatalf("expected 51 cards remaining, got %d", deck.Remaining())
	}
}
