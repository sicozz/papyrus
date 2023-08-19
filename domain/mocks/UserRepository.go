// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/sicozz/papyrus/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

type UserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepository) EXPECT() *UserRepository_Expecter {
	return &UserRepository_Expecter{mock: &_m.Mock}
}

// ChgEmail provides a mock function with given fields: ctx, uname, email
func (_m *UserRepository) ChgEmail(ctx context.Context, uname string, email string) error {
	ret := _m.Called(ctx, uname, email)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, uname, email)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_ChgEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChgEmail'
type UserRepository_ChgEmail_Call struct {
	*mock.Call
}

// ChgEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
//   - email string
func (_e *UserRepository_Expecter) ChgEmail(ctx interface{}, uname interface{}, email interface{}) *UserRepository_ChgEmail_Call {
	return &UserRepository_ChgEmail_Call{Call: _e.mock.On("ChgEmail", ctx, uname, email)}
}

func (_c *UserRepository_ChgEmail_Call) Run(run func(ctx context.Context, uname string, email string)) *UserRepository_ChgEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_ChgEmail_Call) Return(_a0 error) *UserRepository_ChgEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ChgEmail_Call) RunAndReturn(run func(context.Context, string, string) error) *UserRepository_ChgEmail_Call {
	_c.Call.Return(run)
	return _c
}

// ChgLstname provides a mock function with given fields: ctx, uname, nLname
func (_m *UserRepository) ChgLstname(ctx context.Context, uname string, nLname string) error {
	ret := _m.Called(ctx, uname, nLname)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, uname, nLname)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_ChgLstname_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChgLstname'
type UserRepository_ChgLstname_Call struct {
	*mock.Call
}

// ChgLstname is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
//   - nLname string
func (_e *UserRepository_Expecter) ChgLstname(ctx interface{}, uname interface{}, nLname interface{}) *UserRepository_ChgLstname_Call {
	return &UserRepository_ChgLstname_Call{Call: _e.mock.On("ChgLstname", ctx, uname, nLname)}
}

func (_c *UserRepository_ChgLstname_Call) Run(run func(ctx context.Context, uname string, nLname string)) *UserRepository_ChgLstname_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_ChgLstname_Call) Return(_a0 error) *UserRepository_ChgLstname_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ChgLstname_Call) RunAndReturn(run func(context.Context, string, string) error) *UserRepository_ChgLstname_Call {
	_c.Call.Return(run)
	return _c
}

// ChgName provides a mock function with given fields: ctx, uname, nName
func (_m *UserRepository) ChgName(ctx context.Context, uname string, nName string) error {
	ret := _m.Called(ctx, uname, nName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, uname, nName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_ChgName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChgName'
type UserRepository_ChgName_Call struct {
	*mock.Call
}

// ChgName is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
//   - nName string
func (_e *UserRepository_Expecter) ChgName(ctx interface{}, uname interface{}, nName interface{}) *UserRepository_ChgName_Call {
	return &UserRepository_ChgName_Call{Call: _e.mock.On("ChgName", ctx, uname, nName)}
}

func (_c *UserRepository_ChgName_Call) Run(run func(ctx context.Context, uname string, nName string)) *UserRepository_ChgName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_ChgName_Call) Return(_a0 error) *UserRepository_ChgName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ChgName_Call) RunAndReturn(run func(context.Context, string, string) error) *UserRepository_ChgName_Call {
	_c.Call.Return(run)
	return _c
}

// ChgPasswd provides a mock function with given fields: ctx, uuid, nPasswd
func (_m *UserRepository) ChgPasswd(ctx context.Context, uuid string, nPasswd string) error {
	ret := _m.Called(ctx, uuid, nPasswd)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, uuid, nPasswd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_ChgPasswd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChgPasswd'
type UserRepository_ChgPasswd_Call struct {
	*mock.Call
}

// ChgPasswd is a helper method to define mock.On call
//   - ctx context.Context
//   - uuid string
//   - nPasswd string
func (_e *UserRepository_Expecter) ChgPasswd(ctx interface{}, uuid interface{}, nPasswd interface{}) *UserRepository_ChgPasswd_Call {
	return &UserRepository_ChgPasswd_Call{Call: _e.mock.On("ChgPasswd", ctx, uuid, nPasswd)}
}

func (_c *UserRepository_ChgPasswd_Call) Run(run func(ctx context.Context, uuid string, nPasswd string)) *UserRepository_ChgPasswd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_ChgPasswd_Call) Return(_a0 error) *UserRepository_ChgPasswd_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ChgPasswd_Call) RunAndReturn(run func(context.Context, string, string) error) *UserRepository_ChgPasswd_Call {
	_c.Call.Return(run)
	return _c
}

