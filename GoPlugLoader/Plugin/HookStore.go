package Plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Cast"
	"github.com/MickMake/GoPlug/utils/Return"
)

// ---------------------------------------------------------------------------------------------------- //
// HookStore interface and methods

//
// HookStore - Getter/Setter for string map of interfaces{}
// ---------------------------------------------------------------------------------------------------- //
type HookStore interface {
	// NewHookStore - Set up the FuncMap structure.
	NewHookStore() Return.Error

	SetHookPlugin(plugin PluginDataInterface)
	GetHookReference() *HookStruct

	// GetIdentity() *GoPlugLoader.PluginIdentity
	// SetIdentity(identity *GoPlugLoader.PluginIdentity) Return.Error

	GetHookIdentity() string
	SetHookIdentity(identity string) Return.Error

	// HookExists - Check if a key exists.
	HookExists(hook string) bool

	// HookNotExists - Inverse of Exists()
	HookNotExists(hook string) bool

	// GetHook - Get a key's value.
	GetHook(hook string) *Hook
	GetHookName(name string) (string, Return.Error)
	GetHookFunction(name string) (HookFunction, Return.Error)
	GetHookArgs(name string) (HookArgs, Return.Error)

	ValidateHook(args ...any) Return.Error

	// SetHook - Set a key value pair.
	SetHook(name string, function HookFunction, args ...any) Return.Error

	// CountHooks - Return the number of entries.
	CountHooks() int

	// ListHooks - Get HookStruct.
	ListHooks() HookMap

	// PrintHooks - Get HookStruct.
	PrintHooks()

	// String - Stringer method.
	String() string
}

// NewHookStore - Create a HookStore interface structure instance.
//goland:noinspection GoUnusedExportedFunction
func NewHookStore() HookStore {
	return &HookStruct{
		Hooks: make(HookMap),
	}
}

//
// HookStruct
// ---------------------------------------------------------------------------------------------------- //
type HookStruct struct {
	Identity string       `json:"identity,omitempty"`
	Hooks    HookMap      `json:"-"`
	Master   bool         `json:"master,omitempty"`
	Error    Return.Error `json:"-"`
	plugin   PluginDataInterface
}

// NewHookStruct - Create a HookStruct structure instance.
func NewHookStruct() HookStruct {
	return HookStruct{
		Hooks: make(HookMap),
	}
}

// NewHookStore - Create a HookStore interface structure instance.
func (h *HookStruct) NewHookStore() Return.Error {
	h.Error = Return.Ok
	h.Hooks = make(map[string]*Hook)
	return Return.Ok
}

func (h *HookStruct) SetHookPlugin(plugin PluginDataInterface) {
	h.Error = Return.Ok
	h.plugin = plugin
}

func (h *HookStruct) GetHookReference() *HookStruct {
	h.Error = Return.Ok
	return h
}

func (h *HookStruct) SetHookIdentity(identity string) Return.Error {
	h.Error = Return.Ok
	h.Identity = identity
	return Return.Ok
}

func (h *HookStruct) GetHookIdentity() string {
	h.Error = Return.Ok
	return h.Identity
}

// HookExists - Check if a key exists.
func (h *HookStruct) HookExists(name string) bool {
	hook, _ := h.Hooks.Get(name)
	if hook == nil {
		return false
	}
	return true
}

// HookNotExists - Inverse of Exists()
func (h *HookStruct) HookNotExists(name string) bool {
	hook, _ := h.Hooks.Get(name)
	if hook == nil {
		return true
	}
	return false
}

// GetHook - Get a key's value.
func (h *HookStruct) GetHook(name string) *Hook {
	hook, _ := h.Hooks.Get(name)
	return hook
}

// GetHookName - Get a key's value.
func (h *HookStruct) GetHookName(name string) (string, Return.Error) {
	var hook *Hook
	hook, h.Error = h.Hooks.Get(name)
	if h.Error.IsError() {
		return "", h.Error
	}
	return hook.Name, h.Error
}

