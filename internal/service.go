package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppService struct {
	app            *application.App
	window         *application.WebviewWindow
	configStore    *ConfigStore
	meetingStore   *MeetingStore
	whisperManager *WhisperManager
	sttClient      *STTClient
	llmClient      *LLMClient
	mu             sync.Mutex
	isBusy         bool
	cancelStream   context.CancelFunc
}

func NewAppService(cs *ConfigStore, ms *MeetingStore) *AppService {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "interview-copilot")
	wm, _ := NewWhisperManager(appDir)

	return &AppService{
		configStore:    cs,
		meetingStore:   ms,
		whisperManager: wm,
		sttClient:      NewSTTClient(),
		llmClient:      NewLLMClient(),
	}
}

func (s *AppService) SetApp(app *application.App) {
	s.app = app
	if s.whisperManager != nil {
		s.whisperManager.SetApp(app)
	}
}

func (s *AppService) SetWindow(w *application.WebviewWindow) {
	s.window = w
}

func (s *AppService) GetSettings() Settings {
	return s.configStore.Get()
}

func (s *AppService) SaveSettings(settings Settings) error {
	err := s.configStore.Save(settings)
	if err != nil {
		return err
	}
	// Update window protection
	if s.window != nil {
		if hwnd, err := s.window.NativeWindowHandle(); err == nil && hwnd != 0 {
			SetWindowProtection(hwnd, settings.ScreenProtection)
		}
	}
	return nil
}

func (s *AppService) CaptureScreen() (string, error) {
	return CaptureScreen()
}

func (s *AppService) GetMeetingHistory() []Turn {
	return s.meetingStore.GetTurns()
}

func (s *AppService) ClearTranscript() {
	s.meetingStore.Clear()
	s.app.EmitEvent("transcript:cleared", nil)
}

func (s *AppService) ProcessAudio(channel string, base64PCM string) (string, error) {
	if base64PCM == "" {
		return "", nil
	}

	pcmBytes, err := base64.StdEncoding.DecodeString(base64PCM)
	if err != nil {
		return "", fmt.Errorf("invalid base64 pcm audio: %w", err)
	}
	if len(pcmBytes) < 3200 { // less than 0.1s
		return "", nil
	}

	wavBytes, err := PCMToWAV(pcmBytes, 16000, 1, 16)
	if err != nil {
		return "", fmt.Errorf("failed to encode WAV: %w", err)
	}

	settings := s.configStore.Get()
	text, err := s.sttClient.Transcribe(wavBytes, settings, s.whisperManager)
	if err != nil {
		s.app.EmitEvent("stt:error", map[string]string{
			"channel": channel,
			"error":   err.Error(),
		})
		return "", err
	}

	text = strings.TrimSpace(text)
	if text != "" {
		turn := s.meetingStore.AddTurn(channel, text)
		s.app.EmitEvent("transcript:turn", turn)
	}

	return text, nil
}

