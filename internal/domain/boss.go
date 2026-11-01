package domain

// BossID identifies each of the three dread tavern bosses.
type BossID string

const (
	// BossLantern represents The Sorrowful Tracker / Lantern Bearer.
	BossLantern BossID = "the_lanterns"
	// BossPaper represents The Asylum Sovereign / Prisoner of Paper Walls.
	BossPaper BossID = "kingdom_of"
	// BossAbyss represents The Possessed Innkeeper / Abyss Archmage.
	BossAbyss BossID = "let_it_go"
)

// Boss represents an adversary faced across the tavern table.
type Boss struct {
	ID             BossID   `json:"id"`
	Name           string   `json:"name"`
	Title          string   `json:"title"`
	AudioTrack     string   `json:"audioTrack"`
	MaxHP          int      `json:"maxHp"`
	CurrentHP      int      `json:"currentHp"`
	TableTheme     string   `json:"tableTheme"`
	AbilityName    string   `json:"abilityName"`
	AbilityDesc    string   `json:"abilityDesc"`
	StandThreshold int      `json:"standThreshold"`
	Hand           []*Card  `json:"hand"`
	SecondHand     []*Card  `json:"secondHand,omitempty"`
	IsStunned      bool     `json:"isStunned"`
	Dialogue       []string `json:"dialogue"`
}

// Artifact represents a dark passive relic earned after defeating a boss or held from the prologue.
type Artifact struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Active      bool   `json:"active"`
}

// GetSilverMedallionArtifact returns the starting silver medallion given by Alena.
func GetSilverMedallionArtifact() *Artifact {
	return &Artifact{
		ID:          "silver_medallion",
		Name:        "Серебряный медальон",
		Description: "Серебряный медальон Алёны. Хранит тепло её рук. Спасает от гибели, когда гаснет последняя свеча.",
		Icon:        "🥈",
		Active:      true,
	}
}

// Clone creates an isolated copy of the boss for duel initialization.
func (b *Boss) Clone() *Boss {
	copiedHand := make([]*Card, len(b.Hand))
	for i, c := range b.Hand {
		cardCopy := *c
		copiedHand[i] = &cardCopy
	}
	copiedSecondHand := make([]*Card, len(b.SecondHand))
	for i, c := range b.SecondHand {
		cardCopy := *c
		copiedSecondHand[i] = &cardCopy
	}
	return &Boss{
		ID:             b.ID,
		Name:           b.Name,
		Title:          b.Title,
		AudioTrack:     b.AudioTrack,
		MaxHP:          b.MaxHP,
		CurrentHP:      b.CurrentHP,
		TableTheme:     b.TableTheme,
		AbilityName:    b.AbilityName,
		AbilityDesc:    b.AbilityDesc,
		StandThreshold: b.StandThreshold,
		Hand:           copiedHand,
		SecondHand:     copiedSecondHand,
		IsStunned:      b.IsStunned,
		Dialogue:       append([]string(nil), b.Dialogue...),
	}
}

// GetAllBosses returns the ordered catalog of boss encounters.
func GetAllBosses() []*Boss {
	return []*Boss{
		{
			ID:             BossLantern,
			Name:           "Фонарщик",
			Title:          "The Lantern Bearer",
			AudioTrack:     "/audio/grimoras_crypt.mp3",
			MaxHP:          3,
			CurrentHP:      3,
			TableTheme:     "pine_boards",
			AbilityName:    "Слепящий луч",
			AbilityDesc:    "Раз в 2 хода фонарщик ослепляет одну из карт игрока черным силуэтом фонаря, скрывая номинал.",
			StandThreshold: 16,
			Hand:           []*Card{},
			SecondHand:     []*Card{},
			IsStunned:      false,
			Dialogue: []string{
				"Видишь искры на ледяном ветру?.. Это твои последние вздохи.",
				"Тьма шепчет среди сосен... и фонарь требует новую душу.",
				"Стекло треснет, и мы оба станем пеплом.",
			},
		},
		{
			ID:             BossPaper,
			Name:           "Узник",
			Title:          "The Paper Prisoner",
			AudioTrack:     "/audio/grimoras_crypt.mp3",
			MaxHP:          4,
			CurrentHP:      4,
			TableTheme:     "cave_damp",
			AbilityName:    "Чернильные кляксы",
			AbilityDesc:    "Босс подмешивает в колоду чернильные кляксы со случайным номиналом от -3 до +6.",
			StandThreshold: 17,
			Hand:           []*Card{},
			SecondHand:     []*Card{},
			IsStunned:      false,
			Dialogue: []string{
				"Свечи гаснут... а птицы за окном кричат о моем королевстве!",
				"Стены из тонкой бумаги не удержат то, что рвется наружу!",
				"Капля чернил решит твою судьбу быстрее, чем клинок палача.",
			},
		},
		{
			ID:             BossAbyss,
			Name:           "Архимаг",
			Title:          "The Archmage",
			AudioTrack:     "/audio/grimoras_crypt.mp3",
			MaxHP:          5,
			CurrentHP:      5,
			TableTheme:     "charred_oak",
			AbilityName:    "Двуручный вихрь и ожог",
			AbilityDesc:    "Играет двумя руками одновременно. Взяв 3 карты за раунд, игрок теряет 1 огарок свечи от ожога.",
			StandThreshold: 17,
			Hand:           []*Card{},
			SecondHand:     []*Card{},
			IsStunned:      false,
			Dialogue: []string{
				"Пусть этот мир летит к чертям в погребальном пламени!",
				"Ледяные объятия Бездны сомкнутся над твоей головой!",
				"Один неверный шаг — и корчма обратится в прах вместе с тобой!",
			},
		},
	}
}

// GetCatalogArtifacts returns all selectable dark relics in the roguelike pool.
func GetCatalogArtifacts() []*Artifact {
	return []*Artifact{
		{
			ID:          "raven_feather",
			Name:        "Воронье Перо",
			Description: "Шепот ворона позволяет подглядеть верхнюю карту колоды перед добором.",
			Icon:        "🪶",
		},
		{
			ID:          "warding_salt",
			Name:        "Соль Оберега",
			Description: "Первый перебор (>21) в каждой битве прощается и мягко фиксирует счет на 21.",
			Icon:        "🧂",
		},
		{
			ID:          "gambler_bone",
			Name:        "Костяная Кость Игрока",
			Description: "Идеальные 21 очко наносят дополнительную единицу урона боссу.",
			Icon:        "🎲",
		},
		{
			ID:          "soot_flask",
			Name:        "Склянка с Сажей",
			Description: "Победа над каждым боссом вновь зажигает один погасший огарок свечи.",
			Icon:        "🧪",
		},
		{
			ID:          "tallow_charm",
			Name:        "Амулет из Сала",
			Description: "Увеличивает максимальный запас жизни на 1 огарок свечи.",
			Icon:        "🕯️",
		},
		{
			ID:          "splattered_quill",
			Name:        "Чернильное Перо Бунтаря",
			Description: "Иммунитет к ослаблению карт 'Бумажной Клеткой'.",
			Icon:        "✒️",
		},
	}
}
