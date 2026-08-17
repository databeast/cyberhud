package ui

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

// --- Helpers for input dispatcher tests ---

// mockHandler is a test action handler that records calls and returns a fixed result.
type mockHandler struct {
	calls  []action.Action
	result action.Result
}

func (h *mockHandler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	h.calls = append(h.calls, act)
	return h.result
}

// mockInstance implements region.ModeInstance for dispatcher tests.
type mockInstance struct {
	id      string
	handler action.Handler
}

func (m *mockInstance) ID() string                    { return m.id }
func (m *mockInstance) Activate()                     {}
func (m *mockInstance) Deactivate()                   {}
func (m *mockInstance) ActionHandler() action.Handler { return m.handler }
func (m *mockInstance) BuildView() style.ViewData {
	return style.ViewData{}
}

// RenderCacheKey satisfies region.ModeInstance, which requires a uint32 change
// signature. This mock previously returned a string, which no longer matched the
// interface and left the whole ui test package uncompilable.
func (m *mockInstance) RenderCacheKey() uint32 { return 0 }

// inputDispatchTestRegionWithInstance creates a Region with a mock ModeInstance.
func inputDispatchTestRegionWithInstance(name string, handler action.Handler) *region.Region {
	bounds := image.Rect(0, 0, 128, 128)
	surf := surface.New(bounds)
	r := region.NewRegion(name, bounds, surf)

	inst := &mockInstance{id: "test-mode", handler: handler}

	// Configure a ModeFactory that returns our mock instance.
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		return inst, true
	})
	// Trigger SetMode to install the instance.
	_ = r.SetMode("test-mode")

	r.SetInputFocus(true)
	return r
}

// inputDispatchTestRegionNoInstance creates a Region with no active instance (Instance() returns nil).
func inputDispatchTestRegionNoInstance(name string) *region.Region {
	bounds := image.Rect(0, 0, 128, 128)
	surf := surface.New(bounds)
	r := region.NewRegion(name, bounds, surf)
	// No ModeFactory, no SetMode → Instance() returns nil.
	r.SetInputFocus(true)
	return r
}

// recordingMapper returns an InputMapper that records calls and returns the given action.
func recordingMapper(returnAction action.Action) (InputMapper, *[]input.Event) {
	var calls []input.Event
	mapper := func(ev input.Event) action.Action {
		calls = append(calls, ev)
		return returnAction
	}
	return mapper, &calls
}

// --- Tests: Press events map to actions and invoke HandleAction (Req 10.1) ---

func TestRegionInputDispatcher_PressEvent_InvokesHandleAction(t *testing.T) {
	// A press event that maps to a valid action should reach the instance's handler.
	handler := &mockHandler{result: action.Result{}}
	mapper, mapperCalls := recordingMapper(action.Up)

	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	ev := input.Event{Key: input.JoyUp, Type: input.Press}
	dispatcher.Dispatch(r, ev)

	// The mapper should have been called with the event.
	if len(*mapperCalls) != 1 {
		t.Fatalf("expected mapper called once, got %d calls", len(*mapperCalls))
	}
	if (*mapperCalls)[0] != ev {
		t.Fatalf("mapper received %v, want %v", (*mapperCalls)[0], ev)
	}

	// The handler should have been invoked.
	if len(handler.calls) != 1 {
		t.Fatalf("expected handler called once, got %d calls", len(handler.calls))
	}
	if handler.calls[0] != action.Up {
		t.Fatalf("handler received %v, want %v", handler.calls[0], action.Up)
	}
}

func TestRegionInputDispatcher_PressEvent_MapsToDifferentActions(t *testing.T) {
	// Verify that different mapped actions are dispatched correctly.
	actions := []action.Action{action.Up, action.Down, action.Primary, action.Secondary}

	for _, act := range actions {
		t.Run(act.String(), func(t *testing.T) {
			handler := &mockHandler{result: action.Result{}}
			mapper, _ := recordingMapper(act)
			dispatcher := NewRegionInputDispatcher(nil, mapper)

			r := inputDispatchTestRegionWithInstance("test", handler)

			ev := input.Event{Key: input.Key1, Type: input.Press}
			dispatcher.Dispatch(r, ev)

			if len(handler.calls) != 1 {
				t.Fatalf("expected handler called once, got %d calls", len(handler.calls))
			}
			if handler.calls[0] != act {
				t.Fatalf("handler received %v, want %v", handler.calls[0], act)
			}
		})
	}
}

// --- Tests: Release/hold events are discarded (Req 10.6 symmetry) ---

func TestRegionInputDispatcher_ReleaseEvent_Discarded(t *testing.T) {
	// Release events should be filtered without calling the mapper.
	handler := &mockHandler{result: action.Result{}}
	mapper, mapperCalls := recordingMapper(action.Up)

	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	ev := input.Event{Key: input.JoyUp, Type: input.Release}
	dispatcher.Dispatch(r, ev)

	// The mapper should NOT have been called.
	if len(*mapperCalls) != 0 {
		t.Fatalf("expected mapper not called for release event, got %d calls", len(*mapperCalls))
	}
	// The handler should NOT have been called.
	if len(handler.calls) != 0 {
		t.Fatalf("expected handler not called for release event, got %d calls", len(handler.calls))
	}
}

