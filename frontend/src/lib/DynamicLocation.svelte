<script>
  import { getLocationConfig } from './locationThemes.js';

  export let bossId = 'the_lanterns';

  $: config = getLocationConfig(bossId);
</script>

<div class="dynamic-location-root loc-{bossId}">
  <!-- Ambient Gradient Background matching location mood -->
  <div class="loc-backdrop"></div>

  {#if bossId === 'the_lanterns'}
    <!-- LOCATION 1: Cold Forest Highway -->
    <div class="forest-environment">
      <!-- Curved pine silhouettes in stormy distance -->
      <svg class="pines-silhouette" viewBox="0 0 1440 600" preserveAspectRatio="none">
        <path d="M 0,600 L 0,380 C 40,360 80,310 120,380 L 140,290 L 165,370 C 220,320 280,300 330,420 L 370,330 L 400,390 C 460,310 520,340 570,410 L 610,310 L 650,430 C 720,340 780,290 840,400 L 890,260 L 930,390 C 1000,300 1060,330 1120,430 L 1180,310 L 1220,410 C 1280,340 1350,300 1440,360 L 1440,600 Z" fill="#070c09" opacity="0.85" />
        <path d="M 0,600 L 0,440 C 60,400 110,360 180,450 L 220,390 L 250,470 C 340,390 410,420 480,480 L 530,380 L 580,480 C 660,390 750,400 810,480 L 860,390 L 920,490 C 1020,410 1100,430 1190,500 L 1260,420 L 1320,490 C 1380,420 1420,450 1440,480 L 1440,600 Z" fill="#0c140f" opacity="0.9" />
      </svg>

      <!-- Slow undulating cold mist drifts across the ground -->
      <div class="mist-layer m1"></div>
      <div class="mist-layer m2"></div>

      <!-- Rain streaks running down cold pane -->
      <div class="rain-pane">
        {#each Array(14) as _, i}
          <div 
            class="rain-streak" 
            style="left: {i * 7 + (i % 3) * 2}%; animation-delay: {i * 0.28}s; animation-duration: {1.1 + (i % 4) * 0.3}s;"
          ></div>
        {/each}
      </div>

      <!-- Table Decor: Pine moss & amber resin droplets -->
      <div class="table-overlay pine-moss-decor">
        <div class="moss-patch p1" title="Сырой мох"></div>
        <div class="moss-patch p2" title="Влажная кора"></div>
        <div class="resin-droplet r1" title="Капля сосновой смолы"></div>
        <div class="resin-droplet r2"></div>
      </div>
    </div>

  {:else if bossId === 'kingdom_of'}
    <!-- LOCATION 2: Ruined Asylum Chamber -->
    <div class="asylum-environment">
      <!-- Swinging rusty chains suspended from high wooden beams -->
      <div class="hanging-chains">
        <div class="chain c1"></div>
        <div class="chain c2"></div>
        <div class="chain c3"></div>
      </div>

      <!-- Walls of manuscript pages and papers detaching and drifting -->
      <div class="paper-wall-fragments">
        <div class="manuscript-sheet ms1">«Король в лохмотьях»</div>
        <div class="manuscript-sheet ms2">«Палата №21»</div>
        <div class="manuscript-sheet ms3">«Чернила смыты»</div>
        <div class="manuscript-sheet ms4">«Стены рушатся»</div>
      </div>

      <!-- Table Decor: Cracked white stone, ink puddles, broken elixir phials -->
      <div class="table-overlay stone-chamber-decor">
        <div class="stone-fissure sf1"></div>
        <div class="ink-puddle ip1" title="Чернильная лужа"></div>
        <div class="shattered-phial sp1" title="Осколок склянки с микстурой"></div>
      </div>
    </div>

  {:else if bossId === 'let_it_go'}
    <!-- LOCATION 3: Tavern Underworld / Blazing Sacrificial Hearth -->
    <div class="underworld-environment">
      <!-- Blazing rafters and fiery glow -->
      <div class="blazing-rafters">
        <div class="beam b1"></div>
        <div class="beam b2"></div>
      </div>

      <!-- Swirling ash flakes and rising fiery embers -->
      <div class="cinder-vortex">
        {#each Array(18) as _, i}
          <div 
            class="cinder-particle" 
            style="
              left: {4 + i * 5.2}%; 
              animation-delay: {i * 0.22}s; 
              animation-duration: {1.6 + (i % 5) * 0.4}s;
              --drift: {(i % 2 === 0 ? 1 : -1) * (15 + (i * 3))}px;
            "
          ></div>
        {/each}
      </div>

      <!-- Hanging hot chains glowing red -->
      <div class="hell-chains">
        <div class="red-chain rc1"></div>
        <div class="red-chain rc2"></div>
      </div>

      <!-- Table Decor: Magma fissures on charred altar -->
      <div class="table-overlay altar-fissures-decor">
        <div class="magma-artery ma1"></div>
        <div class="magma-artery ma2"></div>
        <div class="charred-sacrificial-rune" title="Демоническое клеймо"></div>
      </div>
    </div>
  {/if}
</div>

<style>
  .dynamic-location-root {
    position: absolute;
    inset: 0;
    pointer-events: none;
    overflow: hidden;
    z-index: 1;
    transition: opacity 0.8s ease;
  }

  .loc-backdrop {
    position: absolute;
    inset: 0;
    transition: background 1.2s ease;
  }

  /* ================== LOCATION 1: FOREST ================== */
  .loc-the_lanterns .loc-backdrop {
    background: radial-gradient(circle at 50% 25%, #0e1c16 0%, #060e0a 60%, #020504 100%);
  }

  .forest-environment {
    position: absolute;
    inset: 0;
  }

  .pines-silhouette {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 60%;
    pointer-events: none;
    filter: blur(1.5px);
  }

  .mist-layer {
    position: absolute;
    bottom: 0;
    left: -20%;
    width: 140%;
    height: 40%;
    background: radial-gradient(ellipse at 50% 100%, rgba(130, 185, 195, 0.16) 0%, rgba(80, 140, 150, 0.08) 45%, transparent 75%);
    pointer-events: none;
  }

  .mist-layer.m1 {
    animation: mistDrift1 18s infinite alternate ease-in-out;
  }

  .mist-layer.m2 {
    bottom: 10%;
    height: 35%;
    opacity: 0.7;
    animation: mistDrift2 24s infinite alternate ease-in-out;
  }

  @keyframes mistDrift1 {
    0% { transform: translateX(0) scaleY(1); }
    100% { transform: translateX(8%) scaleY(1.15); }
  }

  @keyframes mistDrift2 {
    0% { transform: translateX(5%) scaleY(1.1); }
    100% { transform: translateX(-6%) scaleY(0.9); }
  }

  /* Rain streaks */
  .rain-pane {
    position: absolute;
    inset: 0;
    overflow: hidden;
  }

  .rain-streak {
    position: absolute;
    top: -50px;
    width: 1.5px;
    height: 40px;
    background: linear-gradient(180deg, transparent 0%, rgba(180, 220, 240, 0.45) 80%, rgba(255, 255, 255, 0.7) 100%);
    animation: rainFall linear infinite;
    opacity: 0.6;
  }

  @keyframes rainFall {
    0% { transform: translateY(0) rotate(8deg); opacity: 0.3; }
    80% { opacity: 0.75; }
    100% { transform: translateY(105vh) rotate(8deg); opacity: 0; }
  }

  /* Table Pine & Moss Decor */
  .pine-moss-decor {
    position: absolute;
    inset: 0;
  }

  .moss-patch {
    position: absolute;
    background: radial-gradient(circle, #2d451e 0%, #17260e 70%, transparent 100%);
    border-radius: 40% 60% 70% 30% / 50% 30% 70% 50%;
    opacity: 0.65;
    box-shadow: inset 0 0 6px #0d1708;
  }

  .moss-patch.p1 {
    top: 28%;
    left: 22%;
    width: 65px;
    height: 35px;
    transform: rotate(-14deg);
  }

  .moss-patch.p2 {
    bottom: 30%;
    right: 21%;
    width: 80px;
    height: 45px;
    transform: rotate(20deg);
  }

  .resin-droplet {
    position: absolute;
    background: radial-gradient(circle at 35% 35%, #ffd466 0%, #cc7711 60%, #552200 100%);
    border-radius: 50%;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.8), 0 0 8px rgba(255, 170, 40, 0.4);
    opacity: 0.8;
  }

  .resin-droplet.r1 {
    top: 36%;
    left: 28%;
    width: 10px;
    height: 14px;
    border-radius: 40% 40% 60% 60%;
  }

  .resin-droplet.r2 {
    bottom: 34%;
    right: 26%;
    width: 8px;
    height: 8px;
  }

  /* ================== LOCATION 2: ASYLUM ================== */
  .loc-kingdom_of .loc-backdrop {
    background: radial-gradient(circle at 50% 20%, #1c1524 0%, #0c0912 60%, #050407 100%);
  }

  .asylum-environment {
    position: absolute;
    inset: 0;
  }

  .hanging-chains {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 240px;
  }

  .chain {
    position: absolute;
    top: -20px;
    width: 8px;
    background-image: repeating-linear-gradient(
      180deg,
      #3a3227 0px,
      #1f1a14 5px,
      #4f4435 9px,
      #14110c 12px
    );
    border-radius: 3px;
    box-shadow: 2px 4px 8px rgba(0, 0, 0, 0.8);
    transform-origin: top center;
  }

  .chain.c1 {
    left: 18%;
    height: 180px;
    animation: chainSway 7s infinite ease-in-out;
  }

  .chain.c2 {
    right: 20%;
    height: 220px;
    animation: chainSway 9s infinite ease-in-out reverse;
  }

  .chain.c3 {
    left: 48%;
    height: 130px;
    opacity: 0.5;
    animation: chainSway 6s infinite ease-in-out 1s;
  }

  @keyframes chainSway {
    0%, 100% { transform: rotate(-3deg); }
    50% { transform: rotate(4deg); }
  }

  .paper-wall-fragments {
    position: absolute;
    inset: 0;
  }

  .manuscript-sheet {
    position: absolute;
    font-family: var(--font-title);
    font-size: 0.75rem;
    color: #4a3e2e;
    background: #dfd2b9;
    border: 1px solid #7c684d;
    padding: 6px 10px;
    border-radius: 2px;
    box-shadow: 2px 4px 10px rgba(0, 0, 0, 0.8);
    opacity: 0.35;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
    animation: paperFlutter 6s infinite ease-in-out alternate;
  }

  .manuscript-sheet.ms1 {
    top: 14%;
    left: 12%;
    transform: rotate(-12deg);
  }

  .manuscript-sheet.ms2 {
    top: 22%;
    right: 14%;
    transform: rotate(16deg);
    animation-delay: 1.5s;
  }

  .manuscript-sheet.ms3 {
    bottom: 22%;
    left: 15%;
    transform: rotate(8deg);
    animation-delay: 3s;
  }

  .manuscript-sheet.ms4 {
    bottom: 26%;
    right: 18%;
    transform: rotate(-18deg);
    animation-delay: 2s;
  }

  @keyframes paperFlutter {
    0% { transform: translateY(0) rotate(-10deg); }
    50% { transform: translateY(6px) rotate(-6deg); }
    100% { transform: translateY(-4px) rotate(-12deg); }
  }

  /* Table Asylum Decor */
  .stone-chamber-decor {
    position: absolute;
    inset: 0;
  }

  .stone-fissure {
    position: absolute;
    top: 48%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 280px;
    height: 2px;
    background: linear-gradient(90deg, transparent 0%, #000 40%, #000 70%, transparent 100%);
    box-shadow: 0 1px 1px rgba(255, 255, 255, 0.08);
  }

  .ink-puddle {
    position: absolute;
    top: 32%;
    right: 24%;
    width: 55px;
    height: 38px;
    background: #06050a;
    border-radius: 60% 40% 50% 50% / 40% 60% 40% 60%;
    box-shadow: inset 0 0 6px #1c1524, 0 1px 3px rgba(0, 0, 0, 0.9);
    opacity: 0.85;
  }

  .shattered-phial {
    position: absolute;
    bottom: 32%;
    left: 23%;
    width: 20px;
    height: 12px;
    border-right: 3px solid rgba(190, 230, 255, 0.6);
    border-bottom: 2px solid rgba(190, 230, 255, 0.4);
    transform: rotate(35deg);
    filter: drop-shadow(0 0 4px rgba(120, 200, 255, 0.4));
  }

  /* ================== LOCATION 3: UNDERWORLD ================== */
  .loc-let_it_go .loc-backdrop {
    background: radial-gradient(circle at 50% 30%, #300c06 0%, #170402 60%, #070101 100%);
  }

  .underworld-environment {
    position: absolute;
    inset: 0;
  }

  .blazing-rafters {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 160px;
  }

  .beam {
    position: absolute;
    top: -15px;
    height: 24px;
    background: #140503;
    border-bottom: 2px solid #5a1908;
    box-shadow: 0 4px 20px rgba(255, 60, 0, 0.35);
  }

  .beam.b1 {
    left: -10%;
    width: 70%;
    transform: rotate(4deg);
  }

  .beam.b2 {
    right: -10%;
    width: 70%;
    transform: rotate(-5deg);
  }

  /* Swirling Cinders & Embers */
  .cinder-vortex {
    position: absolute;
    inset: 0;
    overflow: hidden;
  }

  .cinder-particle {
    position: absolute;
    bottom: -10px;
    width: 4px;
    height: 4px;
    background: radial-gradient(circle, #ffffff 0%, #ff6600 60%, #b32400 100%);
    border-radius: 50%;
    box-shadow: 0 0 8px #ff4500, 0 0 16px #ff2200;
    animation: cinderRise linear infinite;
  }

  @keyframes cinderRise {
    0% { transform: translateY(0) translateX(0) scale(1); opacity: 0.9; }
    50% { transform: translateY(-50vh) translateX(var(--drift, 20px)) scale(1.2); opacity: 1; }
    100% { transform: translateY(-105vh) translateX(calc(var(--drift, 20px) * 2)) scale(0.3); opacity: 0; }
  }

  .hell-chains {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
  }

  .red-chain {
    position: absolute;
    top: -20px;
    width: 8px;
    height: 190px;
    background-image: repeating-linear-gradient(
      180deg,
      #701e0a 0px,
      #2e0c04 4px,
      #ff4400 8px,
      #1f0602 12px
    );
    border-radius: 3px;
    box-shadow: 0 0 12px rgba(255, 60, 0, 0.6);
  }

  .red-chain.rc1 {
    left: 17%;
    animation: chainSway 5s infinite ease-in-out;
  }

  .red-chain.rc2 {
    right: 17%;
    animation: chainSway 7s infinite ease-in-out reverse;
  }

  /* Table Altar Decor */
  .altar-fissures-decor {
    position: absolute;
    inset: 0;
  }

  .magma-artery {
    position: absolute;
    height: 3px;
    background: linear-gradient(90deg, transparent 0%, #ff4500 30%, #ffaa00 50%, #ff2200 70%, transparent 100%);
    box-shadow: 0 0 10px #ff5500, 0 0 20px #ff2200;
    animation: magmaPulse 2s infinite alternate ease-in-out;
  }

  .magma-artery.ma1 {
    top: 45%;
    left: 20%;
    width: 220px;
    transform: rotate(-16deg);
  }

  .magma-artery.ma2 {
    bottom: 35%;
    right: 22%;
    width: 250px;
    transform: rotate(12deg);
    animation-delay: 1s;
  }

  .charred-sacrificial-rune {
    position: absolute;
    top: 52%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 90px;
    height: 90px;
    border: 1px dashed rgba(255, 60, 0, 0.45);
    border-radius: 50%;
    box-shadow: inset 0 0 15px rgba(255, 40, 0, 0.3), 0 0 15px rgba(255, 60, 0, 0.2);
  }

  @keyframes magmaPulse {
    0% { opacity: 0.6; filter: brightness(1); }
    100% { opacity: 1; filter: brightness(1.5); }
  }
</style>
