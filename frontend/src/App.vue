<template>
  <div id="main-container" :style="{ opacity: settings.opacity }">
    <!-- Top Pill Toolbar -->
    <header class="toolbar glass-panel clickable">
      <!-- Drag Handle -->
      <div class="drag-pill" style="--wails-draggable: drag">
        <span class="drag-dots"></span>
        <span class="app-title">Copilot</span>
      </div>

      <!-- Live Listening Indicator -->
      <span
        class="status-indicator"
        :class="{ active: isListening, busy: isBusy }"
        :title="isListening ? '正在听录会议中...' : '听录未开启'"
      ></span>

      <!-- Action Buttons -->
      <button class="tb-btn" @click="showTutorial = true" title="使用教程与说明">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
          <line x1="12" y1="17" x2="12.01" y2="17"></line>
        </svg>
      </button>

      <button class="tb-btn" @click="openSettings" title="设置中心">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"></circle>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
        </svg>
      </button>

      <div class="tb-divider"></div>

      <!-- Toggle Hide / Collapse -->
      <button class="tb-btn text-btn" @click="isCollapsed = !isCollapsed" title="折叠/展开面板 (Ctrl+Shift+/)">
        <svg class="icon chev" :class="{ rotated: isCollapsed }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="18 15 12 9 6 15"></polyline>
        </svg>
        <span>{{ isCollapsed ? '展开' : '折叠' }}</span>
      </button>

      <div class="tb-divider"></div>

      <!-- Toggle Session / Listening -->
      <button
        class="tb-btn session-btn"
        :class="{ active: isListening }"
        @click="toggleListening"
        :title="isListening ? '停止听录' : '启动听录 (采集麦克风与会议声音)'"
      >
        <span class="mic-dot" :class="{ 'pulsing-live': isListening }"></span>
        <span>{{ isListening ? '听录中' : '开启听录' }}</span>
      </button>

      <div class="tb-divider"></div>

      <!-- Opacity Popover Toggle -->
      <div class="opacity-container">
        <button class="tb-btn" @click="showOpacityPopover = !showOpacityPopover" title="透明度调节">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M12 2a10 10 0 0 1 0 20z" fill="currentColor"></path>
          </svg>
        </button>
        <div v-if="showOpacityPopover" class="opacity-popover glass-panel">
          <span class="pop-title">透明度: {{ Math.round(settings.opacity * 100) }}%</span>
          <input
            type="range"
            min="0.2"
            max="1.0"
            step="0.02"
            v-model.number="settings.opacity"
            @change="saveSettingsSilently"
          />
        </div>
      </div>

      <div class="tb-divider"></div>

      <!-- Close App -->
      <button class="tb-btn close-btn" @click="quitApp" title="退出程序 (Ctrl+Shift+X)">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </header>

    <!-- Main Content Area -->
    <main v-show="!isCollapsed" class="panel-layout">
      <!-- Main Copilot Card -->
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
              @click="showSidebar = !showSidebar"
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

      <!-- Sliding Transcript History Drawer -->
      <aside v-if="showSidebar" class="sidebar glass-panel clickable">
        <div class="sidebar-header">
          <div class="sidebar-title">
            <span>实时对话流</span>
            <span class="turn-count">{{ transcriptTurns.length }} 条记录</span>
          </div>
          <button class="close-sidebar-btn" @click="showSidebar = false">✕</button>
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
    </main>

    <!-- Settings Modal -->
    <div v-if="showSettingsModal" class="modal-backdrop clickable" @click.self="showSettingsModal = false">
      <div class="settings-modal glass-panel">
        <div class="modal-header">
          <h2>设置中心</h2>
          <button class="close-modal-btn" @click="showSettingsModal = false">✕</button>
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
          <button class="save-btn" @click="saveSettings">保存配置</button>
        </div>
      </div>
    </div>

    <!-- Tutorial Modal -->
    <div v-if="showTutorial" class="modal-backdrop clickable" @click.self="showTutorial = false">
      <div class="tutorial-modal glass-panel">
        <div class="modal-header">
          <h2>使用指南 & 隐形说明</h2>
          <button class="close-modal-btn" @click="showTutorial = false">✕</button>
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
          <button class="save-btn" @click="showTutorial = false">我明白了</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { marked } from 'marked'
