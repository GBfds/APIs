# APIs
---
 APIs com golang e mysql

# mysql
---
## conecção com banco de dados mysql
para melhor explicaçao acesse a pagina do criador do [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql#go-mysql-driver)
## criação de tabelas
### tabela da *api01.go*
1. pe
```
CREATE TABLE clientes(id INT PRIMARY KEY auto_increment, nome VARCHAR(80) NOT NULL,email VARCHAR(100) NOT NULL);
```

```
 CREATE TABLE pix_clientes(id_clt INT, pix1 VARCHAR(100),pix2 VARCHAR(100),pix3 VARCHAR(100), FOREIGN KEY(id_clt) REFERENCES clientes(id));

```
