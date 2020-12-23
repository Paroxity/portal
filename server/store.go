package server

import (
	"strings"
	"sync"
)

var (
	defaultGroup *Group
	groups       = make(map[string]*Group)
	groupsMu     sync.Mutex
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
	return defaultGroup
}

func SetDefaultGroup(g *Group) {
	defaultGroup = g
}

func GroupFromName(name string) (*Group, bool) {
	g, ok := groups[strings.ToLower(name)]
	return g, ok
}

func Groups() map[string]*Group {
	return groups
}
