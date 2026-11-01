<script>
  // InventoryScroll: Right column ancient relic vault symmetrical to Rules
  export let artifacts = [];

  $: medallion = artifacts && artifacts[0];
  $: trophy1 = artifacts && artifacts[1];
  $: trophy2 = artifacts && artifacts[2];
</script>

<aside class="scroll-column inventory-scroll-column" aria-label="Полка реликвий">
  <div class="scroll-parchment">
    <!-- Top Iron Mounting Rod -->
    <div class="iron-rod rod-top">
      <div class="rod-finial left"></div>
      <div class="rod-shaft"></div>
      <div class="rod-finial right"></div>
    </div>

    <div class="scroll-content">
      <div class="wax-seal iron-seal" title="Печать Тайника">
        <span class="seal-icon">🗝️</span>
      </div>

      <header class="scroll-header">
        <h2 class="scroll-title">РЕЛИКВИИ</h2>
      </header>

      <div class="scroll-divider"></div>

      <!-- 3 Relic Cards: Slot 1 is always the Story Medallion; Slots 2 & 3 are Boss Trophies -->
      <div class="relics-list-container">
        <!-- 1. The Story Medallion (Silver Charm / Molten Medallion) -->
        {#if medallion}
          <div 
            class="relic-item-panel filled medallion-panel" 
            class:silver-active={medallion.id === 'silver_medallion' && medallion.active !== false}
            class:silver-spent={medallion.id === 'silver_medallion' && medallion.active === false}
            class:molten-active={medallion.id === 'molten_medallion' && medallion.active !== false}
            class:molten-spent={medallion.id === 'molten_medallion' && medallion.active === false}
          >
            <div class="relic-panel-top">
              <div class="medallion-icon-wrap">
                {#if medallion.id === 'molten_medallion'}
                  <svg viewBox="0 0 40 40" class="medallion-svg molten-svg" aria-hidden="true">
                    <circle cx="20" cy="20" r="17" fill="#140608" stroke="#ff2a4b" stroke-width="2" />
                    <circle cx="20" cy="20" r="12" fill="#2d080e" stroke="#ff6b3d" stroke-width="1.5" />
                    <!-- Burning crimson eye rune -->
                    <path d="M 12 20 Q 20 12 28 20 Q 20 28 12 20 Z" fill="#ff1a3c" />
                    <circle cx="20" cy="20" r="3.5" fill="#fff" />
                  </svg>
                {:else}
                  <svg viewBox="0 0 40 40" class="medallion-svg silver-svg" aria-hidden="true">
                    <circle cx="20" cy="20" r="17" fill="#18232c" stroke="#a0c8e6" stroke-width="2" />
                    <circle cx="20" cy="20" r="12" fill="#243747" stroke="#d5e8f5" stroke-width="1.5" />
                    <!-- Luminous silver cross/star rune -->
                    <polygon points="20,11 22.5,17.5 29,20 22.5,22.5 20,29 17.5,22.5 11,20 17.5,17.5" fill="#78d5ff" />
                    <circle cx="20" cy="20" r="2.5" fill="#ffffff" />
                  </svg>
                {/if}
              </div>
              <div class="relic-title-group">
                <span class="relic-slot-badge medallion-badge">
                  {#if medallion.id === 'molten_medallion'}
                    {medallion.active !== false ? 'ОБРЯДОВАЯ РЕЛИКВИЯ' : 'ИСПОЛЬЗОВАНО'}
                  {:else}
                    {medallion.active !== false ? 'СЮЖЕТНЫЙ ОБЕРЕГ' : 'СИЛА ИССЯКЛА'}
                  {/if}
                </span>
                <h3 class="relic-item-name">{medallion.name}</h3>
              </div>
            </div>
            <div class="relic-panel-divider"></div>
            <p class="relic-item-desc">{medallion.description}</p>
            <div class="relic-panel-tag">
              {#if medallion.id === 'molten_medallion'}
                {medallion.active !== false ? '«Проклятый переброс (1 раз)»' : '«Переброс совершен»'}
              {:else}
                {medallion.active !== false ? '«Защита от гибели (1 раз)»' : '«Спасение совершено»'}
              {/if}
            </div>
          </div>
        {/if}

        <!-- 2. Boss 1 Trophy (Lantern Bearer) -->
        {#if trophy1}
          <div class="relic-item-panel filled">
            <div class="relic-panel-top">
              <span class="relic-big-icon">{trophy1.icon}</span>
              <div class="relic-title-group">
                <span class="relic-slot-badge">ТРОФЕЙ ФОНАРЩИКА</span>
                <h3 class="relic-item-name">{trophy1.name}</h3>
              </div>
            </div>
            <div class="relic-panel-divider"></div>
            <p class="relic-item-desc">{trophy1.description}</p>
            <div class="relic-panel-tag">«Трофей побежденного дилера»</div>
          </div>
        {:else}
          <div class="relic-item-panel empty">
            <div class="relic-panel-top">
              <div class="empty-icon-box">
                <span class="empty-glyph">🗝️</span>
              </div>
              <div class="relic-title-group">
                <span class="relic-slot-badge empty-badge">СЛОТ 2</span>
                <h3 class="relic-item-name empty-name">ПУСТАЯ НИША</h3>
              </div>
            </div>
            <div class="relic-panel-divider empty-div"></div>
            <p class="relic-item-desc empty-text">Откроется после победы над Фонарщиком.</p>
            <div class="relic-panel-tag empty-tag">Ожидает победы над боссом</div>
          </div>
        {/if}

        <!-- 3. Boss 2 Trophy (Prisoner) -->
        {#if trophy2}
          <div class="relic-item-panel filled">
            <div class="relic-panel-top">
              <span class="relic-big-icon">{trophy2.icon}</span>
              <div class="relic-title-group">
                <span class="relic-slot-badge">ТРОФЕЙ УЗНИКА</span>
                <h3 class="relic-item-name">{trophy2.name}</h3>
              </div>
            </div>
            <div class="relic-panel-divider"></div>
            <p class="relic-item-desc">{trophy2.description}</p>
            <div class="relic-panel-tag">«Трофей побежденного дилера»</div>
          </div>
        {:else}
          <div class="relic-item-panel empty">
            <div class="relic-panel-top">
              <div class="empty-icon-box">
                <span class="empty-glyph">🗝️</span>
              </div>
              <div class="relic-title-group">
                <span class="relic-slot-badge empty-badge">СЛОТ 3</span>
                <h3 class="relic-item-name empty-name">ПУСТАЯ НИША</h3>
              </div>
            </div>
            <div class="relic-panel-divider empty-div"></div>
            <p class="relic-item-desc empty-text">Откроется после победы над Узником.</p>
            <div class="relic-panel-tag empty-tag">Ожидает победы над боссом</div>
          </div>
        {/if}
      </div>

      <!-- Atmospheric Lore Card: Updated for 3 slots -->
      <div class="relic-lore-card">
        <p class="lore-paragraph">
          {#if !artifacts || artifacts.length <= 1}
            «Медальон Алёны хранит тепло. Впереди бой за реликвии побежденных стражей корчмы.»
          {:else if artifacts.length === 2}
            «Первый трофей покоится в твоих руках. Впереди бой за вторую проклятую реликвию.»
          {:else}
            «Все реликвии собраны. Твой дух готов к встрече с истинным кошмаром корчмы.»
          {/if}
        </p>
      </div>

      <footer class="scroll-footer">
        <span class="whisper-hint">«Каждая реликвия оплачена кровью и воском.»</span>
      </footer>
    </div>

    <!-- Bottom Iron Mounting Rod -->
    <div class="iron-rod rod-bottom">
      <div class="rod-finial left"></div>
      <div class="rod-shaft"></div>
      <div class="rod-finial right"></div>
    </div>
  </div>
</aside>

<style>
  .inventory-scroll-column {
    width: 280px;
    min-width: 280px;
    max-width: 280px;
    height: 100%;
    min-height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    box-sizing: border-box;
    font-variant-numeric: slashed-zero tabular-nums;
  }

  .scroll-parchment {
    position: relative;
    background: radial-gradient(circle at 50% 20%, #1f140e 0%, #140d09 70%, #0a0604 100%);
    border-left: 2px solid #5a3c26;
    border-right: 2px solid #5a3c26;
    border-radius: 4px;
    box-shadow: 
      0 10px 30px rgba(0, 0, 0, 0.95),
      inset 0 0 40px rgba(0, 0, 0, 0.9);
    display: flex;
    flex-direction: column;
    height: 94%;
    max-height: 820px;
  }

  .iron-rod {
    height: 14px;
    display: flex;
    align-items: center;
    position: relative;
    z-index: 5;
  }

  .rod-shaft {
    flex: 1;
    height: 8px;
    background: linear-gradient(180deg, #5a4b41 0%, #2b231d 50%, #15110e 100%);
    border-top: 1px solid #786457;
    border-bottom: 1px solid #000;
  }

  .rod-finial {
    width: 16px;
    height: 14px;
    background: radial-gradient(circle at 40% 40%, #786457 0%, #3a3028 60%, #15110e 100%);
    border-radius: 3px;
    border: 1px solid #1a1410;
    box-shadow: 0 2px 4px rgba(0,0,0,0.8);
  }

  .scroll-content {
    flex: 1;
    padding: 18px 16px;
    display: flex;
    flex-direction: column;
    position: relative;
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: #5a3c26 #120c08;
  }

  .scroll-content::-webkit-scrollbar {
    width: 5px;
  }

  .scroll-content::-webkit-scrollbar-thumb {
    background: #5a3c26;
    border-radius: 3px;
  }

  .wax-seal {
    position: absolute;
    top: 10px;
    right: 12px;
    width: 32px;
    height: 32px;
    background: radial-gradient(circle at 35% 35%, #5c472d 0%, #3a2b19 70%, #1c140a 100%);
    border-radius: 50%;
    border: 1px solid #94754d;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 3px 8px rgba(0,0,0,0.85);
  }

  .seal-icon {
    font-size: 0.95rem;
    filter: drop-shadow(0 1px 1px rgba(0,0,0,0.9));
  }

  .scroll-header {
    text-align: left;
    margin-bottom: 8px;
    padding-right: 36px;
  }

  .scroll-title {
    font-family: var(--font-title);
    font-size: 2rem;
    color: #ffd700;
    margin: 0;
    letter-spacing: 2px;
    text-shadow: 0 0 10px rgba(255, 215, 0, 0.4), 0 2px 4px #000;
  }

  .scroll-divider {
    height: 2px;
    background: linear-gradient(90deg, #d4af37 0%, #8a5e38 70%, transparent 100%);
    margin-bottom: 16px;
  }

  /* List of 2 Relic Cards */
  .relics-list-container {
    display: flex;
    flex-direction: column;
    gap: 14px;
    margin-bottom: 16px;
  }

  .relic-item-panel {
    border-radius: 8px;
    padding: 12px 14px;
    display: flex;
    flex-direction: column;
    text-align: left;
    box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.9);
    transition: all 0.25s ease;
  }

  .relic-item-panel.filled {
    background: radial-gradient(circle at 50% 25%, #2a1a12 0%, #170e09 75%, #0e0805 100%);
    border: 1px solid #a8793b;
    box-shadow: 
      inset 0 1px 3px rgba(255, 215, 0, 0.15),
      0 4px 12px rgba(0, 0, 0, 0.85),
      0 0 10px rgba(212, 175, 55, 0.15);
  }

  .relic-item-panel.filled:hover {
    border-color: #ffd700;
    transform: translateY(-2px);
    box-shadow: 
      0 6px 18px rgba(0, 0, 0, 0.95),
      0 0 16px rgba(255, 215, 0, 0.35);
  }

  .relic-item-panel.empty {
    background: radial-gradient(circle at 50% 40%, #16100c 0%, #0d0806 100%);
    border: 1px dashed #4a3424;
    opacity: 0.85;
  }

  .relic-panel-top {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .relic-big-icon {
    font-size: 2rem;
    line-height: 1;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.9));
  }

  .empty-icon-box {
    width: 36px;
    height: 36px;
    border-radius: 6px;
    background: #110b08;
    border: 1px dashed #4a3525;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .empty-glyph {
    font-size: 1.1rem;
    opacity: 0.45;
  }

  .relic-title-group {
    display: flex;
    flex-direction: column;
    gap: 2px;
    overflow: hidden;
  }

  .relic-slot-badge {
    font-family: var(--font-title);
    font-size: 0.72rem;
    color: #e5be68;
    letter-spacing: 1.5px;
    text-transform: uppercase;
    font-weight: 700;
  }

  .empty-badge {
    color: #7a6047;
  }

  .relic-item-name {
    font-family: var(--font-title);
    font-size: 1.25rem;
    color: #ffd700;
    margin: 0;
    letter-spacing: 0.5px;
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.9);
  }

  .empty-name {
    color: #8c6f52;
    font-size: 1.05rem;
  }

  .relic-panel-divider {
    height: 1px;
    background: linear-gradient(90deg, #8a5e38 0%, #4a3320 80%, transparent 100%);
    margin: 8px 0;
  }

  .empty-div {
    background: #2a1c13;
  }

  .relic-item-desc {
    font-size: 0.92rem;
    color: #faecd9;
    line-height: 1.4;
    margin: 0 0 6px 0;
  }

  .empty-text {
    color: #7d654f;
    font-style: italic;
    font-size: 0.88rem;
  }

  .relic-panel-tag {
    font-size: 0.76rem;
    color: #c49d60;
    font-style: italic;
    letter-spacing: 0.3px;
  }

  .empty-tag {
    color: #5c4430;
  }

  /* Lore Banner */
  .relic-lore-card {
    background: rgba(18, 12, 9, 0.85);
    border: 1px solid #5a3c26;
    border-radius: 6px;
    padding: 12px 14px;
    margin-top: auto;
    box-shadow: inset 0 0 12px rgba(0, 0, 0, 0.9);
  }

  .lore-paragraph {
    font-size: 0.88rem;
    color: #c99a3e;
    line-height: 1.4;
    margin: 0;
    font-style: italic;
    text-align: left;
  }

  .scroll-footer {
    margin-top: 12px;
    padding-top: 8px;
    border-top: 1px solid #4a3220;
    text-align: center;
  }

  .whisper-hint {
    font-size: 0.85rem;
    font-style: italic;
    color: #9c805c;
    line-height: 1.25;
    display: block;
    text-shadow: 0 1px 2px rgba(0,0,0,0.8);
  }

  /* Medallion Specific Panels */
  .medallion-icon-wrap {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .medallion-svg {
    width: 36px;
    height: 36px;
    display: block;
    transition: all 0.3s ease;
  }

  .medallion-panel.silver-active {
    background: radial-gradient(circle at 50% 25%, #182b3d 0%, #0d1721 75%, #060c12 100%);
    border: 1.5px solid #78c5ff;
    box-shadow: 
      inset 0 1px 3px rgba(180, 230, 255, 0.3),
      0 4px 14px rgba(0, 0, 0, 0.9),
      0 0 16px rgba(100, 190, 255, 0.45);
    animation: silverMedallionGlow 3.5s ease-in-out infinite alternate;
  }

  .medallion-panel.silver-active .relic-item-name {
    color: #e0f2fe;
    text-shadow: 0 0 8px rgba(120, 200, 255, 0.7);
  }

  .medallion-panel.silver-active .medallion-badge {
    color: #7dd3fc;
  }

  .medallion-panel.silver-spent {
    background: radial-gradient(circle at 50% 25%, #1a1a1a 0%, #111111 75%, #080808 100%);
    border: 1px solid #444444;
    filter: grayscale(0.9) brightness(0.7);
    opacity: 0.75;
  }

  .medallion-panel.silver-spent .relic-item-name {
    color: #888888;
  }

  .medallion-panel.silver-spent .medallion-badge {
    color: #666666;
  }

  .medallion-panel.molten-active {
    background: radial-gradient(circle at 50% 25%, #380b12 0%, #1c0509 75%, #0b0204 100%);
    border: 1.5px solid #ff3b5c;
    box-shadow: 
      inset 0 1px 3px rgba(255, 80, 100, 0.35),
      0 4px 14px rgba(0, 0, 0, 0.9),
      0 0 18px rgba(255, 30, 60, 0.55);
    animation: moltenMedallionGlow 2.8s ease-in-out infinite alternate;
  }

  .medallion-panel.molten-active .relic-item-name {
    color: #ffd1d8;
    text-shadow: 0 0 10px rgba(255, 50, 80, 0.8);
  }

  .medallion-panel.molten-active .medallion-badge {
    color: #ff6b85;
  }

  .medallion-panel.molten-spent {
    background: radial-gradient(circle at 50% 25%, #190a0d 0%, #0d0406 100%);
    border: 1px solid #4a2228;
    filter: grayscale(0.8) brightness(0.65);
    opacity: 0.7;
  }

  @keyframes silverMedallionGlow {
    0% {
      border-color: #62b0e8;
      box-shadow: 0 0 12px rgba(98, 176, 232, 0.35);
    }
    100% {
      border-color: #b8e2ff;
      box-shadow: 0 0 22px rgba(160, 220, 255, 0.7);
    }
  }

  @keyframes moltenMedallionGlow {
    0% {
      border-color: #e62045;
      box-shadow: 0 0 12px rgba(230, 32, 69, 0.4);
    }
    100% {
      border-color: #ff6b85;
      box-shadow: 0 0 24px rgba(255, 60, 90, 0.75);
    }
  }
</style>
