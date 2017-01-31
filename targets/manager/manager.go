package manager

import (
	"fmt"
	"github.com/soider/d"
	"github.com/soider/schnur/targets"
	"sync"
)

type TargetsManager struct {
	sync.Once
	targets targets.TargetList
	loader  loader
	loadErr error
}

func newNotFoundError(name string) error {
	return fmt.Errorf("Target with name %s not found", name)
}

func (tm *TargetsManager) Target(name string) (targets.Target, error) {
	var r targets.Target
	tm.Do(tm.Load)
	d.D(tm.targets)
	for _, target := range tm.targets {
		if target.Name == name {
			return target, nil
		}
	}
	return r, newNotFoundError(name)
}

func (tm *TargetsManager) Targets() (targets.TargetList, error) {
	tm.Do(tm.Load)
	return tm.targets, tm.loadErr
}

func (tm *TargetsManager) Load() {
	targets, err := tm.loader.Targets()
	if err != nil {
		tm.loadErr = err
		return
	}
	tm.targets = targets
}

type loader interface {
	Targets() (targets.TargetList, error)
}

func New(loader loader) *TargetsManager {
	return &TargetsManager{
		loader: loader,
	}
}
