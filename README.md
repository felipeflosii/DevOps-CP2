# go-api

API REST leve em Go conectada ao MySQL, rodando em Docker.

## Requisitos

- Docker
- Docker Network `flosi-rede` criada

## Como rodar

### 1. Criar a rede
```bash
docker network create flosi-rede
```

### 2. Subir o MySQL
```bash
docker run -d --name mysql-db --network flosi-rede \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=demo \
  -p 3306:3306 --memory="256m" \
  mysql:8.0 --performance_schema=OFF --innodb_buffer_pool_size=64M
```

### 3. Criar a tabela
```bash
docker exec -it mysql-db mysql -uroot -proot demo -e "
CREATE TABLE produtos (
  id bigint NOT NULL AUTO_INCREMENT,
  nome varchar(255) NOT NULL,
  categoria varchar(255),
  preco float NOT NULL,
  PRIMARY KEY (id)
);"
```

### 4. Build da imagem Go
```bash
docker build -t go-api .
```

### 5. Subir a API
```bash
docker run -d --name go-api --network flosi-rede \
  -p 8080:8080 --memory="20m" \
  go-api
```

---

## Endpoints

### Health check
```
GET /health
```
```bash
curl http://localhost:8080/health
```

---

### Listar produtos
```
GET /produtos
```
```bash
curl http://localhost:8080/produtos
```

---

### Criar produto
```
POST /produtos
```
```bash
curl -X POST http://localhost:8080/produtos \
  -H "Content-Type: application/json" \
  -d '{"nome":"Notebook","categoria":"Tech","preco":3500.00}'
```

---

### Buscar produto por ID
```
GET /produtos/{id}
```
```bash
curl http://localhost:8080/produtos/1
```

---

### Deletar produto
```
DELETE /produtos/{id}
```
```bash
curl -X DELETE http://localhost:8080/produtos/1
```

---

## Variáveis de ambiente

| Variável   | Padrão     | Descrição          |
|------------|------------|--------------------|
| DB_HOST    | mysql-db   | Host do banco      |
| DB_PORT    | 3306       | Porta do banco     |
| DB_USER    | root       | Usuário do banco   |
| DB_PASS    | root       | Senha do banco     |
| DB_NAME    | demo       | Nome do banco      |
| PORT       | 8080       | Porta da API       |

---

## Monitorar uso de memória

```bash
docker stats --no-stream
```

### Listar produtos fora da VM
```
GET /produtos
```

```bash
curl http://163.176.225.114:8080/produtos
```
