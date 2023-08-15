// Code generated by mockery. DO NOT EDIT.

package persistermock

import mock "github.com/stretchr/testify/mock"

// RawFilePersister is an autogenerated mock type for the RawFilePersister type
type RawFilePersister struct {
	mock.Mock
}

type RawFilePersister_Expecter struct {
	mock *mock.Mock
}

func (_m *RawFilePersister) EXPECT() *RawFilePersister_Expecter {
	return &RawFilePersister_Expecter{mock: &_m.Mock}
}

// Persist provides a mock function with given fields: fileName, fileContent
func (_m *RawFilePersister) Persist(fileName string, fileContent []byte) error {
	ret := _m.Called(fileName, fileContent)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(fileName, fileContent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RawFilePersister_Persist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Persist'
type RawFilePersister_Persist_Call struct {
	*mock.Call
}

// Persist is a helper method to define mock.On call
//   - fileName string
//   - fileContent []byte
func (_e *RawFilePersister_Expecter) Persist(fileName interface{}, fileContent interface{}) *RawFilePersister_Persist_Call {
	return &RawFilePersister_Persist_Call{Call: _e.mock.On("Persist", fileName, fileContent)}
}

func (_c *RawFilePersister_Persist_Call) Run(run func(fileName string, fileContent []byte)) *RawFilePersister_Persist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte))
	})
	return _c
}

func (_c *RawFilePersister_Persist_Call) Return(_a0 error) *RawFilePersister_Persist_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RawFilePersister_Persist_Call) RunAndReturn(run func(string, []byte) error) *RawFilePersister_Persist_Call {
	_c.Call.Return(run)
	return _c
}

// NewRawFilePersister creates a new instance of RawFilePersister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRawFilePersister(t interface {
	mock.TestingT
	Cleanup(func())
}) *RawFilePersister {
	mock := &RawFilePersister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
