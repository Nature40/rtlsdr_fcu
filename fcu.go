package main

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// FlowCtrlUnit pipes a Reader to a Writer, but terminates prematurely if
// nothing has been sent for longer than a timeout.
type FlowCtrlUnit struct {
	src io.Reader
	dst io.Writer

	timeout time.Duration
	bufSize int

	once      sync.Once
	lastNanos int64
	errChan   chan error
}

// NewFlowCtrlUnit creates and starts a new FlowCtrlUnit.
func NewFlowCtrlUnit(src io.Reader, dst io.Writer, timeout time.Duration, bufSize int) *FlowCtrlUnit {
	fcu := &FlowCtrlUnit{
		src: src,
		dst: dst,

		timeout: timeout,
		bufSize: bufSize,

		lastNanos: time.Now().UnixNano(),
		errChan:   make(chan error),
	}

	go fcu.checkRead()
	go fcu.checkTimeout()

	return fcu
}

// checkRead reads the piped Reader and notifies about data via synChan.
func (fcu *FlowCtrlUnit) checkRead() {
	tee := io.TeeReader(fcu.src, fcu.dst)
	buf := make([]byte, fcu.bufSize)

	for {
		if _, err := tee.Read(buf); err != nil {
			fcu.once.Do(func() { fcu.errChan <- err })
			return
		} else {
			atomic.StoreInt64(&fcu.lastNanos, time.Now().UnixNano())
		}
	}
}

// checkTimeout controls the latest actions and might trigger a timeout.
func (fcu *FlowCtrlUnit) checkTimeout() {
	for range time.Tick(fcu.timeout) {
		lastTime := time.Unix(0, atomic.LoadInt64(&fcu.lastNanos))

		if time.Now().UTC().After(lastTime.Add(fcu.timeout)) {
			fcu.once.Do(func() { fcu.errChan <- fmt.Errorf("timeout") })
			return
		}
	}
}

// Wait until FlowCtrlUnit has finished. This may be due to a timeout or error.
func (fcu *FlowCtrlUnit) Wait() error {
	return <-fcu.errChan
}
