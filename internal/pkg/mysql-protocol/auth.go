package mysql_protocol

import (
	"bytes"
)

func ValidatePassword(expected string, handshakePacket *HandshakePacket, loginPacket *LoginPacket) bool {
	var scrambledPassword = scramblePassword(expected, []byte(string(handshakePacket.Salt1)+string(handshakePacket.Salt2)))

	return bytes.Equal(loginPacket.Password, scrambledPassword)
}
