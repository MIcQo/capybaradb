package mysql_protocol

// How MySQL protocol works:
// Connection Phase:
//  | client 					-> 	server |
//  | server greeting 			-> 	client |
//  | client auth 				-> 	server |
//  | server return OK or ERR 	-> 	client |
// when previous step was successful we can step up to query phase
// | client can send query 		-> 	server |
// | server respond OK or ERR 	-> 	client |
//
// based on these steps we can say we have 3 types of responses:
// HandshakePacket
// OK Response
// ERR Response
// To identify incoming messages we use this:
// Auth
// Query

// Predefined credentials
var (
	allowedUsername = "root"
	allowedPassword = "aa" // Simulate password validation (hashed password comparison is omitted here)
)

type Charset byte

const (
	ServerStatusInTrans            = 1 << iota // A transaction is currently active
	ServerStatusAutocommit                     // Autocommit mode is set
	_                                          // ignored
	ServerMoreResultsExists                    // More results exists (more packets will follow)
	ServerQueryNoGoodIndexUsed                 // Set if EXPLAIN would've shown Range checked for each record
	ServerQueryNoIndexUsed                     // The query did not use an index
	ServerStatusCursorExists                   // When using COM_STMT_FETCH, indicate that current cursor still has result
	ServerStatusLastRowSent                    // When using COM_STMT_FETCH, indicate that current cursor has finished to send results
	ServerStatusDbDropped                      // Database has been dropped
	ServerStatusNoBackslashEscapes             // Current escape mode is "no backslash escape"
	ServerStatusMetadataChanged                // A DDL change did have an impact on an existing PREPARE (an automatic reprepare has been executed)
	ServerQueryWasSlow                         // The query was slower than long_query_time
	ServerPsOutParams                          // This resultset contain stored procedure output parameter
	ServerStatusInTransReadonly                // Current transaction is a read-only transaction
	ServerSessionStateChanged                  // Session state change. See Session change type for more information
)

const (
	Big5ChineseCI   Charset = 1
	Latin2CzechCS           = 2
	Dec8SwedishCI           = 3
	CP850GeneralCI          = 4
	Latin1German1CI         = 5
	HP8EnglishCI            = 6
	Koi8rGeneralCI          = 7
	Latin1SwedishCI         = 8
	Latin2GeneralCI         = 9
	Swe7SwedishCI           = 10
	Utf8GeneralCI           = 33
	Binary                  = 63
)

