package watcher

import (
	"github.com/vlorc/hprose-gateway-core/driver"
	"github.com/vlorc/hprose-gateway-types"
	"io"
)

func NewSnapshotWatcher(
	in types.NamedWatcher,
	router types.NamedRouter,
	manager types.SourceManger,
	source types.Source,
	watcher types.NamedWatcher) types.NamedWatcher {
	snapshot := &SnapshotWatcher{
		NamedWatcher: in,
		router:       router,
		manager:      manager.Append("", source),
		watcher:      watcher,
	}
	snapshot.init()
	return snapshot
}

func (s *SnapshotWatcher) Close() error {
	err := s.NamedWatcher.Close()
	s.watcher.Close()
	return err
}

func (s *SnapshotWatcher) init() {
	go func() {
		for s.__updates() != io.EOF {

		}
	}()
}

func (s *SnapshotWatcher) getSource(service *types.Service) types.Source {
	factory := driver.Query(service.Driver)
	source := factory.Instance(s.router, s.manager)
	source.SetService(service)
	return source
}

func (s *SnapshotWatcher) __updates() error {
	updates, err := s.Pop()
	if nil != err {
		return err
	}
	for _, up := range updates {
		source := s.manager.Resolver(up.Id)
		switch up.Op {
		case types.Add:
			if nil != source {
				source.SetService(up.Service)
			} else {
				source = s.getSource(up.Service)
				s.manager.Append(up.Id, source)
				s.router.Append(types.Named{Id: up.Service.Id, Path: up.Service.Path}, up.Id)
			}
			s.watcher.Push([]types.Update{up})
		case types.Delete:
			if nil != source {
				s.router.Remove(types.Named{Id: source.Service().Id, Path: source.Service().Path}, up.Id)
				s.manager.Remove(up.Id)
				driver.Query(source.Service().Driver).Release(source)
				s.watcher.Push([]types.Update{up})
			}
		}
	}
	return nil
}
