import { writable } from 'svelte/store';

const translations = {
  ru: {
    'app.title': 'Чертово Очко: Страшные Сказки Корчмы',
    'ui.lang': 'Язык',
    'ui.ru': 'Русский',
    'ui.en': 'English',
    'ui.back': 'Назад',
    'ui.next': 'Далее',
    'ui.ok': 'ОК',
    'ui.cancel': 'Отмена',
    'ui.settings': 'Настройки',
    'ui.volume': 'Громкость',
    'phase.title': 'Титульный экран',
    'phase.player_turn': 'Ход игрока',
    'phase.boss_turn': 'Ход босса',
    'phase.round_over': 'Раунд завершен',
    'btn.hit': 'Взять карту',
    'btn.stand': 'Стоп',
    'btn.next_round': 'Следующий раунд',
    'btn.restart': 'Заново',
    'label.score': 'Очки',
    'label.candles': 'Свечи',
    'msg.welcome': 'Ты стоишь на пороге корчмы...',
    'msg.bust': 'Перебор!',
    'msg.crit': 'Критический успех!',
    'msg.victory': 'Победа!',
    'msg.defeat': 'Поражение...',
    'prologue.intro1': 'Они забрали Алёну на исходе осени. Без крика. Без следов. Лишь серебряный медальон остался в холодной дорожной пыли.',
    'prologue.intro2': 'Слухи вели в Проклятую Корчму на краю топей. Говорят, здешний хозяин не берет золота — он играет на годы, память и жизненную силу.',
    'prologue.intro3': 'У меня осталось пять свечей. Когда догорит последняя — я забуду ее лицо навсегда.',
    'prologue.enter': 'ВОЙТИ В ПРОКЛЯТУЮ КОРЧМУ',
    'prologue.skip': 'Пропустить титры ⏩'
  },
  en: {
    'app.title': "Devil's Twenty-One: Scary Tales of the Tavern",
    'ui.lang': 'Language',
    'ui.ru': 'Russian',
    'ui.en': 'English',
    'ui.back': 'Back',
    'ui.next': 'Next',
    'ui.ok': 'OK',
    'ui.cancel': 'Cancel',
    'ui.settings': 'Settings',
    'ui.volume': 'Volume',
    'phase.title': 'Title Screen',
    'phase.player_turn': "Player's Turn",
    'phase.boss_turn': "Boss's Turn",
    'phase.round_over': 'Round Over',
    'btn.hit': 'Hit',
    'btn.stand': 'Stand',
    'btn.next_round': 'Next Round',
    'btn.restart': 'Restart',
    'label.score': 'Score',
    'label.candles': 'Candles',
    'msg.welcome': 'You stand at the threshold of the tavern...',
    'msg.bust': 'Bust!',
    'msg.crit': 'Critical Hit!',
    'msg.victory': 'Victory!',
    'msg.defeat': 'Defeat...',
    'prologue.intro1': 'At autumn\'s edge, they took Alona silently. No cry. No traces. Only the silver medallion remained in the cold road dust.',
    'prologue.intro2': 'Rumors lead to the Cursed Tavern on the edge of the meadows. They say the owner there doesn\'t take gold — he plays for years, memory and life force.',
    'prologue.intro3': 'I have five candles left. When the last one burns out — I will forever forget her face.',
    'prologue.enter': 'ENTER THE CURSED TAVERN',
    'prologue.skip': 'Skip intro ⏩'
  },
};

let currentLang = 'ru';

function getLang() {
  return currentLang;
}

function setLang(lang) {
  if (translations[lang]) {
    currentLang = lang;
    localStorage.setItem('lang', lang);
    langStore.set(lang);
    return true;
  }
  return false;
}

function t(key) {
  const langTranslations = translations[currentLang];
  if (langTranslations && langTranslations[key]) {
    return langTranslations[key];
  }
  const fallback = translations['ru'];
  if (fallback && fallback[key]) {
    return fallback[key];
  }
  return key;
}

const langStore = writable(currentLang);

export { getLang, setLang, t, langStore };