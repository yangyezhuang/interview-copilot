# Interview Copilot (面试辅助副驾驶)

<div align="center">

**参考致敬 [Blueturboguy07/cue](https://github.com/Blueturboguy07/cue) ，使用 Wails v3 + Go + Vue 3 进行全面重构的轻量级隐形 AI 实时面试副驾。**

浮动在屏幕之上，眼观屏幕、耳听会议，在面试与现场编程中为您实时提示最优解。

**在屏幕共享与录屏中全自动隐形 · 极低资源占用 · 自带密钥 (BYOK) 直连主流大模型**

[特性概览](#-核心特性) • [快捷键](#-快捷键指南) • [架构优势](#-为什么使用-wails-重写) • [快速上手](#-快速上手) • [常见问题与隐形设置](#-会议软件隐形说明)

</div>

---

> [!IMPORTANT]
> **使用须知与合规声明**：
> 本工具旨在用于个人学习、面试模拟、技术练习、会议记录与辅助备忘。在正式求职面试或特定监控考试中，请严格遵守招聘方与平台的规定。屏幕共享隐形功能基于系统底层 API 实现，但外部物理设备（如手机拍照、外接采集卡）仍然可以记录屏幕画面。使用者需对自身行为负全部法律与道德责任。

---

## ✨ 核心特性

### 1. 🛡️ 屏幕共享全自动隐形 (Screen-Share Invisible)
- 基于 Windows 底层 `SetWindowDisplayAffinity(hwnd, WDA_EXCLUDEFROMCAPTURE)` API 实现。
- 窗口在 DWM (桌面窗口管理器) 渲染树中被标记为排除捕获：在 **腾讯会议、Zoom、飞书、钉钉、Google Meet、Microsoft Teams、OBS、Windows 录屏** 中完全隐形，共享全屏或特定窗口时，对方屏幕上均不会出现本程序窗口。

### 2. 🎯 三路多模态感知输入
- **屏幕感知 (Screen)**：按下快捷键时自动截取当前屏幕（原生 GDI/DXGI 截图，自动缩放至最佳分辨率并压缩传输），秒级解析算法题与页面内容。
- **麦克风输入 (You)**：采集候选人本人的口述内容，避免 AI 回复与已有阐述产生冲突。
- **会议声音 (Them)**：基于系统音频回环采集面试官的声音，自动分析问题并生成话术建议。

### 3. 🧠 场景化面试上下文引擎
系统自动对面试官的问题进行语义识别与分类（支持中英文双语识别），并自动调取预设的候选人画像信息：
- **行为面试 (Behavioral - STAR)**：自动匹配预备的真实 STAR 故事库，按照「情境 (S) → 任务 (T) → 行动 (A) → 结果 (R)」结构化输出，带量化指标。
- **求职动机 (Motivation)**：结合「为什么选择贵司」和「离职考虑」，输出真诚而得体的高情商回答。
- **情境处置 (Situational)**：展现清晰的决策优先级与结构化逻辑。
- **项目经历 (Experience)**：提取简历中的核心项目，直击底层技术栈与个人贡献。
- **技术与系统设计 (Technical)**：突出架构考量、权衡取舍 (Trade-offs) 与高并发方案。
- **薪酬与反问 (Compensation & Q&A)**：输出自信得体的薪资预期，并提供高水准反问问题清单。

### 4. 💻 算法编程与代码解题 (LeetCode Solve)
- 遇到算法或编程题，一键触发解题模式：
  1. **【题目重述】**：一句话直击要点与边界条件；
  2. **【解题思路】**：讲解最优算法推导（双指针/单调栈/动态规划/图论等）；
  3. **【代码实现】**：输出整洁、带注释、符合规范的代码（自动匹配页面语言）；
  4. **【复杂度分析】**：准确给出时间复杂度与空间复杂度分析。

### 5. 🌐 全面支持主流大模型 (BYOK)
直接配置属于您的 API 密钥，零中间服务器转发，隐私安全完全本地化：
- **阿里通义千问 (Qwen / 阿里百炼 DashScope)**：支持 `qwen-plus`、`qwen-max`、`qwen2.5-coder-32b-instruct` 等，国内免翻直连，编程能力出众。
- **字节豆包 (Doubao / 火山引擎 Ark)**：支持 `doubao-1.5-pro-32k`、`doubao-1.5-lite-32k` 或火山方舟接入点 ID，性价比极高。
- **硅基流动 (SiliconFlow)**：免翻墙高速中转，支持 `deepseek-ai/DeepSeek-V3`、`deepseek-ai/DeepSeek-R1`、千问全系。
- **DeepSeek (深度求索官方)**：支持 `deepseek-chat` 与 `deepseek-reasoner`（实时展现 R1 深度思考推理链）。
- **Kimi (月之暗面 Moonshot)**：支持 `moonshot-v1-8k`、`moonshot-v1-32k`，长文本理解优秀。
- **智谱清言 (GLM-4)**：支持 `glm-4-flash`、`glm-4-plus` 等。
- **Ollama / 本地私有模型**：本地离线运行（如 `qwen2.5:7b`、`deepseek-r1:7b`），无需 API 密钥，零网络费用与数据泄露风险。
- **OpenAI**：支持 `gpt-4o`、`gpt-4o-mini`、`o1` 等官方端点。
- **Anthropic Claude**：支持 `claude-3-5-sonnet`、`claude-3-5-haiku` 等。
- **Google Gemini**：支持 `gemini-2.5-flash`、`gemini-1.5-pro` 等。
- **Groq**：极速 LPU 推理与 Whisper 语音转录。
- **自定义 OpenAI 兼容接口**：支持 OneAPI、NewAPI、vLLM、FastChat 等任意兼容端点。

### 6. 🎙️ 毫秒级语音转文字 (STT)
- **Groq Whisper**：推荐使用，基于 LPU 硬件加速的 `whisper-large-v3-turbo`，转录延迟低至数百毫秒。
- **OpenAI Whisper**：官方 `whisper-1` 高精度语音识别。
- **自定义 STT**：支持任意兼容 `/v1/audio/transcriptions` 的本地或远程端点。

### 7. 🎨 极致的毛玻璃美学设计 (Glassmorphism)
- 悬浮式灵动胶囊工具栏：随意拖拽至屏幕任意位置。
- 无级透明度滑块调节 (20% ~ 100%)，背景穿透不挡视线。
- 动态转录抽屉：展示面试官与自己的对话时间线，支持对面试官的历史提问一键点击「⚡ 回答此题」。
- Markdown 语法高亮与一键复制代码块。

---

## ⌨️ 快捷键指南

所有快捷键均为**系统全局快捷键**，即使程序处于后台或无焦点状态，均可立即响应：

| 快捷键 | 功能 | 说明 |
|---|---|---|
| `Ctrl + Enter` | **我该怎么说? (What to say)** | 依据最新会议对话，输出 1 段自然沉稳的第一人称口述回答 |
| `Ctrl + Shift + Enter` | **智能辅助 (Smart Assist)** | 同时截取当前屏幕 + 最近对话，智能判定题型并给出全套方案 |
| `Ctrl + H` | **算法编程解题 (LeetCode)** | 截取屏幕上的算法题目，输出思路、代码与时空复杂度 |
| `Ctrl + Shift + /` | **折叠 / 展开悬浮窗** | 一键收起为顶部迷你胶囊，或展开完整面板 |
| `Ctrl + Shift + X` | **退出程序** | 快速安全退出程序 |

---

## 🚀 为什么使用 Wails 重写？

原版项目基于 **Electron** 构建，存在以下痛点：
1. **体积臃肿**：包含完整的 Chromium + Node.js 运行时，打包后体积通常在 150MB+。
2. **内存消耗高**：Electron 多进程架构导致开销动辄 200MB ~ 400MB 内存。
3. **本地系统 API 调用受限**：复杂的 Native 模块编译容易产生兼容性问题。

**使用 Wails v3 + Go 重写后的优势**：
- ⚡ **超轻量内存**：采用 Go 原生并发后端 + 系统内置 WebView2 运行时，常驻内存仅需 ~35MB。
- 📦 **编译单文件**：打包输出小巧单执行文件，启动瞬间毫秒级完成。
- 🔒 **原生系统级集成**：Windows API（`SetWindowDisplayAffinity`、`RegisterHotKey`）以 Go 原生系统调用实现，无任何第三方 Node.js C++ 插件依赖，稳定纯粹。
- 🚀 **极高并发吞吐**：音频转码 (PCM -> WAV)、截图压缩与多模态流式解析全部由 Go 协程在后台异步处理，界面流畅丝滑。

---

## 🛠️ 项目结构

```
interview-copilot/
├── main.go               # 应用程序入口与 Wails 窗口初始化
├── service.go            # 暴露给前端的 Wails 核心服务 (AppService)
├── config.go             # 本地配置存储与管理 (Settings)
├── screen.go             # 屏幕截图捕获与自适应压缩 (Vision-ready)
├── wav.go                # 16kHz PCM 音频流包装器
├── stt.go                # 语音识别客户端 (OpenAI/Groq/Custom Whisper)
├── llm.go                # 统一多厂商流式大模型客户端 (SSE解析)
├── prompts.go            # 提示词工程与模式定义 (Assist/Say/LeetCode/Recap等)
├── context.go            # 面试题型分类判断与个人画像上下文组装
├── meetings.go           # 会议对话记忆库与本地持久化
├── window_windows.go     # Windows 原生防录屏隐藏与全局热键
├── window_other.go       # 跨平台兼容桩函数
├── frontend/             # Vue 3 现代化前端界面
│   ├── src/
│   │   ├── App.vue       # 主毛玻璃悬浮界面与交互逻辑
│   │   └── main.js       # 前端入口
│   ├── public/
│   │   └── style.css     # 毛玻璃设计样式表
│   └── package.json
└── README.md             # 中文使用文档
```

---

## 🏁 快速上手

### 环境准备
- **Go**: 1.22 或更高版本（推荐 Go 1.23+）
- **Node.js**: 18.0 或更高版本
- **Wails CLI**: Wails v3 CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- **操作系统**: Windows 10 (Build 19041+) 或 Windows 11

### 1. 克隆并安装依赖
```bash
git clone https://github.com/your-username/interview-copilot.git
cd interview-copilot

# 安装前端依赖
cd frontend
npm install
cd ..
```

### 2. 开发模式运行
在根目录下直接启动开发热重载模式：
```bash
wails3 dev
```
或者分步运行前端与后端：
```bash
# 终端 1: 构建前端
cd frontend && npm run build

# 终端 2: 编译启动 Go 后端
go run .
```

### 3. 生产发布打包
执行以下命令直接生成独立的发布版本二进制：
```bash
# 构建前端生产资源
cd frontend && npm run build && cd ..

# 生成 wails 资源并编译
go build -tags production -trimpath -ldflags="-w -s -H windowsgui" -o bin/interview-copilot.exe .
```
编译产物即位于 `bin/interview-copilot.exe`，双击即可无黑框运行。

---

## ⚙️ 配置指南

启动后点击顶部工具栏的 **设置中心 (齿轮图标)**：

### 1. 模型与密钥 (LLM)
系统已内置各大主流厂商的标准接口端点与推荐模型，只需选择并填入对应的 `API Key`（端点与模型均支持自定义修改）：

| 服务商 | 接口端点 (Base URL) | 推荐极速模型 (Fast) | 推荐深度模型 (Smart) | 密钥获取方式 |
|---|---|---|---|---|
| **DeepSeek (官方)** | `https://api.deepseek.com/v1` | `deepseek-chat` | `deepseek-reasoner` | [platform.deepseek.com](https://platform.deepseek.com) |
| **通义千问 (阿里百炼)** | `https://dashscope.aliyuncs.com/compatible-mode/v1` | `qwen-plus` | `qwen-max` | [bailian.console.aliyun.com](https://bailian.console.aliyun.com) |
| **豆包 (火山引擎 Ark)** | `https://ark.cn-beijing.volces.com/api/v3` | `doubao-1.5-lite-32k` | `doubao-1.5-pro-32k` | [console.volcengine.com/ark](https://console.volcengine.com/ark) (支持填模型名或 Endpoint ID) |
| **硅基流动 (SiliconFlow)** | `https://api.siliconflow.cn/v1` | `deepseek-ai/DeepSeek-V3` | `deepseek-ai/DeepSeek-R1` | [siliconflow.cn](https://siliconflow.cn) (免翻墙，极速响应) |
| **Kimi (月之暗面)** | `https://api.moonshot.cn/v1` | `moonshot-v1-8k` | `moonshot-v1-32k` | [platform.moonshot.cn](https://platform.moonshot.cn) |
| **智谱清言 (GLM-4)** | `https://open.bigmodel.cn/api/paas/v4` | `glm-4-flash` | `glm-4-plus` | [bigmodel.cn](https://bigmodel.cn) |
| **OpenAI (官方)** | `https://api.openai.com/v1` | `gpt-4o-mini` | `gpt-4o` | [platform.openai.com](https://platform.openai.com) |
| **Ollama (本地私有)** | `http://localhost:11434/v1` | `qwen2.5:7b` | `deepseek-r1:7b` | 本地安装 Ollama，无需填 Key |
| **Groq (极速推断)** | `https://api.groq.com/openai/v1` | `llama-3.1-8b-instant` | `llama-3.3-70b-versatile` | [console.groq.com](https://console.groq.com) |
| **自定义兼容端点** | 自定义 (如 `https://your-api.com/v1`) | `default` | `default` | 兼容 OneAPI / NewAPI / vLLM |

### 2. 音频与语音识别 (Audio & STT)
- 推荐使用 **Groq Whisper**，转录模型填入 `whisper-large-v3-turbo`，速度极快。
- 勾选 **「开启会议回环声音捕获」**：可在听录时同时监听面试官的声音。

### 3. 候选人画像与背景 (Profile)
- **个人简历**：填入你的工作履历、技术栈与代表项目，AI 回答时将自动引证你的真实背景。
- **岗位 JD**：粘贴应聘职位的职责与技能要求，AI 会针对性突出人岗匹配度。
- **STAR 故事库**：填入 2~3 个结构完整的项目破局经历，遇到行为面试题时将自动采用该故事。
- **求职动机 / 薪资预期 / 反问问题**：预先准备好这些模块，现场临场发挥不再卡壳。

---

## 💻 会议软件隐形说明

本工具默认启用 **Windows 窗口保护属性 (`WDA_EXCLUDEFROMCAPTURE`)**：

| 会议软件 | 隐形效果 | 配置建议 |
|---|---|---|
| **腾讯会议 (Tencent Meeting)** | 默认完全隐形 ✅ | 无论是「共享整个桌面」还是「共享应用窗口」，腾讯会议均无法捕获本窗口。 |
| **飞书会议 (Feishu)** | 默认完全隐形 ✅ | 完美隐形，屏幕共享和本地录制均不会出现辅助窗口。 |
| **钉钉会议 (DingTalk)** | 默认完全隐形 ✅ | 完美隐形。 |
| **Zoom** | 需注意一项设置 ⚠️ | 请进入 **Zoom 设置 → 共享屏幕 → 高级 → 屏幕捕获模式**，选择 **「使用窗口过滤的高级捕获 (Advanced capture with window filtering)」**。 |
| **Microsoft Teams** | 默认完全隐形 ✅ | 桌面共享时不会录入该窗口。 |
| **OBS Studio / 本地录屏** | 默认完全隐形 ✅ | 使用 Display Capture (桌面采集) 时该窗口会自动被抠除。 |

---

## 🤝 鸣谢与致敬

- [Blueturboguy07/cue](https://github.com/Blueturboguy07/cue) - 原项目的优秀理念与灵感来源。
- [Wails](https://wails.io/) - 强大优雅的 Go + Web 混合桌面应用开发框架。
- [Vue.js](https://vuejs.org/) - 简单灵活的前端框架。

## 📄 License

GPL-3.0-or-later License
