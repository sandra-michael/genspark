package stores

import "app/stores/models"

type DataBase interface {
	Create(u models.User) (*models.User, bool)
	//CreateSimple(string) (models.User, bool)
	Update(int, string) (*models.User, bool)
	Delete(int) (*models.User, bool)
	FetchAll() (map[int]*models.User, bool)
	FetchUser(int) (*models.User, bool)
}
