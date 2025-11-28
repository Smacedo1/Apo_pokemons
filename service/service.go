package service

import (
	"context"
	"gos/domain"
	"gos/repository"
	"gos/response"
)

// Service es la interfaz que expone el servicio de personas
type Service interface {
	//Get(ctx context.Context) ([]domain.Persona, error)//
	GetPokemons(ctx context.Context) ([]domain.Pokemon, error)
	GetPokemon(ctx context.Context, id int) (*response.PokemonResponse, error)
	Post(ctx context.Context, poke *domain.Pokemon) (*domain.Pokemon, error)
	Delete(ctx context.Context, id int) error
	Patch(ctx context.Context, id int, poke *domain.Pokemon) error
	GetType(ctx context.Context) ([]domain.PokemonType, error)
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

func (s *Servicio) GetPokemon(ctx context.Context, id int) (*response.PokemonResponse, error) {
	responsewithtype := response.PokemonResponse{}

	pokemon, err := s.repo.GetPokemon(ctx, id)
	if err != nil {
		return nil, err
	}

	type1, err := s.repo.GetTypeById(ctx, pokemon.Type1_id)
	if err != nil {
		return nil, err
	}

	if pokemon.Type2_id != nil {
		type2, err := s.repo.GetTypeById(ctx, *pokemon.Type2_id)
		if err != nil {
			return nil, err
		}
		responsewithtype.Type2_id = type2

	}
	responsewithtype.ID = pokemon.ID
	responsewithtype.Name = pokemon.Name
	responsewithtype.Type1_id = type1

	return &responsewithtype, nil

}

func (s *Servicio) GetType(ctx context.Context) ([]domain.PokemonType, error) {
	return s.repo.GetType(ctx)
}

/*func (s *Servicio) Get(ctx context.Context) ([]domain.Persona, error) {
	return s.repo.GetAll(ctx)
}*/

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
