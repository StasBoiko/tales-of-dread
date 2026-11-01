<script>
  import { onMount, onDestroy } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  let volume = 0.28;
  let isMuted = false;
  let isPaused = false;
  let isRinging = false;
  let unsubscribe = null;

  onMount(() => {
    volume = audioEngine.getVolume();
    isMuted = audioEngine.isMuted;
    isPaused = audioEngine.isPaused;
    unsubscribe = audioEngine.subscribe((state) => {
      isMuted = state.isMuted;
      isPaused = state.isPaused;
      volume = state.volume;
    });
  });

  onDestroy(() => {
    if (unsubscribe) unsubscribe();
  });

  function handleVolumeInput(e) {
    const val = parseFloat(e.target.value);
    volume = val;
    audioEngine.setVolume(val);
  }

  async function handleMuteToggle() {
    isMuted = await audioEngine.toggleMute();
    isRinging = true;
    setTimeout(() => {
      isRinging = false;
    }, 400);
  }

  async function handlePlayPauseToggle() {
    const isPlaying = await audioEngine.togglePlayPause();
    isPaused = !isPlaying;
  }
</script>

<div class="tavern-acoustics-panel">
  <!-- Play / Pause Button -->
  <button 
    class="acoustics-icon-btn play-pause-btn" 
    class:paused={isPaused}
    on:click={handlePlayPauseToggle}
    title={isPaused ? "Возобновить музыку корчмы (Play)" : "Приостановить музыку (Pause)"}
    aria-label="Пауза / Воспроизведение"
  >
    <span class="btn-glyph">{isPaused ? '▶' : '⏸'}</span>
  </button>

  <!-- Mute / Unmute Bell Button -->
  <button 
    class="acoustics-icon-btn mute-bell-btn" 
    class:ringing={isRinging} 
    class:muted={isMuted}
    on:click={handleMuteToggle}
    title={isMuted ? "Включить звук корчмы" : "Заглушить все звуки (Медный колокольчик)"}
    aria-label="Переключить безмолвие"
  >
    <span class="btn-glyph bell-glyph">{isMuted ? '🔕' : '🔔'}</span>
  </button>

  <div class="slider-wrapper" title="Громкость корчмы: {Math.round(volume * 100)}%">
    <span class="chain-rivet left">⛓</span>
    <input 
      type="range" 
      min="0" 
      max="1" 
      step="0.01" 
      value={volume} 
      on:input={handleVolumeInput}
      class="rust-chain-slider"
      class:slider-muted={isMuted}
      aria-label="Регулятор громкости"
    />
    <span class="chain-rivet right">⛓</span>
  </div>

  <span class="volume-percent">{isMuted ? '0%' : `${Math.round(volume * 100)}%`}</span>
</div>

