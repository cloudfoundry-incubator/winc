// Code generated by counterfeiter. DO NOT EDIT.
package hcsfakes

import (
	"io"
	"sync"
	"time"

	"code.cloudfoundry.org/winc/hcs"
)

type FakeProcess struct {
	PidStub        func() int
	pidMutex       sync.RWMutex
	pidArgsForCall []struct{}
	pidReturns     struct {
		result1 int
	}
	pidReturnsOnCall map[int]struct {
		result1 int
	}
	KillStub        func() error
	killMutex       sync.RWMutex
	killArgsForCall []struct{}
	killReturns     struct {
		result1 error
	}
	killReturnsOnCall map[int]struct {
		result1 error
	}
	WaitStub        func() error
	waitMutex       sync.RWMutex
	waitArgsForCall []struct{}
	waitReturns     struct {
		result1 error
	}
	waitReturnsOnCall map[int]struct {
		result1 error
	}
	WaitTimeoutStub        func(time.Duration) error
	waitTimeoutMutex       sync.RWMutex
	waitTimeoutArgsForCall []struct {
		arg1 time.Duration
	}
	waitTimeoutReturns struct {
		result1 error
	}
	waitTimeoutReturnsOnCall map[int]struct {
		result1 error
	}
	ExitCodeStub        func() (int, error)
	exitCodeMutex       sync.RWMutex
	exitCodeArgsForCall []struct{}
	exitCodeReturns     struct {
		result1 int
		result2 error
	}
	exitCodeReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	ResizeConsoleStub        func(width, height uint16) error
	resizeConsoleMutex       sync.RWMutex
	resizeConsoleArgsForCall []struct {
		width  uint16
		height uint16
	}
	resizeConsoleReturns struct {
		result1 error
	}
	resizeConsoleReturnsOnCall map[int]struct {
		result1 error
	}
	StdioStub        func() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error)
	stdioMutex       sync.RWMutex
	stdioArgsForCall []struct{}
	stdioReturns     struct {
		result1 io.WriteCloser
		result2 io.ReadCloser
		result3 io.ReadCloser
		result4 error
	}
	stdioReturnsOnCall map[int]struct {
		result1 io.WriteCloser
		result2 io.ReadCloser
		result3 io.ReadCloser
		result4 error
	}
	CloseStdinStub        func() error
	closeStdinMutex       sync.RWMutex
	closeStdinArgsForCall []struct{}
	closeStdinReturns     struct {
		result1 error
	}
	closeStdinReturnsOnCall map[int]struct {
		result1 error
	}
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProcess) Pid() int {
	fake.pidMutex.Lock()
	ret, specificReturn := fake.pidReturnsOnCall[len(fake.pidArgsForCall)]
	fake.pidArgsForCall = append(fake.pidArgsForCall, struct{}{})
	fake.recordInvocation("Pid", []interface{}{})
	fake.pidMutex.Unlock()
	if fake.PidStub != nil {
		return fake.PidStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.pidReturns.result1
}

func (fake *FakeProcess) PidCallCount() int {
	fake.pidMutex.RLock()
	defer fake.pidMutex.RUnlock()
	return len(fake.pidArgsForCall)
}

func (fake *FakeProcess) PidReturns(result1 int) {
	fake.PidStub = nil
	fake.pidReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeProcess) PidReturnsOnCall(i int, result1 int) {
	fake.PidStub = nil
	if fake.pidReturnsOnCall == nil {
		fake.pidReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.pidReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeProcess) Kill() error {
	fake.killMutex.Lock()
	ret, specificReturn := fake.killReturnsOnCall[len(fake.killArgsForCall)]
	fake.killArgsForCall = append(fake.killArgsForCall, struct{}{})
	fake.recordInvocation("Kill", []interface{}{})
	fake.killMutex.Unlock()
	if fake.KillStub != nil {
		return fake.KillStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.killReturns.result1
}

func (fake *FakeProcess) KillCallCount() int {
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	return len(fake.killArgsForCall)
}

func (fake *FakeProcess) KillReturns(result1 error) {
	fake.KillStub = nil
	fake.killReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) KillReturnsOnCall(i int, result1 error) {
	fake.KillStub = nil
	if fake.killReturnsOnCall == nil {
		fake.killReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.killReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) Wait() error {
	fake.waitMutex.Lock()
	ret, specificReturn := fake.waitReturnsOnCall[len(fake.waitArgsForCall)]
	fake.waitArgsForCall = append(fake.waitArgsForCall, struct{}{})
	fake.recordInvocation("Wait", []interface{}{})
	fake.waitMutex.Unlock()
	if fake.WaitStub != nil {
		return fake.WaitStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.waitReturns.result1
}

func (fake *FakeProcess) WaitCallCount() int {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return len(fake.waitArgsForCall)
}

func (fake *FakeProcess) WaitReturns(result1 error) {
	fake.WaitStub = nil
	fake.waitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) WaitReturnsOnCall(i int, result1 error) {
	fake.WaitStub = nil
	if fake.waitReturnsOnCall == nil {
		fake.waitReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.waitReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) WaitTimeout(arg1 time.Duration) error {
	fake.waitTimeoutMutex.Lock()
	ret, specificReturn := fake.waitTimeoutReturnsOnCall[len(fake.waitTimeoutArgsForCall)]
	fake.waitTimeoutArgsForCall = append(fake.waitTimeoutArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	fake.recordInvocation("WaitTimeout", []interface{}{arg1})
	fake.waitTimeoutMutex.Unlock()
	if fake.WaitTimeoutStub != nil {
		return fake.WaitTimeoutStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.waitTimeoutReturns.result1
}

func (fake *FakeProcess) WaitTimeoutCallCount() int {
	fake.waitTimeoutMutex.RLock()
	defer fake.waitTimeoutMutex.RUnlock()
	return len(fake.waitTimeoutArgsForCall)
}

func (fake *FakeProcess) WaitTimeoutArgsForCall(i int) time.Duration {
	fake.waitTimeoutMutex.RLock()
	defer fake.waitTimeoutMutex.RUnlock()
	return fake.waitTimeoutArgsForCall[i].arg1
}

func (fake *FakeProcess) WaitTimeoutReturns(result1 error) {
	fake.WaitTimeoutStub = nil
	fake.waitTimeoutReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) WaitTimeoutReturnsOnCall(i int, result1 error) {
	fake.WaitTimeoutStub = nil
	if fake.waitTimeoutReturnsOnCall == nil {
		fake.waitTimeoutReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.waitTimeoutReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) ExitCode() (int, error) {
	fake.exitCodeMutex.Lock()
	ret, specificReturn := fake.exitCodeReturnsOnCall[len(fake.exitCodeArgsForCall)]
	fake.exitCodeArgsForCall = append(fake.exitCodeArgsForCall, struct{}{})
	fake.recordInvocation("ExitCode", []interface{}{})
	fake.exitCodeMutex.Unlock()
	if fake.ExitCodeStub != nil {
		return fake.ExitCodeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.exitCodeReturns.result1, fake.exitCodeReturns.result2
}

func (fake *FakeProcess) ExitCodeCallCount() int {
	fake.exitCodeMutex.RLock()
	defer fake.exitCodeMutex.RUnlock()
	return len(fake.exitCodeArgsForCall)
}

func (fake *FakeProcess) ExitCodeReturns(result1 int, result2 error) {
	fake.ExitCodeStub = nil
	fake.exitCodeReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeProcess) ExitCodeReturnsOnCall(i int, result1 int, result2 error) {
	fake.ExitCodeStub = nil
	if fake.exitCodeReturnsOnCall == nil {
		fake.exitCodeReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.exitCodeReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeProcess) ResizeConsole(width uint16, height uint16) error {
	fake.resizeConsoleMutex.Lock()
	ret, specificReturn := fake.resizeConsoleReturnsOnCall[len(fake.resizeConsoleArgsForCall)]
	fake.resizeConsoleArgsForCall = append(fake.resizeConsoleArgsForCall, struct {
		width  uint16
		height uint16
	}{width, height})
	fake.recordInvocation("ResizeConsole", []interface{}{width, height})
	fake.resizeConsoleMutex.Unlock()
	if fake.ResizeConsoleStub != nil {
		return fake.ResizeConsoleStub(width, height)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.resizeConsoleReturns.result1
}

func (fake *FakeProcess) ResizeConsoleCallCount() int {
	fake.resizeConsoleMutex.RLock()
	defer fake.resizeConsoleMutex.RUnlock()
	return len(fake.resizeConsoleArgsForCall)
}

func (fake *FakeProcess) ResizeConsoleArgsForCall(i int) (uint16, uint16) {
	fake.resizeConsoleMutex.RLock()
	defer fake.resizeConsoleMutex.RUnlock()
	return fake.resizeConsoleArgsForCall[i].width, fake.resizeConsoleArgsForCall[i].height
}

func (fake *FakeProcess) ResizeConsoleReturns(result1 error) {
	fake.ResizeConsoleStub = nil
	fake.resizeConsoleReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) ResizeConsoleReturnsOnCall(i int, result1 error) {
	fake.ResizeConsoleStub = nil
	if fake.resizeConsoleReturnsOnCall == nil {
		fake.resizeConsoleReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.resizeConsoleReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) Stdio() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	fake.stdioMutex.Lock()
	ret, specificReturn := fake.stdioReturnsOnCall[len(fake.stdioArgsForCall)]
	fake.stdioArgsForCall = append(fake.stdioArgsForCall, struct{}{})
	fake.recordInvocation("Stdio", []interface{}{})
	fake.stdioMutex.Unlock()
	if fake.StdioStub != nil {
		return fake.StdioStub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fake.stdioReturns.result1, fake.stdioReturns.result2, fake.stdioReturns.result3, fake.stdioReturns.result4
}

func (fake *FakeProcess) StdioCallCount() int {
	fake.stdioMutex.RLock()
	defer fake.stdioMutex.RUnlock()
	return len(fake.stdioArgsForCall)
}

func (fake *FakeProcess) StdioReturns(result1 io.WriteCloser, result2 io.ReadCloser, result3 io.ReadCloser, result4 error) {
	fake.StdioStub = nil
	fake.stdioReturns = struct {
		result1 io.WriteCloser
		result2 io.ReadCloser
		result3 io.ReadCloser
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeProcess) StdioReturnsOnCall(i int, result1 io.WriteCloser, result2 io.ReadCloser, result3 io.ReadCloser, result4 error) {
	fake.StdioStub = nil
	if fake.stdioReturnsOnCall == nil {
		fake.stdioReturnsOnCall = make(map[int]struct {
			result1 io.WriteCloser
			result2 io.ReadCloser
			result3 io.ReadCloser
			result4 error
		})
	}
	fake.stdioReturnsOnCall[i] = struct {
		result1 io.WriteCloser
		result2 io.ReadCloser
		result3 io.ReadCloser
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeProcess) CloseStdin() error {
	fake.closeStdinMutex.Lock()
	ret, specificReturn := fake.closeStdinReturnsOnCall[len(fake.closeStdinArgsForCall)]
	fake.closeStdinArgsForCall = append(fake.closeStdinArgsForCall, struct{}{})
	fake.recordInvocation("CloseStdin", []interface{}{})
	fake.closeStdinMutex.Unlock()
	if fake.CloseStdinStub != nil {
		return fake.CloseStdinStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.closeStdinReturns.result1
}

func (fake *FakeProcess) CloseStdinCallCount() int {
	fake.closeStdinMutex.RLock()
	defer fake.closeStdinMutex.RUnlock()
	return len(fake.closeStdinArgsForCall)
}

func (fake *FakeProcess) CloseStdinReturns(result1 error) {
	fake.CloseStdinStub = nil
	fake.closeStdinReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) CloseStdinReturnsOnCall(i int, result1 error) {
	fake.CloseStdinStub = nil
	if fake.closeStdinReturnsOnCall == nil {
		fake.closeStdinReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeStdinReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.closeReturns.result1
}

func (fake *FakeProcess) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeProcess) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) CloseReturnsOnCall(i int, result1 error) {
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeProcess) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pidMutex.RLock()
	defer fake.pidMutex.RUnlock()
	fake.killMutex.RLock()
	defer fake.killMutex.RUnlock()
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	fake.waitTimeoutMutex.RLock()
	defer fake.waitTimeoutMutex.RUnlock()
	fake.exitCodeMutex.RLock()
	defer fake.exitCodeMutex.RUnlock()
	fake.resizeConsoleMutex.RLock()
	defer fake.resizeConsoleMutex.RUnlock()
	fake.stdioMutex.RLock()
	defer fake.stdioMutex.RUnlock()
	fake.closeStdinMutex.RLock()
	defer fake.closeStdinMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeProcess) recordInvocation(key string, args []interface{}) {
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

var _ hcs.Process = new(FakeProcess)
