package main

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// FlowCtrlUnit pipes a Reader to a Writer, but terminates prematurely if
// nothing has been sent for longer than a timeout.
type FlowCtrlUnit struct {
	src io.Reader
	dst io.Writer

	timeout time.Duration
	bufSize int

	once    sync.Once
	errChan chan error
	synChan chan struct{}
}

// NewFlowCtrlUnit creates and starts a new FlowCtrlUnit.
func NewFlowCtrlUnit(src io.Reader, dst io.Writer, timeout time.Duration, bufSize int) *FlowCtrlUnit {
	fcu := &FlowCtrlUnit{
		src: src,
		dst: dst,

		timeout: timeout,
		bufSize: bufSize,

		errChan: make(chan error),
		synChan: make(chan struct{}),
	}

	go fcu.checkRead()
	go fcu.checkTimeout()

	return fcu
}

// checkRead reads the piped Reader and notifies about data via synChan.
func (fcu *FlowCtrlUnit) checkRead() {
	tee := io.TeeReader(fcu.src, fcu.dst)

	for {
		buf := make([]byte, fcu.bufSize)

		if _, err := tee.Read(buf); err != nil {
			fcu.once.Do(func() { fcu.errChan <- err })
			return
		} else {
			fcu.synChan <- struct{}{}
		}
	}
}

// checkTimeout controls the latest actions and might trigger a timeout.
func (fcu *FlowCtrlUnit) checkTimeout() {
	lastTime := time.Now()

	for {
		select {
		case <-fcu.synChan:
			lastTime = time.Now()

		case <-time.Tick(fcu.timeout):
			if time.Now().After(lastTime.Add(fcu.timeout)) {
				fcu.once.Do(func() { fcu.errChan <- fmt.Errorf("timeout") })
				return
			}
		}
	}
}

// Wait until FlowCtrlUnit has finished. This may be due to a timeout or error.
func (fcu *FlowCtrlUnit) Wait() error {
	return <-fcu.errChan
}
