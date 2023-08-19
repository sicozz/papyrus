package postgres_test

import (
	"context"
	"testing"

	"github.com/sicozz/papyrus/test_utils/instances"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	userPostgresRepo "github.com/sicozz/papyrus/user/repository/postgres"
)

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	// Mock users
	mUsers := instances.GetUserList()
	rows := sqlmock.NewRows([]string{"uuid", "username", "email", "password", "name", "lastname", "role", "state"}).
		AddRow(mUsers[0].Uuid,
			mUsers[0].Username,
			mUsers[0].Email,
			mUsers[0].Password,
			mUsers[0].Name,
			mUsers[0].Lastname,
			mUsers[0].Role.Code,
			mUsers[0].State.Code,
		).
		AddRow(mUsers[1].Uuid,
			mUsers[1].Username,
			mUsers[1].Email,
			mUsers[1].Password,
			mUsers[1].Name,
			mUsers[1].Lastname,
			mUsers[1].Role.Code,
			mUsers[1].State.Code,
		).
		AddRow(mUsers[2].Uuid,
			mUsers[2].Username,
			mUsers[2].Email,
			mUsers[2].Password,
			mUsers[2].Name,
			mUsers[2].Lastname,
			mUsers[2].Role.Code,
			mUsers[2].State.Code,
		)

	query :=
		`SELECT uuid, username, email, password, name, lastname, role, state
		FROM user_`
	mock.ExpectQuery(query).WillReturnRows(rows)

	r := userPostgresRepo.NewPostgresUserRepository(db)
	users, err := r.GetAll(context.TODO())
	if !assert.NoError(t, err) {
		t.Error("GetAll returned an error -> ", err)
	}
	if !assert.Len(t, users, 3) {
		t.Error("GetAll returned a different number of users -> ", len(users))
	}
	assert.ElementsMatch(t, mUsers, users)
}
