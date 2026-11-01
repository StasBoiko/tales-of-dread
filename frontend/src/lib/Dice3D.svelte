<script>
  import { createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  export let value = 1;
  export let locked = false;
  export let rolling = false;
  export let disabled = false;
  export let isPlague = false;
  export let isWitch = false;
  export let isRevealed = true;
  export let index = 0;
  export let isFlameTarget = false;
  export let isMedallionTarget = false;
  export let isBeingBurned = false;

  const dispatch = createEventDispatcher();

  function handleClick() {
    if (disabled || isWitch) return;
    if (navigator.vibrate) {
      try { navigator.vibrate(35); } catch (e) {}
    }
    audioEngine.playChalkScratch();
    dispatch('toggleLock', { index });
  }

  // Pre-calculated rotation targets for 1..6
  const faceRotations = {
    1: 'rotateX(0deg) rotateY(0deg)',
    2: 'rotateX(-90deg) rotateY(0deg)',
    3: 'rotateY(90deg) rotateX(0deg)',
    4: 'rotateY(-90deg) rotateX(0deg)',
    5: 'rotateX(90deg) rotateY(0deg)',
    6: 'rotateX(180deg) rotateY(0deg)',
  };

  $: rotationStyle = faceRotations[value] || faceRotations[1];
</script>

<div 
  class="dice-container"
  class:locked={locked && !isWitch}
  class:rolling={rolling}
  class:is-plague={isPlague}
  class:is-witch={isWitch}
  class:hidden-witch={isWitch && !isRevealed}
  class:flame-target={isFlameTarget && !isWitch}
  class:medallion-target={isMedallionTarget && !isWitch}
  class:being-burned={isBeingBurned}
  on:click={handleClick}
  role="button"
  tabindex="0"
  on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && handleClick()}
  aria-label={isWitch && !isRevealed ? "Закрытая кость Ведьмы" : (isMedallionTarget ? `Перебросить кость ${value}` : (isFlameTarget ? `Сжечь кость ${value}` : `Кость ${value}, ${locked ? 'зафиксирована' : 'не зафиксирована'}`))}
