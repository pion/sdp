package sdp

type ConnectionRole int

const (
	ConnectionRoleActive ConnectionRole = iota + 1
	ConnectionRolePassive
	ConnectionRoleActPass
	ConnectionRoleHoldConn
)

const (
	connectionRoleActiveStr   = "active"
	connectionRolePassiveStr  = "passive"
	connectionRoleActPassStr  = "actpass"
	connectionRoleHoldConnStr = "holdconn"
)

// NewConnectionRole defines a procedure for creating a new setup role from a
// raw string.
func NewConnectionRole(raw string) ConnectionRole {
	switch raw {
	case connectionRoleActiveStr:
		return ConnectionRoleActive
	case connectionRolePassiveStr:
		return ConnectionRolePassive
	case connectionRoleActPassStr:
		return ConnectionRoleActPass
	case connectionRoleHoldConnStr:
		return ConnectionRoleHoldConn
	default:
		return ConnectionRole(unknown)
	}
}

func (t ConnectionRole) String() string {
	switch t {
	case ConnectionRoleActive:
		return connectionRoleActiveStr
	case ConnectionRolePassive:
		return connectionRolePassiveStr
	case ConnectionRoleActPass:
		return connectionRoleActPassStr
	case ConnectionRoleHoldConn:
		return connectionRoleHoldConnStr
	default:
		return unknownStr
	}
}
