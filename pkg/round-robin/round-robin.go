package round_robin

import (
	"container/list"
	"errors"
	"net/url"
	"sync"
)

var (
	ErrEmptyList = errors.New("no one members are found")
)

type RoundRobin struct {
	mu             sync.Mutex
	members        *list.List
	curr           *list.Element
	healthCallback func(url.URL) bool
}

func NewRoundRobin(f func(url.URL) bool) *RoundRobin {
	l := list.New()
	return &RoundRobin{members: l, healthCallback: f, curr: l.Front()}
}

func (rr *RoundRobin) Next() (url.URL, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if rr.members.Len() == 0 {
		return url.URL{}, ErrEmptyList
	}

	for {
		if !rr.healthCallback(rr.curr.Value.(url.URL)) {
			next := rr.curr.Next()
			rr.members.Remove(rr.curr)
			rr.curr = next
		} else {
			res := rr.curr.Value.(url.URL)
			if rr.curr = rr.curr.Next(); rr.curr == nil {
				rr.curr = rr.members.Front()
			}
			return res, nil
		}
	}

}

func (rr *RoundRobin) Add(next url.URL) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.members.PushFront(next)
}