import { AppService } from '../bindings/interview-copilot/internal'
import { Events } from '@wailsio/runtime'

// State
const defaultBaseUrls = {
  deepseek: 'https://api.deepseek.com/v1',
  qwen: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
  doubao: 'https://ark.cn-beijing.volces.com/api/v3',
  siliconflow: 'https://api.siliconflow.cn/v1',
  moonshot: 'https://api.moonshot.cn/v1',
  zhipu: 'https://open.bigmodel.cn/api/paas/v4',
  openai: 'https://api.openai.com/v1',
  ollama: 'http://localhost:11434/v1',
  groq: 'https://api.groq.com/openai/v1',
  custom: 'https://api.openai.com/v1'
}

const defaultApiKeys = {
  deepseek: '',
  qwen: '',
  doubao: '',
  siliconflow: '',
  moonshot: '',
  zhipu: '',
  openai: '',
  ollama: '',
  anthropic: '',
  gemini: '',
  groq: '',
  custom: ''
}

const defaultModels = {
  deepseek: { fast: 'deepseek-chat', smart: 'deepseek-reasoner' },
  qwen: { fast: 'qwen-plus', smart: 'qwen-max' },
  doubao: { fast: 'doubao-1.5-lite-32k', smart: 'doubao-1.5-pro-32k' },
  siliconflow: { fast: 'deepseek-ai/DeepSeek-V3', smart: 'deepseek-ai/DeepSeek-R1' },
  moonshot: { fast: 'moonshot-v1-8k', smart: 'moonshot-v1-32k' },
  zhipu: { fast: 'glm-4-flash', smart: 'glm-4-plus' },
  openai: { fast: 'gpt-4o-mini', smart: 'gpt-4o' },
  ollama: { fast: 'qwen2.5:7b', smart: 'deepseek-r1:7b' },
  anthropic: { fast: 'claude-3-5-haiku-20241022', smart: 'claude-3-5-sonnet-20241022' },
  gemini: { fast: 'gemini-2.5-flash', smart: 'gemini-1.5-pro' },
  groq: { fast: 'llama-3.1-8b-instant', smart: 'llama-3.3-70b-versatile' },
  custom: { fast: 'default', smart: 'default' }
}

