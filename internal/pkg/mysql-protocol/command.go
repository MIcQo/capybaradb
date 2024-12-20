package mysql_protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	ComSleep           = 0x0
	ComQuit            = 0x01
	ComInitDB          = 0x02
	ComQuery           = 0x03
	ComStats           = 0x09
	ComShutdown        = 0xa
	ComProcessKill     = 0xc
	ComDebug           = 0x0d
	ComPing            = 0x0e
	ComChangeUser      = 0x11
	ComSetOption       = 0x1b
	ComResetConnection = 0x1f
)

type CommandQuery struct {
	Query string
}

// Decode Returns packet
func (c CommandQuery) Decode(data []byte) (Packet, error) {
	var reader = bytes.NewReader(data)
	var query, err = readNullTerminatedString(reader)
	if err != nil {
		return nil, err
	}

	c.Query = query

	return c, nil
}

// Encode returns bytes from packet
func (c CommandQuery) Encode() []byte {
	return []byte(c.Query)
}

// NewCommandQuery returns command query
func NewCommandQuery() *CommandQuery {
	return &CommandQuery{}
}

type CommandQuit struct {
}

func (c CommandQuit) Decode(data []byte) (Packet, error) {
	return c, nil
}

func (c CommandQuit) Encode() []byte {
	return []byte{}
}

func NewCommandQuit() *CommandQuit {
	return &CommandQuit{}
}

func ParseCommandPacket(data []byte) (Packet, error) {
	var reader = bytes.NewReader(data)

	var packetLength = make([]byte, 3)
	reader.Read(packetLength)

	var sequence uint8
	must(binary.Read(reader, binary.LittleEndian, &sequence))

	var command uint8
	must(binary.Read(reader, binary.LittleEndian, &command))

	var restOfData, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	switch command {
	case ComQuery:
		return NewCommandQuery().Decode(restOfData)
	case ComQuit:
		return NewCommandQuit().Decode(restOfData)
	default:
		return nil, errors.New("unknown command")
	}
}
