# APIs
---
 APIs com golang e MySQL

# MySQL
---
## conecção com banco de dados MySQL

- instali o pacote do driver com o comando abaixo
```
$ go get -u github.com/go-sql-driver/mysql
```

- tenha seu banco de dados já criado no seu MySQL
```
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
```
- na terceira linha entre as duas ultimas aspas, basta trocar o __usuario__ pelo seu usuario do MySQL, __senha__ pela sua senha e __nomeDB__ por o nome do seu banco de dados

- se der tudo certo na conecção, basta criar as tabelas dentro do banco de dados que você deseja fazer a conecção 




para melhor explicaçao acesse a pagina do criador do [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql#go-mysql-driver)
## criação de tabelas
### tabela da *api_cliente_pix.go*
1. crie a tabela de clientes com o comando sql abaixo
```
CREATE TABLE clientes(id INT PRIMARY KEY auto_increment, nome VARCHAR(80) NOT NULL,email VARCHAR(100) NOT NULL);
```
2. __depois__ de criar a tabela de clientes, crie a tabela de pix

```
 CREATE TABLE pix_clientes(id_clt INT, pix1 VARCHAR(100),pix2 VARCHAR(100),pix3 VARCHAR(100), FOREIGN KEY(id_clt) REFERENCES clientes(id));

```

### tabela da *api_basica.go*

```
CREATE TABLE pessoas (id INT PRIMARY KEY, nome VARCHAR(80) NOT NULL,sobrenome VARCHAR(80) NOT NULL);
```
### tabela da *api_basica_com_gin.go*
essa API é apenas um melhoria da __api_basica.go__ então é usada a mesma tabela
Lembre-se apenas de instalar o gin com o comando
```
$ go get -u github.com/gin-gonic/gin
```
para melhor explicaçao acesse a pagina do criador do [Gin Web Framework](https://github.com/gin-gonic/gin)