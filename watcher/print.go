package watcher

import (
	"encoding/json"
	"github.com/vlorc/hprose-gateway-types"
)

type printWatcher func(string, ...interface{})

func NewPrintWatcher(output func(string, ...interface{})) types.NamedWatcher {
	return printWatcher(output)
}

func (p printWatcher) Push(event []types.Update) error {
	for i := range event {
		buf, _ := json.MarshalIndent(&event[i], "", "    ")
		p("[%s][%s]",event[i].Op.String(), string(buf))
	}
	return nil
}

func (p printWatcher) Pop() ([]types.Update, error) {
	return nil, nil
}

func (p printWatcher) Close() error {
	return nil
}
