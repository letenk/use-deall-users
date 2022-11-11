package tests

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/letenk/use_deal_user/helper"
	"github.com/letenk/use_deal_user/models/domain"
	"github.com/letenk/use_deal_user/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createRandomUserRepository(t *testing.T) string {
	// Create new sample object user
	user := domain.User{}
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	user.Fullname = fullname
	user.Username = strings.ToLower(helper.RandomPerson())
	user.Role = "admin"

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	helper.ErrLogPanic(err)

	user.Password = string(passwordHash)

	repository := repository.NewUserRepository(ConnTest)

	// Insert
	id, err := repository.Insert(user)
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	return id
}

func TestInsertUserRepository(t *testing.T) {
	t.Parallel()
	createRandomUserRepository(t)
}

func TestGetAllUserRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomUserRepository(t)
			mutex.Unlock()
		}()
	}

	repository := repository.NewUserRepository(ConnTest)
	// Get all
	users, err := repository.GetAll()
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(users))

	for _, data := range users {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Fullname)
		require.NotEmpty(t, data.Username)
		require.NotEmpty(t, data.Password)
		require.NotEmpty(t, data.Role)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
	}
}

func TestGetOneUserRepository(t *testing.T) {
	t.Parallel()
	id := createRandomUserRepository(t)

	repository := repository.NewUserRepository(ConnTest)

	// Get one by id
	user, err := repository.GetOne(id)
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.Equal(t, id, user.ID)
	assert.NotEmpty(t, user.Fullname)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Password)
	assert.NotEmpty(t, user.Role)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
}

func TestGetByUsernameUserRepository(t *testing.T) {
	t.Parallel()
	id := createRandomUserRepository(t)

	repository := repository.NewUserRepository(ConnTest)

	// GetOne by id
	user, err := repository.GetOne(id)
	helper.ErrLogPanic(err)

	// Get one by id
	userByEmail, err := repository.GetOneByUsername(user.Username)
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.Equal(t, id, userByEmail.ID)
	assert.NotEmpty(t, userByEmail.Fullname)
	assert.NotEmpty(t, userByEmail.Username)
	assert.NotEmpty(t, userByEmail.Password)
	assert.NotEmpty(t, userByEmail.Role)
	assert.NotEmpty(t, userByEmail.CreatedAt)
	assert.NotEmpty(t, userByEmail.UpdatedAt)
}

func TestUpdateUserRepository(t *testing.T) {
	t.Parallel()
	id := createRandomUserRepository(t)

	// New data update
	data := domain.User{}
	data.ID = id
	data.Fullname = "updated"

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("newPassword"), bcrypt.MinCost)
	helper.ErrLogPanic(err)
	data.Password = string(passwordHash)

	repository := repository.NewUserRepository(ConnTest)
	// Update
	ok, err := repository.Update(data)
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteUserRepository(t *testing.T) {
	t.Parallel()
	id := createRandomUserRepository(t)

	repository := repository.NewUserRepository(ConnTest)
	// Update
	ok, err := repository.Delete(id)
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.True(t, ok)
}
