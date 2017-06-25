package proxy

import "testing"

func TestReadLenEncodedStringWithValidData(t *testing.T) {
	expectedStr := "ABCDEFGHIKLMONPQRSTYW"
	validStringBytes := []byte{
		0x15, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4b, 0x4c, 0x4d, 0x4f, 0x4e, 0x50,
		0x51, 0x52, 0x53, 0x54, 0x59, 0x57,
	}

	_, str := readLenEncodedString(validStringBytes)

	if str != expectedStr {
		t.Errorf("Expected '%s', got '%s'", expectedStr, str)
	}
}

func TestDecodeComStmtExecuteRequestWithIncorrectPacketType(t *testing.T) {
	invalidTypePacket := []byte{
		0x43, 0x00, 0x00, 0x00, 0x18, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01,
		0xfd, 0x00, 0xfd, 0x00, 0xfd, 0x00, 0x13, 0x31, 0x2e, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x31, 0x30, 0x31, 0x31, 0x31, 0x45, 0x2b, 0x32, 0x31, 0x06, 0x58, 0x59, 0x5a, 0x5a, 0x5a,
		0x5a, 0x15, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4b, 0x4c, 0x4d, 0x4f, 0x4e,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x59, 0x57,
	}

	_, err := DecodeComStmtExecuteRequest(invalidTypePacket, 0)

	if err == nil {
		t.Errorf("Expected '%s', got nil", errInvalidPacketType)
	}
}

func TestDecodeComStmtExecuteRequestWithIncorrectPacketSize(t *testing.T) {
	invalidLengthPacket := []byte{0x43, 0x00, 0x00, 0x00, 0x17}

	_, err := DecodeComStmtExecuteRequest(invalidLengthPacket, 0)

	if err == nil {
		t.Errorf("Expected '%s', got nil", errInvalidPacketLength)
	}
}

func TestDecodeComStmtExecuteRequestCorrectPacketWithStringParams(t *testing.T) {
	validPacketParametersValues := []string{"1.2345678910111E+21", "XYZZZZ", "ABCDEFGHIKLMONPQRSTYW"}
	packetParametersCount := len(validPacketParametersValues)
	validPacket := []byte{
		0x43, 0x00, 0x00, 0x00, 0x17, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01,
		0xfd, 0x00, 0xfd, 0x00, 0xfd, 0x00, 0x13, 0x31, 0x2e, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x31, 0x30, 0x31, 0x31, 0x31, 0x45, 0x2b, 0x32, 0x31, 0x06, 0x58, 0x59, 0x5a, 0x5a, 0x5a,
		0x5a, 0x15, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4b, 0x4c, 0x4d, 0x4f, 0x4e,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x59, 0x57,
	}

	decoded, err := DecodeComStmtExecuteRequest(validPacket, packetParametersCount)

	if err != nil {
		t.Errorf("Expected nil, got: '%s'", err.Error())
	}

	if decoded.StatementID != 1 {
		t.Errorf("Expected %d, got %d", 1, decoded.StatementID)
	}

	if len(decoded.PreparedParameters) != packetParametersCount {
		t.Errorf("Expected: %d, got %d", packetParametersCount, len(decoded.PreparedParameters))
	}

	for i := 0; i < packetParametersCount; i++ {
		if decoded.PreparedParameters[i].Value != validPacketParametersValues[i] {
			t.Errorf("Expected %s, got %s", validPacketParametersValues[i], decoded.PreparedParameters[i].Value)
		}
	}
}

func TestDecodeComStmtExecuteRequestCorrectPacketWithNumericParams(t *testing.T) {
	validPacketParametersValues := []string{"12345", "-12345", "123.450000", "-123.450000"}
	packetParametersCount := len(validPacketParametersValues)
	validPacket := []byte{
		0x34, 0x00, 0x00, 0x00, 0x17, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x00, 0x08, 0x00, 0x05, 0x00, 0x05, 0x00, 0x39, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc7, 0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xcd, 0xcc, 0xcc, 0xcc, 0xcc, 0xdc, 0x5e, 0x40,
		0xcd, 0xcc, 0xcc, 0xcc, 0xcc, 0xdc, 0x5e, 0xc0,
	}

	decoded, err := DecodeComStmtExecuteRequest(validPacket, packetParametersCount)

	if err != nil {
		t.Errorf("Expected nil, got: '%s'", err.Error())
	}

	for i := 0; i < packetParametersCount; i++ {
		if decoded.PreparedParameters[i].Value != validPacketParametersValues[i] {
			t.Errorf("Expected %s, got %s", validPacketParametersValues[i], decoded.PreparedParameters[i].Value)
		}
	}
}

func TestDecodeComStmtPrepareOkResponseWithInvalidPacket(t *testing.T) {
	invalidPacket := []byte{
		0x0c, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x04,
	}

	_, err := DecodeComStmtPrepareOkResponse(invalidPacket)

	if err == nil {
		t.Errorf("Expected '%s', got: nil", err.Error())
	}
}

func TestDecodeComStmtPrepareOkResponseWithValidPacket(t *testing.T) {
	validPacket := []byte{
		0x0c, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
	}

	decoded, err := DecodeComStmtPrepareOkResponse(validPacket)

	if err != nil {
		t.Errorf("Expected nil, got: '%s'", err.Error())
	}

	if decoded.StatementID != 1 {
		t.Errorf("Expected 1, got: %d", decoded.StatementID)
	}

	if decoded.ParametersNum != 4 {
		t.Errorf("Expected 4, got: %d", decoded.ParametersNum)
	}
}
