# Requisitos: API REST de Gestão de Tarefas em Go

Projeto de portfólio para vaga de Desenvolvedor Backend Go Jr.

## 1. Visão geral

API REST para gestão de tarefas organizadas por projeto. O usuário se cadastra, cria projetos e organiza tarefas com status, prioridade, prazo e tags. O objetivo do projeto é demonstrar: Go idiomático, camadas bem separadas, banco relacional, autenticação, testes e ambiente reproduzível com Docker.

## 2. Escopo

**Dentro do escopo:** autenticação JWT, CRUD de projetos, tarefas, tags e comentários, filtros e paginação, testes, Docker, CI e documentação Swagger.

**Fora do escopo (não faça agora):** frontend, upload de arquivos, notificações em tempo real, multi-tenant, microsserviços.

## 3. Stack sugerida

- Go 1.22 ou superior
- Roteador: Chi (ou Gin/Echo)
- Banco: PostgreSQL 16 com pgx (ou sqlx)
- Migrations: golang-migrate
- Auth: golang-jwt/jwt e bcrypt (golang.org/x/crypto)
- Validação: go-playground/validator
- Config: variáveis de ambiente (.env + godotenv ou envconfig)
- Logs: `log/slog`
- Testes: `testing` e testify; testcontainers-go (opcional)
- Docs: swaggo/swag
- Infra: Docker, docker-compose, Makefile, GitHub Actions, golangci-lint
- Opcional: Redis para cache

## 4. Requisitos funcionais

### 4.1 Autenticação e usuários

- **RF01:** cadastro com nome, email e senha. Email único.
- **RF02:** login com email e senha, retornando access token (JWT, 15 min) e refresh token (7 dias).
- **RF03:** renovar o access token usando o refresh token.
- **RF04:** `GET /me` retorna os dados do usuário logado.
- **RF05:** `PATCH /me` atualiza nome e senha do próprio usuário.
- **RF06:** todas as rotas, exceto register, login, refresh e health, exigem token válido.

### 4.2 Projetos

- **RF07:** criar projeto (nome e descrição).
- **RF08:** listar projetos do usuário com paginação.
- **RF09:** buscar projeto por id.
- **RF10:** atualizar projeto.
- **RF11:** excluir projeto (soft delete), ocultando também suas tarefas.
- **RF12:** somente o dono acessa o projeto.

### 4.3 Tarefas

- **RF13:** criar tarefa em um projeto: título (obrigatório), descrição, status (`todo`, `doing`, `done`), prioridade (`low`, `medium`, `high`) e prazo (`due_date`).
- **RF14:** listar tarefas do projeto com paginação (`page`, `page_size`), filtros por status, prioridade, tag e intervalo de prazo, busca por título (`q`) e ordenação (`sort`).
- **RF15:** buscar tarefa por id.
- **RF16:** atualização completa (PUT) e parcial (PATCH).
- **RF17:** alterar somente o status via `PATCH /tasks/{id}/status`.
- **RF18:** excluir tarefa.
- **RF19:** ao marcar como `done`, registrar `completed_at`; ao reabrir, limpar o campo.

### 4.4 Tags

- **RF20:** CRUD de tags do usuário. Nome único por usuário.
- **RF21:** associar e desassociar tags de uma tarefa (relação N:N).

### 4.5 Comentários

- **RF22:** adicionar comentário em uma tarefa.
- **RF23:** listar comentários de uma tarefa.
- **RF24:** excluir o próprio comentário.

### 4.6 Utilitários

- **RF25:** `GET /health` (liveness) e `GET /ready` (checa a conexão com o banco).
- **RF26:** `GET /projects/{id}/stats` retorna contagem de tarefas por status e quantidade de tarefas atrasadas.

## 5. Regras de negócio

- **RN01:** senha com no mínimo 8 caracteres, guardada com bcrypt e nunca retornada em respostas.
- **RN02:** email normalizado (trim e minúsculas) antes de salvar e comparar.
- **RN03:** o usuário só vê e altera os próprios recursos. Acessar recurso de outro usuário retorna 404 (não revela que existe).
- **RN04:** título de tarefa entre 3 e 120 caracteres; descrição até 2000.
- **RN05:** o prazo não pode estar no passado na criação da tarefa.
- **RN06:** transições de status permitidas: `todo` para `doing`, `doing` para `done`, `doing` para `todo`, `done` para `doing`. Ir de `todo` direto para `done` retorna 422.
- **RN07:** tarefa atrasada é a que tem prazo vencido e status diferente de `done`.
- **RN08:** excluir uma tag remove as associações, mas não as tarefas.
- **RN09:** `page_size` padrão 20, máximo 100. Valores inválidos retornam 400.
- **RN10:** refresh token é guardado como hash no banco e pode ser revogado.

