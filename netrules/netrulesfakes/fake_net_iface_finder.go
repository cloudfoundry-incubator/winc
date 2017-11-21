// Code generated by counterfeiter. DO NOT EDIT.
package netrulesfakes

import (
	"net"
	"sync"

	"code.cloudfoundry.org/winc/netrules"
)

type FakeNetIfaceFinder struct {
	ByNameStub        func(string) (*net.Interface, error)
	byNameMutex       sync.RWMutex
	byNameArgsForCall []struct {
		arg1 string
	}
	byNameReturns struct {
		result1 *net.Interface
		result2 error
	}
	byNameReturnsOnCall map[int]struct {
		result1 *net.Interface
		result2 error
	}
	ByIPStub        func(string) (*net.Interface, error)
	byIPMutex       sync.RWMutex
	byIPArgsForCall []struct {
		arg1 string
	}
	byIPReturns struct {
		result1 *net.Interface
		result2 error
	}
	byIPReturnsOnCall map[int]struct {
		result1 *net.Interface
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNetIfaceFinder) ByName(arg1 string) (*net.Interface, error) {
	fake.byNameMutex.Lock()
	ret, specificReturn := fake.byNameReturnsOnCall[len(fake.byNameArgsForCall)]
	fake.byNameArgsForCall = append(fake.byNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ByName", []interface{}{arg1})
	fake.byNameMutex.Unlock()
	if fake.ByNameStub != nil {
		return fake.ByNameStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.byNameReturns.result1, fake.byNameReturns.result2
}

func (fake *FakeNetIfaceFinder) ByNameCallCount() int {
	fake.byNameMutex.RLock()
	defer fake.byNameMutex.RUnlock()
	return len(fake.byNameArgsForCall)
}

func (fake *FakeNetIfaceFinder) ByNameArgsForCall(i int) string {
	fake.byNameMutex.RLock()
	defer fake.byNameMutex.RUnlock()
	return fake.byNameArgsForCall[i].arg1
}

func (fake *FakeNetIfaceFinder) ByNameReturns(result1 *net.Interface, result2 error) {
	fake.ByNameStub = nil
	fake.byNameReturns = struct {
		result1 *net.Interface
		result2 error
	}{result1, result2}
}

func (fake *FakeNetIfaceFinder) ByNameReturnsOnCall(i int, result1 *net.Interface, result2 error) {
	fake.ByNameStub = nil
	if fake.byNameReturnsOnCall == nil {
		fake.byNameReturnsOnCall = make(map[int]struct {
			result1 *net.Interface
			result2 error
		})
	}
	fake.byNameReturnsOnCall[i] = struct {
		result1 *net.Interface
		result2 error
	}{result1, result2}
}

func (fake *FakeNetIfaceFinder) ByIP(arg1 string) (*net.Interface, error) {
	fake.byIPMutex.Lock()
	ret, specificReturn := fake.byIPReturnsOnCall[len(fake.byIPArgsForCall)]
	fake.byIPArgsForCall = append(fake.byIPArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ByIP", []interface{}{arg1})
	fake.byIPMutex.Unlock()
	if fake.ByIPStub != nil {
		return fake.ByIPStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.byIPReturns.result1, fake.byIPReturns.result2
}

func (fake *FakeNetIfaceFinder) ByIPCallCount() int {
	fake.byIPMutex.RLock()
	defer fake.byIPMutex.RUnlock()
	return len(fake.byIPArgsForCall)
}

func (fake *FakeNetIfaceFinder) ByIPArgsForCall(i int) string {
	fake.byIPMutex.RLock()
	defer fake.byIPMutex.RUnlock()
	return fake.byIPArgsForCall[i].arg1
}

func (fake *FakeNetIfaceFinder) ByIPReturns(result1 *net.Interface, result2 error) {
	fake.ByIPStub = nil
	fake.byIPReturns = struct {
		result1 *net.Interface
		result2 error
	}{result1, result2}
}

func (fake *FakeNetIfaceFinder) ByIPReturnsOnCall(i int, result1 *net.Interface, result2 error) {
	fake.ByIPStub = nil
	if fake.byIPReturnsOnCall == nil {
		fake.byIPReturnsOnCall = make(map[int]struct {
			result1 *net.Interface
			result2 error
		})
	}
	fake.byIPReturnsOnCall[i] = struct {
		result1 *net.Interface
		result2 error
	}{result1, result2}
}

func (fake *FakeNetIfaceFinder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.byNameMutex.RLock()
	defer fake.byNameMutex.RUnlock()
	fake.byIPMutex.RLock()
	defer fake.byIPMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeNetIfaceFinder) recordInvocation(key string, args []interface{}) {
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

var _ netrules.NetIfaceFinder = new(FakeNetIfaceFinder)