<template>
  <div class="main-card glass-panel clickable">
    <!-- Status / Header Bar -->
    <div class="card-header">
      <div class="header-tags">
        <span class="badge mode-badge" v-if="currentModeLabel">
          {{ currentModeLabel }}
        </span>
        <span class="badge category-badge" v-if="detectedCategoryLabel">
          {{ detectedCategoryLabel }}
        </span>
        <span class="badge model-badge" :class="{ smart: settings.smart }">
          {{ settings.smart ? '🧠 深度思考' : '⚡ 极速模式' }}
        </span>
      </div>

      <div class="header-actions">
        <button
          v-if="responseContent"
          class="icon-action-btn"
          @click="copyResponse"
          title="复制回答"
        >
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
          <span>{{ copyTip }}</span>
        </button>
      </div>
    </div>

    <!-- Thinking Area (for DeepSeek-R1 / o1 reasoning) -->
    <div v-if="thinkingContent" class="thinking-box">
      <div class="thinking-header" @click="showThinking = !showThinking">
        <span class="spin-dot" :class="{ spinning: isBusy }"></span>
        <span>思考过程 (Reasoning)...</span>
        <span class="toggle-arrow">{{ showThinking ? '▲' : '▼' }}</span>
      </div>
      <div v-show="showThinking" class="thinking-content">
        {{ thinkingContent }}
      </div>
    </div>

    <!-- Response Body -->
    <div class="response-area" ref="responseBox">
      <div v-if="responseContent" class="markdown-body" v-html="renderedMarkdown"></div>
      <div v-else class="welcome-guide">
        <h3 class="welcome-title">面试隐形助手已就绪</h3>
        <p class="welcome-desc">已开启屏幕共享隐形保护，他人无法在共享或录屏中看到此窗口。</p>
        <div class="shortcut-grid">
          <div class="sc-item" @click="triggerFeature('say')">
            <span class="sc-key">Ctrl + Enter</span>
            <span class="sc-desc"><strong>我该怎么说</strong>：实时提词建议</span>
          </div>
          <div class="sc-item" @click="triggerFeature('assist')">
            <span class="sc-key">Ctrl + Shift + Enter</span>
            <span class="sc-desc"><strong>智能辅助</strong>：结合屏幕与上下文解答</span>
          </div>
          <div class="sc-item" @click="triggerFeature('leetcode')">
            <span class="sc-key">Ctrl + H</span>
            <span class="sc-desc"><strong>算法解题</strong>：识别屏幕题目并编写代码</span>
          </div>
          <div class="sc-item" @click="triggerFeature('recap')">
            <span class="sc-key">复盘按钮</span>
            <span class="sc-desc"><strong>会议总结</strong>：结构化梳理全场对话</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Row -->
    <div class="action-row">
      <button class="act-btn primary" @click="triggerFeature('say')" title="Ctrl+Enter">
        <span class="icon-spark">💬</span>
        <span>我该怎么说?</span>
        <span class="key-hint">Ctrl+↵</span>
      </button>
      <button class="act-btn secondary" @click="triggerFeature('assist')" title="Ctrl+Shift+Enter">
        <span class="icon-spark">✨</span>
        <span>智能辅助</span>
        <span class="key-hint">Ctrl+⇧+↵</span>
      </button>
      <button class="act-btn" @click="triggerFeature('leetcode')" title="Ctrl+H">
        <span class="icon-spark">💻</span>
        <span>算法解题</span>
        <span class="key-hint">Ctrl+H</span>
      </button>
      <button class="act-btn" @click="triggerFeature('recap')">
        <span class="icon-spark">📋</span>
        <span>复盘总结</span>
      </button>
    </div>

    <!-- Composer Area -->
    <div class="composer">
      <div class="input-wrap">
        <textarea
          v-model="inputText"
          placeholder="输入你的问题，或按 Ctrl+Shift+Enter 截屏智能解答..."
          rows="1"
          @keydown.enter.exact.prevent="submitQuestion"
          @keydown.shift.enter.stop
          ref="textInput"
        ></textarea>
      </div>
      <div class="composer-toolbar">
        <button
          class="pill-btn"
          :class="{ active: settings.smart }"
          @click="toggleSmartMode"
          title="切换思考大模型"
        >
          <span>{{ settings.smart ? '🧠 Smart 开' : '⚡ Fast 模式' }}</span>
        </button>

        <button
          class="pill-btn sidebar-btn"
          :class="{ active: showSidebar }"
          @click="$emit('toggle-sidebar')"
          title="转录历史对话"
        >
          <span>转录对话 ({{ transcriptTurns.length }})</span>
        </button>

        <div class="flex-spacer"></div>

        <button
          class="send-btn"
          :disabled="!inputText.trim() && !isBusy"
          @click="submitQuestion"
          title="发送"
        >
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="22" y1="2" x2="11" y2="13"></line>
            <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useSettings } from '../composables/useSettings'
import { useCopilot } from '../composables/useCopilot'
import { useTranscript } from '../composables/useTranscript'

