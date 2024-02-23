// Code generated by counterfeiter. DO NOT EDIT.
package connectfakes

import (
	"context"
	"sync"

	"github.com/ivan1993spb/snake-bot/internal/connect"
)

type FakeConnection struct {
	CloseStub        func(context.Context) error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
		arg1 context.Context
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	ReceiveStub        func(context.Context) ([]byte, error)
	receiveMutex       sync.RWMutex
	receiveArgsForCall []struct {
		arg1 context.Context
	}
	receiveReturns struct {
		result1 []byte
		result2 error
	}
	receiveReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SendStub        func(context.Context, interface{}) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 context.Context
		arg2 interface{}
	}
	sendReturns struct {
		result1 error
	}
	sendReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConnection) Close(arg1 context.Context) error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{arg1})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConnection) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeConnection) CloseCalls(stub func(context.Context) error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeConnection) CloseArgsForCall(i int) context.Context {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	argsForCall := fake.closeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConnection) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConnection) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
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

func (fake *FakeConnection) Receive(arg1 context.Context) ([]byte, error) {
	fake.receiveMutex.Lock()
	ret, specificReturn := fake.receiveReturnsOnCall[len(fake.receiveArgsForCall)]
	fake.receiveArgsForCall = append(fake.receiveArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.ReceiveStub
	fakeReturns := fake.receiveReturns
	fake.recordInvocation("Receive", []interface{}{arg1})
	fake.receiveMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeConnection) ReceiveCallCount() int {
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	return len(fake.receiveArgsForCall)
}

func (fake *FakeConnection) ReceiveCalls(stub func(context.Context) ([]byte, error)) {
	fake.receiveMutex.Lock()
	defer fake.receiveMutex.Unlock()
	fake.ReceiveStub = stub
}

func (fake *FakeConnection) ReceiveArgsForCall(i int) context.Context {
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	argsForCall := fake.receiveArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConnection) ReceiveReturns(result1 []byte, result2 error) {
	fake.receiveMutex.Lock()
	defer fake.receiveMutex.Unlock()
	fake.ReceiveStub = nil
	fake.receiveReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeConnection) ReceiveReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.receiveMutex.Lock()
	defer fake.receiveMutex.Unlock()
	fake.ReceiveStub = nil
	if fake.receiveReturnsOnCall == nil {
		fake.receiveReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.receiveReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeConnection) Send(arg1 context.Context, arg2 interface{}) error {
	fake.sendMutex.Lock()
	ret, specificReturn := fake.sendReturnsOnCall[len(fake.sendArgsForCall)]
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 context.Context
		arg2 interface{}
	}{arg1, arg2})
	stub := fake.SendStub
	fakeReturns := fake.sendReturns
	fake.recordInvocation("Send", []interface{}{arg1, arg2})
	fake.sendMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeConnection) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *FakeConnection) SendCalls(stub func(context.Context, interface{}) error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = stub
}

func (fake *FakeConnection) SendArgsForCall(i int) (context.Context, interface{}) {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	argsForCall := fake.sendArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeConnection) SendReturns(result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConnection) SendReturnsOnCall(i int, result1 error) {
	fake.sendMutex.Lock()
	defer fake.sendMutex.Unlock()
	fake.SendStub = nil
	if fake.sendReturnsOnCall == nil {
		fake.sendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConnection) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.receiveMutex.RLock()
	defer fake.receiveMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConnection) recordInvocation(key string, args []interface{}) {
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

var _ connect.Connection = new(FakeConnection)
