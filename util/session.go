package util

import (
	"fmt"
	"sync"
	"time"
)

type sesssionData struct {
	sid        string
	updateTime int64
	data       interface{}
}

// Session provides Session manager
type Session struct {
	data map[string]*sesssionData
	lock sync.Locker
}

// Get the Session
func (s *Session) Get(sid string) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	if d, ok := s.data[sid]; ok {
		return d.data
	}
	return nil
}
func (s *Session) Set(sid string, data interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.Get(sid) == nil {

		sdata := &sesssionData{
			sid:        sid,
			updateTime: time.Now().Unix(),
			data:       data,
		}
		s.data[sid] = sdata

		return nil
	}
	return fmt.Errorf("sid %s exists in the session", sid)

}

func (s *Session) Del(sid string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.data[sid]; ok {
		delete(s.data, sid)
	}
	return nil
}

func NewSession() *Session {
	s := &Session{
		lock: &sync.Mutex{},
	}
	return s
}
