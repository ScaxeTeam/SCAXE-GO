package lua

import (
	"github.com/scaxe/scaxe-go/pkg/logger"
	lua "github.com/yuin/gopher-lua"
)

func registerLoggerAPI(L *lua.LState, p *Plugin) {
	mod := L.NewTable()

	mod.RawSetString("info", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		logger.Info("[" + p.Meta.Name + "] " + msg)
		return 0
	}))

	mod.RawSetString("warn", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		logger.Warn("[" + p.Meta.Name + "] " + msg)
		return 0
	}))

	mod.RawSetString("error", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		logger.Error("[" + p.Meta.Name + "] " + msg)
		return 0
	}))

	mod.RawSetString("debug", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		logger.Debug("[" + p.Meta.Name + "] " + msg)
		return 0
	}))

	L.SetGlobal("logger", mod)
}
