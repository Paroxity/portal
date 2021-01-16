package server

import (
	"strings"
	"sync"
)

var (
	defaultGroup *Group
	groups       = make(map[string]*Group)
	groupsMu     sync.RWMutex
)

func AddGroup(g *Group) {
	groupsMu.Lock()
	groups[strings.ToLower(g.Name())] = g
	groupsMu.Unlock()
}

func RemoveGroup(name string) {
	groupsMu.Lock()
	delete(groups, strings.ToLower(name))
	groupsMu.Unlock()
}

func DefaultGroup() *Group {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	return defaultGroup
}

func SetDefaultGroup(g *Group) {
	groupsMu.Lock()
	defaultGroup = g
	groupsMu.Unlock()
}

func GroupFromName(name string) (*Group, bool) {
	groupsMu.RLock()
	g, ok := groups[strings.ToLower(name)]
	groupsMu.RUnlock()
	return g, ok
}

func Groups() map[string]*Group {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	return groups
}
