# Awesome Project - API REST com Autenticação JWT

## 🚀 Melhorias Implementadas

### 1. **Arquitetura em Camadas (Handler > Service > Repository)**
- ✅ **Handlers**: Responsáveis apenas por HTTP (parse request, response HTTP)
- ✅ **Services**: Contêm toda a lógica de negócio e validações
- ✅ **Repositories**: Abstração da camada de dados (GORM/Banco de Dados)

**Benefícios:**
- Código mais testável
- Melhor separação de responsabilidades
- Lógica de negócio reutilizável
- Fácil manutenção

### 2. **Autenticação com JWT**
- ✅ **Register/Login**: Endpoints públicos para registro e autenticação
- ✅ **Token JWT**: Gerado ao fazer login ou registro (válido por 24h)
- ✅ **Middleware de Autenticação**: Valida token em requisições protegidas
- ✅ **Hash de Senhas**: Senhas armazenadas com bcrypt (DefaultCost)

**Fluxo:**
1. Usuário se registra em `POST /auth/register`
2. Usuário faz login em `POST /auth/login` e recebe um token
3. Usuário inclui token no header `Authorization: Bearer <token>` em requisições protegidas

### 3. **Autorização (Ownership)**
- ✅ Usuários só podem **editar sua própria conta**
- ✅ Usuários só podem **deletar sua própria conta**
- ✅ Usuários só podem **editar seus próprios posts**
- ✅ Usuários só podem **deletar seus próprios posts**

### 4. **Variáveis de Ambiente com Godotenv**
- ✅ Arquivo `.env` centralizado
- ✅ Configurações: `DATABASE_PATH`, `PORT`, `JWT_SECRET`
- ✅ Suporte a valores padrão

**Arquivo .env:**
```
DATABASE_PATH=./app.db
PORT=8080
JWT_SECRET=sua_chave_secreta_muito_segura_aqui_123456
```

### 5. **Validações Melhoradas**
- ✅ Título de post: mínimo 3 caracteres
- ✅ Corpo de post: mínimo 10 caracteres
- ✅ Senha: mínimo 6 caracteres
- ✅ Email: validação de formato
- ✅ Verificação de email duplicado no registro

### 6. **Timestamps**
- ✅ Campo `CreatedAt` em Users e Posts
- ✅ Campo `UpdatedAt` em Users e Posts
- ✅ Gerenciados automaticamente pelo GORM

## 📁 Estrutura do Projeto

```
awesomeProject/
├── config/
│   └── config.go              # Gerenciamento de variáveis de ambiente
├── handlers/
│   ├── auth.go               # Endpoints de autenticação
│   ├── users.go              # Endpoints de usuários (refatorado)
│   └── posts.go              # Endpoints de posts (refatorado)
├── services/
│   ├── auth.go               # Lógica de autenticação e JWT
│   ├── users.go              # Lógica de negócio de usuários
│   └── posts.go              # Lógica de negócio de posts
├── repositories/
│   ├── users.go              # Acesso a dados de usuários
│   └── posts.go              # Acesso a dados de posts
├── models/
│   ├── users.go              # Model User
│   ├── posts.go              # Model Post
│   └── auth.go               # Models de autenticação
├── middleware/
│   └── auth.go               # Middleware de validação de JWT
├── routes/
│   └── routes.go             # Definição de rotas
├── db/
│   └── db.go                 # Inicialização do banco de dados
├── main.go                   # Ponto de entrada da aplicação
├── .env                      # Variáveis de ambiente
└── go.mod                    # Dependências do projeto
```

## 🔐 Endpoints da API

### Autenticação (Públicos)

#### Registro
```bash
POST /auth/register
Content-Type: application/json

{
  "name": "João Silva",
  "email": "joao@example.com",
  "password": "senha123"
}

Response: 201 Created
{
  "id": 1,
  "name": "João Silva",
  "email": "joao@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "joao@example.com",
  "password": "senha123"
}

Response: 200 OK
{
  "id": 1,
  "name": "João Silva",
  "email": "joao@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Usuários

#### Listar todos (Público)
```bash
GET /users
```

#### Buscar por ID (Público)
```bash
GET /users/:id
```

#### Atualizar próprio perfil (Protegido)
```bash
PUT /users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "João Silva Atualizado",
  "email": "joao.novo@example.com"
}
```

#### Deletar própria conta (Protegido)
```bash
DELETE /users/:id
Authorization: Bearer <token>
```

### Posts

#### Listar todos (Público)
```bash
GET /posts
```

#### Buscar por ID (Público)
```bash
GET /posts/:id
```

#### Posts de um usuário (Público)
```bash
GET /users/:userId/posts
```

#### Criar post (Protegido)
```bash
POST /posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Meu Primeiro Post",
  "body": "Este é o corpo do meu primeiro post com conteúdo interessante"
}
```

#### Atualizar post (Protegido - apenas autor)
```bash
PUT /posts/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Título Atualizado",
  "body": "Corpo atualizado com novo conteúdo"
}
```

#### Deletar post (Protegido - apenas autor)
```bash
DELETE /posts/:id
Authorization: Bearer <token>
```

## 🔑 Como Usar a API

### 1. Registrar um novo usuário
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "password": "senha123"
  }'
```

