<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  const dispatch = createEventDispatcher();

  let step = 0; // 0: Doors burst & Smoke, 1: Alena on Throne speaks, 2: Player choice, 3: Alena reveals truth, 4: Witch Transformation
  let playerChoice = '';
  let isTransforming = false;
  let smokeActive = true;

  onMount(() => {
    // Fade out previous combat music smoothly
    audioEngine.fadeTrackOut(2.0);

    // Sequence opening audio effects
    setTimeout(() => {
      audioEngine.playHeavyDoorsKick();
      audioEngine.playGlassShatter();
    }, 400);

    setTimeout(() => {
      audioEngine.playStoneFootsteps();
      smokeActive = false;
      step = 1;
    }, 1800);
  });

  function handlePlayerChoice(choice) {
    playerChoice = choice;
    audioEngine.playWickIgnition(); // burning skin sizzle
    step = 3;
  }

  function handleTransform() {
    isTransforming = true;
    audioEngine.playWitchTransformation();
    step = 4;
  }

  function handleStartBattle() {
    dispatch('startBattle');
  }
</script>

<div class="twist-viewport" class:crimson-burn={step >= 1}>
  <!-- Black Smoke Particle Overlay -->
  <div class="smoke-curtain" class:smoke-clearing={!smokeActive}></div>

  <!-- Red-hot burning vignette from medallion on hero's chest -->
  <div class="medallion-heat-vignette"></div>

  <!-- The Throne Room of Roots & Bones -->
  <div class="throne-chamber-stage">
    <!-- Throne of Roots and Bones -->
    <div class="root-throne">
      <svg viewBox="0 0 400 450" class="throne-svg" aria-hidden="true">
        <!-- Arch of Black Basalt and Roots -->
        <path d="M 50 450 C 40 250, 80 80, 200 40 C 320 80, 360 250, 350 450 Z" fill="#0c070f" stroke="#251228" stroke-width="4" />
        <path d="M 90 450 Q 120 180 200 120 Q 280 180 310 450 Z" fill="#08040a" />

        <!-- Entangled Bone Roots -->
        <path d="M 60 420 Q 110 320 100 200 Q 140 160 200 150 Q 260 160 300 200 Q 290 320 340 420" 
              stroke="#3d2836" stroke-width="14" stroke-linecap="round" fill="none" />
        <path d="M 110 430 Q 150 280 180 240 Q 220 280 290 430" 
              stroke="#593b4f" stroke-width="8" stroke-linecap="round" fill="none" />

        <!-- Bone Skulls embedded in throne back -->
        <circle cx="160" cy="180" r="10" fill="#2d1a29" />
        <circle cx="200" cy="160" r="12" fill="#3b2236" />
        <circle cx="240" cy="180" r="10" fill="#2d1a29" />

        <!-- Throne Seat Platform -->
        <rect x="130" y="320" width="140" height="40" rx="8" fill="#180d1a" stroke="#482545" stroke-width="3" />
      </svg>
    </div>

    <!-- Alena Figure: changes between Maiden (Step 1-3) and Morovaya Witch (Step 4) -->
    <div class="figure-container" class:witch-form={step === 4}>
      {#if step < 4}
        <!-- Alena in Torn Maiden Dress -->
        <div class="maiden-figure">
          <svg viewBox="0 0 200 320" class="figure-svg" aria-hidden="true">
            <!-- Long dark auburn hair -->
            <path d="M 80 50 Q 60 130 55 240 Q 100 250 145 240 Q 140 130 120 50 Z" fill="#2c140d" />
            <!-- Torn grey peasant dress -->
            <path d="M 85 90 L 70 280 L 130 280 L 115 90 Z" fill="#3b3238" stroke="#1d171c" stroke-width="2" />
            <!-- Head & delicate face -->
            <circle cx="100" cy="55" r="22" fill="#d9baa7" />
            <!-- Dark sorrowful eyes -->
            <ellipse cx="93" cy="54" rx="3" ry="2" fill="#20110e" />
            <ellipse cx="107" cy="54" rx="3" ry="2" fill="#20110e" />
            <!-- Faded flower wreath -->
            <path d="M 80 40 Q 100 32 120 40" stroke="#7a4234" stroke-width="4" stroke-linecap="round" fill="none" />
            <circle cx="90" cy="38" r="3" fill="#a8523b" />
            <circle cx="100" cy="35" r="4" fill="#c47854" />
            <circle cx="110" cy="38" r="3" fill="#a8523b" />
          </svg>
        </div>
      {:else}
        <!-- Morovaya Witch Transformed -->
        <div class="witch-figure">
          <!-- Wings of Raven Feathers & Grave Dust -->
          <div class="raven-wings left-wing"></div>
          <div class="raven-wings right-wing"></div>

          <svg viewBox="0 0 240 360" class="figure-svg transformed-svg" aria-hidden="true">
            <!-- Shadow aura -->
            <ellipse cx="120" cy="340" rx="90" ry="16" fill="rgba(80, 0, 100, 0.4)" />

            <!-- Billowing Black Hair into Void -->
            <path d="M 60 10 Q 120 -20 180 10 Q 210 110 190 260 Q 120 280 50 260 Q 30 110 60 10 Z" fill="#080309" />

            <!-- Obsidian Robe with Pulsing Veins -->
            <path d="M 95 100 L 40 340 L 200 340 L 145 100 Z" fill="#130617" stroke="#481854" stroke-width="3" />
            <path d="M 120 100 Q 100 220 70 340 M 120 100 Q 140 220 170 340" stroke="#962985" stroke-width="2" />

            <!-- Witch Face with glowing violet eyes & obsidian horns -->
            <circle cx="120" cy="70" r="26" fill="#1f1424" stroke="#5f246b" stroke-width="2" />
            <ellipse cx="111" cy="68" rx="4" ry="5" fill="#f07aff" class="glowing-witch-eye" />
            <ellipse cx="129" cy="68" rx="4" ry="5" fill="#f07aff" class="glowing-witch-eye" />

            <!-- Sharp Obsidian Claws -->
            <path d="M 65 140 Q 40 180 30 220" stroke="#120616" stroke-width="8" stroke-linecap="round" fill="none" />
            <path d="M 175 140 Q 200 180 210 220" stroke="#120616" stroke-width="8" stroke-linecap="round" fill="none" />
          </svg>
        </div>
      {/if}
    </div>
  </div>

  <!-- Narrative Dialogue Box -->
  <div class="dialogue-card">
    <div class="speaker-tag">
      <span class="speaker-sigil">{step === 4 ? '👑' : '🥀'}</span>
      <span class="speaker-name">{step === 4 ? 'МОРОВАЯ ВЕДЬМА' : 'АЛЁНА'}</span>
    </div>

    <div class="dialogue-body">
      <!-- Medallion Mutation Lore Callout -->
      <div class="medallion-mutation-lore">
        <div class="medallion-badge-glow">
          <svg viewBox="0 0 44 44" class="molten-badge-svg" aria-hidden="true">
            <circle cx="22" cy="22" r="19" fill="#140608" stroke="#ff2a00" stroke-width="2.5" />
            <circle cx="22" cy="22" r="14" fill="#24070a" stroke="#ff5500" stroke-width="1.2" stroke-dasharray="2 2" />
            <path d="M 22 9 L 22 35 M 13 17 L 31 27 M 13 27 L 31 17" stroke="#ff3b20" stroke-width="2.5" stroke-linecap="round" />
            <circle cx="22" cy="22" r="4" fill="#ff7700" />
          </svg>
        </div>
        <div class="medallion-lore-content">
          <div class="medallion-lore-title">🔥 МУТАЦИЯ: РАСКАЛЁННЫЙ МЕДАЛЬОН</div>
          <div class="medallion-lore-desc">
            Серебро обжигает грудь и чернеет. В центре загорается тёмно-багровая руна. Артефакт пробуждается для дуэли на костях!
          </div>
        </div>
      </div>

      {#if step === 1}
        <p class="dialogue-speech">
          «Ты так упорно шел по моим следам, милый... Неужели ты думал, что смертная девка смогла бы выжить в этих топях?»
        </p>
        <div class="dialogue-actions">
          <button class="choice-btn" on:click={() => handlePlayerChoice('what_happened')}>
            «Алёна?.. Что они с тобой сделали?!»
          </button>
          <button class="choice-btn" on:click={() => handlePlayerChoice('you_lied')}>
            «Ты лгала мне с самого начала...»
          </button>
        </div>
      {:else if step === 3}
        <p class="dialogue-speech reveal-speech">
          «Они ничего со мной не сделали. Я создала эту корчму сотни лет назад! Мне нужны были души тех трех стражей, а ты принес мне их силу на блюдечке! Карты были забавой для глупцов. Настоящая цена платится костями!»
        </p>
        <div class="dialogue-actions">
          <button class="choice-btn transform-trigger-btn" on:click={handleTransform}>
            Сжать кулаки и встретить правду ➔
          </button>
        </div>
      {:else if step === 4}
        <div class="witch-reveal-banner">
          <p class="witch-speech">
            «Посмотрим, на что годится твоя бренная душа, когда в дело вступят резные кости!»
          </p>
          <button class="start-part2-btn" on:click={handleStartBattle}>
            <span class="bone-icon">🎲</span>
            <span class="btn-main-title">НАЧАТЬ ОБРЯД НА КОСТЯХ</span>
            <span class="btn-sub">дуэль пяти резных костей D6</span>
          </button>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .twist-viewport {
    position: fixed;
    inset: 0;
    background: radial-gradient(circle at 50% 60%, #15091a 0%, #07030a 70%, #020103 100%);
    z-index: 10000;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    overflow: hidden;
    color: #ebd7be;
  }

  .smoke-curtain {
    position: absolute;
    inset: 0;
    background: #000000;
    z-index: 50;
    pointer-events: none;
    opacity: 1;
    transition: opacity 1.8s ease-out;
  }

  .smoke-curtain.smoke-clearing {
    opacity: 0;
  }

  /* Burning red-hot medallion vignette */
  .medallion-heat-vignette {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 5;
    box-shadow: inset 0 0 80px rgba(255, 30, 0, 0.45);
    animation: pulseBurn 2.5s infinite alternate ease-in-out;
  }

  @keyframes pulseBurn {
    0% { box-shadow: inset 0 0 60px rgba(220, 20, 0, 0.3); }
    100% { box-shadow: inset 0 0 110px rgba(255, 50, 0, 0.65); }
  }

  .throne-chamber-stage {
    position: relative;
    width: 100%;
    max-width: 600px;
    height: 380px;
    display: flex;
    justify-content: center;
    align-items: flex-end;
    margin-top: 10px;
    z-index: 10;
  }

  .root-throne {
    position: absolute;
    bottom: 0;
    width: 380px;
    height: 380px;
    display: flex;
    justify-content: center;
    align-items: flex-end;
  }

  .throne-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 6px 25px rgba(0, 0, 0, 0.95));
  }

  .figure-container {
    position: relative;
    z-index: 15;
    margin-bottom: 20px;
  }

  .maiden-figure {
    width: 140px;
    height: 220px;
    animation: gentleBreathe 3s infinite ease-in-out;
  }

  @keyframes gentleBreathe {
    0%, 100% { transform: scale(1); }
    50% { transform: scale(1.02) translateY(-2px); }
  }

  .witch-figure {
    width: 220px;
    height: 300px;
    position: relative;
    animation: witchHover 2.5s infinite ease-in-out;
  }

  @keyframes witchHover {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-8px); }
  }

  .raven-wings {
    position: absolute;
    top: 30px;
    width: 110px;
    height: 160px;
    background: radial-gradient(circle at 50% 50%, #2b0833 0%, #0c020f 80%);
    border-radius: 80% 0 80% 0;
    border: 2px solid #5a146e;
    filter: drop-shadow(0 0 15px #6e178a);
    opacity: 0.85;
  }

  .left-wing {
    left: -40px;
    transform: rotate(-25deg);
    animation: flapLeft 2.2s infinite ease-in-out;
  }

  .right-wing {
    right: -40px;
    transform: rotate(25deg) scaleX(-1);
    animation: flapRight 2.2s infinite ease-in-out;
  }

  @keyframes flapLeft {
    0%, 100% { transform: rotate(-25deg); }
    50% { transform: rotate(-35deg) scaleX(1.1); }
  }

  @keyframes flapRight {
    0%, 100% { transform: rotate(25deg) scaleX(-1); }
    50% { transform: rotate(35deg) scaleX(-1.1); }
  }

  .glowing-witch-eye {
    animation: eyeFlare 1.4s infinite alternate;
  }

  @keyframes eyeFlare {
    0% { filter: drop-shadow(0 0 3px #ff3df8); }
    100% { filter: drop-shadow(0 0 12px #ffffff) drop-shadow(0 0 18px #ff45ff); }
  }

  /* Dialogue Box */
  .dialogue-card {
    position: relative;
    z-index: 20;
    max-width: 680px;
    width: 100%;
    background: rgba(14, 7, 18, 0.9);
    border: 2px solid #632669;
    border-radius: 12px;
    padding: 20px 24px;
    box-shadow: 0 12px 35px rgba(0, 0, 0, 0.95), 0 0 25px rgba(120, 20, 140, 0.2);
    backdrop-filter: blur(8px);
  }

  .speaker-tag {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: rgba(45, 12, 50, 0.6);
    border: 1px solid #7c2d85;
    border-radius: 14px;
    padding: 3px 12px;
    margin-bottom: 12px;
  }

  .speaker-name {
    font-family: var(--font-title, Georgia, serif);
    font-size: 0.95rem;
    font-weight: bold;
    color: #f7a8e8;
    letter-spacing: 1px;
  }

  .medallion-mutation-lore {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 10px 14px;
    margin-bottom: 14px;
    background: linear-gradient(90deg, rgba(40, 8, 10, 0.9) 0%, rgba(20, 4, 6, 0.95) 100%);
    border: 1px solid #ff4422;
    border-radius: 8px;
    box-shadow: 0 0 16px rgba(255, 50, 0, 0.35);
    animation: moltenPulse 2s infinite alternate ease-in-out;
  }

  @keyframes moltenPulse {
    0% { box-shadow: 0 0 10px rgba(255, 40, 0, 0.3); border-color: #ff3311; }
    100% { box-shadow: 0 0 20px rgba(255, 70, 0, 0.7); border-color: #ff7744; }
  }

  .medallion-badge-glow {
    flex-shrink: 0;
    width: 44px;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .molten-badge-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 0 8px #ff4400);
  }

  .medallion-lore-content {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .medallion-lore-title {
    font-family: var(--font-title, Georgia, serif);
    font-size: 0.85rem;
    font-weight: bold;
    color: #ff8855;
    letter-spacing: 0.05em;
    text-shadow: 0 0 8px rgba(255, 60, 0, 0.6);
  }

  .medallion-lore-desc {
    font-family: var(--font-body, Georgia, serif);
    font-size: 0.85rem;
    line-height: 1.35;
    color: #ffd0c0;
  }

  .dialogue-speech {
    font-family: var(--font-body, Georgia, serif);
    font-size: 1.15rem;
    line-height: 1.6;
    color: #ebd7be;
    margin: 0 0 16px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.9);
  }

  .dialogue-speech.reveal-speech {
    color: #fce1f8;
  }

  .dialogue-actions {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .choice-btn {
    padding: 10px 18px;
    background: linear-gradient(180deg, #301336 0%, #17081a 100%);
    border: 1px solid #783485;
    border-radius: 8px;
    color: #ebd7be;
    font-family: var(--font-title, Georgia, serif);
    font-size: 1.05rem;
    cursor: pointer;
    text-align: left;
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.7);
    transition: all 0.2s ease;
  }

  .choice-btn:hover {
    transform: translateX(4px);
    border-color: #df62f7;
    box-shadow: 0 4px 15px rgba(220, 100, 245, 0.35);
    color: #ffffff;
  }

  .choice-btn.transform-trigger-btn {
    text-align: center;
    border-color: #a338b5;
    background: linear-gradient(180deg, #4d1857 0%, #200824 100%);
    font-weight: bold;
  }

  .witch-reveal-banner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    text-align: center;
  }

  .witch-speech {
    font-family: var(--font-title, Georgia, serif);
    font-size: 1.25rem;
    color: #ff9cf7;
    margin: 0;
    text-shadow: 0 0 10px rgba(255, 100, 240, 0.6);
  }

  .start-part2-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 14px 32px;
    background: linear-gradient(180deg, #5c1866 0%, #28082e 100%);
    border: 2px solid #ca47e0;
    border-radius: 10px;
    color: #ffffff;
    cursor: pointer;
    box-shadow: 0 6px 25px rgba(200, 50, 220, 0.45), inset 0 1px 2px rgba(255, 255, 255, 0.4);
    transition: all 0.2s ease;
    animation: btnThrob 2s infinite alternate;
  }

  .start-part2-btn:hover {
    transform: translateY(-2px) scale(1.04);
    box-shadow: 0 8px 30px rgba(220, 80, 245, 0.65);
  }

  @keyframes btnThrob {
    0% { border-color: #ca47e0; }
    100% { border-color: #ff88ff; box-shadow: 0 6px 30px rgba(255, 100, 250, 0.6); }
  }

  .bone-icon {
    font-size: 1.8rem;
  }

  .btn-main-title {
    font-family: var(--font-title, Georgia, serif);
    font-size: 1.35rem;
    font-weight: 900;
    letter-spacing: 1.5px;
  }

  .btn-sub {
    font-size: 0.78rem;
    color: #dfb2e8;
    letter-spacing: 0.8px;
    text-transform: uppercase;
  }
</style>
