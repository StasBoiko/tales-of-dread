<script>
  import { onMount, onDestroy } from 'svelte';
  import { audioEngine } from './lib/audioEngine.js';
  import CandleBar from './lib/CandleBar.svelte';
  import ChalkScore from './lib/ChalkScore.svelte';
  import BossPortrait from './lib/BossPortrait.svelte';
  import PlayingCard from './lib/PlayingCard.svelte';
  import VolumeControl from './lib/VolumeControl.svelte';
  import RulesScroll from './lib/RulesScroll.svelte';
  import InventoryScroll from './lib/InventoryScroll.svelte';
  import DynamicLocation from './lib/DynamicLocation.svelte';
  import DarkArtifactsModal from './lib/DarkArtifactsModal.svelte';
  import TitleScreen from './lib/TitleScreen.svelte';
  import GameOverModal from './lib/GameOverModal.svelte';
  import VictoryModal from './lib/VictoryModal.svelte';
  import PrologueScreen from './lib/PrologueScreen.svelte';
  import CampfireModal from './lib/CampfireModal.svelte';
  import Part2TwistModal from './lib/Part2TwistModal.svelte';
  import Part2Arena from './lib/Part2Arena.svelte';
  import Part2EndingModal from './lib/Part2EndingModal.svelte';
  import BossCompanionAnimal from './lib/BossCompanionAnimal.svelte';

  // Game state mirroring backend
  let state = {
    phase: 'title',
    playerCandles: 5,
    maxCandles: 5,
    playerHand: [],
    playerScore: 0,
    bossIndex: 0,
    currentBoss: null,
    bossScore: 0,
    bossSecondScore: 0,
    playerHitsThisRound: 0,
    burnTriggeredThisRound: false,
    campfireHealed: false,
    artifacts: [],
    offeredArtifacts: [],
    round: 0,
    roundOutcome: '',
    lastMessage: 'Ты стоишь на пороге корчмы...',
    lastDamageDealt: 0,
    lastDamageTaken: 0,
    isCrit: false,
    screenShake: false,
    lanternTrapActive: false,
    saltTriggeredInBattle: false,
    peekedCard: null,
  };

  let showPrologue = true;
  let showMobileRules = false;
  let showMobileInventory = false;

  let isLoading = false;
  let errorMsg = '';
  let isIdleChoked = false;
  let idleTimer = null;
  let screenCracksActive = false;
  let activeCritShockwave = false;
  let isHoveringHit = false;
  let prevRoundOutcome = '';
  let prevPlayerScore = 0;

  // Inactivity tracking: if idle > 8s, candle chokes and claws appear
  function resetIdleTimer() {
    isIdleChoked = false;
    clearTimeout(idleTimer);
    idleTimer = setTimeout(() => {
      if (state.phase === 'player_turn') {
        isIdleChoked = true;
      }
    }, 8000);
  }

  // Audio & effect synchronization on state updates
  $: {
    if (state && state.currentBoss) {
      if (state.phase !== 'title' && state.phase !== 'prologue' && !showPrologue) {
        audioEngine.startPlaylist();
      }
      audioEngine.setDangerState(state.playerScore);
    }
  }

  // Tinnitus on bust (>21 points)
  $: {
    if (state && state.playerScore > 21 && prevPlayerScore <= 21) {
      audioEngine.playTinnitus();
    }
    prevPlayerScore = state ? state.playerScore : 0;
  }

  $: {
    if (state && state.isCrit) {
      activeCritShockwave = true;
      setTimeout(() => {
        activeCritShockwave = false;
      }, 750);
    }
  }

  $: {
    if (state && state.currentBoss && state.currentBoss.id === 'kingdom_of' && state.lastDamageDealt > 0) {
      screenCracksActive = true;
      setTimeout(() => {
        screenCracksActive = false;
      }, 600);
    }
  }

  let activeMedallionFlash = false;
  let medallionBannerVisible = false;

  $: {
    if (state && state.medallionSavedLife) {
      audioEngine.playMedallionSave();
      activeMedallionFlash = true;
      medallionBannerVisible = true;
      setTimeout(() => {
        activeMedallionFlash = false;
      }, 1200);
      setTimeout(() => {
        medallionBannerVisible = false;
      }, 4500);
    }
  }

  // API calls
  async function fetchAPI(endpoint, method = 'GET', body = null) {
    resetIdleTimer();
    isLoading = true;
    errorMsg = '';
    try {
      const opts = { method, headers: { 'Content-Type': 'application/json' } };
      if (body) {
        opts.body = JSON.stringify(body);
      }
      const res = await fetch(endpoint, opts);
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.error || `HTTP ${res.status}`);
      }
      const data = await res.json();

      // Check if boss was just defeated (went to <= 0 HP or transitioned to campfire/artifact_select/part2_twist)
      const bossJustDied = state.currentBoss && state.currentBoss.currentHp > 0 && 
        ((data.currentBoss && data.currentBoss.currentHp <= 0) || 
         (['artifact_select', 'campfire', 'part2_twist'].includes(data.phase) && state.phase !== data.phase));

      if (bossJustDied) {
        // Keep arena with dying boss for 1.8s
        const defeatedBoss = data.currentBoss ? { ...data.currentBoss, currentHp: 0 } : { ...state.currentBoss, currentHp: 0 };
        state = {
          ...data,
          phase: state.phase, // keep current arena phase during animation
          currentBoss: defeatedBoss
        };
        setTimeout(() => {
          state = data;
        }, 1800);
      } else {
        state = data;
      }
    } catch (err) {
      console.error(`API Error on ${endpoint}:`, err);
      errorMsg = err.message;
    } finally {
      isLoading = false;
    }
  }

  async function handleStart() {
    audioEngine.init().catch(() => {});
    audioEngine.startPlaylist();
    await fetchAPI('/api/start', 'POST');
  }

  async function handleStartNightmare() {
    audioEngine.init().catch(() => {});
    audioEngine.startPlaylist();
    await fetchAPI('/api/nightmare/start', 'POST');
  }

  async function handleEnterFromPrologue() {
    showPrologue = false;
    audioEngine.init().catch(() => {});
    audioEngine.startPlaylist();
    await fetchAPI('/api/start', 'POST');
  }

  async function handleResetToPrologue(skipConfirm = false) {
    if (skipConfirm !== true && !window.confirm('Прервать текущую партию и вернуться к вступительной истории пролога? Весь текущий прогресс будет сброшен.')) {
      return;
    }
    await fetchAPI('/api/reset', 'POST');
    showPrologue = true;
  }

  async function handleHit() {
    audioEngine.playCardRustle();
    await fetchAPI('/api/hit', 'POST');
  }

  async function handleStand() {
    audioEngine.playCardThud();
    await fetchAPI('/api/stand', 'POST');
  }

  async function handleNextRound() {
    audioEngine.playCardRustle();
    await fetchAPI('/api/next-round', 'POST');
  }

  async function handleSelectArtifact(event) {
    const { artifactId } = event.detail;
    await fetchAPI('/api/select-artifact', 'POST', { artifactId });
  }

  function handleCampfireHealed(event) {
    state = event.detail;
  }

  function handleCampfireLeft(event) {
    state = event.detail;
  }

  async function handleStartPart2Battle() {
    await fetchAPI('/api/part2/start', 'POST');
  }

  function handlePart2Update(event) {
    state = event.detail;
  }

  onMount(async () => {
    // Initial fetch of current server state
    await fetchAPI('/api/state');
    if (state.phase !== 'title') {
      showPrologue = false;
    }
    window.addEventListener('mousemove', resetIdleTimer);
    window.addEventListener('keydown', resetIdleTimer);
    resetIdleTimer();
  });

  onDestroy(() => {
    clearTimeout(idleTimer);
    window.removeEventListener('mousemove', resetIdleTimer);
    window.removeEventListener('keydown', resetIdleTimer);
  });
