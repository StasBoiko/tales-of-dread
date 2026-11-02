<script>
  import { langStore } from '../i18n/index.js';
  import { setLang, getLang } from '../i18n/index.js';

  let currentLang = getLang();

  $: currentLang = getLang();

  function switchLang() {
    const newLang = currentLang === 'ru' ? 'en' : 'ru';
    setLang(newLang);
    langStore.set(newLang);
  }
</script>

<div class="language-switch">
  <button
    class="lang-btn"
    class:active={currentLang === 'ru'}
    onclick={() => { setLang('ru'); langStore.set('ru'); }}
    title="Русский"
  >
    RU
  </button>
  <button
    class="lang-btn"
    class:active={currentLang === 'en'}
    onclick={() => { setLang('en'); langStore.set('en'); }}
    title="English"
  >
    EN
  </button>
</div>

<style>
  .language-switch {
    display: flex;
    gap: 6px;
    position: fixed;
    top: 16px;
    right: 16px;
    z-index: 1000;
    background: rgba(0, 0, 0, 0.6);
    padding: 4px;
    border-radius: 8px;
    border: 1px solid #5a3c26;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.8);
  }

  .lang-btn {
    background: #1a1410;
    color: #8c6f52;
    border: 1px solid #4a3320;
    padding: 6px 12px;
    font-family: var(--font-title);
    font-size: 0.85rem;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .lang-btn:hover {
    color: #ffd700;
    border-color: #a8793b;
  }

  .lang-btn.active {
    color: #ffd700;
    background: #2a1a12;
    border-color: #ffd700;
    box-shadow: 0 0 12px rgba(255, 215, 0, 0.3);
  }
</style>
