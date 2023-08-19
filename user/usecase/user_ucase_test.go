package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/mocks"
	"github.com/sicozz/papyrus/test_utils/instances"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userUCase "github.com/sicozz/papyrus/user/usecase"
)

func TestGetAll(t *testing.T) {
	mUserRepo := new(mocks.UserRepository)
	mRoleRepo := new(mocks.RoleRepository)
	mUStateRepo := new(mocks.UserStateRepository)

	t.Run("success", func(t *testing.T) {
		mUsers := instances.GetUserList()
		mRoles := instances.GetRoleList()
		mUStates := instances.GetUStateList()
		mUserRepo.On("GetAll", mock.Anything).Return(mUsers).Once()
		mRoleRepo.On("GetAll", mock.Anything).Return(mRoles).Once()
		mUStateRepo.On("GetAll", mock.Anything).Return(mUStates).Once()

		u := userUCase.NewUserUsecase(mUserRepo, mRoleRepo, mUStateRepo, time.Second*2)
		users, err := u.GetAll(context.TODO())

		mapRoles := map[int64]domain.Role{}
		for _, role := range mRoles { //nolint
			mapRoles[role.Code] = role
		}

		mapStates := map[int64]domain.UserState{}
		for _, state := range mUStates { //nolint
			mapStates[state.Code] = state
		}

		for idx, mU := range mUsers { //nolint
			if r, ok := mapRoles[mU.Role.Code]; ok {
				mUsers[idx].Role = r
			}

			if s, ok := mapStates[mU.State.Code]; ok {
				mUsers[idx].State = s
			}
		}

		if !assert.NoError(t, err) {
			t.Error("GetAll returned an error -> ", err)
		}
		if !assert.Len(t, users, 3) {
			t.Error("GetAll returned a different number of users -> ", len(users))
		}
		assert.ElementsMatch(t, mUsers, users)
	})
}
