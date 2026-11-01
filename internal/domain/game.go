package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

// Standard error definitions for the game domain.
var (
	ErrInvalidPhase      = errors.New("action not allowed in current phase")
	ErrGameOver          = errors.New("game is already over")
	ErrArtifactNotFound  = errors.New("selected artifact not offered")
	ErrInvalidPlayerHand = errors.New("cannot stand with empty hand")
)

// GamePhase identifies the current discrete step of the tavern duel state machine.
type GamePhase string

const (
	// PhaseTitle represents the tavern threshold before the duel begins.
	PhaseTitle GamePhase = "title"
	// PhasePlayerTurn represents player deliberation (Hit or Stand).
	PhasePlayerTurn GamePhase = "player_turn"
	// PhaseBossTurn represents the boss drawing cards and executing abilities.
	PhaseBossTurn GamePhase = "boss_turn"
	// PhaseRoundOver represents the conclusion of a single hand resolution.
	PhaseRoundOver GamePhase = "round_over"
	// PhaseCampfire represents an atmospheric respite between boss duels.
	PhaseCampfire GamePhase = "campfire"
	// PhaseArtifactSelect represents choosing a dark relic after felling a boss.
	PhaseArtifactSelect GamePhase = "artifact_select"
	// PhaseVictory represents final triumph over all three dread bosses.
	PhaseVictory GamePhase = "victory"
	// PhasePart2Twist represents the post-Archmage story cutscene and narrative dialogue with Alena.
	PhasePart2Twist GamePhase = "part2_twist"
	// PhasePart2Battle represents the dice duel on the obsidian altar against the Plague Witch.
	PhasePart2Battle GamePhase = "part2_battle"
	// PhasePart2Ending represents the climactic choice between burning the tavern or claiming the throne.
	PhasePart2Ending GamePhase = "part2_ending"
	// PhaseGameOver represents extinction of all five candle stubs.
	PhaseGameOver GamePhase = "game_over"
)

// GameState represents an immutable, serializable snapshot of the duel.
type GameState struct {
	Phase                  GamePhase   `json:"phase"`
	Part2                  *Part2State `json:"part2,omitempty"`
	PlayerCandles          int         `json:"playerCandles"`
	MaxCandles             int         `json:"maxCandles"`
	PlayerHand             []*Card     `json:"playerHand"`
	PlayerScore            int         `json:"playerScore"`
	BossIndex              int         `json:"bossIndex"`
	CurrentBoss            *Boss       `json:"currentBoss"`
	BossScore              int         `json:"bossScore"`
	BossSecondScore        int         `json:"bossSecondScore,omitempty"`
	PlayerHitsThisRound    int         `json:"playerHitsThisRound"`
	BurnTriggeredThisRound bool        `json:"burnTriggeredThisRound"`
	CampfireHealed         bool        `json:"campfireHealed"`
	Artifacts              []*Artifact `json:"artifacts"`
	OfferedArtifacts       []*Artifact `json:"offeredArtifacts"`
	Round                  int         `json:"round"`
	RoundOutcome           string      `json:"roundOutcome"`
	LastMessage            string      `json:"lastMessage"`
	LastDamageDealt        int         `json:"lastDamageDealt"`
	LastDamageTaken        int         `json:"lastDamageTaken"`
	IsCrit                 bool        `json:"isCrit"`
	ScreenShake            bool        `json:"screenShake"`
	LanternTrapActive      bool        `json:"lanternTrapActive"`
	SaltTriggeredInBattle  bool        `json:"saltTriggeredInBattle"`
	PeekedCard             *Card       `json:"peekedCard,omitempty"`
	IsNightmare            bool        `json:"isNightmare"`
	MedallionSavedLife     bool        `json:"medallionSavedLife,omitempty"`
}

// Game holds the mutable state for a roguelike run through the tavern.
type Game struct {
	mu                     sync.RWMutex
	Phase                  GamePhase
	PlayerCandles          int
	MaxCandles             int
	PlayerHand             []*Card
	PlayerScore            int
	BossIndex              int
	CurrentBoss            *Boss
	BossScore              int
	BossSecondScore        int
	PlayerHitsThisRound    int
	BurnTriggeredThisRound bool
	CampfireHealed         bool
	Deck                   *Deck
	Artifacts              []*Artifact
	OfferedArtifacts       []*Artifact
	Round                  int
	RoundOutcome           string
	LastMessage            string
	LastDamageDealt        int
	LastDamageTaken        int
	IsCrit                 bool
	ScreenShake            bool
	LanternTrapActive      bool
	SaltTriggeredInBattle  bool
	PeekedCard             *Card
	Part2                  *Part2State
	IsNightmare            bool
	MedallionTriggered     bool
	MedallionJustTriggered bool
}

