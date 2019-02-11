package watcher

import (
	"github.com/vlorc/hprose-gateway-types"
	"io"
)

func NewChannelWatcher(length int) types.NamedWatcher {
	return ChannelWatcher(make(chan []types.Update, length))
}

func (c ChannelWatcher) Push(event []types.Update) error {
	c <- event
	return nil
}

func (c ChannelWatcher) Pop() ([]types.Update, error) {
	event, ok := <-c
	if !ok {
		return event, io.EOF
	}
	return event, nil
}

func (c ChannelWatcher) Close() error {
	close(c)
	return nil
}
