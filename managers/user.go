package managers

import (
	"log"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	"gorm.io/gorm"

	"github.com/slavik-o/go-web-app/models"
)

type UserManager struct {
	DB *gorm.DB
}

func (m *UserManager) UpsertAuthUser(authUser goth.User) (*models.User, error) {
	user := new(models.User)

	result := m.DB.Where("provider = ? AND provider_id = ?", authUser.Provider, authUser.UserID).Find(user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected < 1 {
		uuid, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("Failed to generate UUID: %v", err)
		}

		user.ID = uuid.String()
	}

	user.Provider = authUser.Provider
	user.ProviderID = authUser.UserID
	user.Email = authUser.Email
	user.FirstName = authUser.FirstName
	user.LastName = authUser.LastName

	result = m.DB.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (m *UserManager) FindUserByID(ID string) (*models.User, error) {
	user := new(models.User)

	result := m.DB.First(user, "id = ?", ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
