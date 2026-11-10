<template>
  <aside class="sidebar glass-panel clickable">
    <div class="sidebar-header">
      <div class="sidebar-title">
        <span>实时对话流</span>
        <span class="turn-count">{{ transcriptTurns.length }} 条记录</span>
      </div>
      <button class="close-sidebar-btn" @click="$emit('close')">✕</button>
    </div>

    <div class="turns-list" ref="turnsBox">
      <div v-if="transcriptTurns.length === 0" class="empty-turns">
        <p>点击顶部「开启听录」即可实时记录面试官与您的语音对话。</p>
      </div>
      <div
        v-for="(turn, idx) in transcriptTurns"
        :key="idx"
        class="turn-item"
        :class="turn.channel"
      >
        <div class="turn-meta">
          <span class="speaker-tag" :class="turn.channel">
            {{ turn.channel === 'them' ? '面试官' : '我' }}
          </span>
          <span class="turn-time">{{ formatTime(turn.ts) }}</span>
        </div>
        <div class="turn-text">{{ turn.text }}</div>
        <button
          v-if="turn.channel === 'them'"
          class="answer-this-btn"
          @click="answerSpecificQuestion(turn.text)"
          title="让 AI 针对面试官这句话生成回答"
        >
          ⚡ 回答此题
        </button>
      </div>
    </div>

    <div class="sidebar-footer">
      <button class="clear-btn" @click="clearTranscript">清空对话记录</button>
    </div>
  </aside>
</template>

<script setup>
import { useTranscript } from '../composables/useTranscript'
import { useCopilot } from '../composables/useCopilot'

defineEmits(['close'])

const { transcriptTurns, turnsBox, formatTime, clearTranscript } = useTranscript()
const { answerSpecificQuestion } = useCopilot()
</script>

<style scoped>
.sidebar {
  width: 270px;
  height: 100%;
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideIn 0.2s ease;
}
@keyframes slideIn {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}

.sidebar-header {
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-light);
}
.sidebar-title {
  display: flex;
  flex-direction: column;
}
.sidebar-title span:first-child {
  font-size: 13px;
  font-weight: 600;
  color: #fff;
}
.turn-count {
  font-size: 11px;
  color: var(--text-muted);
}
.close-sidebar-btn {
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 14px;
}
.close-sidebar-btn:hover {
  color: #fff;
}

.turns-list {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.empty-turns {
  padding: 20px 10px;
  text-align: center;
  color: var(--text-muted);
  font-size: 12px;
}
.turn-item {
  padding: 8px 10px;
  border-radius: var(--radius-md);
  font-size: 12px;
  line-height: 1.45;
}
.turn-item.them {
  background: rgba(59, 130, 246, 0.12);
  border: 1px solid rgba(59, 130, 246, 0.25);
}
.turn-item.you {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.22);
}
.turn-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}
.speaker-tag {
  font-size: 10.5px;
  font-weight: 600;
}
.speaker-tag.them {
  color: #60a5fa;
}
.speaker-tag.you {
  color: #34d399;
}
.turn-time {
  font-size: 10px;
  color: var(--text-muted);
}
.turn-text {
  color: var(--text-main);
  user-select: text;
  -webkit-user-select: text;
}
.answer-this-btn {
  margin-top: 6px;
  width: 100%;
  border: 1px solid rgba(59, 130, 246, 0.35);
  background: rgba(59, 130, 246, 0.18);
  color: #93c5fd;
  border-radius: 4px;
  font-size: 11px;
  padding: 3px 6px;
  cursor: pointer;
  transition: all 0.15s;
}
.answer-this-btn:hover {
  background: rgba(59, 130, 246, 0.32);
}

.sidebar-footer {
  padding: 8px 12px;
  border-top: 1px solid var(--border-light);
}
.clear-btn {
  width: 100%;
  border: 1px solid var(--border-light);
  background: var(--bg-subtle);
  color: var(--text-muted);
  border-radius: var(--radius-sm);
  padding: 5px;
  font-size: 11.5px;
  cursor: pointer;
}
.clear-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.35);
}
</style>
