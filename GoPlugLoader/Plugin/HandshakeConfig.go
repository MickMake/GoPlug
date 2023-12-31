package Plugin

import (
	plugin "github.com/hashicorp/go-plugin"
)

// HandshakeConfig
// used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  HandshakeProtocol,
	MagicCookieKey:   HandshakeKey,
	MagicCookieValue: HandshakeVersion,
}
