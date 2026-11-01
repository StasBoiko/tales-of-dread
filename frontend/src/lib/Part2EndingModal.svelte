<script>
  import { createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  export let chosenEnding = ''; // "fire" or "throne" or ""
  export let isNightmare = false;
  export let playerCandles = 5;
  export let maxCandles = 5;

  const dispatch = createEventDispatcher();

  async function selectEnding(choice) {
    if (chosenEnding) return;
    chosenEnding = choice;

    if (choice === 'fire') {
      audioEngine.playInfernoRoar();
    } else {
      audioEngine.playMedallionCurse();
    }

    try {
      await fetch('/api/part2/ending', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ choice })
      });
    } catch (e) {
      console.error('Failed to commit ending choice:', e);
    }
  }

  function handleRestart() {
    if (isNightmare) {
      dispatch('restart', { nightmare: false, reset: true });
    } else {
      dispatch('restart', { nightmare: true });
    }
  }
</script>

<div class="ending-viewport">
  <div class="epilogue-stable-card">
    {#if !chosenEnding}
      <!-- 1. Fate Choice View: Simple Single-Click Selection -->
      <div class="choice-view">
        <header class="climax-header">
          <span class="sub-badge">РАЗВЯЗКА ИСТОРИИ</span>
          <h2 class="climax-title">
            {isNightmare ? 'СУДЬБА ТЕНЕВОГО ДВОЙНИКА' : 'СУДЬБА МОРОВОЙ ВЕДЬМЫ'}
          </h2>
          <p class="climax-desc">
            {#if isNightmare}
              Обсидиановый алтарь расколот. Твой Теневой Двойник пал, и зеркало корчмы пошло трещинами. Реши исход схватки с самим собой:
            {:else}
              Обсидиановый стол расколот. Перед тобой на коленях замерла Алёна — испуганная, беззащитная, с глазами, полными слез из сажи. Реши ее и свою судьбу одиночным выбором:
            {/if}
          </p>
        </header>

        <div class="ending-choices-grid">
          <!-- Choice 1: Fire -->
          <button 
            type="button" 
            class="choice-card fire-choice"
            on:click={() => selectEnding('fire')}
          >
            <div class="card-icon">🔥</div>
            <span class="choice-heading">{isNightmare ? 'РАЗБИТЬ ЗЕРКАЛО' : 'ОПРОКИНУТЬ СВЕЧУ'}</span>
            <span class="choice-title">{isNightmare ? 'Пепел Иллюзий' : 'Пепел Корчмы'}</span>
            <p class="choice-desc">
              {#if isNightmare}
                Разбить зеркальную гладь посохом, развеяв тень в пепел и навсегда разорвав связь с кошмаром.
              {:else}
                Смазнуть череп-подсвечник со стола. Очищающий огонь обратит корчму и все проклятия в прах. Ты уйдешь на рассвете.
              {/if}
            </p>
            <div class="click-select-btn fire-btn">ВЫБРАТЬ ЭТОТ ПУТЬ</div>
          </button>

          <!-- Choice 2: Throne -->
          <button 
            type="button" 
            class="choice-card throne-choice"
            on:click={() => selectEnding('throne')}
          >
            <div class="card-icon">🗝️</div>
            <span class="choice-heading">{isNightmare ? 'ПРИНЯТЬ ТЕНЬ' : 'НАДЕТЬ МЕДАЛЬОН'}</span>
            <span class="choice-title">{isNightmare ? 'Владыка Кошмара' : 'Новый Хозяин Топей'}</span>
            <p class="choice-desc">
              {#if isNightmare}
                Впустить темную сущность внутрь себя и обрести вечную власть над бездной корчмы.
              {:else}
                Надеть серебряный амулет на Алёну. Ее проклятие перейдет тебе — ты займешь трон из костей во тьме.
              {/if}
            </p>
            <div class="click-select-btn throne-btn">ВЫБРАТЬ ЭТОТ ПУТЬ</div>
          </button>
        </div>
      </div>

    {:else}
      <!-- 2. Permanent Epilogue View: Single Stable Card with Artwork -->
      <div class="epilogue-view">
        <header class="epilogue-header">
          <span class="sub-badge">ИТОГ ВЕЛИКОЙ ИГРЫ</span>
          <h2 class="epilogue-title">
            {chosenEnding === 'fire' ? 'ФИНАЛ: ПЕПЕЛ КОРЧМЫ' : 'ФИНАЛ: НОВЫЙ ХОЗЯИН ТОПЕЙ'}
          </h2>
        </header>

        <!-- High-res Ending Illustration -->
        <div class="ending-artwork-frame">
          {#if chosenEnding === 'fire'}
            <!-- Illustration: Hero in grey-blue hat walking away at dawn from burning tavern ruins -->
            <svg viewBox="0 0 460 220" class="ending-svg" aria-label="Пепел Корчмы: путник уходит на рассвете">
              <defs>
                <linearGradient id="dawnSky" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#2a0d18" />
                  <stop offset="35%" stop-color="#7a2b1f" />
                  <stop offset="70%" stop-color="#d46a2a" />
                  <stop offset="100%" stop-color="#ffb852" />
                </linearGradient>
                <linearGradient id="fireRuins" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#ff4400" />
                  <stop offset="50%" stop-color="#8a1805" />
                  <stop offset="100%" stop-color="#260603" />
                </linearGradient>
                <linearGradient id="pathGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#2b1a12" />
                  <stop offset="100%" stop-color="#543622" />
                </linearGradient>
              </defs>

              <!-- Dawn Sky & Sun Glow -->
              <rect x="0" y="0" width="460" height="150" fill="url(#dawnSky)" />
              <circle cx="280" cy="115" r="45" fill="#ffe082" opacity="0.85" filter="blur(6px)" />

              <!-- Distant Forest Pine Silhouettes -->
              <polygon points="0,150 20,120 40,150 70,110 100,150 140,115 180,150 220,125 260,150 300,120 340,150 380,110 420,150 460,125 460,150" fill="#14070a" />

              <!-- Burning Tavern Ruins in background left -->
              <g class="burning-tavern">
                <!-- Collapsed Charred Timbers -->
                <polygon points="60,150 85,75 110,85 100,150" fill="#1f0906" stroke="#4a150c" stroke-width="1.5" />
                <polygon points="90,150 140,60 165,70 180,150" fill="#160604" stroke="#4a150c" stroke-width="1.5" />
                <polygon points="135,150 195,85 215,95 230,150" fill="#1f0906" stroke="#4a150c" stroke-width="1.5" />
                <!-- Blazing Flames & Embers -->
                <path d="M 80 150 Q 110 50 135 150 Q 160 40 185 150 Q 200 65 220 150 Z" fill="url(#fireRuins)" opacity="0.9" />
                <path d="M 105 150 Q 135 70 160 150 Q 180 80 195 150 Z" fill="#ffb833" opacity="0.95" />
                <!-- Rising Smoke & Embers -->
                <circle cx="120" cy="50" r="2.5" fill="#ffeaa7" />
                <circle cx="150" cy="35" r="3" fill="#ff7675" />
                <circle cx="175" cy="45" r="2" fill="#ffeaa7" />
                <circle cx="190" cy="25" r="2.5" fill="#fdcb6e" />
              </g>

              <!-- Earth & Winding Morning Path -->
              <polygon points="0,145 460,145 460,220 0,220" fill="#120c08" />
              <path d="M 280 145 Q 310 170 340 190 Q 370 210 400 220 L 460 220 L 460 145 Z" fill="url(#pathGrad)" opacity="0.75" />

              <!-- Hero Walking Away Down Trail (Seen from behind/side) -->
              <g class="hero-walking">
                <!-- Shadow -->
                <ellipse cx="340" cy="195" rx="20" ry="5" fill="rgba(0,0,0,0.6)" />
                <!-- Brown Leather Boots walking stance -->
                <path d="M 330 182 L 328 194 L 322 195 L 333 196 L 334 182 Z" fill="#52321b" stroke="#7a4f2e" stroke-width="1" />
                <path d="M 346 178 L 348 190 L 354 191 L 344 192 L 343 178 Z" fill="#52321b" stroke="#7a4f2e" stroke-width="1" />
                <!-- Dark Traveling Cloak with Gold Edges Billowing in Dawn Wind -->
                <path d="M 326 128 L 318 180 Q 338 184 358 178 L 350 128 Z" fill="#1b1d24" stroke="#3d4352" stroke-width="1.5" />
                <!-- Wooden Staff in Right Hand -->
                <line x1="362" y1="110" x2="368" y2="194" stroke="#5c3e21" stroke-width="3" stroke-linecap="round" />
                <circle cx="362" cy="110" r="3.5" fill="#8f6437" />
                <!-- Shoulders & Back -->
                <path d="M 324 125 Q 338 120 352 125 L 350 145 L 326 145 Z" fill="#292e3b" />
                <!-- Iconic Grey-Blue Wizard Hat -->
                <ellipse cx="338" cy="116" rx="18" ry="4" fill="#3a4859" stroke="#5e758f" stroke-width="1.2" />
                <polygon points="328,115 348,115 344,80 338,76" fill="#44556a" stroke="#68819e" stroke-width="1.2" />
              </g>
            </svg>
          {:else}
            <!-- Illustration: Hero in grey-blue hat with glowing eyes seated on bone throne with shadows -->
            <svg viewBox="0 0 460 220" class="ending-svg" aria-label="Новый Хозяин Топей: герой на костяном троне">
              <defs>
                <radialGradient id="throneVoid" cx="50%" cy="50%" r="55%">
                  <stop offset="0%" stop-color="#240c30" />
                  <stop offset="60%" stop-color="#120419" />
                  <stop offset="100%" stop-color="#040106" />
                </radialGradient>
                <linearGradient id="boneGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#f5e6d0" />
                  <stop offset="60%" stop-color="#bda084" />
                  <stop offset="100%" stop-color="#524233" />
                </linearGradient>
              </defs>

              <!-- Deep Abyssal Vault Background -->
              <rect x="0" y="0" width="460" height="220" fill="url(#throneVoid)" />

              <!-- Flanking Whispering Shadow Tendrils -->
              <g class="shadow-tendrils" opacity="0.6">
                <path d="M 20 220 Q 50 110 110 80 Q 90 140 100 220 Z" fill="#0c0210" />
                <path d="M 440 220 Q 410 110 350 80 Q 370 140 360 220 Z" fill="#0c0210" />
                <path d="M 60 220 Q 90 60 140 30 Q 110 120 120 220 Z" fill="#170621" />
                <path d="M 400 220 Q 370 60 320 30 Q 350 120 340 220 Z" fill="#170621" />
              </g>

              <!-- Massive Throne of Carved Ancient Bones, Horns & Dark Roots -->
              <g class="bone-throne">
                <!-- High Spiked Bone Backrest -->
                <polygon points="180,180 160,50 175,25 190,70" fill="url(#boneGrad)" stroke="#3a2d21" stroke-width="1.5" />
                <polygon points="280,180 300,50 285,25 270,70" fill="url(#boneGrad)" stroke="#3a2d21" stroke-width="1.5" />
                <polygon points="205,180 200,30 215,10 225,60" fill="url(#boneGrad)" stroke="#3a2d21" stroke-width="1.5" />
                <polygon points="255,180 260,30 245,10 235,60" fill="url(#boneGrad)" stroke="#3a2d21" stroke-width="1.5" />
                <polygon points="215,180 230,5 245,180" fill="url(#boneGrad)" stroke="#3a2d21" stroke-width="2" />
                <!-- Carved Skull at apex -->
                <circle cx="230" cy="45" r="14" fill="#dfcdb5" stroke="#4a3726" stroke-width="1.5" />
                <ellipse cx="225" cy="44" rx="3" ry="4" fill="#0c0210" />
                <ellipse cx="235" cy="44" rx="3" ry="4" fill="#0c0210" />
                <!-- Throne Armrests -->
                <rect x="165" y="130" width="22" height="60" rx="4" fill="#30163b" stroke="#7a3494" stroke-width="1.5" />
                <rect x="273" y="130" width="22" height="60" rx="4" fill="#30163b" stroke="#7a3494" stroke-width="1.5" />
                <!-- Skulls on Armrests -->
                <circle cx="176" cy="128" r="9" fill="#dfcdb5" stroke="#4a3726" stroke-width="1.2" />
                <circle cx="284" cy="128" r="9" fill="#dfcdb5" stroke="#4a3726" stroke-width="1.2" />
                <!-- Throne Base & Dais -->
                <polygon points="140,220 160,180 300,180 320,220" fill="#1b0824" stroke="#5c1f73" stroke-width="2" />
              </g>

              <!-- Hero Seated Confidently as the New Tavern Master -->
              <g class="hero-enthroned">
                <!-- Flowing dark cloak draped over throne -->
                <path d="M 175 140 L 155 210 L 305 210 L 285 140 Z" fill="#181a21" stroke="#3b4254" stroke-width="1.5" />
                <!-- Crossed legs in brown leather boots -->
                <path d="M 205 165 L 200 205 L 218 208 L 222 175 Z" fill="#52321b" stroke="#7a4f2e" stroke-width="1.2" />
                <path d="M 230 165 L 255 200 L 240 208 L 220 175 Z" fill="#52321b" stroke="#7a4f2e" stroke-width="1.2" />
                <!-- Torso & Arms resting on armrests -->
                <path d="M 200 110 L 175 135 L 182 145 L 205 125 Z" fill="#2b313d" stroke="#4a556b" stroke-width="1" />
                <path d="M 260 110 L 285 135 L 278 145 L 255 125 Z" fill="#2b313d" stroke="#4a556b" stroke-width="1" />
                <rect x="210" y="105" width="40" height="45" rx="3" fill="#1c212b" stroke="#3e485c" stroke-width="1.5" />
                <!-- Detailed Silver Medallion on Chest glowing with cursed runes -->
                <circle cx="230" cy="122" r="8" fill="#d9e6f2" stroke="#ffffff" stroke-width="1.5" filter="drop-shadow(0 0 6px #78d2ff)" />
                <polygon points="230,116 235,125 225,125" fill="#38a5ff" />

                <!-- Hand holding bone dice cup on armrest -->
                <polygon points="278,126 288,126 286,140 280,140" fill="#dfcdb5" stroke="#3a2d21" stroke-width="1" />

                <!-- Hero Face with Glowing Cyan Eyes from under Hat -->
                <circle cx="230" cy="94" r="13" fill="#08080a" />
                <!-- Glowing Blue Cyan Eyes -->
                <ellipse cx="225" cy="94" rx="3.5" ry="2" fill="#00e5ff" filter="drop-shadow(0 0 5px #00e5ff)" />
                <ellipse cx="235" cy="94" rx="3.5" ry="2" fill="#00e5ff" filter="drop-shadow(0 0 5px #00e5ff)" />

                <!-- Grey-Blue Pointed Wizard Hat with Outer Glow -->
                <ellipse cx="230" cy="85" rx="26" ry="6" fill="#384657" stroke="#607694" stroke-width="1.5" filter="drop-shadow(0 0 8px rgba(56, 165, 255, 0.4))" />
                <polygon points="216,84 244,84 236,44 230,40" fill="#44556a" stroke="#6e87a8" stroke-width="1.5" />
              </g>

              <!-- Atmospheric Glowing Runes Floating in Fog -->
              <text x="70" y="190" fill="#a44ae0" font-size="16" opacity="0.6">ᚲ</text>
              <text x="380" y="190" fill="#a44ae0" font-size="16" opacity="0.6">ᛉ</text>
              <text x="120" y="210" fill="#38a5ff" font-size="14" opacity="0.5">ᚠ</text>
              <text x="330" y="210" fill="#38a5ff" font-size="14" opacity="0.5">ᛏ</text>
            </svg>
          {/if}
        </div>

        <!-- Epilogue Text Content -->
        <div class="epilogue-narrative">
          {#if chosenEnding === 'fire'}
            {#if isNightmare}
              <p class="epilogue-para">
                Рассвет озаряет лесную тропу сквозь дымящиеся руины. Корчмы больше нет. Твой собственный теневой двойник обратился в пепел вместе с вековыми проклятыми стенами.
              </p>
              <p class="epilogue-quote">
                «Ты сжег не только корчму, но и первородную тьму внутри своей души. Двойник пал, и кошмар рассеялся. Но покуда бьется сердце путника, бездна не забывает долги...»
              </p>
            {:else}
              <p class="epilogue-para">
                Рассвет озаряет лесную тропу сквозь дымящиеся руины. Корчмы больше нет. Ты срываешь с груди расплавившийся серебряный амулет и бросаешь его в придорожную траву.
              </p>
              <p class="epilogue-quote">
                «Огонь пожрал грехи этой гиблой чащи. Тьма рассеялась, и рассветное солнце впервые за долгие века греет твою дорогу. Ты свободен... и абсолютно пуст.»
              </p>
            {/if}
          {:else}
            {#if isNightmare}
              <p class="epilogue-para">
                Ты поглотил своего теневого двойника и окончательно занял костяной трон. В твоей руке — обсидиановый стакан с костями. Границы между смертным человеком и вечным Хозяином корчмы стерты навсегда.
              </p>
              <p class="epilogue-quote">
                «Победив собственную тень, ты стал абсолютным владыкой Бездны. Твои глаза горят мертвенно-бирюзовым пламенем. За дверью раздаются робкие шаги... Твоя первая игра начинается прямо сейчас.»
              </p>
            {:else}
              <p class="epilogue-para">
                Ты восседаешь на троне из древних костей перед новым обсидиановым столом. В твоей руке — резной костяной стакан с костями. Сквозь туман за дубовой дверью слышны робкие, неуверенные шаги нового путника.
              </p>
              <p class="epilogue-quote">
                «Ты спас ее жизнь, но забрал ее вечность. Корчма не может пустовать. Твои глаза горят первозданным огнем Бездны... Твоя первая игра в качестве Хозяина вот-вот начнется.»
              </p>
            {/if}
          {/if}

          {#if isNightmare}
            <div class="nightmare-triumph-banner">
              <span class="triumph-crown">👑</span>
              <div class="triumph-content">
                <span class="triumph-title">ИСТИННЫЙ ВЛАДЫКА КОРЧМЫ</span>
                <span class="triumph-sub">РЕЖИМ КОШМАРА ПОКОРЁН НА 100%</span>
              </div>
            </div>
            <div class="nightmare-stats-box">
              <div class="stat-item"><span class="stat-icon">⚔️</span> ПОВЕРЖЕНО БОССОВ: <strong>4 / 4</strong></div>
              <div class="stat-item"><span class="stat-icon">🕯️</span> УЦЕЛЕВШИХ СВЕЧЕЙ: <strong>{playerCandles} / {maxCandles}</strong></div>
              <div class="stat-item"><span class="stat-icon">🗝️</span> СОБРАНО РЕЛИКВИЙ: <strong>2 / 2</strong></div>
            </div>
          {/if}
        </div>

        <footer class="epilogue-actions">
          <button type="button" class="restart-journey-btn" on:click={handleRestart}>
            <span>{isNightmare ? 'НАЧАТЬ ЗАНОВО С ПРОЛОГА 🔄' : 'НОВАЯ ИГРА+ (КОШМАР) 💀'}</span>
          </button>
        </footer>
      </div>
    {/if}
  </div>
</div>

<style>
  .ending-viewport {
    position: fixed;
    inset: 0;
    background: rgba(3, 1, 5, 0.95);
    backdrop-filter: blur(8px);
    z-index: 12000;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow-y: auto;
    padding: 20px;
    box-sizing: border-box;
    font-variant-numeric: slashed-zero tabular-nums;
  }

  /* Single Stable Card: Prevents flickering and card swaps */
  .epilogue-stable-card {
    background: radial-gradient(circle at 50% 30%, #1e131d 0%, #110713 70%, #060208 100%);
    border: 2px solid #75358a;
    border-radius: 12px;
    max-width: 780px;
    width: 100%;
    padding: 28px 32px;
    box-shadow: 
      0 16px 50px rgba(0, 0, 0, 0.98),
      0 0 35px rgba(180, 50, 220, 0.25);
    display: flex;
    flex-direction: column;
    animation: fadeIn 0.4s ease-out;
  }

  .sub-badge {
    font-size: 0.85rem;
    color: #e5be68;
    letter-spacing: 2px;
    text-transform: uppercase;
    font-weight: bold;
    display: block;
    margin-bottom: 4px;
    text-align: center;
  }

  .climax-header {
    text-align: center;
    margin-bottom: 20px;
  }

  .climax-title, .epilogue-title {
    font-family: var(--font-title);
    font-size: 2.2rem;
    color: #ffd700;
    margin: 0 0 10px;
    letter-spacing: 1.5px;
    text-shadow: 0 0 12px rgba(255, 215, 0, 0.4);
    text-align: center;
  }

  .climax-desc {
    font-size: 1.05rem;
    line-height: 1.5;
    color: #e5be68;
    margin: 0 auto;
    max-width: 660px;
    text-align: center;
  }

  .ending-choices-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    margin-top: 10px;
  }

  .choice-card {
    background: linear-gradient(180deg, #240d18 0%, #12040c 100%);
    border: 2px solid #6b243b;
    border-radius: 10px;
    padding: 22px 18px;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    cursor: pointer;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.9);
    transition: all 0.2s ease;
    user-select: none;
  }

  .fire-choice:hover {
    border-color: #ff6a00;
    box-shadow: 0 10px 30px rgba(255, 90, 0, 0.4);
    transform: translateY(-3px);
  }

  .throne-choice {
    background: linear-gradient(180deg, #1f0b29 0%, #0d0312 100%);
    border-color: #63237a;
  }

  .throne-choice:hover {
    border-color: #ba47db;
    box-shadow: 0 10px 30px rgba(186, 71, 219, 0.4);
    transform: translateY(-3px);
  }

  .card-icon {
    font-size: 2.8rem;
    margin-bottom: 8px;
  }

  .choice-heading {
    font-size: 0.85rem;
    letter-spacing: 1.5px;
    color: #ff9955;
    text-transform: uppercase;
    font-weight: bold;
  }

  .throne-choice .choice-heading {
    color: #dca3ff;
  }

  .choice-title {
    font-family: var(--font-title);
    font-size: 1.6rem;
    color: #ffd700;
    margin: 4px 0 10px;
    letter-spacing: 1px;
  }

  .choice-desc {
    font-size: 0.95rem;
    line-height: 1.45;
    color: #e5be68;
    flex-grow: 1;
    margin: 0 0 16px;
  }

  .click-select-btn {
    font-family: var(--font-title);
    font-size: 1.1rem;
    padding: 10px 18px;
    border-radius: 6px;
    letter-spacing: 1px;
    font-weight: bold;
    width: 100%;
    box-sizing: border-box;
    transition: all 0.2s ease;
  }

  .fire-btn {
    background: #6e1c07;
    border: 1px solid #ff7b2b;
    color: #fff;
  }

  .fire-choice:hover .fire-btn {
    background: #a82e0d;
    box-shadow: 0 0 15px rgba(255, 120, 40, 0.5);
  }

  .throne-btn {
    background: #47155e;
    border: 1px solid #ba47db;
    color: #fff;
  }

  .throne-choice:hover .throne-btn {
    background: #6c228f;
    box-shadow: 0 0 15px rgba(190, 70, 220, 0.5);
  }

  /* Artwork Frame */
  .ending-artwork-frame {
    width: 100%;
    height: 220px;
    border-radius: 8px;
    overflow: hidden;
    border: 2px solid #8432a1;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.95);
    margin: 12px 0 18px;
  }

  .ending-svg {
    width: 100%;
    height: 100%;
    display: block;
  }

  .epilogue-narrative {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 20px;
  }

  .epilogue-para {
    font-size: 1.12rem;
    line-height: 1.6;
    color: #f7e8d5;
    margin: 0;
  }

  .epilogue-quote {
    font-style: italic;
    font-size: 1.15rem;
    line-height: 1.6;
    color: #ffd700;
    background: rgba(35, 14, 45, 0.7);
    border-left: 3px solid #ffd700;
    padding: 12px 18px;
    border-radius: 0 8px 8px 0;
    margin: 4px 0 0;
  }

  .epilogue-actions {
    display: flex;
    justify-content: center;
    margin-top: 10px;
  }

  .restart-journey-btn {
    padding: 14px 34px;
    background: linear-gradient(180deg, #5c2075 0%, #290a36 100%);
    border: 2px solid #ffd700;
    border-radius: 8px;
    color: #ffd700;
    font-family: var(--font-title);
    font-size: 1.25rem;
    font-weight: 900;
    letter-spacing: 1.5px;
    cursor: pointer;
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.9);
    transition: all 0.2s ease;
  }

  .restart-journey-btn:hover {
    transform: translateY(-2px);
    background: #7a2b9c;
    color: #fff;
    box-shadow: 0 0 20px rgba(255, 215, 0, 0.4);
  }

  /* Nightmare Mode Triumph & Victory Stats */
  .nightmare-triumph-banner {
    display: flex;
    align-items: center;
    gap: 16px;
    background: linear-gradient(90deg, rgba(74, 43, 10, 0.85) 0%, rgba(120, 80, 20, 0.95) 50%, rgba(74, 43, 10, 0.85) 100%);
    border: 2px solid #ffd700;
    border-radius: 8px;
    padding: 12px 20px;
    margin-top: 10px;
    box-shadow: 0 0 24px rgba(255, 215, 0, 0.35);
    animation: goldShimmer 3s ease-in-out infinite alternate;
  }

  .triumph-crown {
    font-size: 2.2rem;
    filter: drop-shadow(0 0 8px rgba(255, 215, 0, 0.8));
  }

  .triumph-content {
    display: flex;
    flex-direction: column;
    text-align: left;
  }

  .triumph-title {
    font-family: var(--font-title);
    font-size: 1.25rem;
    font-weight: 900;
    letter-spacing: 1.5px;
    color: #fff2a8;
    text-shadow: 0 0 8px rgba(255, 215, 0, 0.6);
  }

  .triumph-sub {
    font-size: 0.85rem;
    font-weight: 700;
    letter-spacing: 2px;
    color: #ffd700;
  }

  .nightmare-stats-box {
    display: flex;
    justify-content: space-around;
    background: rgba(15, 8, 20, 0.85);
    border: 1px solid #75358a;
    border-radius: 8px;
    padding: 10px 14px;
    gap: 12px;
    flex-wrap: wrap;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.9rem;
    letter-spacing: 1px;
    color: #dfcdb5;
  }

  .stat-item strong {
    color: #ffd700;
    font-size: 1rem;
  }

  .stat-icon {
    font-size: 1.1rem;
  }

  @keyframes goldShimmer {
    0% {
      border-color: #ffd700;
      box-shadow: 0 0 16px rgba(255, 215, 0, 0.3);
    }
    100% {
      border-color: #fff099;
      box-shadow: 0 0 28px rgba(255, 215, 0, 0.6);
    }
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @media (max-width: 768px) {
    .ending-choices-grid {
      grid-template-columns: 1fr;
    }
    .epilogue-stable-card {
      padding: 16px;
    }
    .ending-artwork-frame {
      height: 180px;
    }
  }
</style>
