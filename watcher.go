package nasync

import (
	"time"
)

// watcher watches the async.queue channel, and writes the logs to output
func (a *Async) watcher() {
	defer close(a.quit)
	var buf buffer
	for {
		timeout := time.After(time.Second / 2)
		for i := 0; i < a.bufSize; i++ {
			select {
			case req := <-a.taskChan:
				a.flushReq(&buf, req)
			case <-timeout:
				goto ForEnd
			case <-a.quit:
				return

			}
		}
	ForEnd:
		if len(buf.Tasks()) == 0 {
			break
		}
		a.flushBuf(&buf)
	}
}

// flushReq handles the request and writes the result to writer
func (a *Async) flushReq(b *buffer, t *task) {
	//do print for this
	b.Append(t)
}

// flushBuf flushes the content of buffer to out and reset the buffer
func (a *Async) flushBuf(b *buffer) {
	tasks := b.Tasks()
	if len(tasks) > 0 {
		for _, t := range tasks {
			a.wait.Add(1)
			go func(t *task) {
				t.Do()
				a.wait.Done()
			}(t)
		}
		a.wait.Wait()
		b.Reset()
	}
}
