<script>
  import { audioEngine } from './audioEngine.js';

  export let score = 0;
  export let label = 'Очки';
  export let isObscured = false;

  let prevScore = score;
  let isScratching = false;

  $: if (!isObscured && score !== prevScore) {
    isScratching = true;
    audioEngine.playChalkScratch();
    setTimeout(() => {
      prevScore = score;
      isScratching = false;
    }, 280);
  }
</script>

<div 
  class="chalk-container" 
  class:obscured={isObscured}
  class:danger={!isObscured && score >= 18 && score <= 20} 
  class:crit={!isObscured && score === 21} 
  class:bust={!isObscured && score > 21}
>
  <div class="chalk-label">{label}</div>
  <div class="chalk-board">
    {#if isObscured}
      <span class="active-chalk-score obscured-chalk">?</span>
      <div class="lantern-mist-veil"></div>
    {:else}
      {#if isScratching && prevScore > 0}
        <span class="old-chalk-score">
          {#each String(prevScore).split('') as char}
            {#if char === '0'}
              <span class="chalk-zero">0</span>
            {:else}
              <span>{char}</span>
            {/if}
          {/each}
        </span>
        <div class="chalk-strike"></div>
      {/if}
      <span class="active-chalk-score" class:drawing={isScratching}>
        {#each String(score).split('') as char}
          {#if char === '0'}
            <span class="chalk-zero">0</span>
          {:else}
            <span>{char}</span>
          {/if}
        {/each}
      </span>
    {/if}
  </div>
</div>

<style>
  .chalk-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 8px 18px;
    background: radial-gradient(circle at 50% 40%, rgba(20, 15, 12, 0.9) 0%, rgba(8, 6, 5, 0.95) 100%);
    border: 2px solid #543b28;
    border-radius: 6px;
    box-shadow: 
      inset 0 0 14px rgba(0, 0, 0, 0.95),
      0 4px 14px rgba(0, 0, 0, 0.9);
    position: relative;
    min-width: 90px;
    transition: all 0.25s ease;
    font-variant-numeric: slashed-zero tabular-nums;
  }

  .chalk-label {
    font-size: 0.95rem;
    color: #ffd700;
    font-family: var(--font-body);
    font-weight: 700;
    letter-spacing: 2px;
    text-transform: uppercase;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.9);
    margin-bottom: 2px;
  }

  .chalk-board {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    min-width: 60px;
    height: 48px;
  }

  .active-chalk-score {
    font-family: 'Cinzel', 'Pirata One', serif;
    font-size: 2.7rem;
    font-weight: 900;
    line-height: 1;
    color: #ffffff;
    text-shadow: 
      0 0 8px rgba(0, 0, 0, 0.95),
      0 0 16px rgba(0, 0, 0, 0.9),
      2px 2px 0 #000,
      -2px -2px 0 #000,
      2px -2px 0 #000,
      -2px 2px 0 #000,
      0 0 10px rgba(255, 240, 210, 0.4);
    letter-spacing: 1px;
    transition: transform 0.2s ease, opacity 0.2s ease;
  }

  .chalk-zero {
    font-family: 'Cinzel', serif !important;
    font-weight: 900;
    display: inline-block;
    transform: scaleX(1.08);
    letter-spacing: 1px;
  }

  .active-chalk-score.drawing {
    animation: chalkWrite 0.28s ease-out;
  }

  .old-chalk-score {
    position: absolute;
    font-family: 'Cinzel', 'Pirata One', serif;
    font-size: 2.7rem;
    font-weight: 900;
    line-height: 1;
    color: #8c8275;
    opacity: 0.55;
    filter: blur(1px);
    text-shadow: 1px 1px 2px #000;
  }

  .chalk-strike {
    position: absolute;
    width: 120%;
    height: 4px;
    background: #ff5555;
    transform: rotate(-16deg);
    border-radius: 2px;
    box-shadow: 0 0 6px #ff2222;
    animation: strikeThrough 0.25s forwards;
  }

  /* Obscured by Lantern Beam */
  .chalk-container.obscured {
    border-color: #ff9900;
    background: radial-gradient(circle at 50% 50%, rgba(40, 20, 5, 0.95) 0%, rgba(10, 5, 2, 0.98) 100%);
    box-shadow: 
      0 0 20px rgba(255, 140, 0, 0.4),
      inset 0 0 16px rgba(0, 0, 0, 0.98);
  }

  .obscured-chalk {
    color: #ffaa33;
    animation: mysteryGlow 1.8s infinite alternate ease-in-out;
    text-shadow: 0 0 14px #ff7700, 2px 2px 0 #000, -2px -2px 0 #000;
  }

  .lantern-mist-veil {
    position: absolute;
    inset: -6px;
    background: radial-gradient(circle, rgba(255, 150, 0, 0.15) 0%, transparent 70%);
    pointer-events: none;
    animation: mistDrift 2s infinite alternate ease-in-out;
  }

  @keyframes mysteryGlow {
    0% { transform: scale(0.95); opacity: 0.8; filter: drop-shadow(0 0 4px #ff6600); }
    100% { transform: scale(1.08); opacity: 1; filter: drop-shadow(0 0 14px #ffaa00); }
  }

  @keyframes mistDrift {
    0% { transform: rotate(-3deg); }
    100% { transform: rotate(3deg); }
  }

  /* Danger state (18-20) */
  .chalk-container.danger {
    border-color: #d97724;
    box-shadow: 
      0 0 20px rgba(255, 120, 20, 0.35),
      inset 0 0 14px rgba(0, 0, 0, 0.95);
  }

  .chalk-container.danger .active-chalk-score {
    color: #ffb157;
    text-shadow: 
      0 0 8px rgba(0, 0, 0, 0.95),
      0 0 16px #ff7700,
      2px 2px 0 #000,
      -2px -2px 0 #000;
  }

  /* 21 Crit state */
  .chalk-container.crit {
    border-color: #38a5ff;
    box-shadow: 
      0 0 24px rgba(56, 165, 255, 0.6),
      inset 0 0 14px rgba(0, 0, 0, 0.95);
    animation: critGlowPulse 1.2s infinite alternate;
  }

  .chalk-container.crit .active-chalk-score {
    color: #a4e5ff;
    text-shadow: 
      0 0 8px rgba(0, 0, 0, 0.95),
      0 0 18px #00bfff,
      2px 2px 0 #000,
      -2px -2px 0 #000;
  }

  /* Bust state */
  .chalk-container.bust {
    border-color: #cc2222;
    box-shadow: 
      0 0 20px rgba(220, 30, 30, 0.5),
      inset 0 0 14px rgba(0, 0, 0, 0.95);
  }

  .chalk-container.bust .active-chalk-score {
    color: #ff6060;
    text-shadow: 
      0 0 8px rgba(0, 0, 0, 0.95),
      0 0 16px #ff0000,
      2px 2px 0 #000,
      -2px -2px 0 #000;
  }

  @keyframes chalkWrite {
    0% { transform: scale(0.8) rotate(-4deg); opacity: 0; }
    50% { transform: scale(1.15) rotate(2deg); opacity: 0.95; }
    100% { transform: scale(1) rotate(0deg); opacity: 1; }
  }

  @keyframes strikeThrough {
    0% { width: 0; opacity: 0; }
    100% { width: 120%; opacity: 0.9; }
  }

  @keyframes critGlowPulse {
    0% { box-shadow: 0 0 14px rgba(56, 165, 255, 0.4); }
    100% { box-shadow: 0 0 28px rgba(56, 165, 255, 0.8); }
  }
</style>
