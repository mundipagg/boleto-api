package buffer

import (
	"fmt"
	"sync"
	"time"
)

const (
	DefaultCapacity   = 100
	DefaultOnWait     = 100
	DefaultExpiration = 60000
	DefaultBackoff    = 10000
)

type Buffer interface {
	Write(item interface{})
}

type buffer struct {
	sync.Locker
	cap        int
	size       int
	expiration time.Duration
	chunks     chan entry
	items      []interface{}
	backoff    time.Duration
}

func (b *buffer) Write(item interface{}) {
	b.Lock()
	defer b.Unlock()
	b.items[b.size] = item
	b.size++
	if b.size >= b.cap {
		b.clear()
	}
}

func (b *buffer) clear() {
	if b.size > 0 {
		events := b.items[:b.size]
		b.size = 0
		b.items = make([]interface{}, b.cap)
		go func() {
			b.chunks <- entry{
				items: events,
				retries: cap(b.chunks),
			}
		}()
	}
}

func (b *buffer) watcher() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()
	for {
		time.Sleep(b.expiration)
		b.Lock()
		b.clear()
		b.Unlock()
	}
}

type Config struct {
	Cap        int
	OnWait     int
	Expiration time.Duration
	BackOff    time.Duration
	OnOverflow func([]interface{}) error
}

type entry struct {
	items []interface{}
	retries int
}

func New(c Config) Buffer {
	if c.Cap == 0 {
		c.Cap = DefaultCapacity
	}
	if c.Expiration == 0 {
		c.Expiration = DefaultExpiration
	}
	if c.BackOff == 0 {
		c.BackOff = DefaultBackoff
	}
	if c.OnWait == 0 {
		c.OnWait = DefaultOnWait
	}

	b := &buffer{
		Locker:     &sync.Mutex{},
		size:       0,
		cap:        c.Cap,
		expiration: c.Expiration,
		chunks:     make(chan entry, c.OnWait),
		items:      make([]interface{}, c.Cap),
		backoff:    c.BackOff,
	}
	go b.watcher()
	go b.consumer(c)
	return b
}

func (b *buffer) consumer(c Config) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()
	for events := range b.chunks {
		err := c.OnOverflow(events.items)
		if err != nil {
			go func(events entry) {
					events.retries--
					if events.retries >= 0 {
						time.Sleep(b.backoff)
						b.chunks <- events
					}
			}(events)
		}
	}
}
