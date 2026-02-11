package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func registerServerAPI(L *lua.LState, server ServerAPI) {
	mod := L.NewTable()

	mod.RawSetString("broadcast", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		server.BroadcastMessage(msg)
		return 0
	}))

	mod.RawSetString("getOnlineCount", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(server.GetOnlineCount()))
		return 1
	}))

	mod.RawSetString("getMaxPlayers", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(server.GetMaxPlayers()))
		return 1
	}))

	mod.RawSetString("getTPS", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(server.GetTPS()))
		return 1
	}))

	mod.RawSetString("getServerName", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(server.GetServerName()))
		return 1
	}))

	mod.RawSetString("stop", L.NewFunction(func(L *lua.LState) int {
		server.Stop()
		return 0
	}))

	mod.RawSetString("getCurrentTick", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(server.GetCurrentTick()))
		return 1
	}))

	L.SetGlobal("server", mod)
}
