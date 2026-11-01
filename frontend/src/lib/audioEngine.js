// Audio Engine utilizing Web Audio API for Tavern 21: Tales of Dread
// Provides continuous sequential playlist looping, dynamic audio ducking,
// 400Hz danger low-pass with heartbeat, 3000Hz tinnitus burst on bust,
// and procedural ambience (campfire crackle, crickets, rain, wind).

class AudioEngine {
  constructor() {
    this.ctx = null;
    this.masterGain = null;
    this.filterNode = null;
    this.musicDuckingGain = null;
    this.heartbeatGain = null;
    this.heartbeatTimer = null;

    this.musicTrack = null;
    this.musicGain = null;
    this.isMusicPlaying = false;
    this.isPaused = false;
    this.listeners = new Set();
    this.isInitialized = false;

    // Default volume: 0.28 (28% volume, balanced as requested)
    const savedVol = typeof localStorage !== 'undefined' ? localStorage.getItem('tavern_volume') : null;
    this.volume = savedVol !== null ? Math.max(0, Math.min(1, parseFloat(savedVol))) : 0.28;

    const savedMuted = typeof localStorage !== 'undefined' ? localStorage.getItem('tavern_muted') : null;
    this.isMuted = savedMuted === 'true';

    // Ambient nodes
    this.campfireGain = null;
    this.campfireInterval = null;
    this.cricketsGain = null;
    this.cricketsInterval = null;
    this.rainGain = null;
    this.rainSource = null;
    this.isCampfireActive = false;
    this.isRainActive = false;

    this.setupAutoplayUnlock();
  }

  setupAutoplayUnlock() {
    if (typeof window === 'undefined') return;
    const unlock = () => {
      this.init().catch(() => {});
      if (this.ctx && this.ctx.state === 'suspended') {
        this.ctx.resume().catch(() => {});
      }
      window.removeEventListener('pointerdown', unlock);
      window.removeEventListener('keydown', unlock);
      window.removeEventListener('click', unlock);
    };
    window.addEventListener('pointerdown', unlock, { once: true });
    window.addEventListener('keydown', unlock, { once: true });
    window.addEventListener('click', unlock, { once: true });
  }

  getVolume() {
    return this.volume;
  }

  setVolume(newVol) {
    const clamped = Math.max(0, Math.min(1, parseFloat(newVol) || 0));
    this.volume = clamped;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('tavern_volume', clamped.toString());
    }

