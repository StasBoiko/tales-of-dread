<script>
  import { onMount } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  let { candles = 5, maxCandles = 5 } = $props();

  let prevCandles = $state(candles);
  let isCandleHissing = $state(false);

  $effect(() => {
    if (candles < prevCandles) {
      isCandleHissing = true;
      audioEngine.playCandleHiss();
      prevCandles = candles;
    } else if (candles > prevCandles) {
      prevCandles = candles;
    }
  });
</script>

<div class="candle-bar-container" title="Твоя жизненная сила: {candles}/{maxCandles}">
  <div class="candle-row">
    {#each Array(maxCandles) as _, i}
      {@const isLit = i < candles}
      <div class="candle-stub {isLit ? 'lit' : 'extinguished'}">
        {#if isLit}
          <div class="flame-wrapper">
            <div class="flame-halo"></div>
            <div class="flame-core"></div>
            <div class="flame-tip"></div>
          </div>
        {:else}
          <div class="smoke-wisp"></div>
          <div class="burnt-wick"></div>
        {/if}
        <div class="wick"></div>
        <div class="wax-body">
          <div class="wax-drips"></div>
        </div>
      </div>
    {/each}
    <div class="candle-counter-badge">
      <span class="candle-digits-counter">{candles}/{maxCandles}</span>
    </div>
  </div>
</div>

<style>
  .candle-bar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    z-index: 10;
  }

  .candle-counter-badge {
    display: flex;
    align-items: center;
    margin-left: 12px;
    padding: 4px 14px;
    background: radial-gradient(circle at 50% 30%, #26190f 0%, #120b06 100%);
    border: 2px solid #d49344;
    border-radius: 14px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.8), inset 0 1px 1px rgba(255, 230, 160, 0.3);
  }

  .candle-digits-counter {
    font-size: 1.45rem;
    font-weight: 900;
    color: #ffd700;
    letter-spacing: 2px;
    font-family: var(--font-digits, 'Cinzel', serif);
    font-variant-numeric: slashed-zero tabular-nums;
    text-shadow: 0 0 10px rgba(255, 215, 0, 0.7);
  }

  .candle-row {
    display: flex;
    gap: 16px;
    align-items: flex-end;
  }

  .candle-stub {
    position: relative;
    width: 28px;
    height: 52px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
  }

  .wax-body {
    width: 22px;
    height: 32px;
    background: linear-gradient(180deg, #d8cca8 0%, #aa9c7a 70%, #766a50 100%);
    border-radius: 3px 3px 2px 2px;
    box-shadow: inset 2px 0 3px rgba(255, 255, 255, 0.3), inset -2px 0 4px rgba(0, 0, 0, 0.4), 0 3px 8px rgba(0, 0, 0, 0.8);
    position: relative;
    border: 1px solid #4a3e2e;
  }

  .wax-drips {
    position: absolute;
    top: 2px;
    left: 4px;
    width: 3px;
    height: 10px;
    background: #e8ddb5;
    border-radius: 2px;
  }

  .wick {
    width: 2px;
    height: 7px;
    background: #2a2015;
    border-radius: 1px;
    margin-bottom: -1px;
  }

  .flame-wrapper {
    position: absolute;
    top: 0;
    width: 14px;
    height: 20px;
    display: flex;
    justify-content: center;
    align-items: flex-end;
    animation: flameSway 1.8s infinite ease-in-out;
    transform-origin: center bottom;
  }

  .flame-halo {
    position: absolute;
    width: 26px;
    height: 26px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(255, 170, 40, 0.35) 0%, transparent 70%);
    animation: haloPulse 1.2s infinite alternate ease-in-out;
  }

  .flame-core {
    position: absolute;
    width: 10px;
    height: 16px;
    border-radius: 50% 50% 35% 35%;
    background: radial-gradient(ellipse at 50% 80%, #ffffff 0%, #ffdf6d 40%, #ff8a00 80%, #d83d00 100%);
    box-shadow: 0 0 10px #ff9a1f, 0 0 20px #ff5a00;
  }

  .flame-tip {
    position: absolute;
    top: -2px;
    width: 4px;
    height: 6px;
    border-radius: 50%;
    background: #fff8d6;
  }

  .candle-stub.extinguished .wax-body {
    background: linear-gradient(180deg, #5c564c 0%, #3a362f 100%);
    filter: brightness(0.6);
  }

  .burnt-wick {
    position: absolute;
    top: 9px;
    width: 3px;
    height: 4px;
    background: #110d0a;
    border-radius: 2px;
    box-shadow: 0 0 2px #ff3b00;
  }

  .smoke-wisp {
    position: absolute;
    top: -12px;
    width: 8px;
    height: 24px;
    background: radial-gradient(ellipse, rgba(160, 160, 160, 0.5) 0%, transparent 80%);
    border-radius: 50%;
    animation: smokeDrift 2.4s infinite ease-out;
  }

  @keyframes flameSway {
    0%, 100% { transform: rotate(-3deg) scaleY(1); }
    50% { transform: rotate(4deg) scaleY(1.08) scaleX(0.94); }
  }

  @keyframes haloPulse {
    0% { transform: scale(0.9); opacity: 0.8; }
    100% { transform: scale(1.15); opacity: 1; }
  }

  @keyframes smokeDrift {
    0% { transform: translateY(8px) scale(0.4); opacity: 0.8; }
    50% { transform: translateY(-12px) translateX(4px) scale(1); opacity: 0.4; }
    100% { transform: translateY(-30px) translateX(-6px) scale(1.8); opacity: 0; }
  }
</style>
