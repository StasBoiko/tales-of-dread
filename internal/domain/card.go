package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// ErrEmptyDeck is returned when attempting to draw from an exhausted deck.
var ErrEmptyDeck = errors.New("cannot draw card from empty deck")

// Suit represents the thematic card suit in the dread tavern.
type Suit string

const (
	// SuitSkulls represents mortality and raw combat power.
	SuitSkulls Suit = "skulls"
	// SuitBones represents ancient relics and resilient death.
	SuitBones Suit = "bones"
	// SuitCrows represents omen, darkness, and aerial whispers.
	SuitCrows Suit = "crows"
	// SuitChains represents bondage, guilt, and cold iron.
	SuitChains Suit = "chains"
)

// Card represents a single playing card on the dark tavern table.
type Card struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Suit      Suit   `json:"suit"`
	Rank      string `json:"rank"`
	BaseValue int    `json:"baseValue"`
	IsAce     bool   `json:"isAce"`
	IsCaged   bool   `json:"isCaged"`
	IsBlinded bool   `json:"isBlinded"`
	IsInkBlot bool   `json:"isInkBlot"`
}

// EffectiveValue returns the point value taking active status effects into account.
// Caged cards have their point values halved (minimum 1).
func (c *Card) EffectiveValue() int {
	if c.IsCaged {
		half := c.BaseValue / 2
		if half < 1 {
			return 1
		}
		return half
	}
	return c.BaseValue
}

// Deck represents a stack of dread playing cards.
type Deck struct {
	cards []*Card
}

// NewStandardDeck creates a freshly shuffled deck of 52 thematic cards.
func NewStandardDeck() *Deck {
	suits := []Suit{SuitSkulls, SuitBones, SuitCrows, SuitChains}
	ranks := []struct {
		name      string
		rank      string
		baseValue int
		isAce     bool
	}{
		{name: "Two", rank: "2", baseValue: 2, isAce: false},
		{name: "Three", rank: "3", baseValue: 3, isAce: false},
		{name: "Four", rank: "4", baseValue: 4, isAce: false},
		{name: "Five", rank: "5", baseValue: 5, isAce: false},
		{name: "Six", rank: "6", baseValue: 6, isAce: false},
		{name: "Seven", rank: "7", baseValue: 7, isAce: false},
		{name: "Eight", rank: "8", baseValue: 8, isAce: false},
		{name: "Nine", rank: "9", baseValue: 9, isAce: false},
		{name: "Ten", rank: "10", baseValue: 10, isAce: false},
		{name: "Knave", rank: "J", baseValue: 10, isAce: false},
		{name: "Crone", rank: "Q", baseValue: 10, isAce: false},
		{name: "Wraith", rank: "K", baseValue: 10, isAce: false},
		{name: "Dread Ace", rank: "A", baseValue: 11, isAce: true},
	}

	var cards []*Card
	for _, suit := range suits {
		for _, r := range ranks {
			id := fmt.Sprintf("%s-%s", suit, r.rank)
			name := fmt.Sprintf("%s of %s", r.name, suit)
			cards = append(cards, &Card{
				ID:        id,
				Name:      name,
				Suit:      suit,
				Rank:      r.rank,
				BaseValue: r.baseValue,
				IsAce:     r.isAce,
				IsCaged:   false,
				IsBlinded: false,
			})
		}
	}

	deck := &Deck{cards: cards}
	deck.Shuffle()
	return deck
}

// Shuffle randomizes the order of remaining cards in the deck using crypto/rand.
func (d *Deck) Shuffle() {
	n := len(d.cards)
	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		j := int(jBig.Int64())
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

// Draw removes and returns the top card from the deck.
func (d *Deck) Draw() (*Card, error) {
	if len(d.cards) == 0 {
		return nil, ErrEmptyDeck
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	// Create a shallow copy to prevent state contamination across rounds
	clone := *card
	return &clone, nil
}

// Peek returns the top card without removing it from the deck.
func (d *Deck) Peek() (*Card, error) {
	if len(d.cards) == 0 {
		return nil, ErrEmptyDeck
	}
	clone := *d.cards[0]
	return &clone, nil
}

// Remaining returns the count of cards currently available in the deck.
func (d *Deck) Remaining() int {
	return len(d.cards)
}

// InjectInkBlots inserts volatile ink-stained cards into the deck.
func (d *Deck) InjectInkBlots(count int) {
	for i := 0; i < count; i++ {
		valBig, err := rand.Int(rand.Reader, big.NewInt(10))
		val := int(valBig.Int64()) - 3
		if val == 0 || err != nil {
			val = 2
		}
		c := &Card{
			ID:        fmt.Sprintf("ink-blot-%d", i+1),
			Name:      "Чернильная Клякса",
			Suit:      "chains",
			Rank:      fmt.Sprintf("%+d", val),
			BaseValue: val,
			IsAce:     false,
			IsCaged:   false,
			IsBlinded: false,
			IsInkBlot: true,
		}
		d.cards = append(d.cards, c)
	}
	d.Shuffle()
}

// CalculateScore computes the optimal blackjack-style score for a given hand.
// Handles dynamic Ace values (11 or 1) and respects card curses like IsCaged and IsInkBlot.
func CalculateScore(hand []*Card) int {
	total := 0
	aceCount := 0

	for _, c := range hand {
		if c.IsAce {
			aceCount++
		} else {
			total += c.EffectiveValue()
		}
	}

	for i := 0; i < aceCount; i++ {
		// If counting this ace as 11 doesn't bust when all other remaining aces count as 1
		remainingAces := aceCount - 1 - i
		if total+11+remainingAces <= 21 {
			total += 11
		} else {
			total += 1
		}
	}

	if total < 0 {
		return 0
	}

	return total
}
