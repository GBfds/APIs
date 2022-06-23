package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type pessoa struct {
	Id        int
	Nome      string
	Sobrenome string
}

type clt struct {
	Id    int
	Nome  string
	Email string
}

func lerBd(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	registro, erroQuery := db.Query("SELECT * FROM pessoas;")
	if erroQuery != nil {
		log.Println(erroQuery.Error())
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

func criarP(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	Body := c.Request.Body
	corpo, erroRead := ioutil.ReadAll(Body)
	if erroRead != nil {
		fmt.Print(erroRead.Error())
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	var novaPessoa pessoa
	json.Unmarshal(corpo, &novaPessoa)

	_, erroExec := db.Exec("INSERT INTO pessoas() VALUES(?,?,?);", novaPessoa.Id, novaPessoa.Nome, novaPessoa.Sobrenome)
	if erroExec != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(201, gin.H{
		"id":        novaPessoa.Id,
		"nome":      novaPessoa.Nome,
		"sobrenome": novaPessoa.Sobrenome,
	})
}
func princi(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.AbortWithStatus(http.StatusOK)

	c.JSON(200, gin.H{
		"id":    "pessoa",
		"nome":  "aleatoria",
		"email": "pessoaale@",
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
	r.GET("/pessoas", lerBd)
	r.POST("/pessoas", criarP)

	http.ListenAndServe(":8080", r)
}
