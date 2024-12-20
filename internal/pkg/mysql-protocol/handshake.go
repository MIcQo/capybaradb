package mysql_protocol

import (
	"bytes"
	"capybaradb/internal/pkg/version"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	firstSaltLength  = 8
	secondSaltLength = 12
)

type HandshakePacket struct {
	PacketHeader
	ProtocolVersion             uint8
	ServerVersion               string
	ConnectionID                uint32
	Salt1                       []byte
	ServerCapabilities          uint16
	Charset                     Charset
	Status                      uint16
	ExtendedServerCapabilities  uint16
	AuthPluginDataLength        uint8
	MariaDBExtendedCapabilities uint32
	Salt2                       []byte
	AuthPluginName              string
}

func (g *HandshakePacket) Decode(data []byte) (Packet, error) {
	var reader = bytes.NewReader(data)

	var packetLength = make([]byte, 3)
	reader.Read(packetLength)
	g.PacketLength = parse3BytesToUint32(packetLength)

	must(binary.Read(reader, binary.LittleEndian, &g.PacketHeader.PacketSequence))

	logrus.WithField("length", g.PacketLength).
		WithField("sequence", g.PacketSequence).
		Trace("header parsed")

	must(binary.Read(reader, binary.LittleEndian, &g.ProtocolVersion))

	var serverVersion, _ = readNullTerminatedString(reader)
	g.ServerVersion = serverVersion

	must(binary.Read(reader, binary.LittleEndian, &g.ConnectionID))

	var salt1 = make([]byte, 8)
	reader.Read(salt1)
	g.Salt1 = salt1

	reader.ReadByte() // ignore filler

	must(binary.Read(reader, binary.LittleEndian, &g.ServerCapabilities))

	must(binary.Read(reader, binary.LittleEndian, &g.Charset))
	must(binary.Read(reader, binary.LittleEndian, &g.Status))
	must(binary.Read(reader, binary.LittleEndian, &g.ExtendedServerCapabilities))
	must(binary.Read(reader, binary.LittleEndian, &g.AuthPluginDataLength))

	var reserved = make([]byte, 6)
	reader.Read(reserved)

	must(binary.Read(reader, binary.LittleEndian, &g.MariaDBExtendedCapabilities))

	if uint64(g.ServerCapabilities)&SecureConnection == SecureConnection {
		var length = max(12, g.AuthPluginDataLength-9)
		g.Salt2 = make([]byte, length)
		reader.Read(g.Salt2)
		reader.ReadByte()
	}

	var authPluginName, _ = readNullTerminatedString(reader)
	g.AuthPluginName = authPluginName

	logrus.WithFields(map[string]any{
		"protocol":                    g.ProtocolVersion,
		"serverVersion":               g.ServerVersion,
		"connectionID":                g.ConnectionID,
		"salt1":                       string(g.Salt1),
		"serverCapabilities":          g.ServerCapabilities,
		"charset":                     charsetToName[g.Charset],
		"status":                      g.Status,
		"extendedServerCapabilities":  g.ExtendedServerCapabilities,
		"authPluginDataLength":        g.AuthPluginDataLength,
		"mariaDBExtendedCapabilities": g.MariaDBExtendedCapabilities,
		"salt2":                       string(g.Salt2),
		"authPluginName":              g.AuthPluginName,
	}).Trace("packet parsed")

	return g, nil
}

func (g *HandshakePacket) Encode() []byte {
	payload := bytes.Buffer{}

	// write protocol version
	payload.Write([]byte{g.ProtocolVersion})

	// write server version
	binary.Write(&payload, binary.LittleEndian, []byte(g.ServerVersion))
	payload.WriteByte(0)

	// write connection ID
	payload.Write([]byte{
		byte(g.ConnectionID),
		byte(g.ConnectionID >> 8),
		byte(g.ConnectionID >> 16),
		0,
	})

	// write first salt
	binary.Write(&payload, binary.LittleEndian, g.Salt1)

	// write filler
	payload.WriteByte(0)

	// write first part of capabilities
	payload.Write([]byte{
		byte(g.ServerCapabilities),
		byte(g.ServerCapabilities >> 8),
	})

	// write charset
	payload.WriteByte(byte(g.Charset))

	// write status
	payload.Write([]byte{
		byte(g.Status),
		byte(g.Status >> 8),
	})

	// write second part of capabilities
	payload.Write([]byte{
		byte(g.ExtendedServerCapabilities),
		byte(g.ExtendedServerCapabilities >> 8),
	})

	// write auth plugin data length
	payload.WriteByte(g.AuthPluginDataLength)

	// write reserved bytes
	var reserved = make([]byte, 6)
	payload.Write(reserved)

	// write another filler (third capability list)
	payload.Write([]byte{
		byte(g.MariaDBExtendedCapabilities),
		byte(g.MariaDBExtendedCapabilities >> 8),
		byte(g.MariaDBExtendedCapabilities >> 16),
		0,
	})

	// write auth data 2nd (salt2)
	binary.Write(&payload, binary.LittleEndian, g.Salt2)
	payload.WriteByte(0)

	// write auth plugin name
	binary.Write(&payload, binary.LittleEndian, []byte(g.AuthPluginName))

	payload.WriteByte(0)

	var payloadBytes = payload.Bytes()
	var packet = bytes.Buffer{}
	var packetLength = len(payloadBytes)
	packet.Write([]byte{
		byte(packetLength),
		byte(packetLength >> 8),
		byte(packetLength >> 16),
		g.PacketSequence,
	})
	packet.Write(payloadBytes)

	return packet.Bytes()
}

func newHandshakePacket() *HandshakePacket {
	return new(HandshakePacket)
}

func NewDefaultHandshakePacket() *HandshakePacket {
	var p = newHandshakePacket()
	p.ProtocolVersion = 10
	p.ServerVersion = version.Version
	p.ConnectionID = 67786
	p.Salt1 = []byte(RandStringRunes(firstSaltLength))
	p.ServerCapabilities = uint16(Ssl | Compress | Transactions | ClientProtocol41 | FoundRows | ConnectWithDb | LocalFiles | ClientMysql | ClientInteractive | SecureConnection | IgnoreSpace)
	p.Charset = Utf8GeneralCI
	p.Status = ServerStatusAutocommit
	p.ExtendedServerCapabilities = 0x81ff
	p.AuthPluginDataLength = 0x15
	p.Salt2 = []byte(RandStringRunes(secondSaltLength))
	p.MariaDBExtendedCapabilities = 0x1d
	p.AuthPluginName = string(MySQLNativePassword)

	fmt.Println(string(p.Salt1), string(p.Salt2))
	return p
}
