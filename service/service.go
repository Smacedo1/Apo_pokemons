package service

import (
	"context"
	"gos/domain"
	"gos/repository"
)

// Service es la interfaz que expone el servicio de personas
type Service interface {
	Get(ctx context.Context) ([]domain.Persona, error)
	GetPokemons(ctx context.Context) ([]domain.Pokemon, error)
	Post(ctx context.Context, poke *domain.Pokemon) (*domain.Pokemon, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, poke *domain.Pokemon) error
	GetType(ctx context.Context) ([]string, error)
}

// Servicio es la implementación concreta del servicio
type Servicio struct {
	repo repository.Repo
}

func NewService(r repository.Repo) *Servicio {
	return &Servicio{
		repo: r,
	}
}

func (s *Servicio) GetType(ctx context.Context) ([]string, error) {
	return s.repo.GetType(ctx)
}

func (s *Servicio) Get(ctx context.Context) ([]domain.Persona, error) {
	return s.repo.GetAll(ctx)
}

func (s *Servicio) GetPokemons(ctx context.Context) ([]domain.Pokemon, error) {
	return s.repo.GetPokemons(ctx)
}

func (s *Servicio) Post(ctx context.Context, poke *domain.Pokemon) (*domain.Pokemon, error) {
	return s.repo.CreatePokemon(ctx, poke)
}

func (s *Servicio) Delete(ctx context.Context, id int) error {
	return s.repo.DeletePokemon(ctx, id)
}

func (s *Servicio) Patch(ctx context.Context, id int, poke *domain.Pokemon) error {
	return s.repo.PatchPokemon(ctx, id, poke)
}
