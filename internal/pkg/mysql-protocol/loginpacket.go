package mysql_protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
)

type LoginPacket struct {
	PacketHeader

	ClientCapabalities         uint16
	ExtendedClientCapabalities uint16
	MaxPacketSize              uint32
	Charset                    Charset
	reserved                   [23]byte
	Username                   string
	Password                   []byte
	AuthPluginName             string
}

func (l *LoginPacket) Decode(data []byte) (Packet, error) {
	var reader = bytes.NewReader(data)

	var packetLength = make([]byte, 3)
	reader.Read(packetLength)
	l.PacketLength = parse3BytesToUint32(packetLength)

	must(binary.Read(reader, binary.LittleEndian, &l.PacketHeader.PacketSequence))

	logrus.WithFields(map[string]any{
		"length":   l.PacketLength,
		"sequence": l.PacketSequence,
	}).Debug("header parsed")

	// now the body
	must(binary.Read(reader, binary.LittleEndian, &l.ClientCapabalities))
	must(binary.Read(reader, binary.LittleEndian, &l.ExtendedClientCapabalities))
	must(binary.Read(reader, binary.LittleEndian, &l.MaxPacketSize))
	must(binary.Read(reader, binary.LittleEndian, &l.Charset))
	must(binary.Read(reader, binary.LittleEndian, &l.reserved))

	var username, _ = readNullTerminatedString(reader)
	l.Username = username

	reader.ReadByte() // null pointer

	var password, _ = readFixedLengthString(reader, 20)
	l.Password = []byte(password)
	PrintBytesBinary(l.Password)

	var authPluginName, _ = readNullTerminatedString(reader)
	l.AuthPluginName = authPluginName

	logrus.WithFields(map[string]any{
		"clientCapabalities":         l.ClientCapabalities,
		"extendedClientCapabalities": l.ExtendedClientCapabalities,
		"maxPacketSize":              l.MaxPacketSize,
		"charset":                    charsetToName[l.Charset],
		"username":                   username,
		"password":                   password,
		"authPluginName":             authPluginName,
	}).Debug("packet parsed")

	return l, nil
}

func (l *LoginPacket) Encode() []byte {
	//TODO implement me
	panic("implement me")
}

func NewLoginPacket() *LoginPacket {
	return &LoginPacket{}
}
