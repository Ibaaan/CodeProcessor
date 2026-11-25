package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"time"
)

type Session interface {
	Set(key, value any) error
	Get(key any) any
	Delete(key any) error
	SessionID() string
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime time.Duration)
}

type Manager struct {
	provider    Provider
	maxlifetime time.Duration
}

func NewManager(provider Provider, gcInterval, maxlifetime time.Duration) *Manager {
	m := &Manager{provider: provider, maxlifetime: maxlifetime}
	go m.GC(gcInterval)
	return m
}

func (manager *Manager) SessionGetOrCreate(token string) (Session, error) {
	if token == "" {
		sid, err := newSessionId()
		if err != nil {
			return nil, err
		}
		session, err := manager.provider.SessionInit(sid)
		if err != nil {
			return nil, err
		}
		return session, nil
	}

	sid, err := url.QueryUnescape(token)
	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	session, err := manager.provider.SessionRead(sid)
	if err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	}

	return session, nil
}

func newSessionId() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (manager *Manager) GC(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		manager.provider.SessionGC(manager.maxlifetime)
	}
}
