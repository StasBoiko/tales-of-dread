package domain

import (
	"testing"
)

func TestGameLifecycle(t *testing.T) {
	game := NewGame()
	if game.Phase != PhaseTitle {
		t.Fatalf("expected phase %s, got %s", PhaseTitle, game.Phase)
	}

	err := game.StartGame()
	if err != nil {
		t.Fatalf("unexpected error starting game: %v", err)
	}

	if game.Phase != PhasePlayerTurn {
		t.Fatalf("expected phase %s after start, got %s", PhasePlayerTurn, game.Phase)
	}

	if len(game.PlayerHand) != 2 {
		t.Fatalf("expected player hand length 2, got %d", len(game.PlayerHand))
	}

	if len(game.CurrentBoss.Hand) != 1 {
		t.Fatalf("expected boss hand length 1, got %d", len(game.CurrentBoss.Hand))
	}

	snap := game.GetSnapshot()
	if snap.PlayerScore <= 0 {
		t.Fatalf("expected player score > 0, got %d", snap.PlayerScore)
	}
}

func TestGameHitAndStand(t *testing.T) {
	game := NewGame()
	if err := game.StartGame(); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	initialCards := len(game.PlayerHand)
	err := game.Hit()
	if err != nil {
		t.Fatalf("unexpected hit error: %v", err)
	}

	if len(game.PlayerHand) != initialCards+1 {
		t.Fatalf("expected %d cards after hit, got %d", initialCards+1, len(game.PlayerHand))
	}

	// If player did not bust, they can Stand
	if game.Phase == PhasePlayerTurn {
		err = game.Stand()
		if err != nil {
			t.Fatalf("unexpected stand error: %v", err)
		}

		if game.Phase != PhaseRoundOver && game.Phase != PhaseArtifactSelect && game.Phase != PhaseCampfire && game.Phase != PhaseGameOver && game.Phase != PhaseVictory {
			t.Fatalf("unexpected phase after stand: %s", game.Phase)
		}
	}
}

func TestGameCrit21Damage(t *testing.T) {
	game := NewGame()
	if err := game.StartGame(); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	// Force player hand to 21
	game.PlayerHand = []*Card{
		{BaseValue: 11, IsAce: true},
		{BaseValue: 10, IsAce: false},
	}
	game.PlayerScore = 21
	initialBossHP := game.CurrentBoss.CurrentHP

	err := game.Stand()
	if err != nil {
		t.Fatalf("unexpected stand error: %v", err)
	}

	if !game.IsCrit {
		t.Error("expected IsCrit to be true on natural 21")
	}

	expectedDamage := 3
	if game.LastDamageDealt != expectedDamage {
		t.Errorf("expected %d damage on crit, got %d", expectedDamage, game.LastDamageDealt)
	}

	if game.CurrentBoss.CurrentHP != initialBossHP-expectedDamage {
		t.Errorf("expected boss HP %d, got %d", initialBossHP-expectedDamage, game.CurrentBoss.CurrentHP)
	}
}

func TestArtifactSelection(t *testing.T) {
	game := NewGame()
	game.Phase = PhaseArtifactSelect
	game.OfferedArtifacts = []*Artifact{
		{ID: "tallow_charm", Name: "Tallow Charm"},
	}
	initialMax := game.MaxCandles

	err := game.SelectArtifact("tallow_charm")
	if err != nil {
		t.Fatalf("failed selecting artifact: %v", err)
	}

	if game.MaxCandles != initialMax+1 {
		t.Errorf("expected max candles to increase by 1, got %d", game.MaxCandles)
	}

	if game.BossIndex != 1 {
		t.Errorf("expected bossIndex to advance to 1, got %d", game.BossIndex)
	}
}

func TestCampfireRespite(t *testing.T) {
	game := NewGame()
	game.Phase = PhaseCampfire
	game.PlayerCandles = 2
	game.MaxCandles = 5
	game.CampfireHealed = false

	if err := game.CampfireHeal(); err != nil {
		t.Fatalf("unexpected error during campfire heal: %v", err)
	}

	if game.PlayerCandles != 3 {
		t.Errorf("expected candles to increase to 3, got %d", game.PlayerCandles)
	}
	if !game.CampfireHealed {
		t.Error("expected CampfireHealed to be true")
	}

	// Repeated heal must error
	if err := game.CampfireHeal(); err == nil {
		t.Error("expected error on secondary campfire heal")
	}

	// Leave campfire transitions to relic selection
	if err := game.LeaveCampfire(); err != nil {
		t.Fatalf("unexpected error leaving campfire: %v", err)
	}
	if game.Phase != PhaseArtifactSelect {
		t.Errorf("expected phase %s, got %s", PhaseArtifactSelect, game.Phase)
	}
	if len(game.OfferedArtifacts) == 0 {
		t.Error("expected offered artifacts after leaving campfire")
	}
}

