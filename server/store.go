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

// AddGroup stores a group in the groups map to be used for lookups.
func AddGroup(g *Group) {
	groupsMu.Lock()
	groups[strings.ToLower(g.Name())] = g
	groupsMu.Unlock()
}

// RemoveGroup removes a group from the map, meaning it can no longer be used for lookups etc.
func RemoveGroup(name string) {
	groupsMu.Lock()
	delete(groups, strings.ToLower(name))
	groupsMu.Unlock()
}

// DefaultGroup returns the default group that is commonly used for load balancing.
func DefaultGroup() *Group {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	return defaultGroup
}

// SetDefaultGroup sets the default group that is commonly used for loading balancing.
func SetDefaultGroup(g *Group) {
	groupsMu.Lock()
	defaultGroup = g
	groupsMu.Unlock()
}

// GroupFromName attempts to find a group with the name provided. The group and a bool to say if the group is
// found or not is returned.
func GroupFromName(name string) (*Group, bool) {
	groupsMu.RLock()
	g, ok := groups[strings.ToLower(name)]
	groupsMu.RUnlock()
	return g, ok
}

// Groups returns all of the groups registered on the proxy.
func Groups() map[string]*Group {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	return groups
}
