package watcher

import (
	"github.com/vlorc/hprose-gateway-types"
	"sync"
)

type ChannelWatcher chan []types.Update

type EmptyWatcher struct {}

type MultiWatcher struct {
	types.NamedWatcher
	watcher []types.NamedWatcher
}

type SnapshotWatcher struct {
	types.NamedWatcher
	watcher types.NamedWatcher
	router  types.NamedRouter
	manager types.SourceManger
	service sync.Map
}
