export const defaultBaseUrls = {
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

export const defaultApiKeys = {
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

export const defaultModels = {
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

export const providerPresets = {
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

export const modeLabels = {
  say: '💬 我该怎么说',
  assist: '✨ 智能辅助',
  leetcode: '💻 算法解题',
  recap: '📋 会议复盘',
  ask: '❓ 自由提问',
  answerThis: '⚡ 针对提问回答'
}

export const categoryLabels = {
  behavioral: 'STAR 行为题',
  motivation: '求职动机',
  situational: '情境应对',
  experience: '项目经历',
  compensation: '薪资与反问',
  technical: '技术与设计',
  leetcode: '算法编程'
}