func TestRegionInputDispatcher_NonPressEventTypes_AllDiscarded(t *testing.T) {
	// Any event type that is not Press should be discarded.
	handler := &mockHandler{result: action.Result{}}
	mapper, mapperCalls := recordingMapper(action.Down)

	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	nonPressTypes := []input.EventType{input.Release, input.EventType(2), input.EventType(99)}
	for _, evType := range nonPressTypes {
		ev := input.Event{Key: input.Key1, Type: evType}
		dispatcher.Dispatch(r, ev)
	}

	if len(*mapperCalls) != 0 {
		t.Fatalf("expected mapper not called for non-press events, got %d calls", len(*mapperCalls))
	}
	if len(handler.calls) != 0 {
		t.Fatalf("expected handler not called for non-press events, got %d calls", len(handler.calls))
	}
}

// --- Tests: No-action mapping results in no dispatch (Req 10.2 precondition) ---

func TestRegionInputDispatcher_NoActionMapping_NoDispatch(t *testing.T) {
	// When the mapper returns action.None, HandleAction should not be invoked.
	handler := &mockHandler{result: action.Result{}}
	mapper, mapperCalls := recordingMapper(action.None)

	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	ev := input.Event{Key: input.Key3, Type: input.Press}
	dispatcher.Dispatch(r, ev)

	// Mapper is called (to determine the action) ...
	if len(*mapperCalls) != 1 {
		t.Fatalf("expected mapper called once, got %d calls", len(*mapperCalls))
	}
	// ... but handler should not be called since action is None.
	if len(handler.calls) != 0 {
		t.Fatalf("expected handler not called for None action, got %d calls", len(handler.calls))
	}
}

func TestRegionInputDispatcher_NilMapper_NoDispatch(t *testing.T) {
	// When inputMapper is nil, the event maps to action.None and is discarded.
	handler := &mockHandler{result: action.Result{}}
	dispatcher := NewRegionInputDispatcher(nil, nil)

	r := inputDispatchTestRegionWithInstance("main", handler)

	ev := input.Event{Key: input.JoyDown, Type: input.Press}
	dispatcher.Dispatch(r, ev)

	// Handler should not be called since no mapper → action.None.
	if len(handler.calls) != 0 {
		t.Fatalf("expected handler not called with nil mapper, got %d calls", len(handler.calls))
	}
}

// --- Tests: Nil instance discards action without error (Req 10.6) ---

func TestRegionInputDispatcher_NilInstance_EventConsumed(t *testing.T) {
	// When Instance() returns nil, the action should be discarded silently.
	mapper, _ := recordingMapper(action.Up)
	dispatcher := NewRegionInputDispatcher(nil, mapper)

	r := inputDispatchTestRegionNoInstance("main")

	ev := input.Event{Key: input.JoyUp, Type: input.Press}

	// Should not panic even with nil instance.
	dispatcher.Dispatch(r, ev)
}

// --- Tests: Nil ActionHandler discards action without error (Req 10.2) ---

func TestRegionInputDispatcher_NilActionHandler_EventConsumed(t *testing.T) {
	// When ActionHandler() returns nil, the action should be discarded silently.
	mapper, _ := recordingMapper(action.Up)
	dispatcher := NewRegionInputDispatcher(nil, mapper)

	// Create a region with an instance whose ActionHandler returns nil.
	r := inputDispatchTestRegionWithInstance("main", nil)

	ev := input.Event{Key: input.JoyUp, Type: input.Press}

	// Should not panic with nil handler.
	dispatcher.Dispatch(r, ev)
}

// --- Tests: Navigate result triggers SetMode (Req 10.3) ---

func TestRegionInputDispatcher_NavigateResult_CallsSetMode(t *testing.T) {
	// When the handler returns a Navigate result, SetMode should be called
	// on the region with the target mode ID.
	handler := &mockHandler{result: action.Result{Navigate: "clock"}}
	mapper, _ := recordingMapper(action.Primary)
	dispatcher := NewRegionInputDispatcher(nil, mapper)

	// Set up a region with a ModeFactory that tracks mode changes.
	bounds := image.Rect(0, 0, 128, 128)
	surf := surface.New(bounds)
	r := region.NewRegion("nav-test", bounds, surf)

	var setModeRequests []string
	r.SetModeFactory(func(id string, hints textlayout.TextHints) (region.ModeInstance, bool) {
		setModeRequests = append(setModeRequests, id)
		inst := &mockInstance{id: id, handler: handler}
		return inst, true
	})
	_ = r.SetMode("dashboard") // initial mode
	setModeRequests = nil      // reset to only track dispatch-triggered SetMode

	r.SetInputFocus(true)

	ev := input.Event{Key: input.Key1, Type: input.Press}
	dispatcher.Dispatch(r, ev)

	// SetMode should have been called with "clock".
	if len(setModeRequests) != 1 {
		t.Fatalf("expected 1 SetMode call, got %d", len(setModeRequests))
	}
	if setModeRequests[0] != "clock" {
		t.Fatalf("SetMode called with %q, want %q", setModeRequests[0], "clock")
	}
}