// NewGame initializes a fresh roguelike run starting in PhaseTitle.
func NewGame() *Game {
	bosses := GetAllBosses()
	return &Game{
		Phase:                  PhaseTitle,
		PlayerCandles:          5,
		MaxCandles:             5,
		PlayerHand:             []*Card{},
		PlayerScore:            0,
		BossIndex:              0,
		CurrentBoss:            bosses[0].Clone(),
		BossScore:              0,
		BossSecondScore:        0,
		PlayerHitsThisRound:    0,
		BurnTriggeredThisRound: false,
		CampfireHealed:         false,
		Deck:                   NewStandardDeck(),
		Artifacts:              []*Artifact{GetSilverMedallionArtifact()},
		OfferedArtifacts:       []*Artifact{},
		Round:                  0,
		RoundOutcome:           "",
		LastMessage:            "Ты стоишь на пороге корчмы. За дверью завывает буря...",
		IsNightmare:            false,
		MedallionTriggered:     false,
		MedallionJustTriggered: false,
	}
}

// StartGameMode initiates a run in standard or nightmare mode.
func (g *Game) StartGameMode(nightmare bool) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	bosses := GetAllBosses()
	g.BossIndex = 0
	g.CurrentBoss = bosses[0].Clone()
	g.IsNightmare = nightmare
	if g.IsNightmare {
		g.CurrentBoss.MaxHP++
		g.CurrentBoss.CurrentHP = g.CurrentBoss.MaxHP
	}
	g.PlayerCandles = 5
	g.MaxCandles = 5
	g.Artifacts = []*Artifact{GetSilverMedallionArtifact()}
	g.OfferedArtifacts = []*Artifact{}
	g.MedallionTriggered = false
	g.MedallionJustTriggered = false
	g.Round = 0
	g.Deck = NewStandardDeck()
	g.SaltTriggeredInBattle = false
	g.LanternTrapActive = false

	return g.dealNewRoundLocked()
}

// StartGame initiates the run, facing the first boss and dealing initial cards.
func (g *Game) StartGame() error {
	return g.StartGameMode(false)
}

// StartNightmare initiates a Nightmare mode run.
func (g *Game) StartNightmare() error {
	return g.StartGameMode(true)
}

// ResetToTitle resets the entire run and returns the game to PhaseTitle.
func (g *Game) ResetToTitle() {
	g.mu.Lock()
	defer g.mu.Unlock()

	bosses := GetAllBosses()
	g.Phase = PhaseTitle
	g.PlayerCandles = 5
	g.MaxCandles = 5
	g.IsNightmare = false
	g.PlayerHand = []*Card{}
	g.PlayerScore = 0
	g.BossIndex = 0
	g.CurrentBoss = bosses[0].Clone()
	g.BossScore = 0
	g.BossSecondScore = 0
	g.PlayerHitsThisRound = 0
	g.BurnTriggeredThisRound = false
	g.CampfireHealed = false
	g.Deck = NewStandardDeck()
	g.Artifacts = []*Artifact{GetSilverMedallionArtifact()}
	g.OfferedArtifacts = []*Artifact{}
	g.MedallionTriggered = false
	g.MedallionJustTriggered = false
	g.Round = 0
	g.RoundOutcome = ""
	g.LastMessage = "Ты стоишь на пороге корчмы. За дверью завывает буря..."
	g.LastDamageDealt = 0
	g.LastDamageTaken = 0
	g.IsCrit = false
	g.ScreenShake = false
	g.LanternTrapActive = false
	g.SaltTriggeredInBattle = false
	g.PeekedCard = nil
	g.Part2 = nil
}

// Hit draws a card from the deck into the player hand and calculates outcomes.
// checkLethalDamageProtected checks if the player holds the active Silver Medallion and absorbs lethal candle damage.
func (g *Game) checkLethalDamageProtected() bool {
	if g.PlayerCandles <= 0 {
		if g.hasArtifactLocked("silver_medallion") && !g.MedallionTriggered {
			g.MedallionTriggered = true
			g.MedallionJustTriggered = true
			g.PlayerCandles = 1
			for _, art := range g.Artifacts {
				if art.ID == "silver_medallion" {
					art.Active = false
					art.Description = "Серебряный медальон Алёны. Металл потемнел и остыл, спасая твою душу от гибели. Сила оберега иссякла."
					break
				}
			}
			g.LastMessage = "Серебряный медальон Алёны вспыхивает холодным спасительным светом! Смертельный удар поглощен — последняя свеча продолжает гореть."
			return true
		}
		g.PlayerCandles = 0
		return false
	}
	return false
}

