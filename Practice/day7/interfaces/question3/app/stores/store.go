package stores

import "app/stores/models"

type DataBase interface {
	Create(u models.User) error
	CreateSimple(string) error
	Update(string) error
	Delete(int) error
}
