package response

import ("gos/domain")

type PokemonResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type1_id domain.PokemonType   `json:"type1"`
	Type2_id domain.PokemonType    `json:"type2"`
}