// Hit draws an additional card for the player from the deck.
func (g *Game) Hit() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.MedallionJustTriggered = false

	if g.Phase != PhasePlayerTurn {
		return ErrInvalidPhase
	}

	g.PlayerHitsThisRound++

	// Archmage burn curse: taking 3 cards in a round burns fingers for 1 candle
	if g.CurrentBoss.ID == BossAbyss && g.PlayerHitsThisRound >= 3 && !g.BurnTriggeredThisRound {
		g.BurnTriggeredThisRound = true
		g.PlayerCandles--
		g.LastDamageTaken = 1
		g.LastMessage = "Адский жар обжигает пальцы! Третья взятая карта за раунд отнимает 1 огарок свечи!"

		if g.PlayerCandles <= 0 {
			if !g.checkLethalDamageProtected() {
				g.Phase = PhaseGameOver
				g.LastMessage = "Все свечи погасли от адского пламени Архимага. Корчма поглотила твою душу."
				return nil
			}
		}
	}

	card, err := g.drawCardLocked()
	if err != nil {
		return fmt.Errorf("failed to draw card for player: %w", err)
	}

	if g.LanternTrapActive {
		card.IsBlinded = true
	}

	g.PlayerHand = append(g.PlayerHand, card)
	g.PlayerScore = CalculateScore(g.PlayerHand)
	g.updatePeekedCardLocked()

	// Check for bust condition (> 21)
	if g.PlayerScore > 21 {
		// Check Warding Salt artifact protection
		if g.hasArtifactLocked("warding_salt") && !g.SaltTriggeredInBattle {
			g.SaltTriggeredInBattle = true
			g.LastMessage = "Соль Оберега шипит на столе! Твой перебор рассеивается, фиксируя 21 очко."
			g.PlayerScore = 21
			return nil
		}

		// Player busts!
		g.LastDamageTaken = 1
		g.LastDamageDealt = 0
		g.PlayerCandles--
		g.IsCrit = false
		g.ScreenShake = (g.CurrentBoss.ID == BossAbyss)
		g.RoundOutcome = "bust"
		g.LastMessage = "Осечка! Твоя рука превысила 21 очко. Огарок свечи угасает со зловещим шипением..."

		if g.PlayerCandles <= 0 {
			if g.checkLethalDamageProtected() {
				g.Phase = PhaseRoundOver
				return nil
			}
			g.Phase = PhaseGameOver
			g.LastMessage = "Все свечи погасли. Корчма поглотила твою душу навеки."
			return nil
		}

		g.Phase = PhaseRoundOver
		return nil
	}

	if !(g.CurrentBoss.ID == BossAbyss && g.BurnTriggeredThisRound && g.PlayerHitsThisRound == 3) {
		g.LastMessage = fmt.Sprintf("Ты берешь %s. Твой текущий счет: %d.", card.Name, g.PlayerScore)
	}
	return nil
}

// Stand seals the player hand, executes the boss AI turn, and resolves the round.
func (g *Game) Stand() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.MedallionJustTriggered = false

	if g.Phase != PhasePlayerTurn {
		return ErrInvalidPhase
	}

	if len(g.PlayerHand) == 0 {
		return ErrInvalidPlayerHand
	}

	g.Phase = PhaseBossTurn
	g.executeBossTurnLocked()
	g.resolveRoundLocked()
	return nil
}

// NextRound prepares and deals a new hand against the current boss.
func (g *Game) NextRound() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.MedallionJustTriggered = false

	if g.Phase != PhaseRoundOver {
		return ErrInvalidPhase
	}

	return g.dealNewRoundLocked()
}