// --- Tests: CursorDelta and Dirty are handled gracefully (Req 10.4, 10.5) ---

func TestRegionInputDispatcher_CursorDeltaResult_NoExternalAction(t *testing.T) {
	// CursorDelta results should be processed without error. The instance
	// manages its own scroll state internally.
	handler := &mockHandler{result: action.Result{CursorDelta: 1}}
	mapper, _ := recordingMapper(action.Down)
	dispatcher := NewRegionInputDispatcher(nil, mapper)

	r := inputDispatchTestRegionWithInstance("main", handler)

	ev := input.Event{Key: input.JoyDown, Type: input.Press}

	// Should not panic — CursorDelta is handled internally by the instance.
	dispatcher.Dispatch(r, ev)
}

func TestRegionInputDispatcher_DirtyResult_NoExternalAction(t *testing.T) {
	// Dirty results should be processed without error. The render loop will
	// pick up changes on the next tick.
	handler := &mockHandler{result: action.Result{Dirty: true}}
	mapper, _ := recordingMapper(action.Primary)
	dispatcher := NewRegionInputDispatcher(nil, mapper)

	r := inputDispatchTestRegionWithInstance("main", handler)

	ev := input.Event{Key: input.Key1, Type: input.Press}

	// Should not panic — Dirty is a no-op for the dispatcher.
	dispatcher.Dispatch(r, ev)
}

// --- Tests: Multiple events and mixed types ---

func TestRegionInputDispatcher_MultipleEvents_IndependentDispatch(t *testing.T) {
	// Multiple press events should each be independently dispatched.
	handler := &mockHandler{result: action.Result{}}
	callCount := 0
	mapper := func(ev input.Event) action.Action {
		callCount++
		return action.Down
	}
	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	events := []input.Event{
		{Key: input.JoyDown, Type: input.Press},
		{Key: input.JoyDown, Type: input.Press},
		{Key: input.JoyUp, Type: input.Press},
	}

	for _, ev := range events {
		dispatcher.Dispatch(r, ev)
	}

	if callCount != 3 {
		t.Fatalf("expected mapper called 3 times, got %d", callCount)
	}
	if len(handler.calls) != 3 {
		t.Fatalf("expected handler called 3 times, got %d", len(handler.calls))
	}
}

func TestRegionInputDispatcher_MixedEventTypes_OnlyPressDispatched(t *testing.T) {
	// A sequence of mixed event types should only dispatch press events.
	handler := &mockHandler{result: action.Result{}}
	callCount := 0
	mapper := func(ev input.Event) action.Action {
		callCount++
		return action.Primary
	}
	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	events := []input.Event{
		{Key: input.Key1, Type: input.Press},   // dispatched
		{Key: input.Key1, Type: input.Release}, // discarded
		{Key: input.Key2, Type: input.Press},   // dispatched
		{Key: input.Key2, Type: input.Release}, // discarded
		{Key: input.Key3, Type: input.Press},   // dispatched
	}

	for _, ev := range events {
		dispatcher.Dispatch(r, ev)
	}

	if callCount != 3 {
		t.Fatalf("expected 3 press events dispatched, got %d", callCount)
	}
	if len(handler.calls) != 3 {
		t.Fatalf("expected handler called 3 times, got %d", len(handler.calls))
	}
}

func TestRegionInputDispatcher_AllKeys_PressEventsAccepted(t *testing.T) {
	// All defined keys should be accepted as press events.
	keys := []input.Key{
		input.Key1, input.Key2, input.Key3,
		input.JoyUp, input.JoyDown, input.JoyLeft, input.JoyRight, input.JoyPressed,
	}

	handler := &mockHandler{result: action.Result{}}
	callCount := 0
	mapper := func(ev input.Event) action.Action {
		callCount++
		return action.Up
	}
	dispatcher := NewRegionInputDispatcher(nil, mapper)
	r := inputDispatchTestRegionWithInstance("main", handler)

	for _, key := range keys {
		ev := input.Event{Key: key, Type: input.Press}
		dispatcher.Dispatch(r, ev)
	}

	if callCount != len(keys) {
		t.Fatalf("expected %d mapper calls (one per key), got %d", len(keys), callCount)
	}
	if len(handler.calls) != len(keys) {
		t.Fatalf("expected %d handler calls, got %d", len(keys), len(handler.calls))
	}
}
