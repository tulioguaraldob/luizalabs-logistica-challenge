# Desafio Técnico - Vertical Logistica
**luizalabs-logistica-challenges** 🚀

## Tech:

### Ferramentas:

Cada uma das tecnologias utilizadas está listada abaixo de forma resumida, e após a lista há uma descrição de cada uma delas e a razão porque foram escolhidas.

* Go
* Docker
* PostgreSQL
* Testify
* Gomock
* Migrate
* Makefile

#### Go:
Foi escolhido [Go](https://go.dev/) como linguagem de programação para o desafio, devido a sua eficiência, simplicidade, tipagem forte, performance para processar muitas linhas de um file ou request de uma vez e seus compile times rápidos.

Go possui poderosas libs nativa, as quais foram utilizadas em questões de OS e File Handling para melhor leitura dos arquivos enviados por meio das libs "os" e "bufio".
Questões de conexão com o banco de dados SQL (Postgres), por meio da lib nativa "database/sql" a qual foi totalmente utilizada para operações de seleção e inserção no banco, após parse de cada um dos dados nos arquivos.

Utilizou-se também a lib "net/http" nativa do Go, sem o uso de frameworks externos como Gin e Fiber, ou libs de routing como Chi. Por meio apenas da lib nativa foi possível receber o arquivo esperado, filtrar query params como id, receber campos para intervalo de data de compra (data início e data fim) e também fazer a response em JSON como esperado no desafio.

#### Docker:
Foi escolhido [Docker](https://www.docker.com/), devido sua eficiência, praticidade e robustez na questão de gerenciamento e criação de containers para questões de desenvolvimento local como o uso do Postgres como database que veio a ser utilizado na solução final. Também questões de solução final como build do container de nossa aplicação final.

#### PostgreSQL:
Foi escolhido [Postgres](https://www.postgresql.org/) como banco de dados SQL para fazer a camada de persistência de nossa aplicação devido a sua praticidade e robustez na hora de realizar queries para operações de read and write e também a facilidade de migrar as entidades (entities) da nossa aplicação para migrations em arquivos sql para criar ou deletar tabelas baseadas nessas entidades.

#### Testify:
Foi escolhido a lib [testify](https://github.com/stretchr/testify) devido o fato de ser um toolkit bem completo para assertions e mocks que se integra muito bem com a lib nativa "testing" do Go.

#### Gomock:
Foi escolhido o framework [gomock](https://github.com/uber-go/mock), pois permite facilmente gear mocks de interfaces para unit tests e testes E2E. Além da questão de realizar o mock, também é possível realizar o expect dos parametros esperados da dependência a ser mockada e o que a mesma irá retornar.

#### Migrate:
O [migrate](https://github.com/golang-migrate/migrate) foi escolhido pela facilidade em realizar ups (criações) e downs (deleções) em migrations sql e também pela facilidade em gerar os arquivos responsáveis por armazenar tais migrations e também o fato de ter um driver de conexão para quase qualquer banco de dados SQL.

#### Makefile:
O Makefile foi escolhido para abstrair e simplificar comandos como run, test, install e cover.

### Padrão Arquitetural:

O projeto foi construído seguindo uma mescla de Clean Architecture com formas de abordagem de DDD (Domain-Driven Design) e o uso de SOLID.
Ao utilizar essa junção de Clean Architecture com as abordagens de bounded context de DDD, temos o seguinte resultado arquitetural:

* Domain: é o core (núcleo) da nossa aplicação em questões de regras de negócios, separação de domínios de usuários (users), produtos (products), pedidos (orders) e as relações entre eles.
Além do mais, é nessa camada a qual está presente as nossas interfaces de Repository e Service, e também nossos schemas que representam entrada e saída de dados.
Com essa abordagem do Domain, todo o que cerca a nossa aplicação poderia ser facilmente alterado, poderiamos mudar com extrema facilidade de Postgres para MongoDB, ou até mesmo substituir nossa camada HTTP por outra camada de comunicação externa.

* Adapter: são as partes externas a nossa aplicação, e que podem ser alteradas ou adicionadas facilmente, é aqui onde está a nossa conexão com nosso banco de dados Postgres e é aqui onde se encontra as controllers utilizadas com "net/http".

## Executando a aplicação:

Para rodar este projeto local, é ideal que se tenha Docker instalado em sua máquina, juntamente com `docker-compose`

Antes de rodarmos um `make run` ou realizarmos a build de nosso `Dockerfile` devemos configurar nossas váriaveis de ambiente.

Podemos configurar essas váriaveis copiando o comando abaixo e então copiando o conteúdo do exemplo abaixo.

### 1. Create .env:
```bash
cp .env.example .env
```

### 2. Preencher variáveis de ambiente:
```bash
POSTGRES_USER=labs
POSTGRES_PASSWORD=labs1234
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=labsdb
PORT=8080
```

Feito isso precisamos realizar um comando Make para subir nosso Docker container de Postgres:

### 3. Postgres Container:
```bash
make setup
```

Por fim, é só executar a aplicação, podemos executar diretamente com Go instalado em nossa máquina ou podemos realizar a build de um container Docker para isso, conforme nosso Dockerfile.

### 4. Rodando diretamente com Go:
```bash
make run
```

### 5. Rodando com Docker:
```bash
make docker-build
```

## Endpoints:

A aplicação possui 5 endpoints:

* [GET] /healthcheck: Health check para verificar se o servidor está OK
* [GET] /user/{id}: Busca o usuário pelo ID salvo na base
* [POST] /user/upload: Endpoint o qual é feito upload do arquivo de dados por meio do Form Multipart, na key users_data
* [GET] /order/{id}: Busca um pedido específico pelo ID salvo na base
* [GET] /orders: Busca todos os pedidos, com possibilidade de filtrar por id, ou data de início e data final.