const (
	ClientMysql                     uint64 = 1 // Set by older MariaDB versions. MariaDB 10.2 leaves this bit unset to permit MariaDB identification and indicate support for extended capabilities. (MySQL named this CLIENT_LONG_PASSWORD)
	FoundRows                       uint64 = 2
	ConnectWithDb                   uint64 = 8      // One can specify db on connect
	Compress                        uint64 = 32     // Can use compression protocol
	LocalFiles                      uint64 = 128    // Can use LOAD DATA LOCAL
	IgnoreSpace                     uint64 = 256    // Ignore spaces before '('
	ClientProtocol41                uint64 = 1 << 9 //4.1 protocol
	ClientInteractive               uint64 = 1 << 10
	Ssl                             uint64 = 1 << 11 //Can use SSL
	Transactions                    uint64 = 1 << 13
	SecureConnection                uint64 = 1 << 15 //4.1 authentication
	MultiStatements                 uint64 = 1 << 16 //Enable/disable multi-stmt support
	MultiResults                    uint64 = 1 << 17 //Enable/disable multi-results
	PsMultiResults                  uint64 = 1 << 18 //Enable/disable multi-results for PrepareStatement
	PluginAuth                      uint64 = 1 << 19 //Client supports plugin authentication
	ConnectAttrs                    uint64 = 1 << 20 //Client send connection attributes
	PluginAuthLenencClientData      uint64 = 1 << 21 //Enable authentication response packet to be larger than 255 bytes
	ClientCanHandleExpiredPasswords uint64 = 1 << 22 //Client can handle expired passwords
	ClientSessionTrack              uint64 = 1 << 23 //Enable/disable session tracking in OK_Packet
	ClientDeprecateEof              uint64 = 1 << 24 // EOF_Packet deprecation : // * OK_Packet replace EOF_Packet in end of Resulset when in text format //* EOF_Packet between columns definition and resultsetRows is deleted
	ClientOptionalResultsetMetadata uint64 = 1 << 25 //Not use for MariaDB
	ClientZstdCompressionAlgorithm  uint64 = 1 << 26 //Support zstd protocol compression
	ClientCapabilityExtension       uint64 = 1 << 29 //Reserved for future use. (Was CLIENT_PROGRESS Client support progress indicator before 10.2)
	ClientSslVerifyServerCert       uint64 = 1 << 30 //Client verify server certificate. deprecated, client have options to indicate if server certifiate must be verified
	ClientRememberOptions           uint64 = 1 << 31
	MariadbClientProgress           uint64 = 1 << 32 // Client support progress indicator (since 10.2)
	MariadbClientComMulti           uint64 = 1 << 33 // Permit COM_MULTI protocol
	MariadbClientStmtBulkOperations uint64 = 1 << 34 //Permit bulk insert
	MariadbClientExtendedMetadata   uint64 = 1 << 35 // Add extended metadata information
	MariadbClientCacheMetadata      uint64 = 1 << 36 // Permit skipping metadata
	MariadbClientBulkUnitResults    uint64 = 1 << 37 // when enable, indicate that Bulk command can use STMT_BULK_FLAG_SEND_UNIT_RESULTS flag that permit to return a result-set of all affected rows and auto-increment values
)

var charsetToName = map[Charset]string{
	Big5ChineseCI:   "big5_chinese_ci",
	Latin2CzechCS:   "latin2_czech_cs",
	Dec8SwedishCI:   "dec8_swedish_ci",
	CP850GeneralCI:  "cp850_general_ci",
	Latin1German1CI: "latin1_german1_ci",
	HP8EnglishCI:    "hp8_english_ci",
	Koi8rGeneralCI:  "koi8r_general_ci",
	Latin1SwedishCI: "latin1_swedish_ci",
	Latin2GeneralCI: "latin2_general_ci",
	Swe7SwedishCI:   "swe7_swedish_ci",
	Utf8GeneralCI:   "utf8_general_ci",
	Binary:          "binary",
}

var capabilityValueToString = map[uint64]string{
	ClientMysql:                     "CLIENT_MYSQL",
	FoundRows:                       "FOUND_ROWS",
	ConnectWithDb:                   "CONNECT_WITH_DB",
	Compress:                        "COMPRESS",
	LocalFiles:                      "LOCAL_FILES",
	IgnoreSpace:                     "IGNORE_SPACE",
	ClientProtocol41:                "CLIENT_PROTOCOL41",
	ClientInteractive:               "CLIENT_INTERACTIVE",
	Ssl:                             "SSL",
	Transactions:                    "TRANSACTIONS",
	SecureConnection:                "SECURE_CONNECTION",
	MultiStatements:                 "MULTI_STATEMENTS",
	MultiResults:                    "MULTI_RESULTS",
	PsMultiResults:                  "PS_MULTI_RESULTS",
	PluginAuth:                      "PLUGIN_AUTH",
	ConnectAttrs:                    "CONNECT_ATTRS",
	PluginAuthLenencClientData:      "PLUGIN_AUTH_LENENC_CLIENT_DATA",
	ClientCanHandleExpiredPasswords: "CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS",
	ClientSessionTrack:              "CLIENT_SESSION_TRACK",
	ClientDeprecateEof:              "CLIENT_DEPRECATE_EOF",
	ClientOptionalResultsetMetadata: "CLIENT_OPTIONAL_RESULTSET_METADATA",
	ClientZstdCompressionAlgorithm:  "CLIENT_ZSTD_COMPRESSION_ALGORITHM",
	ClientCapabilityExtension:       "CLIENT_CAPABILITY_EXTENSION",
	ClientSslVerifyServerCert:       "CLIENT_SSL_VERIFY_SERVER_CERT",
	ClientRememberOptions:           "CLIENT_REMEMBER_OPTIONS",
	MariadbClientProgress:           "MARIADB_CLIENT_PROGRESS",
	MariadbClientComMulti:           "MARIADB_CLIENT_COM_MULTI",
	MariadbClientStmtBulkOperations: "MARIADB_CLIENT_STMT_BULK_OPERATIONS",
	MariadbClientExtendedMetadata:   "MARIADB_CLIENT_EXTENDED_METADATA",
	MariadbClientCacheMetadata:      "MARIADB_CLIENT_CACHE_METADATA",
	MariadbClientBulkUnitResults:    "MARIADB_CLIENT_BULK_UNIT_RESULTS",
}