const providerPresets = {
  deepseek: {
    name: 'DeepSeek (深度求索)',
    url: 'https://api.deepseek.com/v1',
    keyHint: 'sk-... (在 platform.deepseek.com 获取)',
    fastSuggestions: ['deepseek-chat'],
    smartSuggestions: ['deepseek-reasoner'],
    doc: '官方 API 直连，支持 V3 非思考和 R1 深度思考推理链'
  },
  qwen: {
    name: '通义千问 (阿里百炼)',
    url: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
    keyHint: 'sk-... (在阿里云百炼 DashScope 控制台获取)',
    fastSuggestions: ['qwen-plus', 'qwen-turbo', 'qwen2.5-coder-32b-instruct'],
    smartSuggestions: ['qwen-max', 'qwen-plus', 'qwen2.5-72b-instruct'],
    doc: '兼容 OpenAI 格式，千问中文语境与代码生成表现优异'
  },
  doubao: {
    name: '豆包 (字节跳动 火山引擎 Ark)',
    url: 'https://ark.cn-beijing.volces.com/api/v3',
    keyHint: '火山方舟 API Key (在 console.volcengine.com/ark 获取)',
    fastSuggestions: ['doubao-1.5-lite-32k', 'doubao-pro-4k', 'ep-xxxx-lite'],
    smartSuggestions: ['doubao-1.5-pro-32k', 'doubao-pro-32k', 'ep-xxxx-pro'],
    doc: '支持模型标准名或火山方舟接入点 ID (如 ep-2025xxxx)'
  },
  siliconflow: {
    name: '硅基流动 (SiliconFlow)',
    url: 'https://api.siliconflow.cn/v1',
    keyHint: 'sk-... (在 siliconflow.cn 获取，国内免翻极速)',
    fastSuggestions: ['deepseek-ai/DeepSeek-V3', 'Qwen/Qwen2.5-72B-Instruct', 'Qwen/Qwen2.5-Coder-32B-Instruct'],
    smartSuggestions: ['deepseek-ai/DeepSeek-R1', 'Pro/deepseek-ai/DeepSeek-R1'],
    doc: '国内免翻墙极速算力，低延迟运行 DeepSeek-V3/R1 与千问全系'
  },
  moonshot: {
    name: 'Kimi (月之暗面)',
    url: 'https://api.moonshot.cn/v1',
    keyHint: 'sk-... (在 platform.moonshot.cn 获取)',
    fastSuggestions: ['moonshot-v1-8k', 'moonshot-v1-32k'],
    smartSuggestions: ['moonshot-v1-32k', 'moonshot-v1-128k'],
    doc: '月之暗面长文本模型，中文逻辑推理优秀'
  },
  zhipu: {
    name: '智谱清言 (GLM-4)',
    url: 'https://open.bigmodel.cn/api/paas/v4',
    keyHint: 'API Key (在 bigmodel.cn 获取)',
    fastSuggestions: ['glm-4-flash', 'glm-4-air'],
    smartSuggestions: ['glm-4-plus', 'glm-4-0520'],
    doc: '智谱 AI 原生兼容端点，通用能力出众'
  },
  openai: {
    name: 'OpenAI (官方)',
    url: 'https://api.openai.com/v1',
    keyHint: 'sk-... (在 platform.openai.com 获取)',
    fastSuggestions: ['gpt-4o-mini', 'gpt-3.5-turbo'],
    smartSuggestions: ['gpt-4o', 'o1-mini', 'o1-preview'],
    doc: '业界标杆，原生多模态视觉识别'
  },
  ollama: {
    name: 'Ollama (本地私有)',
    url: 'http://localhost:11434/v1',
    keyHint: '本地服务无需 Key (留空即可)',
    fastSuggestions: ['qwen2.5:7b', 'llama3.2', 'deepseek-coder:6.7b'],
    smartSuggestions: ['deepseek-r1:7b', 'deepseek-r1:14b', 'qwen2.5:14b'],
    doc: '完全本地离线运行，零费用、数据绝不泄露'
  },
  anthropic: {
    name: 'Claude (Anthropic)',
    url: 'https://api.anthropic.com/v1',
    keyHint: 'sk-ant-...',
    fastSuggestions: ['claude-3-5-haiku-20241022'],
    smartSuggestions: ['claude-3-5-sonnet-20241022'],
    doc: '长文思维与架构代码能力极强'
  },
  gemini: {
    name: 'Google Gemini',
    url: 'https://generativelanguage.googleapis.com',
    keyHint: 'AI Studio API Key',
    fastSuggestions: ['gemini-2.5-flash'],
    smartSuggestions: ['gemini-1.5-pro'],
    doc: '超长上下文感知与多模态视觉'
  },
  groq: {
    name: 'Groq (极速推断)',
    url: 'https://api.groq.com/openai/v1',
    keyHint: 'gsk_...',
    fastSuggestions: ['llama-3.1-8b-instant'],
    smartSuggestions: ['llama-3.3-70b-versatile'],
    doc: '基于定制硬件的极致推理速度'
  },
  custom: {
    name: '自定义端点',
    url: 'https://api.openai.com/v1',
    keyHint: '自定义端点 API Key (可选)',
    fastSuggestions: ['default'],
    smartSuggestions: ['default'],
    doc: '支持 OneAPI、NewAPI、vLLM、FastChat 等中转或网关'
  }
}

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

function resetCurrentBaseUrl() {
  const p = settings.value.provider
  if (providerPresets[p]?.url) {
    if (!settings.value.baseUrls) settings.value.baseUrls = {}
    settings.value.baseUrls[p] = providerPresets[p].url
    if (p === 'custom') settings.value.baseUrl = providerPresets[p].url
  }
}

