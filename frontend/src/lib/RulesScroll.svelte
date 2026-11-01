<script>
  // RulesScroll: Left column ancient parchment scroll with legible, bold rules
</script>

<aside class="scroll-column rules-scroll-column" aria-label="Свиток правил">
  <div class="scroll-parchment">
    <!-- Top & Bottom Iron Mounting Rods -->
    <div class="iron-rod rod-top">
      <div class="rod-finial left"></div>
      <div class="rod-shaft"></div>
      <div class="rod-finial right"></div>
    </div>

    <!-- Scroll Body -->
    <div class="scroll-content">
      <div class="wax-seal" title="Печать Корчмы">
        <span class="seal-icon">☠</span>
      </div>

      <header class="scroll-header">
        <h2 class="scroll-title">ПРАВИЛА</h2>
      </header>

      <div class="scroll-divider"></div>

      <!-- Main Tenets with Golden/Yellow Palette & Large Typography -->
      <section class="rules-body">
        <div class="rule-card highlight-card">
          <div class="rule-badge crit-badge">⚡ КРИТ x3</div>
          <p class="rule-lead">
            <strong>Цель:</strong> собери сумму карт максимально близко к <strong>21</strong>.
          </p>
          <p class="rule-detail">
            Ровно <strong>21 очко</strong> — сокрушительный удар, наносящий <strong>3 урона</strong> боссу и оглушающий его на раунд!
          </p>
        </div>

        <div class="rule-card danger-card">
          <div class="rule-badge bust-badge">💀 ПЕРЕБОР</div>
          <p class="rule-lead">
            Сумма <strong>больше 21</strong> — ты ранишь себя сам.
          </p>
          <p class="rule-detail">
            Один огарок свечи угасает мгновенно. Потухнут все 5 свечей — корчма поглотит твою душу.
          </p>
        </div>

        <div class="rule-card trick-card">
          <div class="rule-badge boss-badge">👁️ ДИЛЕР</div>
          <p class="rule-lead">
            Каждый раунд босс пускает в ход <strong>коварные приемы</strong>.
          </p>
          <p class="rule-detail">
            Слепой свет фонаря скрывает карты, стены кромсают силу карт, а пламя сжигает дерзких.
          </p>
        </div>

        <!-- Card Values in identical frame with icon badge -->
        <div class="rule-card values-card">
          <div class="rule-badge values-badge">🎴 НОМИНАЛЫ</div>
          <div class="values-list">
            <div class="val-item">Туз (А) — <strong>11 или 1 (автоматически)</strong></div>
            <div class="val-item">Картинки (K, Q, J) — <strong>10</strong></div>
          </div>
        </div>
      </section>
    </div>

    <div class="iron-rod rod-bottom">
      <div class="rod-finial left"></div>
      <div class="rod-shaft"></div>
      <div class="rod-finial right"></div>
    </div>
  </div>
</aside>

<style>
  .rules-scroll-column {
    width: 280px;
    min-width: 280px;
    max-width: 280px;
    height: 100%;
    min-height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    position: relative;
    z-index: 15;
    padding: 10px 4px;
    box-sizing: border-box;
    font-variant-numeric: slashed-zero tabular-nums;
  }

  .scroll-parchment {
    position: relative;
    background: radial-gradient(circle at 40% 30%, #201711 0%, #17100b 60%, #0d0906 100%);
    border: 2px solid #5a3c26;
    border-radius: 6px;
    box-shadow: 
      0 10px 30px rgba(0, 0, 0, 0.95),
      inset 0 0 40px rgba(0, 0, 0, 0.8),
      0 0 1px rgba(255, 200, 120, 0.2);
    display: flex;
    flex-direction: column;
    height: 100%;
    max-height: 94vh;
    overflow: hidden;
  }

  /* Iron Mounting Rods */
  .iron-rod {
    height: 16px;
    display: flex;
    align-items: center;
    position: relative;
    z-index: 2;
    background: #110d0a;
  }

  .rod-shaft {
    flex-grow: 1;
    height: 8px;
    background: linear-gradient(180deg, #5c4735 0%, #2e2318 50%, #16110a 100%);
    border-top: 1px solid #7a614b;
    border-bottom: 1px solid #000;
    box-shadow: 0 2px 4px rgba(0,0,0,0.8);
  }

  .rod-finial {
    width: 14px;
    height: 16px;
    background: radial-gradient(circle at 35% 35%, #8a6c4f 0%, #3a2b1d 70%, #19120a 100%);
    border-radius: 50%;
    border: 1px solid #9e7f5f;
    box-shadow: 0 2px 5px rgba(0,0,0,0.9);
  }

  .scroll-content {
    flex-grow: 1;
    overflow-y: auto;
    padding: 14px 12px;
    display: flex;
    flex-direction: column;
    position: relative;
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

  /* Wax Seal */
  .wax-seal {
    position: absolute;
    top: 10px;
    right: 12px;
    width: 32px;
    height: 32px;
    background: radial-gradient(circle at 35% 35%, #b51c1c 0%, #750e0e 70%, #3b0505 100%);
    border-radius: 50%;
    border: 1px solid #d43b3b;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 3px 8px rgba(0,0,0,0.85);
  }

  .seal-icon {
    font-size: 0.95rem;
    color: #ffcccc;
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
    margin-bottom: 12px;
  }

  .rules-body {
    display: flex;
    flex-direction: column;
    gap: 10px;
    flex-grow: 1;
  }

  .rule-card {
    background: rgba(18, 12, 9, 0.82);
    border: 1px solid #6b4c2a;
    border-radius: 6px;
    padding: 10px 10px 8px;
    box-shadow: inset 0 0 12px rgba(0,0,0,0.9), 0 2px 6px rgba(0,0,0,0.6);
    position: relative;
  }

  .rule-badge {
    display: inline-block;
    font-family: var(--font-title);
    font-size: 0.95rem;
    padding: 3px 8px;
    border-radius: 3px;
    letter-spacing: 1px;
    margin-bottom: 6px;
    font-weight: bold;
  }

  .crit-badge {
    background: #183b4e;
    color: #8ce1ff;
    border: 1px solid #2f7a9e;
  }

  .bust-badge {
    background: #471919;
    color: #ff9e9e;
    border: 1px solid #8f2f2f;
  }

  .boss-badge {
    background: #3e2617;
    color: #ffd700;
    border: 1px solid #8a6125;
  }

  .values-badge {
    background: #2b1f14;
    color: #ffd700;
    border: 1px solid #996e28;
  }

  /* Golden/Yellow palette typography */
  .rule-lead {
    font-size: 1.12rem;
    line-height: 1.35;
    color: #e5be68;
    margin-bottom: 4px;
    text-shadow: 0 1px 3px rgba(0,0,0,0.95);
  }

  .rule-lead strong {
    color: #ffd700;
    font-weight: bold;
    text-shadow: 0 0 6px rgba(255, 215, 0, 0.35);
  }

  .rule-detail {
    font-size: 1rem;
    line-height: 1.35;
    color: #c99a3e;
  }

  .rule-detail strong {
    color: #fff0a8;
  }

  .values-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-top: 4px;
  }

  .val-item {
    font-size: 1.05rem;
    line-height: 1.3;
    color: #e5be68;
  }

  .val-item strong {
    color: #ffd700;
  }
</style>
