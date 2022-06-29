package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type pessoa struct {
	Id        int    `json:"id"`
	Nome      string `json:"nome"`
	Sobrenome string `json:"sobrenome"`
}

func lerDB(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	registro, erroQuery := db.Query("SELECT * FROM pessoas;")
	if erroQuery != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	for registro.Next() {
		var ps pessoa
		erroScan := registro.Scan(&ps.Id, &ps.Nome, &ps.Sobrenome)
		if erroScan != nil {
			log.Println(erroScan.Error())
			continue
		}
		c.JSON(200, gin.H{
			"id":        ps.Id,
			"nome":      ps.Nome,
			"sobrenome": ps.Sobrenome,
		})
	}
}
func buscarP(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	id := c.Param("id")

	count := db.QueryRow("SELECT COUNT(1) FROM pessoas WHERE id = ?;", id)
	var exists int
	count.Scan(&exists)

	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"resposta": "a pessoa não existe no banco de dados",
		})
		return
	}

	registro := db.QueryRow("SELECT * FROM pessoas WHERE id = ?;", id)
	var ps pessoa
	erroScan := registro.Scan(&ps.Id, &ps.Nome, &ps.Sobrenome)
	if erroScan != nil {
		c.Status(http.StatusBadGateway)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"id":        ps.Id,
		"nome":      ps.Nome,
		"sobrenome": ps.Sobrenome,
	})
}
func criarP(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	body := c.Request.Body
	corpo, erroRead := ioutil.ReadAll(body)
	if erroRead != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var novaPessoa pessoa
	json.Unmarshal(corpo, &novaPessoa)

	_, erroExec := db.Exec("INSERT INTO pessoas() VALUES(?,?,?);", novaPessoa.Id, novaPessoa.Nome, novaPessoa.Sobrenome)
	if erroExec != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        novaPessoa.Id,
		"nome":      novaPessoa.Nome,
		"sobrenome": novaPessoa.Sobrenome,
	})
}
func deletarP(c *gin.Context) {
	id := c.Param("id")

	count := db.QueryRow("SELECT COUNT(1) FROM pessoas WHERE id = ?;", id)
	var exists int
	count.Scan(&exists)

	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"resposta": "a pessoa não existe no banco de dados",
		})
		return
	}

	_, erroExec := db.Exec("DELETE FROM pessoas WHERE id = ?", id)
	if erroExec != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
func atualizarP(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	id := c.Param("id")

	count := db.QueryRow("SELECT COUNT(1) FROM pessoas WHERE id = ?;", id)
	var exists int
	count.Scan(&exists)

	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"resposta": "a pessoa não existe no banco de dados",
		})
		return
	}

	bory := c.Request.Body
	corpo, erroRead := ioutil.ReadAll(bory)
	if erroRead != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var pessoaAtualizada pessoa
	json.Unmarshal(corpo, &pessoaAtualizada)

	_, erroExec := db.Exec("UPDATE pessoas SET id = ?, nome = ?, sobrenome = ? WHERE id = ?", pessoaAtualizada.Id, pessoaAtualizada.Nome, pessoaAtualizada.Sobrenome, id)
	if erroExec != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"id":        pessoaAtualizada.Id,
		"nome":      pessoaAtualizada.Nome,
		"sobrenome": pessoaAtualizada.Sobrenome,
	})
}
func princi(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)

	c.JSON(200, gin.H{
		"menssagem": "insira os dados de pesquisa, crição, exclusão ou atualização dos clientes",
	})
}

func conecDataBase() {
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
	conecDataBase()
	r := gin.Default()
	r.GET("/", princi)
	r.GET("/pessoas", lerDB)
	r.GET("/pessoas/:id", buscarP)
	r.DELETE("pessoas/:id", deletarP)
	r.PUT("pessoas/:id", atualizarP)
	r.POST("/pessoas", criarP)

	http.ListenAndServe(":8080", r)

}
