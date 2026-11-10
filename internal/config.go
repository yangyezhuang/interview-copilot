package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type ModelPair struct {
	Fast  string `json:"fast"`
	Smart string `json:"smart"`
}

type Settings struct {
	Provider         string               `json:"provider"`
	BaseUrl          string               `json:"baseUrl"` // global/custom base url fallback
	BaseUrls         map[string]string    `json:"baseUrls"` // per-provider base url
	ApiKeys          map[string]string    `json:"apiKeys"`
	Models           map[string]ModelPair `json:"models"`
	Smart            bool                 `json:"smart"`
	SttProvider        string               `json:"sttProvider"`
	SttLocalModel      string               `json:"sttLocalModel"`
	SttMirrorSource    string               `json:"sttMirrorSource"`    // "hf-mirror", "huggingface", "custom"
	SttCustomMirrorUrl string               `json:"sttCustomMirrorUrl"` // custom mirror url
	SttWhisperPath     string               `json:"sttWhisperPath"`
	SttKey             string               `json:"sttKey"`
	SttBaseUrl       string               `json:"sttBaseUrl"`
	SttModel         string               `json:"sttModel"`
	SttLanguage      string               `json:"sttLanguage"`
	MeetingAudio     bool                 `json:"meetingAudio"`
	Opacity          float64              `json:"opacity"`
	ScreenProtection bool                 `json:"screenProtection"`

	// Candidate profile & context
	ResumeText     string `json:"resumeText"`
	JobDescription string `json:"jobDescription"`
	StarStories    string `json:"starStories"`
	WhyCompany     string `json:"whyCompany"`
	WhyLeaving     string `json:"whyLeaving"`
	WorkStyle      string `json:"workStyle"`
	SalaryTarget   string `json:"salaryTarget"`
	QuestionsToAsk string `json:"questionsToAsk"`
	AiRules        string `json:"aiRules"`

	// Window
	WindowX int `json:"windowX"`
	WindowY int `json:"windowY"`
}

func DefaultSettings() Settings {
	return Settings{
		Provider: "deepseek",
		BaseUrl:  "",
		BaseUrls: map[string]string{
			"deepseek":    "https://api.deepseek.com/v1",
			"qwen":        "https://dashscope.aliyuncs.com/compatible-mode/v1",
			"doubao":      "https://ark.cn-beijing.volces.com/api/v3",
			"siliconflow": "https://api.siliconflow.cn/v1",
			"moonshot":    "https://api.moonshot.cn/v1",
			"zhipu":       "https://open.bigmodel.cn/api/paas/v4",
			"openai":      "https://api.openai.com/v1",
			"ollama":      "http://localhost:11434/v1",
			"groq":        "https://api.groq.com/openai/v1",
			"custom":      "https://api.openai.com/v1",
		},
		ApiKeys: map[string]string{
			"deepseek":    "",
			"qwen":        "",
			"doubao":      "",
			"siliconflow": "",
			"moonshot":    "",
			"zhipu":       "",
			"openai":      "",
			"ollama":      "",
			"anthropic":   "",
			"gemini":      "",
			"groq":        "",
			"custom":      "",
		},
		Models: map[string]ModelPair{
			"deepseek": {
				Fast:  "deepseek-chat",
				Smart: "deepseek-reasoner",
			},
			"qwen": {
				Fast:  "qwen-plus",
				Smart: "qwen-max",
			},
			"doubao": {
				Fast:  "doubao-1.5-lite-32k",
				Smart: "doubao-1.5-pro-32k",
			},
			"siliconflow": {
				Fast:  "deepseek-ai/DeepSeek-V3",
				Smart: "deepseek-ai/DeepSeek-R1",
			},
			"moonshot": {
				Fast:  "moonshot-v1-8k",
				Smart: "moonshot-v1-32k",
			},
			"zhipu": {
				Fast:  "glm-4-flash",
				Smart: "glm-4-plus",
			},
			"openai": {
				Fast:  "gpt-4o-mini",
				Smart: "gpt-4o",
			},
			"ollama": {
				Fast:  "qwen2.5:7b",
				Smart: "deepseek-r1:7b",
			},
			"anthropic": {
				Fast:  "claude-3-5-haiku-20241022",
				Smart: "claude-3-5-sonnet-20241022",
			},
			"gemini": {
				Fast:  "gemini-2.5-flash",
				Smart: "gemini-1.5-pro",
			},
			"groq": {
				Fast:  "llama-3.1-8b-instant",
				Smart: "llama-3.3-70b-versatile",
			},
			"custom": {
				Fast:  "default",
				Smart: "default",
			},
		},
		Smart:              false,
		SttProvider:        "local",
		SttLocalModel:      "base",
		SttMirrorSource:    "hf-mirror",
		SttCustomMirrorUrl: "",
		SttWhisperPath:     "",
		SttKey:           "",
		SttBaseUrl:       "",
		SttModel:         "whisper-large-v3-turbo",
		SttLanguage:      "zh",
		MeetingAudio:     true,
		Opacity:          0.96,
		ScreenProtection: true,
		WindowX:          -1,
		WindowY:          -1,
	}
}

type ConfigStore struct {
	mu       sync.RWMutex
	filePath string
	settings Settings
}

func NewConfigStore() *ConfigStore {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "interview-copilot")
	_ = os.MkdirAll(appDir, 0755)
	store := &ConfigStore{
		filePath: filepath.Join(appDir, "settings.json"),
		settings: DefaultSettings(),
	}
	store.load()
	return store
}

func (cs *ConfigStore) load() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	data, err := os.ReadFile(cs.filePath)
	if err != nil {
		return
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err == nil {
		def := DefaultSettings()
		if s.BaseUrls == nil {
			s.BaseUrls = def.BaseUrls
		} else {
			for k, v := range def.BaseUrls {
				if _, ok := s.BaseUrls[k]; !ok {
					s.BaseUrls[k] = v
				}
			}
		}
		if s.ApiKeys == nil {
			s.ApiKeys = def.ApiKeys
		} else {
			for k, v := range def.ApiKeys {
				if _, ok := s.ApiKeys[k]; !ok {
					s.ApiKeys[k] = v
				}
			}
		}
		if s.Models == nil {
			s.Models = def.Models
		} else {
			for k, v := range def.Models {
				if _, ok := s.Models[k]; !ok {
					s.Models[k] = v
				}
			}
		}
		if s.Provider == "" {
			s.Provider = def.Provider
		}
		if s.SttMirrorSource == "" {
			s.SttMirrorSource = def.SttMirrorSource
		}
		if s.Opacity <= 0 {
			s.Opacity = def.Opacity
		}
		cs.settings = s
	}
}

func (cs *ConfigStore) Get() Settings {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.settings
}

func (cs *ConfigStore) Save(s Settings) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.settings = s
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cs.filePath, data, 0644)
}
