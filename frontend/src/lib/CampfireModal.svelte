<script>
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  export let bossIndex = 0; // 0 = after Boss 1, 1 = after Boss 2
  export let playerCandles = 3;
  export let maxCandles = 5;
  export let campfireHealed = false;
  export let isNightmare = false;

  const dispatch = createEventDispatcher();

  let isCooking = false;
  let cookProgress = 0;
  let cookInterval = null;
  let cookStatusText = 'Варится похлебка...';
  let healedCandleAnim = false;

  const monologues = [
    '«Старик шептал, что свет фонаря забирает память... Медальон Алёны жжет мне грудь. Держись, я близко.»',
    '«Я нашел обрывок ее кружевного платка в этой сырой яме. Но почему эти твари склоняли головы перед ней?»',
    '«Сквозь руины корчмы веет ледяной бездной... Пора положить конец этому кошмару.»'
  ];

  $: activeMonologue = monologues[bossIndex] || monologues[0];
  $: isCandlesFull = !isNightmare && playerCandles >= maxCandles;

  onMount(() => {
    // Duck music and activate crackling fire + crickets
    audioEngine.setDucking(true);
  });

  onDestroy(() => {
    if (cookInterval) clearInterval(cookInterval);
    audioEngine.setDucking(false);
  });

  function startCookingStew() {
    if (campfireHealed || isCooking || (!isNightmare && playerCandles >= maxCandles)) return;

    isCooking = true;
    cookProgress = 0;
    cookStatusText = 'Разжигается огонь под казаном...';

    const totalSeconds = 15;
    const intervalMs = 100;
    const step = 100 / (totalSeconds * (1000 / intervalMs));

    clearInterval(cookInterval);
    cookInterval = setInterval(async () => {
      cookProgress += step;

      if (cookProgress < 30) {
        cookStatusText = 'Разжигается огонь... Закипает вода в казане...';
      } else if (cookProgress < 65) {
        cookStatusText = 'Варится похлебка из лесных трав и кореньев...';
      } else if (cookProgress < 92) {
        cookStatusText = 'Пар окутывает привал... Силы возвращаются...';
      } else {
        cookStatusText = 'Похлебка готова! Пламя свечей разгорается вновь!';
      }

      if (cookProgress >= 100) {
        clearInterval(cookInterval);
        cookProgress = 100;
        await completeStewHealing();
      }
    }, intervalMs);
  }

  async function completeStewHealing() {
    // Procedural wick ignition & flame whoosh
    audioEngine.playWickIgnition();
    healedCandleAnim = true;

    try {
      const res = await fetch('/api/campfire-heal', { method: 'POST' });
      if (res.ok) {
        const data = await res.json();
        campfireHealed = true;
        dispatch('healed', data);
      }
    } catch (e) {
      console.error('Failed campfire heal:', e);
    } finally {
      isCooking = false;
      setTimeout(() => {
        healedCandleAnim = false;
      }, 1600);
    }
  }

  async function handleContinue() {
    if (cookInterval) clearInterval(cookInterval);
    audioEngine.playIronBolt();
    audioEngine.setDucking(false);

    try {
      const res = await fetch('/api/leave-campfire', { method: 'POST' });
      if (res.ok) {
        const data = await res.json();
        dispatch('leave', data);
      }
    } catch (e) {
      console.error('Failed to leave campfire:', e);
    }
  }
</script>

