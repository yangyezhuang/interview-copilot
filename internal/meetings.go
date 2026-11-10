package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Turn struct {
	Channel   string `json:"channel"` // "you" or "them"
	Text      string `json:"text"`
	Timestamp int64  `json:"ts"`
}

type MeetingSession struct {
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Turns     []Turn `json:"turns"`
}

type MeetingStore struct {
	mu         sync.RWMutex
	filePath   string
	turns      []Turn
	maxTurns   int
	lastActive int64
}

func NewMeetingStore() *MeetingStore {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "interview-copilot")
	_ = os.MkdirAll(appDir, 0755)

	ms := &MeetingStore{
		filePath: filepath.Join(appDir, "meetings.json"),
		turns:    make([]Turn, 0),
		maxTurns: 200,
	}
	ms.loadRecent()
	return ms
}

func (ms *MeetingStore) loadRecent() {
	data, err := os.ReadFile(ms.filePath)
	if err != nil {
		return
	}
	var session MeetingSession
	if err := json.Unmarshal(data, &session); err == nil {
		// Restore if last turn is within 30 minutes
		if len(session.Turns) > 0 {
			lastTurn := session.Turns[len(session.Turns)-1]
			if time.Now().UnixMilli()-lastTurn.Timestamp < 30*60*1000 {
				ms.turns = session.Turns
				ms.lastActive = lastTurn.Timestamp
			}
		}
	}
}

func (ms *MeetingStore) AddTurn(channel string, text string) Turn {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now().UnixMilli()
	turn := Turn{
		Channel:   channel,
		Text:      text,
		Timestamp: now,
	}
	ms.turns = append(ms.turns, turn)
	if len(ms.turns) > ms.maxTurns {
		ms.turns = ms.turns[len(ms.turns)-ms.maxTurns:]
	}
	ms.lastActive = now

	// Save to disk asynchronously
	go ms.saveToDisk()

	return turn
}

func (ms *MeetingStore) GetTurns() []Turn {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	res := make([]Turn, len(ms.turns))
	copy(res, ms.turns)
	return res
}

func (ms *MeetingStore) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.turns = make([]Turn, 0)
	go ms.saveToDisk()
}

func (ms *MeetingStore) saveToDisk() {
	ms.mu.RLock()
	session := MeetingSession{
		EndTime: time.Now().UnixMilli(),
		Turns:   ms.turns,
	}
	if len(ms.turns) > 0 {
		session.StartTime = ms.turns[0].Timestamp
	}
	data, err := json.Marshal(session)
	ms.mu.RUnlock()

	if err == nil {
		_ = os.WriteFile(ms.filePath, data, 0644)
	}
}
