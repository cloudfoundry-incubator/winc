// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/winc/hcs"
	"code.cloudfoundry.org/winc/runtime"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type StateManager struct {
	InitializeStub        func(string) error
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
		arg1 string
	}
	initializeReturns struct {
		result1 error
	}
	initializeReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func() error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct{}
	deleteReturns     struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	SetFailureStub        func() error
	setFailureMutex       sync.RWMutex
	setFailureArgsForCall []struct{}
	setFailureReturns     struct {
		result1 error
	}
	setFailureReturnsOnCall map[int]struct {
		result1 error
	}
	SetSuccessStub        func(hcs.Process) error
	setSuccessMutex       sync.RWMutex
	setSuccessArgsForCall []struct {
		arg1 hcs.Process
	}
	setSuccessReturns struct {
		result1 error
	}
	setSuccessReturnsOnCall map[int]struct {
		result1 error
	}
	StateStub        func() (*specs.State, error)
	stateMutex       sync.RWMutex
	stateArgsForCall []struct{}
	stateReturns     struct {
		result1 *specs.State
		result2 error
	}
	stateReturnsOnCall map[int]struct {
		result1 *specs.State
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StateManager) Initialize(arg1 string) error {
	fake.initializeMutex.Lock()
	ret, specificReturn := fake.initializeReturnsOnCall[len(fake.initializeArgsForCall)]
	fake.initializeArgsForCall = append(fake.initializeArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Initialize", []interface{}{arg1})
	fake.initializeMutex.Unlock()
	if fake.InitializeStub != nil {
		return fake.InitializeStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.initializeReturns.result1
}

func (fake *StateManager) InitializeCallCount() int {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return len(fake.initializeArgsForCall)
}

func (fake *StateManager) InitializeArgsForCall(i int) string {
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	return fake.initializeArgsForCall[i].arg1
}

func (fake *StateManager) InitializeReturns(result1 error) {
	fake.InitializeStub = nil
	fake.initializeReturns = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) InitializeReturnsOnCall(i int, result1 error) {
	fake.InitializeStub = nil
	if fake.initializeReturnsOnCall == nil {
		fake.initializeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initializeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) Delete() error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct{}{})
	fake.recordInvocation("Delete", []interface{}{})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *StateManager) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *StateManager) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) SetFailure() error {
	fake.setFailureMutex.Lock()
	ret, specificReturn := fake.setFailureReturnsOnCall[len(fake.setFailureArgsForCall)]
	fake.setFailureArgsForCall = append(fake.setFailureArgsForCall, struct{}{})
	fake.recordInvocation("SetFailure", []interface{}{})
	fake.setFailureMutex.Unlock()
	if fake.SetFailureStub != nil {
		return fake.SetFailureStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setFailureReturns.result1
}

func (fake *StateManager) SetFailureCallCount() int {
	fake.setFailureMutex.RLock()
	defer fake.setFailureMutex.RUnlock()
	return len(fake.setFailureArgsForCall)
}

func (fake *StateManager) SetFailureReturns(result1 error) {
	fake.SetFailureStub = nil
	fake.setFailureReturns = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) SetFailureReturnsOnCall(i int, result1 error) {
	fake.SetFailureStub = nil
	if fake.setFailureReturnsOnCall == nil {
		fake.setFailureReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setFailureReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) SetSuccess(arg1 hcs.Process) error {
	fake.setSuccessMutex.Lock()
	ret, specificReturn := fake.setSuccessReturnsOnCall[len(fake.setSuccessArgsForCall)]
	fake.setSuccessArgsForCall = append(fake.setSuccessArgsForCall, struct {
		arg1 hcs.Process
	}{arg1})
	fake.recordInvocation("SetSuccess", []interface{}{arg1})
	fake.setSuccessMutex.Unlock()
	if fake.SetSuccessStub != nil {
		return fake.SetSuccessStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setSuccessReturns.result1
}

func (fake *StateManager) SetSuccessCallCount() int {
	fake.setSuccessMutex.RLock()
	defer fake.setSuccessMutex.RUnlock()
	return len(fake.setSuccessArgsForCall)
}

func (fake *StateManager) SetSuccessArgsForCall(i int) hcs.Process {
	fake.setSuccessMutex.RLock()
	defer fake.setSuccessMutex.RUnlock()
	return fake.setSuccessArgsForCall[i].arg1
}

func (fake *StateManager) SetSuccessReturns(result1 error) {
	fake.SetSuccessStub = nil
	fake.setSuccessReturns = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) SetSuccessReturnsOnCall(i int, result1 error) {
	fake.SetSuccessStub = nil
	if fake.setSuccessReturnsOnCall == nil {
		fake.setSuccessReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setSuccessReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *StateManager) State() (*specs.State, error) {
	fake.stateMutex.Lock()
	ret, specificReturn := fake.stateReturnsOnCall[len(fake.stateArgsForCall)]
	fake.stateArgsForCall = append(fake.stateArgsForCall, struct{}{})
	fake.recordInvocation("State", []interface{}{})
	fake.stateMutex.Unlock()
	if fake.StateStub != nil {
		return fake.StateStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.stateReturns.result1, fake.stateReturns.result2
}

func (fake *StateManager) StateCallCount() int {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	return len(fake.stateArgsForCall)
}

func (fake *StateManager) StateReturns(result1 *specs.State, result2 error) {
	fake.StateStub = nil
	fake.stateReturns = struct {
		result1 *specs.State
		result2 error
	}{result1, result2}
}

func (fake *StateManager) StateReturnsOnCall(i int, result1 *specs.State, result2 error) {
	fake.StateStub = nil
	if fake.stateReturnsOnCall == nil {
		fake.stateReturnsOnCall = make(map[int]struct {
			result1 *specs.State
			result2 error
		})
	}
	fake.stateReturnsOnCall[i] = struct {
		result1 *specs.State
		result2 error
	}{result1, result2}
}

func (fake *StateManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.initializeMutex.RLock()
	defer fake.initializeMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.setFailureMutex.RLock()
	defer fake.setFailureMutex.RUnlock()
	fake.setSuccessMutex.RLock()
	defer fake.setSuccessMutex.RUnlock()
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StateManager) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ runtime.StateManager = new(StateManager)
