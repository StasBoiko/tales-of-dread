<script>
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';
  import { setLang, t, langStore } from '../i18n/index.js';

  const dispatch = createEventDispatcher();

  // Reactive storyLines based on current language
  let currentLang = $langStore;
  let storyLines = $derived([
    t('prologue.intro1'),
    t('prologue.intro2'),
    t('prologue.intro3')
  ]);

  let currentLineIndex = 0;
  let displayedText = ['', '', ''];
  let isTyping = true;
  let typeInterval = null;
  let charIndex = 0;
  let isLightningActive = false;
  let lightningTimer = null;
  let prologueRestartRequested = false;

  onMount(() => {
    // Ambient rain during prologue
    audioEngine.startRainAmbience();
    startTypewriter();
    scheduleLightning();
  });

  onDestroy(() => {
    if (typeInterval) clearTimeout(typeInterval);
    if (lightningTimer) clearTimeout(lightningTimer);
    audioEngine.stopRainAmbience();
  });

  // Restart prologue from beginning
  function restartPrologue() {
    currentLineIndex = 0;
    displayedText = ['', '', ''];
    charIndex = 0;
    isTyping = true;
    typeNextChar();
  }

  function startTypewriter() {
    charIndex = 0;
    currentLineIndex = 0;
    typeNextChar();
  }

  function typeNextChar() {
    if (currentLineIndex >= storyLines.length) {
      isTyping = false;
      return;
    }

    const currentLine = storyLines[currentLineIndex];
    if (charIndex < currentLine.length) {
      displayedText[currentLineIndex] = currentLine.slice(0, charIndex + 1);
      displayedText = [...displayedText];
      charIndex++;

      const delay = currentLine[charIndex - 1] === '.' || currentLine[charIndex - 1] === '—' ? 220 : 36;
      typeInterval = setTimeout(typeNextChar, delay);
    } else {
      currentLineIndex++;
      charIndex = 0;
      typeInterval = setTimeout(typeNextChar, 500);
    }
  }

  function skipTypewriter() {
    if (typeInterval) clearTimeout(typeInterval);
    displayedText = [...storyLines];
    currentLineIndex = storyLines.length;
    isTyping = false;
    audioEngine.init().catch(() => {});
  }

  function scheduleLightning() {
    const nextTime = 4000 + Math.random() * 6000;
    lightningTimer = setTimeout(() => {
      triggerLightning();
      scheduleLightning();
    }, nextTime);
  }

  function triggerLightning() {
    isLightningActive = true;
    audioEngine.playThunderRumble();
    setTimeout(() => {
      isLightningActive = false;
    }, 280);
  }

  function handleEnterTavern() {
    audioEngine.stopRainAmbience();
    audioEngine.playIronBolt();
    // Synchronously initiate audio playback immediately inside user gesture
    audioEngine.startMainMusic();
    dispatch('enter');
  }

  function handleLangChange(newLang) {
    if (setLang(newLang)) {
      restartPrologue();
      dispatch('languageChanged', newLang);
    }
  }
</script>