defineProps({
  showSidebar: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-sidebar'])

const { settings, toggleSmartMode } = useSettings()
const {
  isBusy,
  inputText,
  currentModeLabel,
  detectedCategoryLabel,
  responseContent,
  thinkingContent,
  showThinking,
  copyTip,
  responseBox,
  textInput,
  renderedMarkdown,
  triggerFeature,
  submitQuestion,
  copyResponse
} = useCopilot()
const { transcriptTurns } = useTranscript()
</script>

<style scoped>
.main-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: var(--radius-lg);
  overflow: hidden;
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border-light);
  background: rgba(255, 255, 255, 0.02);
}
.header-tags {
  display: flex;
  gap: 6px;
  align-items: center;
}
.badge {
  font-size: 11px;
  padding: 3px 8px;
  border-radius: var(--radius-full);
  font-weight: 500;
}
.badge.mode-badge {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.3);
}
.badge.category-badge {
  background: rgba(168, 85, 247, 0.2);
  color: #c084fc;
  border: 1px solid rgba(168, 85, 247, 0.3);
}
.badge.model-badge {
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-secondary);
}
.badge.model-badge.smart {
  background: rgba(234, 179, 8, 0.18);
  color: #facc15;
  border: 1px solid rgba(234, 179, 8, 0.3);
}

.icon-action-btn {
  border: none;
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-secondary);
  padding: 4px 9px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11.5px;
  cursor: pointer;
  transition: all 0.15s;
}
.icon-action-btn:hover {
  background: rgba(255, 255, 255, 0.14);
  color: #fff;
}
.icon-action-btn .icon {
  width: 13px;
  height: 13px;
}

/* Thinking Box */
.thinking-box {
  background: rgba(15, 23, 42, 0.6);
  border-bottom: 1px solid var(--border-light);
  padding: 8px 16px;
  font-size: 12px;
}
.thinking-header {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #94a3b8;
  cursor: pointer;
  font-weight: 500;
}
.spin-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  border: 2px solid #38bdf8;
  border-top-color: transparent;
}
.spin-dot.spinning {
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.toggle-arrow {
  margin-left: auto;
  font-size: 9px;
}
.thinking-content {
  margin-top: 6px;
  max-height: 100px;
  overflow-y: auto;
  color: #cbd5e1;
  line-height: 1.5;
  white-space: pre-wrap;
  font-family: var(--font-mono);
  font-size: 11.5px;
}

/* Response Area */
.response-area {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  user-select: text;
  -webkit-user-select: text;
}
.welcome-guide {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  padding: 20px;
}
.welcome-title {
  font-size: 17px;
  color: #fff;
  margin-bottom: 6px;
}
.welcome-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 20px;
}
.shortcut-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  width: 100%;
  max-width: 480px;
}
.sc-item {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s;
}
.sc-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--accent);
}
.sc-key {
  font-family: var(--font-mono);
  font-size: 11px;
  color: #60a5fa;
  font-weight: 600;
}
.sc-desc {
  font-size: 12px;
  color: var(--text-secondary);
}

/* Action Row */
.action-row {
  display: flex;
  gap: 6px;
  padding: 8px 14px;
  border-top: 1px solid var(--border-light);
  background: rgba(255, 255, 255, 0.02);
}
.act-btn {
  flex: 1;
  border: 1px solid var(--border-light);
  background: var(--bg-subtle);
  color: var(--text-main);
  padding: 7px 10px;
  border-radius: var(--radius-md);
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.15s;
}
.act-btn:hover {
  background: var(--bg-hover);
  border-color: rgba(255, 255, 255, 0.25);
}
.act-btn.primary {
  background: rgba(59, 130, 246, 0.15);
  border-color: rgba(59, 130, 246, 0.35);
  color: #93c5fd;
}
.act-btn.primary:hover {
  background: rgba(59, 130, 246, 0.28);
}
.act-btn.secondary {
  background: rgba(168, 85, 247, 0.14);
  border-color: rgba(168, 85, 247, 0.35);
  color: #d8b4fe;
}
.act-btn.secondary:hover {
  background: rgba(168, 85, 247, 0.25);
}
.key-hint {
  font-size: 10px;
  opacity: 0.6;
  font-family: var(--font-mono);
}

/* Composer */
.composer {
  padding: 8px 14px 12px;
  background: rgba(0, 0, 0, 0.25);
  border-top: 1px solid var(--border-light);
}
.input-wrap textarea {
  width: 100%;
  border: none;
  background: transparent;
  color: #fff;
  font-family: inherit;
  font-size: 13.5px;
  line-height: 1.5;
  resize: none;
  outline: none;
}
.input-wrap textarea::placeholder {
  color: var(--text-muted);
}
.composer-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
}
.pill-btn {
  border: 1px solid var(--border-light);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-secondary);
  border-radius: var(--radius-full);
  padding: 4px 10px;
  font-size: 11.5px;
  cursor: pointer;
  transition: all 0.15s;
}
.pill-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}
.pill-btn.active {
  background: rgba(234, 179, 8, 0.2);
  color: #facc15;
  border-color: rgba(234, 179, 8, 0.4);
}
.pill-btn.sidebar-btn.active {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border-color: rgba(59, 130, 246, 0.4);
}
.flex-spacer {
  flex: 1;
}
.send-btn {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: none;
  background: var(--accent);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s;
}
.send-btn:hover:not(:disabled) {
  background: var(--accent-hover);
  box-shadow: 0 0 10px var(--accent-glow);
}
.send-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}
.send-btn .icon {
  width: 14px;
  height: 14px;
}
</style>