## 6. Requisitos não funcionais

- **RNF01:** arquitetura em camadas (handler, service, repository), com dependências por interface e injeção manual no `main`.
- **RNF02:** `context.Context` propagado do request até as queries, com timeout.
- **RNF03:** graceful shutdown (SIGINT e SIGTERM) com `http.Server.Shutdown`.
- **RNF04:** configuração por variáveis de ambiente, sem segredos no código. Incluir `.env.example`.
- **RNF05:** logs estruturados com request id, método, rota, status e latência.
- **RNF06:** middlewares: request id, logger, recover, CORS, auth e (opcional) rate limit.
- **RNF07:** erros padronizados (seção 9), sem vazar detalhes internos ao cliente.
- **RNF08:** queries parametrizadas e sem N+1 (carregar tags com JOIN ou query em lote).
- **RNF09:** índices nas colunas usadas em filtros e chaves estrangeiras.
- **RNF10:** transações quando houver mais de uma escrita dependente.
- **RNF11:** cobertura mínima de 70% na camada de service.
- **RNF12:** `golangci-lint` sem warnings.
- **RNF13:** subir todo o ambiente com um único comando (`docker-compose up`).

## 7. Modelo de dados

| Tabela | Campos principais |
| --- | --- |
| users | id (UUID), name, email (único), password\_hash, created\_at, updated\_at |
| refresh\_tokens | id, user\_id, token\_hash, expires\_at, revoked\_at, created\_at |
| projects | id, user\_id, name, description, created\_at, updated\_at, deleted\_at |
| tasks | id, project\_id, title, description, status, priority, due\_date, completed\_at, created\_at, updated\_at |
| tags | id, user\_id, name (único por usuário), created\_at |
| task\_tags | task\_id, tag\_id (chave composta) |
| comments | id, task\_id, user\_id, body, created\_at |

Índices sugeridos: `tasks(project_id, status)`, `tasks(project_id, due_date)`, `projects(user_id)`, `users(email)` único, `tags(user_id, name)` único.

## 8. Endpoints

Prefixo base: `/api/v1`

| Método | Rota | Descrição | Auth |
| --- | --- | --- | --- |
| POST | /auth/register | Cadastro | Não |
| POST | /auth/login | Login | Não |
| POST | /auth/refresh | Renovar token | Não |
| GET | /me | Dados do usuário | Sim |
| PATCH | /me | Atualizar perfil | Sim |
| GET | /projects | Listar projetos | Sim |
| POST | /projects | Criar projeto | Sim |
| GET | /projects/{id} | Buscar projeto | Sim |
| PUT | /projects/{id} | Atualizar projeto | Sim |
| DELETE | /projects/{id} | Excluir projeto | Sim |
| GET | /projects/{id}/stats | Estatísticas | Sim |
| GET | /projects/{id}/tasks | Listar tarefas (filtros) | Sim |
| POST | /projects/{id}/tasks | Criar tarefa | Sim |
| GET | /tasks/{id} | Buscar tarefa | Sim |
| PUT | /tasks/{id} | Atualizar completo | Sim |
| PATCH | /tasks/{id} | Atualizar parcial | Sim |
| PATCH | /tasks/{id}/status | Alterar status | Sim |
| DELETE | /tasks/{id} | Excluir tarefa | Sim |
| GET | /tags | Listar tags | Sim |
| POST | /tags | Criar tag | Sim |
| PUT | /tags/{id} | Atualizar tag | Sim |
| DELETE | /tags/{id} | Excluir tag | Sim |
| PUT | /tasks/{id}/tags/{tagId} | Associar tag | Sim |
| DELETE | /tasks/{id}/tags/{tagId} | Desassociar tag | Sim |
| GET | /tasks/{id}/comments | Listar comentários | Sim |
| POST | /tasks/{id}/comments | Criar comentário | Sim |
| DELETE | /comments/{id} | Excluir comentário | Sim |
| GET | /health | Liveness (fora do prefixo) | Não |
| GET | /ready | Readiness (fora do prefixo) | Não |

