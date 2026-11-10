<template>
  <div class="modal-backdrop clickable" @click.self="$emit('close')">
    <div class="settings-modal glass-panel">
      <div class="modal-header">
        <h2>设置中心</h2>
        <button class="close-modal-btn" @click="$emit('close')">✕</button>
      </div>

      <!-- Tabs -->
      <div class="settings-tabs">
        <button :class="{ active: activeTab === 'llm' }" @click="activeTab = 'llm'">🤖 模型与服务</button>
        <button :class="{ active: activeTab === 'stt' }" @click="activeTab = 'stt'">🎙️ 语音识别</button>
        <button :class="{ active: activeTab === 'profile' }" @click="activeTab = 'profile'">👤 面试背景画像</button>
        <button :class="{ active: activeTab === 'appearance' }" @click="activeTab = 'appearance'">🎨 界面与快捷键</button>
      </div>

      <div class="settings-body">
        <!-- LLM Tab -->
        <div v-show="activeTab === 'llm'" class="tab-pane">
          <div class="form-group">
            <label>大模型服务商</label>
            <select v-model="settings.provider">
              <optgroup label="国内大模型">
                <option value="deepseek">DeepSeek (深度求索 - 官方)</option>
                <option value="qwen">通义千问 (阿里百炼 DashScope)</option>
                <option value="doubao">豆包 (字节跳动 火山方舟 Ark)</option>
                <option value="siliconflow">硅基流动 (SiliconFlow - 极速聚合)</option>
                <option value="moonshot">Kimi (月之暗面 Moonshot)</option>
                <option value="zhipu">智谱清言 (GLM-4)</option>
              </optgroup>
              <optgroup label="海外大模型">
                <option value="openai">OpenAI (GPT-4o / o1)</option>
                <option value="anthropic">Claude (Anthropic)</option>
                <option value="gemini">Google Gemini</option>
                <option value="groq">Groq (超高速推理)</option>
              </optgroup>
              <optgroup label="本地私有化与中转端点">
                <option value="ollama">Ollama (本地私有部署)</option>
                <option value="custom">自定义 OpenAI 兼容接口 (OneAPI/vLLM等)</option>
              </optgroup>
            </select>
            <p class="hint provider-doc" v-if="currentProviderMeta">
              💡 {{ currentProviderMeta.doc }}
            </p>
          </div>

          <div class="form-group">
            <label>API Key</label>
            <input
              type="password"
              v-model="currentApiKey"
              :placeholder="currentProviderMeta ? currentProviderMeta.keyHint : 'sk-...'"
              autocomplete="off"
            />
          </div>

          <div class="form-group" v-if="!['anthropic', 'gemini'].includes(settings.provider)">
            <div class="label-with-action">
              <label>API Base URL (接口端点)</label>
              <button
                type="button"
                class="text-link-btn"
                @click="resetCurrentBaseUrl"
                title="恢复为该服务商的官方标准接口地址"
              >
                恢复官方默认
              </button>
            </div>
            <input
              type="text"
              v-model="currentBaseUrl"
              placeholder="例如: https://dashscope.aliyuncs.com/compatible-mode/v1"
            />
          </div>

          <div class="form-row">
            <div class="form-group flex-1">
              <label>极速模型 (Fast Model)</label>
              <input type="text" v-model="currentModelPair.fast" />
              <div class="model-quick-select" v-if="currentProviderMeta && currentProviderMeta.fastSuggestions">
                <span class="sub-label">推荐:</span>
                <span
                  v-for="m in currentProviderMeta.fastSuggestions"
                  :key="m"
                  class="model-tag"
                  @click="currentModelPair.fast = m"
                >
                  {{ m }}
                </span>
              </div>
            </div>

            <div class="form-group flex-1">
              <label>深度思考模型 (Smart Model)</label>
              <input type="text" v-model="currentModelPair.smart" />
              <div class="model-quick-select" v-if="currentProviderMeta && currentProviderMeta.smartSuggestions">
                <span class="sub-label">推荐:</span>
                <span
                  v-for="m in currentProviderMeta.smartSuggestions"
                  :key="m"
                  class="model-tag"
                  @click="currentModelPair.smart = m"
                >
                  {{ m }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- STT Tab -->
        <div v-show="activeTab === 'stt'" class="tab-pane">
          <div class="form-group">
            <label>语音识别引擎 / 服务商</label>
            <select v-model="settings.sttProvider">
              <optgroup label="本地离线引擎 (推荐 - 零费用/私密)">
                <option value="local">本地 Whisper (whisper.cpp - 离线运行)</option>
              </optgroup>
              <optgroup label="云端极速服务 (API)">
                <option value="groq">Groq Whisper (极速低延迟 - 强烈推荐)</option>
                <option value="siliconflow">硅基流动 SiliconFlow (国内免翻直连 / SenseVoice)</option>
                <option value="openai">OpenAI Whisper (官方 whisper-1)</option>
                <option value="custom">自定义 STT 端点 (兼容 OpenAI 接口)</option>
              </optgroup>
            </select>
          </div>

          <!-- 本地 Whisper 专属配置 -->
          <div v-if="settings.sttProvider === 'local'" class="local-stt-panel">
            <div class="form-group">
              <div class="label-with-badge">
                <label>本地 Whisper 模型 (GGML)</label>
                <span v-if="currentLocalModelInfo?.downloaded" class="model-badge badge-ready">✓ 已就绪</span>
                <span v-else-if="isDownloadingCurrentModel" class="model-badge badge-downloading">⏳ 下载中 {{ downloadProgressText }}</span>
                <span v-else class="model-badge badge-missing">⚠ 需下载 ({{ currentLocalModelInfo?.size || '未知' }})</span>
              </div>
              <select v-model="settings.sttLocalModel">
                <option v-for="m in whisperModels" :key="m.name" :value="m.name">
                  {{ m.displayName }} [{{ m.size }}] {{ m.downloaded ? '✓ 已就绪' : '— 需下载' }} ({{ m.recommendation || (m.speed >= 4 ? '高速' : '高精') }})
                </option>
              </select>
              <p class="hint" v-if="currentLocalModelInfo?.description">
                💡 {{ currentLocalModelInfo.description }}
              </p>
            </div>

            <!-- 模型下载镜像源配置 -->
            <div class="form-group">
              <label>模型下载镜像源</label>
              <select v-model="settings.sttMirrorSource">
                <option value="hf-mirror">国内高速镜像 (hf-mirror.com)</option>
                <option value="huggingface">HuggingFace 官方源 (huggingface.co)</option>
                <option value="custom">自定义镜像源 URL (高级)</option>
              </select>
              <div v-if="settings.sttMirrorSource === 'custom'" class="custom-mirror-box" style="margin-top: 8px;">
                <input
                  type="text"
                  v-model="settings.sttCustomMirrorUrl"
                  placeholder="输入镜像源根地址，例如: https://my-server.com/whisper-models"
                />
                <p class="hint">
                  下载时将自动在末尾拼接文件名: <code>{{ (settings.sttCustomMirrorUrl || 'https://...').replace(/\/+$/, '') }}/{{ currentLocalModelInfo?.filename || 'ggml-base.bin' }}</code>
                </p>
              </div>
            </div>

            <!-- 模型下载与就绪状态卡片 -->
            <div class="model-status-card">
              <!-- 状态 1: 正在下载 -->
              <div v-if="whisperDownloadState.status === 'downloading'" class="download-progress-box">
                <div class="download-info-row">
                  <span class="downloading-title">
                    <span class="spinner-small"></span> 正在下载 {{ downloadingModelTitle }}...
                  </span>
                  <span class="download-speed-badge">{{ whisperDownloadState.speed || '正在连接...' }}</span>
                </div>
                <div class="progress-bar-container">
                  <div class="progress-bar-fill" :style="{ width: whisperDownloadState.progress + '%' }"></div>
                </div>
                <div class="download-meta-row">
                  <span>{{ formatBytes(whisperDownloadState.downloaded) }} / {{ formatBytes(whisperDownloadState.total) }} ({{ whisperDownloadState.progress.toFixed(1) }}%)</span>
                  <button type="button" class="btn-cancel-download" @click="cancelModelDownload">取消下载</button>
                </div>
              </div>

              <!-- 状态 2: 下载失败 -->
              <div v-else-if="whisperDownloadState.status === 'failed' && whisperDownloadState.model === (settings.sttLocalModel || 'base')" class="download-error-box">
                <div class="error-msg-row">
                  <span>⚠️ 下载失败: {{ whisperDownloadState.error }}</span>
                </div>
                <button type="button" class="btn-download-action" @click="startDownloadCurrentModel">重新下载</button>
              </div>

              <!-- 状态 3: 模型未下载，提示下载 -->
              <div v-else-if="!currentLocalModelInfo?.downloaded" class="download-prompt-box">
                <div class="prompt-text">
                  <span class="prompt-icon">📦</span>
                  <div class="prompt-desc">
                    <strong>需要下载该模型后方可使用</strong>
                    <p>模型文件 <code>{{ currentLocalModelInfo?.filename || 'ggml-model.bin' }}</code> 尚未下载（大小约 {{ currentLocalModelInfo?.size }}）。点击下方按钮即可一键下载到本地。</p>
                  </div>
                </div>
                <button type="button" class="btn-download-action" @click="startDownloadCurrentModel">
                  ⬇️ 立即下载该模型 ({{ currentLocalModelInfo?.size }})
                </button>
              </div>

              <!-- 状态 4: 模型已就绪 -->
              <div v-else class="download-ready-box">
                <div class="ready-info">
                  <span class="ready-icon">✓</span>
                  <div>
                    <strong>本地模型已就绪</strong>
                    <p>模型文件 <code>{{ currentLocalModelInfo?.filename }}</code> 存储于本地，转录将完全离线运行。</p>
                  </div>
                </div>
                <button type="button" class="btn-text-danger" title="删除本地模型文件以释放空间" @click="deleteCurrentModel">删除模型</button>
              </div>
            </div>

            <!-- Whisper CLI Runtime 运行环境检测与指引 -->
            <div class="runtime-status-box">
              <div class="runtime-header" @click="showRuntimeConfig = !showRuntimeConfig">
                <div class="runtime-indicator">
                  <span class="status-dot" :class="whisperRuntime.available ? 'dot-green' : 'dot-yellow'"></span>
                  <span v-if="whisperRuntime.available" class="runtime-ok">本地 Whisper 运行时就绪</span>
                  <span v-else class="runtime-warn">未检测到 whisper-cli 执行文件</span>
                </div>
                <button type="button" class="btn-toggle-config">
                  {{ showRuntimeConfig ? '收起配置 ▲' : '配置与说明 ▼' }}
                </button>
              </div>

              <div v-show="showRuntimeConfig" class="runtime-detail">
                <div class="form-group" style="margin-bottom: 8px;">
                  <label>whisper-cli 可执行文件路径 (留空则自动检测)</label>
                  <div class="input-with-button">
                    <input
                      type="text"
                      v-model="settings.sttWhisperPath"
                      placeholder="例如: D:\whisper.cpp\build\bin\Release\whisper-cli.exe"
                    />
                    <button type="button" class="btn-sm-secondary" @click="detectRuntime">检测环境</button>
                  </div>
                </div>
                <p class="hint" v-if="whisperRuntime.available">
                  ✓ 当前使用路径: <code>{{ whisperRuntime.path }}</code>
                </p>
                <p class="hint" v-else>
                  💡 提示：本地转录需配合 <code>whisper-cli.exe</code> 运行。您可以从 <a href="https://github.com/ggerganov/whisper.cpp/releases" target="_blank" style="color: #60a5fa; text-decoration: underline;">whisper.cpp Release</a> 下载 Windows 预编译压缩包，将 <code>whisper-cli.exe</code> 放入本程序 <code>bin/</code> 目录，或在此处填入其完整路径。
                </p>
              </div>
            </div>
          </div>

          <!-- 云端 STT 专属配置 -->
          <div v-else class="cloud-stt-panel">
            <div class="form-group" v-if="settings.sttProvider === 'siliconflow'">
              <label>硅基流动 STT 模型</label>
              <select v-model="settings.sttModel">
                <option value="FunAudioLLM/SenseVoiceSmall">FunAudioLLM/SenseVoiceSmall (阿里开源 - 极速精准 - 推荐)</option>
                <option value="Tele-AI/TeleSpeech-ASR">Tele-AI/TeleSpeech-ASR (电信开源)</option>
                <option value="openai/whisper-large-v3">openai/whisper-large-v3 (云端大模型)</option>
              </select>
              <p class="hint">国内免翻墙高速直连，SenseVoice 转录延迟低至 100ms 左右。</p>
            </div>

            <div class="form-group" v-if="settings.sttProvider === 'custom'">
              <label>自定义 STT Base URL</label>
              <input
                type="text"
                v-model="settings.sttBaseUrl"
                placeholder="例如: https://api.your-domain.com/v1"
              />
              <p class="hint">支持任意符合 OpenAI <code>/audio/transcriptions</code> 规范的 API 端点。</p>
            </div>

            <div class="form-group" v-if="settings.sttProvider === 'custom'">
              <label>自定义 STT 模型名称</label>
              <input
                type="text"
                v-model="settings.sttModel"
                placeholder="例如: whisper-1 或 whisper-large-v3"
              />
            </div>

            <div class="form-group">
              <label>STT API Key</label>
              <input
                type="password"
                v-model="settings.sttKey"
                :placeholder="sttKeyPlaceholder"
              />
              <p class="hint">若留空，将自动复用大模型配置中对应服务商的 API Key。</p>
            </div>
          </div>

          <div class="form-group">
            <label>音频转录语言</label>
            <select v-model="settings.sttLanguage">
              <option value="zh">中文 (zh - 推荐面试使用)</option>
              <option value="en">英文 (en)</option>
              <option value="auto">自动识别 (auto)</option>
              <option value="ja">日语 (ja)</option>
              <option value="ko">韩语 (ko)</option>
            </select>
          </div>

          <div class="form-group checkbox-group">
            <label>
              <input type="checkbox" v-model="settings.meetingAudio" />
              <span>开启会议回环声音捕获 (采集面试官声音)</span>
            </label>
            <p class="hint">通过系统音频回环采集面试官声音，实现对方声音与己方声音双声道独立转录。</p>
          </div>
        </div>

        <!-- Profile Tab -->
        <div v-show="activeTab === 'profile'" class="tab-pane">
          <div class="form-group">
            <label>个人简历与背景 (Resume)</label>
            <textarea
              v-model="settings.resumeText"
              rows="4"
              placeholder="粘贴个人简历核心内容、学历、工作年限与主导项目..."
            ></textarea>
          </div>

          <div class="form-group">
            <label>目标岗位描述 (JD)</label>
            <textarea
              v-model="settings.jobDescription"
              rows="3"
              placeholder="目标岗位职责与要求，AI 将自动围绕 JD 突出匹配度..."
            ></textarea>
          </div>

          <div class="form-group">
            <label>行为面试 STAR 故事库</label>
            <textarea
              v-model="settings.starStories"
              rows="3"
              placeholder="准备 2~3 个核心项目故事：情境(S) -> 任务(T) -> 行动(A) -> 结果与量化数据(R)..."
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group flex-1">
              <label>为什么选贵司</label>
              <input type="text" v-model="settings.whyCompany" placeholder="认可业务/团队/技术栈等" />
            </div>
            <div class="form-group flex-1">
              <label>离职考虑</label>
              <input type="text" v-model="settings.whyLeaving" placeholder="寻求更大技术挑战/平台等" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group flex-1">
              <label>薪资预期</label>
              <input type="text" v-model="settings.salaryTarget" placeholder="如: 30k-40k / 年薪 50w" />
            </div>
            <div class="form-group flex-1">
              <label>反问面试官清单</label>
              <input type="text" v-model="settings.questionsToAsk" placeholder="团队规模、技术挑战、业务目标等" />
            </div>
          </div>

          <div class="form-group">
            <label>自定义 AI Prompt 规则</label>
            <textarea
              v-model="settings.aiRules"
              rows="2"
              placeholder="例如: 语气要诚恳有条理、多用要点列举、避免浮夸词汇..."
            ></textarea>
          </div>
        </div>

        <!-- Appearance Tab -->
        <div v-show="activeTab === 'appearance'" class="tab-pane">
          <div class="form-group">
            <label>全局窗口透明度: {{ Math.round(settings.opacity * 100) }}%</label>
            <input type="range" min="0.2" max="1.0" step="0.02" v-model.number="settings.opacity" />
          </div>

          <div class="form-group checkbox-group">
            <label>
              <input type="checkbox" v-model="settings.screenProtection" />
              <span>防截屏/防录屏隐藏保护 (SetWindowDisplayAffinity)</span>
            </label>
            <p class="hint">在腾讯会议、Zoom、飞书等共享屏幕或录屏时，本窗口对其他人自动隐形。</p>
          </div>

          <div class="shortcuts-table">
            <h4>快捷键一览</h4>
            <div class="sc-row">
              <span class="sc-name">智能辅助 (Smart Assist)</span>
              <span class="sc-badge">Ctrl + Shift + Enter</span>
            </div>
            <div class="sc-row">
              <span class="sc-name">我该怎么说 (What to say)</span>
              <span class="sc-badge">Ctrl + Enter</span>
            </div>
            <div class="sc-row">
              <span class="sc-name">算法解题 (LeetCode Solve)</span>
              <span class="sc-badge">Ctrl + H</span>
            </div>
            <div class="sc-row">
              <span class="sc-name">展开 / 折叠面板</span>
              <span class="sc-badge">Ctrl + Shift + /</span>
            </div>
            <div class="sc-row">
              <span class="sc-name">退出程序</span>
              <span class="sc-badge">Ctrl + Shift + X</span>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <span v-if="saveMessage" class="save-msg">{{ saveMessage }}</span>
        <button class="save-btn" @click="handleSave">保存配置</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useSettings } from '../composables/useSettings'
