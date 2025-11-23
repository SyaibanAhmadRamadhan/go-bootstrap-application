# Authentication & User Management Sample

Sample implementasi lengkap untuk authentication dan user management menggunakan go-bootstrap architecture.

## 📁 Struktur yang Dibuat

### 1. **Domain Layer** (`internal/domain/`)

#### **Auth Domain** (`internal/domain/auth/`)

- `value_object.go` - Token types, user roles, token status
- `dto.go` - Input/Output structs untuk semua operations
- `repository.go` - Repository interfaces untuk token management
- `service.go` - Service interfaces untuk authentication

**Features:**

- Login (email/password)
- Refresh Token
- Logout
- Token Validation
- Token Revocation

#### **User Domain** (`internal/domain/user/`)

- `value_object.go` - User status, roles, gender enums
- `dto.go` - Input/Output structs untuk user operations
- `repository.go` - Repository interfaces untuk user management
- `service.go` - Service interfaces untuk user operations

**Features:**

- User Registration
- Get Profile
- Get List Users (with pagination)
- Update Profile
- Change Password
- Update Status (admin only)

---

### 2. **Module Layer** (`internal/module/`)

#### **Auth Repository** (`internal/module/auth/repository/`)

- `repo_auth.go` - Repository constructor
- `repo_auth_datastore.go` - Database operations implementation
  - CreateToken
  - GetDetailToken
  - RevokeToken
  - DeleteExpiredTokens
- `repo_auth_datastore_test.go` - Unit tests

#### **Auth Service** (`internal/module/auth/service/`)

- `service_auth.go` - Business logic implementation
  - Login with password verification
  - Token generation (access & refresh)
  - Token refresh with validation
  - Logout with token revocation
  - Token validation
- `service_auth_test.go` - Unit tests

#### **User Repository** (`internal/module/user/repository/`)

- `repo_user.go` - Repository constructor
- `repo_user_datastore.go` - Database operations implementation
  - CreateUser
  - GetDetailUser
  - GetListUser (with pagination & filters)
  - UpdateUser
  - UpdatePassword
  - UpdateStatus
- `repo_user_datastore_test.go` - Unit tests

---

### 3. **Worker Layer** (`internal/worker/auth/`)

- `scheduler_token_cleanup.go` - Automated token cleanup
  - Deletes expired tokens older than 24 hours
  - Runs on schedule (configurable via cron)
  - Logs cleanup activities

---

### 4. **API Contracts**

#### **gRPC Proto** (`api/proto/`)

**Auth Proto** (`api/proto/auth/auth.proto`)

```protobuf
service AuthService {
  rpc ApiV1Login(ApiV1LoginRequest) returns (ApiV1LoginResponse)
  rpc ApiV1RefreshToken(ApiV1RefreshTokenRequest) returns (ApiV1RefreshTokenResponse)
  rpc ApiV1Logout(ApiV1LogoutRequest) returns (ApiV1LogoutResponse)
  rpc ApiV1ValidateToken(ApiV1ValidateTokenRequest) returns (ApiV1ValidateTokenResponse)
}
```

**User Proto** (`api/proto/user/user.proto`)

```protobuf
service UserService {
  rpc ApiV1Register(ApiV1RegisterRequest) returns (ApiV1RegisterResponse)
  rpc ApiV1GetProfile(ApiV1GetProfileRequest) returns (ApiV1GetProfileResponse)
  rpc ApiV1GetListUsers(ApiV1GetListUsersRequest) returns (ApiV1GetListUsersResponse)
  rpc ApiV1UpdateProfile(ApiV1UpdateProfileRequest) returns (ApiV1UpdateProfileResponse)
  rpc ApiV1ChangePassword(ApiV1ChangePasswordRequest) returns (ApiV1ChangePasswordResponse)
  rpc ApiV1UpdateStatus(ApiV1UpdateStatusRequest) returns (ApiV1UpdateStatusResponse)
}
```

#### **REST API OpenAPI** (`api/openapi/api.yaml`)

**Auth Endpoints:**

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - User logout

**User Endpoints:**

- `POST /api/v1/users/register` - Register new user
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `GET /api/v1/users` - Get list of users (admin)
- `POST /api/v1/users/change-password` - Change password
- `PUT /api/v1/users/{user_id}/status` - Update user status (admin)

