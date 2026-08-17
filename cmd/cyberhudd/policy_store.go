package main

import (
	"encoding/json"
	"sync"
)

// PolicyStore holds all mode policies in memory. It is the single source of
// truth for current policy state. Updated whenever any mode's SetPolicy is called.
// Keyed by mode ID only — policy is global per mode, not per-region.
type PolicyStore struct {
	mu       sync.RWMutex
	policies map[string]json.RawMessage // modeID → policy JSON
}

// NewPolicyStore creates an initialized PolicyStore ready for use.
func NewPolicyStore() *PolicyStore {
	return &PolicyStore{
		policies: make(map[string]json.RawMessage),
	}
}

// Get returns the saved policy JSON for a mode, or nil if none.
func (ps *PolicyStore) Get(modeID string) json.RawMessage {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.policies[modeID]
}

// Set stores a policy snapshot for a mode.
func (ps *PolicyStore) Set(modeID string, data json.RawMessage) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.policies == nil {
		ps.policies = make(map[string]json.RawMessage)
	}
	ps.policies[modeID] = data
}

// All returns a copy of all stored policies (for serialization to disk).
func (ps *PolicyStore) All() map[string]json.RawMessage {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	cp := make(map[string]json.RawMessage, len(ps.policies))
	for k, v := range ps.policies {
		cp[k] = v
	}
	return cp
}

// LoadFromConfig populates the store from a parsed config file.
func (ps *PolicyStore) LoadFromConfig(cfg *fileConfig) {
	if cfg == nil {
		return
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.policies == nil {
		ps.policies = make(map[string]json.RawMessage)
	}
	for modeID, data := range cfg.Policies {
		ps.policies[modeID] = data
	}
}
