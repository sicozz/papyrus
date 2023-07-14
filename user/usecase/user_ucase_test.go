package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/sicozz/papyrus/domain"
	"github.com/sicozz/papyrus/domain/mocks"
	ucase "github.com/sicozz/papyrus/user/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		Uuid:     "1",
		Username: "tUserName",
		Email:    "tEmail",
		Password: "tPasswd",
		Name:     "tName",
		Lastname: "tLastname",
		Role: domain.Role{
			Code:        0,
			Description: "user",
		},
		State: domain.UserState{
			Code:        0,
			Description: "user",
		},
	}

	mockListUser := make([]domain.User, 0)
	mockListUser = append(mockListUser, mockUser)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("Fetch", mock.Anything).
			Return(mockListUser, nil).
			Once()

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)
		list, err := u.Fetch(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, list, len(mockListUser))
		assert.Equal(t, mockListUser[0].Uuid, list[0].Uuid)

		mockUserRepo.AssertExpectations(t)
	})
}
