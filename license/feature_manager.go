package license

import "sync"

type FeatureManager struct {
	mu      sync.RWMutex
	licensed bool
	owner    string
}

func (f *FeatureManager) Enable(owner string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.licensed = true
	f.owner = owner
}

func (f *FeatureManager) Disable() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.licensed = false
	f.owner = ""
}

func (f *FeatureManager) Licensed() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.licensed
}

func (f *FeatureManager) Enabled(name string) bool {
	// The RealmsLauncher protocol returns no feature list. A successful
	// RealmsLauncher-compatible license therefore enables the browser family.
	return f.Licensed() && name != ""
}

func (f *FeatureManager) Owner() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.owner
}
