<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';
  import Dice3D from './Dice3D.svelte';
  import VolumeControl from './VolumeControl.svelte';
  import CandleBar from './CandleBar.svelte';

  let { part2, playerCandles = 5, maxCandles = 5 } = $props();

  const dispatch = createEventDispatcher();

  let rolling = $state(false);
  let showRulesDrawer = $state(false);
  let showRelicsDrawer = $state(false);
  let showInkModal = $state(false);
  let selectedDieForInk = $state(0);
  let inkTargetFace = $state(6);
  let mouseParallaxX = $state(0);
  let mouseParallaxY = $state(0);
  let isBurningSacrifice = $state(false);
  let flameTargetingMode = $state(false);
  let medallionTargetingMode = $state(false);
  let burningDieIndex = $state(null);
  let isWitchDying = $state(false);
  let showBloodTax = $state(false);
  let prevResolved = $state(false);

  $effect(() => {
    if (part2 && part2.round_resolved && !prevResolved && part2.round_outcome === 'tie') {
      audioEngine.playBloodPayment();
      showBloodTax = true;
      setTimeout(() => {
        showBloodTax = false;
      }, 2400);
    }
    prevResolved = part2 ? part2.round_resolved : false;
  });

  onMount(() => {
    audioEngine.playDiceTumble();
  });

  function handleMouseMove(e) {
    const w = window.innerWidth || 1200;
    const h = window.innerHeight || 800;
    mouseParallaxX = (e.clientX / w - 0.5) * 20;
    mouseParallaxY = (e.clientY / h - 0.5) * 20;
  }

  function handleDieClick(index) {
    if (medallionTargetingMode) {
      confirmMedallionReroll(index);
      return;
    }
    if (flameTargetingMode) {
      confirmSacrifice(index);
      return;
    }
    toggleDieLock(index);
  }

  function toggleDieLock(index) {
    if (part2.round_resolved || part2.rerolls_left <= 0) return;
    const nextLocked = [...part2.locked_dice];
    nextLocked[index] = !nextLocked[index];
    part2.locked_dice = nextLocked;
  }

  function startFlameTargeting() {
    if (part2.sacrificial_flame_used || part2.round_resolved || rolling || isBurningSacrifice) return;
    medallionTargetingMode = false;
    flameTargetingMode = !flameTargetingMode;
    if (flameTargetingMode) {
      audioEngine.playWickIgnition();
    }
  }

  function cancelFlameTargeting() {
    flameTargetingMode = false;
  }

  function startMedallionTargeting() {
    if (part2.molten_medallion_used || part2.round_resolved || rolling || isBurningSacrifice) return;
    flameTargetingMode = false;
    medallionTargetingMode = !medallionTargetingMode;
    if (medallionTargetingMode) {
      audioEngine.playWickIgnition();
    }
  }

  function cancelMedallionTargeting() {
    medallionTargetingMode = false;
  }

  async function confirmMedallionReroll(dieIndex) {
    medallionTargetingMode = false;
    audioEngine.playDiceRattle();
    try {
      const res = await fetch('/api/part2/relic', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          relic: 'molten_medallion',
          die_index: dieIndex
        })
      });
      if (res.ok) {
        const data = await res.json();
        audioEngine.playDiceTumble();
        dispatch('update', data);
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Medallion reroll failed:', e);
    }
  }

  async function confirmSacrifice(dieIndex) {
    flameTargetingMode = false;
    burningDieIndex = dieIndex;
    await handleUseSacrificialFlame(dieIndex);
    burningDieIndex = null;
  }

  async function handleReroll() {
    cancelFlameTargeting();
    cancelMedallionTargeting();
    if (rolling || part2.rerolls_left <= 0 || part2.round_resolved) return;
    rolling = true;
    audioEngine.playDiceRattle();

    try {
      const lockedSlice = (part2.locked_dice || []).slice(0, part2.player_dice ? part2.player_dice.length : 5);
      const [res] = await Promise.all([
        fetch('/api/part2/roll', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ locked: lockedSlice })
        }),
        new Promise((resolve) => setTimeout(resolve, 450))
      ]);

      if (res.ok) {
        const data = await res.json();
        audioEngine.playDiceTumble();
        dispatch('update', data);
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Reroll failed:', e);
    } finally {
      rolling = false;
    }
  }

  async function handleCallBluff(callBluff) {
    cancelFlameTargeting();
    cancelMedallionTargeting();
    if (part2.round_resolved || rolling) return;
    audioEngine.playIronBolt();

    try {
      const res = await fetch('/api/part2/call-bluff', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ call_bluff: callBluff })
      });
      if (res.ok) {
        const data = await res.json();
        // If witch/opponent defeated (WitchHP <= 0), play scream and convulsions before advancing
        if (data.part2 && data.part2.witch_hp <= 0) {
          audioEngine.playWitchScream();
          isWitchDying = true;
          // Temporarily stay in battle so convulsions animation plays
          dispatch('update', { ...data, phase: 'part2_battle' });
          setTimeout(() => {
            isWitchDying = false;
            dispatch('update', data);
          }, 2600);
        } else {
          dispatch('update', data);
        }
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Call bluff failed:', e);
    }
  }

  async function handleNextRound() {
    cancelFlameTargeting();
    cancelMedallionTargeting();
    audioEngine.playDiceTumble();
    try {
      const res = await fetch('/api/part2/next-round', { method: 'POST' });
      if (res.ok) {
        const data = await res.json();
        dispatch('update', data);
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Next round failed:', e);
    }
  }

  async function handleUseEyeOfTruth() {
    if (part2.eye_of_truth_used || part2.round_resolved || rolling) return;
    audioEngine.playWickIgnition();
    try {
      const res = await fetch('/api/part2/relic', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ relic: 'eye_of_truth' })
      });
      if (res.ok) {
        const data = await res.json();
        dispatch('update', data);
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Eye of truth failed:', e);
    }
  }

  function openInkModal() {
    if (part2.ink_mark_used || part2.round_resolved || rolling) return;
    selectedDieForInk = 0;
    inkTargetFace = 6;
    showInkModal = true;
  }

  async function applyInkMark() {
    showInkModal = false;
    audioEngine.playChalkScratch();
    try {
      const res = await fetch('/api/part2/relic', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          relic: 'ink_mark',
          die_index: selectedDieForInk,
          target_face: inkTargetFace
        })
      });
      if (res.ok) {
        const data = await res.json();
        dispatch('update', data);
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Ink mark failed:', e);
    }
  }

  async function handleUseSacrificialFlame(dieIndex = 0) {
    if (part2.sacrificial_flame_used || part2.round_resolved || rolling || isBurningSacrifice) return;
    isBurningSacrifice = true;
    audioEngine.playInfernoRoar();

    try {
      const safeIndex = Math.min(Math.max(0, dieIndex), (part2.player_dice ? part2.player_dice.length - 1 : 0));
      const [res] = await Promise.all([
        fetch('/api/part2/relic', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            relic: 'sacrificial_flame',
            die_index: safeIndex
          })
        }),
        new Promise((resolve) => setTimeout(resolve, 500))
      ]);
      if (res.ok) {
        const data = await res.json();
        dispatch('update', data);
      } else {
        const syncRes = await fetch('/api/state');
        if (syncRes.ok) {
          const syncData = await syncRes.json();
          dispatch('update', syncData);
        }
      }
    } catch (e) {
      console.error('Sacrificial flame failed:', e);
    } finally {
      isBurningSacrifice = false;
      flameTargetingMode = false;
    }
  }
</script>

<div 
  class="part2-stage" 
  role="region" 
  aria-label="Ведьмино Око - Дуэль на костях" 
  on:mousemove={handleMouseMove}
