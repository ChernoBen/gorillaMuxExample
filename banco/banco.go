package banco

import (
	"database/sql"

	//biblioteca/driver do mysql deve ser importada de maneira implicita
	_ "github.com/go-sql-driver/mysql"
)

//funcao para abrir conn
func Conectar() (*sql.DB, error) {
	//instanciar string de conexão
	stringConexao := "sammy:password@/wc?charset=utf8&parseTimeTrue&loc=Local"
	//conectar
	db, err := sql.Open("mysql", stringConexao)
	if err != nil {
		return nil, err
	}
	//verifica se a conexão foi bem sucedida
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
