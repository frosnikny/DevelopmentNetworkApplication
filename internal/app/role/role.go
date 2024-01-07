package role

type Role int

const (
	NotAuthorized Role = iota // 0
	Customer                  // 1
	Moderator                 // 2
)