import { useWhisper } from '../composables/useWhisper'

const emit = defineEmits(['close'])

const activeTab = ref('llm')

const {
  settings,
  saveMessage,
  currentProviderMeta,
  currentApiKey,
  currentBaseUrl,
  currentModelPair,
  sttKeyPlaceholder,
  resetCurrentBaseUrl,
  saveSettings
} = useSettings()

const {
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
  startDownloadCurrentModel,
  cancelModelDownload,
  deleteCurrentModel
} = useWhisper()

onMounted(() => {
  loadWhisperModels()
  detectRuntime()
})

async function handleSave() {
  const ok = await saveSettings()
  if (ok) {
    setTimeout(() => {
      emit('close')
    }, 800)
  }
}
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

.settings-modal {
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

.settings-tabs {
  display: flex;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid var(--border-light);
  overflow-x: auto;
}
.settings-tabs button {
  flex: 1;
  padding: 10px 12px;
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 12.5px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  white-space: nowrap;
}
.settings-tabs button.active {
  color: #fff;
  border-bottom-color: var(--accent);
  font-weight: 500;
}

.settings-body {
  padding: 16px 18px;
  overflow-y: auto;
  flex: 1;
}

.form-group {
  margin-bottom: 14px;
}
.form-group label {
  display: block;
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 5px;
}
.form-group input[type='text'],
.form-group input[type='password'],
.form-group textarea {
  width: 100%;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  color: #fff;
  padding: 8px 10px;
  font-size: 13px;
  outline: none;
}
.form-group select {
  width: 100%;
  background: rgba(26, 32, 48, 0.95);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  color: #f8fafc;
  padding: 8px 10px;
  font-size: 13px;
  outline: none;
  color-scheme: dark;
  cursor: pointer;
}
.form-group select option {
  background-color: #1e2438;
  color: #f8fafc;
  padding: 8px 10px;
}
.form-group select optgroup {
  background-color: #131722;
  color: #60a5fa;
  font-weight: 600;
  padding: 6px 10px;
}
.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: var(--accent);
}

