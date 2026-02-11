package lua

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	lua "github.com/yuin/gopher-lua"
)

type PluginMeta struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
	Main        string `yaml:"main"`
}

type Plugin struct {
	Meta    PluginMeta
	Dir     string
	State   *lua.LState
	Enabled bool

	eventHandlers  map[string][]*lua.LFunction
	schedulerTasks []*schedulerTask
	nextTaskID     int
}

type schedulerTask struct {
	id       int
	callback *lua.LFunction
	interval int64
	delay    int64
	nextRun  int64
	repeat   bool
	cancel   bool
}

func loadPluginMeta(dir string) (*PluginMeta, error) {
	metaPath := filepath.Join(dir, "plugin.yml")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin.yml: %w", err)
	}

	var meta PluginMeta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse plugin.yml: %w", err)
	}

	if meta.Name == "" {
		return nil, fmt.Errorf("plugin name is required in plugin.yml")
	}
	if meta.Main == "" {
		meta.Main = "main.lua"
	}

	return &meta, nil
}

func newPlugin(meta *PluginMeta, dir string) *Plugin {
	return &Plugin{
		Meta:           *meta,
		Dir:            dir,
		Enabled:        false,
		eventHandlers:  make(map[string][]*lua.LFunction),
		schedulerTasks: make([]*schedulerTask, 0),
	}
}

func (p *Plugin) createState(server ServerAPI) *lua.LState {
	L := lua.NewState(lua.Options{
		SkipOpenLibs: false,
	})

	registerServerAPI(L, server)
	registerPlayerAPI(L, server)
	registerLevelAPI(L, server)
	registerEventAPI(L, p)
	registerCommandAPI(L, p, server)
	registerSchedulerAPI(L, p, server)
	registerLoggerAPI(L, p)

	p.State = L
	return L
}

func (p *Plugin) enable() error {
	if p.State == nil {
		return fmt.Errorf("lua state not initialized for plugin %s", p.Meta.Name)
	}

	mainFile := filepath.Join(p.Dir, p.Meta.Main)
	if err := p.State.DoFile(mainFile); err != nil {
		return fmt.Errorf("failed to execute %s: %w", p.Meta.Main, err)
	}

	onEnable := p.State.GetGlobal("onEnable")
	if fn, ok := onEnable.(*lua.LFunction); ok {
		if err := p.State.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}); err != nil {
			return fmt.Errorf("onEnable error: %w", err)
		}
	}

	p.Enabled = true
	return nil
}

func (p *Plugin) disable() {
	if p.State != nil {
		onDisable := p.State.GetGlobal("onDisable")
		if fn, ok := onDisable.(*lua.LFunction); ok {
			_ = p.State.CallByParam(lua.P{
				Fn:      fn,
				NRet:    0,
				Protect: true,
			})
		}
		p.State.Close()
		p.State = nil
	}

	p.Enabled = false
	p.eventHandlers = make(map[string][]*lua.LFunction)
	p.schedulerTasks = nil
}

func (p *Plugin) tick(currentTick int64) {
	if !p.Enabled || p.State == nil {
		return
	}

	for _, task := range p.schedulerTasks {
		if task.cancel {
			continue
		}
		if currentTick >= task.nextRun {
			if err := p.State.CallByParam(lua.P{
				Fn:      task.callback,
				NRet:    0,
				Protect: true,
			}); err != nil {
				fmt.Printf("[Plugin:%s] scheduler error: %v\n", p.Meta.Name, err)
			}
			if task.repeat {
				task.nextRun = currentTick + task.interval
			} else {
				task.cancel = true
			}
		}
	}
}

func (p *Plugin) callEvent(eventName string, eventTable *lua.LTable) bool {
	if !p.Enabled || p.State == nil {
		return false
	}

	handlers, ok := p.eventHandlers[eventName]
	if !ok || len(handlers) == 0 {
		return false
	}

	cancelled := false
	for _, fn := range handlers {
		if err := p.State.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, eventTable); err != nil {
			fmt.Printf("[Plugin:%s] event handler error for %s: %v\n", p.Meta.Name, eventName, err)
			continue
		}

		val := eventTable.RawGetString("cancelled")
		if val == lua.LTrue {
			cancelled = true
		}
	}

	return cancelled
}
