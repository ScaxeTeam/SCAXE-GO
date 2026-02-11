package lua

import (
	"github.com/scaxe/scaxe-go/pkg/command"
	lua "github.com/yuin/gopher-lua"
)

type luaCommand struct {
	command.BaseCommand
	callback   *lua.LFunction
	state      *lua.LState
	pluginName string
}

func (c *luaCommand) Execute(sender command.CommandSender, args []string) bool {
	senderTable := c.state.NewTable()
	senderTable.RawSetString("name", lua.LString(sender.GetName()))
	if sender.IsOp() {
		senderTable.RawSetString("op", lua.LTrue)
	} else {
		senderTable.RawSetString("op", lua.LFalse)
	}

	argsTable := c.state.NewTable()
	for i, arg := range args {
		argsTable.RawSetInt(i+1, lua.LString(arg))
	}

	if err := c.state.CallByParam(lua.P{
		Fn:      c.callback,
		NRet:    0,
		Protect: true,
	}, senderTable, argsTable); err != nil {
		sender.SendMessage("§cCommand error: " + err.Error())
		return false
	}
	return true
}

func registerCommandAPI(L *lua.LState, p *Plugin, server ServerAPI) {
	mod := L.NewTable()

	mod.RawSetString("register", L.NewFunction(func(L *lua.LState) int {
		tbl := L.CheckTable(1)

		name := getStringField(L, tbl, "name")
		if name == "" {
			L.ArgError(1, "command name is required")
			return 0
		}

		desc := getStringField(L, tbl, "description")
		usage := getStringField(L, tbl, "usage")
		perm := getStringField(L, tbl, "permission")
		callback := tbl.RawGetString("callback")

		fn, ok := callback.(*lua.LFunction)
		if !ok {
			L.ArgError(1, "callback function is required")
			return 0
		}

		cmd := &luaCommand{
			BaseCommand: command.BaseCommand{
				Name:        name,
				Description: desc,
				Usage:       usage,
				Permission:  perm,
			},
			callback:   fn,
			state:      L,
			pluginName: p.Meta.Name,
		}

		server.RegisterCommand(cmd)
		return 0
	}))

	L.SetGlobal("commands", mod)
}

func getStringField(L *lua.LState, tbl *lua.LTable, key string) string {
	val := tbl.RawGetString(key)
	if str, ok := val.(lua.LString); ok {
		return string(str)
	}
	return ""
}