.form-row {
  display: flex;
  gap: 10px;
}
.flex-1 {
  flex: 1;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #fff;
}
.hint {
  font-size: 11.5px;
  color: var(--text-muted);
  margin-top: 4px;
  line-height: 1.4;
}

.modal-footer {
  padding: 12px 18px;
  border-top: 1px solid var(--border-light);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}
.save-msg {
  font-size: 12px;
  color: var(--live-green);
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

.shortcuts-table {
  background: rgba(0, 0, 0, 0.2);
  border-radius: var(--radius-md);
  padding: 12px;
  margin-top: 10px;
}
.shortcuts-table h4 {
  font-size: 13px;
  margin-bottom: 8px;
  color: #fff;
}
.sc-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  font-size: 12px;
}
.sc-row:last-child {
  border-bottom: none;
}
.sc-badge {
  font-family: var(--font-mono);
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  color: #60a5fa;
  font-size: 11px;
}

.label-with-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
}
.label-with-action label {
  margin-bottom: 0;
}
.text-link-btn {
  background: transparent;
  border: none;
  color: #60a5fa;
  font-size: 11px;
  cursor: pointer;
  padding: 0;
}
.text-link-btn:hover {
  text-decoration: underline;
  color: #93c5fd;
}
.provider-doc {
  color: #93c5fd;
  margin-top: 5px;
  font-size: 11.5px;
}
.model-quick-select {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 5px;
  margin-top: 6px;
}
.model-quick-select .sub-label {
  font-size: 11px;
  color: var(--text-muted);
}
.model-tag {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #cbd5e1;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 10.5px;
  cursor: pointer;
  font-family: var(--font-mono);
  transition: all 0.15s ease;
}
.model-tag:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: var(--accent);
  color: #93c5fd;
}

