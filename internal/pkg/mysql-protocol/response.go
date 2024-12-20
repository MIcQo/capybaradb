package mysql_protocol

import (
	"bytes"
	"encoding/binary"
)

// OK SimpleCustomPacket
func NewOKPacket() []byte {
	payload := bytes.Buffer{}
	payload.WriteByte(0x00)                                     // OK header
	payload.Write([]byte{0x00, 0x00})                           // Affected rows and last insert ID
	binary.Write(&payload, binary.LittleEndian, uint16(0x0002)) // Server status
	binary.Write(&payload, binary.LittleEndian, uint16(0x0000)) // Warnings

	packet := bytes.Buffer{}
	packetLength := len(payload.Bytes())
	packet.Write([]byte{byte(packetLength), byte(packetLength >> 8), byte(packetLength >> 16), 0})
	packet.Write(payload.Bytes())

	return packet.Bytes()
}

// Error SimpleCustomPacket
func NewErrorPacket(errorMessage string) []byte {
	payload := bytes.Buffer{}
	payload.WriteByte(0xFF)                                   // Error header
	binary.Write(&payload, binary.LittleEndian, uint16(1045)) // Error code
	payload.WriteByte('#')                                    // SQL state marker
	payload.WriteString("28000")                              // SQL state
	payload.WriteString(errorMessage)

	packet := bytes.Buffer{}
	packetLength := len(payload.Bytes())
	packet.Write([]byte{byte(packetLength), byte(packetLength >> 8), byte(packetLength >> 16), 0})
	packet.Write(payload.Bytes())

	return packet.Bytes()
}
