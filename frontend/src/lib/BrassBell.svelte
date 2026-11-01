<script>
  import { audioEngine } from './audioEngine.js';

  let isMuted = false;
  let isRinging = false;

  function handleToggle() {
    isMuted = audioEngine.toggleMute();
    isRinging = true;
    setTimeout(() => {
      isRinging = false;
    }, 400);
  }
</script>

<button 
  class="brass-bell-btn" 
  class:ringing={isRinging} 
  class:muted={isMuted}
  on:click={handleToggle}
  title={isMuted ? "Включить звук корчмы (Медный колокольчик)" : "Заглушить звук корчмы (Медный колокольчик)"}
>
  <div class="bell-icon">🔔</div>
  <div class="bell-state">{isMuted ? "Безмолвие" : "Звук"}</div>
</button>

<style>
  .brass-bell-btn {
    position: absolute;
    top: 14px;
    right: 18px;
    background: radial-gradient(circle at 40% 30%, #d4a359 0%, #855c27 70%, #4a3212 100%);
    border: 2px solid #deb36f;
    border-radius: 8px;
    padding: 6px 12px;
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.8), inset 0 1px 2px rgba(255, 255, 255, 0.4);
    z-index: 100;
    transition: transform 0.15s ease, filter 0.2s ease;
  }

  .brass-bell-btn:hover {
    transform: scale(1.06);
    filter: brightness(1.15);
  }

  .brass-bell-btn.muted {
    background: radial-gradient(circle at 40% 30%, #7d6b53 0%, #4d402f 70%, #292218 100%);
    border-color: #63533d;
    filter: grayscale(0.5);
  }

  .bell-icon {
    font-size: 1.2rem;
    filter: drop-shadow(0 1px 2px rgba(0,0,0,0.6));
  }

  .bell-state {
    font-family: var(--font-title);
    font-size: 0.95rem;
    color: #fff4dc;
    letter-spacing: 1px;
    text-transform: uppercase;
  }

  .brass-bell-btn.ringing {
    animation: bellRing 0.4s ease-in-out;
  }

  @keyframes bellRing {
    0%, 100% { transform: rotate(0deg); }
    25% { transform: rotate(-18deg) scale(1.1); }
    50% { transform: rotate(18deg) scale(1.1); }
    75% { transform: rotate(-8deg); }
  }
</style>