<div class="prologue-viewport" class:lightning={isLightningActive}>
  <div class="rain-overlay"></div>
  <div class="fog-layer"></div>

  <!-- Animated Hero Approaching with Cinematic Step Animation -->
  <div class="hero-stage">
    <div class="hero-walk-container">
      <svg class="hero-svg" viewBox="0 0 240 340" aria-hidden="true">
        <defs>
          <!-- Soft glowing outer rim filter for crisp separation from dark background -->
          <filter id="heroRimGlow" x="-30%" y="-30%" width="160%" height="160%">
            <feDropShadow dx="0" dy="0" stdDeviation="5" flood-color="#8ec5fc" flood-opacity="0.55" />
            <feDropShadow dx="0" dy="0" stdDeviation="1.5" flood-color="#ffffff" flood-opacity="0.75" />
          </filter>

          <filter id="cyanEyeGlow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="0" stdDeviation="4" flood-color="#00e5ff" flood-opacity="0.9" />
          </filter>

          <linearGradient id="staffWood" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#7a4e2f" />
            <stop offset="50%" stop-color="#55331c" />
            <stop offset="100%" stop-color="#361f10" />
          </linearGradient>

          <linearGradient id="wizardHatGrad" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#4d6175" />
            <stop offset="50%" stop-color="#38495a" />
            <stop offset="100%" stop-color="#24303d" />
          </linearGradient>

          <linearGradient id="cloakGrad" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#3b4652" />
            <stop offset="40%" stop-color="#28323c" />
            <stop offset="100%" stop-color="#182028" />
          </linearGradient>

          <linearGradient id="bootLeather" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#804a29" />
            <stop offset="60%" stop-color="#542e17" />
            <stop offset="100%" stop-color="#33190b" />
          </linearGradient>
        </defs>

        <!-- Dynamic Ground Shadow -->
        <ellipse cx="120" cy="322" rx="76" ry="14" fill="rgba(0, 0, 0, 0.75)" class="shadow-pulse" />

        <!-- Stepping Legs & Boots -->
        <g class="hero-legs">
          <!-- Left Leg & Brown Leather Boot -->
          <g class="leg left-leg">
            <path d="M 98 260 L 92 305 L 82 312 L 80 322 L 108 322 L 105 305 L 104 260 Z" fill="url(#bootLeather)" stroke="#221208" stroke-width="1.8" />
            <!-- Leather Cuff & Buckle -->
            <path d="M 88 300 Q 98 296 108 300" stroke="#a06035" stroke-width="3" fill="none" />
            <rect x="94" y="297" width="8" height="5" rx="1" fill="#c0c8cf" stroke="#333" stroke-width="0.8" />
          </g>

          <!-- Right Leg & Brown Leather Boot -->
          <g class="leg right-leg">
            <path d="M 136 260 L 132 305 L 124 312 L 122 322 L 150 322 L 146 305 L 142 260 Z" fill="url(#bootLeather)" stroke="#221208" stroke-width="1.8" />
            <!-- Leather Cuff & Buckle -->
            <path d="M 130 300 Q 140 296 150 300" stroke="#a06035" stroke-width="3" fill="none" />
            <rect x="136" y="297" width="8" height="5" rx="1" fill="#c0c8cf" stroke="#333" stroke-width="0.8" />
          </g>
        </g>

        <!-- Entire Torso, Cloak, Head, and Staff with outer rim glow -->
        <g filter="url(#heroRimGlow)" class="hero-upper-body">
          <!-- Wooden Staff firmly planted -->
          <g class="hero-staff">
            <path d="M 52 80 L 46 324" stroke="url(#staffWood)" stroke-width="7" stroke-linecap="round" />
            <!-- Staff headpiece: carved bone and crystal bound with twine -->
            <path d="M 46 80 Q 42 60 52 48 Q 62 60 56 80 Z" fill="#8c6242" stroke="#2c1a0e" stroke-width="2" />
            <circle cx="51" cy="62" r="7" fill="#69d2e7" opacity="0.85" filter="url(#cyanEyeGlow)" />
            <path d="M 45 74 L 57 70 M 45 78 L 57 74" stroke="#d5b895" stroke-width="2" />
          </g>

          <!-- Flowing Cloak -->
          <path d="M 120 90
                   C 95 92, 70 120, 68 160
                   C 64 210, 56 265, 52 290
                   C 90 298, 150 298, 188 290
                   C 184 265, 176 210, 172 160
                   C 170 120, 145 92, 120 90 Z" 
                fill="url(#cloakGrad)" 
                stroke="#546575" 
                stroke-width="1.5" />

          <!-- Cloak Lapels and Interior Tunic -->
          <path d="M 105 105 L 120 180 L 135 105" stroke="#1d252e" stroke-width="3" fill="#182028" />

          <!-- Detailed Silver Talisman / Medallion on Chest -->
          <g class="chest-medallion">
            <!-- Silver Chain -->
            <path d="M 102 108 Q 120 135 138 108" stroke="#d8e3ec" stroke-width="2" stroke-dasharray="3,2" fill="none" />
            <!-- Medallion Center -->
            <circle cx="120" cy="136" r="9" fill="#eaf2f8" stroke="#7b92a5" stroke-width="2" />
            <circle cx="120" cy="136" r="6" fill="#adc2d1" />
            <!-- Ornate Runic Star -->
            <path d="M 120 130 L 120 142 M 114 136 L 126 136" stroke="#ffffff" stroke-width="1.8" />
            <circle cx="120" cy="136" r="2.5" fill="#ffffff" filter="url(#cyanEyeGlow)" />
          </g>

          <!-- Left Arm & Hand grasping staff firmly -->
          <g class="staff-arm">
            <!-- Sleeve -->
            <path d="M 96 112 Q 72 125 58 145" stroke="url(#cloakGrad)" stroke-width="16" stroke-linecap="round" fill="none" />
            <path d="M 96 112 Q 72 125 58 145" stroke="#546575" stroke-width="1.5" fill="none" />
            <!-- Detailed Leather Gauntlet / Hand wrapping fingers around staff -->
            <circle cx="50" cy="148" r="8" fill="#6e4225" stroke="#2c170a" stroke-width="1.5" />
            <!-- Four Knuckles / Fingers -->
            <rect x="44" y="142" width="11" height="4" rx="2" fill="#875632" stroke="#1f1007" stroke-width="1" />
            <rect x="44" y="147" width="12" height="4" rx="2" fill="#875632" stroke="#1f1007" stroke-width="1" />
            <rect x="44" y="152" width="11" height="4" rx="2" fill="#875632" stroke="#1f1007" stroke-width="1" />
          </g>

          <!-- Head, Face Shadow & Glowing Blue Eyes -->
          <g class="hero-head">
            <!-- Shaded Face under hat -->
            <ellipse cx="120" cy="85" rx="18" ry="16" fill="#11161a" />

            <!-- Glowing Blue Eyes peering from darkness -->
            <g filter="url(#cyanEyeGlow)">
              <ellipse cx="112" cy="84" rx="3.5" ry="2" fill="#5ce1e6" />
              <ellipse cx="128" cy="84" rx="3.5" ry="2" fill="#5ce1e6" />
              <!-- Bright Eye Cores -->
              <circle cx="112" cy="84" r="1.5" fill="#ffffff" />
              <circle cx="128" cy="84" r="1.5" fill="#ffffff" />
            </g>

            <!-- Grey-Blue Pointed Wizard/Wanderer Hat -->
            <!-- Wide Curved Brim -->
            <ellipse cx="120" cy="74" rx="42" ry="11" fill="url(#wizardHatGrad)" stroke="#688094" stroke-width="1.8" />
            <!-- Hat Band with silver buckle -->
            <path d="M 94 70 Q 120 74 146 70 L 144 63 Q 120 67 96 63 Z" fill="#1b242e" />
            <rect x="116" y="64" width="8" height="7" rx="1" fill="#b0c0cf" stroke="#333" stroke-width="0.8" />
            <!-- Pointed Arching Crown -->
            <path d="M 96 64 
                     C 96 48, 102 24, 114 10 
                     C 126 2, 134 8, 132 18 
                     C 130 28, 138 48, 144 64 Z" 
                  fill="url(#wizardHatGrad)" 
                  stroke="#688094" 
                  stroke-width="1.5" />
          </g>
        </g>
      </svg>
    </div>
  </div>

  <!-- Narrative Section with 1.5x Larger Candle Badge and Centered Action Button -->
  <div class="narrative-container">
    <!-- Language Switch (top-right of narrative) -->
    <div class="prologue-lang-switch">
      <button
        class="lang-btn"
        class:active={currentLang === 'ru'}
        onclick={() => handleLangChange('ru')}
        title="Русский"
      >RU</button>
      <button
        class="lang-btn"
        class:active={currentLang === 'en'}
        onclick={() => handleLangChange('en')}
        title="English"
      >EN</button>
    </div>

    <!-- Candle Badge: 1.5x Larger and High Contrast -->
    <div class="candle-count-badge">
      <span class="candle-flame-icon">🕯️</span>
      <span class="candle-count-text">ОСТАЛОСЬ СВЕЧЕЙ СУДЬБЫ: <strong>5 / 5</strong></span>
    </div>

    <!-- Typewriter Story Paragraphs -->
    <div class="text-block">
      {#each displayedText as line, idx}
        {#if line}
          <p class="story-paragraph" class:active-typing={idx === currentLineIndex && isTyping}>
            {line}
          </p>
        {/if}
      {/each}
    </div>

    <!-- Actions Area: Skip on side, Centered Big Enter Button -->
    <div class="actions-wrapper">
      {#if isTyping}
        <button class="skip-btn" onclick={skipTypewriter}>
          {t('prologue.skip')}
        </button>
      {/if}

      <!-- Clean Centered Button without icons above text -->
      <button 
        class="enter-threshold-btn" 
        onclick={handleEnterTavern}
        class:glow-pulse={!isTyping}
      >
        {t('prologue.enter')}
      </button>
    </div>
  </div>
</div>

<style>
  .prologue-viewport {
    position: fixed;
    inset: 0;
    background: radial-gradient(circle at 50% 75%, #181c20 0%, #0a0c0e 55%, #020304 100%);
    z-index: 9999;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    padding: 24px 20px 36px;
    overflow: hidden;
    color: #ded1be;
    transition: background-color 0.1s ease;
  }

  .prologue-viewport.lightning {
    background-color: #2b3846;
    filter: brightness(1.6);
  }

  .rain-overlay {
    position: absolute;
    inset: 0;
    pointer-events: none;
    background: repeating-linear-gradient(
      105deg,
      rgba(255, 255, 255, 0.06) 0px,
      rgba(255, 255, 255, 0.06) 1px,
      transparent 2px,
      transparent 18px
    );
    animation: rainFall 0.35s linear infinite;
    opacity: 0.85;
  }

  @keyframes rainFall {
    0% { transform: translateY(-50px) translateX(-10px); }
    100% { transform: translateY(50px) translateX(10px); }
  }

  .fog-layer {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 380px;
    background: linear-gradient(0deg, rgba(8, 10, 12, 0.96) 0%, rgba(8, 10, 12, 0) 100%);
    pointer-events: none;
  }

  /* Cinematic "Шаг навстречу" Hero Stage */
  .hero-stage {
    position: relative;
    width: 260px;
    height: 300px;
    margin-top: 10px;
    display: flex;
    justify-content: center;
    align-items: flex-end;
    z-index: 2;
  }

  .hero-walk-container {
    width: 100%;
    height: 100%;
    animation: heroWalkCycle 3.6s ease-in-out infinite alternate, heroStepZoom 16s ease-out infinite alternate;
    transform-origin: 50% 90%;
  }

  /* Shoulder sway and slight vertical bobbing during stride */
  @keyframes heroWalkCycle {
    0% {
      transform: translateY(0) rotate(-1.2deg);
    }
    50% {
      transform: translateY(-8px) rotate(1.2deg);
    }
    100% {
      transform: translateY(0) rotate(-1.2deg);
    }
  }

  /* Smooth scale approaching the screen ("Шаг навстречу") */
  @keyframes heroStepZoom {
    0% {
      transform: scale(0.92);
    }
    100% {
      transform: scale(1.16);
    }
  }

  .hero-svg {
    width: 100%;
    height: 100%;
  }

  /* Stepping legs subtle swing */
  .left-leg {
    animation: legStrideLeft 1.8s ease-in-out infinite alternate;
    transform-origin: 95px 260px;
  }

  .right-leg {
    animation: legStrideRight 1.8s ease-in-out infinite alternate;
    transform-origin: 135px 260px;
  }

  @keyframes legStrideLeft {
    0% { transform: translateY(0) scaleY(1); }
    100% { transform: translateY(-6px) scaleY(0.96); }
  }

  @keyframes legStrideRight {
    0% { transform: translateY(-6px) scaleY(0.96); }
    100% { transform: translateY(0) scaleY(1); }
  }

  .narrative-container {
    position: relative;
    z-index: 5;
    max-width: 760px;
    width: 100%;
    background: rgba(14, 18, 22, 0.88);
    border: 1px solid rgba(130, 160, 190, 0.35);
    border-radius: 12px;
    padding: 24px 32px;
    box-shadow: 0 14px 45px rgba(0, 0, 0, 0.95), inset 0 1px 1px rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(8px);
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  /* 1.5x Larger, High-Contrast Candle Badge */
  .candle-count-badge {
    display: inline-flex;
    align-items: center;
    gap: 12px;
    padding: 8px 22px;
    background: linear-gradient(180deg, #322212 0%, #1c1208 100%);
    border: 2px solid #e5a742;
    border-radius: 20px;
    font-size: 1.35rem; /* 1.5x larger */
    color: #ffe699; /* Bright contrast */
    text-shadow: 0 0 10px rgba(255, 180, 50, 0.65);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.8), inset 0 1px 1px rgba(255, 230, 160, 0.3);
    margin-bottom: 18px;
    letter-spacing: 1.5px;
  }

  .candle-flame-icon {
    font-size: 1.6rem;
    filter: drop-shadow(0 0 8px #ff9900);
  }

  .text-block {
    min-height: 130px;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .story-paragraph {
    font-family: var(--font-body, 'Cinzel', Georgia, serif);
    font-size: 1.15rem;
    line-height: 1.65;
    color: #edd9b9;
    margin: 0;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.95);
  }

  .story-paragraph.active-typing::after {
    content: '▋';
    display: inline-block;
    margin-left: 4px;
    color: #e5a742;
    animation: blink 0.75s infinite;
  }

  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0; }
  }

  .actions-wrapper {
    margin-top: 20px;
    width: 100%;
    max-width: 720px;
    display: flex;
    flex-direction: row;
    align-items: stretch;
    justify-content: center;
    gap: 16px;
  }

  .skip-btn {
    flex: 1;
    min-width: 220px;
    height: 50px;
    background: linear-gradient(180deg, #2b231c 0%, #17130f 100%);
    border: 2px solid #6b533d;
    color: #cfbda8;
    padding: 0 18px;
    border-radius: 8px;
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 1rem;
    font-weight: 700;
    letter-spacing: 1.5px;
    text-transform: uppercase;
    cursor: pointer;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.7);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
    transition: all 0.2s ease;
  }

  .skip-btn:hover {
    color: #fff1d0;
    border-color: #a87948;
    background: linear-gradient(180deg, #382d24 0%, #1f1a14 100%);
    transform: translateY(-1px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.85);
  }

  /* Matching sized Action Button */
  .enter-threshold-btn {
    flex: 1.2;
    min-width: 260px;
    height: 50px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 24px;
    background: linear-gradient(180deg, #3d2314 0%, #1c0f08 100%);
    border: 2px solid #c98838;
    border-radius: 8px;
    color: #fff1d0;
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 1.05rem;
    font-weight: 700;
    letter-spacing: 1.5px;
    text-transform: uppercase;
    cursor: pointer;
    box-shadow: 0 4px 18px rgba(0, 0, 0, 0.8), inset 0 1px 2px rgba(255, 230, 180, 0.3);
    box-sizing: border-box;
    transition: all 0.25s ease;
  }

  .enter-threshold-btn:hover {
    transform: translateY(-1px);
    border-color: #ffb852;
    box-shadow: 0 6px 24px rgba(255, 160, 50, 0.35), inset 0 1px 3px rgba(255, 255, 255, 0.4);
    color: #ffffff;
  }

  .enter-threshold-btn.glow-pulse {
    animation: enterBtnGlow 2.2s ease-in-out infinite;
  }

  @keyframes enterBtnGlow {
    0%, 100% {
      border-color: #c98838;
      box-shadow: 0 4px 18px rgba(0, 0, 0, 0.8);
    }
    50% {
      border-color: #ffbf66;
      box-shadow: 0 4px 26px rgba(255, 175, 60, 0.5);
    }
  }

  .prologue-lang-switch {
    position: absolute;
    top: 20px;
    right: 20px;
    z-index: 60;
    display: flex;
    gap: 6px;
    background: rgba(0, 0, 0, 0.5);
    padding: 4px;
    border-radius: 8px;
    border: 1px solid rgba(200, 136, 56, 0.3);
  }

  .lang-btn {
    min-width: 38px;
    height: 36px;
    background: linear-gradient(180deg, #2b231c 0%, #17130f 100%);
    border: 2px solid #6b533d;
    color: #cfbda8;
    padding: 0 8px;
    border-radius: 6px;
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 0.85rem;
    font-weight: 700;
    letter-spacing: 1px;
    text-transform: uppercase;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.7);
    transition: all 0.2s ease;
  }

  .lang-btn:hover {
    color: #fff1d0;
    border-color: #c98838;
    transform: translateY(-1px);
  }

  .lang-btn.active {
    color: #ffffff;
    border-color: #c98838;
    background: linear-gradient(180deg, #3d2314 0%, #1c0f08 100%);
    box-shadow: 0 4px 18px rgba(0, 0, 0, 0.8);
  }

  @media (max-width: 650px) {
    .prologue-lang-switch {
      position: static;
      width: auto;
      padding: 8px 0;
    }
    .lang-btn {
      min-width: auto;
      flex: 1;
      height: auto;
      padding: 6px 10px;
    }
  }

  @media (max-width: 650px) {
    .narrative-container {
      padding: 16px;
    }
    .candle-count-badge {
      font-size: 1.05rem;
      padding: 6px 14px;
    }
    .story-paragraph {
      font-size: 1rem;
    }
    .actions-wrapper {
      position: static;
    }
    .skip-btn {
      position: static;
      transform: none;
      margin-bottom: 8px;
    }
    .enter-threshold-btn {
      width: 100%;
      padding: 14px 20px;
      font-size: 1.1rem;
    }
  }
</style>