## 9. Padrões de resposta

**Sucesso:** JSON com o recurso. Listagens retornam os itens em `data` e, em `meta`, os campos `page`, `page_size`, `total` e `total_pages`.

**Erro:** JSON com um objeto `error` contendo `code` (ex.: `validation_error`, `not_found`, `unauthorized`), `message` e, em erros de validação, `details` com uma lista de campo e mensagem.

**Status codes:**

- 200 OK, 201 Created (com header `Location`), 204 No Content
- 400 parâmetros inválidos, 401 sem token ou token inválido, 404 não encontrado
- 409 conflito (email ou tag duplicada), 422 regra de negócio violada
- 500 erro interno (logado, sem detalhes para o cliente)

## 10. Estrutura de pastas

```
cmd/api/main.go
internal/
  config/
  domain/        (entidades e erros de domínio)
  handler/       (HTTP, DTOs, validação)
  service/       (regras de negócio)
  repository/    (acesso ao banco)
  middleware/
  auth/          (JWT e hash de senha)
migrations/
docs/            (swagger gerado)
Dockerfile
docker-compose.yml
Makefile
.env.example
.github/workflows/ci.yml
README.md
```

## 11. Testes

- **Unitários:** services com repositories mockados, em table-driven tests.
- **Handlers:** `net/http/httptest` cobrindo status e corpo da resposta.
- **Integração:** repositories contra um Postgres real (testcontainers ou compose de teste).
- **Casos obrigatórios:** email duplicado no cadastro, login inválido, acesso sem token, acesso a recurso de outro usuário, transição de status inválida, paginação e filtros combinados.
- Rodar `go test -race ./...` no CI.

## 12. Docker, Makefile e CI

- **Dockerfile** multi-stage, imagem final enxuta (alpine ou distroless) e usuário não root.
- **docker-compose:** serviços `api` e `postgres` (com healthcheck); `redis` se usar cache.
- **Makefile:** `run`, `build`, `test`, `lint`, `migrate-up`, `migrate-down`, `swag`.
- **GitHub Actions:** lint, testes e build a cada push e pull request.

## 13. Documentação

O README deve ter: descrição do projeto, stack, como rodar (passo a passo), variáveis de ambiente, exemplos de `curl`, diagrama ER, decisões técnicas e o que ficou fora. Swagger disponível em `/swagger/index.html`.

## 14. Roadmap sugerido

1. **Fundação (2 a 3 dias):** estrutura de pastas, config, Postgres no compose, migrations, health, logger e graceful shutdown.
2. **Autenticação (2 a 3 dias):** register, login, refresh, middleware JWT e `/me`, com testes.
3. **Projetos e tarefas (3 a 4 dias):** CRUDs, regras de negócio e transições de status.
4. **Filtros, tags e comentários (3 dias):** paginação, busca, ordenação, N:N e stats.
5. **Qualidade (2 a 3 dias):** cobertura de testes, lint, CI e tratamento de erros padronizado.
6. **Polimento (1 a 2 dias):** Swagger, README, Dockerfile final e revisão de commits.

## 15. Critérios de pronto

- [ ] Todos os RF e RN implementados
- [ ] `docker-compose up` sobe a API e o banco sem passos manuais
- [ ] Testes e lint passando no CI
- [ ] Swagger gerado e README completo
- [ ] Nenhum segredo commitado
- [ ] Commits pequenos e com mensagens claras

## 16. Diferenciais opcionais

- Cache do endpoint de stats com Redis e invalidação ao alterar tarefas
- Rate limiting por IP ou usuário
- Worker em goroutine com `time.Ticker` e `context` que varre tarefas com prazo próximo e registra lembretes
- Métricas Prometheus em `/metrics`
- Refresh token rotativo com detecção de reuso

## 17. Perguntas de entrevista sobre este projeto

- Por que escolheu esse roteador e como funcionam os middlewares?
- Como evitou SQL injection e o problema de N+1?
- Como o JWT funciona e por que existe refresh token?
- Por que usar interfaces nos repositories?
- Como funciona o graceful shutdown e o que acontece sem ele?
- Por que retornar 404 em vez de 403 para recurso de outro usuário?
- Como você escalaria a API para mais tráfego?
- Onde usou `context` e por quê?