---

### 5. **Database Schema** (`migrations/`)

**Users Table:**

```sql
- id (bigint, PK)
- email (varchar, unique)
- password_hash (varchar)
- name (varchar)
- role (enum: admin, user)
- status (enum: active, inactive, suspended)
- phone (varchar, nullable)
- gender (enum: male, female, other, nullable)
- created_at, updated_at (timestamp)
```

**Auth Tokens Table:**

```sql
- id (bigint, PK)
- user_id (bigint, FK)
- token (text)
- token_type (enum: access, refresh)
- status (enum: active, revoked, expired)
- expires_at (timestamp)
- refresh_token (text, nullable)
- created_at, updated_at (timestamp)
```

---

## 🎯 Naming Convention yang Diterapkan

### Repository Layer

- **Input (write):** `Create<Entity>Params`
- **Input (read):** `<MethodName><Entity>Filters`
- **Output:** `<MethodName><Entity>Result`

### Service Layer

- **Input:** `<MethodName>Input`
- **Output:** `<MethodName>Output`

**Prinsip:** Setiap method memiliki struct param sendiri yang independent.

---

## 🚀 Implementation Status

### ✅ Completed

1. **Domain Layer** - Auth & User interfaces dan DTOs
2. **Repository Layer** - Auth & User repository implementations dengan Squirrel query builder
3. **Service Layer** - Auth service implementation (User service perlu dilengkapi)
4. **Worker Layer** - Token cleanup scheduler
5. **Transport Layer Skeleton** - REST API handlers di `internal/transport/auth/` dan `internal/transport/user/` (dengan TODO comments)
6. **Database Schema** - Migration file tersedia di `migrations/001_create_auth_and_user_tables.sql`

### 🔄 In Progress / TODO

1. **User Service Implementation** - Complete remaining business logic di `internal/module/user/service/`
2. **Transport Layer Implementation** - Implement handlers di:
   - `internal/transport/auth/restapi_auth.go` (Login, Logout, Refresh)
   - `internal/transport/user/restapi_user.go` (Register, Profile, List, Change Password, Update Status)
3. **Middleware** - JWT authentication middleware
4. **Code Generation** - Generate gRPC & OpenAPI code:

   ```bash
   make generate
   ```

5. **Integration Tests** - End-to-end API tests
6. **Database Migrations** - Run migrations:

   ```bash
   mysql -u root -p database_name < migrations/001_create_auth_and_user_tables.sql
   ```

---

## 💡 Key Features

✅ **Complete Authentication Flow**

- Login with email/password
- JWT-like token management (custom tokens)
- Refresh token rotation
- Token validation for inter-service calls
- Secure logout with token revocation

✅ **User Management**

- Registration with password hashing (bcrypt)
- Profile management
- Password change with verification
- User status management (admin)
- Paginated user listing with filters

✅ **Security Best Practices**

- Password hashing with bcrypt
- Token expiration (15 min for access, 7 days for refresh)
- Token revocation on logout
- Automatic cleanup of expired tokens
- Status-based access control

✅ **Observability**

- Structured logging with slog
- Cleanup worker with monitoring
- Error tracking

---

## 📝 Usage Example

### 1. Login

```bash
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

### 2. Get Profile (with token)

```bash
GET /api/v1/users/profile
Authorization: Bearer {access_token}

Response:
{
  "user": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "status": "active"
  }
}
```

### 3. Register New User

```bash
POST /api/v1/users/register
{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "New User"
}
```

---

## 🔧 Configuration

### Token Expiration (customizable in service)

- Access Token: 15 minutes
- Refresh Token: 7 days

### Worker Schedule (untuk cleanup)

Tambahkan di scheduler configuration:

```go
cron.AddFunc("0 0 * * *", tokenCleanupWorker.CleanupExpiredTokens) // Daily at midnight
```

---

Sample ini mendemonstrasikan:

- ✅ Clean Architecture dengan strict naming conventions
- ✅ Separation of Concerns (Auth vs User domain)
- ✅ Repository pattern dengan interface
- ✅ Service layer dengan business logic
- ✅ API contracts (Proto & OpenAPI)
- ✅ Database schema design
- ✅ Background workers
- ✅ Security best practices
