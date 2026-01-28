package event

const (
	PriorityMonitor = 0
	PriorityHighest = 1
	PriorityHigh    = 2
	PriorityNormal  = 3
	PriorityLow     = 4
	PriorityLowest  = 5
)

type Event interface {
	Name() string
	GetHandlers() *HandlerList
}

type Cancellable interface {
	Event
	IsCancelled() bool
	SetCancelled(cancelled bool)
}

type BaseEvent struct {
	name      string
	cancelled bool
}

func NewBaseEvent(name string) *BaseEvent {
	return &BaseEvent{name: name}
}

func (e *BaseEvent) Name() string {
	return e.name
}

func (e *BaseEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *BaseEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

type Handler func(Event)

type RegisteredListener struct {
	handler         Handler
	priority        int
	ignoreCancelled bool
	pluginName      string
}

func NewRegisteredListener(handler Handler, priority int, ignoreCancelled bool, pluginName string) *RegisteredListener {
	return &RegisteredListener{
		handler:         handler,
		priority:        priority,
		ignoreCancelled: ignoreCancelled,
		pluginName:      pluginName,
	}
}

func (r *RegisteredListener) Priority() int {
	return r.priority
}

func (r *RegisteredListener) PluginName() string {
	return r.pluginName
}

func (r *RegisteredListener) Call(event Event) {

	if r.ignoreCancelled {
		if c, ok := event.(Cancellable); ok && c.IsCancelled() {
			return
		}
	}
	r.handler(event)
}

type HandlerList struct {
	slots    map[int][]*RegisteredListener
	handlers []*RegisteredListener
	dirty    bool
}

func NewHandlerList() *HandlerList {
	return &HandlerList{
		slots: map[int][]*RegisteredListener{
			PriorityLowest:  {},
			PriorityLow:     {},
			PriorityNormal:  {},
			PriorityHigh:    {},
			PriorityHighest: {},
			PriorityMonitor: {},
		},
		dirty: true,
	}
}

func (h *HandlerList) Register(listener *RegisteredListener) {
	priority := listener.Priority()
	if priority < PriorityMonitor || priority > PriorityLowest {
		priority = PriorityNormal
	}
	h.slots[priority] = append(h.slots[priority], listener)
	h.dirty = true
}

func (h *HandlerList) Unregister(listener *RegisteredListener) {
	for priority, listeners := range h.slots {
		for i, l := range listeners {
			if l == listener {
				h.slots[priority] = append(listeners[:i], listeners[i+1:]...)
				h.dirty = true
				return
			}
		}
	}
}

func (h *HandlerList) UnregisterPlugin(pluginName string) {
	for priority, listeners := range h.slots {
		filtered := make([]*RegisteredListener, 0, len(listeners))
		for _, l := range listeners {
			if l.PluginName() != pluginName {
				filtered = append(filtered, l)
			}
		}
		if len(filtered) != len(listeners) {
			h.slots[priority] = filtered
			h.dirty = true
		}
	}
}

func (h *HandlerList) Clear() {
	for priority := range h.slots {
		h.slots[priority] = []*RegisteredListener{}
	}
	h.handlers = nil
	h.dirty = true
}

func (h *HandlerList) GetRegisteredListeners() []*RegisteredListener {
	if h.dirty {
		h.bake()
	}
	return h.handlers
}

func (h *HandlerList) bake() {
	h.handlers = make([]*RegisteredListener, 0)

	for priority := PriorityLowest; priority >= PriorityMonitor; priority-- {
		h.handlers = append(h.handlers, h.slots[priority]...)
	}
	h.dirty = false
}

type EventManager struct {
	handlers map[string]*HandlerList
}

func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string]*HandlerList),
	}
}

var globalManager = NewEventManager()

func GetGlobalManager() *EventManager {
	return globalManager
}

func (m *EventManager) RegisterHandler(eventName string, handler Handler, priority int, pluginName string) {
	m.RegisterHandlerEx(eventName, handler, priority, false, pluginName)
}

func (m *EventManager) RegisterHandlerEx(eventName string, handler Handler, priority int, ignoreCancelled bool, pluginName string) {
	if _, ok := m.handlers[eventName]; !ok {
		m.handlers[eventName] = NewHandlerList()
	}
	listener := NewRegisteredListener(handler, priority, ignoreCancelled, pluginName)
	m.handlers[eventName].Register(listener)
}

func (m *EventManager) UnregisterPlugin(pluginName string) {
	for _, handlerList := range m.handlers {
		handlerList.UnregisterPlugin(pluginName)
	}
}

func (m *EventManager) Call(event Event) {
	handlerList, ok := m.handlers[event.Name()]
	if !ok {
		return
	}

	for _, listener := range handlerList.GetRegisteredListeners() {
		listener.Call(event)
	}
}

func Call(event Event) {
	globalManager.Call(event)
}

func Register(eventName string, handler Handler, priority int, pluginName string) {
	globalManager.RegisterHandler(eventName, handler, priority, pluginName)
}
