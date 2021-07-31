package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// requisição com letra minuscula para nao ser utilizado fora do pacote
type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

//user controller
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	//ioutil para leitura e escrita na requisição
	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//escrevendo response
		w.Write([]byte("Falha ao ler copo da requisição"))
		return
	}

	//inserir requisição na struct local usuario
	var usuario usuario
	if err = json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		w.Write([]byte("Erro ao converter usuario p/struct"))
		return
	}

	//apos o processamento da requisição iniciar inserção no banco
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha na conexão com banco"))
		return
	}
	defer db.Close()

	// prepare statement / criar comando de inserção para evitar sql injection
	statement, err := db.Prepare("INSERT INTO usuarios (nome, email) VALUES(?, ?)")
	if err != nil {
		w.Write([]byte("Erro ao criar statement"))
		return
	}
	//assim como o banco o statement deve ser fechado
	defer statement.Close()
	//statement construir agora executar inserção
	insercao, err := statement.Exec(usuario.Nome, usuario.Email)
	if err != nil {
		w.Write([]byte("Erro ao executar statement"))
		return
	}
	//se passar desse ponto entao o usuario ja foi inserido
	//pegar id inserido e devolver na response
	idInserido, err := insercao.LastInsertId()
	if err != nil {
		w.Write([]byte("erro ao obter id inserido"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso id: %d", idInserido)))
}

//obter toda a listagem de usuarios
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	//abrir conn com banco
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("falha ao conectar com banco de dados"))
		return
	}
	//fechado conn
	defer db.Close()

	//criando select
	linhas, err := db.Query("SELECT * FROM usuarios")
	if err != nil {
		w.Write([]byte("Falha ao buscar usuarios"))
		return
	}
	//fechando conn
	defer linhas.Close()

	// criando slice de usuarios/tipo um array de objetos
	var usuarios []usuario
	//para cada linha retornada executar uma iteração
	for linhas.Next() {
		var usuario usuario
		//verifica a ordem de cada info da linha com o tipo de dado da struct
		if err := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("falha ao escanear usuario"))
			return
		}
		//inserir dados da tabela usuarios no slice de usuarios
		usuarios = append(usuarios, usuario)
	}
	w.WriteHeader(http.StatusOK)
	//tranformar slice de struct em json
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Erro ao converter resposta em json"))
		return
	}
}

//buscar usuario especifico
func BuscaUsuario(w http.ResponseWriter, r *http.Request) {
	//ler parametro vindo da rota
	parametros := mux.Vars(r) //passar a r(request como parametro em vars())
	//obter paremetro id da rota e converter para inteiro
	ID, err := strconv.ParseInt(parametros["id"], 10, 32) // param["id"] <- string a ser convertida / 10 <- base / 32 <-bits id será um tipo int32
	if err != nil {
		w.Write([]byte("Erro ao converter id em inteiro, verifique o id passado!"))
	}
	// abrir conn
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("falha ao conectar com banco de dados"))
		return
	}
	//fechado conn
	defer db.Close()

	//criando select
	linha, err := db.Query("SELECT * FROM usuarios WHERE id = ?", ID)
	if err != nil {
		w.Write([]byte("Falha ao buscar usuario!"))
		return
	}
	//fechando conn
	defer linha.Close()

	var usuario usuario
	if linha.Next() {
		//se der erro em algum dos campor id,nome ou email faça
		if err := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Erro ao processar/escanear dados do usuário"))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		w.Write([]byte("Erro ao converter dados do usuário"))
		return
	}

}
