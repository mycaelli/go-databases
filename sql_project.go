package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("zzZZZzz", err)
	}
	
	stmt, err := db.Prepare("CREATE DATABASE IF NOT EXISTS hello")

	if err != nil {
		log.Fatal("Erro na preparação da query", err)
	}
	
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("Erro ao criar o banco de dados: ", err)
	}

	defer stmt.Close()
	
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		nome VARCHAR(50),
		idade INT)`)
	if err != nil {
		log.Fatal("Erro na preparação da query", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil{
		log.Fatal("Erro ao criar tabela:", err)
	}

	_, err = db.Exec("USE hello")
	if err != nil {
		log.Fatal("Erro ao selecionar o banco de dados: ", err)
	}

	stmt, err = db.Prepare(`INSERT INTO users (nome, idade) VALUES (?, ?)`)
	if err != nil {
		log.Fatal("Erro na preparação da query", err)
	}

	defer stmt.Close()

	usuarios := []struct {
		nome string
		idade int
	}{
		{"Carol", 34},
		{"Ana", 29},
		{"Maria", 41},
	}

	for _, usuario := range usuarios {
		res, err := stmt.Exec(usuario.nome, usuario.idade)
		if err != nil {
			log.Fatal("Erro na inserção de dados", err)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		rowCount, err := res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID = %d, affected = %d\n", lastId, rowCount)
	}

	var (
		id int
		nome string
	)
	stmt, err = db.Prepare("SELECT id, nome FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		err := rows.Scan(&id, &nome)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, nome)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}