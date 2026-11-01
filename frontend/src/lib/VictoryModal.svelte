<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function restart() {
    dispatch('restart');
  }
</script>

<div class="victory-backdrop">
  <div class="victory-modal">
    <div class="triumph-icon">⚡</div>
    <h2 class="title">Корчма Обращена в Прах!</h2>
    <p class="subtitle">Ты покорил все три темных сказания</p>
    <div class="lore-ending">
      <p>
        Фонарь Скорбного Следопыта разбит вдребезги. Стены безумия рассыпались бумажным прахом. 
        А пламя Архимага Бездны поглотило проклятую корчму, освободив твою душу из вечной петли.
      </p>
    </div>

    <button class="dread-btn triumph-btn" on:click={restart}>
      Войти в легенды корчмы (Начать Заново)
    </button>
  </div>
</div>

<style>
  .victory-backdrop {
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 50% 50%, rgba(255, 230, 180, 0.95) 0%, rgba(20, 15, 10, 0.96) 80%);
    backdrop-filter: blur(10px);
    z-index: 95;
    display: flex;
    align-items: center;
    justify-content: center;
    animation: whiteFlashIn 1.5s cubic-bezier(0.1, 0.9, 0.2, 1);
  }

  .victory-modal {
    max-width: 640px;
    background: radial-gradient(circle at 50% 20%, #291d12 0%, #150f09 80%, #0a0704 100%);
    border: 2px solid #deb36f;
    border-radius: 8px;
    padding: 40px 32px;
    text-align: center;
    box-shadow: 0 0 70px rgba(255, 180, 50, 0.5), inset 0 0 30px rgba(0, 0, 0, 0.9);
  }

  .triumph-icon {
    font-size: 3.5rem;
    color: #ffd899;
    filter: drop-shadow(0 0 16px #ffae19);
    margin-bottom: 8px;
  }

  .title {
    font-family: var(--font-title);
    font-size: 3rem;
    color: #fff2d1;
    letter-spacing: 2px;
    margin-bottom: 6px;
    text-shadow: 0 0 20px rgba(255, 180, 50, 0.8);
  }

  .subtitle {
    font-size: 1.1rem;
    color: #e5bd87;
    margin-bottom: 20px;
    letter-spacing: 2px;
    text-transform: uppercase;
  }

  .lore-ending {
    background: rgba(10, 8, 5, 0.7);
    border: 1px solid #483520;
    border-radius: 6px;
    padding: 18px 24px;
    margin-bottom: 28px;
  }

  .lore-ending p {
    color: #dfcdba;
    line-height: 1.5;
    font-size: 0.95rem;
  }

  .triumph-btn {
    font-size: 1.4rem;
    padding: 14px 38px;
    border-color: #ffd899;
    background: linear-gradient(180deg, #573a1e 0%, #29190c 100%);
    box-shadow: 0 0 25px rgba(255, 180, 50, 0.4);
  }

  @keyframes whiteFlashIn {
    0% { opacity: 0; background-color: #ffffff; }
    30% { opacity: 1; background-color: #ffffff; }
    100% { opacity: 1; background-color: transparent; }
  }
</style>