const isListening = ref(false)
const isBusy = ref(false)
const isCollapsed = ref(false)
const showSidebar = ref(false)
const showSettingsModal = ref(false)
const showTutorial = ref(false)
const showOpacityPopover = ref(false)
const activeTab = ref('llm')
const saveMessage = ref('')
const copyTip = ref('复制')

const inputText = ref('')
const currentMode = ref('')
const detectedCategory = ref('')
const responseContent = ref('')
const thinkingContent = ref('')
const showThinking = ref(true)
const transcriptTurns = ref([])

const responseBox = ref(null)
const turnsBox = ref(null)

// Audio context references
let audioCtx = null
let micStream = null
let micProcessor = null
let systemStream = null
let systemProcessor = null

const renderedMarkdown = computed(() => {
  if (!responseContent.value) return ''
  return marked.parse(responseContent.value)
})

const currentModeLabel = computed(() => {
  const map = {
    say: '💬 我该怎么说',
    assist: '✨ 智能辅助',
    leetcode: '💻 算法解题',
    recap: '📋 会议复盘',
    ask: '❓ 自由提问',
    answerThis: '⚡ 针对提问回答'
  }
  return map[currentMode.value] || ''
})

const detectedCategoryLabel = computed(() => {
  const map = {
    behavioral: 'STAR 行为题',
    motivation: '求职动机',
    situational: '情境应对',
    experience: '项目经历',
    compensation: '薪资与反问',
    technical: '技术与设计',
    leetcode: '算法编程'
  }
  return map[detectedCategory.value] || ''
})

// Whisper Models & STT State
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

const currentLocalModelInfo = computed(() => {
  const name = settings.value.sttLocalModel || 'base'
  return whisperModels.value.find(m => m.name === name) || null
})

const isDownloadingCurrentModel = computed(() => {
  const current = settings.value.sttLocalModel || 'base'
  return whisperDownloadState.value.status === 'downloading' && whisperDownloadState.value.model === current
})

const downloadProgressText = computed(() => {
  return `${(whisperDownloadState.value.progress || 0).toFixed(0)}%`
})

