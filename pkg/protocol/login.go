package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type LoginPacket struct {
	BasePacket
	Protocol      int32
	Edition       byte
	ChainData     string
	ClientData    string
	ClientUUID    string
	ClientID      int64
	Username      string
	SkinID        string
	SkinData      []byte
	ServerAddress string
	LanguageCode  string
	DeviceOS      int32
	DeviceModel   string
}

func NewLoginPacket() *LoginPacket {
	return &LoginPacket{
		BasePacket: BasePacket{PacketID: IDLogin},
	}
}

func (p *LoginPacket) Name() string {
	return "LoginPacket"
}

func (p *LoginPacket) Encode(stream *BinaryStream) error {
	logger.Debug("LoginPacket.Encode", "username", p.Username, "protocol", p.Protocol)
	EncodeHeader(stream, p.ID())
	stream.WriteString16(p.Username)
	stream.WriteInt(p.Protocol)
	stream.WriteInt(p.Protocol)
	stream.WriteLong(p.ClientID)
	stream.WriteBytes([]byte(p.ClientUUID))

	stream.WriteString16(p.ServerAddress)

	stream.WriteString16("")

	stream.WriteString16(p.SkinID)

	stream.WriteString16(string(p.SkinData))

	return nil
}

func (p *LoginPacket) Decode(stream *BinaryStream) error {
	logger.Debug("LoginPacket.Decode", "start", true, "len", stream.Len())

	var err error

	p.Username, err = stream.ReadString16()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read username", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "username", p.Username)

	p.Protocol, err = stream.ReadInt()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read protocol1", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "protocol1", p.Protocol)

	protocol2, err := stream.ReadInt()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read protocol2", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "protocol2", protocol2)

	p.ClientID, err = stream.ReadLong()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read clientId", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "clientId", p.ClientID)

	uuidBytes, err := stream.ReadBytes(16)
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read UUID", "err", err)
		return err
	}
	p.ClientUUID = formatUUID(uuidBytes)
	logger.Debug("LoginPacket.Decode", "uuid", p.ClientUUID)

	p.ServerAddress, err = stream.ReadString16()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read serverAddress", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "serverAddress", p.ServerAddress)

	clientSecret, err := stream.ReadString16()
	if err != nil {

		logger.Error("LoginPacket.Decode", "error", "failed to read clientSecret", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "clientSecret", clientSecret)

	p.SkinID, err = stream.ReadString16()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read skinName", "err", err)
		return err
	}
	logger.Debug("LoginPacket.Decode", "skinName", p.SkinID)

	skinData, err := stream.ReadString16()
	if err != nil {
		logger.Error("LoginPacket.Decode", "error", "failed to read skinData", "err", err)

		return err
	}
	p.SkinData = []byte(skinData)
	logger.Debug("LoginPacket.Decode", "skinDataLen", len(p.SkinData))

	logger.Info("LoginPacket.Decode", "success", true, "username", p.Username, "protocol", p.Protocol, "uuid", p.ClientUUID)
	return nil
}

func formatUUID(data []byte) string {
	if len(data) < 16 {
		return ""
	}

	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		data[0], data[1], data[2], data[3],
		data[4], data[5],
		data[6], data[7],
		data[8], data[9],
		data[10], data[11], data[12], data[13], data[14], data[15])
}

func (p *LoginPacket) parseChainData() {
	logger.Debug("LoginPacket.parseChainData", "chainLength", len(p.ChainData))

	var chainWrapper struct {
		Chain []string `json:"chain"`
	}

	if err := json.Unmarshal([]byte(p.ChainData), &chainWrapper); err != nil {
		logger.Error("LoginPacket.parseChainData", "error", "failed to parse chain wrapper", "err", err)
		return
	}

	for _, token := range chainWrapper.Chain {

		payload := extractJWTPayload(token)
		if payload == "" {
			continue
		}

		var claims map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &claims); err != nil {
			continue
		}

		if extraData, ok := claims["extraData"].(map[string]interface{}); ok {
			if displayName, ok := extraData["displayName"].(string); ok {
				p.Username = displayName
				logger.Debug("LoginPacket.parseChainData", "foundUsername", displayName)
			}
			if identity, ok := extraData["identity"].(string); ok {
				p.ClientUUID = identity
				logger.Debug("LoginPacket.parseChainData", "foundUUID", identity)
			}
		}
	}

	p.parseClientData()
}

func (p *LoginPacket) parseClientData() {
	payload := extractJWTPayload(p.ClientData)
	if payload == "" {
		return
	}

	var claims map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &claims); err != nil {
		logger.Error("LoginPacket.parseClientData", "error", err)
		return
	}

	if skinID, ok := claims["SkinId"].(string); ok {
		p.SkinID = skinID
		logger.Debug("LoginPacket.parseClientData", "skinID", skinID)
	}

	if serverAddr, ok := claims["ServerAddress"].(string); ok {
		p.ServerAddress = serverAddr
		logger.Debug("LoginPacket.parseClientData", "serverAddress", serverAddr)
	}

	if langCode, ok := claims["LanguageCode"].(string); ok {
		p.LanguageCode = langCode
		logger.Debug("LoginPacket.parseClientData", "languageCode", langCode)
	}

	if deviceOS, ok := claims["DeviceOS"].(float64); ok {
		p.DeviceOS = int32(deviceOS)
		logger.Debug("LoginPacket.parseClientData", "deviceOS", p.DeviceOS)
	}

	if deviceModel, ok := claims["DeviceModel"].(string); ok {
		p.DeviceModel = deviceModel
		logger.Debug("LoginPacket.parseClientData", "deviceModel", deviceModel)
	}
}

func extractJWTPayload(token string) string {

	parts := splitJWT(token)
	if len(parts) < 2 {
		return ""
	}

	decoded, err := base64URLDecode(parts[1])
	if err != nil {
		logger.Error("extractJWTPayload", "error", err)
		return ""
	}

	return string(decoded)
}

func splitJWT(token string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	parts = append(parts, token[start:])
	return parts
}

func base64URLDecode(s string) ([]byte, error) {

	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '-':
			result = append(result, '+')
		case '_':
			result = append(result, '/')
		default:
			result = append(result, s[i])
		}
	}

	return base64Decode(string(result))
}

func base64Decode(s string) ([]byte, error) {
	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	charToVal := make(map[byte]int)
	for i := 0; i < len(base64Chars); i++ {
		charToVal[base64Chars[i]] = i
	}

	var result []byte
	var accumulator int
	var bits int

	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			break
		}
		val, ok := charToVal[s[i]]
		if !ok {
			continue
		}
		accumulator = (accumulator << 6) | val
		bits += 6

		if bits >= 8 {
			bits -= 8
			result = append(result, byte(accumulator>>bits))
			accumulator &= (1 << bits) - 1
		}
	}

	return result, nil
}

func init() {
	RegisterPacket(IDLogin, func() DataPacket { return NewLoginPacket() })
	RegisterPacket(IDContainerSetContent, func() DataPacket { return &ContainerSetContentPacket{} })
	RegisterPacket(IDContainerSetSlot, func() DataPacket { return &ContainerSetSlotPacket{} })
	RegisterPacket(IDUseItem, func() DataPacket { return &UseItemPacket{} })
	RegisterPacket(IDUpdateBlock, func() DataPacket { return &UpdateBlockPacket{} })
}