// GetHookFunction - Get a key's value.
func (h *HookStruct) GetHookFunction(name string) (HookFunction, Return.Error) {
	var hook *Hook
	hook, h.Error = h.Hooks.Get(name)
	if h.Error.IsError() {
		return nil, h.Error
	}
	return hook.function, h.Error
}

// GetHookArgs - Get a key's value.
func (h *HookStruct) GetHookArgs(name string) (HookArgs, Return.Error) {
	var hook *Hook
	hook, h.Error = h.Hooks.Get(name)
	if h.Error.IsError() {
		return HookArgs{}, h.Error
	}
	return hook.Args, h.Error
}

// SetHook - Set a key value pair.
func (h *HookStruct) SetHook(name string, function HookFunction, args ...any) Return.Error {
	h.Error = Return.Ok
	fp, fm := utils.GetPackageAndFunctionNameFromPointer(function)
	name = strings.TrimSpace(name)
	if name == "" {
		name = fm
	}

	var hookArgs HookArgs
	for _, a := range args {
		hookArgs = append(hookArgs, NewHookArg(a))
	}
	hook := &Hook{
		Name:     fp + "." + fm,
		function: function,
		Args:     hookArgs,
	}
	h.Hooks[name] = hook
	return h.Error
}

func (h *HookStruct) CallHook(name string, args ...any) (HookResponse, Return.Error) {
	h.Error = Return.Ok
	var resp HookResponse
	for range Only.Once {
		h.Error.SetPrefix("hook[%s]", name)

		hook := h.GetHook(name)
		if hook == nil {
			h.Error.SetError("hook '%s' not found", name)
			break
		}

		h.Error = hook.Args.Validate(args...)
		if h.Error.IsError() {
			break
		}

		resp, h.Error = hook.function(*h, args...)
	}
	return resp, h.Error
}

// CountHooks - Return the number of entries.
func (h *HookStruct) CountHooks() int {
	h.Error = Return.Ok
	return len(h.Hooks)
}

func (h *HookStruct) ListHooks() HookMap {
	h.Error = Return.Ok
	return h.Hooks
}

func (h *HookStruct) PrintHooks() {
	h.Error = Return.Ok
	fmt.Print(h.String())
}

// StringHooks - Stringer interface.
func (h HookStruct) String() string {
	var ret string
	ret += fmt.Sprintf("# Available function hooks from plugin '%s'\n", h.Identity)
	for name, hook := range h.Hooks {
		ret += fmt.Sprintf("\t[%s]: %s\n", name, hook)
	}
	return ret
}

// ValidateHook - .
func (h *HookStruct) ValidateHook(args ...any) Return.Error {
	for range Only.Once {
		h.Error = Return.Ok
		name := utils.GetCallerFunctionName(1)
		hook := h.GetHook(name)
		if hook == nil {
			h.Error.SetError("hook function mismatch: looking for %s", name)
			break
		}

		if hook.function == nil {
			h.Error.SetError("hook function is nil: looking for %s", name)
			break
		}

		h.Error = hook.Args.Validate(args...)
		if h.Error.IsError() {
			break
		}
	}
	return h.Error
}

//
// HookMap
// ---------------------------------------------------------------------------------------------------- //
type HookMap map[string]*Hook

func (m *HookMap) Get(name string) (*Hook, Return.Error) {
	name = strings.TrimSpace(name)
	if value, ok := (*m)[name]; ok {
		return value, Return.Ok
	}
	return nil, Return.NewError("hook '%s' not found", name)
}

//
// Hook
// ---------------------------------------------------------------------------------------------------- //
type Hook struct {
	Name     string `json:"name,omitempty"`
	function HookFunction
	Args     HookArgs `json:"args,omitempty"`
}

