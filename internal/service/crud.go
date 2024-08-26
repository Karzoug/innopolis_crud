package service

import (
	"crud/internal/domain"
	"crud/internal/repository/recipedb"

	"github.com/google/uuid"
)

var recipes recipedb.DB

func Init(DB recipedb.DB) {
	recipes = DB
}

func Get(id string) (*domain.Recipe, error) {
	return recipes.Get(id)
}

func GetAuthorID(id string) (string, error) {
	r, err := recipes.Get(id)
	if err != nil {
		return "", err
	}
	return r.AuthorID, nil
}

func Delete(id string) error {
	return recipes.Delete(id)
}

func AddOrUpd(r *domain.Recipe) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}

	return recipes.Set(r.ID, r)
}
