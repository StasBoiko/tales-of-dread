<script>
  import { onMount, onDestroy } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  export let bossId = 'the_lanterns'; // 'the_lanterns', 'kingdom_of', 'let_it_go'

  let soundInterval = null;
  let isInteracting = false;
  let speechBubbleText = '';
  let speechTimeout = null;

  const bossReactions = {
    the_lanterns: 'Фонарщик: «Не трогай птицу, она считает твои последние вздохи!»',
    kingdom_of: 'Узник: «Кот чует твою дрожь... и помнит твои карты.»',
    let_it_go: 'Архимаг: «Шипение змея предрекает твой прах!»'
  };

  onMount(() => {
    // Ambient sound trigger every 14 seconds
    soundInterval = setInterval(() => {
      triggerAnimalSound(false);
    }, 14000);
  });

  onDestroy(() => {
    if (soundInterval) clearInterval(soundInterval);
    if (speechTimeout) clearTimeout(speechTimeout);
  });

  function triggerAnimalSound(isClick = true) {
    isInteracting = true;
    setTimeout(() => {
      isInteracting = false;
    }, 1100);

    if (isClick) {
      speechBubbleText = bossReactions[bossId] || '';
      clearTimeout(speechTimeout);
      speechTimeout = setTimeout(() => {
        speechBubbleText = '';
      }, 3200);
    }

    if (bossId === 'the_lanterns') {
      audioEngine.playRavenCaw();
    } else if (bossId === 'kingdom_of') {
      audioEngine.playCatMeow();
    } else {
      audioEngine.playViperHiss();
    }
  }
</script>

<!-- Boss Companion Animal Widget -->
<button 
  type="button"
  class="companion-widget" 
  class:interacting={isInteracting}
  on:click={() => triggerAnimalSound(true)}
  title="Фамильяр босса (кликни для взаимодействия)"
  aria-label="Фамильяр босса"
