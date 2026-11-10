import { ref, computed } from 'vue'
import { AppService } from '../../bindings/interview-copilot/internal'
import { useSettings } from './useSettings'

const whisperModels = ref([])
const whisperDownloadState = ref({
  model: '',
  status: 'idle',
  downloaded: 0,
  total: 0,
  progress: 0,
  speed: '',
  error: ''
})
const whisperRuntime = ref({
  available: false,
  path: '',
  version: '',
  error: ''
})
const showRuntimeConfig = ref(false)
let downloadPollTimer = null

export function useWhisper() {
  const { settings } = useSettings()

  const currentLocalModelInfo = computed(() => {
    const name = settings.value.sttLocalModel || 'base'
    return whisperModels.value.find(m => m.name === name) || null
  })

  const isDownloadingCurrentModel = computed(() => {
    const current = settings.value.sttLocalModel || 'base'
    return (
      whisperDownloadState.value.status === 'downloading' &&
      whisperDownloadState.value.model === current
    )
  })

  const downloadProgressText = computed(() => {
    return `${(whisperDownloadState.value.progress || 0).toFixed(0)}%`
  })

  const downloadingModelTitle = computed(() => {
    const m = whisperModels.value.find(x => x.name === whisperDownloadState.value.model)
    return m ? m.displayName : (whisperDownloadState.value.model || '模型')
  })

  function formatBytes(bytes) {
    if (!bytes || bytes <= 0) return '0 B'
    const units = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
  }

  async function loadWhisperModels() {
    try {
      const list = await AppService.GetWhisperModels()
      if (list && list.length) {
        whisperModels.value = list
      }
    } catch (err) {
      console.error('Failed to load whisper models:', err)
    }
  }

  async function detectRuntime() {
    try {
      const res = await AppService.CheckWhisperRuntime()
      if (res) {
        whisperRuntime.value = res
      }
    } catch (err) {
      whisperRuntime.value = { available: false, error: err.message }
    }
  }

  function startDownloadPolling() {
    if (downloadPollTimer) clearInterval(downloadPollTimer)
    downloadPollTimer = setInterval(async () => {
      try {
        const state = await AppService.GetWhisperModelDownloadState()
        if (state) {
          whisperDownloadState.value = state
          if (state.status === 'completed' || state.status === 'failed' || state.status === 'cancelled') {
            clearInterval(downloadPollTimer)
            downloadPollTimer = null
            await loadWhisperModels()
          }
        }
      } catch (e) {}
    }, 400)
  }

  async function startDownloadCurrentModel() {
    const name = settings.value.sttLocalModel || 'base'
    try {
      // 确保当前选中的镜像源等配置已保存
      try {
        await AppService.SaveSettings(settings.value)
      } catch (saveErr) {
        console.warn('Auto-save before download warning:', saveErr)
      }
      whisperDownloadState.value = {
        model: name,
        status: 'downloading',
        downloaded: 0,
        total: currentLocalModelInfo.value?.sizeBytes || 0,
        progress: 0,
        speed: '连接中...',
        error: ''
      }
      await AppService.DownloadWhisperModel(name)
      startDownloadPolling()
    } catch (err) {
      whisperDownloadState.value = {
        model: name,
        status: 'failed',
        error: err.message || '启动下载失败'
      }
    }
  }

  async function cancelModelDownload() {
    try {
      const model = whisperDownloadState.value.model || (settings.value.sttLocalModel || 'base')
      await AppService.CancelWhisperModelDownload(model)
      whisperDownloadState.value.status = 'cancelled'
      if (downloadPollTimer) {
        clearInterval(downloadPollTimer)
        downloadPollTimer = null
      }
      await loadWhisperModels()
    } catch (err) {
      console.error('Cancel download error:', err)
    }
  }

  async function deleteCurrentModel() {
    const name = settings.value.sttLocalModel || 'base'
    if (!confirm(`确定要删除本地 Whisper 模型 ${name} 吗？`)) return
    try {
      await AppService.DeleteWhisperModel(name)
      await loadWhisperModels()
    } catch (err) {
      alert('删除模型失败: ' + err.message)
    }
  }

  return {
    whisperModels,
    whisperDownloadState,
    whisperRuntime,
    showRuntimeConfig,
    currentLocalModelInfo,
    isDownloadingCurrentModel,
    downloadProgressText,
    downloadingModelTitle,
    formatBytes,
    loadWhisperModels,
    detectRuntime,
    startDownloadPolling,
    startDownloadCurrentModel,
    cancelModelDownload,
    deleteCurrentModel
  }
}