// SelectArtifact applies the chosen dark relic and transitions to the next tale.
func (g *Game) SelectArtifact(artifactID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhaseArtifactSelect {
		return ErrInvalidPhase
	}

	var chosen *Artifact
	for _, art := range g.OfferedArtifacts {
		if art.ID == artifactID {
			chosen = art
			break
		}
	}

	if chosen == nil {
		return ErrArtifactNotFound
	}

	g.Artifacts = append(g.Artifacts, chosen)
	g.OfferedArtifacts = []*Artifact{}

	// Apply immediate passive benefits
	if chosen.ID == "tallow_charm" {
		g.MaxCandles++
		g.PlayerCandles++
	}

	// Advance to next boss
	bosses := GetAllBosses()
	g.BossIndex++
	if g.BossIndex >= len(bosses) {
		g.Phase = PhaseVictory
		g.LastMessage = "Все три сказания завершены. Ты победил тьму корчмы!"
		return nil
	}

	g.CurrentBoss = bosses[g.BossIndex].Clone()
	if g.IsNightmare {
		g.CurrentBoss.MaxHP++
		g.CurrentBoss.CurrentHP = g.CurrentBoss.MaxHP
	}
	g.SaltTriggeredInBattle = false
	g.LanternTrapActive = false
	g.Deck = NewStandardDeck()

	return g.dealNewRoundLocked()
}

// GetSnapshot returns a thread-safe copy of the game state for presentation.
func (g *Game) GetSnapshot() *GameState {
	g.mu.RLock()
	defer g.mu.RUnlock()

	copiedHand := make([]*Card, len(g.PlayerHand))
	for i, c := range g.PlayerHand {
		copied := *c
		copiedHand[i] = &copied
	}

	var bossCopy *Boss
	if g.CurrentBoss != nil {
		bossCopy = g.CurrentBoss.Clone()
	}

	return &GameState{
		Phase:                  g.Phase,
		PlayerCandles:          g.PlayerCandles,
		MaxCandles:             g.MaxCandles,
		PlayerHand:             copiedHand,
		PlayerScore:            g.PlayerScore,
		BossIndex:              g.BossIndex,
		CurrentBoss:            bossCopy,
		BossScore:              g.BossScore,
		BossSecondScore:        g.BossSecondScore,
		PlayerHitsThisRound:    g.PlayerHitsThisRound,
		BurnTriggeredThisRound: g.BurnTriggeredThisRound,
		CampfireHealed:         g.CampfireHealed,
		Artifacts:              append([]*Artifact(nil), g.Artifacts...),
		OfferedArtifacts:       append([]*Artifact(nil), g.OfferedArtifacts...),
		Round:                  g.Round,
		RoundOutcome:           g.RoundOutcome,
		LastMessage:            g.LastMessage,
		LastDamageDealt:        g.LastDamageDealt,
		LastDamageTaken:        g.LastDamageTaken,
		IsCrit:                 g.IsCrit,
		ScreenShake:            g.ScreenShake,
		LanternTrapActive:      g.LanternTrapActive,
		SaltTriggeredInBattle:  g.SaltTriggeredInBattle,
		PeekedCard:             g.PeekedCard,
		Part2:                  g.Part2,
		IsNightmare:            g.IsNightmare,
		MedallionSavedLife:     g.MedallionJustTriggered,
	}
}

// dealNewRoundLocked initializes a new round with fresh initial cards.
func (g *Game) dealNewRoundLocked() error {
	g.Round++
	g.PlayerHand = []*Card{}
	g.CurrentBoss.Hand = []*Card{}
	g.CurrentBoss.SecondHand = []*Card{}
	g.IsCrit = false
	g.ScreenShake = false
	g.LastDamageDealt = 0
	g.LastDamageTaken = 0
	g.RoundOutcome = ""
	g.PlayerHitsThisRound = 0
	g.BurnTriggeredThisRound = false

	// Ensure deck has enough cards for a full deal
	if g.Deck.Remaining() < 14 {
		g.Deck = NewStandardDeck()
		if g.CurrentBoss.ID == BossPaper {
			g.Deck.InjectInkBlots(5)
		}
	} else if g.CurrentBoss.ID == BossPaper && g.Round == 1 {
		g.Deck.InjectInkBlots(5)
	}

	// Deal 2 cards to player, ensuring no starting hand busts (> 21)
	for attempt := 0; attempt < 10; attempt++ {
		g.PlayerHand = []*Card{}
		for i := 0; i < 2; i++ {
			c, err := g.drawCardLocked()
			if err != nil {
				return fmt.Errorf("failed dealing player card: %w", err)
			}
			g.PlayerHand = append(g.PlayerHand, c)
		}
		if CalculateScore(g.PlayerHand) <= 21 {
			break
		}
	}

	// Deal 1 starting card to boss main hand
	bossCard, err := g.drawCardLocked()
	if err != nil {
		return fmt.Errorf("failed dealing boss card: %w", err)
	}
	g.CurrentBoss.Hand = append(g.CurrentBoss.Hand, bossCard)

	// Archmage plays with two hands simultaneously
	if g.CurrentBoss.ID == BossAbyss {
		secondCard, err := g.drawCardLocked()
		if err != nil {
			return fmt.Errorf("failed dealing boss second hand card: %w", err)
		}
		g.CurrentBoss.SecondHand = append(g.CurrentBoss.SecondHand, secondCard)
	}

	// Apply boss passive / active special abilities
	g.applyBossPreRoundAbilitiesLocked()

	g.PlayerScore = CalculateScore(g.PlayerHand)
	g.BossScore = CalculateScore(g.CurrentBoss.Hand)
	g.BossSecondScore = CalculateScore(g.CurrentBoss.SecondHand)
	g.updatePeekedCardLocked()

	g.Phase = PhasePlayerTurn
	g.LastMessage = fmt.Sprintf("Раунд %d. Карты сданы на стол.", g.Round)
	return nil
}

