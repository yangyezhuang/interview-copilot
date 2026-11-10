import { ref } from 'vue'
import { AppService } from '../../bindings/interview-copilot/internal'
import { useSettings } from './useSettings'

const isListening = ref(false)

let audioCtx = null
let micStream = null
let micProcessor = null
let systemStream = null
let systemProcessor = null

export function useAudioCapture() {
  const { settings } = useSettings()

  async function toggleListening() {
    if (isListening.value) {
      stopAudioCapture()
    } else {
      await startAudioCapture()
    }
  }

  async function startAudioCapture() {
    try {
      audioCtx = new (window.AudioContext || window.webkitAudioContext)({ sampleRate: 16000 })

      // 1. Mic capture ("you")
      micStream = await navigator.mediaDevices.getUserMedia({
        audio: {
          channelCount: 1,
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true
        }
      })
      setupAudioProcessing(micStream, 'you')

      // 2. Loopback capture ("them") if enabled
      if (settings.value.meetingAudio && navigator.mediaDevices.getDisplayMedia) {
        try {
          systemStream = await navigator.mediaDevices.getDisplayMedia({
            video: true,
            audio: true
          })
          systemStream.getVideoTracks().forEach(t => t.stop())
          const audioTracks = systemStream.getAudioTracks()
          if (audioTracks.length > 0) {
            setupAudioProcessing(systemStream, 'them')
          }
        } catch (err) {
          console.warn('System loopback audio not captured:', err)
        }
      }

      isListening.value = true
    } catch (err) {
      console.error('Failed to start audio recording:', err)
      alert('麦克风启动失败，请检查系统设置中的麦克风权限！\n' + err.message)
    }
  }

  function setupAudioProcessing(stream, channel) {
    const source = audioCtx.createMediaStreamSource(stream)
    const processor = audioCtx.createScriptProcessor(4096, 1, 1)

    let pcmChunks = []
    let silentFrames = 0
    let speaking = false
    let lastFlush = Date.now()

    processor.onaudioprocess = (e) => {
      if (!isListening.value) return
      const input = e.inputBuffer.getChannelData(0)

      let sum = 0
      for (let i = 0; i < input.length; i++) {
        sum += input[i] * input[i]
      }
      const rms = Math.sqrt(sum / input.length)

      const pcm = new Int16Array(input.length)
      for (let i = 0; i < input.length; i++) {
        const s = Math.max(-1, Math.min(1, input[i]))
        pcm[i] = s < 0 ? s * 0x8000 : s * 0x7fff
      }
      pcmChunks.push(pcm)

      const THRESHOLD = 0.015
      if (rms > THRESHOLD) {
        silentFrames = 0
        speaking = true
      } else {
        if (speaking) {
          silentFrames++
        }
      }

      const elapsed = Date.now() - lastFlush
      const isPause = speaking && silentFrames > 7
      const isLongEnough = pcmChunks.length > 8
      const isMaxDuration = elapsed > 3500

      if ((isPause || isMaxDuration) && isLongEnough) {
        flushAudio(channel, pcmChunks)
        pcmChunks = []
        silentFrames = 0
        speaking = false
        lastFlush = Date.now()
      }
    }

    source.connect(processor)
    processor.connect(audioCtx.destination)

    if (channel === 'you') {
      micProcessor = processor
    } else {
      systemProcessor = processor
    }
  }

  function flushAudio(channel, chunks) {
    let totalLen = 0
    for (const c of chunks) totalLen += c.length
    const merged = new Int16Array(totalLen)
    let offset = 0
    for (const c of chunks) {
      merged.set(c, offset)
      offset += c.length
    }

    const bytes = new Uint8Array(merged.buffer)
    let binary = ''
    const len = bytes.byteLength
    for (let i = 0; i < len; i++) {
      binary += String.fromCharCode(bytes[i])
    }
    const base64 = window.btoa(binary)

    AppService.ProcessAudio(channel, base64).catch(err => {
      console.warn('Process audio error:', err)
    })
  }

  function stopAudioCapture() {
    isListening.value = false
    if (micStream) {
      micStream.getTracks().forEach(t => t.stop())
      micStream = null
    }
    if (systemStream) {
      systemStream.getTracks().forEach(t => t.stop())
      systemStream = null
    }
    if (audioCtx) {
      audioCtx.close().catch(() => {})
      audioCtx = null
    }
  }

  return {
    isListening,
    toggleListening,
    startAudioCapture,
    stopAudioCapture
  }
}
