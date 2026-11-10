import { ref, nextTick } from 'vue'
import { AppService } from '../../bindings/interview-copilot/internal'

const transcriptTurns = ref([])
const turnsBox = ref(null)

function formatTime(ts) {
  if (!ts) return ''
  const d = new Date(ts)
  return d.toTimeString().slice(0, 8)
}

async function loadHistory() {
  try {
    const turns = await AppService.GetMeetingHistory()
    if (turns) {
      transcriptTurns.value = turns
    }
  } catch (err) {
    console.error('Failed to load history:', err)
  }
}

function clearTranscript() {
  AppService.ClearTranscript()
  transcriptTurns.value = []
}

function addTurn(turn) {
  transcriptTurns.value.push(turn)
  nextTick(() => {
    if (turnsBox.value) {
      turnsBox.value.scrollTop = turnsBox.value.scrollHeight
    }
  })
}

export function useTranscript() {
  return {
    transcriptTurns,
    turnsBox,
    formatTime,
    loadHistory,
    clearTranscript,
    addTurn
  }
}