>
  {#if speechBubbleText}
    <div class="companion-speech-bubble" role="status">
      <span class="speech-quote">💬</span>
      <span class="speech-text">{speechBubbleText}</span>
    </div>
  {/if}
  {#if bossId === 'the_lanterns'}
    <!-- ============================================== -->
    <!-- BOSS 1 COMPANION: Raven (Ворон Фонарщика)      -->
    <!-- ============================================== -->
    <div class="animal-container raven-container">
      <svg class="animal-svg" viewBox="0 0 120 120">
        <defs>
          <radialGradient id="ravenFeather" cx="40%" cy="35%" r="65%">
            <stop offset="0%" stop-color="#2a3342" />
            <stop offset="50%" stop-color="#141a22" />
            <stop offset="100%" stop-color="#070a0e" />
          </radialGradient>
          <filter id="ravenAura" x="-20%" y="-20%" width="140%" height="140%">
            <feDropShadow dx="0" dy="2" stdDeviation="4" flood-color="#000" flood-opacity="0.8" />
            <feDropShadow dx="0" dy="0" stdDeviation="2" flood-color="#385470" flood-opacity="0.4" />
          </filter>
        </defs>

        <!-- Wooden Perch Stick with Nails -->
        <path d="M 15 95 L 105 92" stroke="#4a3321" stroke-width="7" stroke-linecap="round" />
        <ellipse cx="60" cy="98" rx="42" ry="4" fill="rgba(0,0,0,0.6)" />

        <g class="raven-body-group" filter="url(#ravenAura)">
          <!-- Tail Feathers -->
          <path d="M 35 88 L 18 108 L 30 92 L 24 112 L 38 92 Z" fill="#0c1015" stroke="#1c2430" stroke-width="1.2" />

          <!-- Claws gripping perch -->
          <path d="M 46 90 L 44 96 M 52 90 L 52 97 M 58 90 L 60 96" stroke="#221710" stroke-width="2.5" stroke-linecap="round" />

          <!-- Wing with Flapping Animation -->
          <g class="raven-wing">
            <path d="M 45 48 C 30 52, 22 70, 32 88 C 42 78, 56 70, 52 52 Z" fill="url(#ravenFeather)" stroke="#1a222c" stroke-width="1.5" />
            <path d="M 38 60 Q 25 75 34 85" stroke="#334254" stroke-width="1.2" fill="none" />
            <path d="M 44 56 Q 32 68 40 80" stroke="#334254" stroke-width="1.2" fill="none" />
          </g>

          <!-- Main Body -->
          <path d="M 48 42 C 40 52, 42 75, 52 86 C 62 88, 70 78, 68 56 C 66 45, 58 38, 48 42 Z" fill="url(#ravenFeather)" stroke="#0e1318" stroke-width="1.5" />

          <!-- Head & Throat with Breathing Bobbing -->
          <g class="raven-head">
            <circle cx="62" cy="38" r="14" fill="url(#ravenFeather)" stroke="#141a22" stroke-width="1.2" />
            <!-- Keen Beak -->
            <path d="M 72 35 L 94 42 L 72 45 Z" fill="#1b1d22" stroke="#0a0c0e" stroke-width="1.2" />
            <!-- Glowing Amber Eye -->
            <circle cx="66" cy="35" r="3.2" fill="#ffb700" />
            <circle cx="66" cy="35" r="1.4" fill="#000000" />
            <circle cx="65.2" cy="34" r="0.8" fill="#ffffff" />
          </g>
        </g>
      </svg>
    </div>

  {:else if bossId === 'kingdom_of'}
    <!-- ============================================== -->
    <!-- BOSS 2 COMPANION: Black Cat (Кошка Узника)     -->
    <!-- ============================================== -->
    <div class="animal-container cat-container">
      <svg class="animal-svg" viewBox="0 0 120 120">
        <defs>
          <radialGradient id="catFur" cx="45%" cy="35%" r="65%">
            <stop offset="0%" stop-color="#2a2327" />
            <stop offset="60%" stop-color="#141113" />
            <stop offset="100%" stop-color="#080708" />
          </radialGradient>
          <filter id="catAura" x="-20%" y="-20%" width="140%" height="140%">
            <feDropShadow dx="0" dy="2" stdDeviation="4" flood-color="#000" flood-opacity="0.85" />
            <feDropShadow dx="0" dy="0" stdDeviation="2" flood-color="#2a774d" flood-opacity="0.3" />
          </filter>
        </defs>

        <!-- Ground Shadow -->
        <ellipse cx="60" cy="100" rx="38" ry="7" fill="rgba(0,0,0,0.65)" />

        <g class="cat-body-group" filter="url(#catAura)">
          <!-- Swishing Tail -->
          <g class="cat-tail">
            <path d="M 32 82 C 15 80, 8 60, 16 46 C 20 40, 26 44, 22 52 C 16 62, 22 75, 34 78 Z" fill="#141113" stroke="#221b1f" stroke-width="1.2" />
          </g>

          <!-- Back Paws & Hips -->
          <ellipse cx="42" cy="84" rx="14" ry="12" fill="url(#catFur)" stroke="#221b1f" stroke-width="1.2" />
          <ellipse cx="38" cy="94" rx="8" ry="4" fill="#0d0b0c" />

          <!-- Torso -->
          <path d="M 42 74 C 40 60, 52 52, 64 56 C 72 60, 75 75, 72 88 C 65 94, 48 94, 42 74 Z" fill="url(#catFur)" stroke="#221b1f" stroke-width="1.2" />

          <!-- Front Paws -->
          <rect x="62" y="76" width="7" height="20" rx="3.5" fill="#141113" stroke="#241d21" stroke-width="1" />
          <rect x="71" y="76" width="7" height="20" rx="3.5" fill="#141113" stroke="#241d21" stroke-width="1" />

          <!-- Head with Twitching Ears -->
          <g class="cat-head">
            <circle cx="66" cy="46" r="15" fill="url(#catFur)" stroke="#241d21" stroke-width="1.2" />
            <!-- Left Ear -->
            <path d="M 54 38 L 50 20 L 62 33 Z" fill="#1c161a" stroke="#2b2027" stroke-width="1.2" />
            <path d="M 54 35 L 52 25 L 59 33 Z" fill="#613b4e" opacity="0.7" />
            <!-- Right Ear -->
            <path d="M 70 33 L 80 20 L 78 38 Z" fill="#1c161a" stroke="#2b2027" stroke-width="1.2" />
            <path d="M 72 33 L 78 25 L 76 35 Z" fill="#613b4e" opacity="0.7" />

            <!-- Glowing Emerald Eyes -->
            <ellipse cx="61" cy="45" rx="3.5" ry="4.5" fill="#00ff88" />
            <ellipse cx="71" cy="45" rx="3.5" ry="4.5" fill="#00ff88" />
            <ellipse cx="61" cy="45" rx="1.2" ry="4" fill="#000" />
            <ellipse cx="71" cy="45" rx="1.2" ry="4" fill="#000" />
            <circle cx="60.2" cy="43.5" r="0.8" fill="#fff" />
            <circle cx="70.2" cy="43.5" r="0.8" fill="#fff" />

            <!-- Little Nose & Whiskers -->
            <polygon points="66,51 64.5,49 67.5,49" fill="#c4789a" />
            <path d="M 58 52 L 40 50 M 58 54 L 38 56 M 74 52 L 92 50 M 74 54 L 94 56" stroke="#8e8074" stroke-width="0.8" opacity="0.7" />
          </g>
        </g>
      </svg>
    </div>

  {:else}
    <!-- ============================================== -->
    <!-- BOSS 3 COMPANION: Fire Viper & Bat             -->
    <!-- ============================================== -->
    <div class="animal-container viper-container">
      <svg class="animal-svg" viewBox="0 0 120 120">
        <defs>
          <linearGradient id="viperScales" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stop-color="#5a1812" />
            <stop offset="50%" stop-color="#2c0c08" />
            <stop offset="100%" stop-color="#120403" />
          </linearGradient>
          <radialGradient id="viperEye" cx="40%" cy="35%" r="60%">
            <stop offset="0%" stop-color="#ff3300" />
            <stop offset="80%" stop-color="#990000" />
          </radialGradient>
          <filter id="viperAura" x="-20%" y="-20%" width="140%" height="140%">
            <feDropShadow dx="0" dy="2" stdDeviation="5" flood-color="#ff2600" flood-opacity="0.5" />
          </filter>
        </defs>

        <!-- Burning Embers beneath coiling serpent -->
        <ellipse cx="60" cy="98" rx="42" ry="8" fill="rgba(255, 60, 0, 0.2)" />
        <ellipse cx="60" cy="98" rx="34" ry="5" fill="rgba(0, 0, 0, 0.7)" />

        <g class="viper-group" filter="url(#viperAura)">
          <!-- Bat Wing Flapping on Serpent's Back -->
          <g class="bat-wings">
            <path d="M 45 45 Q 20 25 10 38 Q 24 50 20 65 Q 35 60 45 60 Z" fill="#200a12" stroke="#521526" stroke-width="1.2" />
            <path d="M 75 45 Q 100 25 110 38 Q 96 50 100 65 Q 85 60 75 60 Z" fill="#200a12" stroke="#521526" stroke-width="1.2" />
          </g>

          <!-- Coiling Body -->
          <path d="M 28 92 C 20 75, 45 70, 60 76 C 75 82, 95 75, 90 90 C 85 102, 35 104, 28 92 Z" fill="url(#viperScales)" stroke="#6b1d16" stroke-width="1.8" />
          <!-- Rising Neck & Hood -->
          <path d="M 48 78 C 42 62, 54 48, 62 42 C 68 38, 76 45, 70 60 C 66 70, 58 75, 48 78 Z" fill="url(#viperScales)" stroke="#6b1d16" stroke-width="1.8" />

          <!-- Serpent Head -->
          <g class="viper-head">
            <polygon points="62,38 78,34 82,46 66,48" fill="#420f0a" stroke="#872318" stroke-width="1.2" />
            <!-- Glowing Ruby Eye -->
            <circle cx="73" cy="40" r="2.8" fill="url(#viperEye)" />
            <ellipse cx="73" cy="40" rx="0.8" ry="2.2" fill="#000" />
            <circle cx="72.2" cy="39" r="0.6" fill="#fff" />
            <!-- Flicking Forked Tongue -->
            <g class="forked-tongue">
              <path d="M 82 43 L 95 44 L 102 40 M 95 44 L 101 47" stroke="#ff2a00" stroke-width="1.2" stroke-linecap="round" fill="none" />
            </g>
          </g>
        </g>
      </svg>
    </div>
  {/if}
</button>

<style>
  .companion-widget {
    background: none;
    border: none;
    padding: 0;
    margin: 0;
    cursor: pointer;
    position: relative;
    outline: none;
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    transition: transform 0.25s cubic-bezier(0.18, 0.89, 0.32, 1.28);
  }

  .companion-widget:hover {
    transform: scale(1.08) translateY(-4px);
  }

  .companion-widget:active, .companion-widget.interacting {
    transform: scale(1.15) translateY(-6px);
  }

  .animal-container {
    width: 95px;
    height: 95px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    position: relative;
  }

  .animal-svg {
    width: 100%;
    height: 100%;
    overflow: visible;
  }



  /* ================== RAVEN ANIMATIONS ================== */
  .raven-container .raven-body-group {
    transform-origin: 60px 92px;
    animation: ravenPerchIdle 4.2s ease-in-out infinite;
  }

  .raven-container .raven-wing {
    transform-origin: 50px 52px;
    animation: wingShuffle 5.5s ease-in-out infinite;
  }

  .raven-container .raven-head {
    transform-origin: 62px 42px;
    animation: ravenHeadLook 3.8s ease-in-out infinite;
  }

  .interacting .raven-container .raven-wing {
    animation: ravenFlapQuick 0.4s ease-in-out 3;
  }

  @keyframes ravenPerchIdle {
    0%, 100% { transform: translateY(0) rotate(0deg); }
    50% { transform: translateY(-3px) rotate(1.5deg); }
  }

  @keyframes wingShuffle {
    0%, 90%, 100% { transform: rotate(0deg); }
    93% { transform: rotate(-8deg) scaleY(1.05); }
    96% { transform: rotate(4deg); }
  }

  @keyframes ravenHeadLook {
    0%, 70%, 100% { transform: rotate(0deg); }
    75% { transform: rotate(-12deg); }
    85% { transform: rotate(10deg); }
  }

  @keyframes ravenFlapQuick {
    0%, 100% { transform: rotate(0deg) scale(1); }
    50% { transform: rotate(-25deg) scale(1.15); }
  }

  /* ================== CAT ANIMATIONS ================== */
  .cat-container .cat-tail {
    transform-origin: 32px 82px;
    animation: catTailSwish 3.5s ease-in-out infinite;
  }

  .cat-container .cat-head {
    transform-origin: 66px 52px;
    animation: catHeadTilt 4.8s ease-in-out infinite;
  }

  .interacting .cat-container .cat-body-group {
    animation: catArchBack 0.7s ease-in-out;
  }

  @keyframes catTailSwish {
    0%, 100% { transform: rotate(0deg); }
    50% { transform: rotate(-14deg) scaleY(1.04); }
  }

  @keyframes catHeadTilt {
    0%, 80%, 100% { transform: rotate(0deg); }
    85% { transform: rotate(-8deg); }
    92% { transform: rotate(6deg); }
  }

  @keyframes catArchBack {
    0%, 100% { transform: translateY(0) scaleY(1); }
    50% { transform: translateY(-6px) scaleY(1.12); }
  }

  /* ================== VIPER ANIMATIONS ================== */
  .viper-container .viper-group {
    transform-origin: 60px 85px;
    animation: viperSway 4s ease-in-out infinite;
  }

  .viper-container .bat-wings {
    transform-origin: 60px 48px;
    animation: batWingFlap 3s ease-in-out infinite;
  }

  .viper-container .forked-tongue {
    transform-origin: 82px 43px;
    animation: tongueFlick 2.5s ease-in-out infinite;
  }

  .interacting .viper-container .viper-head {
    animation: viperStrike 0.6s cubic-bezier(0.18, 0.89, 0.32, 1.28);
  }

  @keyframes viperSway {
    0%, 100% { transform: rotate(0deg) scale(1); }
    50% { transform: rotate(3deg) scale(1.03); }
  }

  @keyframes batWingFlap {
    0%, 100% { transform: scaleY(1); }
    50% { transform: scaleY(0.65) rotate(4deg); }
  }

  @keyframes tongueFlick {
    0%, 80%, 100% { opacity: 0; transform: scaleX(0.2); }
    85%, 95% { opacity: 1; transform: scaleX(1); }
  }

  @keyframes viperStrike {
    0%, 100% { transform: translate(0, 0); }
    40% { transform: translate(8px, -4px) scale(1.1); }
  }

  /* Companion Speech Bubble */
  .companion-speech-bubble {
    position: absolute;
    top: calc(100% + 8px);
    right: 50%;
    transform: translateX(50%);
    background: rgba(16, 10, 6, 0.95);
    border: 1px solid #d49344;
    border-radius: 6px;
    padding: 6px 12px;
    display: flex;
    align-items: center;
    gap: 6px;
    white-space: nowrap;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.9), 0 0 10px rgba(212, 147, 68, 0.35);
    z-index: 100;
    pointer-events: none;
    animation: speechPop 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  .speech-quote {
    font-size: 0.95rem;
  }

  .speech-text {
    font-family: var(--font-body, serif);
    font-size: 0.84rem;
    color: #ffeed1;
    font-style: italic;
    text-shadow: 0 1px 2px #000;
  }

  @keyframes speechPop {
    0% { opacity: 0; transform: translateX(50%) translateY(-6px) scale(0.9); }
    100% { opacity: 1; transform: translateX(50%) translateY(0) scale(1); }
  }
</style>
