package ioms

type sesssionData struct {
	SID        string
	updateTime int64
	data       interface{}
}

// Session provides Session manager
type Session struct {
	data map[string]interface{}
}

// Get the Session
func (s *Session) Get(sid string) interface{} {
	if d, ok := s[sid]; ok {
		return d
	}
	return nil
}
