package user

import (
	"bookmark-management/internal/domain"
	"bookmark-management/pkg/sqldb"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var dbError = errors.New("test error")

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		setupDb     func(t *testing.T) *gorm.DB
		inputModel  *domain.User
		expectedErr error
		verifyFunc  func(db *gorm.DB, username string)
	}{
		{
			name: "success",
			setupDb: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDatabase(t)
				db.AutoMigrate(&domain.User{})
				return db
			},
			inputModel: &domain.User{
				ID:          "684f979e-672c-4f7a-8199-4274c0f8d83a",
				Username:    "testuser",
				Password:    "testpassword",
				DisplayName: "testuser",
			},
			expectedErr: nil,
			verifyFunc: func(db *gorm.DB, username string) {
				user := &domain.User{}
				result := db.Where("username = ?", username).First(user)
				assert.Nil(t, result.Error)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name: "db error",
			setupDb: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDatabase(t)
				db.AutoMigrate(&domain.User{})
				db.Error = dbError
				return db
			},
			inputModel: &domain.User{
				ID:          "684f979e-672c-4f7a-8199-4274c0f8d83a",
				Username:    "testuser",
				Password:    "testpassword",
				DisplayName: "testuser",
			},
			expectedErr: dbError,
		},
		{
			name: "username already exists",
			setupDb: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDatabase(t)
				db.AutoMigrate(&domain.User{})
				db.Create(&domain.User{
					ID:       "00000001-0001-0001-0001-000000000012",
					Username: "testuser",
					Password: "testpassword",
					Email:    "a1@gmail.com",
				})
				return db
			},
			inputModel: &domain.User{
				ID:       "684f979e-672c-4f7a-8199-4274c0f8d83a",
				Username: "testuser",
				Password: "testpassword",
				Email:    "a2@gmail.com",
			},
			expectedErr: domain.ErrAlreadyExists,
		},
		{
			name: "email already exists",
			setupDb: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDatabase(t)
				db.AutoMigrate(&domain.User{})
				db.Create(&domain.User{
					ID:       "00000001-0001-0001-0001-000000000012",
					Username: "a1",
					Password: "testpassword",
					Email:    "test@gmail.com",
				})
				return db
			},
			inputModel: &domain.User{
				ID:       "684f979e-672c-4f7a-8199-4274c0f8d83a",
				Username: "a1",
				Password: "testpassword",
				Email:    "test@gmail.com",
			},
			expectedErr: domain.ErrAlreadyExists,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDb(t)
			repo := NewRepository(db)
			user, err := repo.CreateUser(t.Context(), tc.inputModel)

			switch {
			case tc.expectedErr != nil:
				assert.ErrorIs(t, err, tc.expectedErr)
			default:
				assert.NoError(t, err)
				if tc.verifyFunc != nil {
					tc.verifyFunc(db, user.Username)
				}
			}
		})
	}
}
