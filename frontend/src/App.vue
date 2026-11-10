<template>
  <div id="main-container" :style="{ opacity: settings.opacity }">
    <!-- Top Pill Toolbar -->
    <TopToolbar
      :is-collapsed="isCollapsed"
      @open-tutorial="showTutorial = true"
      @open-settings="showSettingsModal = true"
      @toggle-collapse="isCollapsed = !isCollapsed"
    />

    <!-- Main Content Area -->
    <main v-show="!isCollapsed" class="panel-layout">
      <!-- Main Copilot Card -->
      <MainCard
        :show-sidebar="showSidebar"
        @toggle-sidebar="showSidebar = !showSidebar"
      />

      <!-- Sliding Transcript History Drawer -->
      <TranscriptDrawer
        v-if="showSidebar"
        @close="showSidebar = false"
      />
    </main>

    <!-- Settings Modal -->
    <SettingsModal
      v-if="showSettingsModal"
      @close="showSettingsModal = false"
    />

    <!-- Tutorial Modal -->
    <TutorialModal
      v-if="showTutorial"
      @close="showTutorial = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Events } from '@wailsio/runtime'
import { useSettings } from './composables/useSettings'
import { useTranscript } from './composables/useTranscript'
import { useWhisper } from './composables/useWhisper'
import { useCopilot } from './composables/useCopilot'

import TopToolbar from './components/TopToolbar.vue'
import MainCard from './components/MainCard.vue'
import TranscriptDrawer from './components/TranscriptDrawer.vue'
import SettingsModal from './components/SettingsModal.vue'
import TutorialModal from './components/TutorialModal.vue'

const { settings, loadSettings } = useSettings()
const { loadHistory, addTurn, clearTranscript } = useTranscript()
const { loadWhisperModels, whisperDownloadState } = useWhisper()
const { onLlmStart, onLlmThinking, onLlmToken, onLlmDone, onLlmError } = useCopilot()

const isCollapsed = ref(false)
const showSidebar = ref(false)
const showSettingsModal = ref(false)
const showTutorial = ref(false)

onMounted(async () => {
  await loadSettings()
  await loadHistory()
  await loadWhisperModels()

  // Wails Event Listeners
  Events.On('stt:model_download_progress', (data) => {
    const payload = data.data || data
    whisperDownloadState.value = payload
    if (payload.status === 'completed' || payload.status === 'failed') {
      loadWhisperModels()
    }
  })

  Events.On('llm:start', (data) => onLlmStart(data.data || data))
  Events.On('llm:thinking', (data) => onLlmThinking(data.data || data))
  Events.On('llm:token', (data) => onLlmToken(data.data || data))
  Events.On('llm:done', () => onLlmDone())
  Events.On('llm:error', (data) => onLlmError(data.data || data))

  Events.On('transcript:turn', (data) => addTurn(data.data || data))
  Events.On('transcript:cleared', () => clearTranscript())

  Events.On('window:toggle_hide', () => {
    isCollapsed.value = !isCollapsed.value
  })
})
</script>

<style scoped>
#main-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  height: 100%;
  padding: 10px;
  position: relative;
  transition: opacity 0.2s ease;
}

.panel-layout {
  display: flex;
  gap: 12px;
  width: 100%;
  max-width: 740px;
  height: calc(100% - 50px);
  position: relative;
}
</style>
