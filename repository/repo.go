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
	GetAll(ctx context.Context) ([]domain.Persona, error) //crea la query hace la sentencia sql
	GetPokemons(ctx context.Context) ([]domain.Pokemon, error)
	GetType(ctx context.Context) ([]string, error)
	CreatePokemon(ctx context.Context, poke *domain.Pokemon) (*domain.Pokemon, error)
	DeletePokemon(ctx context.Context, id int) error
	PatchPokemon(ctx context.Context, id int, poke *domain.Pokemon) error
}

type Repository struct { //Implementacion real del repositorio

	Db *sql.DB
}

func NewRepository(db *sql.DB) *Repository { //Un constructor para poder inyectar dependencias
	return &Repository{
		Db: db,
	}
}

func (r *Repository) GetType(ctx context.Context) ([]string, error) {
	const qy = `
		SELECT name
		FROM type
		ORDER BY id;
	`
	rows, err := r.Db.QueryContext(ctx, qy)
	if err != nil { 
		return nil, fmt.Errorf("repo pokemones: query GetType: %w", err)
	} 
	defer rows.Close() //cierra la consulta a la bd y libera los recursos el defer pospone esto
	var types []string
	for rows.Next() {
		var t string 
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("repo pokemones: scan: %w", err)
		}
		types = append(types, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo pokemones: rows err: %w", err)
	}
	/*if len(types) == 0 {
		return nil, errors.New("no hay tipos")
	}*/
	return types, nil
}



func (r *Repository) PatchPokemon(ctx context.Context, id int, poke *domain.Pokemon) error {
	const updateQ = `
	UPDATE pokemones 
	SET tipo = ?, nombre = ?, nivel = ?
	WHERE id = ?;
	`
	_, err := r.Db.ExecContext(ctx, updateQ, poke.Tipo, poke.Nombre, poke.Nivel, id)
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
	result, err := r.Db.ExecContext(ctx, insertQ, poke.Tipo, poke.Nombre, poke.Nivel)
	if err != nil {
		return nil, fmt.Errorf("repo pokemones: insert: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("repo pokemones: getting last insert id: %w", err)
	}
	pokemon := &domain.Pokemon{
		ID:     int(id),
		Tipo:   poke.Tipo,
		Nombre: poke.Nombre,
		Nivel:  poke.Nivel,
	}
	return pokemon, nil

}

func (r *Repository) GetAll(ctx context.Context) ([]domain.Persona, error) { //metodo que va a usar el servicio
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
}

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
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Tipo, &p.Nivel); err != nil {
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
