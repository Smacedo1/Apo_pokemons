package repository

import (
	"context"
	"database/sql"
	"testing"

	"gos/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetPokemons(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "nombre", "tipo", "nivel"}).
		AddRow(1, "Pikachu", 1, 10).
		AddRow(2, "Charizard", 2, 20)

	mock.ExpectQuery(`SELECT id, nombre, tipo, nivel FROM pokemones ORDER BY id`).
		WillReturnRows(rows)

	pokemons, err := repo.GetPokemons(context.Background())

	assert.NoError(t, err)
	assert.Len(t, pokemons, 2)
	assert.Equal(t, 1, pokemons[0].ID)
	assert.Equal(t, "Pikachu", pokemons[0].Name)
	assert.Equal(t, 1, pokemons[0].Type1_id)
	assert.Equal(t, 10, *pokemons[0].Type2_id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetPokemon(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	row := sqlmock.NewRows([]string{"id", "name", "type1_id", "type2_id"}).
		AddRow(1, "Pikachu", 1, sql.NullInt64{Int64: 2, Valid: true})

	mock.ExpectQuery(`SELECT id, name, type1_id, type2_id FROM pokemon WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(row)

	pokemon, err := repo.GetPokemon(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, pokemon.ID)
	assert.Equal(t, "Pikachu", pokemon.Name)
	assert.Equal(t, 1, pokemon.Type1_id)
	assert.Equal(t, 2, *pokemon.Type2_id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetType(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Electric").
		AddRow(2, "Fire")

	mock.ExpectQuery(`SELECT id, name FROM type ORDER BY id`).
		WillReturnRows(rows)

	types, err := repo.GetType(context.Background())

	assert.NoError(t, err)
	assert.Len(t, types, 2)
	assert.Equal(t, 1, types[0].ID)
	assert.Equal(t, "Electric", types[0].Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreatePokemon(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	poke := &domain.Pokemon{Name: "Pikachu", Type1_id: 1, Type2_id: &[]int{2}[0]}

	mock.ExpectExec(`INSERT INTO pokemones \(nombre, tipo, nivel\) VALUES \(\?, \?, \?\)`).
		WithArgs(poke.Name, poke.Type1_id, *poke.Type2_id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	created, err := repo.CreatePokemon(context.Background(), poke)

	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePokemon(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec(`DELETE FROM pokemones WHERE id = \?`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeletePokemon(context.Background(), 1)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_PatchPokemon(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	poke := &domain.Pokemon{ID: 1, Name: "Pikachu", Type1_id: 1, Type2_id: &[]int{2}[0]}

	mock.ExpectExec(`UPDATE pokemon SET name = \?, type1_id = \?, type2_id = \? WHERE id = \?`).
		WithArgs(poke.Name, poke.Type1_id, *poke.Type2_id, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.PatchPokemon(context.Background(), 1, poke)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetTypeById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	row := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Electric")

	mock.ExpectQuery(`SELECT id, name FROM type WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(row)

	typ, err := repo.GetTypeById(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, typ.ID)
	assert.Equal(t, "Electric", typ.Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}
