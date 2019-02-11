package watcher

import (
	"github.com/vlorc/hprose-gateway-types"
)

func (e EmptyWatcher) Push([]types.Update) error {
	return nil
}

func (e EmptyWatcher) Pop() ([]types.Update, error) {
	return nil,nil
}

func (e EmptyWatcher) Close() error {
	return nil
}