func TestNightmareCampfireAndMode(t *testing.T) {
	game := NewGame()
	if err := game.StartNightmare(); err != nil {
		t.Fatalf("failed starting nightmare: %v", err)
	}
	if !game.IsNightmare {
		t.Error("expected IsNightmare to be true")
	}
	// Boss 1 in nightmare has MaxHP 4 (3 + 1)
	if game.CurrentBoss.MaxHP != 4 {
		t.Errorf("expected Boss 1 MaxHP 4 in nightmare, got %d", game.CurrentBoss.MaxHP)
	}

	// Test campfire skip penalty in nightmare
	game.Phase = PhaseCampfire
	game.CampfireHealed = false
	game.PlayerCandles = 3
	if err := game.LeaveCampfire(); err != nil {
		t.Fatalf("unexpected error leaving campfire: %v", err)
	}
	if game.PlayerCandles != 2 {
		t.Errorf("expected 2 candles after unhealed campfire leave in nightmare, got %d", game.PlayerCandles)
	}
}

func TestArchmageBurnAndDualHands(t *testing.T) {
	game := NewGame()
	if err := game.StartGame(); err != nil {
		t.Fatalf("failed starting game: %v", err)
	}

	// Fast forward to Archmage (BossIndex = 2)
	game.BossIndex = 2
	game.CurrentBoss = GetAllBosses()[2].Clone()
	game.PlayerCandles = 3
	game.PlayerHitsThisRound = 0
	game.BurnTriggeredThisRound = false
	game.Phase = PhasePlayerTurn
	game.PlayerHand = []*Card{{Rank: "2", Suit: "hearts", BaseValue: 2}}
	game.PlayerScore = 2

	// Perform 3 hits
	for i := 0; i < 3; i++ {
		// Replace deck with safe cards to avoid bust
		game.Deck = &Deck{
			cards: []*Card{{Rank: "2", Suit: "diamonds", BaseValue: 2}},
		}
		if err := game.Hit(); err != nil {
			t.Fatalf("hit %d failed: %v", i+1, err)
		}
	}

	if game.PlayerCandles != 2 {
		t.Errorf("expected 1 candle burnt upon 3rd hit, remaining: %d", game.PlayerCandles)
	}
	if !game.BurnTriggeredThisRound {
		t.Error("expected BurnTriggeredThisRound to be true")
	}

	// Test 1: Player beats hand 1 (19 > 18) but loses to hand 2 (19 < 20) -> Boss win
	game.Phase = PhasePlayerTurn
	game.PlayerScore = 19
	game.BossScore = 18
	game.BossSecondScore = 20
	game.PlayerHand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "9", BaseValue: 9}}
	game.CurrentBoss.Hand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "8", BaseValue: 8}}
	game.CurrentBoss.SecondHand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "10", BaseValue: 10}}
	prevCandles := game.PlayerCandles
	if err := game.Stand(); err != nil {
		t.Fatalf("stand failed: %v", err)
	}
	if game.RoundOutcome != "boss_win" {
		t.Errorf("expected boss_win when player fails to beat both hands, got %s", game.RoundOutcome)
	}
	if game.PlayerCandles != prevCandles-1 {
		t.Errorf("expected candle loss, got %d", game.PlayerCandles)
	}

	// Test 2: Player beats BOTH hands (20 > 18 and 20 > 19) -> Player win
	game.Phase = PhasePlayerTurn
	game.PlayerScore = 20
	game.BossScore = 18
	game.BossSecondScore = 19
	game.PlayerHand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "10", BaseValue: 10}}
	game.CurrentBoss.Hand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "8", BaseValue: 8}}
	game.CurrentBoss.SecondHand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "9", BaseValue: 9}}
	prevHp := game.CurrentBoss.CurrentHP
	if err := game.Stand(); err != nil {
		t.Fatalf("stand failed: %v", err)
	}
	if game.RoundOutcome != "player_win" {
		t.Errorf("expected player_win when player beats both hands, got %s", game.RoundOutcome)
	}
	if game.CurrentBoss.CurrentHP != prevHp-1 {
		t.Errorf("expected boss HP loss, got %d", game.CurrentBoss.CurrentHP)
	}
}

func TestResetToTitle(t *testing.T) {
	game := NewGame()
	if err := game.StartGame(); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	if game.Phase != PhasePlayerTurn {
		t.Fatalf("expected PhasePlayerTurn, got %s", game.Phase)
	}

	game.ResetToTitle()
	if game.Phase != PhaseTitle {
		t.Errorf("expected PhaseTitle after reset, got %s", game.Phase)
	}
	if game.PlayerCandles != 5 {
		t.Errorf("expected 5 candles restored, got %d", game.PlayerCandles)
	}
	if len(game.PlayerHand) != 0 {
		t.Errorf("expected empty player hand after reset, got %d cards", len(game.PlayerHand))
	}
}

