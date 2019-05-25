package users

// User is a representation of a user as stored in a /etc/passwd file
type User struct {
	Name    string `json:"name"`
	UID     int64  `json:"uid"` // Use int64 for UID and GID even if they are only stored as int32
	GID     int64  `json:"gid"` // in order to maintain forward compatibility.
	Comment string `json:"comment"`
	Home    string `json:"home"`
	Shell   string `json:"shell"`
}