// ChgRole provides a mock function with given fields: ctx, uname, ro
func (_m *UserRepository) ChgRole(ctx context.Context, uname string, ro domain.Role) error {
	ret := _m.Called(ctx, uname, ro)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Role) error); ok {
		r0 = rf(ctx, uname, ro)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_ChgRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChgRole'
type UserRepository_ChgRole_Call struct {
	*mock.Call
}

// ChgRole is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
//   - ro domain.Role
func (_e *UserRepository_Expecter) ChgRole(ctx interface{}, uname interface{}, ro interface{}) *UserRepository_ChgRole_Call {
	return &UserRepository_ChgRole_Call{Call: _e.mock.On("ChgRole", ctx, uname, ro)}
}

func (_c *UserRepository_ChgRole_Call) Run(run func(ctx context.Context, uname string, ro domain.Role)) *UserRepository_ChgRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.Role))
	})
	return _c
}

func (_c *UserRepository_ChgRole_Call) Return(_a0 error) *UserRepository_ChgRole_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ChgRole_Call) RunAndReturn(run func(context.Context, string, domain.Role) error) *UserRepository_ChgRole_Call {
	_c.Call.Return(run)
	return _c
}

// ChgState provides a mock function with given fields: ctx, uname, st
func (_m *UserRepository) ChgState(ctx context.Context, uname string, st domain.UserState) error {
	ret := _m.Called(ctx, uname, st)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.UserState) error); ok {
		r0 = rf(ctx, uname, st)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_ChgState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChgState'
type UserRepository_ChgState_Call struct {
	*mock.Call
}

// ChgState is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
//   - st domain.UserState
func (_e *UserRepository_Expecter) ChgState(ctx interface{}, uname interface{}, st interface{}) *UserRepository_ChgState_Call {
	return &UserRepository_ChgState_Call{Call: _e.mock.On("ChgState", ctx, uname, st)}
}

func (_c *UserRepository_ChgState_Call) Run(run func(ctx context.Context, uname string, st domain.UserState)) *UserRepository_ChgState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.UserState))
	})
	return _c
}

func (_c *UserRepository_ChgState_Call) Return(_a0 error) *UserRepository_ChgState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ChgState_Call) RunAndReturn(run func(context.Context, string, domain.UserState) error) *UserRepository_ChgState_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, uuid
func (_m *UserRepository) Delete(ctx context.Context, uuid string) error {
	ret := _m.Called(ctx, uuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - uuid string
func (_e *UserRepository_Expecter) Delete(ctx interface{}, uuid interface{}) *UserRepository_Delete_Call {
	return &UserRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, uuid)}
}

func (_c *UserRepository_Delete_Call) Run(run func(ctx context.Context, uuid string)) *UserRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_Delete_Call) Return(_a0 error) *UserRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Delete_Call) RunAndReturn(run func(context.Context, string) error) *UserRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// ExistByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepository) ExistByEmail(ctx context.Context, email string) bool {
	ret := _m.Called(ctx, email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UserRepository_ExistByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistByEmail'
type UserRepository_ExistByEmail_Call struct {
	*mock.Call
}

// ExistByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *UserRepository_Expecter) ExistByEmail(ctx interface{}, email interface{}) *UserRepository_ExistByEmail_Call {
	return &UserRepository_ExistByEmail_Call{Call: _e.mock.On("ExistByEmail", ctx, email)}
}

func (_c *UserRepository_ExistByEmail_Call) Run(run func(ctx context.Context, email string)) *UserRepository_ExistByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_ExistByEmail_Call) Return(_a0 bool) *UserRepository_ExistByEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ExistByEmail_Call) RunAndReturn(run func(context.Context, string) bool) *UserRepository_ExistByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// ExistByUname provides a mock function with given fields: ctx, uname
func (_m *UserRepository) ExistByUname(ctx context.Context, uname string) bool {
	ret := _m.Called(ctx, uname)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, uname)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UserRepository_ExistByUname_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistByUname'
type UserRepository_ExistByUname_Call struct {
	*mock.Call
}

// ExistByUname is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
func (_e *UserRepository_Expecter) ExistByUname(ctx interface{}, uname interface{}) *UserRepository_ExistByUname_Call {
	return &UserRepository_ExistByUname_Call{Call: _e.mock.On("ExistByUname", ctx, uname)}
}

func (_c *UserRepository_ExistByUname_Call) Run(run func(ctx context.Context, uname string)) *UserRepository_ExistByUname_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_ExistByUname_Call) Return(_a0 bool) *UserRepository_ExistByUname_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ExistByUname_Call) RunAndReturn(run func(context.Context, string) bool) *UserRepository_ExistByUname_Call {
	_c.Call.Return(run)
	return _c
}