// applyBossPreRoundAbilitiesLocked triggers boss-specific twists at round start.
func (g *Game) applyBossPreRoundAbilitiesLocked() {
	switch g.CurrentBoss.ID {
	case BossLantern:
		// Lantern trap: every 2 rounds blinds 1 of player's cards under lantern silhouette
		if g.Round%2 == 0 {
			g.LanternTrapActive = true
			if len(g.PlayerHand) > 0 {
				g.PlayerHand[0].IsBlinded = true
			}
			g.LastMessage = "Скорбный Фонарщик ставит фонарь в центр стола! Ослепляющий луч скрывает номинал карты."
		} else {
			g.LanternTrapActive = false
		}
	case BossPaper:
		// Paper Cage: if player lacks splattered_quill, ink cross out highest value card
		if !g.hasArtifactLocked("splattered_quill") && len(g.PlayerHand) > 0 {
			highest := g.PlayerHand[0]
			for _, c := range g.PlayerHand {
				if c.BaseValue > highest.BaseValue {
					highest = c
				}
			}
			highest.IsCaged = true
		}
	case BossAbyss:
		// Abyss Archmage passive aura
		g.ScreenShake = false
	}
}

// executeBossTurnLocked processes boss AI drawing logic.
func (g *Game) executeBossTurnLocked() {
	if g.CurrentBoss.IsStunned {
		g.CurrentBoss.IsStunned = false
		g.LastMessage = fmt.Sprintf("%s оглушен твоим прошлым ударом и не может добирать карты!", g.CurrentBoss.Name)
		g.BossScore = CalculateScore(g.CurrentBoss.Hand)
		g.BossSecondScore = CalculateScore(g.CurrentBoss.SecondHand)
		return
	}

	// Draw for primary hand
	for {
		score := CalculateScore(g.CurrentBoss.Hand)
		if score >= g.CurrentBoss.StandThreshold || score >= 21 {
			break
		}
		card, err := g.drawCardLocked()
		if err != nil {
			break
		}
		g.CurrentBoss.Hand = append(g.CurrentBoss.Hand, card)
	}
	g.BossScore = CalculateScore(g.CurrentBoss.Hand)

	// Draw for secondary hand if Archmage
	if g.CurrentBoss.ID == BossAbyss {
		for {
			score := CalculateScore(g.CurrentBoss.SecondHand)
			if score >= g.CurrentBoss.StandThreshold || score >= 21 {
				break
			}
			card, err := g.drawCardLocked()
			if err != nil {
				break
			}
			g.CurrentBoss.SecondHand = append(g.CurrentBoss.SecondHand, card)
		}
		g.BossSecondScore = CalculateScore(g.CurrentBoss.SecondHand)
	}
}

