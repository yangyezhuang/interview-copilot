<template>
  <div class="modal-backdrop clickable" @click.self="$emit('close')">
    <div class="tutorial-modal glass-panel">
      <div class="modal-header">
        <h2>使用指南 & 隐形说明</h2>
        <button class="close-modal-btn" @click="$emit('close')">✕</button>
      </div>
      <div class="tutorial-body">
        <div class="tut-section">
          <h3>🎯 三路多模态输入</h3>
          <p><strong>1. 屏幕感知：</strong>每次触发辅助时自动截取当前屏幕，代码题秒级解析。</p>
          <p><strong>2. 麦克风输入 (You)：</strong>实时监听你的回答，防止 AI 产生重复或冲突的回复。</p>
          <p><strong>3. 会议声音 (Them)：</strong>通过系统回环监听面试官提问，自动生成应答话术。</p>
        </div>

        <div class="tut-section">
          <h3>🛡️ 屏幕共享隐形原理</h3>
          <p>本程序采用 Windows 底层 <code>SetWindowDisplayAffinity</code> API，窗口在 DWM 桌面合成层中被剔除，<strong>Zoom、腾讯会议、飞书、Google Meet、Teams</strong> 等软件共享全屏或窗口时，他人看到的都是纯净桌面，完全隐形！</p>
        </div>

        <div class="tut-section">
          <h3>⚡ 极速提词技巧</h3>
          <p>• 遇到突发问题不知如何开口：按 <code>Ctrl + Enter</code>，AI 直接给出你可以立即口述的第一人称回答。</p>
          <p>• 遇到算法代码题：打开题目页面按 <code>Ctrl + H</code>，AI 将自动分析思路、输出代码与复杂度。</p>
          <p>• 在右侧「转录对话」中，点击面试官任一条发言下方的「回答此题」，可进行单点精准突击。</p>
        </div>
      </div>
      <div class="modal-footer">
        <button class="save-btn" @click="$emit('close')">我明白了</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineEmits(['close'])
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.tutorial-modal {
  width: 90%;
  max-width: 580px;
  max-height: 85vh;
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6), 0 0 0 1px var(--border-highlight);
}

.modal-header {
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-light);
}
.modal-header h2 {
  font-size: 16px;
  color: #fff;
}
.close-modal-btn {
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 16px;
  cursor: pointer;
}
.close-modal-btn:hover {
  color: #fff;
}

.tutorial-body {
  padding: 16px 18px;
  overflow-y: auto;
  flex: 1;
}

.tut-section {
  margin-bottom: 16px;
}
.tut-section h3 {
  font-size: 14px;
  color: #60a5fa;
  margin-bottom: 6px;
}
.tut-section p {
  font-size: 12.5px;
  line-height: 1.6;
  color: var(--text-secondary);
  margin-bottom: 4px;
}
.tut-section code {
  background: rgba(255, 255, 255, 0.08);
  padding: 1px 4px;
  border-radius: 3px;
  color: #93c5fd;
  font-family: monospace;
}

.modal-footer {
  padding: 12px 18px;
  border-top: 1px solid var(--border-light);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}
.save-btn {
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  padding: 7px 16px;
  font-size: 13px;
  cursor: pointer;
  font-weight: 500;
}
.save-btn:hover {
  background: var(--accent-hover);
}
</style>