</script>

<!-- 1. CINEMATIC PROLOGUE SCREEN -->
{#if showPrologue}
  <PrologueScreen on:enter={handleEnterFromPrologue} />
{/if}

<main 
  class="tavern-table theme-{state.currentBoss ? state.currentBoss.tableTheme : 'pine_boards'}"
  class:shake-screen={state.screenShake}
>
  <!-- Dynamic Boss-Specific Environmental Location -->
  <DynamicLocation bossId={state.currentBoss ? state.currentBoss.id : 'the_lanterns'} />

  <!-- Ambient Candle Lighting & Shadows -->
  <div class="candle-lighting" class:candle-choked={isIdleChoked}></div>
  <div class="creeping-claws" class:visible={isIdleChoked}></div>
  <div class="screen-cracks" class:active={screenCracksActive}></div>

  <!-- Cold Silver Flash & Life-Save Banner on Medallion Trigger -->
  {#if activeMedallionFlash}
    <div class="silver-medallion-flash-screen" aria-hidden="true"></div>
  {/if}

  {#if medallionBannerVisible}
    <div class="silver-medallion-alert-banner" role="alert">
      <div class="alert-icon-ring">
        <span class="alert-icon">🥈</span>
      </div>
      <div class="alert-content">
        <span class="alert-title">СЕРЕБРЯНЫЙ ОБЕРЕГ СПАС ВАШУ ЖИЗНЬ!</span>
        <span class="alert-sub">Смертельный удар поглощен. Последняя свеча судьбы продолжает гореть.</span>
      </div>
    </div>
  {/if}

  {#if activeCritShockwave}
    <div class="crit-shockwave"></div>
  {/if}

  <!-- TOP BAR: Mobile drawers toggles + Rustic Tavern Acoustics / Volume Slider -->
  {#if !['part2_twist', 'part2_battle', 'part2_ending'].includes(state.phase)}
    <div class="top-nav-bar">
      <div class="nav-actions-left">
        <button 
          type="button" 
          class="prologue-nav-btn"
          on:click={() => handleResetToPrologue(false)}
          aria-label="В пролог (заново)"
          title="Перейти к вступительной истории и начать заново"
        >
          <span class="nav-icon">📖</span>
          <span class="nav-label">В пролог</span>
        </button>

        <div class="mobile-drawer-triggers">
          <button 
            type="button" 
            class="mobile-nav-btn rules-toggle"
            on:click={() => showMobileRules = !showMobileRules}
            aria-label="Открыть правила"
          >
            <span class="nav-icon">📜</span>
            <span class="nav-label">Правила</span>
          </button>

          <button 
            type="button" 
            class="mobile-nav-btn inventory-toggle"
            on:click={() => showMobileInventory = !showMobileInventory}
            aria-label="Открыть реликвии"
          >
            <span class="nav-icon">🗝️</span>
            <span class="nav-label">Реликвии ({state.artifacts ? state.artifacts.length : 0})</span>
          </button>
        </div>
      </div>

      <VolumeControl />
    </div>
  {/if}

  <!-- Mobile Rules Drawer / Bottom Sheet -->
  {#if showMobileRules}
    <div 
      class="mobile-drawer-backdrop" 
      role="button"
      tabindex="0"
      on:click={() => showMobileRules = false}
      on:keydown={(e) => e.key === 'Escape' && (showMobileRules = false)}
    >
      <div 
        class="mobile-drawer-sheet" 
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        on:click|stopPropagation
        on:keydown|stopPropagation
      >
        <div class="drawer-header">
          <span class="drawer-title">📜 Свиток Правил Корчмы</span>
          <button class="drawer-close-btn" on:click={() => showMobileRules = false}>✕</button>
        </div>
        <div class="drawer-content">
          <RulesScroll />
        </div>
      </div>
    </div>
  {/if}

  <!-- Mobile Inventory Drawer / Bottom Sheet -->
  {#if showMobileInventory}
    <div 
      class="mobile-drawer-backdrop" 
      role="button"
      tabindex="0"
      on:click={() => showMobileInventory = false}
      on:keydown={(e) => e.key === 'Escape' && (showMobileInventory = false)}
    >
      <div 
        class="mobile-drawer-sheet" 
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        on:click|stopPropagation
        on:keydown|stopPropagation
      >
        <div class="drawer-header">
          <span class="drawer-title">🗝️ Шкатулка Темных Реликвий</span>
          <button class="drawer-close-btn" on:click={() => showMobileInventory = false}>✕</button>
        </div>
        <div class="drawer-content">
          <InventoryScroll artifacts={state.artifacts} />
        </div>
      </div>
    </div>
  {/if}

  <!-- Modals & Fullscreen Screens -->
  {#if state.phase === 'title' && !showPrologue}
    <TitleScreen on:start={handleStart} />
  {/if}

  {#if state.phase === 'campfire'}
    <CampfireModal 
      bossIndex={state.bossIndex}
      playerCandles={state.playerCandles}
      maxCandles={state.maxCandles}
      campfireHealed={state.campfireHealed}
      isNightmare={state.isNightmare}
      on:healed={handleCampfireHealed}
      on:leave={handleCampfireLeft}
    />
  {/if}

  {#if state.phase === 'artifact_select'}
    <DarkArtifactsModal 
      offeredArtifacts={state.offeredArtifacts} 
      bossName={state.currentBoss ? state.currentBoss.name : 'Босс'}
      on:select={handleSelectArtifact}
    />
  {/if}

  {#if state.phase === 'game_over'}
    <GameOverModal on:restart={handleStart} />
  {/if}

  {#if state.phase === 'victory'}
    <VictoryModal on:restart={handleStart} />
  {/if}

  {#if state.phase === 'part2_twist'}
    <Part2TwistModal on:startBattle={handleStartPart2Battle} />
  {/if}

  {#if state.phase === 'part2_battle' && state.part2}
    <Part2Arena 
      part2={state.part2} 
      playerCandles={state.playerCandles} 
      maxCandles={state.maxCandles} 
      on:update={handlePart2Update} 
    />
  {/if}

  {#if state.phase === 'part2_ending'}
    <Part2EndingModal 
      chosenEnding={state.part2 ? state.part2.chosen_ending : ''} 
      isNightmare={state.isNightmare}
      playerCandles={state.playerCandles}
      maxCandles={state.maxCandles}
      on:restart={(e) => {
        if (e.detail && e.detail.reset) {
          handleResetToPrologue(true);
        } else if (e.detail && e.detail.nightmare) {
          handleStartNightmare();
        } else {
          handleStart();
        }
      }} 
    />
  {/if}

  <!-- Active Duel Arena with Symmetrical 3-Column Layout -->
  {#if state.currentBoss && state.phase !== 'title' && !['part2_twist', 'part2_battle', 'part2_ending'].includes(state.phase)}
    <div class="tavern-arena-grid">
      <!-- LEFT COLUMN (Desktop): Symmetrical Ancient Rules Scroll -->
      <div class="desktop-column left-column">
        <RulesScroll />
      </div>

      <!-- CENTER COLUMN: Boss Arena, Card Table & Focal Hearth -->
      <div class="arena-center-stage">
        <!-- UPPER: Boss Section & Dealer Cards -->
        <div class="top-duel-area">
          <div class="boss-flank-row">
            <div class="boss-left-candle">
              <div class="central-wick">
                <div class="candle-flame"></div>
                <div class="candle-pedestal"></div>
              </div>
            </div>

            <BossPortrait 
              boss={state.currentBoss}
              playerScore={state.playerScore}
              isRoundOver={state.phase === 'round_over'}
              roundOutcome={state.roundOutcome}
              isCrit={state.isCrit}
              isHoveringHit={isHoveringHit}
              isIdleChoked={isIdleChoked}
            />

            <!-- Symmetrical Right Companion Animal with interactive animation & sounds -->
            <div class="boss-right-companion">
              <BossCompanionAnimal bossId={state.currentBoss.id} />
            </div>
          </div>

          <!-- Boss Hands & Scores: Primary Hand and Optional Second Hand (Archmage) -->
          <div class="dealer-hands-container">
            <!-- Primary Hand -->
            <div class="dealer-hand-wrapper">
              <ChalkScore score={state.bossScore} label={state.currentBoss.id === 'let_it_go' ? "Рука 1" : "Рука Босса"} />
              <div class="cards-fan">
                {#each state.currentBoss.hand as card, i (card.id || i)}
                  <PlayingCard {card} index={i} isDealer={true} />
                {/each}
              </div>
            </div>

            <!-- Boss 3 (Archmage) Second Hand -->
            {#if state.currentBoss.id === 'let_it_go' && state.currentBoss.secondHand && state.currentBoss.secondHand.length > 0}
              <div class="dealer-hand-wrapper second-hand">
                <ChalkScore score={state.bossSecondScore} label="Рука 2" />
                <div class="cards-fan">
                  {#each state.currentBoss.secondHand as card, i (`second-${card.id || i}`)}
                    <PlayingCard {card} index={i} isDealer={true} />
                  {/each}
                </div>
              </div>
            {/if}
          </div>
        </div>

        <!-- MIDDLE: Clean Table Center & Action Message Banner -->
        <div class="table-center-hearth">
          <div class="action-banner">
            <div class="last-message-box">
              {state.lastMessage}
            </div>

            {#if state.lanternTrapActive}
              <div class="status-warning lantern-warning">
                🏮 Слепящий луч фонаря скрыл одну из твоих карт! Держи счет в памяти.
              </div>
            {/if}

            {#if state.burnTriggeredThisRound}
              <div class="status-warning burn-warning">
                🔥 Ожог пальцев! Взяв 3 карты у Архимага, ты потерял 1 свечу.
              </div>
            {/if}
          </div>
        </div>

        <!-- LOWER: Player Area with Cards, Controls & Candle Health Bar -->
        <div class="bottom-duel-area">
          <!-- Player Cards & Chalk Score -->
          <div class="player-hand-wrapper">
            <ChalkScore 
              score={state.playerScore} 
              label="Твой Счет" 
              isObscured={state.lanternTrapActive || (state.playerHand && state.playerHand.some(c => c.isBlinded))}
            />

            <div class="cards-fan">
              {#each state.playerHand as card, i (card.id || i)}
                <PlayingCard {card} index={i} isDealer={false} />
              {/each}
            </div>

            <!-- Raven's Feather Peeked Card Preview -->
            {#if state.peekedCard}
              <div class="peek-container" title="Воронье Перо: верхняя карта колоды">
                <span class="peek-badge">🪶 Подглядеть</span>
                <span class="peek-card-name">{state.peekedCard.name} ({state.peekedCard.baseValue})</span>
              </div>
            {/if}
          </div>

          <!-- Action Buttons (Massive Thumb-Friendly on Mobile) -->
          <div class="controls-row">
            {#if state.phase === 'player_turn'}
              <button 
                class="dread-btn hit-btn action-bar-btn" 
                on:click={handleHit}
                on:mouseenter={() => isHoveringHit = true}
                on:mouseleave={() => isHoveringHit = false}
                disabled={isLoading}
              >
                <span class="btn-primary-text">ВЗЯТЬ КАРТУ</span>
              </button>

              <button 
                class="dread-btn stand-btn action-bar-btn" 
                on:click={handleStand} 
                disabled={isLoading}
              >
                <span class="btn-primary-text">МНЕ ХВАТИТ</span>
              </button>
            {:else if state.phase === 'round_over'}
              <button 
                class="dread-btn next-btn action-bar-btn full-span" 
                on:click={handleNextRound} 
                disabled={isLoading}
              >
                <span class="btn-primary-text">СЛЕДУЮЩИЙ РАУНД ➔</span>
              </button>
            {/if}
          </div>

          <!-- Bottom Status & Candle Life Bar -->
          <div class="status-footer">
            <CandleBar candles={state.playerCandles} maxCandles={state.maxCandles} />
          </div>
        </div>
      </div>

      <!-- RIGHT COLUMN (Desktop): Symmetrical Vault Inventory Scroll -->
      <div class="desktop-column right-column">
        <InventoryScroll artifacts={state.artifacts} />
      </div>
    </div>
  {/if}
</main>

<style>
  .top-nav-bar {
    position: absolute;
    top: 10px;
    left: 14px;
    right: 14px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    z-index: 100;
    pointer-events: none;
  }

  .top-nav-bar > * {
    pointer-events: auto;
  }

  .nav-actions-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .prologue-nav-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: linear-gradient(180deg, #2b1b13 0%, #150d09 100%);
    border: 1px solid #7d4e2d;
    border-radius: 18px;
    color: #ebd7be;
    padding: 6px 14px;
    font-family: var(--font-title);
    font-size: 0.95rem;
    cursor: pointer;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.85);
    transition: transform 0.15s ease, border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .prologue-nav-btn:hover {
    transform: scale(1.05);
    border-color: #ffaa54;
    box-shadow: 0 4px 14px rgba(255, 170, 84, 0.3);
  }

  .mobile-drawer-triggers {
    display: none; /* hidden on desktop > 900px */
    gap: 8px;
  }

  .mobile-nav-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: linear-gradient(180deg, #2b1b13 0%, #150d09 100%);
    border: 1px solid #6b4327;
    border-radius: 18px;
    color: #ebd7be;
    padding: 6px 12px;
    font-family: var(--font-title);
    font-size: 0.95rem;
    cursor: pointer;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.85);
    transition: transform 0.15s ease, border-color 0.2s ease;
  }

  .mobile-nav-btn:hover {
    transform: scale(1.05);
    border-color: #ba763f;
  }

  .nav-icon {
    font-size: 1.1rem;
  }

  /* Mobile Drawer Bottom Sheet */
  .mobile-drawer-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(3, 2, 2, 0.85);
    z-index: 2000;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    backdrop-filter: blur(5px);
    animation: fadeIn 0.25s ease-out;
  }

  .mobile-drawer-sheet {
    background: linear-gradient(180deg, #22150f 0%, #0d0805 100%);
    border-top: 2px solid #85532d;
    border-radius: 20px 20px 0 0;
    max-height: 82vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 -8px 30px rgba(0, 0, 0, 0.95);
    animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    overflow: hidden;
  }

  .drawer-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 14px 20px;
    border-bottom: 1px solid #4a2d18;
  }

  .drawer-title {
    font-family: var(--font-title);
    font-size: 1.25rem;
    color: #ffda96;
  }

  .drawer-close-btn {
    background: transparent;
    border: 1px solid #6a4427;
    color: #dfc8ad;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    cursor: pointer;
    font-size: 1.1rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .drawer-content {
    overflow-y: auto;
    padding: 16px;
    display: flex;
    justify-content: center;
  }

  @keyframes slideUp {
    0% { transform: translateY(100%); }
    100% { transform: translateY(0); }
  }

  @keyframes fadeIn {
    0% { opacity: 0; }
    100% { opacity: 1; }
  }

  .top-duel-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
    z-index: 10;
    width: 100%;
  }

  .boss-flank-row {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 100%;
    max-width: 820px;
  }

  .boss-left-candle {
    position: absolute;
    left: calc(50% - 215px);
    top: 52%;
    transform: translateY(-50%);
    z-index: 15;
    filter: drop-shadow(0 0 12px rgba(255, 140, 30, 0.6));
  }

  .boss-right-companion {
    position: absolute;
    right: calc(50% - 225px);
    top: 52%;
    transform: translateY(-50%);
    z-index: 16;
  }

  .dealer-hands-container {
    display: flex;
    align-items: center;
    gap: 18px;
    flex-wrap: wrap;
    justify-content: center;
    margin-top: 0;
  }

  .dealer-hand-wrapper {
    display: flex;
    align-items: center;
    gap: 14px;
  }

  .dealer-hand-wrapper.second-hand {
    border-left: 2px dashed rgba(160, 60, 40, 0.5);
    padding-left: 14px;
  }

  .table-center-hearth {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    position: relative;
    z-index: 10;
    margin: 2px 0;
  }

  .central-wick {
    position: relative;
    width: 36px;
    height: 38px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
  }

  .candle-flame {
    width: 16px;
    height: 26px;
    background: radial-gradient(ellipse at 50% 80%, #fff 0%, #ffe480 30%, #ff8c00 70%, #d83d00 100%);
    border-radius: 50% 50% 35% 35%;
    box-shadow: 0 0 18px #ff9a1f, 0 0 38px #ff6600;
    animation: flameSway 1.6s infinite ease-in-out;
  }

  .candle-pedestal {
    width: 32px;
    height: 12px;
    background: #473a2c;
    border-radius: 4px;
    border: 1px solid #755b41;
    box-shadow: 0 3px 6px rgba(0, 0, 0, 0.9);
  }

  .action-banner {
    display: flex;
    flex-direction: column;
    align-items: center;
    max-width: 680px;
    margin-top: 4px;
  }

  .last-message-box {
    background: rgba(14, 10, 8, 0.88);
    border: 2px solid #5a3e28;
    border-radius: 6px;
    padding: 8px 22px;
    color: #f7e6ce;
    font-size: 1.15rem;
    line-height: 1.35;
    text-align: center;
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.9), inset 0 0 10px rgba(0,0,0,0.8);
    text-shadow: 0 1px 3px rgba(0, 0, 0, 0.95);
  }

  .status-warning {
    font-size: 0.92rem;
    font-weight: bold;
    margin-top: 4px;
    padding: 4px 14px;
    border-radius: 4px;
    text-align: center;
  }

  .status-warning.lantern-warning {
    color: #ffd269;
    background: rgba(40, 20, 5, 0.8);
    border: 1px solid #9e6423;
    text-shadow: 0 0 10px rgba(255, 180, 40, 0.8);
  }

  .status-warning.burn-warning {
    color: #ff9977;
    background: rgba(50, 10, 10, 0.8);
    border: 1px solid #a63228;
    text-shadow: 0 0 10px rgba(255, 70, 50, 0.8);
  }

  .bottom-duel-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
    z-index: 10;
    margin-top: -10px;
    padding-bottom: 4px;
    width: 100%;
  }

  .player-hand-wrapper {
    display: flex;
    align-items: center;
    gap: 16px;
    position: relative;
  }

  .cards-fan {
    display: flex;
    gap: 10px;
    min-height: 130px;
    align-items: center;
    flex-wrap: nowrap;
    overflow-x: auto;
    overflow-y: visible;
    padding: 6px 2px;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .cards-fan::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
  }

  .peek-container {
    background: rgba(22, 16, 28, 0.92);
    border: 2px dashed #9371b8;
    border-radius: 6px;
    padding: 8px 14px;
    display: flex;
    flex-direction: column;
    align-items: center;
    box-shadow: 0 4px 12px rgba(0,0,0,0.85);
  }

  .peek-badge {
    font-size: 0.85rem;
    color: #d1b5f2;
    text-transform: uppercase;
    letter-spacing: 1px;
    font-weight: bold;
  }

  .peek-card-name {
    font-size: 1.05rem;
    color: #ffffff;
    font-weight: 900;
    margin-top: 2px;
    text-shadow: 0 0 6px #a56eff;
  }

  .controls-row {
    display: flex;
    gap: 18px;
    margin: 10px 0 6px;
    width: 100%;
    max-width: 520px;
    justify-content: center;
  }

  .action-bar-btn {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 10px 20px;
    min-height: 52px;
  }

  .action-bar-btn.full-span {
    flex: 2;
  }

  .btn-primary-text {
    font-size: 1.35rem;
    font-weight: 900;
    letter-spacing: 1.5px;
  }

  .btn-badge-sub {
    font-size: 0.72rem;
    opacity: 0.8;
    letter-spacing: 1px;
    text-transform: uppercase;
  }

  .hit-btn {
    border-color: #b06f3e;
  }

  .hit-btn:hover {
    border-color: #fca85d;
    box-shadow: 0 0 20px rgba(255, 140, 40, 0.5);
  }

  .stand-btn {
    border-color: #6e5241;
  }

  .next-btn {
    border-color: #3f7b4c;
    background: linear-gradient(180deg, #1f4227 0%, #0d2413 100%);
  }

  .next-btn:hover {
    border-color: #5dd478;
    box-shadow: 0 0 24px rgba(60, 200, 80, 0.5);
  }

  .status-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 4px;
  }

  @keyframes flameSway {
    0%, 100% { transform: rotate(-4deg); }
    50% { transform: rotate(5deg) scaleY(1.06); }
  }

  /* Responsive Rules: Mobile screens (< 900px) */
  @media (max-width: 900px) {
    .top-nav-bar {
      position: relative;
      top: 0;
      left: 0;
      right: 0;
      padding: 8px 10px;
      background: rgba(15, 10, 7, 0.95);
      border-bottom: 1px solid #4a2f1b;
      display: flex;
      justify-content: space-between;
      align-items: center;
      z-index: 50;
      pointer-events: auto;
    }

    .mobile-drawer-triggers {
      display: flex;
      gap: 8px;
    }

    .desktop-column {
      display: none !important;
    }

    .tavern-arena-grid {
      grid-template-columns: 1fr;
      padding: 8px 6px 14px;
      gap: 0;
      height: calc(100vh - 56px);
      overflow-y: auto;
    }

    .arena-center-stage {
      padding: 0 4px;
      width: 100%;
      height: 100%;
      justify-content: space-between;
    }

    .top-duel-area {
      transform: scale(0.92);
      transform-origin: top center;
      margin-bottom: -15px;
    }

    .dealer-hands-container {
      gap: 8px;
    }

    .table-center-hearth {
      margin: 2px 0;
      transform: scale(0.92);
    }

    .bottom-duel-area {
      width: 100%;
    }

    .player-hand-wrapper {
      transform: scale(0.95);
      gap: 10px;
    }

    .cards-fan {
      gap: 6px;
      justify-content: center;
      min-height: 140px;
    }

    .controls-row {
      max-width: 100%;
      width: 100%;
      padding: 6px 8px;
      gap: 10px;
      margin: 6px 0;
    }

    .action-bar-btn {
      min-height: 54px;
      padding: 8px 14px;
    }

    .btn-primary-text {
      font-size: 1.18rem;
    }

    .last-message-box {
      font-size: 0.94rem;
      padding: 5px 12px;
    }
  }

  /* Silver Medallion Life Save Screen Flash & Banner */
  .silver-medallion-flash-screen {
    position: fixed;
    inset: 0;
    background: radial-gradient(circle at 50% 50%, rgba(210, 240, 255, 0.9) 0%, rgba(130, 195, 255, 0.5) 50%, rgba(40, 100, 160, 0.2) 80%, transparent 100%);
    z-index: 15000;
    pointer-events: none;
    animation: silverFlashFade 1.2s cubic-bezier(0.1, 0.9, 0.2, 1) forwards;
  }

  .silver-medallion-alert-banner {
    position: fixed;
    top: 60px;
    left: 50%;
    transform: translateX(-50%);
    background: linear-gradient(90deg, rgba(14, 28, 44, 0.95) 0%, rgba(24, 48, 72, 0.98) 50%, rgba(14, 28, 44, 0.95) 100%);
    border: 2px solid #78d2ff;
    border-radius: 10px;
    padding: 12px 24px;
    display: flex;
    align-items: center;
    gap: 16px;
    box-shadow: 
      0 10px 35px rgba(0, 0, 0, 0.9),
      0 0 25px rgba(120, 210, 255, 0.5);
    z-index: 15001;
    pointer-events: none;
    animation: silverBannerSlide 0.4s ease-out, silverBannerGlow 2s ease-in-out infinite alternate;
  }

  .alert-icon-ring {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: radial-gradient(circle, #2a4c6a 0%, #132433 100%);
    border: 1.5px solid #a6e3ff;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 0 12px rgba(166, 227, 255, 0.7);
    flex-shrink: 0;
  }

  .alert-icon {
    font-size: 1.6rem;
  }

  .alert-content {
    display: flex;
    flex-direction: column;
    text-align: left;
    gap: 3px;
  }

  .alert-title {
    font-family: var(--font-title);
    font-size: 1.15rem;
    font-weight: 900;
    letter-spacing: 1.5px;
    color: #e0f4ff;
    text-shadow: 0 0 10px rgba(120, 210, 255, 0.8);
  }

  .alert-sub {
    font-size: 0.88rem;
    color: #a5d8ff;
    letter-spacing: 0.5px;
  }

  @keyframes silverFlashFade {
    0% { opacity: 0; }
    15% { opacity: 1; }
    100% { opacity: 0; }
  }

  @keyframes silverBannerSlide {
    from { opacity: 0; transform: translate(-50%, -20px); }
    to { opacity: 1; transform: translate(-50%, 0); }
  }

  @keyframes silverBannerGlow {
    0% {
      border-color: #78d2ff;
      box-shadow: 0 10px 35px rgba(0, 0, 0, 0.9), 0 0 20px rgba(120, 210, 255, 0.4);
    }
    100% {
      border-color: #c4eeff;
      box-shadow: 0 10px 35px rgba(0, 0, 0, 0.9), 0 0 32px rgba(160, 230, 255, 0.75);
    }
  }
</style>
