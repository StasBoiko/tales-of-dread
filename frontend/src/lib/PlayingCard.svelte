<script>
  import { onMount } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  export let card;
  export let index = 0;
  export let isDealer = false;

  const suitIcons = {
    skulls: '☠',
    bones: '🦴',
    crows: '🪶',
    chains: '⛓',
    ink: '🖋️',
  };

  const suitColors = {
    skulls: '#991b1b',
    bones: '#856f54',
    crows: '#1e283d',
    chains: '#3f3933',
    ink: '#1c152e',
  };

  onMount(() => {
    audioEngine.playCardThud();
  });
</script>

<div 
  class="dread-card"
  class:caged={card.isCaged}
  class:blinded={card.isBlinded}
  class:ink-blot={card.isInkBlot}
  class:dealer-card={isDealer}
  style="--deal-delay: {index * 0.08}s;"
  aria-label="{card.name}, номинал {card.baseValue}"
>
  <!-- Card Face -->
  <div class="card-face" class:ink-splashed={card.isInkBlot}>
    {#if card.isBlinded}
      <!-- Blinded by Lantern silhouette -->
      <div class="lantern-blinded-curtain">
        <div class="lantern-silhouette-icon">🏮</div>
        <div class="blinded-label">СЛЕПЯЩИЙ ЛУЧ</div>
        <div class="blinded-subtext">номинал скрыт</div>
      </div>
    {:else if card.isInkBlot}
      <!-- Ink Blot Card: Single central ink splatter with large score -->
      <div class="ink-blot-container">
        <div class="corner top-left">
          <span class="card-rank ink-val">{card.baseValue >= 0 ? `+${card.baseValue}` : `${card.baseValue}`}</span>
        </div>

        <div class="ink-splatter-graphic">
          <svg viewBox="0 0 100 100" class="ink-svg">
            <path d="M50 20 C60 10, 75 25, 70 40 C85 45, 90 60, 78 72 C80 85, 65 92, 50 82 C35 95, 20 85, 22 70 C10 60, 15 42, 30 38 C25 22, 40 10, 50 20 Z" fill="#140f24" />
            <circle cx="25" cy="20" r="4" fill="#140f24" />
            <circle cx="82" cy="28" r="3" fill="#140f24" />
            <circle cx="78" cy="88" r="5" fill="#140f24" />
            <circle cx="15" cy="78" r="3.5" fill="#140f24" />
          </svg>
          <div class="ink-score-badge">
            {card.baseValue >= 0 ? `+${card.baseValue}` : `${card.baseValue}`}
          </div>
          <div class="card-name-banner ink-name">КЛЯКСА</div>
        </div>
      </div>
    {:else}
      <!-- Standard Card: Single large central illustration, corner number top-left only -->
      <div class="corner top-left">
        <span class="card-rank">{card.rank}</span>
      </div>

      <div class="card-center">
        <div class="suit-emblem" style="color: {suitColors[card.suit] || '#333'}">
          {suitIcons[card.suit] || '♠'}
        </div>
        <div class="card-name-banner">{card.name}</div>
      </div>
    {/if}

    <!-- Paper Cage Effect -->
    {#if card.isCaged}
      <div class="paper-cage-ink">
        <div class="cross-line line1"></div>
        <div class="cross-line line2"></div>
        <div class="halved-stamp">½ СИЛЫ</div>
      </div>
    {/if}
  </div>
</div>

<style>
  .dread-card {
    position: relative;
    width: 96px;
    height: 126px;
    border-radius: 7px;
    perspective: 800px;
    backface-visibility: hidden;
    will-change: transform, opacity;
    animation: cardDeal 0.35s cubic-bezier(0.18, 0.89, 0.32, 1.28) forwards;
    animation-delay: var(--deal-delay, 0s);
    opacity: 0;
    transform: translateY(-40px) scale(0.85);
    box-shadow: 
      0 6px 16px rgba(0, 0, 0, 0.9),
      0 0 1px rgba(0, 0, 0, 0.95);
    transition: transform 0.2s cubic-bezier(0.2, 0.8, 0.3, 1), box-shadow 0.2s ease;
    user-select: none;
    font-variant-numeric: slashed-zero tabular-nums;
  }

  .dread-card:hover {
    transform: translateY(-6px) scale(1.05);
    box-shadow: 
      0 12px 24px rgba(0, 0, 0, 0.95),
      0 0 12px rgba(255, 180, 80, 0.45);
    z-index: 25;
  }

  .dealer-card {
    filter: brightness(0.92) contrast(1.05);
  }

  .card-face {
    width: 100%;
    height: 100%;
    background: radial-gradient(circle at 50% 50%, #f6ecd6 0%, #dfcd9f 70%, #c4af7e 100%);
    border: 2px solid #6b5338;
    border-radius: 7px;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    padding: 4px 6px;
    position: relative;
    overflow: hidden;
    box-shadow: 
      inset 0 0 14px rgba(70, 45, 18, 0.35),
      inset 0 0 1px rgba(255, 255, 255, 0.5);
  }

  .card-face.ink-splashed {
    background: radial-gradient(circle at 45% 45%, #e1d8cf 0%, #c7bcb0 70%, #a4978c 100%);
    border-color: #3f3148;
  }

  .corner {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    font-family: 'Cinzel', serif;
    line-height: 0.9;
    z-index: 2;
  }

  /* Big Bold Numerals with clear open vertical void for 0 */
  .card-rank {
    font-family: 'Cinzel', serif;
    font-size: 1.7rem;
    font-weight: 900;
    color: #160e08;
    text-shadow: 
      0 1px 1px rgba(255, 255, 255, 0.6),
      0 0 2px rgba(0, 0, 0, 0.5);
    letter-spacing: -0.5px;
  }

  .card-rank.ink-val {
    font-size: 1.4rem;
    color: #3b1654;
    text-shadow: 0 1px 1px rgba(255, 255, 255, 0.7);
  }

  .card-center {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    position: relative;
    flex: 1;
    margin-top: -6px;
    z-index: 2;
  }

  /* Single large central illustration */
  .suit-emblem {
    font-size: 2.5rem;
    line-height: 1;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.55));
    transition: transform 0.2s ease;
  }

  .dread-card:hover .suit-emblem {
    transform: scale(1.08);
  }

  .card-name-banner {
    font-size: 0.56rem;
    color: #2b1f14;
    font-family: var(--font-body);
    font-weight: 800;
    letter-spacing: 0.4px;
    text-align: center;
    margin-top: 2px;
    max-width: 100%;
    text-transform: uppercase;
    text-shadow: 0 1px 0 rgba(255, 255, 255, 0.4);
    white-space: normal;
    line-height: 1.1;
    word-break: break-word;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* Blinded by Lantern silhouette */
  .lantern-blinded-curtain {
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 50% 45%, #1f1308 0%, #0c0805 70%, #030201 100%);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 10px;
    z-index: 15;
    border: 2px solid #543315;
  }

  .lantern-silhouette-icon {
    font-size: 2.8rem;
    filter: drop-shadow(0 0 10px #ff8800) drop-shadow(0 0 20px #ff3b00);
    animation: lanternSway 2s ease-in-out infinite alternate;
  }

  @keyframes lanternSway {
    0% { transform: rotate(-5deg) scale(0.95); filter: drop-shadow(0 0 8px #ff7700); }
    100% { transform: rotate(5deg) scale(1.05); filter: drop-shadow(0 0 16px #ffa200); }
  }

  .blinded-label {
    font-family: var(--font-title);
    font-size: 0.8rem;
    color: #ffaa33;
    font-weight: bold;
    letter-spacing: 1px;
    margin-top: 6px;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.95);
    text-align: center;
  }

  .blinded-subtext {
    font-size: 0.7rem;
    color: #9c826b;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  /* Ink Blot Card Layout */
  .ink-blot-container {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    position: relative;
    z-index: 3;
  }

  .ink-splatter-graphic {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
  }

  .ink-svg {
    position: absolute;
    width: 82px;
    height: 82px;
    filter: drop-shadow(0 2px 5px rgba(0, 0, 0, 0.6));
  }

  .ink-score-badge {
    position: relative;
    font-family: var(--font-title);
    font-size: 2.3rem;
    font-weight: 900;
    color: #ffffff;
    z-index: 5;
    text-shadow: 0 2px 6px rgba(0, 0, 0, 0.95);
  }

  .card-name-banner.ink-name {
    position: relative;
    z-index: 5;
    color: #fff4e0;
    font-size: 0.75rem;
    background: rgba(20, 10, 35, 0.8);
    padding: 1px 6px;
    border-radius: 4px;
    margin-top: 2px;
  }

  /* Paper Cage Inked Cross */
  .paper-cage-ink {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 20;
  }

  .cross-line {
    position: absolute;
    top: 50%;
    left: 8%;
    width: 84%;
    height: 6px;
    background: #060506;
    box-shadow: 0 0 3px rgba(0, 0, 0, 0.95);
    border-radius: 3px;
  }

  .cross-line.line1 {
    transform: rotate(44deg);
  }

  .cross-line.line2 {
    transform: rotate(-44deg);
  }

  .halved-stamp {
    position: absolute;
    bottom: 14px;
    left: 50%;
    transform: translateX(-50%) rotate(-8deg);
    background: #991b1b;
    color: #fff;
    font-family: var(--font-title);
    font-size: 0.85rem;
    font-weight: bold;
    padding: 3px 8px;
    border-radius: 3px;
    letter-spacing: 1.2px;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.9);
    border: 1px solid #ff7b7b;
  }

  @keyframes cardDeal {
    0% { opacity: 0; transform: translateY(-40px) scale(0.85); }
    100% { opacity: 1; transform: translateY(0) scale(1); }
  }
</style>
