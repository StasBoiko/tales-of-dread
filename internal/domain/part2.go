package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
)

// DiceCombination represents the poker hand rank for five six-sided dice.
type DiceCombination int

const (
	// HighDie represents no matching dice; value determined by the highest single die.
	HighDie DiceCombination = iota
	// OnePair represents exactly two matching dice.
	OnePair
	// TwoPair represents two distinct pairs of matching dice.
	TwoPair
	// ThreeOfAKind represents three matching dice.
	ThreeOfAKind
	// StraightLow represents a consecutive sequence from 1 through 5.
	StraightLow
	// StraightHigh represents a consecutive sequence from 2 through 6.
	StraightHigh
	// FullHouse represents three of a kind combined with a distinct pair.
	FullHouse
	// FourOfAKind represents four matching dice.
	FourOfAKind
	// WitchEye represents all five dice sharing the exact same value.
	WitchEye
)

// String returns the human-readable Russian localized name of the poker combination.
func (dc DiceCombination) String() string {
	switch dc {
	case HighDie:
		return "Старшая кость"
	case OnePair:
		return "Пара"
	case TwoPair:
		return "Две Пары"
	case ThreeOfAKind:
		return "Тройка"
	case StraightLow:
		return "Малый Стрит — от 1 до 5"
	case StraightHigh:
		return "Большой Стрит — от 2 до 6"
	case FullHouse:
		return "Фулл-Хаус — три плюс два"
	case FourOfAKind:
		return "Каре — четыре одинаковых"
	case WitchEye:
		return "Ведьмино Око — пять одинаковых"
	default:
		return "Неизвестная рука"
	}
}

