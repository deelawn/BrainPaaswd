package models

// Group is a representation of a user as stored in a /etc/group file
type Group struct {
	Name    string `json:"name"`  // Use int64 for GID even if it is only stored as int32
	GID     int64  `json:"gid"`   // in order to maintain forward compatibility.
	Members []string `json:"members"`
	members map[string]interface{} // To prevent iterating over the entire Members array when querying groups my members
}

func (g *Group) AddMember(member string) {

	g.Members = append(g.Members, member)

	if g.members == nil {
		g.members = make(map[string]interface{})
	}
	
	g.members[member] = nil
}

func (g *Group) ContainsMember(member string) bool {

	_, exists := g.members[member]
	return exists
}