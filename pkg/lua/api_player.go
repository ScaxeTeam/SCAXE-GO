package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func registerPlayerAPI(L *lua.LState, server ServerAPI) {
	mod := L.NewTable()

	mod.RawSetString("getByName", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		p := server.GetPlayer(name)
		if p == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(playerToTable(L, p))
		return 1
	}))

	mod.RawSetString("getAll", L.NewFunction(func(L *lua.LState) int {
		players := server.GetOnlinePlayers()
		tbl := L.NewTable()
		for i, p := range players {
			tbl.RawSetInt(i+1, playerToTable(L, p))
		}
		L.Push(tbl)
		return 1
	}))

	mod.RawSetString("sendMessage", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		msg := L.CheckString(2)
		p := server.GetPlayer(name)
		if p != nil {
			p.SendMessage(msg)
		}
		return 0
	}))

	mod.RawSetString("kick", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		reason := L.OptString(2, "Kicked by plugin")
		server.KickPlayer(name, reason)
		return 0
	}))

	mod.RawSetString("getPosition", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		p := server.GetPlayer(name)
		if p == nil {
			L.Push(lua.LNil)
			return 1
		}
		x, y, z := p.GetPosition()
		tbl := L.NewTable()
		tbl.RawSetString("x", lua.LNumber(x))
		tbl.RawSetString("y", lua.LNumber(y))
		tbl.RawSetString("z", lua.LNumber(z))
		L.Push(tbl)
		return 1
	}))

	mod.RawSetString("teleport", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		x := L.CheckNumber(2)
		y := L.CheckNumber(3)
		z := L.CheckNumber(4)
		p := server.GetPlayer(name)
		if p != nil {
			p.SetPosition(float64(x), float64(y), float64(z))
		}
		return 0
	}))

	mod.RawSetString("setGamemode", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		mode := L.CheckInt(2)
		p := server.GetPlayer(name)
		if p != nil {
			p.SetGamemode(mode)
		}
		return 0
	}))

	mod.RawSetString("getGamemode", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		p := server.GetPlayer(name)
		if p == nil {
			L.Push(lua.LNumber(-1))
			return 1
		}
		L.Push(lua.LNumber(p.GetGamemode()))
		return 1
	}))

	mod.RawSetString("isOp", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		p := server.GetPlayer(name)
		if p == nil {
			L.Push(lua.LFalse)
			return 1
		}
		if p.IsOp() {
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LFalse)
		}
		return 1
	}))

	mod.RawSetString("getHealth", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		p := server.GetPlayer(name)
		if p == nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		L.Push(lua.LNumber(p.GetHealth()))
		return 1
	}))

	mod.RawSetString("setHealth", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		health := L.CheckInt(2)
		p := server.GetPlayer(name)
		if p != nil {
			p.SetHealth(health)
		}
		return 0
	}))

	L.SetGlobal("player", mod)
}

func playerToTable(L *lua.LState, p PlayerAPI) *lua.LTable {
	tbl := L.NewTable()
	tbl.RawSetString("name", lua.LString(p.GetName()))
	x, y, z := p.GetPosition()
	tbl.RawSetString("x", lua.LNumber(x))
	tbl.RawSetString("y", lua.LNumber(y))
	tbl.RawSetString("z", lua.LNumber(z))
	tbl.RawSetString("health", lua.LNumber(p.GetHealth()))
	tbl.RawSetString("gamemode", lua.LNumber(p.GetGamemode()))
	if p.IsOp() {
		tbl.RawSetString("op", lua.LTrue)
	} else {
		tbl.RawSetString("op", lua.LFalse)
	}
	tbl.RawSetString("entityId", lua.LNumber(p.GetEntityID()))
	return tbl
}
