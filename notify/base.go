package notify

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ---- Event Definitions ----
type Event interface {
	Name() string
	Data() any
}

type BaseEvent struct {
	typeName string
	data     any
}

func (e *BaseEvent) Name() string { return e.typeName }
func (e *BaseEvent) Data() any    { return e.data }

// ---- Trigger Interface ----
type Trigger interface {
	Check() (Event, bool)
}

// ---- Listener Interface ----
type Listener interface {
	Start()
	Name() string
}

// ---- Event Handler Interface ----
type EventHandler interface {
	Handle(event Event)
}

// ---- Notifier Interface ----
type Notifier interface {
	Notify(userID string, content string)
}

// ---- Engine ----
type Engine struct {
	listeners []Listener
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) RegisterListener(l Listener) {
	e.listeners = append(e.listeners, l)
}

func (e *Engine) StartAll() {
	for _, l := range e.listeners {
		go l.Start()
	}
}

// ---- Notifier Implementation ----
type EmailNotifier struct{}

func (e *EmailNotifier) Notify(userID string, content string) {
	fmt.Printf("[email notify]] user[%s]: %s\n", userID, content)
}

// ---- Event Handler Implementation ----
type DefaultHandler struct {
	Subscribers []string
	Notifier    Notifier
}

func (h *DefaultHandler) Handle(event Event) {
	for _, user := range h.Subscribers {
		h.Notifier.Notify(user, fmt.Sprintf("event: %s, data: %v", event.Name(), event.Data()))
	}
}

// ---- URL Change Trigger ----
type URLChangeTrigger struct {
	URL      string
	lastHash string
}

func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func hash(content string) string {
	h := md5.Sum([]byte(content))
	return hex.EncodeToString(h[:])
}

func (t *URLChangeTrigger) Check() (Event, bool) {
	body := httpGet(t.URL)
	newHash := hash(body)
	if newHash != t.lastHash {
		t.lastHash = newHash
		return &BaseEvent{"URLChanged", t.URL}, true
	}
	return nil, false
}

// ---- URL Listener ----
type URLListener struct {
	trigger  Trigger
	handler  EventHandler
	interval time.Duration
}

func NewURLListener() *URLListener {
	handler := DefaultHandler{
		Subscribers: []string{"user1", "user2"},
		Notifier:    &EmailNotifier{},
	}

	urlTrigger := &URLChangeTrigger{URL: "https://example.com"}
	l := &URLListener{
		trigger:  urlTrigger,
		interval: 30 * time.Minute,
		handler:  &handler,
	}
	return l
}

func (l *URLListener) Start() {
	for {
		if event, ok := l.trigger.Check(); ok {
			l.handler.Handle(event)
		}
		time.Sleep(l.interval)
	}
}

func (l *URLListener) Name() string { return "URLListener" }

// ---- Festival Trigger ----
type FestivalTrigger struct {
	dates []time.Time
}

func (f *FestivalTrigger) Check() (Event, bool) {
	today := time.Now()
	for _, d := range f.dates {
		diff := d.Sub(today).Hours() / 24
		if diff <= 3 && diff >= 0 {
			return &BaseEvent{"FestivalComing", d.Format("2006-01-02")}, true
		}
	}
	return nil, false
}

// ---- Festival Listener ----
type FestivalListener struct {
	trigger Trigger
	handler EventHandler
}

func (l *FestivalListener) Start() {
	for {
		if event, ok := l.trigger.Check(); ok {
			l.handler.Handle(event)
		}
		time.Sleep(24 * time.Hour)
	}
}

func (l *FestivalListener) Name() string { return "FestivalListener" }