func (s *AppService) RunFeature(mode string, userText string) error {
	s.mu.Lock()
	if s.isBusy {
		// Cancel existing stream if running
		if s.cancelStream != nil {
			s.cancelStream()
		}
	}
	s.isBusy = true
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelStream = cancel
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.isBusy = false
		s.cancelStream = nil
		s.mu.Unlock()
	}()

	def, exists := Modes[mode]
	if !exists {
		def = Modes["ask"]
	}

	turns := s.meetingStore.GetTurns()
	if def.TranscriptRequired && len(turns) == 0 {
		errMsg := "当前尚未转录任何对话，请先开启听录或进行对话后再试。"
		s.app.EmitEvent("llm:error", map[string]string{"message": errMsg})
		return fmt.Errorf(errMsg)
	}

	settings := s.configStore.Get()
	category := DetectCategory(turns)
	if mode == "leetcode" {
		category = "leetcode"
	}

	userBubble := userText
	if mode == "say" {
		userBubble = "我该怎么说？"
	} else if mode == "assist" {
		userBubble = "智能辅助 (Smart Assist)"
	} else if mode == "recap" {
		userBubble = "会议与面试复盘总结"
	} else if mode == "leetcode" {
		userBubble = "解答当前屏幕中的算法编程题"
	} else if mode == "answerThis" {
		userBubble = fmt.Sprintf("针对此题回答: %s", userText)
	}

	s.app.EmitEvent("llm:start", map[string]interface{}{
		"mode":       mode,
		"userBubble": userBubble,
		"category":   category,
	})

	var imageDataUrl string
	if def.NeedsScreen {
		img, err := CaptureScreen()
		if err != nil {
			s.app.EmitEvent("llm:status", map[string]string{
				"message": fmt.Sprintf("屏幕截图失败: %v", err),
			})
		} else {
			imageDataUrl = img
		}
	}

	ctxBlock := BuildInterviewContext(settings, mode, turns)
	sysPrompt := def.SystemPrompt(ctxBlock, settings.AiRules)
	userPrompt := def.BuildUserPrompt(turns, userText)

	s.llmClient.Stream(ctx, settings, sysPrompt, userPrompt, imageDataUrl, LLMStreamCallbacks{
		OnThinking: func(thought string) {
			s.app.EmitEvent("llm:thinking", map[string]string{"thought": thought})
		},
		OnToken: func(token string) {
			s.app.EmitEvent("llm:token", map[string]string{"token": token})
		},
		OnDone: func(fullText string) {
			s.app.EmitEvent("llm:done", map[string]string{"fullText": fullText})
		},
		OnError: func(err error) {
			if ctx.Err() != context.Canceled {
				s.app.EmitEvent("llm:error", map[string]string{"message": err.Error()})
			}
		},
	})

	return nil
}

func (s *AppService) SetProtection(enable bool) bool {
	if s.window == nil {
		return false
	}
	hwnd, err := s.window.NativeWindowHandle()
	if err != nil || hwnd == 0 {
		return false
	}
	ok := SetWindowProtection(hwnd, enable)
	settings := s.configStore.Get()
	settings.ScreenProtection = enable
	_ = s.configStore.Save(settings)
	return ok
}

func (s *AppService) ToggleHide() {
	s.app.EmitEvent("window:toggle_hide", nil)
}

func (s *AppService) QuitApp() {
	s.app.Quit()
}

func (s *AppService) GetWhisperModels() []WhisperModelInfo {
	if s.whisperManager == nil {
		return nil
	}
	return s.whisperManager.GetModels()
}

func (s *AppService) DownloadWhisperModel(modelName string) error {
	if s.whisperManager == nil {
		return fmt.Errorf("本地模型管理器未初始化")
	}
	settings := s.configStore.Get()
	return s.whisperManager.StartDownload(modelName, settings.SttMirrorSource, settings.SttCustomMirrorUrl)
}

func (s *AppService) CancelWhisperModelDownload(modelName string) error {
	if s.whisperManager == nil {
		return fmt.Errorf("本地模型管理器未初始化")
	}
	return s.whisperManager.CancelDownload(modelName)
}

func (s *AppService) DeleteWhisperModel(modelName string) error {
	if s.whisperManager == nil {
		return fmt.Errorf("本地模型管理器未初始化")
	}
	return s.whisperManager.DeleteModel(modelName)
}

func (s *AppService) GetWhisperModelDownloadState() WhisperDownloadState {
	if s.whisperManager == nil {
		return WhisperDownloadState{Status: "idle"}
	}
	return s.whisperManager.GetDownloadState()
}

func (s *AppService) CheckWhisperRuntime() WhisperRuntimeStatus {
	if s.whisperManager == nil {
		return WhisperRuntimeStatus{Available: false, Error: "本地模型管理器未初始化"}
	}
	settings := s.configStore.Get()
	return s.whisperManager.CheckRuntime(settings.SttWhisperPath)
}

func (s *AppService) GetModelsDirectory() string {
	if s.whisperManager == nil {
		return ""
	}
	return s.whisperManager.GetModelsDir()
}
