package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func registerLevelAPI(L *lua.LState, server ServerAPI) {
	mod := L.NewTable()

	mod.RawSetString("getBlock", L.NewFunction(func(L *lua.LState) int {
		x := int32(L.CheckInt(1))
		y := int32(L.CheckInt(2))
		z := int32(L.CheckInt(3))
		lvl := server.GetLevel()
		if lvl == nil {
			L.Push(lua.LNil)
			return 1
		}
		id, meta := lvl.GetBlock(x, y, z)
		tbl := L.NewTable()
		tbl.RawSetString("id", lua.LNumber(id))
		tbl.RawSetString("meta", lua.LNumber(meta))
		L.Push(tbl)
		return 1
	}))

	mod.RawSetString("setBlock", L.NewFunction(func(L *lua.LState) int {
		x := int32(L.CheckInt(1))
		y := int32(L.CheckInt(2))
		z := int32(L.CheckInt(3))
		id := uint8(L.CheckInt(4))
		meta := uint8(L.OptInt(5, 0))
		lvl := server.GetLevel()
		if lvl != nil {
			lvl.SetBlock(x, y, z, id, meta)
		}
		return 0
	}))

	mod.RawSetString("getTime", L.NewFunction(func(L *lua.LState) int {
		lvl := server.GetLevel()
		if lvl == nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		L.Push(lua.LNumber(lvl.GetTime()))
		return 1
	}))

	mod.RawSetString("setTime", L.NewFunction(func(L *lua.LState) int {
		t := int64(L.CheckNumber(1))
		lvl := server.GetLevel()
		if lvl != nil {
			lvl.SetTime(t)
		}
		return 0
	}))

	mod.RawSetString("getSeed", L.NewFunction(func(L *lua.LState) int {
		lvl := server.GetLevel()
		if lvl == nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		L.Push(lua.LNumber(lvl.GetSeed()))
		return 1
	}))

	mod.RawSetString("getSpawnPosition", L.NewFunction(func(L *lua.LState) int {
		lvl := server.GetLevel()
		if lvl == nil {
			L.Push(lua.LNil)
			return 1
		}
		x, y, z := lvl.GetSpawnLocation()
		tbl := L.NewTable()
		tbl.RawSetString("x", lua.LNumber(x))
		tbl.RawSetString("y", lua.LNumber(y))
		tbl.RawSetString("z", lua.LNumber(z))
		L.Push(tbl)
		return 1
	}))

	L.SetGlobal("level", mod)
}
