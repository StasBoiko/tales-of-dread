// Boss dialogue triggers and quotes for atmospheric speech bubbles and whispers

export const BOSS_QUOTES = {
  // Triggered when player draws cards in the high-danger zone (18-20 points)
  dangerHit: [
    '«Тяни еще... могила не бывает наполовину пустой!»',
    '«Рука дрожит, смертный?»',
    '«Жадность слаще страха, верно?..»',
    '«Один шаг в могилу... сделай его!»',
  ],

  // Triggered on exact 21 (Critical Strike)
  crit21: [
    '«Проклятье!.. Откуда у тебя эта масть?!»',
    '«Ты шулер или сам Сатана?!»',
    '«Пламя ослепляет меня!.. Как ты посмел?!»',
    '«Кости трещат!.. Это не может быть правдой!»',
  ],

  // Triggered when player busts (>21 points)
  bust: [
    '«Твоя кровь слаще меда!»',
    '«Минус палец, минус жизнь!»',
    '«Жадность сгубила тебя!»',
    '«Ха-ха-ха! Захлебнись своей спешкой!»',
  ],

  // Triggered when player is idle > 8 seconds
  idle: [
    '«Свеча коптит... ты тратишь мои последние секунды. Ходи!»',
    '«Свеча тает... и твое время тоже.»',
    '«Боишься шевельнуться, червь?»',
    '«Тьма уже в твоем затылке... ходи!»',
  ],

  // Specific Boss idle whispers / intro taunts
  bossSpecific: {
    the_lanterns: {
      intro: '«Скорбный свет укажет путь твоей душе... во тьму.»',
      hitReact: '«Стекло не треснет от твоих жалких уловок!»',
      taunt: '«Мой огонь выжжет твои карты дотла...»',
    },
    kingdom_of: {
      intro: '«Королевство бумаги ждет нового пленника...»',
      hitReact: '«Чернила текут! Листы рвутся! Ха-ха-ха!»',
      taunt: '«Каждая твоя мысль уже записана на этих стенах!»',
    },
    let_it_go: {
      intro: '«Врата корчмы ведут лишь в преисподнюю!»',
      hitReact: '«Пепел и магма поглотят тебя целиком!»',
      taunt: '«Пусть всё здесь сгорит к чертям собачьим!»',
    },
  },
};

/**
 * Returns a random quote from the specified category
 * @param {keyof typeof BOSS_QUOTES} category
 * @returns {string}
 */
export function getRandomQuote(category) {
  const list = BOSS_QUOTES[category];
  if (!Array.isArray(list) || list.length === 0) return '';
  const idx = Math.floor(Math.random() * list.length);
  return list[idx];
}
