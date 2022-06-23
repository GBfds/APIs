package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//variaveis e structs
var db *sql.DB

type pessoa struct {
	Id        int    `json:"id"`
	Nome      string `json:"nome"`
	Sobrenome string `json:"sobrenome"`
}

//hundles
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "td ok")
}
func criarPss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	corpo, erroReadAll := ioutil.ReadAll(r.Body)
	if erroReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var novaPessoa pessoa
	json.Unmarshal(corpo, &novaPessoa)

	_, erroExec := db.Exec("INSERT INTO pessoas() VALUES(?,?,?);", novaPessoa.Id, novaPessoa.Nome, novaPessoa.Sobrenome)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func atualizarPss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[2])
	if erroSplit != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	corpo, erroReadAll := ioutil.ReadAll(r.Body)
	if erroReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var pessoaAtualizada pessoa
	json.Unmarshal(corpo, &pessoaAtualizada)

	_, erroExec := db.Exec("UPDATE pessoas SET id = ?, nome = ?, email = ? WHERE id = ?", pessoaAtualizada.Id, pessoaAtualizada.Nome, pessoaAtualizada.Sobrenome, id)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pessoaAtualizada)
}
func lerBD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	registro, erroQuery := db.Query("SELECT * FROM pessoas;")
	if erroQuery != nil {
		log.Println(erroQuery.Error())
		return
	}
	var pss []pessoa = make([]pessoa, 0)
	for registro.Next() {
		var ps pessoa
		erroScan := registro.Scan(&ps.Id, &ps.Nome, &ps.Sobrenome)
		if erroScan != nil {
			log.Println(erroScan.Error())
			continue
		}

		pss = append(pss, ps)
	}
	json.NewEncoder(w).Encode(pss)
}
func buscarPss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[2])
	if erroSplit != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	registro := db.QueryRow("SELECT * FROM pessoas WHERE id = ?", id)
	var ps pessoa
	erroScan := registro.Scan(&ps.Id, &ps.Nome, &ps.Sobrenome)
	if erroScan != nil {
		log.Println(erroScan.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(ps)
}
func deletarPss(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[2])
	if erroSplit != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	registro := db.QueryRow("SELECT * FROM pessoas WHERE id = ?;", id)
	var ps pessoa
	erroScan := registro.Scan(&ps.Id, &ps.Nome, &ps.Sobrenome)
	if erroScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, erroExec := db.Exec("DELETE FROM pessoas WHERE id = ?", id)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//configurações
func rotas() {
	http.HandleFunc("/", root)
	http.HandleFunc("/pessoas", rotearDD)
	http.HandleFunc("/pessoas/", rotearDD)
}
func rotearDD(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")

	if len(partes) == 2 || len(partes) == 3 && partes[2] == "" {
		if r.Method == "GET" {
			lerBD(w, r)
		} else if r.Method == "POST" {
			criarPss(w, r)
		}
	} else if len(partes) == 3 || len(partes) == 4 && partes[3] == "" {
		if r.Method == "GET" {
			buscarPss(w, r)
		} else if r.Method == "PUT" {
			atualizarPss(w, r)
		} else if r.Method == "DELETE" {
			deletarPss(w, r)
		}
	}
}
func conectarDb() {
	var erroDeConecçao error
	db, erroDeConecçao = sql.Open("mysql", "usuario:senha/nomeDB")
	if erroDeConecçao != nil {
		log.Println("erro na conecção do banco de dados" + erroDeConecçao.Error())
		log.Fatal()
	}

	erroPing := db.Ping()
	if erroPing != nil {
		log.Println("erro na conecção do banco de dados" + erroPing.Error())
		log.Fatal()
	}
}
func main() {
	conectarDb()
	rotas()
	fmt.Println("ok...")
	http.ListenAndServe(":8080", nil)
}