func TestSilverMedallionSavesLife(t *testing.T) {
	tests := []struct {
		name       string
		bustDamage bool
	}{
		{name: "saves on bust lethal damage", bustDamage: true},
		{name: "saves on stand defeat lethal damage", bustDamage: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame()
			if err := game.StartGame(); err != nil {
				t.Fatalf("failed starting game: %v", err)
			}

			// Verify medallion is present and active
			if len(game.Artifacts) == 0 || game.Artifacts[0].ID != "silver_medallion" || !game.Artifacts[0].Active {
				t.Fatalf("expected active silver_medallion at index 0")
			}

			// Drop candles to 1 so the next hit/defeat is lethal
			game.PlayerCandles = 1

			if tt.bustDamage {
				// Force bust (>21)
				game.PlayerHand = []*Card{
					{Rank: "10", BaseValue: 10},
					{Rank: "10", BaseValue: 10},
				}
				game.Deck = &Deck{cards: []*Card{{Rank: "5", BaseValue: 5}}}
				if err := game.Hit(); err != nil {
					t.Fatalf("unexpected hit error: %v", err)
				}
			} else {
				// Stand against a higher dealer score
				game.PlayerHand = []*Card{{Rank: "7", BaseValue: 7}}
				game.PlayerScore = 7
				game.CurrentBoss.Hand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "9", BaseValue: 9}}
				game.CurrentBoss.StandThreshold = 15
				if err := game.Stand(); err != nil {
					t.Fatalf("unexpected stand error: %v", err)
				}
			}

			// Player should NOT be game over; medallion should have triggered and kept 1 candle
			if game.Phase == PhaseGameOver {
				t.Fatalf("expected player to survive with medallion protection, got PhaseGameOver")
			}
			if game.PlayerCandles != 1 {
				t.Errorf("expected 1 candle preserved by medallion, got %d", game.PlayerCandles)
			}
			if !game.MedallionTriggered {
				t.Errorf("expected MedallionTriggered to be true")
			}
			if game.Artifacts[0].Active {
				t.Errorf("expected silver_medallion to be inactive after triggering")
			}

			snap := game.GetSnapshot()
			if !snap.MedallionSavedLife {
				t.Errorf("expected MedallionSavedLife in snapshot to be true on trigger round")
			}

			// Second lethal damage: medallion is spent, player should perish
			game.PlayerCandles = 1
			game.PlayerHand = []*Card{{Rank: "7", BaseValue: 7}}
			game.PlayerScore = 7
			game.CurrentBoss.Hand = []*Card{{Rank: "10", BaseValue: 10}, {Rank: "9", BaseValue: 9}}
			game.CurrentBoss.CurrentHP = 3
			game.Phase = PhasePlayerTurn
			if err := game.Stand(); err != nil {
				t.Fatalf("unexpected second stand error: %v", err)
			}

			if game.Phase != PhaseGameOver {
				t.Errorf("expected PhaseGameOver on second lethal damage, got %s", game.Phase)
			}
			if game.PlayerCandles != 0 {
				t.Errorf("expected 0 candles on game over, got %d", game.PlayerCandles)
			}
		})
	}
}

func TestMoltenMedallionPart2Reroll(t *testing.T) {
	game := NewGame()
	game.Phase = PhasePart2Battle
	game.Artifacts = []*Artifact{GetSilverMedallionArtifact()}
	game.JumpToPart2(false)

	// Verify medallion mutated
	if game.Artifacts[0].ID != "molten_medallion" {
		t.Fatalf("expected molten_medallion, got %s", game.Artifacts[0].ID)
	}

	// Preset player dice
	game.Part2.PlayerDice = []int{1, 1, 1, 2, 3}
	game.Part2.PlayerCombo, game.Part2.PlayerScore = EvaluateDice(game.Part2.PlayerDice)

	// Use molten medallion on die at index 3 (value 2)
	err := game.Part2Relic("molten_medallion", 3, 0)
	if err != nil {
		t.Fatalf("unexpected error using molten_medallion: %v", err)
	}

	if !game.Part2.MoltenMedallionUsed {
		t.Errorf("expected MoltenMedallionUsed to be true")
	}
	if game.Artifacts[0].Active {
		t.Errorf("expected molten_medallion artifact Active to be false after use")
	}

	// Repeated usage must fail
	err = game.Part2Relic("molten_medallion", 1, 0)
	if err == nil {
		t.Errorf("expected error using molten_medallion second time")
	}
}
