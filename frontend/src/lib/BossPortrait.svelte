<script>
  import { onMount, onDestroy } from 'svelte';
  import { audioEngine } from './audioEngine.js';
  import { BOSS_QUOTES, getRandomQuote } from './bossQuotes.js';

  let { boss, playerScore = 0, isRoundOver = false, roundOutcome = '', isCrit = false, isHoveringHit = false, isIdleChoked = false } = $props();

  let prevHP = $state(boss ? boss.currentHp : 0);
  let isHit = $state(false);

  // Mouse tracking look-at & parallax state
  let targetLookX = 0;
  let targetLookY = 0;
  let smoothLookX = 0;
  let smoothLookY = 0;
  let animFrameId = null;

  // Speech bubble state
  let currentQuote = $state('');
  let showQuote = $state(false);
  let quoteTimer = null;
  let prevPlayerScore = $state(playerScore);

  function triggerQuote(quoteText, duration = 3800) {
    if (!quoteText) return;
    clearTimeout(quoteTimer);
    currentQuote = quoteText;
    showQuote = true;
    quoteTimer = setTimeout(() => {
      showQuote = false;
    }, duration);
  }

  // Watch HP change
  $effect(() => {
    if (boss && boss.currentHp < prevHP) {
      isHit = true;
      audioEngine.playBossStrike();
      if (boss.currentHp <= 0) {
        if (boss.id === 'the_lanterns') {
          audioEngine.playGlassShatter();
        } else {
          audioEngine.playBossStrike();
        }
        const deathQuotes = {
          the_lanterns: "Фонарь разбился... свет тонет в тине...",
          kingdom_of: "Стены рушатся... рукописи горят... свобода?!",
          let_it_go: "Пламя Бездны... ты поглотил его... или оно тебя?!"
        };
        triggerQuote(deathQuotes[boss.id] || "Это еще не конец...", 3500);
      } else {
        const bossLines = BOSS_QUOTES.bossSpecific[boss.id];
        if (bossLines && bossLines.hitReact) {
          triggerQuote(bossLines.hitReact, 2600);
        }
      }
      setTimeout(() => {
        isHit = false;
      }, 600);
      prevHP = boss.currentHp;
    } else if (boss) {
      prevHP = boss.currentHp;
    }
  });

  // Watch Player Score changes for trigger lines
  $effect(() => {
    if (playerScore !== prevPlayerScore) {
      if (playerScore === 21) {
        triggerQuote(getRandomQuote('crit21'), 4200);
      } else if (playerScore > 21) {
        triggerQuote(getRandomQuote('bust'), 4200);
      } else if (playerScore >= 18 && playerScore <= 20 && playerScore > prevPlayerScore) {
        triggerQuote(getRandomQuote('dangerHit'), 3600);
      }
      prevPlayerScore = playerScore;
    }
  });

  // Watch Idle trigger
  $effect(() => {
    if (isIdleChoked) {
      triggerQuote(getRandomQuote('idle'), 4000);
    }
  });

  function handleMouseMove(e) {
    const w = window.innerWidth || 1200;
    const h = window.innerHeight || 800;
    targetLookX = (e.clientX / w) - 0.5; // -0.5 to +0.5
    targetLookY = (e.clientY / h) - 0.5; // -0.5 to +0.5
  }

  function updateLookAt() {
    // Smooth lerp
    smoothLookX += (targetLookX - smoothLookX) * 0.08;
    smoothLookY += (targetLookY - smoothLookY) * 0.08;
    animFrameId = requestAnimationFrame(updateLookAt);
  }

  onMount(() => {
    window.addEventListener('mousemove', handleMouseMove);
    animFrameId = requestAnimationFrame(updateLookAt);

    // Initial taunt if present
    if (boss && BOSS_QUOTES.bossSpecific[boss.id]) {
      setTimeout(() => {
        triggerQuote(BOSS_QUOTES.bossSpecific[boss.id].intro, 3600);
      }, 1000);
    }
  });

  onDestroy(() => {
    window.removeEventListener('mousemove', handleMouseMove);
    if (animFrameId) cancelAnimationFrame(animFrameId);
    clearTimeout(quoteTimer);
  });
</script>

