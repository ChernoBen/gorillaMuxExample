package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// CREATE, READ, UPDATE, DELETE
	// utilizando pacote mux

	// Primeiro passo criar um router
	router := mux.NewRouter()
	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", servidor.BuscarUsuarios).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", servidor.BuscaUsuario).Methods("GET")
	router.HandleFunc("/usuarios/{id}", servidor.AtualizarUsuario).Methods("PUT")
	// utilizar pacote hhtp para subir o servidor junto com o router

	fmt.Println("Escutando porta 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
