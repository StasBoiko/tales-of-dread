<script>
  import { createEventDispatcher } from 'svelte';
  import { audioEngine } from './audioEngine.js';

  export let offeredArtifacts = [];
  export let bossName = '';

  const dispatch = createEventDispatcher();

  function select(artifact) {
    audioEngine.playCardThud();
    dispatch('select', { artifactId: artifact.id });
  }
</script>

<div class="artifact-backdrop">
  <div class="artifact-modal">
    <div class="modal-header">
      <span class="sub-title">Награда за победу</span>
      <h2 class="title">{bossName} повержен</h2>
      <p class="description">
        ВЫБЕРИТЕ ОДИН ИЗ ТРЕХ ПРОКЛЯТЫХ АРТЕФАКТОВ ДЛЯ ПРОДОЛЖЕНИЯ ПУТИ:
      </p>
    </div>

    <div class="artifacts-grid">
      {#each offeredArtifacts as art}
        <button class="artifact-card" on:click={() => select(art)}>
          <div class="art-icon">{art.icon}</div>
          <h3 class="art-name">{art.name}</h3>
          <p class="art-desc">{art.description}</p>
          <div class="claim-btn">Взять реликвию</div>
        </button>
      {/each}
    </div>
  </div>
</div>

<style>
  .artifact-backdrop {
    position: absolute;
    inset: 0;
    background: rgba(3, 2, 2, 0.88);
    backdrop-filter: blur(6px);
    z-index: 80;
    display: flex;
    align-items: center;
    justify-content: center;
    animation: fadeIn 0.4s ease-out;
  }

  .artifact-modal {
    width: 90%;
    max-width: 820px;
    background: radial-gradient(circle at 50% 20%, #1e140d 0%, #100b08 80%, #070403 100%);
    border: 2px solid #5a3e26;
    border-radius: 8px;
    padding: 30px;
    box-shadow: 0 0 50px rgba(0, 0, 0, 0.95), 0 0 25px rgba(255, 140, 30, 0.15);
    text-align: center;
  }

  .modal-header .sub-title {
    font-size: 0.95rem;
    color: #e5be68;
    letter-spacing: 2px;
    text-transform: uppercase;
    font-weight: 700;
  }

  .modal-header .title {
    font-family: var(--font-title);
    font-size: 2.4rem;
    color: #ffd700;
    margin: 4px 0 10px;
    letter-spacing: 1.5px;
    text-shadow: 0 0 12px rgba(255, 215, 0, 0.4);
  }

  .modal-header .description {
    color: #ffd700;
    font-size: 0.92rem;
    font-weight: 700;
    letter-spacing: 0.8px;
    margin-bottom: 24px;
    white-space: nowrap;
    text-align: center;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.9);
  }

  .artifacts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 20px;
  }

  .artifact-card {
    background: linear-gradient(180deg, #241710 0%, #150d09 100%);
    border: 1px solid #6b4c2a;
    border-radius: 6px;
    padding: 22px 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    cursor: pointer;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.8);
    transition: transform 0.25s cubic-bezier(0.2, 0.8, 0.3, 1), border-color 0.25s ease, box-shadow 0.25s ease;
  }

  .artifact-card:hover {
    transform: translateY(-6px);
    border-color: #ffd700;
    box-shadow: 0 12px 28px rgba(0, 0, 0, 0.95), 0 0 20px rgba(255, 215, 0, 0.35);
  }

  .art-icon {
    font-size: 3rem;
    margin-bottom: 12px;
    filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.8));
  }

  .art-name {
    font-family: var(--font-title);
    font-size: 1.4rem;
    color: #ffd700;
    margin-bottom: 8px;
    letter-spacing: 1px;
  }

  .art-desc {
    font-family: var(--font-body);
    font-size: 0.92rem;
    color: #e5be68;
    line-height: 1.4;
    margin-bottom: 16px;
    flex-grow: 1;
  }

  .claim-btn {
    font-family: var(--font-title);
    font-size: 1.05rem;
    color: #ffd700;
    border: 1px solid #996e28;
    border-radius: 4px;
    padding: 6px 16px;
    background: rgba(0, 0, 0, 0.5);
    text-transform: uppercase;
    letter-spacing: 1px;
    transition: all 0.2s ease;
  }

  .artifact-card:hover .claim-btn {
    background: #503118;
    border-color: #ffd700;
    color: #fff;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.96); }
    to { opacity: 1; transform: scale(1); }
  }
</style>