/* Local STT & Whisper Models */
.local-stt-panel,
.cloud-stt-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 8px;
}

.label-with-badge {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 5px;
}
.label-with-badge label {
  margin-bottom: 0;
}

.model-badge {
  font-size: 11px;
  font-weight: 500;
  padding: 2px 7px;
  border-radius: 10px;
}

.badge-ready {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
  border: 1px solid rgba(34, 197, 94, 0.3);
}
.badge-missing {
  background: rgba(234, 179, 8, 0.15);
  color: #facc15;
  border: 1px solid rgba(234, 179, 8, 0.3);
}
.badge-downloading {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.model-status-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
}

.download-prompt-box {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.prompt-text {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}
.prompt-icon {
  font-size: 20px;
  line-height: 1;
}
.prompt-desc strong {
  display: block;
  font-size: 13px;
  color: #f8fafc;
  margin-bottom: 4px;
}
.prompt-desc p {
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.4;
  margin: 0;
}
.prompt-desc code {
  background: rgba(255, 255, 255, 0.08);
  padding: 1px 4px;
  border-radius: 3px;
  font-family: monospace;
  color: #93c5fd;
}

.btn-download-action {
  background: linear-gradient(135deg, #2563eb, #3b82f6);
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  padding: 9px 14px;
  font-size: 12.5px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}
.btn-download-action:hover {
  filter: brightness(1.1);
  transform: translateY(-1px);
}

.download-progress-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.download-info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #f8fafc;
}
.download-speed-badge {
  background: rgba(59, 130, 246, 0.2);
  color: #93c5fd;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
}
.progress-bar-container {
  width: 100%;
  height: 8px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 4px;
  overflow: hidden;
}
.progress-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6, #60a5fa);
  border-radius: 4px;
  transition: width 0.2s ease;
}
.download-meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11.5px;
  color: var(--text-muted);
}
.btn-cancel-download {
  background: transparent;
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.3);
  border-radius: var(--radius-sm);
  padding: 2px 8px;
  font-size: 11px;
  cursor: pointer;
}
.btn-cancel-download:hover {
  background: rgba(248, 113, 113, 0.1);
}
.download-error-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.error-msg-row {
  color: #f87171;
  font-size: 12px;
}
.download-ready-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.ready-info {
  display: flex;
  align-items: center;
  gap: 8px;
}
.ready-icon {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(34, 197, 94, 0.2);
  color: #4ade80;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}
.ready-info strong {
  display: block;
  font-size: 12.5px;
  color: #f8fafc;
}
.ready-info p {
  font-size: 11.5px;
  color: var(--text-muted);
  margin: 2px 0 0 0;
}
.ready-info code {
  background: rgba(255, 255, 255, 0.08);
  padding: 1px 4px;
  border-radius: 3px;
  color: #93c5fd;
}
.btn-text-danger {
  background: transparent;
  color: #f87171;
  border: none;
  font-size: 12px;
  text-decoration: underline;
  cursor: pointer;
  padding: 4px;
}
.btn-text-danger:hover {
  color: #ef4444;
}

.runtime-status-box {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
}
.runtime-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}
.runtime-indicator {
  display: flex;
  align-items: center;
  font-size: 12px;
}
.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 6px;
}
.dot-green {
  background: #22c55e;
  box-shadow: 0 0 6px rgba(34, 197, 94, 0.6);
}
.dot-yellow {
  background: #eab308;
  box-shadow: 0 0 6px rgba(234, 179, 8, 0.6);
}
.runtime-ok {
  color: #4ade80;
  font-weight: 500;
}
.runtime-warn {
  color: #facc15;
  font-weight: 500;
}
.btn-toggle-config {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 11.5px;
  cursor: pointer;
}
.btn-toggle-config:hover {
  color: #fff;
}
.runtime-detail {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed rgba(255, 255, 255, 0.1);
}
.input-with-button {
  display: flex;
  gap: 8px;
}
.btn-sm-secondary {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid var(--border-light);
  color: #fff;
  border-radius: var(--radius-sm);
  padding: 0 12px;
  font-size: 12px;
  white-space: nowrap;
  cursor: pointer;
}
.btn-sm-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.spinner-small {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  vertical-align: middle;
  margin-right: 4px;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
