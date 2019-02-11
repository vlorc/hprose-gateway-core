package watcher

import (
	"github.com/vlorc/hprose-gateway-types"
	"io"
)

func NewMultiWatcher(in types.NamedWatcher, out ...types.NamedWatcher) types.NamedWatcher {
	multi := &MultiWatcher{
		NamedWatcher: in,
		watcher:      out,
	}
	multi.init()
	return multi
}

func (m *MultiWatcher) init() {
	go func() {
		for m.__updates() != io.EOF {

		}
	}()
}

func (m *MultiWatcher) __updates() error {
	up, err := m.Pop()
	if nil != err {
		return err
	}
	for _, w := range m.watcher {
		w.Push(up)
	}
	return nil
}

func (m *MultiWatcher) Close() error {
	err := m.NamedWatcher.Close()
	for _, w := range m.watcher {
		w.Close()
	}
	return err
}
