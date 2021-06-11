// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errgroup provides synchronization, error propagation, and Context
// cancellation for groups of goroutines working on subtasks of a common task.

package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type fc func(ctx context.Context) error
type IGroup interface {
	maxProcess(int)
	Go(fc)
	Wait() error
}

type CancelGroup struct {
	cancel  func()
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	ctx     context.Context

	ch  chan fc
	chs []fc
}

func NewCancel(ctx context.Context, ops ...IOption) IGroup {
	ctx, cancel := context.WithCancel(ctx)
	g := &CancelGroup{cancel: cancel, ctx: ctx}
	for _, op := range ops {
		op(g)
	}
	return g
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *CancelGroup) Go(f fc) {
	g.wg.Add(1)
	if g.ch == nil {
		go g.do(f)
		return
	}
	select {
	case g.ch <- f:
		go g.do(f)
	default:
		g.chs = append(g.chs, f)
	}
}

func (g *CancelGroup) do(f fc) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<16)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		if g.ch != nil {
			<-g.ch
		}
		g.wg.Done()
	}()
	err = f(g.ctx)
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *CancelGroup) Wait() error {
	if len(g.chs) > 0 {
		for _, f := range g.chs {
			g.ch <- f
			go g.do(f)
		}
	}
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	if g.ch != nil {
		close(g.ch)
	}
	return g.err
}

func (g *CancelGroup) maxProcess(n int) {
	g.ch = make(chan fc, n)
}

type ContinueGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	ctx     context.Context

	ch  chan fc
	chs []fc
}

func NewContinue(ctx context.Context, ops ...IOption) IGroup {
	g := &ContinueGroup{ctx: ctx}
	for _, op := range ops {
		op(g)
	}
	return g
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *ContinueGroup) Go(f fc) {
	g.wg.Add(1)
	if g.ch == nil {
		go g.do(f)
		return
	}
	select {
	case g.ch <- f:
		go g.do(f)
	default:
		g.chs = append(g.chs, f)
	}
}

func (g *ContinueGroup) do(f fc) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1<<16)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() { g.err = err })
		}
		if g.ch != nil {
			<-g.ch
		}
		g.wg.Done()
	}()
	err = f(g.ctx)
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *ContinueGroup) Wait() error {
	if len(g.chs) > 0 {
		for _, f := range g.chs {
			g.ch <- f
			go g.do(f)
		}
	}
	g.wg.Wait()
	if g.ch != nil {
		close(g.ch)
	}
	return g.err
}

func (g *ContinueGroup) maxProcess(n int) {
	g.ch = make(chan fc, n)
}
