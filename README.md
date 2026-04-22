# Awesome Project - API REST com Autenticação JWT

Uma API REST modular e segura construída em Go com autenticação JWT, arquitetura em camadas e banco de dados SQLite.

## ✨ Funcionalidades Implementadas

### 1. **Arquitetura em Camadas**
- ✅ **Handlers**: Responsáveis apenas por HTTP (parse de requisições e respostas)
- ✅ **Services**: Contêm toda a lógica de negócio e validações
- ✅ **Repositories**: Abstração da camada de dados (GORM/SQLite)
- ✅ **Middleware**: Autenticação e logging estruturado

### 2. **Sistema de Autenticação JWT**
- ✅ Registro de novos usuários com validação
- ✅ Login e geração de token JWT (válido por 24h)
- ✅ Hash de senhas com bcrypt
- ✅ Middleware de autenticação para rotas protegidas

### 3. **Gerenciamento de Usuários**
- ✅ Listar todos os usuários (público)
- ✅ Buscar usuário por ID (público)
- ✅ Atualizar próprio perfil (protegido)
- ✅ Deletar própria conta (protegido)

### 4. **Sistema de Posts**
- ✅ Criar posts (protegido)
- ✅ Listar todos os posts (público)
- ✅ Buscar post por ID (público)
- ✅ Listar posts de um usuário específico (público)
- ✅ Atualizar próprio post (protegido)
- ✅ Deletar próprio post (protegido)

### 5. **Autorização por Ownership**
- ✅ Usuários só podem editar sua própria conta
- ✅ Usuários só podem deletar sua própria conta
- ✅ Usuários só podem editar seus próprios posts
- ✅ Usuários só podem deletar seus próprios posts

### 6. **Validações**
- ✅ Email: validação de formato obrigatória
- ✅ Senha: mínimo 6 caracteres
- ✅ Nome: mínimo 3 caracteres
- ✅ Verificação de email duplicado no registro

### 7. **Configuração com Variáveis de Ambiente**
- ✅ Arquivo `.env` centralizado
- ✅ Variáveis: `DATABASE_PATH`, `PORT`, `JWT_SECRET`
- ✅ Suporte a valores padrão

### 8. **Banco de Dados**
- ✅ SQLite com GORM ORM
- ✅ Migrations automáticas
- ✅ Timestamps automáticos (CreatedAt, UpdatedAt)

### 9. **Tratamento de Erros**
- ✅ Respostas de erro estruturadas e consistentes
- ✅ Códigos de erro customizados
- ✅ Mensagens de erro informativas

### 10. **Logging Estruturado**
- ✅ Middleware de logging em todas as requisições
- ✅ Rastreamento de informações de requisição e resposta

## 📁 Estrutura do Projeto

```
awesomeProject/
├── config/
│   └── config.go              # Gerenciamento de variáveis de ambiente
├── db/
│   └── db.go                  # Inicialização do banco de dados
├── errors/
│   ├── errors.go              # Definição de erros customizados
│   └── handler.go             # Tratamento centralizado de erros
├── handlers/
│   ├── auth.go                # Endpoints de autenticação
│   ├── users.go               # Endpoints de usuários
│   └── posts.go               # Endpoints de posts
├── mappers/
│   └── user.go                # Mapeamento de dados de usuário
├── middleware/
│   └── auth.go                # Validação de JWT e logging
├── models/
│   ├── auth.go                # Models de autenticação
│   ├── users.go               # Model User
│   └── posts.go               # Model Post
├── repositories/
│   ├── users.go               # Acesso a dados de usuários
│   └── posts.go               # Acesso a dados de posts
├── routes/
│   └── routes.go              # Definição de rotas
├── services/
│   ├── auth.go                # Lógica de autenticação
│   ├── users.go               # Lógica de negócio de usuários
│   └── posts.go               # Lógica de negócio de posts
├── data/                      # Dados do banco de dados
├── main.go                    # Ponto de entrada da aplicação
├── docker-compose.yml         # Composição de serviços Docker
├── Dockerfile                 # Imagem Docker
├── .env                       # Variáveis de ambiente
└── go.mod                     # Dependências do projeto
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

## 📦 Dependências

```
github.com/gin-gonic/gin v1.10.0           - Framework HTTP
github.com/golang-jwt/jwt/v5 v5.3.1        - Geração e validação de JWT
golang.org/x/crypto v0.50.0                - Hash de senhas com bcrypt
github.com/joho/godotenv v1.5.1            - Carregamento de .env
gorm.io/driver/sqlite v1.6.0               - Driver SQLite para GORM
gorm.io/gorm v1.31.1                       - ORM para banco de dados
```

## ⚙️ Configuração

### Variáveis de Ambiente (.env)
```env
# Caminho do banco de dados SQLite
DATABASE_PATH=./app.db

# Porta do servidor
PORT=:8080

# Chave secreta para assinar tokens JWT (use uma chave forte em produção!)
JWT_SECRET=sua_chave_secreta_muito_segura_aqui_123456

# Driver do banco de dados
DB_DRIVER=sqlite
```

## 🚀 Como Executar

### Localmente (sem Docker)

1. **Clone/Prepare o projeto**
   ```bash
   cd C:\Users\Renan\GolandProjects\awesomeProject
   ```

2. **Configure o arquivo .env**
   ```bash
   # Já foi criado com valores padrão
   # Edite conforme necessário
   ```

3. **Download de dependências**
   ```bash
   go mod download
   ```

4. **Compile o projeto**
   ```bash
   go build -o awesomeProject.exe
   ```

5. **Execute**
   ```bash
   ./awesomeProject.exe
   ```

6. **Servidor iniciará em**
   ```
   http://localhost:8080
   ```

### Com Docker

```bash
docker-compose up --build
```

O servidor iniciará em `http://localhost:8080`

## 🔄 Fluxo de Autenticação

```
1. Usuário faz POST /auth/register ou /auth/login
   ↓
2. Auth Service valida credenciais
   - Valida formato de email
   - Verifica senha mínima (6 caracteres)
   - Hash da senha com bcrypt
   ↓
3. JWT Token é gerado
   - Válido por 24 horas
   ↓
4. Usuário recebe token e o armazena
   ↓
5. Usuário inclui token em requisições protegidas
   - Header: Authorization: Bearer <token>
   ↓
6. Middleware RequireAuth() valida o token
   - Verifica assinatura
   - Valida expiração
   - Extrai ID do usuário
   ↓
7. Se válido, requisição é processada
   Se inválido, erro 401 Unauthorized
```

## 🧪 Testando a API

Use ferramentas como **Postman**, **Insomnia** ou **cURL**:

1. ✅ Registrar novo usuário
2. ✅ Fazer login e obter token
3. ✅ Listar posts públicos (sem autenticação)
4. ✅ Tentar criar post sem token (deve falhar com 401)
5. ✅ Criar post com token válido (deve funcionar)
6. ✅ Listar posts do usuário
7. ✅ Tentar atualizar post de outro usuário (deve falhar com 403)
8. ✅ Atualizar próprio post (deve funcionar)
9. ✅ Deletar próprio post

## 📝 Próximas Melhorias Opcionais

- [ ] Refresh tokens para melhor segurança
- [ ] Rate limiting por IP/usuário
- [ ] Paginação em listagens
- [ ] Soft delete (marcar como deletado ao invés de remover)
- [ ] Testes unitários e de integração
- [ ] Documentação OpenAPI/Swagger
- [ ] Cache com Redis
- [ ] CI/CD com GitHub Actions

---

**Projeto desenvolvido com Go, Gin, JWT e GORM 🚀**
