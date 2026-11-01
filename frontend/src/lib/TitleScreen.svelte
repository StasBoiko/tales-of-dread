<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function start() {
    dispatch('start');
  }
</script>

<div class="title-screen-container">
  <div class="title-glow-aura"></div>
  <div class="title-content">
    <div class="sub-header">Пошаговый карточный рогалик-хоррор</div>
    <h1 class="main-title">ЧЕРТОВО ОЧКО</h1>
    <div class="title-subtitle">СТРАШНЫЕ СКАЗКИ КОРЧМЫ</div>

    <div class="lore-box">
      <p>
        За порогом бушует ледяная метель, а в полумраке корчмы на скрипучем столе 
        тебя ждут три жутких противника. Держись ближе к 21 очку, остерегайся проклятий 
        и не дай пяти огаркам свечей угаснуть в холодной тьме.
      </p>
    </div>

    <div class="bosses-preview">
      <div class="preview-card">
        <span class="preview-num">Сказка I</span>
        <h4>Скорбный Фонарщик</h4>
        <p class="track-tag">🎵 The Lantern's Curse</p>
      </div>
      <div class="preview-card">
        <span class="preview-num">Сказка II</span>
        <h4>Узник Бумажных Стен</h4>
        <p class="track-tag">🎵 Kingdom of the Paper Walls</p>
      </div>
      <div class="preview-card">
        <span class="preview-num">Сказка III</span>
        <h4>Архимаг Бездны</h4>
        <p class="track-tag">🎵 Let It Go to Hell</p>
      </div>
    </div>

    <button class="dread-btn step-button" on:click={start}>
      Шагнуть во тьму корчмы
    </button>
  </div>
</div>

<style>
  .title-screen-container {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 90;
    background: radial-gradient(circle at 50% 50%, #150e09 0%, #0a0604 60%, #030202 100%);
    overflow: hidden;
  }

  .title-glow-aura {
    position: absolute;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(255, 140, 30, 0.15) 0%, transparent 70%);
    border-radius: 50%;
    animation: pulseAura 4s infinite alternate ease-in-out;
  }

  .title-content {
    position: relative;
    z-index: 2;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    max-width: 760px;
    padding: 20px;
  }

  .sub-header {
    font-size: 0.85rem;
    color: #a48c73;
    letter-spacing: 3px;
    text-transform: uppercase;
    margin-bottom: 8px;
  }

  .main-title {
    font-family: var(--font-title);
    font-size: 4rem;
    color: #ffd899;
    letter-spacing: 4px;
    line-height: 1;
    text-shadow: 
      0 0 10px rgba(255, 150, 40, 0.8),
      0 4px 15px rgba(0, 0, 0, 0.95);
  }

  .title-subtitle {
    font-family: var(--font-title);
    font-size: 1.8rem;
    color: #c94028;
    letter-spacing: 6px;
    margin-bottom: 18px;
    text-shadow: 0 0 8px rgba(200, 50, 30, 0.6);
  }

  .lore-box {
    background: rgba(14, 10, 8, 0.7);
    border: 1px solid #483222;
    border-radius: 6px;
    padding: 16px 24px;
    margin-bottom: 24px;
    box-shadow: inset 0 0 20px rgba(0, 0, 0, 0.8);
  }

  .lore-box p {
    color: #d1beaa;
    line-height: 1.5;
    font-size: 0.95rem;
  }

  .bosses-preview {
    display: flex;
    gap: 16px;
    margin-bottom: 30px;
    width: 100%;
  }

  .preview-card {
    flex: 1;
    background: linear-gradient(180deg, #1f140e 0%, #110a07 100%);
    border: 1px solid #3d2719;
    border-radius: 4px;
    padding: 12px 8px;
  }

  .preview-num {
    font-size: 0.7rem;
    color: #8c7661;
    text-transform: uppercase;
    letter-spacing: 1px;
    display: block;
  }

  .preview-card h4 {
    font-family: var(--font-title);
    font-size: 1.05rem;
    color: #edd6bb;
    margin: 4px 0;
  }

  .track-tag {
    font-size: 0.75rem;
    color: #d68945;
  }

  .step-button {
    font-size: 1.6rem;
    padding: 14px 44px;
    border-radius: 6px;
    animation: glowBtn 2s infinite alternate ease-in-out;
  }

  @keyframes pulseAura {
    0% { transform: scale(0.9); opacity: 0.5; }
    100% { transform: scale(1.15); opacity: 0.8; }
  }

  @keyframes glowBtn {
    0% { box-shadow: 0 0 10px rgba(255, 150, 40, 0.3); }
    100% { box-shadow: 0 0 28px rgba(255, 170, 50, 0.7); }
  }
</style>