// resolveRoundLocked evaluates player vs boss scores and applies damage.
func (g *Game) resolveRoundLocked() {
	pScore := g.PlayerScore
	bScore := g.BossScore

	g.LastDamageDealt = 0
	g.LastDamageTaken = 0
	g.IsCrit = false

	// Critical Strike condition: Perfect 21
	if pScore == 21 {
		dmg := 3
		if g.hasArtifactLocked("gambler_bone") {
			dmg++
		}
		g.IsCrit = true
		g.LastDamageDealt = dmg
		g.CurrentBoss.CurrentHP -= dmg
		g.CurrentBoss.IsStunned = true
		g.RoundOutcome = "crit"
		g.LastMessage = fmt.Sprintf("КРИТИЧЕСКИЙ УДАР! Ровно 21 очко! Ты наносишь %d урона и оглушаешь босса!", dmg)
	} else if g.CurrentBoss.ID == BossAbyss {
		// Archmage Dual Hand Resolution: player must beat BOTH hands to win
		bScore1 := g.BossScore
		bScore2 := g.BossSecondScore

		beatHand1 := (bScore1 > 21) || (pScore > bScore1)
		beatHand2 := (bScore2 > 21) || (pScore > bScore2)

		if beatHand1 && beatHand2 {
			g.LastDamageDealt = 1
			g.CurrentBoss.CurrentHP--
			g.RoundOutcome = "player_win"
			g.LastMessage = fmt.Sprintf("Ты превзошел обе руки Архимага (%d против %d и %d)! 1 урон нанесен.", pScore, bScore1, bScore2)
		} else if (beatHand1 && pScore == bScore2) || (beatHand2 && pScore == bScore1) || (pScore == bScore1 && pScore == bScore2) {
			g.RoundOutcome = "tie"
			g.LastMessage = fmt.Sprintf("Ничья с рукой Архимага (%d против %d и %d). Морок сгущается...", pScore, bScore1, bScore2)
		} else {
			g.LastDamageTaken = 1
			g.PlayerCandles--
			g.RoundOutcome = "boss_win"
			g.LastMessage = fmt.Sprintf("Ты не смог превзойти обе руки Архимага (%d и %d против %d). Огарок свечи угасает...", bScore1, bScore2, pScore)
		}
	} else if bScore > 21 {
		// Boss busted
		g.LastDamageDealt = 1
		g.CurrentBoss.CurrentHP--
		g.RoundOutcome = "player_win"
		g.LastMessage = fmt.Sprintf("Босс перебрал (%d очков)! Ты наносишь 1 урон.", bScore)
	} else if pScore > bScore {
		// Player higher score
		g.LastDamageDealt = 1
		g.CurrentBoss.CurrentHP--
		g.RoundOutcome = "player_win"
		g.LastMessage = fmt.Sprintf("Твой счет (%d) превзошел счет врага (%d)! 1 урон нанесен.", pScore, bScore)
	} else if bScore > pScore {
		// Boss higher score
		g.LastDamageTaken = 1
		g.PlayerCandles--
		g.RoundOutcome = "boss_win"
		g.LastMessage = fmt.Sprintf("Босс превзошел тебя (%d против %d). Один огарок свечи угасает...", bScore, pScore)
	} else {
		// Tie
		g.RoundOutcome = "tie"
		g.LastMessage = fmt.Sprintf("Ничья (%d очков). Холодный взгляд сквозь мерцающий дым...", pScore)
	}

	// Check if Boss fell
	if g.CurrentBoss.CurrentHP <= 0 {
		g.CurrentBoss.CurrentHP = 0
		g.handleBossDefeatLocked()
		return
	}

	// Check if Player perished
	if g.PlayerCandles <= 0 {
		if g.checkLethalDamageProtected() {
			g.Phase = PhaseRoundOver
			return
		}
		g.Phase = PhaseGameOver
		g.LastMessage = "Твои свечи истлели. Холодная тьма поглотила твой рассудок."
		return
	}

	g.Phase = PhaseRoundOver
}

// handleBossDefeatLocked manages post-boss victory transitions to campfire or final victory.
func (g *Game) handleBossDefeatLocked() {
	// Soot Flask artifact restores 1 candle on boss victory
	if g.hasArtifactLocked("soot_flask") && g.PlayerCandles < g.MaxCandles {
		g.PlayerCandles++
	}

	// If this was the final boss (Boss 3: Archmage), transition to Part II twist cutscene
	if g.BossIndex == 2 {
		g.Phase = PhasePart2Twist
		for _, art := range g.Artifacts {
			if art.ID == "silver_medallion" {
				art.ID = "molten_medallion"
				art.Name = "Раскалённый медальон"
				art.Description = "Связь с ведьмой жжёт грудь. Позволяет один раз перебросить любую кость в дуэли."
				art.Icon = "🪬"
				art.Active = true
				break
			}
		}
		g.Part2 = NewPart2State(g.Artifacts, g.IsNightmare)
		if g.IsNightmare {
			g.LastMessage = "Архимаг Бездны обращен в пепел! Но из зеркальной глади зала выступает твоя собственная тень..."
		} else {
			g.LastMessage = "Архимаг обращен в пепел! Черный дым заволакивает зал, из темноты доносятся гулкие шаги..."
		}
		return
	}

	// Otherwise, transition to Campfire respite
	g.CampfireHealed = false
	g.Phase = PhaseCampfire
	if g.BossIndex == 0 {
		g.LastMessage = "Фонарщик повержен. Впереди мерцает костер привала..."
	} else {
		g.LastMessage = "Узник пал. Ты находишь убежище у тусклого костра..."
	}
}

