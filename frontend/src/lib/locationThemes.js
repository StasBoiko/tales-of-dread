// Configuration for dynamic locations and environmental parameters per boss encounter

export const LOCATION_THEMES = {
  the_lanterns: {
    id: 'the_lanterns',
    name: 'Холодный Лесной Тракт',
    subtitle: 'The Cold Forest Highway',
    tableType: 'pine_moss',
    tableDesc: 'Стол из сырых сосновых бревен, поросший влажным мхом и каплями янтарной смолы',
    ambientColor: 'rgba(12, 28, 25, 0.65)',
    candleColor: 'rgba(255, 180, 70, 0.28)',
    elements: {
      fog: true,
      rainOnGlass: true,
      pinesSilhouette: true,
      chains: false,
      flyingManuscripts: false,
      ashStorm: false,
      magmaCracks: false,
    },
    tableDecorations: [
      { type: 'moss', style: 'top: 15%; left: 8%; transform: rotate(-10deg);' },
      { type: 'resin', style: 'top: 60%; right: 12%;' },
      { type: 'bark_crack', style: 'bottom: 20%; left: 30%;' },
    ],
  },

  kingdom_of: {
    id: 'kingdom_of',
    name: 'Разрушенная Палата (Обитель Безумия)',
    subtitle: 'The Ruined Asylum Chamber',
    tableType: 'cracked_stone',
    tableDesc: 'Белый потрескавшийся каменный стол, истерзанный чернилами, разбитыми склянками и черным воском',
    ambientColor: 'rgba(18, 14, 24, 0.72)',
    candleColor: 'rgba(220, 190, 255, 0.22)',
    elements: {
      fog: false,
      rainOnGlass: false,
      pinesSilhouette: false,
      chains: true,
      flyingManuscripts: true,
      ashStorm: false,
      magmaCracks: false,
    },
    tableDecorations: [
      { type: 'ink_puddle', style: 'top: 25%; left: 14%;' },
      { type: 'shattered_phial', style: 'bottom: 25%; right: 18%;' },
      { type: 'stone_fissure', style: 'top: 50%; left: 45%;' },
    ],
  },

  let_it_go: {
    id: 'let_it_go',
    name: 'Преисподняя Корчмы',
    subtitle: 'The Tavern Underworld',
    tableType: 'sacrificial_altar',
    tableDesc: 'Обугленный жертвенный алтарь с раскаленными магматическими трещинами и снопом горящих искр',
    ambientColor: 'rgba(40, 10, 4, 0.78)',
    candleColor: 'rgba(255, 60, 0, 0.35)',
    elements: {
      fog: false,
      rainOnGlass: false,
      pinesSilhouette: false,
      chains: true,
      flyingManuscripts: false,
      ashStorm: true,
      magmaCracks: true,
    },
    tableDecorations: [
      { type: 'magma_glow', style: 'top: 40%; left: 20%;' },
      { type: 'cinders', style: 'top: 15%; right: 25%;' },
      { type: 'charred_rune', style: 'bottom: 30%; left: 60%;' },
    ],
  },
};

/**
 * Retrieves the location config for a given boss ID or default
 * @param {string} bossId
 */
export function getLocationConfig(bossId) {
  return LOCATION_THEMES[bossId] || LOCATION_THEMES.the_lanterns;
}
