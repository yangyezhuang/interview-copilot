<template>
  <header class="toolbar glass-panel clickable">
    <!-- Drag Handle -->
    <div class="drag-pill" style="--wails-draggable: drag">
      <span class="drag-dots"></span>
      <span class="app-title">Copilot</span>
    </div>

    <!-- Live Listening Indicator -->
    <span
      class="status-indicator"
      :class="{ active: isListening, busy: isBusy }"
      :title="isListening ? '正在听录会议中...' : '听录未开启'"
    ></span>

    <!-- Action Buttons -->
    <button class="tb-btn" @click="$emit('open-tutorial')" title="使用教程与说明">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
        <line x1="12" y1="17" x2="12.01" y2="17"></line>
      </svg>
    </button>

    <button class="tb-btn" @click="$emit('open-settings')" title="设置中心">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="3"></circle>
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
      </svg>
    </button>

    <div class="tb-divider"></div>

    <!-- Toggle Hide / Collapse -->
    <button class="tb-btn text-btn" @click="$emit('toggle-collapse')" title="折叠/展开面板 (Ctrl+Shift+/)">
      <svg class="icon chev" :class="{ rotated: isCollapsed }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="18 15 12 9 6 15"></polyline>
      </svg>
      <span>{{ isCollapsed ? '展开' : '折叠' }}</span>
    </button>

    <div class="tb-divider"></div>

    <!-- Toggle Session / Listening -->
    <button
      class="tb-btn session-btn"
      :class="{ active: isListening }"
      @click="toggleListening"
      :title="isListening ? '停止听录' : '启动听录 (采集麦克风与会议声音)'"
    >
      <span class="mic-dot" :class="{ 'pulsing-live': isListening }"></span>
      <span>{{ isListening ? '听录中' : '开启听录' }}</span>
    </button>

    <div class="tb-divider"></div>

    <!-- Opacity Popover Toggle -->
    <div class="opacity-container">
      <button class="tb-btn" @click="showOpacityPopover = !showOpacityPopover" title="透明度调节">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M12 2a10 10 0 0 1 0 20z" fill="currentColor"></path>
        </svg>
      </button>
      <div v-if="showOpacityPopover" class="opacity-popover glass-panel">
        <span class="pop-title">透明度: {{ Math.round(settings.opacity * 100) }}%</span>
        <input
          type="range"
          min="0.2"
          max="1.0"
          step="0.02"
          v-model.number="settings.opacity"
          @change="saveSettingsSilently"
        />
      </div>
    </div>

    <div class="tb-divider"></div>

    <!-- Close App -->
    <button class="tb-btn close-btn" @click="quitApp" title="退出程序 (Ctrl+Shift+X)">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { AppService } from '../../bindings/interview-copilot/internal'
import { useSettings } from '../composables/useSettings'
import { useAudioCapture } from '../composables/useAudioCapture'
import { useCopilot } from '../composables/useCopilot'

defineProps({
  isCollapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['open-tutorial', 'open-settings', 'toggle-collapse'])

const { settings, saveSettingsSilently } = useSettings()
const { isListening, toggleListening } = useAudioCapture()
const { isBusy } = useCopilot()

const showOpacityPopover = ref(false)

function quitApp() {
  AppService.QuitApp()
}
</script>

<style scoped>
.toolbar {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: var(--radius-full);
  margin-bottom: 8px;
  z-index: 50;
}

.drag-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 8px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-full);
  cursor: grab;
}
.drag-pill:active {
  cursor: grabbing;
}
.drag-dots {
  width: 12px;
  height: 8px;
  background-image: radial-gradient(circle, rgba(255, 255, 255, 0.9) 1px, transparent 1.2px);
  background-size: 4px 4px;
}
.app-title {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.5px;
  color: #fff;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  transition: all 0.3s;
}
.status-indicator.active {
  background: var(--live-green);
  box-shadow: 0 0 8px var(--live-green);
}
.status-indicator.busy {
  background: #f59e0b;
  box-shadow: 0 0 8px #f59e0b;
}

.tb-btn {
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-full);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 5px 7px;
  font-size: 12px;
  transition: all 0.15s ease;
}
.tb-btn:hover {
  background: var(--bg-hover);
  color: #fff;
}

.tb-btn .icon {
  width: 15px;
  height: 15px;
}
.tb-btn.text-btn {
  padding: 4px 10px;
  gap: 4px;
}
.tb-btn.session-btn {
  padding: 4px 12px;
  gap: 6px;
  font-weight: 500;
}
.tb-btn.session-btn.active {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}
.mic-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-muted);
}
.mic-dot.pulsing-live {
  background: var(--live-green);
}

.tb-divider {
  width: 1px;
  height: 16px;
  background: rgba(255, 255, 255, 0.14);
  margin: 0 2px;
}

.tb-btn.close-btn:hover {
  background: rgba(239, 68, 68, 0.25);
  color: #ef4444;
}

.chev {
  transition: transform 0.2s ease;
}
.chev.rotated {
  transform: rotate(180deg);
}

.opacity-container {
  position: relative;
}
.opacity-popover {
  position: absolute;
  top: 110%;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 14px;
  border-radius: var(--radius-md);
  min-width: 170px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  z-index: 100;
}
.pop-title {
  font-size: 12px;
  color: var(--text-secondary);
}
.opacity-popover input[type='range'] {
  accent-color: var(--accent);
  cursor: pointer;
}
</style>