<div class="campfire-backdrop">
  <div class="campfire-container">
    <!-- Header title -->
    <div class="campfire-header">
      <span class="moon-icon">🌙</span>
      <h2 class="campfire-title">Ночной Привал</h2>
      <span class="moon-icon">🔥</span>
    </div>

    <!-- Large Central Scene: Hero + Fire + Cauldron on Tripod with Swirling Steam -->
    <div class="campfire-scene">
      <!-- Embers Emitter -->
      <div class="embers-emitter">
        {#each Array(10) as _, i}
          <span class="ember-spark" style="--i: {i}; --rnd: {0.15 + (i * 0.1)}s;"></span>
        {/each}
      </div>

      <!-- Steam puffs from cauldron -->
      <div class="steam-emitter">
        <span class="steam-cloud s1"></span>
        <span class="steam-cloud s2"></span>
        <span class="steam-cloud s3"></span>
      </div>

      <svg class="scene-svg" viewBox="0 0 540 280" aria-hidden="true">
        <defs>
          <filter id="campfireHeroGlow" x="-20%" y="-20%" width="140%" height="140%">
            <feDropShadow dx="0" dy="0" stdDeviation="4.5" flood-color="#8ec5fc" flood-opacity="0.6" />
            <feDropShadow dx="0" dy="0" stdDeviation="1" flood-color="#ffffff" flood-opacity="0.8" />
          </filter>

          <filter id="fireAuraGlow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="0" stdDeviation="12" flood-color="#ff6a00" flood-opacity="0.65" />
          </filter>

          <linearGradient id="cfBootGrad" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#804a29" />
            <stop offset="60%" stop-color="#542e17" />
            <stop offset="100%" stop-color="#33190b" />
          </linearGradient>

          <linearGradient id="cfStaffGrad" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#7a4e2f" />
            <stop offset="100%" stop-color="#361f10" />
          </linearGradient>

          <linearGradient id="cfHatGrad" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#4d6175" />
            <stop offset="60%" stop-color="#38495a" />
            <stop offset="100%" stop-color="#24303d" />
          </linearGradient>

          <linearGradient id="cfCauldronGrad" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#45454b" />
            <stop offset="50%" stop-color="#222328" />
            <stop offset="100%" stop-color="#121214" />
          </linearGradient>
        </defs>

        <!-- Forest Ground Mound -->
        <path d="M 0 250 Q 270 230 540 250 L 540 280 L 0 280 Z" fill="#0c0e12" />

        <!-- Wooden Log that Hero sits on -->
        <rect x="70" y="210" width="130" height="32" rx="8" fill="#321e13" stroke="#1c1009" stroke-width="2.5" />
        <ellipse cx="70" cy="226" rx="6" ry="16" fill="#42291b" />

        <!-- HERO (Sitting on Log beside Campfire, adhering to visual standard) -->
        <g filter="url(#campfireHeroGlow)">
          <!-- Visible Brown Leather Boots -->
          <path d="M 130 225 L 142 258 L 134 266 L 158 266 L 152 254 L 140 225 Z" fill="url(#cfBootGrad)" stroke="#221208" stroke-width="1.8" />
          <path d="M 110 225 L 120 258 L 112 266 L 136 266 L 130 254 L 120 225 Z" fill="url(#cfBootGrad)" stroke="#221208" stroke-width="1.8" />
          <rect x="142" y="252" width="7" height="4" rx="1" fill="#c0c8cf" stroke="#333" stroke-width="0.7" />
          <rect x="120" y="252" width="7" height="4" rx="1" fill="#c0c8cf" stroke="#333" stroke-width="0.7" />

          <!-- Cloak & Body -->
          <path d="M 132 145 
                   C 115 145, 102 165, 105 190 
                   C 107 205, 110 220, 105 235 
                   L 155 235 
                   C 160 220, 158 205, 160 190 
                   C 165 165, 150 145, 132 145 Z" 
                fill="#2c3742" 
                stroke="#526678" 
                stroke-width="1.8" />

          <!-- Arm firmly holding Wooden Staff -->
          <path d="M 85 105 L 80 266" stroke="url(#cfStaffGrad)" stroke-width="6" stroke-linecap="round" />
          <circle cx="85" cy="100" r="5" fill="#69d2e7" opacity="0.9" />
          <!-- Hand gripping staff -->
          <circle cx="83" cy="180" r="7" fill="#6e4225" stroke="#2c170a" stroke-width="1.2" />

          <!-- Chest with Glowing Silver Talisman/Medallion -->
          <circle cx="134" cy="168" r="6" fill="#eaf2f8" stroke="#7b92a5" stroke-width="1.5" />
          <circle cx="134" cy="168" r="3.5" fill="#9db4c7" />
          <circle cx="134" cy="168" r="1.5" fill="#ffffff" />

          <!-- Head, Face Shadow & Glowing Blue Eyes -->
          <ellipse cx="132" cy="128" rx="15" ry="13" fill="#11161a" />
          <ellipse cx="127" cy="127" rx="2.5" ry="1.5" fill="#5ce1e6" />
          <ellipse cx="139" cy="127" rx="2.5" ry="1.5" fill="#5ce1e6" />

          <!-- Grey-Blue Pointed Wizard/Wanderer Hat -->
          <ellipse cx="132" cy="118" rx="34" ry="9" fill="url(#cfHatGrad)" stroke="#688094" stroke-width="1.6" />
          <path d="M 112 110 C 112 95, 116 75, 126 62 C 136 55, 142 60, 140 68 C 138 76, 144 95, 150 110 Z" fill="url(#cfHatGrad)" stroke="#688094" stroke-width="1.4" />
        </g>

        <!-- CAMPFIRE (Center-Right) -->
        <g class="campfire-core" filter="url(#fireAuraGlow)">
          <!-- Stones Ring -->
          <ellipse cx="350" cy="242" rx="46" ry="15" fill="#1b140f" stroke="#0e0a07" stroke-width="2.5" />
          <!-- Crossed Logs -->
          <path d="M 325 244 L 375 232 M 330 232 L 370 244" stroke="#4a2c19" stroke-width="8" stroke-linecap="round" />
          <!-- Glowing Ember Bed -->
          <ellipse cx="350" cy="235" rx="30" ry="10" fill="#ff4500" opacity="0.85" />
          <!-- Animated Flame Tongues -->
          <g class="fire-flames">
            <path d="M 334 238 Q 346 175 350 155 Q 358 190 366 238 Z" fill="#ff4d00" class="flame-outer" />
            <path d="M 340 236 Q 348 185 351 170 Q 355 200 362 236 Z" fill="#ffa200" class="flame-mid" />
            <path d="M 345 234 Q 350 200 351 185 Q 354 205 357 234 Z" fill="#fff5b3" class="flame-inner" />
          </g>
        </g>

        <!-- IRON TRIPOD & CAULDRON (Казан на треноге) -->
        <g class="cauldron-assembly">
          <!-- Iron Tripod Legs -->
          <path d="M 350 110 L 310 248" stroke="#1d1e22" stroke-width="5" stroke-linecap="round" />
          <path d="M 350 110 L 390 248" stroke="#1d1e22" stroke-width="5" stroke-linecap="round" />
          <path d="M 350 110 L 350 244" stroke="#111215" stroke-width="4.5" stroke-linecap="round" />

          <!-- Top Iron Ring / Shackle -->
          <circle cx="350" cy="110" r="6" fill="none" stroke="#5a5d68" stroke-width="3" />

          <!-- Iron Hanging Chain -->
          <line x1="350" y1="116" x2="350" y2="160" stroke="#727682" stroke-width="3.5" stroke-dasharray="4,2" />

          <!-- Cast-Iron Cauldron / Pot -->
          <g class="cauldron-pot">
            <!-- Belly -->
            <path d="M 324 180 
                     C 320 215, 335 228, 350 228 
                     C 365 228, 380 215, 376 180 Z" 
                  fill="url(#cfCauldronGrad)" 
                  stroke="#5a5d68" 
                  stroke-width="2.5" />
            <!-- Rim -->
            <ellipse cx="350" cy="180" rx="27" ry="6" fill="#1a1a1d" stroke="#686b76" stroke-width="2.2" />
            <!-- Bubbling Stew inside -->
            <ellipse cx="350" cy="180" rx="23" ry="4.5" fill="#8b4513" />
            <circle cx="344" cy="179" r="2" fill="#d2691e" />
            <circle cx="355" cy="181" r="1.8" fill="#ff7f50" />
            <!-- Handle -->
            <path d="M 324 180 Q 350 156 376 180" fill="none" stroke="#727682" stroke-width="3" />
          </g>
        </g>
      </svg>
    </div>

    <!-- Candle Health Bar with Large Contrast Numbers -->
    <div class="candles-status-row">
      <span class="candles-label">ОГОНЬ ТВОЕЙ ДУШИ:</span>
      <div class="candles-display" class:igniting={healedCandleAnim}>
        {#each Array(maxCandles) as _, i}
          <span 
            class="respite-candle" 
            class:lit={i < playerCandles}
            class:just-healed={healedCandleAnim && i === playerCandles - 1}
          >
            {i < playerCandles ? '🕯️' : '🌑'}
          </span>
        {/each}
        <span class="candles-fraction"><strong>{playerCandles}</strong> / {maxCandles}</span>
      </div>
    </div>

    <!-- Narrative monologue from the wanderer -->
    <div class="hero-monologue-bubble">
      <span class="quill-icon">📜</span>
      <p class="monologue-text">{activeMonologue}</p>
    </div>

    <!-- Interactive Respite Panel with 15s Progress Bar -->
    <div class="interactive-respite-panel">
      {#if !campfireHealed && !isCooking}
        <p class="prompt-text">Тепло огня и аромат похлебки манят отдохнуть. Выбери действие:</p>
        <div class="respite-options-grid">
          <!-- Option 1: Cook Stew (15-sec progress bar) -->
          <button 
            type="button" 
            class="respite-choice-btn cook-btn"
            class:disabled={isCandlesFull}
            disabled={isCandlesFull}
            on:click={startCookingStew}
          >
            <span class="choice-icon">🍲</span>
            <div class="choice-content">
              <span class="choice-title">
                {isCandlesFull ? 'Свечи полны (Лечение не требуется)' : 'Сварить похлебку / Набраться сил'}
              </span>
              <span class="choice-desc">
                {#if isCandlesFull}
                  Пламя всех твоих свечей горит ярко — отдых у костра не требуется.
                {:else if isNightmare}
                  15 секунд отдыха у костра защищают твои свечи от угасания
                {:else}
                  15 секунд отдыха у костра восстанавливают 1 свечу
                {/if}
              </span>
            </div>
          </button>

          <!-- Option 2: Continue into Darkness without wait or bonus -->
          <button 
            type="button" 
            class="respite-choice-btn skip-forward-btn"
            on:click={handleContinue}
          >
            <span class="choice-icon">🚪</span>
            <div class="choice-content">
              <span class="choice-title">Продолжить путь в темноту</span>
              <span class="choice-desc">
                {isNightmare ? 'Мгновенный уход в темноту гасит 1 свечу (смертельно при 1 свече!)' : 'Мгновенный переход вперед без ожидания и без бонуса'}
              </span>
            </div>
          </button>
        </div>

      {:else if isCooking}
        <!-- 15-second Progress Bar without time numbers -->
        <div class="cooking-progress-card">
          <div class="cooking-status-line">
            <span class="cooking-pot-icon">🍲</span>
            <span class="cooking-text">{cookStatusText}</span>
          </div>

          <div class="progress-track">
            <div class="progress-fill" style="width: {cookProgress}%;"></div>
          </div>

          <p class="cooking-subtext">Подожди, пока кипящий отвар вернет тепло твоим замерзшим пальцам...</p>
        </div>

      {:else}
        <!-- Healed Confirmation -->
        <div class="healed-confirmation-card">
          <span class="healed-star">✨</span>
          <span class="healed-message">
            {isNightmare ? 'Похлебка съедена, пламя защищено от ледяного дыхания кошмара! Твоя свеча спасена.' : 'Похлебка съедена, зажглась 1 новая свеча! Силы восстановлены.'}
          </span>
          <span class="healed-star">✨</span>
        </div>

        <div class="campfire-footer">
          <button class="continue-journey-btn" on:click={handleContinue}>
            <span>Продолжить путь в темноту ➔</span>
          </button>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .campfire-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(3, 4, 6, 0.95);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(8px);
    padding: 16px;
    box-sizing: border-box;
    animation: fadeIn 0.4s ease-out;
  }

  .campfire-container {
    width: 100%;
    max-width: 780px;
    background: radial-gradient(circle at 50% 15%, #241a14 0%, #140e0b 60%, #090605 100%);
    border: 2px solid #8f5c35;
    border-radius: 12px;
    padding: 24px 30px;
    box-shadow: 0 0 50px rgba(0, 0, 0, 0.95), 0 0 30px rgba(255, 110, 20, 0.2);
    display: flex;
    flex-direction: column;
    align-items: center;
    box-sizing: border-box;
    max-height: 94vh;
    overflow-y: auto;
  }

  .campfire-header {
    display: flex;
    align-items: center;
    gap: 14px;
    margin-bottom: 8px;
  }

  .campfire-title {
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 2.2rem;
    color: #ffd899;
    letter-spacing: 2px;
    text-shadow: 0 0 16px rgba(255, 140, 40, 0.7);
    margin: 0;
  }

  .moon-icon {
    font-size: 1.8rem;
    filter: drop-shadow(0 0 8px rgba(255, 180, 50, 0.5));
  }

  /* Large Central Scene */
  .campfire-scene {
    position: relative;
    width: 100%;
    max-width: 560px;
    height: 250px;
    margin: 6px 0 14px 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .scene-svg {
    width: 100%;
    height: 100%;
  }

  /* Swirling steam puffs above cauldron */
  .steam-emitter {
    position: absolute;
    top: 70px;
    right: 175px;
    pointer-events: none;
    z-index: 5;
  }

  .steam-cloud {
    position: absolute;
    width: 22px;
    height: 22px;
    background: radial-gradient(circle, rgba(230, 240, 255, 0.35) 0%, rgba(200, 220, 240, 0) 70%);
    border-radius: 50%;
    filter: blur(3px);
    animation: steamRise 3s ease-out infinite;
  }

  .steam-cloud.s1 { animation-delay: 0s; left: 0; }
  .steam-cloud.s2 { animation-delay: 1s; left: 8px; }
  .steam-cloud.s3 { animation-delay: 2s; left: -6px; }

  @keyframes steamRise {
    0% { transform: translateY(0) scale(0.6); opacity: 0; }
    30% { opacity: 0.6; }
    100% { transform: translateY(-55px) scale(1.8); opacity: 0; }
  }

  /* Flame Tongues Animation */
  .flame-outer {
    animation: flameTongue 1.4s ease-in-out infinite alternate;
    transform-origin: 350px 238px;
  }
  .flame-mid {
    animation: flameTongue 1.1s ease-in-out infinite alternate-reverse;
    transform-origin: 350px 238px;
  }
  .flame-inner {
    animation: flameTongue 0.8s ease-in-out infinite alternate;
    transform-origin: 350px 238px;
  }

  @keyframes flameTongue {
    0% { transform: scaleY(0.92) scaleX(0.96) rotate(-1deg); }
    100% { transform: scaleY(1.1) scaleX(1.04) rotate(1.5deg); }
  }

  /* Embers particles */
  .embers-emitter {
    position: absolute;
    bottom: 40px;
    right: 180px;
    pointer-events: none;
  }

  .ember-spark {
    position: absolute;
    width: 4px;
    height: 4px;
    background: #ffaa00;
    border-radius: 50%;
    box-shadow: 0 0 6px #ff4500;
    animation: emberFloat 2.6s linear infinite;
    animation-delay: var(--rnd);
  }

  @keyframes emberFloat {
    0% { transform: translateY(0) translateX(0) scale(1); opacity: 1; }
    100% { transform: translateY(-90px) translateX(calc(sin(var(--i)) * 30px)) scale(0.2); opacity: 0; }
  }

  /* Candles row with large contrast font */
  .candles-status-row {
    display: flex;
    align-items: center;
    gap: 14px;
    margin-bottom: 12px;
    background: rgba(18, 14, 11, 0.7);
    padding: 8px 18px;
    border-radius: 20px;
    border: 1px solid #5a3c26;
  }

  .candles-label {
    font-size: 1.15rem;
    font-weight: 700;
    color: #e5bd87;
    letter-spacing: 1.5px;
    text-transform: uppercase;
  }

  .candles-display {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .respite-candle {
    font-size: 1.5rem;
    filter: drop-shadow(0 0 6px #ff8c00);
  }

  .respite-candle:not(.lit) {
    filter: grayscale(1);
    opacity: 0.45;
  }

  .candles-fraction {
    font-size: 1.45rem; /* Large contrasting number */
    font-weight: 700;
    color: #ffd700;
    margin-left: 8px;
    text-shadow: 0 0 8px rgba(255, 215, 0, 0.6);
  }

  .hero-monologue-bubble {
    width: 100%;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    background: rgba(10, 8, 7, 0.75);
    border: 1px solid #483424;
    border-radius: 8px;
    padding: 12px 18px;
    margin-bottom: 18px;
  }

  .quill-icon {
    font-size: 1.3rem;
  }

  .monologue-text {
    font-size: 1.05rem;
    font-style: italic;
    color: #edd9b9;
    line-height: 1.5;
    margin: 0;
  }

  /* Respite panel */
  .interactive-respite-panel {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .prompt-text {
    font-size: 1.1rem;
    color: #d6be9f;
    margin-bottom: 14px;
    text-align: center;
  }

  .respite-options-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    width: 100%;
  }

  .respite-choice-btn {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 16px 20px;
    background: linear-gradient(180deg, #2b1c12 0%, #170e08 100%);
    border: 2px solid #734827;
    border-radius: 8px;
    cursor: pointer;
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.8), inset 0 1px 1px rgba(255, 220, 160, 0.2);
    transition: all 0.25s ease;
    text-align: left;
  }

  .respite-choice-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    border-color: #d49344;
    box-shadow: 0 6px 20px rgba(212, 147, 68, 0.35);
  }

  .respite-choice-btn:disabled,
  .respite-choice-btn.disabled {
    opacity: 0.45;
    cursor: not-allowed;
    border-color: #4a3828;
    filter: grayscale(0.5);
    box-shadow: none;
    transform: none;
  }

  .cook-btn {
    border-color: #a86c35;
  }

  .choice-icon {
    font-size: 2.2rem;
  }

  .choice-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .choice-title {
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 1.15rem;
    font-weight: 700;
    color: #fff2d1;
    letter-spacing: 0.5px;
  }

  .choice-desc {
    font-size: 0.88rem;
    color: #b0957b;
    line-height: 1.35;
  }

  /* 15-second Cooking Progress Bar without time digits */
  .cooking-progress-card {
    width: 100%;
    background: rgba(18, 14, 11, 0.9);
    border: 2px solid #d49344;
    border-radius: 10px;
    padding: 20px 24px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    box-shadow: 0 0 24px rgba(212, 147, 68, 0.35);
    animation: glowPulse 2s ease-in-out infinite;
  }

  @keyframes glowPulse {
    0%, 100% { border-color: #d49344; box-shadow: 0 0 20px rgba(212, 147, 68, 0.3); }
    50% { border-color: #ffb852; box-shadow: 0 0 30px rgba(255, 184, 82, 0.5); }
  }

  .cooking-status-line {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .cooking-pot-icon {
    font-size: 1.8rem;
    animation: potShake 0.4s ease-in-out infinite alternate;
  }

  @keyframes potShake {
    0% { transform: rotate(-3deg); }
    100% { transform: rotate(3deg); }
  }

  .cooking-text {
    font-size: 1.25rem;
    font-weight: 700;
    color: #ffe699;
    letter-spacing: 1px;
    text-shadow: 0 0 8px rgba(255, 180, 50, 0.6);
  }

  .progress-track {
    width: 100%;
    height: 16px;
    background: #140d09;
    border: 2px solid #5d3d25;
    border-radius: 10px;
    overflow: hidden;
    position: relative;
    box-shadow: inset 0 2px 6px rgba(0, 0, 0, 0.9);
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #ff4500 0%, #ff8c00 50%, #ffd700 100%);
    box-shadow: 0 0 14px rgba(255, 140, 0, 0.85);
    transition: width 0.1s linear;
  }

  .cooking-subtext {
    font-size: 0.92rem;
    color: #c9b398;
    margin: 0;
    font-style: italic;
  }

  /* Healed confirmation */
  .healed-confirmation-card {
    display: flex;
    align-items: center;
    gap: 12px;
    background: rgba(30, 24, 16, 0.9);
    border: 2px solid #e5a742;
    border-radius: 8px;
    padding: 14px 22px;
    color: #fff1d0;
    font-size: 1.15rem;
    font-weight: 700;
    text-shadow: 0 0 8px rgba(229, 167, 66, 0.6);
    margin-bottom: 16px;
  }

  .healed-star {
    font-size: 1.5rem;
  }

  .campfire-footer {
    width: 100%;
    display: flex;
    justify-content: center;
  }

  .continue-journey-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 14px 36px;
    background: linear-gradient(180deg, #3d2315 0%, #1e1109 100%);
    border: 2px solid #9e6435;
    border-radius: 8px;
    color: #ffe6ba;
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 1.15rem;
    font-weight: 700;
    letter-spacing: 1.5px;
    text-transform: uppercase;
    cursor: pointer;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.85);
    transition: all 0.25s ease;
  }

  .continue-journey-btn:hover {
    transform: scale(1.03);
    border-color: #e5a742;
    box-shadow: 0 6px 24px rgba(229, 167, 66, 0.4);
    color: #ffffff;
  }

  @media (max-width: 650px) {
    .campfire-container {
      padding: 18px;
    }
    .respite-options-grid {
      grid-template-columns: 1fr;
    }
    .campfire-title {
      font-size: 1.7rem;
    }
    .campfire-scene {
      height: 200px;
    }
  }
</style>
