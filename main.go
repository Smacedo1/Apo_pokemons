package main

import (
	"context"
	"database/sql"
	"gos/handler"
	"gos/middleware"
	"gos/repository"
	"gos/service"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env, usando variables del sistema")
	}

	// Configurar conexión a la base de datos
	dsn := "root:1304@tcp(127.0.0.1:3306)/mi_primera_base"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer db.Close()

	// Verificar la conexión
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("Error al conectarse a la base de datos: %v", err)
	}

	repo := repository.NewRepository(db)
	srv := service.NewService(repo)
	h := handler.NewHandler(srv)

	// Crear router con Gorilla Mux
	r := mux.NewRouter()

	// Aplicar middleware de autenticación a todas las rutas excepto /login
	protected := r.PathPrefix("/pokedex").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/pokemons", h.GetPokemons).Methods("GET")
	protected.HandleFunc("/pokemons", h.CreatePokemon).Methods("POST")
	protected.HandleFunc("/pokemons/{id}", h.GetPokemon).Methods("GET")
	protected.HandleFunc("/pokemons/{id}", h.UpdatePokemon).Methods("PATCH")
	protected.HandleFunc("/pokemons/{id}", h.DeletePokemon).Methods("DELETE")
	protected.HandleFunc("/types", h.GetTypes).Methods("GET")

	// /login sin auth
	r.HandleFunc("/login", h.Login).Methods("POST")

	log.Printf("Servidor escuchando en :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
