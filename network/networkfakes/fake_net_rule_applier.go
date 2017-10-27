// Code generated by counterfeiter. DO NOT EDIT.
package networkfakes

import (
	"sync"

	"code.cloudfoundry.org/winc/netrules"
	"code.cloudfoundry.org/winc/network"
	"github.com/Microsoft/hcsshim"
)

type FakeNetRuleApplier struct {
	InStub        func(netrules.NetIn) (hcsshim.NatPolicy, error)
	inMutex       sync.RWMutex
	inArgsForCall []struct {
		arg1 netrules.NetIn
	}
	inReturns struct {
		result1 hcsshim.NatPolicy
		result2 error
	}
	inReturnsOnCall map[int]struct {
		result1 hcsshim.NatPolicy
		result2 error
	}
	OutStub        func(netrules.NetOut, hcsshim.HNSEndpoint) error
	outMutex       sync.RWMutex
	outArgsForCall []struct {
		arg1 netrules.NetOut
		arg2 hcsshim.HNSEndpoint
	}
	outReturns struct {
		result1 error
	}
	outReturnsOnCall map[int]struct {
		result1 error
	}
	MTUStub        func(string, int) error
	mTUMutex       sync.RWMutex
	mTUArgsForCall []struct {
		arg1 string
		arg2 int
	}
	mTUReturns struct {
		result1 error
	}
	mTUReturnsOnCall map[int]struct {
		result1 error
	}
	CleanupStub        func() error
	cleanupMutex       sync.RWMutex
	cleanupArgsForCall []struct{}
	cleanupReturns     struct {
		result1 error
	}
	cleanupReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNetRuleApplier) In(arg1 netrules.NetIn) (hcsshim.NatPolicy, error) {
	fake.inMutex.Lock()
	ret, specificReturn := fake.inReturnsOnCall[len(fake.inArgsForCall)]
	fake.inArgsForCall = append(fake.inArgsForCall, struct {
		arg1 netrules.NetIn
	}{arg1})
	fake.recordInvocation("In", []interface{}{arg1})
	fake.inMutex.Unlock()
	if fake.InStub != nil {
		return fake.InStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.inReturns.result1, fake.inReturns.result2
}

func (fake *FakeNetRuleApplier) InCallCount() int {
	fake.inMutex.RLock()
	defer fake.inMutex.RUnlock()
	return len(fake.inArgsForCall)
}

func (fake *FakeNetRuleApplier) InArgsForCall(i int) netrules.NetIn {
	fake.inMutex.RLock()
	defer fake.inMutex.RUnlock()
	return fake.inArgsForCall[i].arg1
}

func (fake *FakeNetRuleApplier) InReturns(result1 hcsshim.NatPolicy, result2 error) {
	fake.InStub = nil
	fake.inReturns = struct {
		result1 hcsshim.NatPolicy
		result2 error
	}{result1, result2}
}

func (fake *FakeNetRuleApplier) InReturnsOnCall(i int, result1 hcsshim.NatPolicy, result2 error) {
	fake.InStub = nil
	if fake.inReturnsOnCall == nil {
		fake.inReturnsOnCall = make(map[int]struct {
			result1 hcsshim.NatPolicy
			result2 error
		})
	}
	fake.inReturnsOnCall[i] = struct {
		result1 hcsshim.NatPolicy
		result2 error
	}{result1, result2}
}

func (fake *FakeNetRuleApplier) Out(arg1 netrules.NetOut, arg2 hcsshim.HNSEndpoint) error {
	fake.outMutex.Lock()
	ret, specificReturn := fake.outReturnsOnCall[len(fake.outArgsForCall)]
	fake.outArgsForCall = append(fake.outArgsForCall, struct {
		arg1 netrules.NetOut
		arg2 hcsshim.HNSEndpoint
	}{arg1, arg2})
	fake.recordInvocation("Out", []interface{}{arg1, arg2})
	fake.outMutex.Unlock()
	if fake.OutStub != nil {
		return fake.OutStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.outReturns.result1
}

func (fake *FakeNetRuleApplier) OutCallCount() int {
	fake.outMutex.RLock()
	defer fake.outMutex.RUnlock()
	return len(fake.outArgsForCall)
}

func (fake *FakeNetRuleApplier) OutArgsForCall(i int) (netrules.NetOut, hcsshim.HNSEndpoint) {
	fake.outMutex.RLock()
	defer fake.outMutex.RUnlock()
	return fake.outArgsForCall[i].arg1, fake.outArgsForCall[i].arg2
}

func (fake *FakeNetRuleApplier) OutReturns(result1 error) {
	fake.OutStub = nil
	fake.outReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNetRuleApplier) OutReturnsOnCall(i int, result1 error) {
	fake.OutStub = nil
	if fake.outReturnsOnCall == nil {
		fake.outReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.outReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeNetRuleApplier) MTU(arg1 string, arg2 int) error {
	fake.mTUMutex.Lock()
	ret, specificReturn := fake.mTUReturnsOnCall[len(fake.mTUArgsForCall)]
	fake.mTUArgsForCall = append(fake.mTUArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("MTU", []interface{}{arg1, arg2})
	fake.mTUMutex.Unlock()
	if fake.MTUStub != nil {
		return fake.MTUStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.mTUReturns.result1
}

func (fake *FakeNetRuleApplier) MTUCallCount() int {
	fake.mTUMutex.RLock()
	defer fake.mTUMutex.RUnlock()
	return len(fake.mTUArgsForCall)
}

func (fake *FakeNetRuleApplier) MTUArgsForCall(i int) (string, int) {
	fake.mTUMutex.RLock()
	defer fake.mTUMutex.RUnlock()
	return fake.mTUArgsForCall[i].arg1, fake.mTUArgsForCall[i].arg2
}

func (fake *FakeNetRuleApplier) MTUReturns(result1 error) {
	fake.MTUStub = nil
	fake.mTUReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNetRuleApplier) MTUReturnsOnCall(i int, result1 error) {
	fake.MTUStub = nil
	if fake.mTUReturnsOnCall == nil {
		fake.mTUReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.mTUReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeNetRuleApplier) Cleanup() error {
	fake.cleanupMutex.Lock()
	ret, specificReturn := fake.cleanupReturnsOnCall[len(fake.cleanupArgsForCall)]
	fake.cleanupArgsForCall = append(fake.cleanupArgsForCall, struct{}{})
	fake.recordInvocation("Cleanup", []interface{}{})
	fake.cleanupMutex.Unlock()
	if fake.CleanupStub != nil {
		return fake.CleanupStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.cleanupReturns.result1
}

func (fake *FakeNetRuleApplier) CleanupCallCount() int {
	fake.cleanupMutex.RLock()
	defer fake.cleanupMutex.RUnlock()
	return len(fake.cleanupArgsForCall)
}

func (fake *FakeNetRuleApplier) CleanupReturns(result1 error) {
	fake.CleanupStub = nil
	fake.cleanupReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNetRuleApplier) CleanupReturnsOnCall(i int, result1 error) {
	fake.CleanupStub = nil
	if fake.cleanupReturnsOnCall == nil {
		fake.cleanupReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.cleanupReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeNetRuleApplier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.inMutex.RLock()
	defer fake.inMutex.RUnlock()
	fake.outMutex.RLock()
	defer fake.outMutex.RUnlock()
	fake.mTUMutex.RLock()
	defer fake.mTUMutex.RUnlock()
	fake.cleanupMutex.RLock()
	defer fake.cleanupMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeNetRuleApplier) recordInvocation(key string, args []interface{}) {
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

var _ network.NetRuleApplier = new(FakeNetRuleApplier)