func (h *Hook) Validate(args ...any) Return.Error {
	var err Return.Error

	for range Only.Once {
		if h == nil {
			err.SetError("Hook struct is nil")
			break
		}

		if h.function == nil {
			err.SetError("Hook function not defined")
			break
		}

		err = h.Args.Validate(args...)
	}

	return err
}

func (h Hook) String() string {
	// name := utils.GetPackageAndFunctionNameFromPointer(h.Function)
	return fmt.Sprintf("Function: %s(%s)", h.Name, h.Args)
}

//
// HookCallArgs
// ---------------------------------------------------------------------------------------------------- //
type HookCallArgs struct {
	Name string `json:"name,omitempty"`
	Args []any  `json:"args,omitempty"`
}

//
// HookFunction
// ---------------------------------------------------------------------------------------------------- //
type HookFunction func(hook HookStruct, args ...any) (HookResponse, Return.Error)

//
// HookArgs
// ---------------------------------------------------------------------------------------------------- //
type HookArgs []HookArg

//goland:noinspection GoUnusedExportedFunction
func NewHookArgs(args ...any) HookArgs {
	var ret HookArgs
	for _, arg := range args {
		ret = append(ret, HookArg(utils.GetTypeName(arg)))
	}
	return ret
}

func (a *HookArgs) Append(args ...HookArg) {
	*a = append(*a, args...)
}

func (a *HookArgs) Validate(args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		nargs := len(args)
		cargs := len(*a)
		if nargs > cargs {
			err.SetError("too many args, should be %d", cargs)
			break
		}
		if nargs < cargs {
			err.SetError("not enough args, should be %d", cargs)
			break
		}
		for index, arg := range args {
			targ := utils.GetTypeName(arg)
			if targ != string((*a)[index]) {
				err.SetError("args at position %d should be of type %s, not %s", index, (*a)[index], targ)
				break
			}
		}
	}
	return err
}

func (a *HookArgs) Count() int {
	return len(*a)
}

func (a HookArgs) String() string {
	var ret []string
	for _, arg := range a {
		ret = append(ret, string(arg))
	}
	return strings.Join(ret, ", ")
}

//
// HookArg
// ---------------------------------------------------------------------------------------------------- //
type HookArg json.RawMessage

func NewHookArg(arg any) HookArg {
	return HookArg(utils.GetTypeName(arg))
}

func (a HookArg) String() string {
	return string(a)
}

//
// HookResponse
// ---------------------------------------------------------------------------------------------------- //
type HookResponse struct {
	Value any    `json:"-"`
	Type  string `json:"type,omitempty"`
}

func NewHookResponse(response any) (HookResponse, Return.Error) {
	var ret HookResponse
	var err Return.Error
	ret.Value = response
	ret.Type = utils.GetTypeName(response)
	return ret, err
}

func HookResponseNil() (HookResponse, Return.Error) {
	return HookResponse{}, Return.Ok
}

func (r HookResponse) String() string {
	if r.Value == nil {
		return ""
	}
	return fmt.Sprintf("%v", r.Value)
}

func (r *HookResponse) Print() {
	if r.Value == nil {
		return
	}
	fmt.Println(r.String())
}

func (r *HookResponse) TypeMatches(Type string) bool {
	if r.Type == Type {
		return true
	}
	return false
}

func (r *HookResponse) AsString() string {
	ret := Cast.ToString(r.Value)
	return ret
}

func (r *HookResponse) AsByteArray() []byte {
	return r.Value.([]byte)
}

// ---------------------------------------------------------------------------------------------------- //

//goland:noinspection GoUnusedExportedFunction
func HookArgAsString(arg any) *string {
	value, err := Cast.ToStringE(arg)
	if err == nil {
		return &value
	}
	return nil
}

//goland:noinspection GoUnusedExportedFunction
func HookArgAsInt(arg any) *int {
	value, err := Cast.ToIntE(arg)
	if err == nil {
		return &value
	}
	return nil
}
