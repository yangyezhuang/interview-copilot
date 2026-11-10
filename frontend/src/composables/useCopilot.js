import { ref, computed, nextTick } from 'vue'
import { marked } from 'marked'
import { AppService } from '../../bindings/interview-copilot/internal'
import { modeLabels, categoryLabels } from '../constants/presets'

const isBusy = ref(false)
const inputText = ref('')
const currentMode = ref('')
const detectedCategory = ref('')
const responseContent = ref('')
const thinkingContent = ref('')
const showThinking = ref(true)
const copyTip = ref('复制')
const responseBox = ref(null)
const textInput = ref(null)

const renderedMarkdown = computed(() => {
  if (!responseContent.value) return ''
  return marked.parse(responseContent.value)
})

const currentModeLabel = computed(() => {
  return modeLabels[currentMode.value] || ''
})

const detectedCategoryLabel = computed(() => {
  return categoryLabels[detectedCategory.value] || ''
})

function triggerFeature(mode) {
  currentMode.value = mode
  responseContent.value = ''
  thinkingContent.value = ''
  isBusy.value = true
  AppService.RunFeature(mode, '').catch(err => {
    isBusy.value = false
    console.error(err)
  })
}

function answerSpecificQuestion(questionText) {
  currentMode.value = 'answerThis'
  responseContent.value = ''
  thinkingContent.value = ''
  isBusy.value = true
  AppService.RunFeature('answerThis', questionText).catch(err => {
    isBusy.value = false
    console.error(err)
  })
}

function submitQuestion() {
  const text = inputText.value.trim()
  if (!text) return
  inputText.value = ''
  currentMode.value = 'ask'
  responseContent.value = ''
  thinkingContent.value = ''
  isBusy.value = true
  AppService.RunFeature('ask', text).catch(err => {
    isBusy.value = false
    console.error(err)
  })
}

function copyResponse() {
  if (!responseContent.value) return
  navigator.clipboard.writeText(responseContent.value)
  copyTip.value = '已复制!'
  setTimeout(() => {
    copyTip.value = '复制'
  }, 1500)
}

function onLlmStart(payload) {
  isBusy.value = true
  currentMode.value = payload.mode || currentMode.value
  detectedCategory.value = payload.category || ''
  responseContent.value = ''
  thinkingContent.value = ''
}

function onLlmThinking(payload) {
  if (payload.thought) {
    thinkingContent.value += payload.thought
  }
}

function onLlmToken(payload) {
  if (payload.token) {
    responseContent.value += payload.token
    nextTick(() => {
      if (responseBox.value) {
        responseBox.value.scrollTop = responseBox.value.scrollHeight
      }
    })
  }
}

function onLlmDone() {
  isBusy.value = false
}

function onLlmError(payload) {
  isBusy.value = false
  responseContent.value += `\n\n> ⚠️ **错误提示**: ${payload.message}`
}

export function useCopilot() {
  return {
    isBusy,
    inputText,
    currentMode,
    detectedCategory,
    responseContent,
    thinkingContent,
    showThinking,
    copyTip,
    responseBox,
    textInput,
    renderedMarkdown,
    currentModeLabel,
    detectedCategoryLabel,
    triggerFeature,
    answerSpecificQuestion,
    submitQuestion,
    copyResponse,
    onLlmStart,
    onLlmThinking,
    onLlmToken,
    onLlmDone,
    onLlmError
  }
}