// CampfireHeal restores 1 candle stub in standard mode, or prevents candle loss in Nightmare mode.
func (g *Game) CampfireHeal() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhaseCampfire {
		return ErrInvalidPhase
	}
	if g.CampfireHealed {
		return errors.New("already rested at this campfire")
	}

	g.CampfireHealed = true
	if g.IsNightmare {
		g.LastMessage = "Ты переждал бурю у тлеющих углей. Твоя свеча спасена от ледяного дыхания кошмара."
		return nil
	}

	if g.PlayerCandles < g.MaxCandles {
		g.PlayerCandles++
		g.LastMessage = "Ты вкусил теплой похлебки и эля у костра. Зажглась 1 свеча."
	} else {
		g.LastMessage = "Все твои свечи уже горят ярким пламенем."
	}
	return nil
}

// LeaveCampfire concludes the rest and transitions to dark relic selection.
// In Nightmare mode, skipping the 15-second rest extinguishes 1 candle and can cause death.
func (g *Game) LeaveCampfire() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhaseCampfire {
		return ErrInvalidPhase
	}

	if g.IsNightmare && !g.CampfireHealed {
		g.PlayerCandles--
		if g.PlayerCandles <= 0 {
			g.PlayerCandles = 0
			g.Phase = PhaseGameOver
			g.LastMessage = "Ты покинул костер без отдыха. Ледяной мрак кошмара задул твою последнюю свечу..."
			return nil
		}
		g.LastMessage = "Ты покинул костер без отдыха. Ледяной мрак кошмара задул 1 свечу!"
	}

	// Offer 3 random unowned artifacts
	catalog := GetCatalogArtifacts()
	var available []*Artifact
	for _, art := range catalog {
		if !g.hasArtifactLocked(art.ID) {
			available = append(available, art)
		}
	}

	n := len(available)
	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		j := int(jBig.Int64())
		available[i], available[j] = available[j], available[i]
	}

	pickCount := 3
	if len(available) < pickCount {
		pickCount = len(available)
	}

	g.OfferedArtifacts = available[:pickCount]
	g.Phase = PhaseArtifactSelect
	g.LastMessage = fmt.Sprintf("%s повержен! Выбери один из трех темных артефактов перед следующим сказанием.", g.CurrentBoss.Name)
	return nil
}

// drawCardLocked retrieves a card from the deck or reshuffles if exhausted.
func (g *Game) drawCardLocked() (*Card, error) {
	if g.Deck.Remaining() == 0 {
		g.Deck = NewStandardDeck()
	}
	return g.Deck.Draw()
}

// updatePeekedCardLocked updates the preview card if player possesses Raven's Feather.
func (g *Game) updatePeekedCardLocked() {
	if g.hasArtifactLocked("raven_feather") && g.Deck.Remaining() > 0 {
		peek, err := g.Deck.Peek()
		if err == nil {
			g.PeekedCard = peek
			return
		}
	}
	g.PeekedCard = nil
}

// hasArtifactLocked checks whether the player currently holds a specific artifact.
func (g *Game) hasArtifactLocked(id string) bool {
	for _, art := range g.Artifacts {
		if art.ID == id {
			return true
		}
	}
	return false
}

// StartPart2 transitions from the narrative twist into the dice duel against the Plague Witch.
func (g *Game) StartPart2() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhasePart2Twist {
		return ErrInvalidPhase
	}

	for _, art := range g.Artifacts {
		if art.ID == "silver_medallion" {
			art.ID = "molten_medallion"
			art.Name = "Раскалённый медальон"
			art.Description = "Связь с ведьмой жжёт грудь. Позволяет один раз перебросить любую кость в дуэли."
			art.Icon = "🪬"
			art.Active = true
			break
		}
	}

	if g.Part2 == nil {
		g.Part2 = NewPart2State(g.Artifacts, g.IsNightmare)
	}
	g.Phase = PhasePart2Battle
	g.LastMessage = "Обряд на Костях начался. Брось резные кости на обсидиановый алтарь."
	return nil
}

