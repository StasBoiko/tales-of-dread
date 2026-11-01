<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function restart() {
    dispatch('restart');
  }
</script>

<div class="game-over-backdrop">
  <div class="game-over-modal">
    <div class="skull-icon">☠</div>
    <h2 class="title">Свечи Погасли</h2>
    <p class="subtitle">Тьма корчмы поглотила твой рассудок</p>
    <p class="desc">
      Холодный ветер разметал карты по столу. Никто не выходит из этой корчмы прежним.
    </p>

    <button class="dread-btn dread-btn-danger restart-btn" on:click={restart}>
      Зажечь свечи вновь
    </button>
  </div>
</div>

<style>
  .game-over-backdrop {
    position: absolute;
    inset: 0;
    background: rgba(4, 2, 2, 0.92);
    backdrop-filter: blur(8px);
    z-index: 85;
    display: flex;
    align-items: center;
    justify-content: center;
    animation: fadeIn 0.5s ease-out;
  }

  .game-over-modal {
    max-width: 520px;
    background: radial-gradient(circle at 50% 20%, #1a0a0a 0%, #0d0404 80%, #050202 100%);
    border: 2px solid #5a1919;
    border-radius: 8px;
    padding: 36px 28px;
    text-align: center;
    box-shadow: 0 0 50px rgba(180, 20, 20, 0.4);
  }

  .skull-icon {
    font-size: 3.5rem;
    color: #e04545;
    margin-bottom: 8px;
    filter: drop-shadow(0 0 10px rgba(220, 40, 40, 0.7));
  }

  .title {
    font-family: var(--font-title);
    font-size: 2.8rem;
    color: #ff6666;
    letter-spacing: 2px;
    margin-bottom: 6px;
  }

  .subtitle {
    font-family: var(--font-body);
    font-size: 1rem;
    color: #d99f9f;
    margin-bottom: 16px;
    letter-spacing: 1px;
  }

  .desc {
    color: #a88585;
    font-size: 0.9rem;
    line-height: 1.4;
    margin-bottom: 28px;
  }

  .restart-btn {
    font-size: 1.3rem;
    padding: 12px 32px;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.94); }
    to { opacity: 1; transform: scale(1); }
  }
</style>
