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

//Variaveis e Tipos
var db *sql.DB

type cliente struct {
	Id    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type cltPix struct {
	IdClt int    `json:"idClt"`
	Nome  string `json:"nome"`
	Pix1  string `json:"pix1"`
	Pix2  string `json:"pix2"`
	Pix3  string `json:"pix3"`
}

type texto struct {
	Pix string `json:"pix"`
}
type respostaEmTexto struct {
	Resposta string `json:"resposta"`
}

//Hundles
func pix(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[3])
	if erroSplit != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	registro := db.QueryRow("SELECT pix_clientes.id_clt,clientes.nome,pix_clientes.pix1,pix_clientes.pix2,pix_clientes.pix3 FROM pix_clientes INNER JOIN clientes ON clientes.id = pix_clientes.id_clt WHERE clientes.id = ?;", id)
	var clt cltPix
	erroScan := registro.Scan(&clt.IdClt, &clt.Nome, &clt.Pix1, &clt.Pix2, &clt.Pix3)
	if erroScan != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(respostaEmTexto{"o cliente não existe"})
		return
	}

	json.NewEncoder(w).Encode(clt)
}
func atualizarPix(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[3])
	if erroSplit != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	corpo, erroReadAll := ioutil.ReadAll(r.Body)
	if erroReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var novoPix texto
	json.Unmarshal(corpo, &novoPix)

	registro := db.QueryRow("SELECT pix_clientes.id_clt,clientes.nome,pix_clientes.pix1,pix_clientes.pix2,pix_clientes.pix3 FROM pix_clientes INNER JOIN clientes ON clientes.id = pix_clientes.id_clt WHERE clientes.id = ?;", id)
	var clt cltPix
	erroScan := registro.Scan(&clt.IdClt, &clt.Nome, &clt.Pix1, &clt.Pix2, &clt.Pix3)
	if erroScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if clt.Pix1 == "" {
		clt.Pix1 = novoPix.Pix
	} else if clt.Pix2 == "" {
		clt.Pix2 = novoPix.Pix
	} else if clt.Pix3 == "" {
		clt.Pix3 = novoPix.Pix
	} else {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(respostaEmTexto{"Todas as chaves já foram cadastradas"})
		return
	}

	_, erroExec := db.Exec("UPDATE pix_clientes SET pix1 = ?, pix2 = ?, pix3 = ? WHERE id_clt = ?;", clt.Pix1, clt.Pix2, clt.Pix3, id)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(clt)
}

func raiz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "insira os dados de pesquisa, crição, exclusão ou atualização dos clientes")
}
func buscarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[2])
	if erroSplit != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	registro := db.QueryRow("SELECT * FROM clientes WHERE id = ?", id)
	var clt cliente
	erroScan := registro.Scan(&clt.Id, &clt.Nome, &clt.Email)
	if erroScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(clt)
}
func lerBancoDeDados(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	registro, erroQuery := db.Query("SELECT * FROM clientes;")
	if erroQuery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var clts []cliente = make([]cliente, 0)
	for registro.Next() {
		var clt cliente
		erroScan := registro.Scan(&clt.Id, &clt.Nome, &clt.Email)
		if erroScan != nil {
			log.Println(erroScan.Error())
			continue
		}

		clts = append(clts, clt)
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(clts)
}
func deletarCliente(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[2])
	if erroSplit != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//verificação se o cliente existe
	registro := db.QueryRow("SELECT * FROM clientes WHERE id = ?;", id)
	var clt cliente
	erroScan := registro.Scan(&clt.Id, &clt.Nome, &clt.Email)
	if erroScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//...
	// deletar primeiro da tabela de chaves pix
	_, erroExec2 := db.Exec("DELETE FROM pix_clientes WHERE id_clt = ?", id)
	if erroExec2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, erroExec := db.Exec("DELETE FROM clientes WHERE id = ?", id)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func atualizarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(r.URL.Path, "/")
	id, erroSplit := strconv.Atoi(partes[2])
	if erroSplit != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//verificação se o cliente existe
	registro := db.QueryRow("SELECT * FROM clientes WHERE id = ?;", id)
	var clt cliente
	erroScan := registro.Scan(&clt.Id, &clt.Nome, &clt.Email)
	if erroScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//...

	corpo, erroReadAll := ioutil.ReadAll(r.Body)
	if erroReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var clienteAtualizado cliente
	json.Unmarshal(corpo, &clienteAtualizado)

	_, erroExec := db.Exec("UPDATE clientes SET id = ?, nome = ?, email = ? WHERE id = ?", clienteAtualizado.Id, clienteAtualizado.Nome, clienteAtualizado.Email, id)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(clienteAtualizado)
}
func criarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	corpo, erroReadAll := ioutil.ReadAll(r.Body)
	if erroReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var novoCliente cliente
	json.Unmarshal(corpo, &novoCliente)

	//confirmação se o cliente existe, usando o email como chave unica
	registro := db.QueryRow("SELECT EXISTS(SELECT * FROM clientes WHERE email=?);", novoCliente.Email)
	var SN bool
	registro.Scan(&SN)
	if SN == true {
		json.NewEncoder(w).Encode(respostaEmTexto{"o email já está em uso"})
		return
	}
	//..

	_, erroExec := db.Exec("INSERT INTO clientes(nome,email) VALUES(?,?);", novoCliente.Nome, novoCliente.Email)
	if erroExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//adicionando novo cliente na tabela de chaves pix
	var add cliente
	addPix := db.QueryRow("SELECT * FROM clientes WHERE email = ?;", novoCliente.Email)
	addPix.Scan(&add.Id, &add.Nome, &add.Email)

	var strVazia string = ""
	_, errerroExec2 := db.Exec("INSERT INTO pix_clientes() VALUES(?,?,?,?);", add.Id, strVazia, strVazia, strVazia)
	if errerroExec2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//encaminhamento de rotas
func RotearDados(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")

	if len(partes) == 2 || len(partes) == 3 && partes[2] == "" {
		if r.Method == "GET" {
			lerBancoDeDados(w, r)
		} else if r.Method == "POST" {
			criarCliente(w, r)
		}
	} else if len(partes) == 3 || len(partes) == 4 && partes[3] == "" {
		if r.Method == "GET" {
			buscarCliente(w, r)
		} else if r.Method == "PUT" {
			atualizarCliente(w, r)
		} else if r.Method == "DELETE" {
			deletarCliente(w, r)
		}
	}
}
func RotearPix(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	if len(partes) == 3 || len(partes) == 4 && partes[3] == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respostaEmTexto{"Insira o ID do cliente"})
	} else if len(partes) == 4 || len(partes) == 5 && partes[4] == "" {
		if r.Method == "GET" {
			pix(w, r)
		} else if r.Method == "PUT" {
			atualizarPix(w, r)
		}
	}
}

func handles() {
	http.HandleFunc("/", raiz)
	http.HandleFunc("/clientes", RotearDados)
	http.HandleFunc("/clientes/", RotearDados)
	http.HandleFunc("/clientes/pix", RotearPix)
	http.HandleFunc("/clientes/pix/", RotearPix)
}

//configurações da API
func conectarDataBase() {
	var erroDeConecçao error
	db, erroDeConecçao = sql.Open("mysql", "usuario:senha@/nomeDB")
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
	conectarDataBase()
	handles()
	fmt.Println("ok")
	http.ListenAndServe(":8080", nil)
}
