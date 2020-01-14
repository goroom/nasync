package nasync

import (
	"sync"
	"time"
)

const (
	//DefaultReqSize Default max goroutine created
	DefaultReqSize = 1000
	//DefaultBufSize Default task when on goroutine
	DefaultBufSize = 1000
)

//DefaultAsync default instance when you run Do(...)
var DefaultAsync *Async

//Do use nasync do some functions
func Do(handler interface{}, params ...interface{}) {
	if DefaultAsync == nil {
		DefaultAsync = New(DefaultReqSize, DefaultBufSize)
	}
	DefaultAsync.Do(handler, params...)
}

//Async  async model
type Async struct {
	quit      chan bool  // quit signal for the watcher to quit
	taskChan  chan *task // queue used in non-runtime  tasks
	bufSize   int
	forceQuit chan bool
	wait      *sync.WaitGroup
}

//New custom your async
func New(ReqSize int, BufSzie int) *Async {
	as := Async{
		quit:      make(chan bool),
		forceQuit: make(chan bool),
		taskChan:  make(chan *task, ReqSize),
		bufSize:   BufSzie,
		wait:      &sync.WaitGroup{},
	}

	go as.watcher()
	return &as
}

//Do do some functions
func (a *Async) Do(handler interface{}, params ...interface{}) {
	t := newTask(handler, params...)
	a.taskChan <- t
}

// Close sends quit signal to watcher and releases all the resources.
// Wait for all tasks complete to close
func (a *Async) Close(wait bool) {
	for {
		select {
		case <-a.quit:
			return
		case <-a.forceQuit:
			close(a.quit)
			return
		default:
			if !wait {
				close(a.quit)
			}
			time.Sleep(time.Second/2)
		}
	}

}
