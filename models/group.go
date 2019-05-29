package models

import (
	"strings"
)

// Group is a representation of a user as stored in a /etc/group file
type Group struct {
	Name    string                 `json:"name"` // Use int64 for GID even if it is only stored as int32
	GID     int64                  `json:"gid"`  // in order to maintain forward compatibility.
	Members []string               `json:"members"`
	members map[string]interface{} // To prevent iterating over the entire Members array when querying groups my members
}

// AddMember will add a member to the Members list and the members map
func (g *Group) AddMember(member string) {

	member = strings.TrimSpace(member)

	if len(member) == 0 {
		return
	}

	g.Members = append(g.Members, member)

	if g.members == nil {
		g.members = make(map[string]interface{})
	}

	g.members[member] = nil
}

// ContainsMember returns true if the group contains the provided member
func (g *Group) ContainsMember(member string) bool {

	_, exists := g.members[member]
	return exists
}
