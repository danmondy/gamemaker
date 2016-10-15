package data

import (
	"fmt"
	"sync"
	"time"
)

var SessionCache *sessionCache

const EXPIRE_DURATION = 1 //minutes

type Session struct {
	Id          string
	LastTouched time.Time
	User        *User
}
type sessionCache struct {
	lock     *sync.RWMutex
	Sessions map[string]Session
}

func (sc *sessionCache) Add(s Session) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	s.LastTouched = time.Now()
	sc.Sessions[s.Id] = s
}
func (sc *sessionCache) Remove(id string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	delete(sc.Sessions, id)
}
func (sc *sessionCache) Get(id string) (Session, bool) {
	sc.lock.Lock()
	defer sc.lock.Unlock()	
	s, ok := sc.Sessions[id]
	if ok {
		s.LastTouched = time.Now()
		fmt.Println("Last Touched:",s.LastTouched)
	}
	return s, ok
}

func NewSessionCache() *sessionCache {
	return &sessionCache{
		&sync.RWMutex{},
		make(map[string]Session, 1000), //1000 users during peak hours (chosen arbitrarily)
	}
}
func init() {
	SessionCache = NewSessionCache()
	go SessionCache.RemoveOldItems() //run forever
}

func (sc *sessionCache) RemoveOldItems() {
	for {
		time.Sleep(time.Minute)
		now := time.Now()
		for _, s := range sc.Sessions {
			if s.LastTouched.Sub(now) > (time.Minute * EXPIRE_DURATION) {
				sc.Remove(s.Id)
				fmt.Println("Carlton removed:", s.User.Email)
			}
		}
	}
}
