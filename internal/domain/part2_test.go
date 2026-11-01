package domain

import (
	"testing"
)

func TestEvaluateDice(t *testing.T) {
	tests := []struct {
		name          string
		dice          []int
		expectedCombo DiceCombination
	}{
		{
			name:          "Witch Eye (Five of a kind)",
			dice:          []int{4, 4, 4, 4, 4},
			expectedCombo: WitchEye,
		},
		{
			name:          "Four of a Kind",
			dice:          []int{2, 6, 2, 2, 2},
			expectedCombo: FourOfAKind,
		},
		{
			name:          "Full House",
			dice:          []int{3, 5, 3, 5, 3},
			expectedCombo: FullHouse,
		},
		{
			name:          "Straight High (2-6)",
			dice:          []int{5, 2, 4, 6, 3},
			expectedCombo: StraightHigh,
		},
		{
			name:          "Straight Low (1-5)",
			dice:          []int{3, 1, 4, 5, 2},
			expectedCombo: StraightLow,
		},
		{
			name:          "Three of a Kind",
			dice:          []int{1, 6, 1, 4, 1},
			expectedCombo: ThreeOfAKind,
		},
		{
			name:          "Two Pair",
			dice:          []int{2, 3, 5, 2, 3},
			expectedCombo: TwoPair,
		},
		{
			name:          "One Pair",
			dice:          []int{6, 1, 4, 2, 6},
			expectedCombo: OnePair,
		},
		{
			name:          "High Die",
			dice:          []int{1, 3, 4, 2, 6},
			expectedCombo: HighDie,
		},
		{
			name:          "Four of a Kind (4 dice)",
			dice:          []int{3, 3, 3, 3},
			expectedCombo: FourOfAKind,
		},
		{
			name:          "Three of a Kind (4 dice)",
			dice:          []int{5, 2, 5, 5},
			expectedCombo: ThreeOfAKind,
		},
		{
			name:          "Two Pair (4 dice)",
			dice:          []int{4, 1, 4, 1},
			expectedCombo: TwoPair,
		},
		{
			name:          "One Pair (4 dice)",
			dice:          []int{6, 3, 2, 6},
			expectedCombo: OnePair,
		},
		{
			name:          "High Die (4 dice)",
			dice:          []int{1, 4, 5, 2},
			expectedCombo: HighDie,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			combo, score := EvaluateDice(tt.dice)
			if combo != tt.expectedCombo {
				t.Errorf("expected combo %v (%s), got %v (%s)", tt.expectedCombo, tt.expectedCombo.String(), combo, combo.String())
			}
			if score <= 0 {
				t.Errorf("expected positive score, got %d", score)
			}
		})
	}
}

func TestPart2LifecycleAndReroll(t *testing.T) {
	p2 := NewPart2State(nil)
	if !p2.Active {
		t.Error("expected Part2 to be active")
	}
	if p2.WitchHP != 6 {
		t.Errorf("expected Witch HP 6, got %d", p2.WitchHP)
	}
	if p2.RerollsLeft != 1 {
		t.Errorf("expected 1 reroll, got %d", p2.RerollsLeft)
	}

	// Lock 2 dice and reroll 3
	locked := []bool{true, true, false, false, false}
	kept0 := p2.PlayerDice[0]
	kept1 := p2.PlayerDice[1]

	err := p2.RerollUnlocked(locked)
	if err != nil {
		t.Fatalf("unexpected reroll error: %v", err)
	}

	if p2.PlayerDice[0] != kept0 || p2.PlayerDice[1] != kept1 {
		t.Errorf("locked dice changed! [0]=%d vs %d, [1]=%d vs %d", p2.PlayerDice[0], kept0, p2.PlayerDice[1], kept1)
	}
	if p2.RerollsLeft != 0 {
		t.Errorf("expected 0 rerolls remaining, got %d", p2.RerollsLeft)
	}

	// Second reroll must fail
	if err := p2.RerollUnlocked(locked); err == nil {
		t.Error("expected error on exhausted rerolls")
	}
}

func TestPart2MutatedRelics(t *testing.T) {
	p2 := NewPart2State(nil)

	// 1. Eye of truth
	if err := p2.UseEyeOfTruth(); err != nil {
		t.Fatalf("failed using eye of truth: %v", err)
	}
	if len(p2.RevealedWitchIndices) != 3 {
		t.Errorf("expected 3 revealed dice indices, got %d", len(p2.RevealedWitchIndices))
	}
	// Re-use must error
	if err := p2.UseEyeOfTruth(); err == nil {
		t.Error("expected error on re-using eye of truth in same round")
	}

	// 2. Ink Mark
	err := p2.UseInkMark(0, 6)
	if err != nil {
		t.Fatalf("failed using ink mark: %v", err)
	}
	if p2.PlayerDice[0] != 6 {
		t.Errorf("expected die[0] to be 6, got %d", p2.PlayerDice[0])
	}

	// 3. Sacrificial Flame
	initialRerolls := p2.RerollsLeft
	err = p2.UseSacrificialFlame(4)
	if err != nil {
		t.Fatalf("failed using sacrificial flame: %v", err)
	}
	if p2.RerollsLeft != initialRerolls+1 {
		t.Errorf("expected rerolls to increase by 1, got %d", p2.RerollsLeft)
	}
	if len(p2.PlayerDice) != 4 {
		t.Errorf("expected 4 player dice after sacrificial flame, got %d", len(p2.PlayerDice))
	}
	if !p2.SacrificialFlameActive {
		t.Error("expected SacrificialFlameActive to be true")
	}

	// In next round, 5 dice must return
	p2.NextDiceRound()
	if len(p2.PlayerDice) != 5 {
		t.Errorf("expected 5 player dice in next round, got %d", len(p2.PlayerDice))
	}
	if p2.SacrificialFlameActive {
		t.Error("expected SacrificialFlameActive to be reset to false in next round")
	}
	// Relic must remain consumed for the entire duel
	if !p2.SacrificialFlameUsed {
		t.Error("expected SacrificialFlameUsed to remain true for the duel")
	}
}

func TestPart2ResolveBluff(t *testing.T) {
	p2 := NewPart2State(nil)
	p2.IsWitchBluffing = true
	p2.WitchClaim = FullHouse
	p2.WitchCombo = HighDie
	p2.WitchClaimName = "Фулл-Хаус"
	p2.WitchComboName = "Старшая кость"

	candles := 5
	dmgWitch, dmgPlayer := p2.ResolveRound(true, &candles)
	if dmgWitch != 2 {
		t.Errorf("expected 2 damage to Witch for caught bluff, got %d", dmgWitch)
	}
	if dmgPlayer != 0 {
		t.Errorf("expected 0 damage to player, got %d", dmgPlayer)
	}
	if p2.RoundOutcome != "witch_bluff_caught" {
		t.Errorf("expected outcome 'witch_bluff_caught', got %s", p2.RoundOutcome)
	}
}

func TestPart2EndingChoice(t *testing.T) {
	p2 := NewPart2State(nil)
	p2.WitchHP = 0
	p2.GameOver = true

	if err := p2.ApplyEnding("invalid"); err == nil {
		t.Error("expected error for invalid ending")
	}

	if err := p2.ApplyEnding("fire"); err != nil {
		t.Fatalf("failed applying fire ending: %v", err)
	}
	if p2.ChosenEnding != "fire" {
		t.Errorf("expected chosen ending 'fire', got '%s'", p2.ChosenEnding)
	}
}