const downloadingModelTitle = computed(() => {
  const m = whisperModels.value.find(x => x.name === whisperDownloadState.value.model)
  return m ? m.displayName : (whisperDownloadState.value.model || '模型')
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

// Methods
function formatTime(ts) {
  if (!ts) return ''
  const d = new Date(ts)
  return d.toTimeString().slice(0, 8)
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

function openSettings() {
  showSettingsModal.value = true
  loadWhisperModels()
  detectRuntime()
}

async function saveSettings() {
  try {
    await AppService.SaveSettings(settings.value)
    saveMessage.value = '✓ 已保存'
    setTimeout(() => {
      saveMessage.value = ''
      showSettingsModal.value = false
    }, 800)
  } catch (err) {
    saveMessage.value = '保存失败: ' + err.message
  }
}

async function saveSettingsSilently() {
  try {
    await AppService.SaveSettings(settings.value)
  } catch (err) {
    console.error(err)
  }
}

function toggleSmartMode() {
  settings.value.smart = !settings.value.smart
  saveSettingsSilently()
}

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

function clearTranscript() {
  AppService.ClearTranscript()
  transcriptTurns.value = []
}

function quitApp() {
  AppService.QuitApp()
}

// Audio Capture Implementation
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

// Lifecycle & Wails Events
onMounted(async () => {
  await loadSettings()
  await loadHistory()
  await loadWhisperModels()
  detectRuntime()

  Events.On('stt:model_download_progress', (data) => {
    const payload = data.data || data
    whisperDownloadState.value = payload
    if (payload.status === 'completed' || payload.status === 'failed') {
      loadWhisperModels()
    }
  })

  // Register Wails Event Listeners
  Events.On('llm:start', (data) => {
    isBusy.value = true
    const payload = data.data || data
    currentMode.value = payload.mode || currentMode.value
    detectedCategory.value = payload.category || ''
    responseContent.value = ''
    thinkingContent.value = ''
  })

  Events.On('llm:thinking', (data) => {
    const payload = data.data || data
    if (payload.thought) {
      thinkingContent.value += payload.thought
    }
  })

  Events.On('llm:token', (data) => {
    const payload = data.data || data
    if (payload.token) {
      responseContent.value += payload.token
      nextTick(() => {
        if (responseBox.value) {
          responseBox.value.scrollTop = responseBox.value.scrollHeight
        }
      })
    }
  })

  Events.On('llm:done', (data) => {
    isBusy.value = false
  })

  Events.On('llm:error', (data) => {
    isBusy.value = false
    const payload = data.data || data
    responseContent.value += `\n\n> ⚠️ **错误提示**: ${payload.message}`
  })

  Events.On('transcript:turn', (data) => {
    const turn = data.data || data
    transcriptTurns.value.push(turn)
    nextTick(() => {
      if (turnsBox.value) {
        turnsBox.value.scrollTop = turnsBox.value.scrollHeight
      }
    })
  })

  Events.On('transcript:cleared', () => {
    transcriptTurns.value = []
  })

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

/* Toolbar */
.toolbar {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  border-radius: var(--radius-full);
  margin-bottom: 8px;
  z-index: 50;
}

.drag-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 8px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-full);
  cursor: grab;
}
.drag-pill:active {
  cursor: grabbing;
}
.drag-dots {
  width: 12px;
  height: 8px;
  background-image: radial-gradient(circle, rgba(255, 255, 255, 0.9) 1px, transparent 1.2px);
  background-size: 4px 4px;
}
.app-title {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.5px;
  color: #fff;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  transition: all 0.3s;
}
.status-indicator.active {
  background: var(--live-green);
  box-shadow: 0 0 8px var(--live-green);
}
.status-indicator.busy {
  background: #f59e0b;
  box-shadow: 0 0 8px #f59e0b;
}

.tb-btn {
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: var(--radius-full);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 5px 7px;
  font-size: 12px;
  transition: all 0.15s ease;
}
.tb-btn:hover {
  background: var(--bg-hover);
  color: #fff;
}

.tb-btn .icon {
  width: 15px;
  height: 15px;
}
.tb-btn.text-btn {
  padding: 4px 10px;
  gap: 4px;
}
.tb-btn.session-btn {
  padding: 4px 12px;
  gap: 6px;
  font-weight: 500;
}
.tb-btn.session-btn.active {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}
.mic-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-muted);
}
.mic-dot.pulsing-live {
  background: var(--live-green);
}

.tb-divider {
  width: 1px;
  height: 16px;
  background: rgba(255, 255, 255, 0.14);
  margin: 0 2px;
}

.tb-btn.close-btn:hover {
  background: rgba(239, 68, 68, 0.25);
  color: #ef4444;
}

.chev {
  transition: transform 0.2s ease;
}
.chev.rotated {
  transform: rotate(180deg);
}

/* Opacity Popover */
.opacity-container {
  position: relative;
}
.opacity-popover {
  position: absolute;
  top: 110%;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 14px;
  border-radius: var(--radius-md);
  min-width: 170px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  z-index: 100;
}
.pop-title {
  font-size: 12px;
  color: var(--text-secondary);
}
.opacity-popover input[type='range'] {
  accent-color: var(--accent);
  cursor: pointer;
}

/* Panel Layout */
.panel-layout {
  display: flex;
  gap: 12px;
  width: 100%;
  max-width: 740px;
  height: calc(100% - 50px);
  position: relative;
}

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
  box-shadow: 0 0 12px var(--accent-glow);
}
.send-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}
.send-btn .icon {
  width: 14px;
  height: 14px;
}

/* Sidebar Drawer */
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
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
}

/* Modals */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}
.settings-modal, .tutorial-modal {
  width: 90%;
  max-width: 580px;
  max-height: 85vh;
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
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
.settings-body, .tutorial-body {
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
