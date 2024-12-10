## Como rodar o projeto

### Pré-requisitos
- [Go](https://golang.org/doc/install) instalado.
- [Docker](https://docs.docker.com/get-docker/) instalado.

### Passos

1. Inicie o banco de dados usando Docker Compose:
   ```bash```
    docker-compose up -d

2. Rode o projeto na pasta cmd: 
    ```bash```
    go run main.go
    
### Notas:
1. Certifique-se de que o arquivo `docker-compose.yml` está configurado corretamente para o banco de dados.
2. Se o projeto usa variáveis de ambiente, inclua uma seção no `README.md` explicando como configurá-las (ex.: criar um arquivo `.env`).

Se precisar ajustar ou adicionar mais informações ao README, me avise! 😊