<style>
  .tavern-acoustics-panel {
    position: relative;
    display: flex;
    align-items: center;
    gap: 10px;
    background: linear-gradient(180deg, #2a1b12 0%, #150d09 100%);
    border: 2px solid #5d3d25;
    border-radius: 20px;
    padding: 4px 14px 4px 8px;
    box-shadow: 
      0 4px 14px rgba(0, 0, 0, 0.9),
      inset 0 1px 1px rgba(255, 220, 160, 0.2),
      inset 0 -1px 3px rgba(0, 0, 0, 0.8);
    z-index: 100;
    user-select: none;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .tavern-acoustics-panel:hover {
    border-color: #9e693f;
    box-shadow: 
      0 6px 18px rgba(0, 0, 0, 0.95),
      0 0 14px rgba(230, 140, 50, 0.25);
  }

  .acoustics-icon-btn {
    border-radius: 50%;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.8), inset 0 1px 1px rgba(255, 255, 255, 0.4);
    transition: transform 0.15s ease, filter 0.2s ease, border-color 0.2s ease;
    user-select: none;
  }

  .play-pause-btn {
    background: radial-gradient(circle at 40% 30%, #466075 0%, #223442 70%, #121c24 100%);
    border: 1px solid #7ea3c2;
    color: #e2effa;
    font-size: 0.85rem;
  }

  .play-pause-btn:hover {
    transform: scale(1.1);
    filter: brightness(1.2);
    border-color: #a8cded;
  }

  .play-pause-btn.paused {
    background: radial-gradient(circle at 40% 30%, #52473c 0%, #302821 70%, #1a1511 100%);
    border-color: #635445;
    color: #a19280;
  }

  .mute-bell-btn {
    background: radial-gradient(circle at 40% 30%, #deb36f 0%, #8c5d26 70%, #4a3012 100%);
    border: 1px solid #ecc987;
  }

  .mute-bell-btn:hover {
    transform: scale(1.1);
    filter: brightness(1.2);
  }

  .mute-bell-btn.muted {
    background: radial-gradient(circle at 40% 30%, #52473c 0%, #302821 70%, #1a1511 100%);
    border-color: #635445;
    filter: grayscale(0.7);
  }

  .mute-bell-btn.ringing {
    animation: bellRing 0.4s ease-in-out;
  }

  .bell-glyph {
    font-size: 1.05rem;
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.8));
  }

  .slider-wrapper {
    display: flex;
    align-items: center;
    gap: 4px;
    position: relative;
  }

  .chain-rivet {
    font-size: 0.8rem;
    color: #6a492f;
    opacity: 0.85;
    text-shadow: 0 1px 1px rgba(0, 0, 0, 0.9);
  }

  /* Stylized Tavern Range Slider */
  .rust-chain-slider {
    -webkit-appearance: none;
    appearance: none;
    width: 95px;
    height: 10px;
    background: #0d0806;
    border-radius: 5px;
    outline: none;
    cursor: pointer;
    margin: 0;
    position: relative;
    border: 1px solid #4a2f1b;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.9), 0 1px 0 rgba(255, 255, 255, 0.1);
    background-image: repeating-linear-gradient(
      90deg,
      #1f130b 0px,
      #1f130b 6px,
      #332014 7px,
      #1f130b 8px
    );
  }

  .rust-chain-slider.slider-muted {
    opacity: 0.45;
    filter: grayscale(0.8);
  }

  /* WebKit Thumb - Skull / Iron Rivet with Amber Glow */
  .rust-chain-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: radial-gradient(circle at 35% 35%, #ffd480 0%, #b86b24 55%, #421d0a 100%);
    border: 1px solid #ffe6a6;
    cursor: grab;
    box-shadow: 
      0 2px 6px rgba(0, 0, 0, 0.95),
      0 0 8px rgba(255, 140, 30, 0.6),
      inset 0 1px 2px rgba(255, 255, 255, 0.8);
    transition: transform 0.15s ease, box-shadow 0.15s ease;
  }

  .rust-chain-slider::-webkit-slider-thumb:hover {
    transform: scale(1.22);
    box-shadow: 
      0 3px 8px rgba(0, 0, 0, 0.95),
      0 0 14px rgba(255, 170, 50, 0.9);
  }

  .rust-chain-slider::-webkit-slider-thumb:active {
    cursor: grabbing;
    transform: scale(1.1);
  }

  /* Firefox Thumb */
  .rust-chain-slider::-moz-range-thumb {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: radial-gradient(circle at 35% 35%, #ffd480 0%, #b86b24 55%, #421d0a 100%);
    border: 1px solid #ffe6a6;
    cursor: grab;
    box-shadow: 
      0 2px 6px rgba(0, 0, 0, 0.95),
      0 0 8px rgba(255, 140, 30, 0.6);
    transition: transform 0.15s ease;
  }

  .rust-chain-slider::-moz-range-thumb:hover {
    transform: scale(1.22);
  }

  .volume-percent {
    font-family: var(--font-title);
    font-size: 0.95rem;
    color: #dfcaa8;
    min-width: 32px;
    text-align: right;
    letter-spacing: 0.5px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.9);
  }

  @keyframes bellRing {
    0%, 100% { transform: rotate(0deg); }
    25% { transform: rotate(-18deg) scale(1.1); }
    50% { transform: rotate(18deg) scale(1.1); }
    75% { transform: rotate(-8deg); }
  }
</style>
