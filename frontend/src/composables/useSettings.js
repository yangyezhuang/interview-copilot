import { ref, computed } from 'vue'
import { AppService } from '../../bindings/interview-copilot/internal'
import {
  defaultBaseUrls,
  defaultApiKeys,
  defaultModels,
  providerPresets
} from '../constants/presets'

const settings = ref({
  provider: 'deepseek',
  baseUrl: '',
  baseUrls: { ...defaultBaseUrls },
  apiKeys: { ...defaultApiKeys },
  models: { ...defaultModels },
  smart: false,
  sttProvider: 'local',
  sttLocalModel: 'base',
  sttMirrorSource: 'hf-mirror',
  sttCustomMirrorUrl: '',
  sttWhisperPath: '',
  sttModel: 'whisper-large-v3-turbo',
  sttBaseUrl: '',
  sttKey: '',
  sttLanguage: 'zh',
  meetingAudio: true,
  opacity: 0.96,
  screenProtection: true,
  resumeText: '',
  jobDescription: '',
  starStories: '',
  whyCompany: '',
  whyLeaving: '',
  workStyle: '',
  salaryTarget: '',
  questionsToAsk: '',
  aiRules: ''
})

const saveMessage = ref('')

const currentProviderMeta = computed(() => {
  return providerPresets[settings.value.provider] || null
})

const currentApiKey = computed({
  get() {
    if (!settings.value.apiKeys) settings.value.apiKeys = {}
    return settings.value.apiKeys[settings.value.provider] || ''
  },
  set(val) {
    if (!settings.value.apiKeys) settings.value.apiKeys = {}
    settings.value.apiKeys[settings.value.provider] = val
  }
})

const currentBaseUrl = computed({
  get() {
    if (!settings.value.baseUrls) settings.value.baseUrls = {}
    const p = settings.value.provider
    if (settings.value.baseUrls[p]) return settings.value.baseUrls[p]
    return providerPresets[p]?.url || ''
  },
  set(val) {
    if (!settings.value.baseUrls) settings.value.baseUrls = {}
    settings.value.baseUrls[settings.value.provider] = val
    if (settings.value.provider === 'custom') {
      settings.value.baseUrl = val
    }
  }
})

const currentModelPair = computed({
  get() {
    if (!settings.value.models) settings.value.models = {}
    const p = settings.value.provider
    if (!settings.value.models[p]) {
      settings.value.models[p] = defaultModels[p] || { fast: '', smart: '' }
    }
    return settings.value.models[p]
  },
  set(val) {
    if (!settings.value.models) settings.value.models = {}
    settings.value.models[settings.value.provider] = val
  }
})

const sttKeyPlaceholder = computed(() => {
  if (settings.value.sttProvider === 'siliconflow') {
    return settings.value.apiKeys?.siliconflow ? '已配置 (留空复用硅基流动 Key)' : '输入 SiliconFlow API Key'
  }
  if (settings.value.sttProvider === 'groq') {
    return settings.value.apiKeys?.groq ? '已配置 (留空复用 Groq Key)' : '输入 Groq API Key'
  }
  if (settings.value.sttProvider === 'openai') {
    return settings.value.apiKeys?.openai ? '已配置 (留空复用 OpenAI Key)' : '输入 OpenAI API Key'
  }
  return '输入对应的 STT API Key'
})

function resetCurrentBaseUrl() {
  const p = settings.value.provider
  if (providerPresets[p]?.url) {
    if (!settings.value.baseUrls) settings.value.baseUrls = {}
    settings.value.baseUrls[p] = providerPresets[p].url
    if (p === 'custom') settings.value.baseUrl = providerPresets[p].url
  }
}

async function loadSettings() {
  try {
    const s = await AppService.GetSettings()
    if (s) {
      if (!s.baseUrls) s.baseUrls = {}
      if (!s.apiKeys) s.apiKeys = {}
      if (!s.models) s.models = {}
      for (const [k, v] of Object.entries(defaultBaseUrls)) {
        if (!s.baseUrls[k]) s.baseUrls[k] = v
      }
      for (const [k, v] of Object.entries(defaultApiKeys)) {
        if (s.apiKeys[k] === undefined) s.apiKeys[k] = v
      }
      for (const [k, v] of Object.entries(defaultModels)) {
        if (!s.models[k]) s.models[k] = { ...v }
      }
      if (!s.sttMirrorSource) s.sttMirrorSource = 'hf-mirror'
      if (s.sttCustomMirrorUrl === undefined) s.sttCustomMirrorUrl = ''
      settings.value = s
    }
  } catch (err) {
    console.error('Failed to load settings:', err)
  }
}

async function saveSettings() {
  try {
    await AppService.SaveSettings(settings.value)
    saveMessage.value = '✓ 已保存'
    setTimeout(() => {
      saveMessage.value = ''
    }, 1500)
    return true
  } catch (err) {
    saveMessage.value = '保存失败: ' + err.message
    return false
  }
}

async function saveSettingsSilently() {
  try {
    await AppService.SaveSettings(settings.value)
  } catch (err) {
    console.error('Failed to save settings silently:', err)
  }
}

function toggleSmartMode() {
  settings.value.smart = !settings.value.smart
  saveSettingsSilently()
}

export function useSettings() {
  return {
    settings,
    saveMessage,
    currentProviderMeta,
    currentApiKey,
    currentBaseUrl,
    currentModelPair,
    sttKeyPlaceholder,
    resetCurrentBaseUrl,
    loadSettings,
    saveSettings,
    saveSettingsSilently,
    toggleSmartMode
  }
}