var (
	MySQLNativePassword = []byte("mysql_native_password")
)

type PacketHeader struct {
	PacketLength   uint32
	PacketSequence uint8
}

type Packet interface {
	Decode(data []byte) (Packet, error)
	Encode() []byte
}

// 0xff = 255 = uint8
// 0xff 0xff = 65355 = uint16
// 0xff 0xff 0xff 0xff = xxx = uint32

// Handshake SimpleCustomPacket
//func createHandshakePacket() []byte {
//	protocolVersion := byte(10)
//	serverVersion := []byte("8.0.0-mock-server")
//	connectionID := uint32(12345)
//	salt1 := []byte("random_salt1")
//	capabilities := uint32(0xA6FF) // CLIENT_PLUGIN_AUTH, CLIENT_LONG_PASSWORD, etc.
//	charset := byte(33)            // utf8_general_ci
//	status := uint16(0x0002)       // SERVER_STATUS_AUTOCOMMIT
//	salt2 := []byte("random_salt2")
//	authPlugin := []byte("mysql_native_password")
//
//	payload := bytes.Buffer{}
//	payload.WriteByte(protocolVersion)
//	payload.Write(serverVersion)
//	payload.WriteByte(0) // Null-terminator for server version
//	binary.Write(&payload, binary.LittleEndian, connectionID)
//	payload.Write(salt1)
//	payload.WriteByte(0) // Null terminator for salt1
//	binary.Write(&payload, binary.LittleEndian, uint16(capabilities))
//	binary.Write(&payload, binary.LittleEndian, charset)
//	binary.Write(&payload, binary.LittleEndian, status)
//	payload.Write(make([]byte, 10))                               // Reserved (all 0s)
//	binary.Write(&payload, binary.LittleEndian, capabilities>>16) // Extended capabilities
//	payload.WriteByte(byte(len(authPlugin)))
//	payload.WriteByte(0) // Reserved
//	payload.Write(salt2)
//	payload.WriteByte(0) // Null terminator for salt2
//	payload.Write(authPlugin)
//	payload.WriteByte(0) // Null terminator for auth plugin
//
//	packet := bytes.Buffer{}
//	packetLength := len(payload.Bytes())
//	packet.Write([]byte{byte(packetLength), byte(packetLength >> 8), byte(packetLength >> 16), 0})
//	packet.Write(payload.Bytes())
//
//	return packet.Bytes()
//}

//
//func parseAuthPacket(data []byte) (username string, password string, err error) {
//	payload := data[4:] // Skip the 4-byte header
//	nullTerminatorIndex := bytes.IndexByte(payload[32:], 0)
//	if nullTerminatorIndex == -1 {
//		return "", "", fmt.Errorf("invalid auth packet format")
//	}
//
//	username = string(payload[32 : 32+nullTerminatorIndex])
//	password = string(payload[32+nullTerminatorIndex+21:]) // Skip the rest for simplicity
//	return username, password, nil
//}
