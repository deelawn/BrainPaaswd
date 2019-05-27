package models

// Group is a representation of a user as stored in a /etc/group file
type Group struct {
	Name    string `json:"name"`  // Use int64 for GID even if it is only stored as int32
	GID     int64  `json:"gid"`   // in order to maintain forward compatibility.
	Members []string `json:"members"`
}