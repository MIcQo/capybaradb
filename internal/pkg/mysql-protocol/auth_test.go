package mysql_protocol

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	var handshakePacket = NewDefaultHandshakePacket()
	var loginPacket = NewLoginPacket()

	loginPacket.Password = []byte{0xd8, 0xd7, 0x3e, 0x65, 0x88, 0x80, 0xf7, 0x31, 0x5f, 0xa0, 0x8a, 0x5a, 0xa2, 0xa1, 0x6e, 0x1, 0x93, 0x3c, 0x50, 0x3c}

	type args struct {
		expected        []byte
		handshakePacket *HandshakePacket
		loginPacket     *LoginPacket
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "expect valid password", args: args{expected: []byte("aa"), handshakePacket: handshakePacket, loginPacket: loginPacket}, want: true},
		{name: "expect invalid password", args: args{expected: []byte{}, handshakePacket: handshakePacket, loginPacket: loginPacket}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ValidatePassword(string(tt.args.expected), tt.args.handshakePacket, tt.args.loginPacket), "ValidatePassword(%v, %v, %v)", tt.args.expected, tt.args.handshakePacket, tt.args.loginPacket)
		})
	}
}