>
  {#if isFlameTarget && !isWitch}
    <div class="flame-target-indicator" title="Нажмите, чтобы сжечь эту кость">
      <span class="flame-icon">🔥</span>
      <span class="flame-text">СЖЕЧЬ</span>
    </div>
  {/if}
  {#if isMedallionTarget && !isWitch}
    <div class="medallion-target-indicator" title="Нажмите, чтобы перебросить эту кость">
      <span class="medallion-icon">🎲</span>
      <span class="medallion-text">ПЕРЕБРОС</span>
    </div>
  {/if}
  {#if isWitch && !isRevealed}
    <!-- Hidden Witch Bone under claw / skull cup -->
    <div class="hidden-skull-die">
      <span class="skull-glyph">☠</span>
      <span class="rune-question">?</span>
    </div>
  {:else}
    <div 
      class="dice-cube" 
      style="transform: {rolling ? 'rotateX(720deg) rotateY(1080deg)' : rotationStyle};"
    >
      <!-- Face 1 (Front) -->
      <div class="dice-face face-1">
        <span class="pip pip-center">✦</span>
      </div>
      <!-- Face 2 (Bottom) -->
      <div class="dice-face face-2">
        <span class="pip pip-tl">✦</span>
        <span class="pip pip-br">✦</span>
      </div>
      <!-- Face 3 (Right) -->
      <div class="dice-face face-3">
        <span class="pip pip-tl">✦</span>
        <span class="pip pip-center">✦</span>
        <span class="pip pip-br">✦</span>
      </div>
      <!-- Face 4 (Left) -->
      <div class="dice-face face-4">
        <span class="pip pip-tl">✦</span>
        <span class="pip pip-tr">✦</span>
        <span class="pip pip-bl">✦</span>
        <span class="pip pip-br">✦</span>
      </div>
      <!-- Face 5 (Top) -->
      <div class="dice-face face-5">
        <span class="pip pip-tl">✦</span>
        <span class="pip pip-tr">✦</span>
        <span class="pip pip-center">✦</span>
        <span class="pip pip-bl">✦</span>
        <span class="pip pip-br">✦</span>
      </div>
      <!-- Face 6 (Back) -->
      <div class="dice-face face-6">
        <span class="pip pip-tl">✦</span>
        <span class="pip pip-tr">✦</span>
        <span class="pip pip-ml">✦</span>
        <span class="pip pip-mr">✦</span>
        <span class="pip pip-bl">✦</span>
        <span class="pip pip-br">✦</span>
      </div>
    </div>
  {/if}

  {#if locked && !isWitch}
    <div class="lock-indicator">
      <span class="lock-icon">🔒</span>
    </div>
  {/if}

  {#if isPlague}
    <div class="plague-curse-badge" title="Кость Мора: штраф к счету!">
      <span>☣</span>
    </div>
  {/if}
</div>

<style>
  .dice-container {
    width: 64px;
    height: 64px;
    position: relative;
    perspective: 800px;
    cursor: pointer;
    user-select: none;
    transition: transform 0.2s cubic-bezier(0.18, 0.89, 0.32, 1.28);
  }

  .dice-container:hover:not(.is-witch) {
    transform: translateY(-4px) scale(1.08);
  }

  .dice-container.locked {
    filter: drop-shadow(0 0 10px #ffb326);
  }

  .dice-cube {
    width: 100%;
    height: 100%;
    position: relative;
    transform-style: preserve-3d;
    transition: transform 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  .dice-container.rolling .dice-cube {
    transition: transform 0.5s cubic-bezier(0.4, 0.0, 0.2, 1);
  }

  .dice-face {
    position: absolute;
    width: 64px;
    height: 64px;
    background: radial-gradient(circle at 40% 40%, #f4ead5 0%, #d8c39e 65%, #b59868 100%);
    border: 2px solid #5a3c22;
    border-radius: 10px;
    box-shadow: 
      inset 0 0 8px rgba(70, 40, 15, 0.4),
      inset 0 1px 1px rgba(255, 255, 255, 0.6),
      0 4px 10px rgba(0, 0, 0, 0.7);
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-template-rows: repeat(3, 1fr);
    padding: 6px;
    box-sizing: border-box;
    backface-visibility: hidden;
  }

  .is-witch .dice-face {
    background: radial-gradient(circle at 40% 40%, #301f2e 0%, #1c0e1b 65%, #0e050d 100%);
    border-color: #662b59;
  }

  .is-plague .dice-face {
    background: radial-gradient(circle at 40% 40%, #1e331e 0%, #102110 65%, #060e06 100%);
    border-color: #438234;
  }

  /* 3D transforms for 6 faces (Z-offset is half cube width = 32px) */
  .face-1 { transform: rotateY(0deg) translateZ(32px); }
  .face-6 { transform: rotateX(180deg) translateZ(32px); }
  .face-3 { transform: rotateY(-90deg) translateZ(32px); }
  .face-4 { transform: rotateY(90deg) translateZ(32px); }
  .face-5 { transform: rotateX(-90deg) translateZ(32px); }
  .face-2 { transform: rotateX(90deg) translateZ(32px); }

  /* Pip styling */
  .pip {
    font-size: 1.15rem;
    color: #4a1c12;
    display: flex;
    align-items: center;
    justify-content: center;
    line-height: 1;
    filter: drop-shadow(0 1px 1px rgba(0, 0, 0, 0.5));
  }

  .is-witch .pip {
    color: #f7a8e8;
    filter: drop-shadow(0 0 4px #b8339c);
  }

  .is-plague .pip {
    color: #8be065;
    filter: drop-shadow(0 0 4px #44a81e);
  }

  .pip-center { grid-column: 2; grid-row: 2; }
  .pip-tl { grid-column: 1; grid-row: 1; }
  .pip-tr { grid-column: 3; grid-row: 1; }
  .pip-ml { grid-column: 1; grid-row: 2; }
  .pip-mr { grid-column: 3; grid-row: 2; }
  .pip-bl { grid-column: 1; grid-row: 3; }
  .pip-br { grid-column: 3; grid-row: 3; }

  /* Hidden Witch Die */
  .hidden-skull-die {
    width: 64px;
    height: 64px;
    background: radial-gradient(circle at 50% 50%, #20121f 0%, #0d060d 75%, #050205 100%);
    border: 2px dashed #7a3a6b;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.85);
  }

  .skull-glyph {
    font-size: 1.4rem;
    color: #b05c9b;
    filter: drop-shadow(0 0 6px #8f2e7b);
  }

  .rune-question {
    font-size: 0.8rem;
    color: #e8a2d5;
    font-weight: bold;
  }

  /* Lock Badge */
  .lock-indicator {
    position: absolute;
    top: -6px;
    right: -6px;
    background: #ffb833;
    border: 1px solid #ffffff;
    border-radius: 50%;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.9);
    z-index: 10;
  }

  .lock-icon {
    font-size: 0.72rem;
  }

  .plague-curse-badge {
    position: absolute;
    bottom: -6px;
    left: -6px;
    background: #2b7524;
    border: 1px solid #75e865;
    border-radius: 50%;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #ffffff;
    font-size: 0.75rem;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.9);
    z-index: 10;
  }

  /* Flame Targeting Mode */
  .dice-container.flame-target {
    cursor: pointer;
    animation: flamePulse 1.1s ease-in-out infinite alternate;
  }

  .dice-container.flame-target:hover {
    transform: translateY(-8px) scale(1.08);
    filter: drop-shadow(0 0 16px rgba(255, 80, 0, 0.95)) drop-shadow(0 0 28px rgba(255, 20, 0, 0.75));
  }

  /* Medallion Targeting Mode */
  .dice-container.medallion-target {
    cursor: pointer;
    animation: medallionPulse 1.1s ease-in-out infinite alternate;
  }

  .dice-container.medallion-target:hover {
    transform: translateY(-8px) scale(1.08);
    filter: drop-shadow(0 0 18px rgba(255, 40, 20, 0.95)) drop-shadow(0 0 28px rgba(255, 90, 0, 0.8));
  }

  @keyframes medallionPulse {
    0% {
      filter: drop-shadow(0 0 6px rgba(255, 60, 20, 0.6));
    }
    100% {
      filter: drop-shadow(0 0 18px rgba(255, 40, 0, 0.95)) drop-shadow(0 0 26px rgba(255, 120, 0, 0.7));
    }
  }

  .medallion-target-indicator {
    position: absolute;
    top: -16px;
    left: 50%;
    transform: translateX(-50%);
    background: radial-gradient(circle, #ff2a00 0%, #700600 100%);
    border: 1px solid #ffaa44;
    border-radius: 12px;
    padding: 2px 8px;
    display: flex;
    align-items: center;
    gap: 3px;
    justify-content: center;
    box-shadow: 0 0 14px rgba(255, 60, 0, 0.9);
    z-index: 12;
    animation: flameBounce 0.8s ease-in-out infinite alternate;
    pointer-events: none;
  }

  .medallion-icon {
    font-size: 0.85rem;
  }

  .medallion-text {
    font-family: 'Cinzel', serif;
    font-size: 0.65rem;
    font-weight: 700;
    color: #fff4d0;
    letter-spacing: 0.05em;
    text-shadow: 0 1px 2px #000;
  }

  @keyframes flamePulse {
    0% {
      filter: drop-shadow(0 0 6px rgba(255, 120, 0, 0.6));
    }
    100% {
      filter: drop-shadow(0 0 18px rgba(255, 60, 0, 0.95)) drop-shadow(0 0 26px rgba(255, 0, 0, 0.7));
    }
  }

  .flame-target-indicator {
    position: absolute;
    top: -16px;
    left: 50%;
    transform: translateX(-50%);
    background: radial-gradient(circle, #ff5500 0%, #991100 100%);
    border: 1px solid #ffcc00;
    border-radius: 12px;
    padding: 2px 8px;
    display: flex;
    align-items: center;
    gap: 3px;
    justify-content: center;
    box-shadow: 0 0 12px rgba(255, 80, 0, 0.9);
    z-index: 12;
    animation: flameBounce 0.8s ease-in-out infinite alternate;
    pointer-events: none;
  }

  .flame-icon {
    font-size: 0.85rem;
  }

  .flame-text {
    font-family: 'Cinzel', serif;
    font-size: 0.65rem;
    font-weight: 700;
    color: #fff4d0;
    letter-spacing: 0.05em;
    text-shadow: 0 1px 2px #000;
  }

  @keyframes flameBounce {
    0% { transform: translateX(-50%) translateY(0); }
    100% { transform: translateX(-50%) translateY(-4px); }
  }

  /* Disintegration burning animation */
  .dice-container.being-burned {
    animation: disintegrateBurn 0.5s forwards ease-in-out;
  }

  @keyframes disintegrateBurn {
    0% { transform: scale(1); filter: drop-shadow(0 0 20px #ff3300) brightness(1.5); opacity: 1; }
    50% { transform: scale(1.15) rotate(12deg); filter: drop-shadow(0 0 35px #ff0000) brightness(3); opacity: 0.8; }
    100% { transform: scale(0.15) rotate(45deg); filter: drop-shadow(0 0 45px #ff0000) brightness(5); opacity: 0; }
  }

  @media (max-width: 900px) {
    .dice-container {
      width: 58px;
      height: 58px;
    }
    .dice-face {
      width: 58px;
      height: 58px;
    }
    .face-1 { transform: rotateY(0deg) translateZ(29px); }
    .face-6 { transform: rotateX(180deg) translateZ(29px); }
    .face-3 { transform: rotateY(-90deg) translateZ(29px); }
    .face-4 { transform: rotateY(90deg) translateZ(29px); }
    .face-5 { transform: rotateX(-90deg) translateZ(29px); }
    .face-2 { transform: rotateX(90deg) translateZ(29px); }
    .hidden-skull-die {
      width: 58px;
      height: 58px;
    }
  }
</style>