>
  <!-- Basalt Vaults Background -->
  <div class="basalt-vault-bg">
    <div 
      class="ash-swirl-layer" 
      style="transform: translate({mouseParallaxX}px, {mouseParallaxY}px);"
    ></div>
  </div>

  <!-- Top Bar for Navigation & Volume Control -->
  <div class="part2-top-bar">
    <div class="top-spacer"></div>
    <VolumeControl />
  </div>

  <!-- Symmetrical 3-Column Grid -->
  <div class="part2-grid">
    <!-- LEFT: Stone Tablet with Dice Poker Rules (Large Readable Typography) -->
    <div class="altar-column left-tablet">
      <div class="tablet-card">
        <div class="tablet-header">
          <span class="tablet-sigil">📜</span>
          <h3 class="tablet-title">ОБРЯД НА КОСТЯХ</h3>
        </div>

        <div class="rules-scroll-body">
          <div class="rule-chip gold">
            <span class="chip-bold">Покер на 5 костях D6</span>
            <span class="chip-desc">Собери комбинацию выше противника или разоблачи его ложь.</span>
          </div>

          <div class="combinations-hierarchy">
            <div class="combo-item rank-8"><span class="c-rank">8</span> <strong>Ведьмино Око</strong> — 5 одинаковых</div>
            <div class="combo-item rank-7"><span class="c-rank">7</span> <strong>Каре</strong> — 4 одинаковых</div>
            <div class="combo-item rank-6"><span class="c-rank">6</span> <strong>Фулл-Хаус</strong> — 3 плюс 2</div>
            <div class="combo-item rank-5"><span class="c-rank">5</span> <strong>Большой Стрит</strong> — от 2 до 6</div>
            <div class="combo-item rank-4"><span class="c-rank">4</span> <strong>Малый Стрит</strong> — от 1 до 5</div>
            <div class="combo-item rank-3"><span class="c-rank">3</span> <strong>Тройка</strong> — 3 одинаковых</div>
            <div class="combo-item rank-2"><span class="c-rank">2</span> <strong>Две Пары</strong></div>
            <div class="combo-item rank-1"><span class="c-rank">1</span> <strong>Пара</strong></div>
            <div class="combo-item rank-0"><span class="c-rank">0</span> <strong>Старшая кость</strong></div>
          </div>

          <div class="rule-chip blood">
            <span class="chip-bold">Поражение:</span> гасит свечу.<br />
            <span class="chip-bold">Ничья:</span> заплати кровью.
          </div>
        </div>
      </div>
    </div>

    <!-- CENTER: Alena Silhouette, Obsidian Table, Dice -->
    <div class="altar-center">
      <!-- Upper Section: Expressive Alena Figure & Witch Dice -->
      <div class="witch-battle-zone">
        <div class="witch-avatar-plate" class:phase-2={part2.witch_phase === 2} class:phase-3={part2.witch_phase === 3} class:witch-convulsing={isWitchDying}>
          <div class="witch-title-badge">
            {#if part2.is_nightmare}
              {#if part2.witch_phase === 1}
                <span class="phase-title">Фаза I: Тень Прошлого</span>
              {:else if part2.witch_phase === 2}
                <span class="phase-title plague-text">Фаза II: Зеркальный Демон — Пепел Мора</span>
              {:else}
                <span class="phase-title agony-text">Фаза III: Поглощение Души — Ставки x2!</span>
              {/if}
            {:else}
              {#if part2.witch_phase === 1}
                <span class="phase-title">Фаза I: Голос Любимой</span>
              {:else if part2.witch_phase === 2}
                <span class="phase-title plague-text">Фаза II: Моровая Тварь — Кость Мора</span>
              {:else}
                <span class="phase-title agony-text">Фаза III: Агония Бездны — Ставки x2!</span>
              {/if}
            {/if}
          </div>

          <!-- Expressive Detailed Character Stage (Alena in Normal, Shadow Doppelganger in Nightmare) -->
          <div class="alena-character-stage" class:shrieking={isWitchDying}>
            {#if part2.is_nightmare}
              <!-- SHADOW DOPPELGANGER: Dark Mirror Hero in Hat with Cyan Glowing Eyes -->
              <svg viewBox="0 0 220 180" class="alena-svg doppelganger-svg" aria-label="Теневой Двойник">
                <defs>
                  <filter id="shadowAuraGlow" x="-30%" y="-30%" width="160%" height="160%">
                    <feDropShadow dx="0" dy="0" stdDeviation="6" flood-color="#00e5ff" flood-opacity="0.75" />
                    <feDropShadow dx="0" dy="0" stdDeviation="1.5" flood-color="#ffffff" flood-opacity="0.9" />
                  </filter>
                  <linearGradient id="doppelHat" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="#2a3848" />
                    <stop offset="100%" stop-color="#0f161e" />
                  </linearGradient>
                  <linearGradient id="doppelCloak" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="#1f2833" />
                    <stop offset="60%" stop-color="#10151c" />
                    <stop offset="100%" stop-color="#06090c" />
                  </linearGradient>
                </defs>

                <ellipse cx="110" cy="150" rx="80" ry="18" fill="rgba(0, 229, 255, 0.25)" filter="blur(8px)" />

                <g filter="url(#shadowAuraGlow)">
                  <!-- Hero Boots -->
                  <path d="M 92 145 L 88 172 L 102 172 L 98 145 Z" fill="#42291b" stroke="#1d120a" stroke-width="1.5" />
                  <path d="M 122 145 L 118 172 L 132 172 L 128 145 Z" fill="#42291b" stroke="#1d120a" stroke-width="1.5" />

                  <!-- Cloak with Jagged Shadows -->
                  <path d="M 110 70 Q 75 110 65 155 Q 110 162 155 155 Q 145 110 110 70 Z" fill="url(#doppelCloak)" stroke="#2b475e" stroke-width="1.8" />

                  <!-- Shadow Staff -->
                  <path d="M 68 35 L 60 168" stroke="#523a27" stroke-width="4.5" stroke-linecap="round" />
                  <circle cx="68" cy="32" r="6" fill="#00e5ff" filter="drop-shadow(0 0 6px #00e5ff)" />

                  <!-- Glowing Talisman on Chest -->
                  <circle cx="110" cy="95" r="7" fill="#d9e6f2" stroke="#ffffff" stroke-width="1.2" />
                  <polygon points="110,90 114,98 106,98" fill="#00e5ff" />

                  <!-- Face Oval in deep shadow -->
                  <circle cx="110" cy="58" r="14" fill="#090b0e" />
                  <!-- Cyan Piercing Eyes -->
                  <ellipse cx="104" cy="58" rx="3.5" ry="2" fill="#00e5ff" filter="drop-shadow(0 0 5px #00e5ff)" />
                  <ellipse cx="116" cy="58" rx="3.5" ry="2" fill="#00e5ff" filter="drop-shadow(0 0 5px #00e5ff)" />

                  <!-- Wizard Hat -->
                  <ellipse cx="110" cy="50" rx="30" ry="7" fill="url(#doppelHat)" stroke="#45607a" stroke-width="1.5" />
                  <polygon points="94,49 126,49 118,12 110,8" fill="url(#doppelHat)" stroke="#45607a" stroke-width="1.5" />
                </g>
              </svg>
            {:else}
              <!-- ALENA: Feminine Silhouette, Soft Violet Glow & Crisp Outer Contour, Slender Waist, High-Heeled Boots -->
              <svg viewBox="0 0 220 180" class="alena-svg" aria-label="Алёна Моровая Ведьма">
                <defs>
                  <filter id="alenaGlowContour" x="-30%" y="-30%" width="160%" height="160%">
                    <feDropShadow dx="0" dy="0" stdDeviation="5.5" flood-color="#cf5aff" flood-opacity="0.75" />
                    <feDropShadow dx="0" dy="0" stdDeviation="1.8" flood-color="#ffffff" flood-opacity="0.85" />
                  </filter>
                  <radialGradient id="alenaSkin" cx="50%" cy="40%" r="60%">
                    <stop offset="0%" stop-color="#fff5f0" />
                    <stop offset="60%" stop-color="#e8cfc2" />
                    <stop offset="100%" stop-color="#ba9c8e" />
                  </radialGradient>
                  <linearGradient id="alenaDress" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="#4a1954" />
                    <stop offset="40%" stop-color="#2c0c35" />
                    <stop offset="80%" stop-color="#18041f" />
                    <stop offset="100%" stop-color="#09010d" />
                  </linearGradient>
                  <linearGradient id="alenaHair" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stop-color="#2b0e35" />
                    <stop offset="50%" stop-color="#18061e" />
                    <stop offset="100%" stop-color="#08010b" />
                  </linearGradient>
                  <linearGradient id="alenaShoe" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stop-color="#280c2e" />
                    <stop offset="100%" stop-color="#100314" />
                  </linearGradient>
                </defs>

                <!-- Soft Floor Mist -->
                <ellipse cx="110" cy="148" rx="85" ry="20" fill="rgba(160, 40, 220, 0.3)" filter="blur(7px)" />

                <g filter="url(#alenaGlowContour)">
                  <!-- Slender Legs in Mesh Stockings & Pointed High-Heeled Gothic Shoes -->
                  <g class="alena-legs">
                    <!-- Left Leg & High Heel Stiletto -->
                    <path d="M 94 135 C 93 148, 92 158, 91 168 L 84 171 L 81 176 L 98 176 L 96 166 L 98 135 Z" fill="#1e1022" stroke="#5c246b" stroke-width="1.2" />
                    <!-- Left Stiletto Heel & Buckle -->
                    <line x1="83" y1="171" x2="83" y2="176" stroke="#d49aff" stroke-width="1.8" />
                    <rect x="88" y="167" width="5" height="3" rx="0.8" fill="#e8c7ff" stroke="#441a52" stroke-width="0.6" />

                    <!-- Right Leg & High Heel Stiletto -->
                    <path d="M 120 135 C 121 148, 122 158, 123 168 L 130 171 L 133 176 L 116 176 L 118 166 L 116 135 Z" fill="#1e1022" stroke="#5c246b" stroke-width="1.2" />
                    <!-- Right Stiletto Heel & Buckle -->
                    <line x1="131" y1="171" x2="131" y2="176" stroke="#d49aff" stroke-width="1.8" />
                    <rect x="121" y="167" width="5" height="3" rx="0.8" fill="#e8c7ff" stroke="#441a52" stroke-width="0.6" />
                  </g>

                  <!-- Billowing Gothic Gown with Accentuated Hourglass Slender Waist -->
                  <g class="alena-gown">
                    <!-- Flared Skirt with High Front Slit -->
                    <path d="M 98 94 Q 72 120 62 154 Q 110 162 158 154 Q 148 120 122 94 Z" fill="url(#alenaDress)" stroke="#87309c" stroke-width="1.5" />
                    <!-- Elegant Draped Folds -->
                    <path d="M 98 94 Q 92 125 84 154" stroke="#1f0724" stroke-width="2" fill="none" />
                    <path d="M 110 94 Q 110 130 110 158" stroke="#5d1e6e" stroke-width="1.8" fill="none" />
                    <path d="M 122 94 Q 128 125 136 154" stroke="#1f0724" stroke-width="2" fill="none" />

                    <!-- Cinched Feminine Corset Waist with Hourglass Taper -->
                    <path d="M 97 68 C 94 78, 93 86, 97 94 C 103 96, 117 96, 123 94 C 127 86, 126 78, 123 68 Z" fill="#3a1042" stroke="#aa3ec4" stroke-width="1.5" />
                    <!-- Gold/Violet Corset Ribs -->
                    <line x1="104" y1="72" x2="116" y2="76" stroke="#f1b5ff" stroke-width="1.2" />
                    <line x1="104" y1="80" x2="116" y2="84" stroke="#f1b5ff" stroke-width="1.2" />
                    <line x1="105" y1="88" x2="115" y2="90" stroke="#f1b5ff" stroke-width="1.2" />
                  </g>

                  <!-- Slender Pale Arms with Feminine Poise -->
                  <g class="alena-arms">
                    <path d="M 96 70 Q 75 84 80 112 Q 85 114 89 108 Q 85 92 98 78 Z" fill="url(#alenaSkin)" stroke="#6e2d4d" stroke-width="1" />
                    <!-- Left Hand & Violet Talons -->
                    <path d="M 80 112 Q 76 120 74 125 Q 78 126 82 120 Z" fill="#280c33" />
                    <path d="M 124 70 Q 145 84 140 112 Q 135 114 131 108 Q 135 92 122 78 Z" fill="url(#alenaSkin)" stroke="#6e2d4d" stroke-width="1" />
                    <!-- Right Hand & Violet Talons -->
                    <path d="M 140 112 Q 144 120 146 125 Q 142 126 138 120 Z" fill="#280c33" />
                  </g>

                  <!-- Long Voluminous Cascading Hair with Expressive Tresses -->
                  <path d="M 84 42 C 55 65, 48 105, 54 135 C 64 142, 74 125, 72 105 C 70 78, 88 60, 94 52 Z" fill="url(#alenaHair)" stroke="#4e1761" stroke-width="1" />
                  <path d="M 136 42 C 165 65, 172 105, 166 135 C 156 142, 146 125, 148 105 C 150 78, 132 60, 126 52 Z" fill="url(#alenaHair)" stroke="#4e1761" stroke-width="1" />
                  <!-- Flowing Highlights -->
                  <path d="M 82 55 Q 60 85 66 122" stroke="#9d42ba" stroke-width="1.5" fill="none" opacity="0.8" />
                  <path d="M 138 55 Q 160 85 154 122" stroke="#9d42ba" stroke-width="1.5" fill="none" opacity="0.8" />

                  <!-- Slender Neck & Velvet Choker -->
                  <path d="M 105 55 L 105 68 Q 110 70 115 68 L 115 55 Z" fill="url(#alenaSkin)" />
                  <rect x="105" y="61" width="10" height="3" fill="#16051c" stroke="#aa3ec4" stroke-width="0.5" />

                  <!-- Feminine Face Oval -->
                  <ellipse cx="110" cy="44" rx="15" ry="18" fill="url(#alenaSkin)" stroke="#5c1f43" stroke-width="1" />

                  <!-- Seductive / Mocking Facial Features -->
                  <g class="alena-face">
                    <path d="M 100 35 Q 104 33 108 36" stroke="#2b0a33" stroke-width="1.4" fill="none" stroke-linecap="round" />
                    <path d="M 112 36 Q 116 32 120 34" stroke="#2b0a33" stroke-width="1.4" fill="none" stroke-linecap="round" />
                    <!-- Slanted Glowing Eyes with Blinking Highlights -->
                    <ellipse cx="104" cy="40" rx="3.5" ry="2.2" fill="#d944ff" filter="drop-shadow(0 0 3px #d944ff)" />
                    <ellipse cx="116" cy="40" rx="3.5" ry="2.2" fill="#d944ff" filter="drop-shadow(0 0 3px #d944ff)" />
                    <ellipse cx="104" cy="40" rx="1" ry="2" fill="#08020a" />
                    <ellipse cx="116" cy="40" rx="1" ry="2" fill="#08020a" />
                    <circle cx="103.2" cy="39" r="0.8" fill="#fff" />
                    <circle cx="115.2" cy="39" r="0.8" fill="#fff" />

                    <!-- Petite Nose & Sarcastic Carmine Lips -->
                    <path d="M 110 40 L 110 47 L 112 48" stroke="#966453" stroke-width="1" fill="none" stroke-linecap="round" />
                    <path d="M 105 52 Q 110 55 115 52 Q 118 49 115 50 Q 110 52 105 52 Z" fill="#b0103b" />
                  </g>

                  <!-- Bangs & Forehead Strands -->
                  <path d="M 96 34 Q 106 45 102 53 Q 98 43 96 34 Z" fill="url(#alenaHair)" />
                  <path d="M 124 34 Q 114 45 118 53 Q 122 43 124 34 Z" fill="url(#alenaHair)" />
                  <path d="M 97 30 Q 110 22 123 30 Q 117 26 110 26 Q 103 26 97 30 Z" fill="url(#alenaHair)" />
                </g>
              </svg>
            {/if}
          </div>

          <!-- Witch Health Skulls -->
          <div class="witch-hp-bar">
            {#each Array(part2.witch_max_hp) as _, i}
              <span class="witch-hp-skull" class:dead={i >= part2.witch_hp}>
                {i < part2.witch_hp ? '💀' : '🌑'}
              </span>
            {/each}
            <span class="witch-hp-num">{part2.witch_hp} / {part2.witch_max_hp} HP</span>
          </div>

          <!-- Witch Dice Hand (Hidden or Revealed) -->
          <div class="witch-dice-row">
            {#each part2.witch_dice as val, i}
              <Dice3D 
                value={val}
                isWitch={true}
                isRevealed={part2.round_resolved || (part2.revealed_witch_indices && part2.revealed_witch_indices.includes(i))}
                disabled={true}
              />
            {/each}
          </div>

          <!-- Witch Speech Bubble Claim -->
          <div class="witch-claim-bubble" class:blood-demand={part2.round_resolved && part2.round_outcome === 'tie'}>
            <span class="speech-sigil">{part2.round_resolved && part2.round_outcome === 'tie' ? '🩸' : '💬'}</span>
            <span class="claim-text">
              {#if part2.round_resolved}
                {#if part2.round_outcome === 'tie'}
                  {#if part2.is_nightmare}
                    Твой Двойник скалится: <strong class="blood-demand-quote">«Ничья?! Заплати нашей общей кровью за каждый вздох!»</strong>
                  {:else}
                    Алёна злобно скалится: <strong class="blood-demand-quote">«Ничья не спасет тебя! Плати горячей кровью прямо на алтарь!»</strong>
                  {/if}
                {:else}
                  На костях {part2.is_nightmare ? 'Двойника' : 'Ведьмы'}: <strong>{part2.witch_combo_name}</strong>
                {/if}
              {:else}
                {part2.is_nightmare ? 'Твой Двойник ухмыляется' : 'Алёна усмехается'}: <strong>«У меня {part2.witch_claim_name}... рискнешь проверить?»</strong>
              {/if}
            </span>
          </div>
        </div>
      </div>

      <!-- Center Altar: Cracked Skull Candlestick & Action Message -->
      <div class="obsidian-hearth">
        <div class="cracked-skull-candlestick">
          <div class="skull-socket-flame left-flame"></div>
          <div class="skull-socket-flame right-flame"></div>
          <span class="skull-candle-glyph">💀</span>
        </div>

        <div class="part2-message-banner">
          <p class="part2-msg-text">{part2.last_message}</p>
        </div>
      </div>

      <!-- Lower Section: Player Hand, Combination Name, Action Bar -->
      <div 
        class="player-battle-zone" 
        class:sacrifice-burn={isBurningSacrifice || part2.sacrificial_flame_active}
        class:mortal-peril={(part2.witch_phase === 3 || part2.is_nightmare) && playerCandles <= 2}
      >
        {#if (part2.witch_phase === 3 || part2.is_nightmare) && playerCandles <= 2}
          <div class="mortal-stakes-banner" role="alert">
            <span class="mortal-icon">⚠️</span>
            <span class="mortal-text">СМЕРТЕЛЬНЫЙ КОН — СТАВКИ x2! ПОРАЖЕНИЕ ПОГАСИТ ВСЕ ОСТАВШИЕСЯ СВЕЧИ!</span>
          </div>
        {/if}

        <div class="player-combo-banner">
          <span class="combo-label">Твоя рука:</span>
          <span class="combo-value">{part2.player_combo_name}</span>
          {#if part2.sacrificial_flame_active}
            <span class="combo-multiplier">🔥 x2 УРОН ПО ВЕДЬМЕ (4 КОСТИ)</span>
          {/if}
        </div>

        {#if flameTargetingMode}
          <div class="flame-targeting-banner" role="alert">
            <span class="flame-banner-icon">🔥</span>
            <span class="flame-banner-text">ВЫБЕРИТЕ КОСТЬ ДЛЯ ЖЕРТВЫ: нажмите на любую кость ниже</span>
            <button class="flame-banner-cancel-btn" on:click={cancelFlameTargeting}>✕ Отмена</button>
          </div>
        {/if}

        {#if medallionTargetingMode}
          <div class="medallion-targeting-banner" role="alert">
            <span class="medallion-banner-icon">🎲</span>
            <span class="medallion-banner-text">ВЫБЕРИТЕ КОСТЬ ДЛЯ ПЕРЕБРОСА: нажмите на любую кость ниже</span>
            <button class="medallion-banner-cancel-btn" on:click={cancelMedallionTargeting}>✕ Отмена</button>
          </div>
        {/if}

        <!-- 5 or 4 D6 Player Dice Arc -->
        <div class="player-dice-arc" class:hand-burned={part2.sacrificial_flame_active} class:targeting-active={flameTargetingMode || medallionTargetingMode}>
          {#each part2.player_dice as die, i}
            <Dice3D 
              value={die}
              locked={part2.locked_dice[i]}
              rolling={rolling}
              disabled={part2.round_resolved || (isBurningSacrifice && burningDieIndex !== i)}
              isFlameTarget={flameTargetingMode}
              isMedallionTarget={medallionTargetingMode}
              isBeingBurned={isBurningSacrifice && burningDieIndex === i}
              index={i}
              on:toggleLock={(e) => handleDieClick(e.detail.index)}
            />
          {/each}
        </div>

        <!-- Action Buttons: "ПЕРЕБРОС", "ВСКРЫТЬ КОСТИ", "ТЫ ЛЖЕШЬ!" -->
        <div class="part2-controls-row">
          {#if !part2.round_resolved}
            <button 
              class="dread-btn reroll-btn" 
              on:click={handleReroll}
              disabled={rolling || part2.rerolls_left <= 0}
            >
              <span class="btn-title">ПЕРЕБРОС</span>
              <span class="btn-sub">Осталось: {part2.rerolls_left}</span>
            </button>

            <button 
              class="dread-btn showdown-btn" 
              on:click={() => handleCallBluff(false)}
              disabled={rolling}
            >
              <span class="btn-title">ВСКРЫТЬ КОСТИ</span>
              <span class="btn-sub">Принять бой</span>
            </button>

            <button 
              class="dread-btn bluff-btn" 
              on:click={() => handleCallBluff(true)}
              disabled={rolling}
            >
              <span class="btn-title">ТЫ ЛЖЕШЬ!</span>
              <span class="btn-sub">Уличить в блефе (x2 урон)</span>
            </button>
          {:else}
            <button class="dread-btn next-round-btn" on:click={handleNextRound}>
              <span>СЛЕДУЮЩИЙ КОН ➔</span>
            </button>
          {/if}
        </div>

        <!-- Bottom Candle Life Bar -->
        <div class="player-candles-footer">
          <CandleBar candles={playerCandles} maxCandles={maxCandles} />
        </div>
      </div>
    </div>

    <!-- RIGHT: Mutated Relics Cabinet (Strictly 1 use per battle, crossed out with red graphic cross) -->
    <div class="altar-column right-relics" class:drawer-open={showRelicsDrawer}>
      <div class="tablet-card">
        <div class="tablet-header">
          <span class="tablet-sigil">🗝️</span>
          <h3 class="tablet-title">РЕЛИКВИИ</h3>
          {#if showRelicsDrawer}
            <button class="drawer-close" on:click={() => showRelicsDrawer = false}>✕</button>
          {/if}
        </div>

        <div class="relics-list">
          <!-- 0. Molten Medallion -->
          <div class="relic-item-card medallion-relic" class:used={part2.molten_medallion_used}>
            {#if part2.molten_medallion_used}
              <div class="red-cross-overlay">
                <svg viewBox="0 0 100 100" class="red-cross-svg" aria-label="Использовано">
                  <line x1="12" y1="12" x2="88" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                  <line x1="88" y1="12" x2="12" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                </svg>
              </div>
            {/if}
            <div class="relic-item-top">
              <span class="relic-item-icon">📿</span>
              <div class="relic-item-name">Раскалённый медальон</div>
            </div>
            <p class="relic-item-desc">Связь с ведьмой жжёт грудь. Позволяет один раз перебросить любую кость в дуэли.</p>
            <button 
              class="use-relic-btn" 
              class:targeting={medallionTargetingMode}
              on:click={startMedallionTargeting}
              disabled={part2.molten_medallion_used || part2.round_resolved || rolling || isBurningSacrifice}
            >
              {#if part2.molten_medallion_used}
                ИСПОЛЬЗОВАНО
              {:else if medallionTargetingMode}
                ✕ ОТМЕНА
              {:else}
                Перебросить кость
              {/if}
            </button>
          </div>

          <!-- 1. Eye of Truth -->
          <div class="relic-item-card" class:used={part2.eye_of_truth_used}>
            {#if part2.eye_of_truth_used}
              <div class="red-cross-overlay">
                <svg viewBox="0 0 100 100" class="red-cross-svg" aria-label="Использовано">
                  <line x1="12" y1="12" x2="88" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                  <line x1="88" y1="12" x2="12" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                </svg>
              </div>
            {/if}
            <div class="relic-item-top">
              <span class="relic-item-icon">👁️</span>
              <div class="relic-item-name">Око Истины</div>
            </div>
            <p class="relic-item-desc">Раскрывает 3 закрытые кости противника до принятия решения.</p>
            <button 
              class="use-relic-btn" 
              on:click={handleUseEyeOfTruth}
              disabled={part2.eye_of_truth_used || part2.round_resolved}
            >
              {part2.eye_of_truth_used ? 'ИСПОЛЬЗОВАНО' : 'Подсмотреть'}
            </button>
          </div>

          <!-- 2. Ink Mark -->
          <div class="relic-item-card" class:used={part2.ink_mark_used}>
            {#if part2.ink_mark_used}
              <div class="red-cross-overlay">
                <svg viewBox="0 0 100 100" class="red-cross-svg" aria-label="Использовано">
                  <line x1="12" y1="12" x2="88" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                  <line x1="88" y1="12" x2="12" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                </svg>
              </div>
            {/if}
            <div class="relic-item-top">
              <span class="relic-item-icon">🖋️</span>
              <div class="relic-item-name">Чернильная Метка</div>
            </div>
            <p class="relic-item-desc">Перевернуть любой свой кубик на нужную грань от 1 до 6.</p>
            <button 
              class="use-relic-btn" 
              on:click={openInkModal}
              disabled={part2.ink_mark_used || part2.round_resolved}
            >
              {part2.ink_mark_used ? 'ИСПОЛЬЗОВАНО' : 'Нанести Метку'}
            </button>
          </div>

          <!-- 3. Sacrificial Flame -->
          <div class="relic-item-card" class:used={part2.sacrificial_flame_used}>
            {#if part2.sacrificial_flame_used}
              <div class="red-cross-overlay">
                <svg viewBox="0 0 100 100" class="red-cross-svg" aria-label="Использовано">
                  <line x1="12" y1="12" x2="88" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                  <line x1="88" y1="12" x2="12" y2="88" stroke="#e01010" stroke-width="14" stroke-linecap="round" />
                </svg>
              </div>
            {/if}
            <div class="relic-item-top">
              <span class="relic-item-icon">🔥</span>
              <div class="relic-item-name">Жертвенное Пламя</div>
            </div>
            <p class="relic-item-desc">Сжечь кость: дает дополнительный переброс и удвоенный урон по врагу.</p>
            <button 
              class="use-relic-btn" 
              class:targeting={flameTargetingMode}
              on:click={startFlameTargeting}
              disabled={part2.sacrificial_flame_used || part2.round_resolved || rolling}
            >
              {#if part2.sacrificial_flame_used}
                ИСПОЛЬЗОВАНО
              {:else if flameTargetingMode}
                ✕ ОТМЕНА
              {:else}
                Сжечь кость
              {/if}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Ink Mark Modal Popup (NO DROPDOWNS: Two Rows of Clickable Dice Cards) -->
  {#if showInkModal}
    <div 
      class="ink-modal-backdrop" 
      role="button" 
      tabindex="0" 
      on:click={() => showInkModal = false}
      on:keydown={(e) => e.key === 'Escape' && (showInkModal = false)}
    >
      <div 
        class="ink-modal-content" 
        role="dialog" 
        aria-modal="true" 
        tabindex="-1"
        on:click|stopPropagation 
        on:keydown|stopPropagation
      >
        <h4 class="ink-modal-title">🖋️ ЧЕРНИЛЬНАЯ МЕТКА</h4>
        <p class="ink-modal-desc">Выбери кость на столе и укажи желаемую грань:</p>

        <!-- ROW 1: Current player dice on the table -->
        <div class="ink-block">
          <div class="ink-section-label">РЯД 1: ТЕКУЩИЕ КОСТИ НА СТОЛЕ</div>
          <div class="dice-selection-row">
            {#each part2.player_dice as die, i}
              <button 
                type="button"
                class="dice-pick-card" 
                class:selected={selectedDieForInk === i} 
                on:click={() => selectedDieForInk = i}
              >
                <span class="dice-tag">Кость #{i + 1}</span>
                <div class="dice-box val-{die}">{die}</div>
              </button>
            {/each}
          </div>
        </div>

        <!-- ROW 2: 6 dice values from 1 to 6 in order -->
        <div class="ink-block">
          <div class="ink-section-label">РЯД 2: ЖЕЛАЕМАЯ ГРАНЬ (ОТ 1 ДО 6)</div>
          <div class="faces-selection-row">
            {#each [1, 2, 3, 4, 5, 6] as face}
              <button 
                type="button"
                class="dice-pick-card face-card" 
                class:selected={inkTargetFace === face} 
                on:click={() => inkTargetFace = face}
              >
                <div class="dice-box val-{face}">{face}</div>
              </button>
            {/each}
          </div>
        </div>

        <div class="ink-modal-actions">
          <button class="ink-cancel-btn" on:click={() => showInkModal = false}>Отмена</button>
          <button class="ink-apply-btn" on:click={applyInkMark}>Перевернуть на {inkTargetFace}</button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Blood Tax Visual Effect on Tie -->
  {#if showBloodTax}
    <div class="blood-tax-overlay" aria-hidden="true">
      <div class="blood-flash-screen"></div>
      <div class="blood-slash-cut"></div>
      <div class="blood-drips-row">
        <span class="blood-drop d1"></span>
        <span class="blood-drop d2"></span>
        <span class="blood-drop d3"></span>
        <span class="blood-drop d4"></span>
      </div>
      <div class="blood-tax-toast">
        <span class="toast-blood-icon">🩸</span>
        <span class="toast-blood-msg">НИЧЬЯ: ЗАПЛАТИ КРОВЬЮ!</span>
      </div>
    </div>
  {/if}
</div>

<style>
  .part2-stage {
    position: relative;
    width: 100%;
    height: 100%;
    background: radial-gradient(circle at 50% 50%, #150a1a 0%, #07030a 70%, #020104 100%);
    overflow: hidden;
    color: #ffd700;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    font-variant-numeric: slashed-zero tabular-nums;
  }

  .basalt-vault-bg {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 1;
    background-image: 
      radial-gradient(ellipse at 50% 20%, rgba(90, 15, 110, 0.25) 0%, transparent 60%),
      radial-gradient(circle at 10% 80%, rgba(50, 0, 70, 0.3) 0%, transparent 50%),
      radial-gradient(circle at 90% 80%, rgba(50, 0, 70, 0.3) 0%, transparent 50%);
  }

  .ash-swirl-layer {
    position: absolute;
    inset: -30px;
    background-image: repeating-linear-gradient(
      115deg,
      rgba(180, 50, 220, 0.03) 0px,
      rgba(180, 50, 220, 0.03) 1px,
      transparent 2px,
      transparent 24px
    );
    opacity: 0.8;
  }

  .part2-top-bar {
    position: relative;
    z-index: 50;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
    background: rgba(10, 4, 12, 0.85);
    border-bottom: 1px solid #4a1c52;
  }

  .top-nav-buttons {
    display: flex;
    gap: 8px;
  }

  .tablet-btn {
    background: linear-gradient(180deg, #2b1133 0%, #14071a 100%);
    border: 1px solid #6b267d;
    border-radius: 16px;
    color: #ffd700;
    padding: 6px 14px;
    font-family: var(--font-title);
    font-size: 1rem;
    cursor: pointer;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.7);
    transition: all 0.2s ease;
  }

  .tablet-btn:hover {
    border-color: #ba47db;
    transform: translateY(-1px);
  }

  /* Grid Layout: Fixed 280px side columns matching part 1 */
  .part2-grid {
    position: relative;
    z-index: 10;
    display: grid;
    grid-template-columns: 280px minmax(560px, 1fr) 280px;
    width: 100%;
    height: calc(100vh - 52px);
    padding: 8px 14px;
    gap: 12px;
    box-sizing: border-box;
  }

  .altar-column {
    height: 100%;
    width: 280px;
    min-width: 280px;
    max-width: 280px;
    overflow: hidden;
  }

  .tablet-card {
    background: linear-gradient(180deg, #1b0e21 0%, #0d0612 100%);
    border: 2px solid #5a266b;
    border-radius: 8px;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 14px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.9);
  }

  .tablet-header {
    display: flex;
    align-items: center;
    gap: 8px;
    border-bottom: 1px solid #75298a;
    padding-bottom: 8px;
    margin-bottom: 10px;
  }

  .tablet-title {
    font-family: var(--font-title);
    font-size: 1.6rem;
    color: #ffd700;
    margin: 0;
    letter-spacing: 1.5px;
    flex-grow: 1;
    text-shadow: 0 0 10px rgba(255, 215, 0, 0.35);
  }

  .drawer-close {
    background: transparent;
    border: 1px solid #6a2e7a;
    color: #ebd7be;
    border-radius: 50%;
    width: 26px;
    height: 26px;
    cursor: pointer;
  }

  .rules-scroll-body {
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .rule-chip {
    padding: 8px 10px;
    border-radius: 6px;
    font-size: 0.95rem;
    line-height: 1.35;
  }

  .rule-chip.gold {
    background: rgba(60, 25, 75, 0.6);
    border: 1px solid #82329e;
    color: #f7e6ce;
  }

  .rule-chip.blood {
    background: rgba(60, 10, 30, 0.6);
    border: 1px solid #8f2347;
    color: #ff9e9e;
  }

  .chip-bold {
    font-weight: bold;
    color: #ffd700;
  }

  .chip-desc {
    display: block;
    margin-top: 2px;
    color: #e5be68;
  }

  /* Large legible font sizes in poker hierarchy */
  .combinations-hierarchy {
    display: flex;
    flex-direction: column;
    gap: 5px;
    font-size: 0.95rem;
  }

  .combo-item {
    background: rgba(25, 12, 32, 0.75);
    border: 1px solid #481a54;
    border-radius: 4px;
    padding: 6px 10px;
    display: flex;
    align-items: center;
    gap: 8px;
    color: #e5be68;
  }

  .combo-item strong {
    color: #ffd700;
  }

  .c-rank {
    font-family: var(--font-title);
    font-weight: bold;
    font-size: 1.15rem;
    color: #f7a8ff;
    width: 16px;
  }

  /* Center Stage */
  .altar-center {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    height: 100%;
    padding: 0 8px;
  }

  .witch-battle-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

  .witch-avatar-plate {
    background: rgba(20, 8, 26, 0.88);
    border: 2px solid #662278;
    border-radius: 12px;
    padding: 8px 18px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.95);
    position: relative;
  }

  .phase-title {
    font-family: var(--font-title);
    font-size: 1.1rem;
    color: #ffd700;
    letter-spacing: 1px;
  }

  .plague-text {
    color: #8ce065;
    text-shadow: 0 0 8px #3b8a1c;
  }

  .agony-text {
    color: #ff55bb;
    text-shadow: 0 0 10px #ff2299;
    animation: agonyPulse 1s infinite alternate;
  }

  @keyframes agonyPulse {
    0% { transform: scale(1); }
    100% { transform: scale(1.05); }
  }

  /* Alena Character Stage */
  .alena-character-stage {
    width: 220px;
    height: 155px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    animation: alenaBreathe 4s infinite ease-in-out;
  }

  .alena-character-stage.shrieking {
    animation: witchDeathConvulsions 0.1s infinite alternate !important;
  }

  .witch-avatar-plate.witch-convulsing {
    animation: witchDeathPlateShake 0.1s infinite alternate !important;
  }

  .alena-svg {
    width: 100%;
    height: 100%;
    filter: drop-shadow(0 4px 14px rgba(0, 0, 0, 0.95)) drop-shadow(0 0 20px rgba(200, 60, 255, 0.45));
  }

  .alena-svg.doppelganger-svg {
    filter: drop-shadow(0 4px 14px rgba(0, 0, 0, 0.95)) drop-shadow(0 0 20px rgba(0, 229, 255, 0.5));
  }

  @keyframes alenaBreathe {
    0%, 100% { transform: translateY(0) scale(1); }
    50% { transform: translateY(-4px) scale(1.02); }
  }

  @keyframes witchDeathConvulsions {
    0% { transform: scale(1.08) translate(-6px, 4px) rotate(-4deg); filter: brightness(2) contrast(1.6) drop-shadow(0 0 35px #ff0055); }
    100% { transform: scale(1.12) translate(6px, -4px) rotate(4deg); filter: brightness(2.8) contrast(1.8) drop-shadow(0 0 50px #bf00ff); }
  }

  @keyframes witchDeathPlateShake {
    0% { transform: translateY(-2px) rotate(-1deg); }
    100% { transform: translateY(3px) rotate(1deg); }
  }

  .witch-hp-bar {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .witch-hp-skull {
    font-size: 1.35rem;
    filter: drop-shadow(0 0 6px #ba30c4);
  }

  .witch-hp-skull.dead {
    filter: grayscale(1);
    opacity: 0.35;
  }

  .witch-hp-num {
    font-family: var(--font-title);
    font-size: 1.15rem;
    color: #ffd700;
    margin-left: 6px;
    font-weight: bold;
  }

  .witch-dice-row {
    display: flex;
    gap: 10px;
    margin: 2px 0;
  }

  .witch-claim-bubble {
    background: rgba(35, 12, 45, 0.92);
    border: 1px solid #9432ad;
    border-radius: 18px;
    padding: 6px 18px;
    font-size: 1rem;
    display: flex;
    align-items: center;
    gap: 8px;
    color: #ffdafc;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.7);
  }

  .witch-claim-bubble strong {
    color: #ffd700;
  }

  /* Obsidian Hearth with Skull Candle */
  .obsidian-hearth {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    margin: 2px 0;
  }

  .cracked-skull-candlestick {
    position: relative;
    width: 44px;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .skull-candle-glyph {
    font-size: 2.2rem;
    filter: drop-shadow(0 0 14px #a028d6);
  }

  .skull-socket-flame {
    position: absolute;
    width: 8px;
    height: 14px;
    background: radial-gradient(ellipse at 50% 80%, #fff 0%, #c45eff 50%, #750fa8 100%);
    border-radius: 50% 50% 30% 30%;
    box-shadow: 0 0 10px #c45eff;
    animation: flameFlicker 0.25s infinite alternate ease-in-out;
  }

  .left-flame { top: 8px; left: 11px; }
  .right-flame { top: 8px; right: 11px; }

  @keyframes flameFlicker {
    0% { transform: scale(0.9); opacity: 0.8; }
    100% { transform: scale(1.15) translateY(-2px); opacity: 1; }
  }

  .part2-message-banner {
    background: rgba(18, 7, 24, 0.92);
    border: 1px solid #6b267d;
    border-radius: 6px;
    padding: 6px 20px;
    max-width: 620px;
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.85);
  }

  .part2-msg-text {
    font-size: 1.08rem;
    margin: 0;
    text-align: center;
    color: #ffd700;
  }

  /* Player Battle Zone */
  .player-battle-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    gap: 6px;
    transition: all 0.3s ease;
  }

  .player-combo-banner {
    display: flex;
    align-items: center;
    gap: 8px;
    background: rgba(25, 10, 32, 0.85);
    border: 1px solid #75298a;
    border-radius: 14px;
    padding: 5px 16px;
  }

  .combo-label {
    font-size: 0.95rem;
    color: #e5be68;
  }

  .combo-value {
    font-family: var(--font-title);
    font-size: 1.3rem;
    color: #ffd700;
    font-weight: bold;
  }

  .combo-multiplier {
    font-size: 0.92rem;
    color: #ff5500;
    font-weight: bold;
    text-shadow: 0 0 8px #ff4400;
    animation: flameMultiplier 1s infinite alternate;
  }

  @keyframes flameMultiplier {
    0% { transform: scale(1); }
    100% { transform: scale(1.08); }
  }

  .player-dice-arc {
    display: flex;
    gap: 14px;
    align-items: center;
    justify-content: center;
    min-height: 74px;
    padding: 4px 0;
    transition: all 0.3s ease;
  }

  .part2-controls-row {
    display: flex;
    gap: 12px;
    width: 100%;
    max-width: 580px;
    justify-content: center;
    margin-top: 4px;
  }

  .reroll-btn {
    border-color: #8c3b9e;
    flex: 1;
  }

  .showdown-btn {
    border-color: #4a7542;
    background: linear-gradient(180deg, #1f3b1c 0%, #0d1f0d 100%);
    flex: 1;
  }

  .showdown-btn:hover {
    border-color: #72cc62;
    box-shadow: 0 0 20px rgba(90, 200, 70, 0.4);
  }

  .bluff-btn {
    border-color: #992440;
    background: linear-gradient(180deg, #3d121c 0%, #1f070d 100%);
    flex: 1;
  }

  .bluff-btn:hover {
    border-color: #f24e72;
    box-shadow: 0 0 20px rgba(240, 70, 110, 0.45);
  }

  .next-round-btn {
    flex: 2;
    border-color: #3b8ca8;
    background: linear-gradient(180deg, #163642 0%, #08171f 100%);
  }

  .btn-title {
    font-size: 1.25rem;
    font-weight: 900;
    display: block;
    color: #ffd700;
  }

  .btn-sub {
    font-size: 0.72rem;
    opacity: 0.9;
    text-transform: uppercase;
  }

  .player-candles-footer {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Right Relics Column */
  .relics-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    overflow-y: auto;
  }

  .relic-item-card {
    background: rgba(25, 12, 33, 0.85);
    border: 1px solid #6b267d;
    border-radius: 8px;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 5px;
    position: relative;
    overflow: hidden;
    transition: all 0.25s ease;
  }

  /* Red graphic cross when used */
  .relic-item-card.used {
    border-color: #701111;
    background: rgba(18, 5, 8, 0.9);
  }

  .red-cross-overlay {
    position: absolute;
    inset: 0;
    z-index: 20;
    pointer-events: none;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(10, 2, 4, 0.45);
  }

  .red-cross-svg {
    width: 85%;
    height: 85%;
    filter: drop-shadow(0 0 8px #ff0000);
  }

  .relic-item-top {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .relic-item-icon {
    font-size: 1.5rem;
  }

  .relic-item-name {
    font-family: var(--font-title);
    font-size: 1.2rem;
    color: #ffd700;
    font-weight: bold;
    letter-spacing: 0.8px;
  }

  .relic-item-desc {
    font-size: 0.88rem;
    line-height: 1.35;
    color: #e5be68;
    margin: 2px 0 6px;
  }

  .use-relic-btn {
    padding: 8px 14px;
    background: linear-gradient(180deg, #451b54 0%, #1f0a26 100%);
    border: 1px solid #8432a1;
    border-radius: 6px;
    color: #ffd700;
    font-family: var(--font-title);
    font-size: 1.05rem;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.2s ease;
    text-transform: uppercase;
  }

  .use-relic-btn:hover:not(:disabled) {
    border-color: #ffd700;
    box-shadow: 0 0 14px rgba(255, 215, 0, 0.4);
    transform: translateY(-1px);
  }

  .use-relic-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
    border-color: #551111;
    color: #884444;
  }

  .use-relic-btn.targeting {
    background: linear-gradient(180deg, #991e00 0%, #4a0c00 100%);
    border-color: #ff6600;
    color: #fff099;
    box-shadow: 0 0 16px rgba(255, 80, 0, 0.8);
    animation: flameButtonPulse 1s infinite alternate;
  }

  @keyframes flameButtonPulse {
    0% { box-shadow: 0 0 8px rgba(255, 60, 0, 0.5); }
    100% { box-shadow: 0 0 18px rgba(255, 120, 0, 0.95); }
  }

  /* Flame Targeting Banner */
  .flame-targeting-banner {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 8px 16px;
    background: linear-gradient(90deg, rgba(80, 15, 0, 0.95), rgba(160, 40, 0, 0.95), rgba(80, 15, 0, 0.95));
    border: 1px solid #ff6600;
    border-radius: 8px;
    box-shadow: 0 0 22px rgba(255, 60, 0, 0.6);
    margin: 0 auto 12px;
    max-width: 650px;
    animation: flameBannerGlow 1.2s infinite alternate;
  }

  @keyframes flameBannerGlow {
    0% { border-color: #ff4400; box-shadow: 0 0 12px rgba(255, 50, 0, 0.4); }
    100% { border-color: #ffbb00; box-shadow: 0 0 24px rgba(255, 90, 0, 0.85); }
  }

  .flame-banner-icon {
    font-size: 1.3rem;
  }

  .flame-banner-text {
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 0.92rem;
    font-weight: bold;
    color: #fff3c4;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.9);
    letter-spacing: 0.04em;
  }

  .flame-banner-cancel-btn {
    background: rgba(0, 0, 0, 0.65);
    border: 1px solid #ff7744;
    color: #ffccaa;
    border-radius: 4px;
    padding: 4px 10px;
    font-size: 0.8rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .flame-banner-cancel-btn:hover {
    background: #ff3300;
    color: #ffffff;
    box-shadow: 0 0 10px #ff3300;
  }

  /* Medallion Targeting Banner */
  .medallion-targeting-banner {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 8px 16px;
    background: linear-gradient(90deg, rgba(70, 10, 5, 0.95), rgba(150, 25, 0, 0.95), rgba(70, 10, 5, 0.95));
    border: 1px solid #ff4400;
    border-radius: 8px;
    box-shadow: 0 0 22px rgba(255, 40, 0, 0.65);
    margin: 0 auto 12px;
    max-width: 650px;
    animation: medallionBannerGlow 1.2s infinite alternate;
  }

  @keyframes medallionBannerGlow {
    0% { border-color: #ff3300; box-shadow: 0 0 12px rgba(255, 40, 0, 0.4); }
    100% { border-color: #ff8833; box-shadow: 0 0 26px rgba(255, 70, 0, 0.9); }
  }

  .medallion-banner-icon {
    font-size: 1.3rem;
  }

  .medallion-banner-text {
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 0.92rem;
    font-weight: bold;
    color: #fff0d8;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.9);
    letter-spacing: 0.04em;
  }

  .medallion-banner-cancel-btn {
    background: rgba(0, 0, 0, 0.65);
    border: 1px solid #ff6633;
    color: #ffddbb;
    border-radius: 4px;
    padding: 4px 10px;
    font-size: 0.8rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .medallion-banner-cancel-btn:hover {
    background: #ff2200;
    color: #ffffff;
  }

  .relic-item-card.medallion-relic {
    border-color: #7a2214;
    background: linear-gradient(135deg, rgba(40, 10, 12, 0.9) 0%, rgba(20, 5, 8, 0.95) 100%);
  }

  .relic-item-card.medallion-relic:hover:not(.used) {
    border-color: #ff4411;
    box-shadow: 0 0 18px rgba(255, 50, 0, 0.4);
  }

  /* Mortal Stakes x2 Warning */
  .mortal-stakes-banner {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 6px 14px;
    background: linear-gradient(90deg, rgba(120, 10, 10, 0.95), rgba(190, 20, 20, 0.95), rgba(120, 10, 10, 0.95));
    border: 1px solid #ff3333;
    border-radius: 6px;
    margin: 0 auto 10px;
    max-width: 650px;
    animation: mortalPulse 1s ease-in-out infinite alternate;
    box-shadow: 0 0 18px rgba(255, 0, 0, 0.7);
  }

  .mortal-icon {
    font-size: 1.1rem;
    animation: flicker 0.4s infinite alternate;
  }

  .mortal-text {
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 0.85rem;
    font-weight: 800;
    color: #fff2d1;
    letter-spacing: 0.05em;
    text-shadow: 0 1px 3px #000;
  }

  @keyframes mortalPulse {
    0% { border-color: #ff3333; box-shadow: 0 0 10px rgba(255, 0, 0, 0.5); }
    100% { border-color: #ff9999; box-shadow: 0 0 25px rgba(255, 0, 0, 0.9); }
  }

  .player-battle-zone.mortal-peril {
    border-color: rgba(255, 50, 50, 0.5);
    box-shadow: 0 0 24px rgba(255, 0, 0, 0.25);
  }

  /* Ink Modal: Two Rows of Clickable Dice Cards */
  .ink-modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(4, 2, 6, 0.9);
    backdrop-filter: blur(5px);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    animation: fadeIn 0.2s ease-out;
  }

  .ink-modal-content {
    background: linear-gradient(180deg, #220e2b 0%, #100615 100%);
    border: 2px solid #a83dd4;
    border-radius: 12px;
    padding: 24px;
    max-width: 540px;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.98), 0 0 25px rgba(180, 50, 220, 0.3);
  }

  .ink-modal-title {
    font-family: var(--font-title);
    font-size: 1.6rem;
    color: #ffd700;
    margin: 0;
    letter-spacing: 1.2px;
    text-align: center;
  }

  .ink-modal-desc {
    font-size: 0.95rem;
    color: #e5be68;
    text-align: center;
    margin: 0;
  }

  .ink-block {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .ink-section-label {
    font-family: var(--font-title);
    font-size: 1.05rem;
    color: #ffda96;
    letter-spacing: 1px;
  }

  .dice-selection-row, .faces-selection-row {
    display: flex;
    gap: 10px;
    justify-content: center;
    flex-wrap: wrap;
  }

  .dice-pick-card {
    background: #1a0a24;
    border: 2px solid #5a266b;
    border-radius: 8px;
    padding: 8px 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    cursor: pointer;
    transition: all 0.2s ease;
    min-width: 60px;
  }

  .dice-pick-card:hover {
    border-color: #ba47db;
    transform: translateY(-2px);
  }

  .dice-pick-card.selected {
    border-color: #ffd700;
    background: #3b144a;
    box-shadow: 0 0 16px rgba(255, 215, 0, 0.5), inset 0 0 8px rgba(255, 215, 0, 0.3);
    transform: scale(1.08);
  }

  .dice-tag {
    font-size: 0.75rem;
    color: #e5be68;
    margin-bottom: 4px;
    font-weight: bold;
  }

  .dice-box {
    width: 44px;
    height: 44px;
    background: radial-gradient(circle at 40% 30%, #f6ecd6 0%, #dfcd9f 80%, #b89d6c 100%);
    border: 2px solid #4a341f;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: var(--font-title);
    font-size: 1.8rem;
    font-weight: 900;
    color: #170d06;
    box-shadow: inset 0 0 6px rgba(0,0,0,0.4), 0 3px 6px rgba(0,0,0,0.8);
  }

  .dice-pick-card.face-card {
    min-width: 54px;
    padding: 10px 8px;
  }

  .ink-modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 10px;
  }

  .ink-cancel-btn {
    background: transparent;
    border: 1px solid #5a266b;
    color: #a48c9e;
    padding: 8px 18px;
    border-radius: 6px;
    font-family: var(--font-title);
    font-size: 1rem;
    cursor: pointer;
  }

  .ink-apply-btn {
    background: linear-gradient(180deg, #611e75 0%, #300c3b 100%);
    border: 2px solid #ffd700;
    color: #ffd700;
    padding: 8px 22px;
    border-radius: 6px;
    font-family: var(--font-title);
    font-size: 1.1rem;
    font-weight: bold;
    cursor: pointer;
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.8);
    transition: all 0.2s ease;
  }

  .ink-apply-btn:hover {
    background: #7d2696;
    color: #fff;
    box-shadow: 0 0 16px rgba(255, 215, 0, 0.4);
    transform: translateY(-1px);
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  /* Mobile Adaptation */
  @media (max-width: 900px) {
    .part2-grid {
      grid-template-columns: 1fr;
      padding: 6px;
    }

    .altar-column {
      display: none;
    }

    .altar-column.drawer-open {
      display: flex;
      position: fixed;
      inset: 0;
      z-index: 500;
      background: rgba(3, 1, 5, 0.9);
      padding: 20px;
      align-items: center;
      justify-content: center;
    }

    .altar-column.drawer-open .tablet-card {
      max-width: 360px;
      width: 100%;
      max-height: 85vh;
    }

    .player-dice-arc {
      gap: 6px;
    }

    .part2-controls-row {
      max-width: 100%;
      padding: 0 4px;
      gap: 6px;
    }

    .btn-title {
      font-size: 1.05rem;
    }

    .btn-sub {
      font-size: 0.65rem;
    }
  }

  /* ================== BLOOD TAX VISUAL EFFECTS ================== */
  .witch-claim-bubble.blood-demand {
    border-color: #ff2222;
    background: radial-gradient(circle at 50% 50%, #3d0808 0%, #1c0303 80%, #0d0101 100%);
    box-shadow: 0 0 20px rgba(255, 34, 34, 0.6), inset 0 0 10px rgba(255, 0, 0, 0.4);
    animation: bloodPulse 0.8s ease-in-out infinite alternate;
  }

  .blood-demand-quote {
    color: #ff4d4d;
    text-shadow: 0 0 10px #ff0000;
  }

  @keyframes bloodPulse {
    from { transform: scale(1); box-shadow: 0 0 16px rgba(255, 34, 34, 0.45); }
    to { transform: scale(1.02); box-shadow: 0 0 28px rgba(255, 34, 34, 0.8); }
  }

  .blood-tax-overlay {
    position: fixed;
    inset: 0;
    pointer-events: none;
    z-index: 99999;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .blood-flash-screen {
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 50% 50%, rgba(180, 0, 0, 0.3) 0%, rgba(120, 0, 0, 0.7) 65%, rgba(60, 0, 0, 0.95) 100%);
    animation: bloodFlashAnim 2.4s ease-out forwards;
  }

  .blood-slash-cut {
    position: absolute;
    width: 65%;
    max-width: 600px;
    height: 4px;
    background: linear-gradient(90deg, transparent 0%, #ff1a1a 30%, #ffffff 50%, #ff1a1a 70%, transparent 100%);
    transform: rotate(-12deg);
    box-shadow: 0 0 20px #ff0000, 0 0 40px #ff0000;
    animation: slashCutAnim 2.4s ease-out forwards;
  }

  .blood-drips-row {
    position: absolute;
    top: 25%;
    display: flex;
    gap: 40px;
  }

  .blood-drop {
    width: 12px;
    height: 34px;
    background: linear-gradient(180deg, #ff1a1a 0%, #8b0000 70%, #400000 100%);
    border-radius: 0 0 8px 8px;
    box-shadow: 0 0 8px rgba(255, 0, 0, 0.8);
    animation: bloodDrip 1.8s ease-in forwards;
  }

  .blood-drop.d1 { animation-delay: 0.1s; height: 42px; }
  .blood-drop.d2 { animation-delay: 0.25s; height: 60px; }
  .blood-drop.d3 { animation-delay: 0.15s; height: 50px; }
  .blood-drop.d4 { animation-delay: 0.3s; height: 35px; }

  .blood-tax-toast {
    position: relative;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 12px;
    background: rgba(20, 2, 2, 0.92);
    border: 2px solid #ff2222;
    border-radius: 8px;
    padding: 12px 24px;
    box-shadow: 0 0 30px rgba(255, 0, 0, 0.7), inset 0 0 15px rgba(255, 0, 0, 0.5);
    animation: toastPop 2.4s ease-out forwards;
  }

  .toast-blood-icon {
    font-size: 2rem;
    filter: drop-shadow(0 0 8px #ff0000);
  }

  .toast-blood-msg {
    font-family: var(--font-title, 'Cinzel', serif);
    font-size: 1.4rem;
    font-weight: 900;
    color: #fff;
    letter-spacing: 2px;
    text-shadow: 0 0 12px #ff2222, 0 2px 4px #000;
  }

  @keyframes bloodFlashAnim {
    0% { opacity: 0; }
    15% { opacity: 1; }
    60% { opacity: 0.8; }
    100% { opacity: 0; }
  }

  @keyframes slashCutAnim {
    0% { transform: scaleX(0) rotate(-12deg); opacity: 0; }
    10% { transform: scaleX(1) rotate(-12deg); opacity: 1; }
    70% { transform: scaleX(1) rotate(-12deg); opacity: 0.9; }
    100% { transform: scaleX(1.1) rotate(-12deg); opacity: 0; }
  }

  @keyframes bloodDrip {
    0% { transform: translateY(0); opacity: 0; }
    20% { opacity: 1; }
    100% { transform: translateY(180px); opacity: 0; }
  }

  @keyframes toastPop {
    0% { transform: scale(0.8); opacity: 0; }
    15% { transform: scale(1.05); opacity: 1; }
    25% { transform: scale(1); opacity: 1; }
    80% { transform: scale(1); opacity: 1; }
    100% { transform: scale(0.95); opacity: 0; }
  }
</style>
