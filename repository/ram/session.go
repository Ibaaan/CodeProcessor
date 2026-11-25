package ram

import (
	"codeproccessor/session"
	"container/list"
	"sync"
	"time"
)

var pder = &MemoryProvider{list: list.New(), sessions: make(map[string]*list.Element, 0)}

type SessionImpl struct {
	sid          string
	timeAccessed time.Time
	value        map[any]any
}

func (s *SessionImpl) Set(key, value any) error {
	s.value[key] = value
	pder.SessionUpdate(s.sid)
	return nil
}

func (s *SessionImpl) Get(key any) any {
	pder.SessionUpdate(s.sid)
	if v, ok := s.value[key]; ok {
		return v
	}
	return nil
}

func (s *SessionImpl) Delete(key any) error {
	delete(s.value, key)
	pder.SessionUpdate(s.sid)
	return nil
}

func (s *SessionImpl) SessionID() string {
	return s.sid
}

var _ session.Provider = (*MemoryProvider)(nil)

type MemoryProvider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (pder *MemoryProvider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[any]any, 0)
	newsess := &SessionImpl{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *MemoryProvider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionImpl), nil
	}
	sess, err := pder.SessionInit(sid)
	return sess, err
}

func (pder *MemoryProvider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *MemoryProvider) SessionGC(maxlifetime time.Duration) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		session := element.Value.(*SessionImpl)
		if time.Since(session.timeAccessed) <= maxlifetime {
			return
		}

		pder.list.Remove(element)
		delete(pder.sessions, session.sid)
	}
}

func (pder *MemoryProvider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionImpl).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func GetProvider() *MemoryProvider {
	return pder
}