// MutatedArtifact represents a dark relic upgraded for the Part II bone ritual.
type MutatedArtifact struct {
	ID          string `json:"id"`
	OriginalID  string `json:"original_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Part2State holds the state and progression for Part II: Witch's Eye Dice Duel.
type Part2State struct {
	Active                 bool               `json:"active"`
	DialogueStep           int                `json:"dialogue_step"`
	PlayerDice             []int              `json:"player_dice"`
	LockedDice             []bool             `json:"locked_dice"`
	RerollsLeft            int                `json:"rerolls_left"`
	WitchDice              []int              `json:"witch_dice"`
	RevealedWitchIndices   []int              `json:"revealed_witch_indices"`
	WitchClaim             DiceCombination    `json:"witch_claim"`
	WitchClaimName         string             `json:"witch_claim_name"`
	IsWitchBluffing        bool               `json:"is_witch_bluffing"`
	WitchHP                int                `json:"witch_hp"`
	WitchMaxHP             int                `json:"witch_max_hp"`
	WitchPhase             int                `json:"witch_phase"`
	PlayerCombo            DiceCombination    `json:"player_combo"`
	PlayerComboName        string             `json:"player_combo_name"`
	PlayerScore            int                `json:"player_score"`
	WitchCombo             DiceCombination    `json:"witch_combo"`
	WitchComboName         string             `json:"witch_combo_name"`
	WitchScore             int                `json:"witch_score"`
	Round                  int                `json:"round"`
	RoundOutcome           string             `json:"round_outcome"`
	LastMessage            string             `json:"last_message"`
	RoundResolved          bool               `json:"round_resolved"`
	GameOver               bool               `json:"game_over"`
	ChosenEnding           string             `json:"chosen_ending"`
	MutatedArtifacts       []*MutatedArtifact `json:"mutated_artifacts"`
	MoltenMedallionUsed    bool               `json:"molten_medallion_used"`
	EyeOfTruthUsed         bool               `json:"eye_of_truth_used"`
	InkMarkUsed            bool               `json:"ink_mark_used"`
	SacrificialFlameUsed   bool               `json:"sacrificial_flame_used"`
	SacrificialFlameActive bool               `json:"sacrificial_flame_active"`
	PlaguePenalty          int                `json:"plague_penalty"`
	IsNightmare            bool               `json:"is_nightmare"`
	WitchName              string             `json:"witch_name"`
}

// RandomDie returns a cryptographically secure random integer between 1 and 6.
func RandomDie() int {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		return 1
	}
	return int(n.Int64()) + 1
}

// EvaluateDice evaluates a slice of 4 or 5 dice and returns the combination rank and tie-break score.
func EvaluateDice(dice []int) (DiceCombination, int) {
	nDice := len(dice)
	if nDice != 5 && nDice != 4 {
		return HighDie, 0
	}

	counts := make(map[int]int)
	for _, d := range dice {
		counts[d]++
	}

	type countPair struct {
		val   int
		count int
	}
	pairs := make([]countPair, 0, len(counts))
	for val, count := range counts {
		pairs = append(pairs, countPair{val: val, count: count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].val > pairs[j].val
		}
		return pairs[i].count > pairs[j].count
	})

	sortedDice := make([]int, nDice)
	copy(sortedDice, dice)
	sort.Ints(sortedDice)

	isStraightLow := false
	isStraightHigh := false
	if nDice == 5 {
		isStraightLow = sortedDice[0] == 1 && sortedDice[1] == 2 && sortedDice[2] == 3 && sortedDice[3] == 4 && sortedDice[4] == 5
		isStraightHigh = sortedDice[0] == 2 && sortedDice[1] == 3 && sortedDice[2] == 4 && sortedDice[3] == 5 && sortedDice[4] == 6
	}

	// Five of a kind (Witch Eye) - requires 5 dice
	if nDice == 5 && pairs[0].count == 5 {
		score := int(WitchEye)*1000000 + pairs[0].val*10000
		return WitchEye, score
	}

	// Four of a kind
	if pairs[0].count == 4 {
		kicker := 0
		if len(pairs) > 1 {
			kicker = pairs[1].val
		}
		score := int(FourOfAKind)*1000000 + pairs[0].val*10000 + kicker
		return FourOfAKind, score
	}

	// Full House (3 + 2) - requires 5 dice
	if nDice == 5 && pairs[0].count == 3 && pairs[1].count == 2 {
		score := int(FullHouse)*1000000 + pairs[0].val*10000 + pairs[1].val*100
		return FullHouse, score
	}

	// Straights - requires 5 dice
	if isStraightHigh {
		score := int(StraightHigh)*1000000 + 6*10000
		return StraightHigh, score
	}
	if isStraightLow {
		score := int(StraightLow)*1000000 + 5*10000
		return StraightLow, score
	}

	// Three of a kind
	if pairs[0].count == 3 {
		kicker1 := 0
		kicker2 := 0
		if len(pairs) > 1 {
			kicker1 = pairs[1].val
		}
		if len(pairs) > 2 {
			kicker2 = pairs[2].val
		}
		score := int(ThreeOfAKind)*1000000 + pairs[0].val*10000 + kicker1*100 + kicker2
		return ThreeOfAKind, score
	}

	// Two pair
	if len(pairs) >= 2 && pairs[0].count == 2 && pairs[1].count == 2 {
		highPair := pairs[0].val
		lowPair := pairs[1].val
		if lowPair > highPair {
			highPair, lowPair = lowPair, highPair
		}
		kicker := 0
		if len(pairs) > 2 {
			kicker = pairs[2].val
		}
		score := int(TwoPair)*1000000 + highPair*10000 + lowPair*100 + kicker
		return TwoPair, score
	}

	// One pair
	if pairs[0].count == 2 {
		pairVal := pairs[0].val
		kickerScore := 0
		multiplier := 100
		for i := 1; i < len(pairs); i++ {
			kickerScore += pairs[i].val * multiplier
			multiplier /= 10
		}
		score := int(OnePair)*1000000 + pairVal*10000 + kickerScore
		return OnePair, score
	}

	// High die
	score := int(HighDie) * 1000000
	multiplier := 10000
	for i := nDice - 1; i >= 0; i-- {
		score += sortedDice[i] * multiplier
		multiplier /= 10
	}
	return HighDie, score
}

// DecideWitchAction determines whether the Witch will tell the truth or bluff based on her hand.
func DecideWitchAction(witchDice []int, playerCombo DiceCombination) (DiceCombination, bool) {
	trueCombo, _ := EvaluateDice(witchDice)

	// If the Witch already has a dominant hand (Full House or higher), she speaks true
	if trueCombo >= FullHouse || trueCombo > playerCombo {
		// 80% chance truth, 20% taunt bluff
		n, _ := rand.Int(rand.Reader, big.NewInt(100))
		if n.Int64() < 80 {
			return trueCombo, false
		}
	}

	// If Witch has a weak hand (HighDie, OnePair, TwoPair), she bluffs ~55% of the time
	bluffChance, _ := rand.Int(rand.Reader, big.NewInt(100))
	if bluffChance.Int64() < 55 {
		// Claim a combination 1-2 ranks higher than her actual hand
		boost, _ := rand.Int(rand.Reader, big.NewInt(2))
		claimed := trueCombo + DiceCombination(boost.Int64()+1)
		if claimed > WitchEye {
			claimed = WitchEye
		}
		return claimed, true
	}

	return trueCombo, false
}

// NewPart2State initializes a fresh Part II battle state (6 HP normal, 9 HP nightmare).
func NewPart2State(ownedArtifacts []*Artifact, isNightmare ...bool) *Part2State {
	nightmare := len(isNightmare) > 0 && isNightmare[0]
	witchName := "Алёна"
	if nightmare {
		witchName = "Теневой Двойник"
	}

	// Guarantee that the obsidian altar always provides the full set of mutated relics
	mutated := []*MutatedArtifact{
		{
			ID:          "molten_medallion",
			Name:        "Раскалённый медальон",
			Description: "Связь с ведьмой жжёт грудь. Позволяет один раз перебросить любую кость в дуэли.",
			Icon:        "🪬",
		},
		{
			ID:          "eye_of_truth",
			Name:        "Око Истины",
			Description: "Позволяет подсмотреть 3 закрытые кости противника до принятия решения.",
			Icon:        "👁️",
		},
		{
			ID:          "ink_mark",
			Name:        "Чернильная Метка",
			Description: "Позволяет перевернуть один любой кубик на выбранную грань от 1 до 6.",
			Icon:        "🖋️",
		},
		{
			ID:          "sacrificial_flame",
			Name:        "Жертвенное Пламя",
			Description: "Сжигает один кубик, даруя бесплатный переброс оставшихся с удвоением счета.",
			Icon:        "🔥",
		},
	}

	maxHP := 6
	if nightmare {
		maxHP = 9
	}

	p2 := &Part2State{
		Active:               true,
		DialogueStep:         0,
		PlayerDice:           []int{1, 1, 1, 1, 1},
		LockedDice:           []bool{false, false, false, false, false},
		RerollsLeft:          1,
		WitchDice:            []int{1, 1, 1, 1, 1},
		RevealedWitchIndices: []int{},
		WitchHP:              maxHP,
		WitchMaxHP:           maxHP,
		WitchPhase:           1,
		Round:                1,
		MutatedArtifacts:     mutated,
		LastMessage:          "Обряд на Костях начался. Брось резные кости на обсидиановый алтарь.",
		IsNightmare:          nightmare,
		WitchName:            witchName,
	}
	p2.rollInitialDice()
	return p2
}

func (p *Part2State) rollInitialDice() {
	p.PlayerDice = make([]int, 5)
	p.WitchDice = make([]int, 5)
	p.LockedDice = make([]bool, 5)
	p.RerollsLeft = 1
	p.SacrificialFlameActive = false
	p.RevealedWitchIndices = []int{}
	p.RoundResolved = false
	p.PlaguePenalty = 0

	for i := 0; i < 5; i++ {
		p.PlayerDice[i] = RandomDie()
		p.WitchDice[i] = RandomDie()
	}

	p.PlayerCombo, p.PlayerScore = EvaluateDice(p.PlayerDice)
	p.PlayerComboName = p.PlayerCombo.String()

	p.WitchCombo, p.WitchScore = EvaluateDice(p.WitchDice)
	p.WitchComboName = p.WitchCombo.String()

	p.WitchClaim, p.IsWitchBluffing = DecideWitchAction(p.WitchDice, p.PlayerCombo)
	p.WitchClaimName = p.WitchClaim.String()

	p.updateWitchPhase()
}

func (p *Part2State) updateWitchPhase() {
	if p.IsNightmare {
		if p.WitchHP > 6 {
			p.WitchPhase = 1
		} else if p.WitchHP > 3 {
			p.WitchPhase = 2
			p.PlaguePenalty = RandomDie() / 2
			if p.PlaguePenalty < 1 {
				p.PlaguePenalty = 1
			}
		} else {
			p.WitchPhase = 3
		}
		return
	}

	if p.WitchHP > 4 {
		p.WitchPhase = 1 // 100% - 66%: Voice of the Beloved
	} else if p.WitchHP > 2 {
		p.WitchPhase = 2 // 66% - 33%: Plague Fiend
		// In phase 2, plague penalty can subtract a die's value
		p.PlaguePenalty = RandomDie() / 2
		if p.PlaguePenalty < 1 {
			p.PlaguePenalty = 1
		}
	} else {
		p.WitchPhase = 3 // 33% - 0%: Abyss Agony
	}
}

// RerollUnlocked re-rolls all dice that the player has not locked.
func (p *Part2State) RerollUnlocked(locked []bool) error {
	if p.RerollsLeft <= 0 {
		return errors.New("no rerolls remaining for this hand")
	}
	if p.RoundResolved {
		return errors.New("current hand already resolved")
	}
	if len(locked) > 0 {
		p.LockedDice = make([]bool, len(p.PlayerDice))
		for i := 0; i < len(p.PlayerDice) && i < len(locked); i++ {
			p.LockedDice[i] = locked[i]
		}
	}

	for i := 0; i < len(p.PlayerDice); i++ {
		if !p.LockedDice[i] {
			p.PlayerDice[i] = RandomDie()
		}
	}
	p.RerollsLeft--

	p.PlayerCombo, p.PlayerScore = EvaluateDice(p.PlayerDice)
	p.PlayerComboName = p.PlayerCombo.String()
	p.LastMessage = fmt.Sprintf("Кости переброшены. Твоя текущая рука: %s.", p.PlayerComboName)
	return nil
}

// UseMoltenMedallion re-rolls a single targeted die randomly (1..6).
func (p *Part2State) UseMoltenMedallion(dieIndex int) error {
	if p.MoltenMedallionUsed {
		return errors.New("molten medallion already used in this duel")
	}
	if dieIndex < 0 || dieIndex >= len(p.PlayerDice) {
		return errors.New("invalid die index")
	}

	p.MoltenMedallionUsed = true
	oldVal := p.PlayerDice[dieIndex]
	newVal := RandomDie()
	p.PlayerDice[dieIndex] = newVal
	p.PlayerCombo, p.PlayerScore = EvaluateDice(p.PlayerDice)
	p.PlayerComboName = p.PlayerCombo.String()
	p.LastMessage = fmt.Sprintf("Раскалённый медальон перебросил кость [%d] -> [%d]! Новая рука: %s.", oldVal, newVal, p.PlayerComboName)
	return nil
}

// UseEyeOfTruth reveals three hidden dice of the opponent.
func (p *Part2State) UseEyeOfTruth() error {
	if p.EyeOfTruthUsed {
		return errors.New("eye of truth already used in this duel")
	}
	p.EyeOfTruthUsed = true
	p.RevealedWitchIndices = []int{0, 1, 2}
	targetName := p.WitchName
	if targetName == "" {
		targetName = "соперника"
	}
	p.LastMessage = fmt.Sprintf("Око Истины рассеивает морок: три кости %s равны [%d], [%d] и [%d].", targetName, p.WitchDice[0], p.WitchDice[1], p.WitchDice[2])
	return nil
}

// UseInkMark flips a specific die to a targeted face value (1..6).
func (p *Part2State) UseInkMark(dieIndex int, targetFace int) error {
	if p.InkMarkUsed {
		return errors.New("ink mark already used in this duel")
	}
	if dieIndex < 0 || dieIndex >= len(p.PlayerDice) {
		return errors.New("invalid die index")
	}
	if targetFace < 1 || targetFace > 6 {
		return errors.New("target face must be between 1 and 6")
	}

	p.InkMarkUsed = true
	p.PlayerDice[dieIndex] = targetFace
	p.PlayerCombo, p.PlayerScore = EvaluateDice(p.PlayerDice)
	p.PlayerComboName = p.PlayerCombo.String()
	p.LastMessage = fmt.Sprintf("Чернильная Метка изменила кость на [%d]! Новая рука: %s.", targetFace, p.PlayerComboName)
	return nil
}

// UseSacrificialFlame burns one die, leaving 4 dice for remainder of round and granting 1 extra reroll with x2 damage.
func (p *Part2State) UseSacrificialFlame(dieIndex int) error {
	if p.SacrificialFlameUsed {
		return errors.New("sacrificial flame already used in this duel")
	}
	if dieIndex < 0 || dieIndex >= len(p.PlayerDice) {
		return errors.New("invalid die index")
	}

	p.SacrificialFlameUsed = true
	p.SacrificialFlameActive = true
	p.RerollsLeft++

	burnedVal := p.PlayerDice[dieIndex]

	// Burn the chosen die out of player's hand, leaving remaining 4 dice
	p.PlayerDice = append(p.PlayerDice[:dieIndex], p.PlayerDice[dieIndex+1:]...)
	if len(p.LockedDice) > dieIndex {
		p.LockedDice = append(p.LockedDice[:dieIndex], p.LockedDice[dieIndex+1:]...)
	} else {
		p.LockedDice = make([]bool, len(p.PlayerDice))
	}

	p.PlayerCombo, p.PlayerScore = EvaluateDice(p.PlayerDice)
	p.PlayerComboName = p.PlayerCombo.String()
	p.LastMessage = fmt.Sprintf("Жертвенное Пламя пожрало кость [%d]! Осталось 4 кости, получен экстра-переброс с множителем урона x2.", burnedVal)
	return nil
}

// ResolveRound processes either calling the Witch's bluff or directly comparing poker combinations.
func (p *Part2State) ResolveRound(callBluff bool, playerCandles *int) (damageToWitch int, damageToPlayer int) {
	if p.RoundResolved {
		return 0, 0
	}
	p.RoundResolved = true

	candleDamageBase := 1
	if p.WitchPhase == 3 {
		candleDamageBase = 2 // Phase 3 double stakes
	}

	damageMultiplier := 1
	if p.SacrificialFlameActive {
		damageMultiplier = 2
	}

	if callBluff {
		// Player challenges the Witch
		if p.IsWitchBluffing {
			// Witch lied! She takes double damage (multiplied by sacrificial flame if active)
			damageToWitch = 2 * damageMultiplier
			p.WitchHP -= damageToWitch
			if p.WitchHP < 0 {
				p.WitchHP = 0
			}
			p.RoundOutcome = "witch_bluff_caught"
			p.LastMessage = fmt.Sprintf("Ты уличил соперника во лжи! %s заявляет: %s, но на костях лишь: %s. Нанесено урона: %d!", p.WitchName, p.WitchClaimName, p.WitchComboName, damageToWitch)
		} else {
			// Witch was telling the truth! Player is punished
			damageToPlayer = candleDamageBase
			*playerCandles -= damageToPlayer
			p.RoundOutcome = "witch_truth_punish"
			p.LastMessage = fmt.Sprintf("%s говорит правду! Ее рука истинна: %s. Ментальный шок гасит твою свечу! Потеряно свечей: %d.", p.WitchName, p.WitchComboName, damageToPlayer)
		}
	} else {
		// Showdown: Direct poker combination comparison
		effectivePlayerScore := p.PlayerScore
		if p.WitchPhase == 2 && p.PlaguePenalty > 0 {
			effectivePlayerScore -= p.PlaguePenalty * 10000
			if effectivePlayerScore < 0 {
				effectivePlayerScore = 0
			}
		}

		if effectivePlayerScore > p.WitchScore {
			// Player wins
			damageToWitch = 1 * damageMultiplier
			p.WitchHP -= damageToWitch
			if p.WitchHP < 0 {
				p.WitchHP = 0
			}
			p.RoundOutcome = "player_dice_win"
			p.LastMessage = fmt.Sprintf("Твоя комбинация %s превзошла кости соперника %s! Нанесено урона: %d.", p.PlayerComboName, p.WitchComboName, damageToWitch)
		} else if effectivePlayerScore < p.WitchScore {
			// Witch wins
			damageToPlayer = candleDamageBase
			*playerCandles -= damageToPlayer
			p.RoundOutcome = "witch_dice_win"
			p.LastMessage = fmt.Sprintf("Кости противника %s оказались сильнее твоих %s! Потеряно свечей: %d.", p.WitchComboName, p.PlayerComboName, damageToPlayer)
		} else {
			// Tie
			p.RoundOutcome = "tie"
			if p.IsNightmare {
				p.LastMessage = fmt.Sprintf("Ничья на комбинации %s! Заплати своей кровью за право продолжать бой.", p.PlayerComboName)
			} else {
				p.LastMessage = fmt.Sprintf("Ничья на комбинации %s! Алёна скалится: «Ничья?! Заплати горячей кровью!»", p.PlayerComboName)
			}
		}
	}

	p.updateWitchPhase()
	if p.WitchHP <= 0 {
		p.GameOver = true
		if p.IsNightmare {
			p.LastMessage = "Твой Теневой Двойник рассыпается в черный пепел! Проклятие зеркала разорвано..."
		} else {
			p.LastMessage = "Моровая Ведьма сломлена! Стол раскалывается, Алёна падает на колени перед тобой..."
		}
	}
	return damageToWitch, damageToPlayer
}

// NextDiceRound sets up the next duel round if neither combatant is defeated.
func (p *Part2State) NextDiceRound() {
	if p.WitchHP <= 0 {
		return
	}
	p.Round++
	p.rollInitialDice()
	p.LastMessage = fmt.Sprintf("Кон %d. Стаканы из черепов сотрясают кости...", p.Round)
}

// ApplyEnding sets the story climax choice ("fire" or "throne").
func (p *Part2State) ApplyEnding(choice string) error {
	if choice != "fire" && choice != "throne" {
		return errors.New("invalid ending choice, must be 'fire' or 'throne'")
	}
	p.ChosenEnding = choice
	return nil
}