    if (clamped > 0 && this.isMuted) {
      this.isMuted = false;
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem('tavern_muted', 'false');
      }
    } else if (clamped === 0) {
      this.isMuted = true;
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem('tavern_muted', 'true');
      }
    }

    if (this.ctx && this.ctx.state === 'suspended') {
      this.ctx.resume().catch(() => {});
    }

    if (this.masterGain && this.ctx) {
      const now = this.ctx.currentTime;
      this.masterGain.gain.cancelScheduledValues(now);
      const effectiveGain = this.isMuted ? 0 : this.volume;
      this.masterGain.gain.setValueAtTime(effectiveGain, now);
    }
    if (this.musicTrack) {
      this.musicTrack.muted = this.isMuted;
    }
    this.notifyState();
    return this.volume;
  }

  async init() {
    if (this.isInitialized && this.ctx) {
      if (this.ctx.state === 'suspended') {
        try {
          await this.ctx.resume();
        } catch (e) {
          console.warn('AudioContext resume failed:', e);
        }
      }
      return;
    }

    const AudioContextClass = window.AudioContext || window['webkitAudioContext'];
    if (!AudioContextClass) return;

    this.ctx = new AudioContextClass();
    if (this.ctx.state === 'suspended') {
      try {
        await this.ctx.resume();
      } catch (e) {
        console.warn('AudioContext initial resume failed:', e);
      }
    }

    // Master Gain
    this.masterGain = this.ctx.createGain();
    const effectiveGain = this.isMuted ? 0 : this.volume;
    this.masterGain.gain.setValueAtTime(effectiveGain, this.ctx.currentTime);
    this.masterGain.connect(this.ctx.destination);

    // Filter Node (Danger Low-Pass at 400Hz, defaults to 20000Hz)
    this.filterNode = this.ctx.createBiquadFilter();
    this.filterNode.type = 'lowpass';
    this.filterNode.frequency.setValueAtTime(20000, this.ctx.currentTime);
    this.filterNode.Q.setValueAtTime(1.2, this.ctx.currentTime);
    this.filterNode.connect(this.masterGain);

    // Music Ducking Node (drops volume during campfire/story)
    this.musicDuckingGain = this.ctx.createGain();
    this.musicDuckingGain.gain.setValueAtTime(1.0, this.ctx.currentTime);
    this.musicDuckingGain.connect(this.filterNode);

    // Heartbeat Bus connected directly to master
    this.heartbeatGain = this.ctx.createGain();
    this.heartbeatGain.gain.setValueAtTime(0.0, this.ctx.currentTime);
    this.heartbeatGain.connect(this.masterGain);

    // Prepare single unified soundtrack: Grimoras Crypt
    const audio = new Audio('/audio/grimoras_crypt.mp3');
    audio.loop = true;
    audio.crossOrigin = 'anonymous';
    audio.muted = this.isMuted;

    const source = this.ctx.createMediaElementSource(audio);
    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.0001, this.ctx.currentTime);

    source.connect(gain);
    gain.connect(this.musicDuckingGain);

    this.musicTrack = audio;
    this.musicGain = gain;

    this.isInitialized = true;
    this.startHeartbeatLoop();
  }

  // Smooth fade-in start from 0.0 to 0.85 over 2.8 seconds
  async startMainMusic() {
    if (!this.isInitialized) {
      await this.init();
    }
    if (this.ctx && this.ctx.state === 'suspended') {
      try {
        await this.ctx.resume();
      } catch (e) {}
    }
    if (!this.musicTrack || !this.musicGain) return;

    if (this.isMusicPlaying && !this.isPaused) return;

    if (this.isPaused) {
      try {
        await this.musicTrack.play();
        this.isPaused = false;
        this.isMusicPlaying = true;
        this.notifyState();
      } catch (e) {
        console.warn('Audio resume error:', e);
      }
      return;
    }

    const now = this.ctx.currentTime;
    const fadeDuration = 2.8;

    this.musicGain.gain.cancelScheduledValues(now);
    this.musicGain.gain.setValueAtTime(0.0001, now);
    this.musicGain.gain.linearRampToValueAtTime(0.85, now + fadeDuration);

    this.musicTrack.currentTime = 0;
    try {
      await this.musicTrack.play();
      this.isMusicPlaying = true;
      this.isPaused = false;
      this.notifyState();
    } catch (e) {
      console.warn('Music playback error:', e);
    }
  }

  // Smoothly fade out music track over specified seconds
  fadeTrackOut(durationSeconds = 2.0) {
    if (!this.musicGain || !this.ctx) return;
    try {
      const now = this.ctx.currentTime;
      this.musicGain.gain.cancelScheduledValues(now);
      this.musicGain.gain.setValueAtTime(this.musicGain.gain.value, now);
      this.musicGain.gain.linearRampToValueAtTime(0.0001, now + durationSeconds);
      setTimeout(() => {
        if (this.musicTrack && this.isMusicPlaying) {
          this.musicTrack.pause();
          this.isMusicPlaying = false;
          this.notifyState();
        }
      }, durationSeconds * 1000);
    } catch (e) {
      console.warn('fadeTrackOut error:', e);
    }
  }

  // Seamless alias for backwards compatibility
  async startPlaylist() {
    return this.startMainMusic();
  }

  async togglePlayPause() {
    if (!this.isInitialized) {
      await this.init();
    }
    if (this.ctx && this.ctx.state === 'suspended') {
      try {
        await this.ctx.resume();
      } catch (e) {}
    }
    if (!this.musicTrack) return false;

    if (this.isPaused || !this.isMusicPlaying) {
      try {
        if (this.musicGain && this.ctx) {
          const now = this.ctx.currentTime;
          this.musicGain.gain.cancelScheduledValues(now);
          this.musicGain.gain.setValueAtTime(0.85, now);
        }
        await this.musicTrack.play();
        this.isMusicPlaying = true;
        this.isPaused = false;
      } catch (e) {
        console.warn('Audio play error:', e);
      }
    } else {
      this.musicTrack.pause();
      this.isPaused = true;
    }
    this.notifyState();
    return !this.isPaused;
  }

  subscribe(listener) {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }

  notifyState() {
    for (const fn of this.listeners) {
      try {
        fn({ isPaused: this.isPaused, isMuted: this.isMuted, volume: this.volume });
      } catch {}
    }
  }

  // Duck music by 65% during campfire and story interludes
  setDucking(isDucked) {
    if (!this.ctx || !this.musicDuckingGain) return;
    const now = this.ctx.currentTime;
    this.musicDuckingGain.gain.cancelScheduledValues(now);

    if (isDucked) {
      // 65% reduction -> 35% remaining
      this.musicDuckingGain.gain.setTargetAtTime(0.35, now, 0.35);
      this.startCampfireAmbience();
    } else {
      this.musicDuckingGain.gain.setTargetAtTime(1.0, now, 0.4);
      this.stopCampfireAmbience();
    }
  }

  // Set adrenaline filter based on player score (18-20 = danger zone: 400 Hz cut + heartbeat)
  setDangerState(score) {
    if (!this.ctx || !this.filterNode || !this.heartbeatGain) return;
    const now = this.ctx.currentTime;

    if (score >= 18 && score <= 20) {
      // Cut highs at 400 Hz and boost heartbeat
      this.filterNode.frequency.cancelScheduledValues(now);
      this.filterNode.frequency.setTargetAtTime(400, now, 0.3);

      this.heartbeatGain.gain.cancelScheduledValues(now);
      this.heartbeatGain.gain.setTargetAtTime(0.85, now, 0.25);
    } else {
      // Normal clear audio
      this.filterNode.frequency.cancelScheduledValues(now);
      this.filterNode.frequency.setTargetAtTime(20000, now, 0.4);

      this.heartbeatGain.gain.cancelScheduledValues(now);
      this.heartbeatGain.gain.setTargetAtTime(0.0, now, 0.35);
    }
  }

  // Tinnitus FX on bust (>21): sharp 3000 Hz tone decaying over 1.5s with music volume dip
  playTinnitus() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Dip music momentarily
    if (this.musicDuckingGain) {
      this.musicDuckingGain.gain.cancelScheduledValues(now);
      this.musicDuckingGain.gain.setValueAtTime(0.15, now);
      this.musicDuckingGain.gain.exponentialRampToValueAtTime(1.0, now + 1.5);
    }

    // Piercing 3000 Hz sine oscillator
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(3000, now);
    osc.frequency.exponentialRampToValueAtTime(2800, now + 1.5);

    gain.gain.setValueAtTime(0.001, now);
    gain.gain.linearRampToValueAtTime(0.35, now + 0.05);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 1.5);

    osc.connect(gain);
    gain.connect(this.masterGain);

    osc.start(now);
    osc.stop(now + 1.55);
  }

  // Campfire procedural ambience (crackle, crickets, water drops)
  startCampfireAmbience() {
    if (!this.ctx || this.isCampfireActive) return;
    this.isCampfireActive = true;

    // Crackle burst generator
    this.campfireInterval = setInterval(() => {
      if (!this.ctx || this.isMuted || !this.isCampfireActive) return;
      if (Math.random() > 0.4) {
        this.triggerCampfireCrackle();
      }
      if (Math.random() > 0.75) {
        this.triggerWaterDrop();
      }
    }, 280);

    // Crickets chirp generator
    this.cricketsInterval = setInterval(() => {
      if (!this.ctx || this.isMuted || !this.isCampfireActive) return;
      if (Math.random() > 0.5) {
        this.triggerCricketChirp();
      }
    }, 900);
  }

  stopCampfireAmbience() {
    this.isCampfireActive = false;
    if (this.campfireInterval) clearInterval(this.campfireInterval);
    if (this.cricketsInterval) clearInterval(this.cricketsInterval);
  }

  triggerCampfireCrackle() {
    const now = this.ctx.currentTime;
    const bufferSize = Math.floor(this.ctx.sampleRate * 0.04);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);

    for (let i = 0; i < bufferSize; i++) {
      data[i] = (Math.random() * 2 - 1) * Math.exp(-i / (this.ctx.sampleRate * 0.008));
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'bandpass';
    filter.frequency.setValueAtTime(1800 + Math.random() * 1200, now);
    filter.Q.setValueAtTime(3.0, now);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.18 + Math.random() * 0.15, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.04);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  triggerCricketChirp() {
    const now = this.ctx.currentTime;
    for (let i = 0; i < 3; i++) {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      const t = now + i * 0.06;

      osc.type = 'sine';
      osc.frequency.setValueAtTime(4500 + i * 150, t);

      gain.gain.setValueAtTime(0.001, t);
      gain.gain.linearRampToValueAtTime(0.06, t + 0.015);
      gain.gain.exponentialRampToValueAtTime(0.0001, t + 0.045);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(t);
      osc.stop(t + 0.05);
    }
  }

  triggerWaterDrop() {
    const now = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.type = 'sine';
    const baseFreq = 900 + Math.random() * 400;
    osc.frequency.setValueAtTime(baseFreq, now);
    osc.frequency.exponentialRampToValueAtTime(baseFreq * 0.7, now + 0.12);

    gain.gain.setValueAtTime(0.001, now);
    gain.gain.linearRampToValueAtTime(0.12, now + 0.01);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.12);

    osc.connect(gain);
    gain.connect(this.masterGain);

    osc.start(now);
    osc.stop(now + 0.14);
  }

  // Rain and wind procedural ambience for prologue / storm
  startRainAmbience() {
    if (!this.ctx || this.isRainActive) return;
    this.isRainActive = true;

    const bufferSize = this.ctx.sampleRate * 2;
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);
    let lastOut = 0.0;

    for (let i = 0; i < bufferSize; i++) {
      const white = Math.random() * 2 - 1;
      data[i] = (lastOut + 0.02 * white) / 1.02;
      lastOut = data[i];
    }

    this.rainSource = this.ctx.createBufferSource();
    this.rainSource.buffer = buffer;
    this.rainSource.loop = true;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(800, this.ctx.currentTime);

    this.rainGain = this.ctx.createGain();
    this.rainGain.gain.setValueAtTime(0.001, this.ctx.currentTime);
    this.rainGain.gain.linearRampToValueAtTime(0.35, this.ctx.currentTime + 1.5);

    this.rainSource.connect(filter);
    filter.connect(this.rainGain);
    this.rainGain.connect(this.masterGain);

    this.rainSource.start();
  }

  stopRainAmbience() {
    if (!this.isRainActive || !this.rainGain || !this.ctx) return;
    const now = this.ctx.currentTime;
    this.rainGain.gain.linearRampToValueAtTime(0.001, now + 1.0);
    setTimeout(() => {
      if (this.rainSource) {
        try { this.rainSource.stop(); } catch (e) {}
      }
      this.isRainActive = false;
    }, 1100);
  }

  // Wick ignition / candle lighting sound
  playWickIgnition() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Match strike crackle
    const bufferSize = Math.floor(this.ctx.sampleRate * 0.12);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      data[i] = (Math.random() * 2 - 1) * Math.exp(-i / (this.ctx.sampleRate * 0.03));
    }
    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;
    const noiseFilter = this.ctx.createBiquadFilter();
    noiseFilter.type = 'highpass';
    noiseFilter.frequency.setValueAtTime(2200, now);

    const noiseGain = this.ctx.createGain();
    noiseGain.gain.setValueAtTime(0.3, now);
    noiseGain.gain.exponentialRampToValueAtTime(0.001, now + 0.12);

    noise.connect(noiseFilter);
    noiseFilter.connect(noiseGain);
    noiseGain.connect(this.masterGain);
    noise.start(now);

    // Warm flame whoosh
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    osc.type = 'sine';
    osc.frequency.setValueAtTime(120, now + 0.05);
    osc.frequency.exponentialRampToValueAtTime(260, now + 0.35);

    gain.gain.setValueAtTime(0.001, now + 0.05);
    gain.gain.linearRampToValueAtTime(0.25, now + 0.12);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.4);

    osc.connect(gain);
    gain.connect(this.masterGain);
    osc.start(now + 0.05);
    osc.stop(now + 0.45);
  }

  // Procedural 55 Hz heartbeat pulse
  startHeartbeatLoop() {
    const playBeats = () => {
      if (!this.ctx || this.isMuted) return;

      const now = this.ctx.currentTime;
      // Lub
      this.triggerSubPulse(now, 55, 0.12);
      // Dub (slightly higher pitch, slightly softer)
      this.triggerSubPulse(now + 0.22, 62, 0.09);
    };

    // ~110 BPM heartbeat loop
    this.heartbeatTimer = setInterval(playBeats, 550);
  }

  triggerSubPulse(time, freq, duration) {
    if (!this.ctx) return;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(freq, time);
    osc.frequency.exponentialRampToValueAtTime(30, time + duration);

    gain.gain.setValueAtTime(0.001, time);
    gain.gain.linearRampToValueAtTime(0.9, time + 0.02);
    gain.gain.exponentialRampToValueAtTime(0.001, time + duration);

    osc.connect(gain);
    gain.connect(this.heartbeatGain);

    osc.start(time);
    osc.stop(time + duration + 0.05);
  }

  // Procedural Sound: Dull dry wooden thud when a card hits the table
  playCardThud() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    const filter = this.ctx.createBiquadFilter();

    osc.type = 'triangle';
    osc.frequency.setValueAtTime(140, now);
    osc.frequency.exponentialRampToValueAtTime(35, now + 0.09);

    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(300, now);

    gain.gain.setValueAtTime(0.7, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.1);

    osc.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    osc.start(now);
    osc.stop(now + 0.12);
  }

  // Procedural Sound: Crisp parchment paper rustle when drawing a card
  playCardRustle() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;
    const bufferSize = this.ctx.sampleRate * 0.15;
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);

    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'bandpass';
    filter.frequency.setValueAtTime(3200, now);
    filter.Q.setValueAtTime(2.0, now);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.01, now);
    gain.gain.linearRampToValueAtTime(0.4, now + 0.03);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.15);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  // Procedural Sound: Candle extinction hiss
  playCandleHiss() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;
    const duration = 0.7;
    const bufferSize = Math.floor(this.ctx.sampleRate * duration);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);

    for (let i = 0; i < bufferSize; i++) {
      data[i] = (Math.random() * 2 - 1) * Math.exp(-i / (this.ctx.sampleRate * 0.25));
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'highpass';
    filter.frequency.setValueAtTime(1500, now);
    filter.frequency.linearRampToValueAtTime(4500, now + duration);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.5, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + duration);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  // Procedural Sound: Chalk scratch on wooden table
  playChalkScratch() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;
    const duration = 0.18;
    const bufferSize = Math.floor(this.ctx.sampleRate * duration);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);

    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'bandpass';
    filter.frequency.setValueAtTime(2400, now);
    filter.Q.setValueAtTime(4.0, now);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.25, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + duration);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  // Procedural Sound: Metallic strike / clatter
  playBossStrike() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    const freqs = [280, 410, 830];
    freqs.forEach((f) => {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = 'sawtooth';
      osc.frequency.setValueAtTime(f, now);
      osc.frequency.exponentialRampToValueAtTime(f * 0.4, now + 0.35);

      gain.gain.setValueAtTime(0.35, now);
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.35);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(now);
      osc.stop(now + 0.4);
    });
  }

  // Procedural Sound: Shattering lantern glass
  playGlassShatter() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    for (let i = 0; i < 4; i++) {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      const timeOffset = now + i * 0.03;
      const freq = 1800 + i * 900;

      osc.type = 'sine';
      osc.frequency.setValueAtTime(freq, timeOffset);

      gain.gain.setValueAtTime(0.4 / (i + 1), timeOffset);
      gain.gain.exponentialRampToValueAtTime(0.001, timeOffset + 0.25);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(timeOffset);
      osc.stop(timeOffset + 0.3);
    }
  }

  // Procedural Sound: Thunder rumble
  playThunderRumble() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;
    const duration = 2.2;
    const bufferSize = Math.floor(this.ctx.sampleRate * duration);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);

    let lastOut = 0.0;
    for (let i = 0; i < bufferSize; i++) {
      const white = Math.random() * 2 - 1;
      data[i] = (lastOut + 0.02 * white) / 1.02;
      lastOut = data[i];
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(95, now);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.01, now);
    gain.gain.linearRampToValueAtTime(0.8, now + 0.3);
    gain.gain.exponentialRampToValueAtTime(0.001, now + duration);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  // Procedural Sound: Iron locking bolt
  playIronBolt() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.type = 'square';
    osc.frequency.setValueAtTime(160, now);
    osc.frequency.exponentialRampToValueAtTime(80, now + 0.2);

    gain.gain.setValueAtTime(0.3, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.2);

    osc.connect(gain);
    gain.connect(this.masterGain);

    osc.start(now);
    osc.stop(now + 0.25);
  }

  // Antique brass bell chime for mute toggle
  playBellChime() {
    if (!this.ctx) return;
    const now = this.ctx.currentTime;

    const partials = [1200, 2400, 3600];
    partials.forEach((p, idx) => {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = 'sine';
      osc.frequency.setValueAtTime(p, now);

      gain.gain.setValueAtTime(0.25 / (idx + 1), now);
      gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.8 + idx * 0.2);

      osc.connect(gain);
      gain.connect(this.ctx.destination);

      osc.start(now);
      osc.stop(now + 1.2);
    });
  }

  // Toggle Mute with Bell feedback
  async toggleMute() {
    if (!this.isInitialized) {
      await this.init();
    }
    if (this.ctx && this.ctx.state === 'suspended') {
      await this.ctx.resume();
    }
    this.isMuted = !this.isMuted;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('tavern_muted', this.isMuted ? 'true' : 'false');
    }
    this.playBellChime();

    if (this.masterGain && this.ctx) {
      const now = this.ctx.currentTime;
      this.masterGain.gain.cancelScheduledValues(now);
      this.masterGain.gain.setValueAtTime(this.isMuted ? 0 : this.volume, now);
    }
    if (this.musicTrack) {
      this.musicTrack.muted = this.isMuted;
    }

    this.notifyState();
    return this.isMuted;
  }

  // Part II: Heavy iron door kick
  playHeavyDoorsKick() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Sub-bass heavy thump
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    osc.type = 'triangle';
    osc.frequency.setValueAtTime(110, now);
    osc.frequency.exponentialRampToValueAtTime(28, now + 0.6);

    gain.gain.setValueAtTime(0.9, now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.7);

    osc.connect(gain);
    gain.connect(this.masterGain);
    osc.start(now);
    osc.stop(now + 0.75);

    // Iron latch clash
    this.playIronBolt();
  }

  // Part II: Echoing stone footsteps
  playStoneFootsteps() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;
    for (let step = 0; step < 3; step++) {
      const t = now + step * 0.45;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      const filter = this.ctx.createBiquadFilter();

      osc.type = 'sine';
      osc.frequency.setValueAtTime(180, t);
      osc.frequency.exponentialRampToValueAtTime(70, t + 0.12);

      filter.type = 'lowpass';
      filter.frequency.setValueAtTime(450, t);

      gain.gain.setValueAtTime(0.35, t);
      gain.gain.exponentialRampToValueAtTime(0.001, t + 0.14);

      osc.connect(filter);
      filter.connect(gain);
      gain.connect(this.masterGain);

      osc.start(t);
      osc.stop(t + 0.16);
    }
  }

  // Part II: Demonic Witch Transformation roar and dark sweep
  playWitchTransformation() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Deep sub-drone rising
    const subOsc = this.ctx.createOscillator();
    const subGain = this.ctx.createGain();
    subOsc.type = 'sawtooth';
    subOsc.frequency.setValueAtTime(55, now);
    subOsc.frequency.exponentialRampToValueAtTime(140, now + 1.8);

    subGain.gain.setValueAtTime(0.01, now);
    subGain.gain.linearRampToValueAtTime(0.5, now + 0.4);
    subGain.gain.exponentialRampToValueAtTime(0.001, now + 2.0);

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(300, now);
    filter.frequency.linearRampToValueAtTime(1200, now + 1.5);

    subOsc.connect(filter);
    filter.connect(subGain);
    subGain.connect(this.masterGain);

    subOsc.start(now);
    subOsc.stop(now + 2.1);

    // Eerie high pitch resonance
    const highOsc = this.ctx.createOscillator();
    const highGain = this.ctx.createGain();
    highOsc.type = 'sine';
    highOsc.frequency.setValueAtTime(660, now + 0.2);
    highOsc.frequency.linearRampToValueAtTime(880, now + 1.6);

    highGain.gain.setValueAtTime(0.001, now + 0.2);
    highGain.gain.linearRampToValueAtTime(0.3, now + 0.7);
    highGain.gain.exponentialRampToValueAtTime(0.001, now + 1.9);

    highOsc.connect(highGain);
    highGain.connect(this.masterGain);

    highOsc.start(now + 0.2);
    highOsc.stop(now + 2.0);
  }

  // Part II: Bone dice rattling inside skull cup
  playDiceRattle() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    for (let i = 0; i < 6; i++) {
      const t = now + i * 0.045;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      osc.type = 'triangle';
      osc.frequency.setValueAtTime(620 + Math.random() * 300, t);
      osc.frequency.exponentialRampToValueAtTime(240, t + 0.03);

      gain.gain.setValueAtTime(0.35, t);
      gain.gain.exponentialRampToValueAtTime(0.001, t + 0.035);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(t);
      osc.stop(t + 0.04);
    }
  }

  // Part II: Bone dice tumbling onto obsidian stone
  playDiceTumble() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    for (let i = 0; i < 5; i++) {
      const t = now + i * 0.05 + Math.random() * 0.02;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = 'triangle';
      osc.frequency.setValueAtTime(480 + Math.random() * 250, t);
      osc.frequency.exponentialRampToValueAtTime(180, t + 0.04);

      gain.gain.setValueAtTime(0.45, t);
      gain.gain.exponentialRampToValueAtTime(0.001, t + 0.045);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(t);
      osc.stop(t + 0.05);
    }
  }

  // Part II Ending 1: Roaring purifying flames
  playInfernoRoar() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;
    const duration = 3.5;
    const bufferSize = Math.floor(this.ctx.sampleRate * duration);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);

    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(240, now);
    filter.frequency.linearRampToValueAtTime(800, now + 1.2);
    filter.frequency.exponentialRampToValueAtTime(180, now + duration);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.05, now);
    gain.gain.linearRampToValueAtTime(0.75, now + 0.8);
    gain.gain.exponentialRampToValueAtTime(0.001, now + duration);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  // Part II Ending 2: Medallion curse transfer & chain rattle
  playMedallionCurse() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Metallic chain rattle
    for (let i = 0; i < 4; i++) {
      const t = now + i * 0.08;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      osc.type = 'sawtooth';
      osc.frequency.setValueAtTime(950 - i * 120, t);

      gain.gain.setValueAtTime(0.3, t);
      gain.gain.exponentialRampToValueAtTime(0.001, t + 0.1);

      osc.connect(gain);
      gain.connect(this.masterGain);

      osc.start(t);
      osc.stop(t + 0.12);
    }

    // Ominous sub tone
    const sub = this.ctx.createOscillator();
    const subGain = this.ctx.createGain();
    sub.type = 'sine';
    sub.frequency.setValueAtTime(65, now + 0.3);
    sub.frequency.exponentialRampToValueAtTime(35, now + 2.5);

    subGain.gain.setValueAtTime(0.001, now + 0.3);
    subGain.gain.linearRampToValueAtTime(0.7, now + 0.8);
    subGain.gain.exponentialRampToValueAtTime(0.001, now + 2.8);

    sub.connect(subGain);
    subGain.connect(this.masterGain);

    sub.start(now + 0.3);
    sub.stop(now + 3.0);
  }

  // Silver Medallion Life Save: Dull metallic strike followed by a fading resonant silver chime (no glass shatter)
  playMedallionSave() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // 1. Dull metallic clang / impact
    const strikeOsc = this.ctx.createOscillator();
    const strikeGain = this.ctx.createGain();
    strikeOsc.type = 'triangle';
    strikeOsc.frequency.setValueAtTime(240, now);
    strikeOsc.frequency.exponentialRampToValueAtTime(120, now + 0.22);

    strikeGain.gain.setValueAtTime(0.7, now);
    strikeGain.gain.exponentialRampToValueAtTime(0.001, now + 0.25);

    strikeOsc.connect(strikeGain);
    strikeGain.connect(this.masterGain);
    strikeOsc.start(now);
    strikeOsc.stop(now + 0.28);

    // 2. Resonant pure silver chime harmonics (G5, D6, G6)
    const silverFreqs = [784, 1174, 1568, 2349];
    silverFreqs.forEach((freq, idx) => {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(freq, now + 0.03);
      // Soft shimmering modulation
      osc.frequency.linearRampToValueAtTime(freq + (idx % 2 === 0 ? 3 : -3), now + 2.2);

      const peakGain = 0.35 / (idx + 1);
      gain.gain.setValueAtTime(0.001, now);
      gain.gain.linearRampToValueAtTime(peakGain, now + 0.05 + idx * 0.02);
      gain.gain.exponentialRampToValueAtTime(0.0001, now + 2.6);

      osc.connect(gain);
      gain.connect(this.masterGain);
      osc.start(now);
      osc.stop(now + 2.7);
    });
  }

  // Boss Companion Animals & Witch Sounds
  playRavenCaw() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Dual raspy crow caw (sawtooth with rapid pitch fall and bandpass)
    [0, 0.28].forEach((offset) => {
      const t = now + offset;
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      const filter = this.ctx.createBiquadFilter();

      osc.type = 'sawtooth';
      osc.frequency.setValueAtTime(680, t);
      osc.frequency.exponentialRampToValueAtTime(320, t + 0.22);

      filter.type = 'bandpass';
      filter.frequency.setValueAtTime(900, t);
      filter.Q.setValueAtTime(3.5, t);

      gain.gain.setValueAtTime(0.01, t);
      gain.gain.linearRampToValueAtTime(0.35, t + 0.04);
      gain.gain.exponentialRampToValueAtTime(0.001, t + 0.24);

      osc.connect(filter);
      filter.connect(gain);
      gain.connect(this.masterGain);

      osc.start(t);
      osc.stop(t + 0.25);
    });
  }

  playCatMeow() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    osc.type = 'sine';

    // Feline meow / purr pitch arc
    osc.frequency.setValueAtTime(450, now);
    osc.frequency.exponentialRampToValueAtTime(780, now + 0.25);
    osc.frequency.exponentialRampToValueAtTime(420, now + 0.6);

    gain.gain.setValueAtTime(0.001, now);
    gain.gain.linearRampToValueAtTime(0.28, now + 0.15);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.65);

    osc.connect(gain);
    gain.connect(this.masterGain);

    osc.start(now);
    osc.stop(now + 0.7);
  }

  playViperHiss() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Filtered white noise hiss burst
    const bufferSize = Math.floor(this.ctx.sampleRate * 0.45);
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }

    const noise = this.ctx.createBufferSource();
    noise.buffer = buffer;

    const filter = this.ctx.createBiquadFilter();
    filter.type = 'highpass';
    filter.frequency.setValueAtTime(3800, now);
    filter.frequency.linearRampToValueAtTime(2400, now + 0.4);

    const gain = this.ctx.createGain();
    gain.gain.setValueAtTime(0.01, now);
    gain.gain.linearRampToValueAtTime(0.3, now + 0.08);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.45);

    noise.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    noise.start(now);
  }

  playWitchScream() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // Shrieking agonizing death scream with convulsing frequency modulation
    const duration = 2.4;
    const osc1 = this.ctx.createOscillator();
    const osc2 = this.ctx.createOscillator();
    const mod = this.ctx.createOscillator();
    const modGain = this.ctx.createGain();
    const gain = this.ctx.createGain();
    const filter = this.ctx.createBiquadFilter();

    osc1.type = 'sawtooth';
    osc2.type = 'sine';
    mod.type = 'sawtooth';

    // Vibrato / convulsion LFO
    mod.frequency.setValueAtTime(14, now);
    mod.frequency.linearRampToValueAtTime(26, now + duration);
    modGain.gain.setValueAtTime(60, now);
    modGain.gain.linearRampToValueAtTime(140, now + duration * 0.7);

    mod.connect(modGain);
    modGain.connect(osc1.frequency);
    modGain.connect(osc2.frequency);

    // Shriek pitch descent
    osc1.frequency.setValueAtTime(1150, now);
    osc1.frequency.exponentialRampToValueAtTime(320, now + duration);
    osc2.frequency.setValueAtTime(1185, now);
    osc2.frequency.exponentialRampToValueAtTime(310, now + duration);

    filter.type = 'bandpass';
    filter.frequency.setValueAtTime(1400, now);
    filter.frequency.exponentialRampToValueAtTime(450, now + duration);
    filter.Q.setValueAtTime(2.0, now);

    gain.gain.setValueAtTime(0.01, now);
    gain.gain.linearRampToValueAtTime(0.65, now + 0.15);
    gain.gain.setValueAtTime(0.6, now + 1.2);
    gain.gain.exponentialRampToValueAtTime(0.001, now + duration);

    osc1.connect(filter);
    osc2.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterGain);

    mod.start(now);
    osc1.start(now);
    osc2.start(now);
    mod.stop(now + duration + 0.1);
    osc1.stop(now + duration + 0.1);
    osc2.stop(now + duration + 0.1);
  }

  // Painful visceral wet slash & blood splatter impact when paying blood tax on tie
  playBloodPayment() {
    if (!this.ctx || this.isMuted) return;
    const now = this.ctx.currentTime;

    // 1. Visceral Impact Thud (Pain)
    const thudOsc = this.ctx.createOscillator();
    const thudGain = this.ctx.createGain();
    thudOsc.type = 'triangle';
    thudOsc.frequency.setValueAtTime(120, now);
    thudOsc.frequency.exponentialRampToValueAtTime(32, now + 0.35);

    thudGain.gain.setValueAtTime(0.7, now);
    thudGain.gain.exponentialRampToValueAtTime(0.001, now + 0.4);

    thudOsc.connect(thudGain);
    thudGain.connect(this.masterGain);
    thudOsc.start(now);
    thudOsc.stop(now + 0.42);

    // 2. Wet Blood Splatter (Filtered Noise Burst)
    const bufferSize = Math.floor(this.ctx.sampleRate * 0.4);
    const noiseBuffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const output = noiseBuffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      output[i] = Math.random() * 2 - 1;
    }

    const whiteNoise = this.ctx.createBufferSource();
    whiteNoise.buffer = noiseBuffer;

    const noiseFilter = this.ctx.createBiquadFilter();
    noiseFilter.type = 'bandpass';
    noiseFilter.frequency.setValueAtTime(800, now);
    noiseFilter.frequency.linearRampToValueAtTime(350, now + 0.35);
    noiseFilter.Q.setValueAtTime(3.0, now);

    const noiseGain = this.ctx.createGain();
    noiseGain.gain.setValueAtTime(0.01, now);
    noiseGain.gain.linearRampToValueAtTime(0.55, now + 0.04);
    noiseGain.gain.exponentialRampToValueAtTime(0.001, now + 0.38);

    whiteNoise.connect(noiseFilter);
    noiseFilter.connect(noiseGain);
    noiseGain.connect(this.masterGain);

    whiteNoise.start(now);
    whiteNoise.stop(now + 0.4);

    // 3. Sizzling Blood on Obsidian Altar (High Frequency Sizzle)
    const sizzleBuffer = this.ctx.createBuffer(1, Math.floor(this.ctx.sampleRate * 0.6), this.ctx.sampleRate);
    const sData = sizzleBuffer.getChannelData(0);
    for (let i = 0; i < sData.length; i++) {
      sData[i] = (Math.random() * 2 - 1) * Math.exp(-i / (this.ctx.sampleRate * 0.15));
    }
    const sizzle = this.ctx.createBufferSource();
    sizzle.buffer = sizzleBuffer;
    const sFilter = this.ctx.createBiquadFilter();
    sFilter.type = 'highpass';
    sFilter.frequency.setValueAtTime(2800, now);

    const sGain = this.ctx.createGain();
    sGain.gain.setValueAtTime(0.3, now + 0.05);
    sGain.gain.exponentialRampToValueAtTime(0.001, now + 0.55);

    sizzle.connect(sFilter);
    sFilter.connect(sGain);
    sGain.connect(this.masterGain);

    sizzle.start(now + 0.05);
    sizzle.stop(now + 0.6);
  }
}

export const audioEngine = new AudioEngine();
