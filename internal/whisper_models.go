package internal

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	HFBrowseBaseURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main"
	HFMirrorBaseURL = "https://hf-mirror.com/ggerganov/whisper.cpp/resolve/main"
)

type WhisperModelSpec struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	Filename       string `json:"filename"`
	SizeBytes      int64  `json:"sizeBytes"`
	Size           string `json:"size"`
	EnglishOnly    bool   `json:"englishOnly"`
	Speed          int    `json:"speed"`    // 1-5
	Accuracy       int    `json:"accuracy"` // 1-5
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
}

type WhisperModelInfo struct {
	WhisperModelSpec
	Downloaded bool `json:"downloaded"`
}

type WhisperDownloadState struct {
	Model      string  `json:"model"`
	Status     string  `json:"status"` // "idle", "downloading", "completed", "failed", "cancelled"
	Downloaded int64   `json:"downloaded"`
	Total      int64   `json:"total"`
	Progress   float64 `json:"progress"` // 0.0 - 100.0
	Speed      string  `json:"speed"`    // e.g. "4.5 MB/s"
	Error      string  `json:"error"`
}

type WhisperRuntimeStatus struct {
	Available bool   `json:"available"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Error     string `json:"error"`
}

var WhisperCatalog = []WhisperModelSpec{
	{
		Name:           "base",
		DisplayName:    "Whisper Base (多语言)",
		Filename:       "ggml-base.bin",
		SizeBytes:      147951465,
		Size:           "148 MB",
		EnglishOnly:    false,
		Speed:          4,
		Accuracy:       3,
		Description:    "日常首选，平衡了转录速度与准确率，响应迅速，资源占用低",
		Recommendation: "推荐首选",
	},
	{
		Name:           "tiny",
		DisplayName:    "Whisper Tiny (多语言)",
		Filename:       "ggml-tiny.bin",
		SizeBytes:      77691713,
		Size:           "78 MB",
		EnglishOnly:    false,
		Speed:          5,
		Accuracy:       2,
		Description:    "极速轻量，占用资源极小，适合老旧电脑或追求极致低延迟",
		Recommendation: "最快速度",
	},
	{
		Name:           "small",
		DisplayName:    "Whisper Small (多语言)",
		Filename:       "ggml-small.bin",
		SizeBytes:      487601967,
		Size:           "488 MB",
		EnglishOnly:    false,
		Speed:          3,
		Accuracy:       4,
		Description:    "更高准确度，对中英文混合及口音有更好的识别效果",
		Recommendation: "高精度",
	},
	{
		Name:           "medium",
		DisplayName:    "Whisper Medium (多语言)",
		Filename:       "ggml-medium.bin",
		SizeBytes:      1533763059,
		Size:           "1.5 GB",
		EnglishOnly:    false,
		Speed:          2,
		Accuracy:       5,
		Description:    "专业级精度，适合复杂专业术语与嘈杂环境（推荐独立显卡或强劲 CPU）",
		Recommendation: "专业精度",
	},
	{
		Name:           "large-v3-turbo",
		DisplayName:    "Whisper Large v3 Turbo (多语言)",
		Filename:       "ggml-large-v3-turbo.bin",
		SizeBytes:      1629853963,
		Size:           "1.6 GB",
		EnglishOnly:    false,
		Speed:          3,
		Accuracy:       5,
		Description:    "OpenAI 官方 Turbo 优化版，具备顶级 Large v3 精度且推理速度提升近一倍",
		Recommendation: "顶级精度+加速",
	},
	{
		Name:           "tiny.en",
		DisplayName:    "Whisper Tiny (纯英文)",
		Filename:       "ggml-tiny.en.bin",
		SizeBytes:      77704715,
		Size:           "78 MB",
		EnglishOnly:    true,
		Speed:          5,
		Accuracy:       2,
		Description:    "纯英文面试极速专用，资源占用极小",
		Recommendation: "英文极速",
	},
	{
		Name:           "base.en",
		DisplayName:    "Whisper Base (纯英文)",
		Filename:       "ggml-base.en.bin",
		SizeBytes:      147964211,
		Size:           "148 MB",
		EnglishOnly:    true,
		Speed:          4,
		Accuracy:       3,
		Description:    "纯英文面试推荐，速度与准确率良好平衡",
		Recommendation: "英文推荐",
	},
	{
		Name:           "small.en",
		DisplayName:    "Whisper Small (纯英文)",
		Filename:       "ggml-small.en.bin",
		SizeBytes:      487614201,
		Size:           "488 MB",
		EnglishOnly:    true,
		Speed:          3,
		Accuracy:       4,
		Description:    "纯英文高精度转录，准确率优异",
		Recommendation: "英文高精度",
	},
}

func LookupWhisperModel(name string) (WhisperModelSpec, bool) {
	for _, spec := range WhisperCatalog {
		if spec.Name == name {
			return spec, true
		}
	}
	return WhisperModelSpec{}, false
}

type WhisperManager struct {
	appDir         string
	modelsDir      string
	app            *application.App
	mu             sync.Mutex
	isDownloading  bool
	cancelDownload context.CancelFunc
	downloadState  WhisperDownloadState
}

func NewWhisperManager(appDir string) (*WhisperManager, error) {
	modelsDir := filepath.Join(appDir, "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	return &WhisperManager{
		appDir:    appDir,
		modelsDir: modelsDir,
		downloadState: WhisperDownloadState{
			Status: "idle",
		},
	}, nil
}

func (m *WhisperManager) SetApp(app *application.App) {
	m.app = app
}

func (m *WhisperManager) GetModelsDir() string {
	return m.modelsDir
}

func (m *WhisperManager) GetModelPath(modelName string) string {
	spec, ok := LookupWhisperModel(modelName)
	if ok {
		return filepath.Join(m.modelsDir, spec.Filename)
	}
	return filepath.Join(m.modelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
}

func (m *WhisperManager) IsModelDownloaded(modelName string) bool {
	spec, ok := LookupWhisperModel(modelName)
	if !ok {
		return false
	}
	path := filepath.Join(m.modelsDir, spec.Filename)
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	// Check size (allow tolerance of 1KB in case of minor metadata variation)
	diff := info.Size() - spec.SizeBytes
	if diff < -1024 || diff > 1024 {
		return false
	}
	return true
}

func (m *WhisperManager) GetModels() []WhisperModelInfo {
	res := make([]WhisperModelInfo, len(WhisperCatalog))
	for i, spec := range WhisperCatalog {
		res[i] = WhisperModelInfo{
			WhisperModelSpec: spec,
			Downloaded:       m.IsModelDownloaded(spec.Name),
		}
	}
	return res
}

func (m *WhisperManager) GetDownloadState() WhisperDownloadState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.downloadState
}

func (m *WhisperManager) emitState() {
	if m.app != nil {
		m.app.EmitEvent("stt:model_download_progress", m.downloadState)
	}
}

func (m *WhisperManager) StartDownload(modelName string, mirrorSource string, customMirrorUrl string) error {
	spec, ok := LookupWhisperModel(modelName)
	if !ok {
		return fmt.Errorf("未知的模型名称: %s", modelName)
	}

	m.mu.Lock()
	if m.isDownloading {
		m.mu.Unlock()
		return fmt.Errorf("当前已有模型正在下载中，请等待完成或取消")
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.isDownloading = true
	m.cancelDownload = cancel
	m.downloadState = WhisperDownloadState{
		Model:      modelName,
		Status:     "downloading",
		Downloaded: 0,
		Total:      spec.SizeBytes,
		Progress:   0,
		Speed:      "0 MB/s",
		Error:      "",
	}
	m.mu.Unlock()

	m.emitState()

	go m.runDownload(ctx, spec, mirrorSource, customMirrorUrl)
	return nil
}

func (m *WhisperManager) runDownload(ctx context.Context, spec WhisperModelSpec, mirrorSource string, customMirrorUrl string) {
	defer func() {
		m.mu.Lock()
		m.isDownloading = false
		m.cancelDownload = nil
		m.mu.Unlock()
		m.emitState()
	}()

	targetPath := filepath.Join(m.modelsDir, spec.Filename)
	tempPath := filepath.Join(m.modelsDir, spec.Filename+".tmp")
	_ = os.Remove(tempPath)

	var urls []string
	mirrorSource = strings.TrimSpace(strings.ToLower(mirrorSource))
	customMirrorUrl = strings.TrimSpace(customMirrorUrl)

	switch mirrorSource {
	case "custom":
		if customMirrorUrl != "" {
			customBase := strings.TrimRight(customMirrorUrl, "/")
			urls = append(urls, fmt.Sprintf("%s/%s", customBase, spec.Filename))
		}
		urls = append(urls,
			fmt.Sprintf("%s/%s", HFMirrorBaseURL, spec.Filename),
			fmt.Sprintf("%s/%s", HFBrowseBaseURL, spec.Filename),
		)
	case "huggingface":
		urls = append(urls,
			fmt.Sprintf("%s/%s", HFBrowseBaseURL, spec.Filename),
			fmt.Sprintf("%s/%s", HFMirrorBaseURL, spec.Filename),
		)
	case "hf-mirror":
		fallthrough
	default:
		urls = append(urls,
			fmt.Sprintf("%s/%s", HFMirrorBaseURL, spec.Filename),
			fmt.Sprintf("%s/%s", HFBrowseBaseURL, spec.Filename),
		)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 25 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   40 * time.Minute,
	}

	var resp *http.Response
	var fetchErr error

	for _, downloadURL := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
		if err != nil {
			fetchErr = err
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		r, err := client.Do(req)
		if err != nil {
			fetchErr = err
			continue
		}
		if r.StatusCode != http.StatusOK {
			r.Body.Close()
			fetchErr = fmt.Errorf("服务器响应状态码异常: %d (URL: %s)", r.StatusCode, downloadURL)
			continue
		}
		resp = r
		fetchErr = nil
		break
	}

	if fetchErr != nil || resp == nil {
		m.mu.Lock()
		if ctx.Err() == context.Canceled {
			m.downloadState.Status = "cancelled"
		} else {
			m.downloadState.Status = "failed"
			m.downloadState.Error = fmt.Sprintf("下载失败: %v", fetchErr)
		}
		m.mu.Unlock()
		_ = os.Remove(tempPath)
		return
	}
	defer resp.Body.Close()

	total := spec.SizeBytes
	if resp.ContentLength > 0 {
		total = resp.ContentLength
	}

	out, err := os.Create(tempPath)
	if err != nil {
		m.mu.Lock()
		m.downloadState.Status = "failed"
		m.downloadState.Error = fmt.Sprintf("创建临时文件失败: %v", err)
		m.mu.Unlock()
		return
	}

	buf := make([]byte, 64*1024)
	var downloaded int64
	startTime := time.Now()
	lastEmitTime := time.Now()
	lastDownloaded := int64(0)

	for {
		select {
		case <-ctx.Done():
			out.Close()
			_ = os.Remove(tempPath)
			m.mu.Lock()
			m.downloadState.Status = "cancelled"
			m.mu.Unlock()
			return
		default:
		}

		n, rErr := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				out.Close()
				_ = os.Remove(tempPath)
				m.mu.Lock()
				m.downloadState.Status = "failed"
				m.downloadState.Error = fmt.Sprintf("写入文件失败: %v", wErr)
				m.mu.Unlock()
				return
			}
			downloaded += int64(n)

			now := time.Now()
			if now.Sub(lastEmitTime) >= 200*time.Millisecond {
				elapsedSec := now.Sub(lastEmitTime).Seconds()
				var speedStr string
				if elapsedSec > 0 {
					speedMBps := float64(downloaded-lastDownloaded) / (1024 * 1024) / elapsedSec
					speedStr = fmt.Sprintf("%.2f MB/s", speedMBps)
				}
				progress := float64(downloaded) / float64(total) * 100
				if progress > 99.9 {
					progress = 99.9
				}

				m.mu.Lock()
				m.downloadState.Downloaded = downloaded
				m.downloadState.Total = total
				m.downloadState.Progress = progress
				m.downloadState.Speed = speedStr
				m.mu.Unlock()

				m.emitState()

				lastEmitTime = now
				lastDownloaded = downloaded
			}
		}

		if rErr == io.EOF {
			break
		}
		if rErr != nil {
			out.Close()
			_ = os.Remove(tempPath)
			m.mu.Lock()
			if ctx.Err() == context.Canceled {
				m.downloadState.Status = "cancelled"
			} else {
				m.downloadState.Status = "failed"
				m.downloadState.Error = fmt.Sprintf("下载中断: %v", rErr)
			}
			m.mu.Unlock()
			return
		}
	}

	_ = out.Sync()
	out.Close()

	// Verify file size
	if info, err := os.Stat(tempPath); err == nil && info.Size() > 0 {
		_ = os.Remove(targetPath)
		if err := os.Rename(tempPath, targetPath); err != nil {
			m.mu.Lock()
			m.downloadState.Status = "failed"
			m.downloadState.Error = fmt.Sprintf("替换目标文件失败: %v", err)
			m.mu.Unlock()
			return
		}
	} else {
		m.mu.Lock()
		m.downloadState.Status = "failed"
		m.downloadState.Error = "下载的文件为空或损坏"
		m.mu.Unlock()
		return
	}

	totalElapsed := time.Since(startTime).Seconds()
	avgSpeed := ""
	if totalElapsed > 0 {
		avgSpeed = fmt.Sprintf("%.2f MB/s", float64(downloaded)/(1024*1024)/totalElapsed)
	}

	m.mu.Lock()
	m.downloadState.Status = "completed"
	m.downloadState.Downloaded = downloaded
	m.downloadState.Total = total
	m.downloadState.Progress = 100.0
	m.downloadState.Speed = avgSpeed
	m.downloadState.Error = ""
	m.mu.Unlock()
}

func (m *WhisperManager) CancelDownload(modelName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.isDownloading {
		return nil
	}
	if modelName != "" && m.downloadState.Model != modelName {
		return fmt.Errorf("当前正在下载的模型是 %s，而非 %s", m.downloadState.Model, modelName)
	}
	if m.cancelDownload != nil {
		m.cancelDownload()
	}
	return nil
}

func (m *WhisperManager) DeleteModel(modelName string) error {
	m.mu.Lock()
	if m.isDownloading && m.downloadState.Model == modelName {
		m.mu.Unlock()
		return fmt.Errorf("该模型正在下载中，无法删除")
	}
	m.mu.Unlock()

	spec, ok := LookupWhisperModel(modelName)
	if !ok {
		return fmt.Errorf("未知的模型名称: %s", modelName)
	}
	path := filepath.Join(m.modelsDir, spec.Filename)
	_ = os.Remove(path)
	_ = os.Remove(path + ".tmp")
	return nil
}

func (m *WhisperManager) FindWhisperBinary(customPath string) string {
	// 1. User custom path
	if customPath != "" {
		if info, err := os.Stat(customPath); err == nil && !info.IsDir() {
			return customPath
		}
		if p, err := exec.LookPath(customPath); err == nil {
			return p
		}
	}

	names := []string{"whisper-cli.exe", "whisper-cli", "whisper-server.exe", "whisper-server", "main.exe"}

	// 2. App bin directory
	for _, name := range names {
		p := filepath.Join(m.appDir, "bin", name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}

	// 3. Current executable directory / bin
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		for _, name := range names {
			p1 := filepath.Join(execDir, name)
			if info, err := os.Stat(p1); err == nil && !info.IsDir() {
				return p1
			}
			p2 := filepath.Join(execDir, "bin", name)
			if info, err := os.Stat(p2); err == nil && !info.IsDir() {
				return p2
			}
		}
	}

	// 4. Relative ./bin or ./
	for _, name := range names {
		p1 := filepath.Join("bin", name)
		if info, err := os.Stat(p1); err == nil && !info.IsDir() {
			return p1
		}
		if info, err := os.Stat(name); err == nil && !info.IsDir() {
			return name
		}
	}

	// 5. Common Windows local app paths
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		for _, name := range names {
			p := filepath.Join(localAppData, "Programs", "whisper.cpp", name)
			if info, err := os.Stat(p); err == nil && !info.IsDir() {
				return p
			}
		}
	}
	if progFiles := os.Getenv("ProgramFiles"); progFiles != "" {
		for _, name := range names {
			p := filepath.Join(progFiles, "whisper.cpp", "bin", name)
			if info, err := os.Stat(p); err == nil && !info.IsDir() {
				return p
			}
		}
	}

	// 6. System PATH
	for _, name := range names {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}

	return ""
}

func (m *WhisperManager) CheckRuntime(customPath string) WhisperRuntimeStatus {
	binPath := m.FindWhisperBinary(customPath)
	if binPath == "" {
		return WhisperRuntimeStatus{
			Available: false,
			Path:      "",
			Error:     "未检测到 whisper-cli。请将 whisper-cli.exe 放入程序 bin 目录或配置路径",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "-h")
	HideWindow(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil && ctx.Err() == context.DeadlineExceeded {
		return WhisperRuntimeStatus{
			Available: false,
			Path:      binPath,
			Error:     "运行检测超时",
		}
	}

	version := "whisper.cpp"
	outStr := string(out)
	if strings.Contains(outStr, "whisper") {
		version = "whisper-cli (ready)"
	}

	return WhisperRuntimeStatus{
		Available: true,
		Path:      binPath,
		Version:   version,
	}
}