// JumpToPart2 directly initializes and transitions the game to Part 2 battle phase.
func (g *Game) JumpToPart2(nightmare bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.IsNightmare = nightmare
	g.Phase = PhasePart2Battle
	g.PlayerCandles = 5
	g.MaxCandles = 5
	g.Artifacts = []*Artifact{{
		ID:          "molten_medallion",
		Name:        "Раскалённый медальон",
		Description: "Связь с ведьмой жжёт грудь. Позволяет один раз перебросить любую кость в дуэли.",
		Icon:        "🪬",
		Active:      true,
	}}
	g.Part2 = NewPart2State(g.Artifacts, nightmare)
	if nightmare {
		g.LastMessage = "Кошмар: Твой Теневой Двойник поднимает обсидиановый стакан с костями."
	} else {
		g.LastMessage = "Обряд на Костях начался. Брось резные кости на обсидиановый алтарь."
	}
}

// Part2Reroll re-rolls unheld dice in the Witch's Eye duel.
func (g *Game) Part2Reroll(locked []bool) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhasePart2Battle || g.Part2 == nil {
		return ErrInvalidPhase
	}

	return g.Part2.RerollUnlocked(locked)
}

// Part2CallBluff resolves the hand by either challenging the Witch's claim or showdown.
func (g *Game) Part2CallBluff(callBluff bool) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhasePart2Battle || g.Part2 == nil {
		return ErrInvalidPhase
	}

	_, damageToPlayer := g.Part2.ResolveRound(callBluff, &g.PlayerCandles)
	if damageToPlayer > 0 {
		g.LastDamageTaken = damageToPlayer
	}

	if g.PlayerCandles <= 0 {
		g.PlayerCandles = 0
		g.Phase = PhaseGameOver
		g.LastMessage = "Последняя свеча твоей души угасла. Тьма поглотила тебя навсегда..."
		return nil
	}

	if g.Part2.WitchHP <= 0 {
		g.Phase = PhasePart2Ending
		g.LastMessage = "Моровая Ведьма сломлена! Стол раскалывается, Алёна падает на колени перед тобой..."
	}

	return nil
}

// Part2NextRound begins the next dice duel round against the Witch.
func (g *Game) Part2NextRound() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhasePart2Battle || g.Part2 == nil {
		return ErrInvalidPhase
	}

	g.Part2.NextDiceRound()
	return nil
}

// Part2Relic triggers a mutated relic power in the bone ritual.
func (g *Game) Part2Relic(relic string, dieIndex int, targetFace int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhasePart2Battle || g.Part2 == nil {
		return ErrInvalidPhase
	}

	switch relic {
	case "molten_medallion":
		err := g.Part2.UseMoltenMedallion(dieIndex)
		if err == nil {
			for _, art := range g.Artifacts {
				if art.ID == "molten_medallion" {
					art.Active = false
					break
				}
			}
		}
		return err
	case "eye_of_truth":
		return g.Part2.UseEyeOfTruth()
	case "ink_mark":
		return g.Part2.UseInkMark(dieIndex, targetFace)
	case "sacrificial_flame":
		return g.Part2.UseSacrificialFlame(dieIndex)
	default:
		return fmt.Errorf("unknown mutated relic: %s", relic)
	}
}

// Part2Ending locks in the player's fateful ending choice ("fire" or "throne").
func (g *Game) Part2Ending(choice string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Phase != PhasePart2Ending || g.Part2 == nil {
		return ErrInvalidPhase
	}

	return g.Part2.ApplyEnding(choice)
}

// TriggerPart2ForDev transitions the game directly to Part 2 for testing.
func (g *Game) TriggerPart2ForDev() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.BossIndex = 2
	g.Phase = PhasePart2Twist
	g.Artifacts = []*Artifact{
		{ID: "lantern_oil", Name: "Огарок Фонарщика"},
		{ID: "quill_feather", Name: "Перо Узника"},
		{ID: "archmage_ash", Name: "Пепел Архимага"},
	}
	g.Part2 = NewPart2State(g.Artifacts)
	g.LastMessage = "Архимаг Бездны обращен в пепел! Черный дым заволакивает зал, из темноты доносятся гулкие шаги..."
}

// TriggerEndingForDev transitions the game directly to the ending screen for testing.
func (g *Game) TriggerEndingForDev(choice string, nightmare bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.IsNightmare = nightmare
	g.Phase = PhasePart2Ending
	g.Part2 = NewPart2State(g.Artifacts, nightmare)
	g.Part2.GameOver = true
	g.Part2.ChosenEnding = choice
	if nightmare {
		g.LastMessage = "Кошмар преодолен: Твой Теневой Двойник растворяется в черном тумане."
	} else {
		g.LastMessage = "Моровая Ведьма сломлена! Судьба корчмы в твоих руках."
	}
}