// ExistByUuid provides a mock function with given fields: ctx, uuid
func (_m *UserRepository) ExistByUuid(ctx context.Context, uuid string) bool {
	ret := _m.Called(ctx, uuid)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UserRepository_ExistByUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistByUuid'
type UserRepository_ExistByUuid_Call struct {
	*mock.Call
}

// ExistByUuid is a helper method to define mock.On call
//   - ctx context.Context
//   - uuid string
func (_e *UserRepository_Expecter) ExistByUuid(ctx interface{}, uuid interface{}) *UserRepository_ExistByUuid_Call {
	return &UserRepository_ExistByUuid_Call{Call: _e.mock.On("ExistByUuid", ctx, uuid)}
}

func (_c *UserRepository_ExistByUuid_Call) Run(run func(ctx context.Context, uuid string)) *UserRepository_ExistByUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_ExistByUuid_Call) Return(_a0 bool) *UserRepository_ExistByUuid_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_ExistByUuid_Call) RunAndReturn(run func(context.Context, string) bool) *UserRepository_ExistByUuid_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	ret := _m.Called(ctx)

	var r0 []domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type UserRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserRepository_Expecter) GetAll(ctx interface{}) *UserRepository_GetAll_Call {
	return &UserRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *UserRepository_GetAll_Call) Run(run func(ctx context.Context)) *UserRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserRepository_GetAll_Call) Return(_a0 []domain.User, _a1 error) *UserRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetAll_Call) RunAndReturn(run func(context.Context) ([]domain.User, error)) *UserRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUsername provides a mock function with given fields: ctx, uname
func (_m *UserRepository) GetByUsername(ctx context.Context, uname string) (domain.User, error) {
	ret := _m.Called(ctx, uname)

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, error)); ok {
		return rf(ctx, uname)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, uname)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uname)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUsername'
type UserRepository_GetByUsername_Call struct {
	*mock.Call
}

// GetByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
func (_e *UserRepository_Expecter) GetByUsername(ctx interface{}, uname interface{}) *UserRepository_GetByUsername_Call {
	return &UserRepository_GetByUsername_Call{Call: _e.mock.On("GetByUsername", ctx, uname)}
}

func (_c *UserRepository_GetByUsername_Call) Run(run func(ctx context.Context, uname string)) *UserRepository_GetByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_GetByUsername_Call) Return(_a0 domain.User, _a1 error) *UserRepository_GetByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByUsername_Call) RunAndReturn(run func(context.Context, string) (domain.User, error)) *UserRepository_GetByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUuid provides a mock function with given fields: ctx, uuid
func (_m *UserRepository) GetByUuid(ctx context.Context, uuid string) (domain.User, error) {
	ret := _m.Called(ctx, uuid)

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, error)); ok {
		return rf(ctx, uuid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_GetByUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUuid'
type UserRepository_GetByUuid_Call struct {
	*mock.Call
}

// GetByUuid is a helper method to define mock.On call
//   - ctx context.Context
//   - uuid string
func (_e *UserRepository_Expecter) GetByUuid(ctx interface{}, uuid interface{}) *UserRepository_GetByUuid_Call {
	return &UserRepository_GetByUuid_Call{Call: _e.mock.On("GetByUuid", ctx, uuid)}
}

func (_c *UserRepository_GetByUuid_Call) Run(run func(ctx context.Context, uuid string)) *UserRepository_GetByUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_GetByUuid_Call) Return(_a0 domain.User, _a1 error) *UserRepository_GetByUuid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_GetByUuid_Call) RunAndReturn(run func(context.Context, string) (domain.User, error)) *UserRepository_GetByUuid_Call {
	_c.Call.Return(run)
	return _c
}

// Login provides a mock function with given fields: ctx, uname, passwd
func (_m *UserRepository) Login(ctx context.Context, uname string, passwd string) (domain.User, error) {
	ret := _m.Called(ctx, uname, passwd)

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (domain.User, error)); ok {
		return rf(ctx, uname, passwd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.User); ok {
		r0 = rf(ctx, uname, passwd)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, uname, passwd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type UserRepository_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - ctx context.Context
//   - uname string
//   - passwd string
func (_e *UserRepository_Expecter) Login(ctx interface{}, uname interface{}, passwd interface{}) *UserRepository_Login_Call {
	return &UserRepository_Login_Call{Call: _e.mock.On("Login", ctx, uname, passwd)}
}

func (_c *UserRepository_Login_Call) Run(run func(ctx context.Context, uname string, passwd string)) *UserRepository_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *UserRepository_Login_Call) Return(_a0 domain.User, _a1 error) *UserRepository_Login_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Login_Call) RunAndReturn(run func(context.Context, string, string) (domain.User, error)) *UserRepository_Login_Call {
	_c.Call.Return(run)
	return _c
}

// Store provides a mock function with given fields: ctx, u
func (_m *UserRepository) Store(ctx context.Context, u *domain.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRepository_Store_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Store'
type UserRepository_Store_Call struct {
	*mock.Call
}

// Store is a helper method to define mock.On call
//   - ctx context.Context
//   - u *domain.User
func (_e *UserRepository_Expecter) Store(ctx interface{}, u interface{}) *UserRepository_Store_Call {
	return &UserRepository_Store_Call{Call: _e.mock.On("Store", ctx, u)}
}

func (_c *UserRepository_Store_Call) Run(run func(ctx context.Context, u *domain.User)) *UserRepository_Store_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.User))
	})
	return _c
}

func (_c *UserRepository_Store_Call) Return(_a0 error) *UserRepository_Store_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserRepository_Store_Call) RunAndReturn(run func(context.Context, *domain.User) error) *UserRepository_Store_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
