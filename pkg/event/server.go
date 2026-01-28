package event

type ServerEvent struct {
	*BaseEvent
}

func NewServerEvent(name string) *ServerEvent {
	return &ServerEvent{
		BaseEvent: NewBaseEvent(name),
	}
}

type CommandEvent struct {
	*ServerEvent
	Command string
	Sender  string
}

var commandEventHandlers = NewHandlerList()

func NewCommandEvent(command, sender string) *CommandEvent {
	return &CommandEvent{
		ServerEvent: NewServerEvent("CommandEvent"),
		Command:     command,
		Sender:      sender,
	}
}

func (e *CommandEvent) GetHandlers() *HandlerList {
	return commandEventHandlers
}

func (e *CommandEvent) GetCommand() string {
	return e.Command
}

func (e *CommandEvent) SetCommand(cmd string) {
	e.Command = cmd
}

type PluginEnableEvent struct {
	*ServerEvent
	PluginName string
}

var pluginEnableHandlers = NewHandlerList()

func NewPluginEnableEvent(pluginName string) *PluginEnableEvent {
	return &PluginEnableEvent{
		ServerEvent: NewServerEvent("PluginEnableEvent"),
		PluginName:  pluginName,
	}
}

func (e *PluginEnableEvent) GetHandlers() *HandlerList {
	return pluginEnableHandlers
}

type PluginDisableEvent struct {
	*ServerEvent
	PluginName string
}

var pluginDisableHandlers = NewHandlerList()

func NewPluginDisableEvent(pluginName string) *PluginDisableEvent {
	return &PluginDisableEvent{
		ServerEvent: NewServerEvent("PluginDisableEvent"),
		PluginName:  pluginName,
	}
}

func (e *PluginDisableEvent) GetHandlers() *HandlerList {
	return pluginDisableHandlers
}

type DataPacketReceiveEvent struct {
	*ServerEvent
	PlayerID int64
	PacketID int
}

var dataPacketReceiveHandlers = NewHandlerList()

func NewDataPacketReceiveEvent(playerID int64, packetID int) *DataPacketReceiveEvent {
	return &DataPacketReceiveEvent{
		ServerEvent: NewServerEvent("DataPacketReceiveEvent"),
		PlayerID:    playerID,
		PacketID:    packetID,
	}
}

func (e *DataPacketReceiveEvent) GetHandlers() *HandlerList {
	return dataPacketReceiveHandlers
}

type DataPacketSendEvent struct {
	*ServerEvent
	PlayerID int64
	PacketID int
}

var dataPacketSendHandlers = NewHandlerList()

func NewDataPacketSendEvent(playerID int64, packetID int) *DataPacketSendEvent {
	return &DataPacketSendEvent{
		ServerEvent: NewServerEvent("DataPacketSendEvent"),
		PlayerID:    playerID,
		PacketID:    packetID,
	}
}

func (e *DataPacketSendEvent) GetHandlers() *HandlerList {
	return dataPacketSendHandlers
}

type QueryRegenerateEvent struct {
	*ServerEvent
	ServerName    string
	WorldName     string
	OnlinePlayers int
	MaxPlayers    int
	GameMode      string
	Version       string
}

var queryRegenerateHandlers = NewHandlerList()

func NewQueryRegenerateEvent(serverName, worldName string, online, max int) *QueryRegenerateEvent {
	return &QueryRegenerateEvent{
		ServerEvent:   NewServerEvent("QueryRegenerateEvent"),
		ServerName:    serverName,
		WorldName:     worldName,
		OnlinePlayers: online,
		MaxPlayers:    max,
	}
}

func (e *QueryRegenerateEvent) GetHandlers() *HandlerList {
	return queryRegenerateHandlers
}