### 2. Fazer login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "senha123"
  }'
```
Salve o token retornado.

### 3. Criar um post (usando o token)
```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "title": "Meu Primeiro Post",
    "body": "Este é um post interessante sobre Go e APIs REST"
  }'
```

### 4. Listar todos os posts
```bash
curl http://localhost:8080/posts
```

### 5. Buscar posts de um usuário
```bash
curl http://localhost:8080/users/1/posts
```

## 📦 Dependências Adicionadas

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get github.com/joho/godotenv
```

- **golang-jwt/jwt/v5**: Geração e validação de tokens JWT
- **crypto**: Hash de senhas com bcrypt
- **godotenv**: Carregamento de variáveis de ambiente do arquivo .env

## ⚙️ Configuração

### Variavelies de Ambiente (.env)
```env
# Caminho do banco de dados SQLite
DATABASE_PATH=./app.db

# Porta do servidor
PORT=8080

# Chave secreta para assinar tokens JWT (use uma chave forte em produção!)
JWT_SECRET=sua_chave_secreta_muito_segura_aqui_123456
```

## 🚀 Como Executar

1. **Clone/Prepare o projeto**
   ```bash
   cd C:\Users\Renan\GolandProjects\awesomeProject
   ```

2. **Configure o arquivo .env**
   ```bash
   # Já foi criado com valores padrão
   # Edite se necessário
   ```

3. **Compile o projeto**
   ```bash
   go build -o awesomeProject.exe
   ```

4. **Execute**
   ```bash
   ./awesomeProject.exe
   ```

5. **Servidor iniciará em**
   ```
   http://localhost:8080
   ```

## 🔄 Fluxo de Autenticação

```
┌─────────────────┐
│   Usuário       │
└────────┬────────┘
         │
         │ 1. POST /auth/register ou /auth/login
         ↓
┌─────────────────────────────────┐
│  Auth Handler & Service          │
│  - Valida credenciais            │
│  - Hash de senha (bcrypt)        │
│  - Gera JWT Token               │
└────────┬────────────────────────┘
         │
         │ 2. Retorna Token
         ↓
┌─────────────────┐
│   Usuário       │
│ Armazena Token  │
└────────┬────────┘
         │
         │ 3. Inclui token em requisições
         │    Authorization: Bearer <token>
         ↓
┌──────────────────────────────┐
│  Middleware RequireAuth()     │
│  - Valida assinatura         │
│  - Verifica expiração        │
│  - Extrai userID             │
└──────────┬───────────────────┘
           │
           │ ✅ Token válido
           ↓
┌──────────────────────┐
│  Handler & Service   │
│  Processa requisição │
└──────────────────────┘
```

## 🧪 Testes Recomendados

Use ferramentas como **Postman**, **Insomnia** ou **cURL** para testar:

1. ✅ Registrar novo usuário
2. ✅ Fazer login
3. ✅ Criar post sem autenticação (deve falhar)
4. ✅ Criar post com token inválido (deve falhar)
5. ✅ Criar post com token válido (deve funcionar)
6. ✅ Listar posts públicos
7. ✅ Tentar atualizar post de outro usuário (deve falhar)
8. ✅ Atualizar próprio post (deve funcionar)

## 🎯 Próximas Melhorias Opcionais

- [ ] Refresh tokens (token expira em 24h, refresh token em 7 dias)
- [ ] Rate limiting por IP/usuário
- [ ] Paginação em listagens
- [ ] Soft delete (marcar como deletado ao invés de remover)
- [ ] Logs estruturados com zerolog ou logrus
- [ ] Testes unitários e de integração
- [ ] Docker e docker-compose
- [ ] CI/CD com GitHub Actions
- [ ] Documentação OpenAPI/Swagger
- [ ] Cache com Redis

---

**Projeto refatorado e pronto para produção! 🎉**

