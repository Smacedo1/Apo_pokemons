package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gos/domain"
)

/*todo tipo que implemente esta interfaz debe tener estas funciones con estas firmas, contrato que dice qué funciones deben existi
Esto significa que cualquier struct que tenga un método GetAll(ctx context.Context) ya es un Repo*/

type Repo interface {
	//GetAll(ctx context.Context) ([]domain.Persona, error) crea la query hace la sentencia sql //
	GetPokemons(ctx context.Context) ([]domain.Pokemon, error)
	GetPokemon(ctx context.Context, id int) (*domain.Pokemon, error)
	GetType(ctx context.Context) ([]domain.PokemonType, error)
	CreatePokemon(ctx context.Context, poke *domain.Pokemon) (*domain.Pokemon, error)
	DeletePokemon(ctx context.Context, id int) error
	PatchPokemon(ctx context.Context, id int, poke *domain.Pokemon) error
	GetTypeById(ctx context.Context, id int) (domain.PokemonType, error)
}

type Repository struct { //Implementacion real del repositorio

	Db *sql.DB
}

func NewRepository(db *sql.DB) *Repository { //Un constructor para poder inyectar dependencias
	return &Repository{
		Db: db,
	}
}

func (r *Repository) GetPokemon(ctx context.Context, id int) (*domain.Pokemon, error) {
	const q = `
		SELECT id, name, type1_id, type2_id 
		FROM pokemon 
		WHERE id = ?;`

	row := r.Db.QueryRowContext(ctx, q, id)

	var p domain.Pokemon

	if err := row.Scan(&p.ID, &p.Name, &p.Type1_id, &p.Type2_id); err != nil {
		return nil, fmt.Errorf("repo pokemon: scan: %w", err)
	}

	return &p, nil
}

func (r *Repository) GetTypeById(ctx context.Context, id int) (domain.PokemonType, error) {
	const q = `
		SELECT *
		FROM type
		WHERE id = ?;
		
	`
	row := r.Db.QueryRowContext(ctx, q, id)

	var t domain.PokemonType
	if err := row.Scan(&t); err != nil {
		return domain.PokemonType{}, fmt.Errorf("repo type: scan: %w", err)
	}
	return t, nil
}

func (r *Repository) GetType(ctx context.Context) ([]domain.PokemonType, error) {
	const q = `
		SELECT name
		FROM type
		ORDER BY id;
	`
	rows, err := r.Db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repo type: query GetType: %w", err)
	}
	defer rows.Close()

	var types []domain.PokemonType
	for rows.Next() {
		var t domain.PokemonType
		if err := rows.Scan(&t, "repo type: scan: %w", err); err != nil {
			return nil, fmt.Errorf("repo type: rows err: %w", err)
		}
		types = append(types, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo type: rows err: %w", err)
	}
	/*if len(types) == 0 {
		return nil, errors.New("no hay tipos")
	}*/
	return types, nil
}

func (r *Repository) PatchPokemon(ctx context.Context, id int, poke *domain.Pokemon) error {
	const updateQ = `
	UPDATE pokemon 
	SET tipo = ?, Name = ?, Type1_id = ? Type2_id = ?
	WHERE id = ?;
	`
	_, err := r.Db.ExecContext(ctx, updateQ, poke.ID, poke.Name, poke.Type1_id, id)
	if err != nil {
		return fmt.Errorf("repo pokemones: update: %w", err)
	}
	return nil
}

func (r *Repository) DeletePokemon(ctx context.Context, id int) error {
	const deleteQ = `
	 DELETE FROM pokemones
	 WHERE id = ?;
	 `
	_, err := r.Db.ExecContext(ctx, deleteQ, id)
	if err != nil {
		return fmt.Errorf("repo pokemones: delete: %w", err)
	}
	return nil
}

func (r *Repository) CreatePokemon(ctx context.Context, poke *domain.Pokemon) (*domain.Pokemon, error) {
	// TODO: Implementar la inserción del usuario
	const insertQ = `
		INSERT INTO pokemones (tipo, nombre, nivel)
		VALUES (?, ?, ?);
	`
	result, err := r.Db.ExecContext(ctx, insertQ, poke.ID, poke.Name, poke.Type1_id, poke.Type2_id)
	if err != nil {
		return nil, fmt.Errorf("repo pokemones: insert: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("repo pokemones: getting last insert id: %w", err)
	}
	pokemon := &domain.Pokemon{
		ID: int(id),
	}
	return pokemon, nil

}

/*func (r *Repository) GetAll(ctx context.Context) ([]domain.Persona, error) { //metodo que va a usar el servicio
	// Ajustá el nombre de la tabla/columnas a tu esquema real: personas o usuarios
	const q = `
		SELECT dni, nombre, apellido
		FROM usuarios
		-- descomentá si usás soft-delete:
		-- WHERE deleted_at IS NULL
		ORDER BY apellido, nombre;
	`

	rows, err := r.Db.QueryContext(ctx, q) //llamamos
	if err != nil {
		return nil, fmt.Errorf("repo personas: query GetAll: %w", err)
	}
	defer rows.Close()

	var out []domain.Persona //creamos listado de personas y con for recorremos linea por linea
	for rows.Next() {
		var p domain.Persona                                              //creamos persona vacia
		if err := rows.Scan(&p.Dni, &p.Nombre, &p.Apellido); err != nil { //el Scan transforma todo los datos que tenga la linea que trajo la base de dato lo parsea al objeto, y lo deja en la direccion del objeto &
			return nil, fmt.Errorf("repo personas: scan: %w", err) //en caso de que haya error
		}
		out = append(out, p) //le pasamos lo que creamos a la nueva lista, apendeamos
	}
	if err := rows.Err(); err != nil { //pregunta si alguna linea tuvo error o esta corrupta
		return nil, fmt.Errorf("repo personas: rows err: %w", err) //hay que controlarlo por las dudas
	}
	if len(out) == 0 {
		return nil, errors.New("no hay personas") //error que la base no tenga personas
	}
	return out, nil //en caso de que todo este bien devolvemos la lista y nil
}*/

// GetPokemons devuelve la lista de pokemones desde la base de datos
func (r *Repository) GetPokemons(ctx context.Context) ([]domain.Pokemon, error) {
	const q = `
		SELECT id, nombre, tipo, nivel
		FROM pokemones
		ORDER BY id;
	`

	rows, err := r.Db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repo pokemones: query GetPokemons: %w", err)
	}
	defer rows.Close()

	var out []domain.Pokemon
	for rows.Next() {
		var p domain.Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Type1_id, &p.Type2_id); err != nil {
			return nil, fmt.Errorf("repo pokemones: scan: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo pokemones: rows err: %w", err)
	}
	if len(out) == 0 {
		return nil, errors.New("no hay pokemones")
	}
	return out, nil
}