<div class="boss-section" class:boss-dead={boss && boss.currentHp <= 0}>
  <!-- Boss Health Runes / Bone Beads -->
  <div class="boss-health-bar">
    <div class="boss-name-banner">
      <h2 class="boss-name">{boss.name}</h2>
    </div>

    <div class="bone-beads-container" title="Жизни босса: {boss.currentHp} из {boss.maxHp}">
      {#each Array(boss.maxHp) as _, i}
        {@const intact = i < boss.currentHp}
        <div class="rune-bead {intact ? 'intact' : 'splintered'}" title={intact ? 'Живое сердце босса' : 'Разбитое сердце'}>
          <svg class="heart-life-icon" viewBox="0 0 24 24" width="20" height="20" aria-hidden="true">
            {#if intact}
              <defs>
                <linearGradient id="bossHeartGrad_{i}" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#ff4d4d" />
                  <stop offset="50%" stop-color="#cc1111" />
                  <stop offset="100%" stop-color="#7a0000" />
                </linearGradient>
                <filter id="bossHeartGlow_{i}" x="-20%" y="-20%" width="140%" height="140%">
                  <feDropShadow dx="0" dy="1" stdDeviation="2" flood-color="#ff1a1a" flood-opacity="0.85" />
                </filter>
              </defs>
              <path 
                d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" 
                fill="url(#bossHeartGrad_{i})" 
                stroke="#ffe6e6" 
                stroke-width="0.8"
                filter="url(#bossHeartGlow_{i})"
              />
              <ellipse cx="7" cy="7" rx="2.5" ry="1.2" fill="#ffffff" opacity="0.5" transform="rotate(-30 7 7)" />
            {:else}
              <path 
                d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" 
                fill="#241916" 
                stroke="#4a3028" 
                stroke-width="0.8"
                opacity="0.6"
              />
              <path d="M12 4 L10 11 L14 14 L11 20" stroke="#000000" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.8" />
            {/if}
          </svg>
        </div>
      {/each}
    </div>
  </div>

  <!-- Speech Bubble / Bloody Whisper Callout -->
  <div class="speech-bubble-anchor">
    {#if showQuote}
      <div class="boss-speech-bubble" class:boss-taunt-crit={playerScore === 21} class:boss-taunt-bust={playerScore > 21}>
        <div class="bubble-stem"></div>
        <div class="bubble-text">{currentQuote}</div>
      </div>
    {/if}
  </div>

  <!-- Animated Interactive Boss Stage with Look-At Parallax -->
  <div 
    class="boss-stage"
    class:hover-hit={isHoveringHit}
  >
    {#if boss.id === 'the_lanterns'}
      <!-- ============================================== -->
      <!-- BOSS 1: Скорбный Фонарщик (The Lantern Bearer) -->
      <!-- ============================================== -->
      <div 
        class="boss-figure tracker-container"
        class:hit-reaction={isHit}
        class:high-danger={playerScore > 17}
        style="
          --look-x: {smoothLookX * 24}px;
          --look-y: {smoothLookY * 16}px;
          --look-rot: {smoothLookX * 12}deg;
          --pupil-x: {smoothLookX * 9}px;
          --pupil-y: {smoothLookY * 7}px;
        "
      >
        <!-- Mist aura when danger is high -->
        <div class="cold-mist-aura" class:active={playerScore > 17}></div>

        <!-- Layered Ragged Cloak & Torso with Breathing Idle -->
        <div class="lantern-torso">
          <svg class="cloak-svg" viewBox="0 0 200 160">
            <defs>
              <linearGradient id="hoodGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#1c1611" />
                <stop offset="50%" stop-color="#100b08" />
                <stop offset="100%" stop-color="#050302" />
              </linearGradient>
              <linearGradient id="cloakRagged" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#221912" />
                <stop offset="70%" stop-color="#120c08" />
                <stop offset="100%" stop-color="#060403" />
              </linearGradient>
            </defs>
            <!-- Wide Ragged Shoulders -->
            <path d="M 20,160 L 35,90 L 70,60 L 130,60 L 165,90 L 180,160 C 150,150 135,160 100,155 C 65,160 50,150 20,160 Z" fill="url(#cloakRagged)" stroke="#3a281a" stroke-width="2" />
            <!-- Torn Cloak Strands -->
            <path d="M 25,140 Q 30,165 40,155 T 60,160 T 90,155 T 120,160 T 150,155 T 175,145" stroke="#100a06" stroke-width="4" fill="none" />
          </svg>
        </div>

        <!-- Head / Hood with Look-at Tracking -->
        <div class="lantern-head">
          <svg class="hood-svg" viewBox="0 0 120 120">
            <!-- Hood Rim and Folds -->
            <path d="M 20,110 C 15,60 30,15 60,12 C 90,15 105,60 100,110 C 85,115 70,110 60,112 C 50,110 35,115 20,110 Z" fill="url(#hoodGrad)" stroke="#3f2b1d" stroke-width="2" />
            <!-- Black Abyss Void inside hood -->
            <ellipse cx="60" cy="65" rx="32" ry="34" fill="#020202" />
          </svg>

          <!-- Glowing Amber Eyes tracking cursor -->
          <div class="amber-eyes-socket">
            <div class="amber-eye left">
              <div class="pupil-core"></div>
            </div>
            <div class="amber-eye right">
              <div class="pupil-core"></div>
            </div>
          </div>
        </div>

        <!-- Clawed Skeletal Hand (raises on Hit hover) -->
        <div class="clawed-hand-left" class:poised-to-strike={isHoveringHit}>
          <div class="claw finger-1"></div>
          <div class="claw finger-2"></div>
          <div class="claw finger-3"></div>
        </div>

        <!-- Stained Glass Forged Iron Lantern (Right arm with physics swing) -->
        <div class="lantern-swing-arm">
          <div class="chain-suspension">
            <div class="iron-link l1"></div>
            <div class="iron-link l2"></div>
            <div class="iron-link l3"></div>
          </div>

          <div class="stained-lantern-housing" class:blue-flame={playerScore > 17}>
            <div class="lantern-dome">
              <div class="dome-ring"></div>
            </div>
            <div class="glass-chamber">
              <!-- Slotted Iron Lattice Grate -->
              <div class="iron-grate vertical"></div>
              <div class="iron-grate horizontal"></div>
              <!-- Dancing Flame & Shimmer Rays -->
              <div class="flame-spark"></div>
              <div class="radiant-rays"></div>
            </div>
            <div class="lantern-base"></div>
          </div>
        </div>
      </div>

    {:else if boss.id === 'kingdom_of'}
      <!-- ============================================== -->
      <!-- BOSS 2: Узник Бумажных Стен (The Asylum Sovereign) -->
      <!-- ============================================== -->
      <div 
        class="boss-figure sovereign-container"
        class:hit-reaction={isHit}
        class:mad-laughter={roundOutcome === 'boss_win' || isRoundOver}
        style="
          --look-x: {smoothLookX * 22}px;
          --look-y: {smoothLookY * 16}px;
          --look-rot: {smoothLookX * 14}deg;
          --pupil-x: {smoothLookX * 10}px;
          --pupil-y: {smoothLookY * 8}px;
        "
      >
        <!-- Fluttering quills and manuscript scraps from back -->
        <div class="back-quills-aura">
          <div class="quill-scrap qs1">✒️</div>
          <div class="quill-scrap qs2">📜</div>
          <div class="quill-scrap qs3">🪶</div>
          <div class="quill-scrap qs4">📄</div>
        </div>

        <!-- Straitjacket Body with Leather Buckles -->
        <div class="straitjacket-torso">
          <svg class="jacket-svg" viewBox="0 0 190 150">
            <!-- Torn Canvas Jacket -->
            <path d="M 25,150 L 40,80 L 75,55 L 115,55 L 150,80 L 165,150 Z" fill="#b8ad97" stroke="#4a4235" stroke-width="2.5" />
            <!-- Stitched Patches & Ink Blots -->
            <path d="M 50,95 Q 60,110 55,130" stroke="#2c271f" stroke-width="2" stroke-dasharray="3,3" />
            <path d="M 135,90 Q 125,115 130,135" stroke="#2c271f" stroke-width="2" stroke-dasharray="3,3" />
            <circle cx="85" cy="110" r="14" fill="#15121b" opacity="0.6" />
          </svg>

          <!-- Heavy Leather Straps with Metallic Buckles -->
          <div class="leather-straps">
            <div class="strap s1"><span class="buckle"></span></div>
            <div class="strap s2"><span class="buckle"></span></div>
            <div class="strap s3"><span class="buckle"></span></div>
          </div>
        </div>

        <!-- Deranged Pale Head with Twitching Mouth & Sunken Eyes -->
        <div class="deranged-head" class:jaw-open={isHoveringHit}>
          <div class="pale-skull">
            <!-- Sunken Dark Eye Sockets -->
            <div class="sunken-sockets">
              <div class="maniac-eye left">
                <div class="maniac-pupil"></div>
              </div>
              <div class="maniac-eye right">
                <div class="maniac-pupil"></div>
              </div>
            </div>

            <!-- Manic Grin / Twitching Jaw -->
            <div class="twitching-mouth">
              <div class="teeth-row top"></div>
              <div class="mouth-void"></div>
              <div class="teeth-row bottom"></div>
            </div>
          </div>
        </div>

        <!-- Bound arms with twitching clawed fingers -->
        <div class="bound-hands-wrapper">
          <div class="twitching-paw left"></div>
          <div class="twitching-paw right"></div>
        </div>
      </div>

    {:else if boss.id === 'let_it_go'}
      <!-- ============================================== -->
      <!-- BOSS 3: Одержимый Архимаг (The Abyss Archmage)  -->
      <!-- ============================================== -->
      <div 
        class="boss-figure archmage-container"
        class:hit-reaction={isHit}
        class:crit-stun={isCrit}
        style="
          --look-x: {smoothLookX * 26}px;
          --look-y: {smoothLookY * 18}px;
          --look-rot: {smoothLookX * 10}deg;
          --pupil-x: {smoothLookX * 10}px;
          --pupil-y: {smoothLookY * 8}px;
        "
      >
        <!-- Swirling Ash & Inferno Aura -->
        <div class="archmage-inferno-glow"></div>

        <!-- Colossal Demonic Horns & Shoulders with Spikes -->
        <div class="demon-horns-spikes">
          <svg class="horns-svg" viewBox="0 0 240 120">
            <!-- Left Horn -->
            <path d="M 70,80 C 45,50 20,20 5,2 C 28,18 52,38 78,65 Z" fill="#1f0905" stroke="#7a2512" stroke-width="2" />
            <!-- Right Horn -->
            <path d="M 170,80 C 195,50 220,20 235,2 C 212,18 188,38 162,65 Z" fill="#1f0905" stroke="#7a2512" stroke-width="2" />
          </svg>
        </div>

        <!-- Fissured Magma Torso with Molten Veins -->
        <div class="colossal-torso">
          <svg class="torso-svg" viewBox="0 0 220 160">
            <path d="M 15,160 L 35,70 L 80,45 L 140,45 L 185,70 L 205,160 Z" fill="#140604" stroke="#52170b" stroke-width="3" />
            <!-- Glowing Magma Veins -->
            <path d="M 60,65 Q 85,90 80,120 T 70,150" stroke="#ff4500" stroke-width="3.5" fill="none" filter="drop-shadow(0 0 6px #ff2200)" />
            <path d="M 150,65 Q 130,95 135,120 T 145,155" stroke="#ff5500" stroke-width="3" fill="none" filter="drop-shadow(0 0 6px #ff2200)" />
            <path d="M 105,75 L 110,130" stroke="#ffaa00" stroke-width="2.5" fill="none" />
          </svg>
        </div>

        <!-- Demonic Head with Burning Crimson Eyes & Ash-Spewing Maw -->
        <div class="archmage-head" class:jaw-open={isHoveringHit}>
          <div class="obsidian-face">
            <div class="crimson-eyes-row">
              <div class="fiery-eye left">
                <div class="eye-ember"></div>
              </div>
              <div class="fiery-eye right">
                <div class="eye-ember"></div>
              </div>
            </div>

            <!-- Ash Spewing Jaw -->
            <div class="ash-maw">
              <div class="fanged-teeth"></div>
              <div class="ash-smoke-stream"></div>
            </div>
          </div>
        </div>

        <!-- Bone Staff with Smoking Skull & Blazing Ruby -->
        <div class="bone-staff-wrapper">
          <div class="staff-shaft"></div>
          <div class="staff-skull-head">
            <div class="skull-bone">💀</div>
            <div class="staff-flame-ruby"></div>
            <div class="smoke-wisp"></div>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .boss-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
    z-index: 10;
    width: 100%;
    max-width: 720px;
  }

  /* Boss Health Runes Banner */
  .boss-health-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 95%;
    max-width: 680px;
    background: radial-gradient(circle at 50% 30%, #20150e 0%, #100b07 80%, #080504 100%);
    border: 2px solid #543924;
    border-radius: 6px;
    padding: 6px 20px;
    box-shadow: 
      0 6px 18px rgba(0, 0, 0, 0.95),
      inset 0 1px 1px rgba(255, 220, 160, 0.15);
  }

  .boss-name-banner {
    display: flex;
    flex-direction: column;
  }

  .boss-name {
    font-family: var(--font-title);
    font-size: 1.6rem;
    color: #fae7cf;
    letter-spacing: 1.5px;
    margin: 0;
    text-shadow: 0 0 10px rgba(0, 0, 0, 0.9);
  }

  .bone-beads-container {
    display: flex;
    gap: 10px;
  }

  .rune-bead {
    width: 32px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #241412;
    border: 1px solid #5a2626;
    border-radius: 4px;
    box-shadow: inset 0 0 8px rgba(0, 0, 0, 0.9);
    transition: all 0.3s ease;
  }

  .rune-bead.intact {
    border-color: #ba3838;
    background: linear-gradient(180deg, #421616 0%, #210a0a 100%);
    box-shadow: 0 0 10px rgba(255, 40, 40, 0.4), inset 0 1px 2px rgba(255, 180, 180, 0.25);
  }

  .rune-bead.splintered {
    border-color: #211a14;
    background: #0f0b08;
    filter: grayscale(1);
    transform: scale(0.85) rotate(-6deg);
    opacity: 0.5;
  }

  .heart-life-icon {
    display: block;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.8));
    transition: transform 0.2s ease;
  }

  .rune-bead.intact:hover .heart-life-icon {
    transform: scale(1.15);
  }

  /* Speech Bubble Anchor & Floating Box */
  .speech-bubble-anchor {
    position: relative;
    width: 100%;
    height: 28px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 4px;
    z-index: 30;
  }

  .boss-speech-bubble {
    position: absolute;
    bottom: -6px;
    background: radial-gradient(circle at 50% 30%, #2b180e 0%, #170d07 100%);
    border: 2px solid #854f2c;
    border-radius: 6px;
    padding: 6px 18px;
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.95), 0 0 14px rgba(255, 120, 20, 0.3);
    animation: bubblePop 0.25s cubic-bezier(0.18, 0.89, 0.32, 1.28);
    max-width: 540px;
    text-align: center;
  }

  .bubble-stem {
    position: absolute;
    bottom: -7px;
    left: 50%;
    transform: translateX(-50%);
    width: 0;
    height: 0;
    border-left: 7px solid transparent;
    border-right: 7px solid transparent;
    border-top: 7px solid #854f2c;
  }

  .bubble-text {
    font-family: var(--font-title);
    font-size: 1.15rem;
    color: #ffe0b2;
    letter-spacing: 1px;
    text-shadow: 0 0 6px rgba(0, 0, 0, 0.95);
  }

  .boss-speech-bubble.boss-taunt-crit {
    border-color: #38a5ff;
    box-shadow: 0 0 25px rgba(56, 165, 255, 0.6);
  }

  .boss-speech-bubble.boss-taunt-crit .bubble-text {
    color: #aae5ff;
    text-shadow: 0 0 8px #00aaff;
  }

  .boss-speech-bubble.boss-taunt-bust {
    border-color: #e03232;
    box-shadow: 0 0 25px rgba(220, 30, 30, 0.6);
  }

  .boss-speech-bubble.boss-taunt-bust .bubble-text {
    color: #ffa4a4;
    text-shadow: 0 0 8px #ff2222;
  }

  @keyframes bubblePop {
    0% { transform: scale(0.7) translateY(10px); opacity: 0; }
    100% { transform: scale(1) translateY(0); opacity: 1; }
  }

  /* Stage & Boss Figure Containers */
  .boss-stage {
    position: relative;
    width: 320px;
    height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: -26px;
    z-index: 5;
  }

  .boss-figure {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: filter 0.3s ease;
  }

  /* ================== BOSS 1: TRACKER ================== */
  .tracker-container {
    animation: breathingIdle 4s infinite ease-in-out;
  }

  .cold-mist-aura {
    position: absolute;
    inset: -30px;
    background: radial-gradient(circle, rgba(120, 210, 255, 0.3) 0%, transparent 70%);
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.8s ease;
  }

  .cold-mist-aura.active {
    opacity: 1;
    animation: mistPulse 2s infinite alternate ease-in-out;
  }

  .lantern-torso {
    position: absolute;
    bottom: 5px;
    width: 180px;
    height: 140px;
  }

  .cloak-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 6px 14px rgba(0, 0, 0, 0.9));
  }

  .lantern-head {
    position: absolute;
    top: 15px;
    width: 100px;
    height: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
    transform: translate(var(--look-x), var(--look-y)) rotate(var(--look-rot));
    transition: transform 0.06s ease-out;
    z-index: 5;
  }

  .hood-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 4px 10px rgba(0, 0, 0, 0.95));
  }

  .amber-eyes-socket {
    position: absolute;
    top: 50px;
    width: 44px;
    display: flex;
    justify-content: space-between;
    pointer-events: none;
  }

  .amber-eye {
    width: 10px;
    height: 6px;
    background: radial-gradient(circle, #fff 0%, #ffaa00 70%, #aa5500 100%);
    border-radius: 50%;
    box-shadow: 0 0 10px #ff9900, 0 0 20px #ff6600;
    position: relative;
    overflow: hidden;
    animation: eyeBlink 4.5s infinite;
  }

  .pupil-core {
    position: absolute;
    width: 4px;
    height: 4px;
    background: #000;
    border-radius: 50%;
    top: 1px;
    left: 3px;
    transform: translate(var(--pupil-x), var(--pupil-y));
    transition: transform 0.05s ease-out;
  }

  /* Clawed Hand (striking posture on Hit hover) */
  .clawed-hand-left {
    position: absolute;
    bottom: 25px;
    left: 25px;
    display: flex;
    gap: 4px;
    transform: rotate(20deg);
    transition: transform 0.25s cubic-bezier(0.18, 0.89, 0.32, 1.28);
    z-index: 6;
  }

  .claw {
    width: 4px;
    height: 24px;
    background: linear-gradient(180deg, #70553e 0%, #291c13 80%, #0d0805 100%);
    border-radius: 2px 2px 50% 50%;
    box-shadow: 0 1px 3px rgba(0,0,0,0.9);
  }

  .claw.finger-2 { height: 28px; }

  .clawed-hand-left.poised-to-strike {
    transform: translateY(-16px) rotate(-15deg) scale(1.2);
  }

  /* Stained Glass Lantern with physics swing */
  .lantern-swing-arm {
    position: absolute;
    top: 35px;
    right: 15px;
    display: flex;
    flex-direction: column;
    align-items: center;
    transform-origin: top center;
    animation: lanternPhysicsSwing 3.6s infinite ease-in-out;
    z-index: 7;
  }

  .chain-suspension {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
  }

  .iron-link {
    width: 3px;
    height: 8px;
    background: #523f2f;
    border-radius: 2px;
  }

  .stained-lantern-housing {
    position: relative;
    width: 36px;
    height: 52px;
    background: #1f150e;
    border: 2px solid #6e4e31;
    border-radius: 5px;
    display: flex;
    flex-direction: column;
    align-items: center;
    box-shadow: 0 0 18px rgba(255, 160, 40, 0.4);
    transition: box-shadow 0.5s ease, border-color 0.5s ease;
  }

  .lantern-dome {
    width: 28px;
    height: 8px;
    background: #3d2919;
    border-radius: 4px 4px 0 0;
    border-bottom: 1px solid #140d07;
  }

  .glass-chamber {
    position: relative;
    width: 26px;
    height: 34px;
    background: radial-gradient(circle at 50% 60%, rgba(255, 230, 160, 0.3) 0%, rgba(255, 140, 30, 0.15) 100%);
    border: 1px solid rgba(255, 170, 50, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .iron-grate.vertical {
    position: absolute;
    width: 2px;
    height: 100%;
    background: #3f2a1b;
  }

  .iron-grate.horizontal {
    position: absolute;
    height: 2px;
    width: 100%;
    background: #3f2a1b;
  }

  .flame-spark {
    width: 9px;
    height: 15px;
    background: radial-gradient(circle at 50% 80%, #ffffff 0%, #ffaa00 60%, #cc4400 100%);
    border-radius: 50% 50% 30% 30%;
    box-shadow: 0 0 16px #ff7700, 0 0 30px #ff4400;
    animation: flamePulse 1.1s infinite alternate;
  }

  .radiant-rays {
    position: absolute;
    inset: -15px;
    background: radial-gradient(circle, rgba(255, 180, 50, 0.25) 0%, transparent 70%);
    pointer-events: none;
    animation: rayGleam 2s infinite alternate;
  }

  .stained-lantern-housing.blue-flame {
    box-shadow: 0 0 28px #00bfff;
    border-color: #38b0de;
  }

  .stained-lantern-housing.blue-flame .flame-spark {
    background: radial-gradient(circle at 50% 80%, #ffffff 0%, #00bfff 60%, #0055ff 100%);
    box-shadow: 0 0 20px #00bfff, 0 0 35px #0055ff;
  }

  /* ================== BOSS 2: ASYLUM SOVEREIGN ================== */
  .sovereign-container {
    animation: maniacNervousTwitch 2.6s infinite ease-in-out;
  }

  .back-quills-aura {
    position: absolute;
    top: 10px;
    inset: 0;
    pointer-events: none;
  }

  .quill-scrap {
    position: absolute;
    font-size: 1.4rem;
    animation: scrapFlutter 3s infinite alternate ease-in-out;
  }

  .quill-scrap.qs1 { top: 15px; left: 45px; transform: rotate(-25deg); }
  .quill-scrap.qs2 { top: 10px; right: 50px; transform: rotate(30deg); animation-delay: 1s; }
  .quill-scrap.qs3 { bottom: 65px; left: 35px; transform: rotate(-15deg); animation-delay: 1.8s; }
  .quill-scrap.qs4 { bottom: 70px; right: 40px; transform: rotate(20deg); animation-delay: 2.4s; }

  .straitjacket-torso {
    position: absolute;
    bottom: 5px;
    width: 170px;
    height: 135px;
  }

  .jacket-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 6px 14px rgba(0, 0, 0, 0.9));
  }

  .leather-straps {
    position: absolute;
    inset: 50px 25px 15px;
    display: flex;
    flex-direction: column;
    justify-content: space-around;
  }

  .strap {
    height: 6px;
    background: #3b2a1c;
    border-radius: 2px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.9);
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .buckle {
    width: 12px;
    height: 10px;
    background: #736149;
    border: 1px solid #bda88a;
    border-radius: 2px;
  }

  .deranged-head {
    position: absolute;
    top: 15px;
    width: 90px;
    height: 90px;
    display: flex;
    align-items: center;
    justify-content: center;
    transform: translate(var(--look-x), var(--look-y)) rotate(var(--look-rot));
    transition: transform 0.06s ease-out;
    z-index: 5;
  }

  .pale-skull {
    width: 74px;
    height: 80px;
    background: radial-gradient(circle at 45% 40%, #e2d7c5 0%, #beb199 70%, #7d725d 100%);
    border-radius: 40% 40% 50% 50%;
    border: 2px solid #4f4739;
    box-shadow: inset 0 0 12px rgba(0,0,0,0.6), 0 4px 12px rgba(0,0,0,0.9);
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .sunken-sockets {
    margin-top: 24px;
    width: 48px;
    display: flex;
    justify-content: space-between;
  }

  .maniac-eye {
    width: 13px;
    height: 11px;
    background: #110e0b;
    border-radius: 50%;
    box-shadow: inset 0 0 4px #000;
    position: relative;
    overflow: hidden;
  }

  .maniac-pupil {
    position: absolute;
    width: 5px;
    height: 5px;
    background: radial-gradient(circle, #ffffff 0%, #aa1111 60%, #440000 100%);
    border-radius: 50%;
    top: 3px;
    left: 4px;
    transform: translate(var(--pupil-x), var(--pupil-y));
    transition: transform 0.05s ease-out;
  }

  .twitching-mouth {
    margin-top: 12px;
    width: 32px;
    height: 10px;
    background: #1f0b0b;
    border-radius: 4px;
    border: 1px solid #3d1b1b;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    transition: height 0.2s ease;
  }

  .deranged-head.jaw-open .twitching-mouth {
    height: 18px;
    background: #420808;
  }

  .teeth-row {
    height: 2px;
    background: #ede6d8;
  }

  .sovereign-container.mad-laughter {
    transform: rotate(-6deg) translateY(-8px);
    filter: drop-shadow(0 0 16px rgba(180, 20, 20, 0.7));
  }

  /* ================== BOSS 3: ARCHMAGE ================== */
  .archmage-container {
    animation: infernoBreathing 3.6s infinite ease-in-out;
  }

  .archmage-inferno-glow {
    position: absolute;
    inset: -30px;
    background: radial-gradient(circle, rgba(255, 60, 0, 0.28) 0%, transparent 72%);
    animation: magmaStormPulse 2.4s infinite alternate ease-in-out;
    pointer-events: none;
  }

  .demon-horns-spikes {
    position: absolute;
    top: -5px;
    width: 220px;
    height: 110px;
    z-index: 4;
  }

  .horns-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 4px 10px rgba(0,0,0,0.95));
  }

  .colossal-torso {
    position: absolute;
    bottom: 5px;
    width: 200px;
    height: 145px;
  }

  .torso-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 8px 18px rgba(0,0,0,0.95));
  }

  .archmage-head {
    position: absolute;
    top: 15px;
    width: 90px;
    height: 90px;
    display: flex;
    align-items: center;
    justify-content: center;
    transform: translate(var(--look-x), var(--look-y)) rotate(var(--look-rot));
    transition: transform 0.06s ease-out;
    z-index: 6;
  }

  .obsidian-face {
    width: 76px;
    height: 82px;
    background: radial-gradient(circle at 45% 35%, #240c08 0%, #120503 70%, #050101 100%);
    border-radius: 40% 40% 45% 45%;
    border: 2px solid #5a190f;
    box-shadow: inset 0 0 12px #ff2200, 0 4px 14px rgba(0,0,0,0.95);
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .crimson-eyes-row {
    margin-top: 24px;
    width: 46px;
    display: flex;
    justify-content: space-between;
  }

  .fiery-eye {
    width: 12px;
    height: 8px;
    background: radial-gradient(circle, #fff 0%, #ff4500 70%, #660000 100%);
    border-radius: 50%;
    box-shadow: 0 0 10px #ff3300, 0 0 20px #ff0000;
    position: relative;
    overflow: hidden;
  }

  .eye-ember {
    position: absolute;
    width: 4px;
    height: 4px;
    background: #000;
    border-radius: 50%;
    top: 2px;
    left: 4px;
    transform: translate(var(--pupil-x), var(--pupil-y));
    transition: transform 0.05s ease-out;
  }

  .ash-maw {
    margin-top: 12px;
    width: 34px;
    height: 12px;
    background: #080201;
    border: 1px solid #7a2213;
    border-radius: 3px;
    box-shadow: inset 0 0 8px #ff3300;
    position: relative;
    transition: height 0.2s ease;
  }

  .archmage-head.jaw-open .ash-maw {
    height: 22px;
    box-shadow: inset 0 0 14px #ff5500, 0 0 10px #ff2200;
  }

  /* Bone Staff with Smoking Skull */
  .bone-staff-wrapper {
    position: absolute;
    top: 10px;
    right: 15px;
    width: 32px;
    height: 180px;
    display: flex;
    flex-direction: column;
    align-items: center;
    z-index: 7;
  }

  .staff-shaft {
    width: 6px;
    height: 140px;
    background: linear-gradient(180deg, #8a7a6c 0%, #473d34 70%, #211c18 100%);
    border-radius: 3px;
    box-shadow: 0 0 8px rgba(0,0,0,0.9);
  }

  .staff-skull-head {
    position: absolute;
    top: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .skull-bone {
    font-size: 1.8rem;
    filter: drop-shadow(0 0 8px #ff4500);
  }

  .staff-flame-ruby {
    position: absolute;
    top: 4px;
    width: 14px;
    height: 14px;
    background: radial-gradient(circle, #fff 0%, #ff2200 60%, #880000 100%);
    border-radius: 50%;
    box-shadow: 0 0 14px #ff3300, 0 0 26px #ff0000;
    animation: flamePulse 1.2s infinite alternate;
  }

  /* Reactions */
  .hit-reaction {
    animation: violentHitReaction 0.4s ease-out;
    filter: brightness(1.7) contrast(1.4);
  }

  .boss-dead {
    animation: bossDeathVanish 2s forwards ease-in-out;
  }

  @keyframes bossDeathVanish {
    0% { transform: scale(1) translateY(0); opacity: 1; filter: brightness(1.8) contrast(1.5); }
    25% { transform: scale(1.06) translateY(-8px) rotate(-2deg); opacity: 0.95; filter: drop-shadow(0 0 25px #ff2200) brightness(2.2); }
    60% { transform: scale(0.85) translateY(12px) rotate(3deg); opacity: 0.45; filter: grayscale(0.8) blur(3px); }
    100% { transform: scale(0.55) translateY(30px); opacity: 0.05; filter: grayscale(1) blur(10px); }
  }

  /* Keyframe Animations */
  @keyframes breathingIdle {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-5px); }
  }

  @keyframes maniacNervousTwitch {
    0%, 100% { transform: translate(0, 0); }
    25% { transform: translate(-3px, 2px) rotate(-1deg); }
    50% { transform: translate(3px, -2px) rotate(1.5deg); }
    75% { transform: translate(-2px, -2px); }
  }

  @keyframes infernoBreathing {
    0%, 100% { transform: translateY(0) scale(1); }
    50% { transform: translateY(-6px) scale(1.02); }
  }

  @keyframes lanternPhysicsSwing {
    0%, 100% { transform: rotate(-8deg); }
    50% { transform: rotate(9deg); }
  }

  @keyframes flamePulse {
    0% { transform: scale(0.9); }
    100% { transform: scale(1.18); }
  }

  @keyframes rayGleam {
    0% { opacity: 0.4; transform: scale(0.9); }
    100% { opacity: 0.95; transform: scale(1.15); }
  }

  @keyframes mistPulse {
    0% { opacity: 0.4; }
    100% { opacity: 0.9; }
  }

  @keyframes eyeBlink {
    0%, 96%, 100% { opacity: 1; transform: scaleY(1); }
    98% { opacity: 0; transform: scaleY(0.1); }
  }

  @keyframes violentHitReaction {
    0%, 100% { transform: translate(0, 0); }
    25% { transform: translate(-10px, 8px); }
    50% { transform: translate(10px, -8px); }
    75% { transform: translate(-5px, 5px); }
  }

  @keyframes scrapFlutter {
    0% { transform: rotate(-10deg) translateY(0); }
    100% { transform: rotate(15deg) translateY(-8px); }
  }

  @keyframes magmaStormPulse {
    0% { opacity: 0.6; transform: scale(0.95); }
    100% { opacity: 1; transform: scale(1.1); }
  }
</style>
