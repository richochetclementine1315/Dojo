# Dojo Backend API Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Tech Stack](#tech-stack)
3. [Database Schema](#database-schema)
4. [Environment Setup](#environment-setup)
5. [Authentication](#authentication)
6. [API Modules](#api-modules)
   - [Module 1: Auth API](#module-1-auth-api)
   - [Module 2: User API](#module-2-user-api)
7. [Error Handling](#error-handling)
8. [Testing Guide](#testing-guide)

---

## Project Overview

**Dojo** is a collaborative competitive programming platform built with Go and Fiber framework. The platform enables users to:
- Authenticate via email/password or OAuth (Google, GitHub)
- Manage user profiles with platform integrations (LeetCode, Codeforces, CodeChef, GFG)
- Track coding statistics across multiple platforms
- Collaborate in real-time rooms (upcoming)
- Create and share problem sheets (upcoming)
- Participate in contests (upcoming)

**Base URL:** `http://localhost:8080`

---

## Tech Stack

### Backend Framework
- **Go 1.21+** - Programming language
- **Fiber v2.52.0** - Web framework
- **GORM v1.25.5** - ORM for database operations
- **PostgreSQL 15** - Primary database

### Authentication & Security
- **JWT (golang-jwt/jwt v5.2.0)** - Token-based authentication
- **OAuth2 (golang.org/x/oauth2)** - Google and GitHub OAuth
- **bcrypt (golang.org/x/crypto)** - Password hashing
- **Validator v10.16.0** - Request validation

### Utilities
- **UUID (google/uuid v1.5.0)** - Primary key generation
- **godotenv** - Environment variable management

---

## Database Schema

### Core Tables (Implemented)

#### users
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | User unique identifier |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User email address |
| username | VARCHAR(50) | UNIQUE, NOT NULL | User username |
| password_hash | VARCHAR(255) | NULLABLE | Hashed password (null for OAuth-only) |
| avatar_url | VARCHAR(255) | DEFAULT '' | Profile picture URL |
| is_verified | BOOLEAN | DEFAULT false | Email verification status |
| created_at | TIMESTAMP | NOT NULL | Account creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

#### user_profiles
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Profile unique identifier |
| user_id | UUID | FOREIGN KEY, UNIQUE | Reference to users table |
| bio | TEXT | DEFAULT '' | User biography |
| location | VARCHAR(100) | DEFAULT '' | User location |
| website | VARCHAR(255) | DEFAULT '' | Personal website URL |
| leetcode_username | VARCHAR(50) | DEFAULT '' | LeetCode handle |
| codeforces_username | VARCHAR(50) | DEFAULT '' | Codeforces handle |
| codechef_username | VARCHAR(50) | DEFAULT '' | CodeChef handle |
| gfg_username | VARCHAR(50) | DEFAULT '' | GeeksforGeeks handle |
| total_solved | INTEGER | DEFAULT 0 | Total problems solved |
| easy_solved | INTEGER | DEFAULT 0 | Easy problems solved |
| medium_solved | INTEGER | DEFAULT 0 | Medium problems solved |
| hard_solved | INTEGER | DEFAULT 0 | Hard problems solved |
| created_at | TIMESTAMP | NOT NULL | Profile creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

#### auth_accounts
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Auth account identifier |
| user_id | UUID | FOREIGN KEY | Reference to users table |
| provider | VARCHAR(50) | NOT NULL | Auth provider (email/google/github) |
| provider_user_id | VARCHAR(255) | NOT NULL | Provider-specific user ID |
| access_token | TEXT | NULLABLE | OAuth access token |
| refresh_token | TEXT | NULLABLE | OAuth refresh token |
| expires_at | TIMESTAMP | NULLABLE | Token expiration time |
| created_at | TIMESTAMP | NOT NULL | Account creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

**Unique Constraint:** `(provider, provider_user_id)`

#### refresh_tokens
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Token identifier |
| user_id | UUID | FOREIGN KEY | Reference to users table |
| token | VARCHAR(255) | UNIQUE, NOT NULL | JWT refresh token |
| expires_at | TIMESTAMP | NOT NULL | Token expiration time |
| created_at | TIMESTAMP | NOT NULL | Token creation time |

#### user_platform_stats
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Stat identifier |
| user_id | UUID | FOREIGN KEY | Reference to users table |
| platform | VARCHAR(50) | NOT NULL | Platform name (leetcode/codeforces/etc) |
| rating | INTEGER | NULLABLE | Current rating |
| max_rating | INTEGER | NULLABLE | Maximum rating achieved |
| problems_solved | INTEGER | DEFAULT 0 | Total problems solved |
| easy_problems_solved | INTEGER | DEFAULT 0 | Easy problems solved |
| med_problems_solved | INTEGER | DEFAULT 0 | Medium problems solved |
| hard_problems_solved | INTEGER | DEFAULT 0 | Hard problems solved |
| contests_attended | INTEGER | DEFAULT 0 | Contests participated |
| global_rank | INTEGER | NULLABLE | Global ranking |
| last_synced | TIMESTAMP | NOT NULL | Last sync time |
| created_at | TIMESTAMP | NOT NULL | Record creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

**Unique Constraint:** `(user_id, platform)`

---

## Environment Setup

### Required Environment Variables (.env)

```env
# Server Configuration
PORT=8080
FRONTEND_URL=http://localhost:5173

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=DojoDB

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_min_32_chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/auth/github/callback
```

### Running the Server

```bash
# Navigate to backend directory
cd Backend

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

Server will start on `http://localhost:8080`

---

## Authentication

### JWT Token System

The API uses JWT-based authentication with two token types:

1. **Access Token**
   - Expires in 15 minutes
   - Used for API authorization
   - Sent in `Authorization` header as `Bearer <token>`

2. **Refresh Token**
   - Expires in 7 days
   - Used to generate new access tokens
   - Stored in database with user association

### Protected Routes

Protected routes require the `Authorization` header:

```
Authorization: Bearer <access_token>
```

Example:
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## API Modules

## Module 1: Auth API

### 1.1 Register (Email/Password)

**Endpoint:** `POST /api/auth/register`

**Description:** Create a new user account with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123",
  "leetcode_username": "john_leetcode"
}
```

**Validation Rules:**
- `email`: Required, valid email format
- `username`: Required, 3-50 characters
- `password`: Required, minimum 8 characters
- `leetcode_username`: Optional, 3-50 characters

**Success Response (201):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Error Responses:**

- **400 Bad Request** - Validation error
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "email is required"
}
```

- **409 Conflict** - Email or username already exists
```json
{
  "success": false,
  "message": "Email already in use",
  "error": "user already exists"
}
```

---

### 1.2 Login (Email/Password)

**Endpoint:** `POST /api/auth/login`

**Description:** Authenticate user with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Error Responses:**

- **401 Unauthorized** - Invalid credentials
```json
{
  "success": false,
  "message": "Invalid credentials",
  "error": "incorrect password"
}
```

- **404 Not Found** - User not found
```json
{
  "success": false,
  "message": "User not found",
  "error": "no user with this email"
}
```

---

### 1.3 Google OAuth Login

**Endpoint:** `GET /api/auth/google/login`

**Description:** Redirect to Google OAuth consent screen.

**Response:** HTTP 302 Redirect to Google authorization URL

**Usage:**
```html
<a href="http://localhost:8080/api/auth/google/login">Login with Google</a>
```

---

### 1.4 Google OAuth Callback

**Endpoint:** `GET /api/auth/google/callback`

**Description:** Handle Google OAuth callback and create/login user.

**Query Parameters:**
- `code`: Authorization code from Google
- `state`: CSRF protection token

**Success Response (200):**
```json
{
  "success": true,
  "message": "Google authentication successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

---

### 1.5 GitHub OAuth Login

**Endpoint:** `GET /api/auth/github/login`

**Description:** Redirect to GitHub OAuth consent screen.

**Response:** HTTP 302 Redirect to GitHub authorization URL

**Usage:**
```html
<a href="http://localhost:8080/api/auth/github/login">Login with GitHub</a>
```

---

### 1.6 GitHub OAuth Callback

**Endpoint:** `GET /api/auth/github/callback`

**Description:** Handle GitHub OAuth callback and create/login user.

**Query Parameters:**
- `code`: Authorization code from GitHub
- `state`: CSRF protection token

**Success Response (200):**
```json
{
  "success": true,
  "message": "GitHub authentication successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

---

### 1.7 Refresh Access Token

**Endpoint:** `POST /api/auth/refresh`

**Description:** Generate new access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Error Responses:**

- **401 Unauthorized** - Invalid or expired token
```json
{
  "success": false,
  "message": "Invalid refresh token",
  "error": "token expired"
}
```

---

### 1.8 Logout

**Endpoint:** `POST /api/auth/logout`

**Description:** Invalidate refresh token and logout user.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

---

## Module 2: User API

**All User API endpoints require authentication.**

### 2.1 Get User Profile

**Endpoint:** `GET /api/users/profile`

**Description:** Retrieve authenticated user's complete profile with platform stats.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "avatar_url": "https://example.com/avatar.jpg",
    "is_verified": false,
    "created_at": "2026-01-24T10:30:00Z",
    "profile": {
      "bio": "Competitive programmer and software engineer",
      "location": "San Francisco, CA",
      "website": "https://johndoe.dev",
      "leetcode_username": "john_leetcode",
      "codeforces_username": "john_cf",
      "codechef_username": "john_cc",
      "gfg_username": "john_gfg",
      "total_solved": 450,
      "easy_solved": 180,
      "medium_solved": 200,
      "hard_solved": 70,
      "platform_stats": [
        {
          "platform": "leetcode",
          "rating": 1850,
          "max_rating": 1900,
          "solved_count": 320,
          "contest_rating": 0,
          "global_rank": 12543,
          "last_synced_at": "2026-01-24T09:15:00Z"
        },
        {
          "platform": "codeforces",
          "rating": 1654,
          "max_rating": 1720,
          "solved_count": 215,
          "contest_rating": 1654,
          "global_rank": 45231,
          "last_synced_at": "2026-01-24T09:20:00Z"
        }
      ]
    }
  }
}
```

**Error Responses:**

- **401 Unauthorized** - Missing or invalid token
- **404 Not Found** - User not found

---

### 2.2 Update User Profile

**Endpoint:** `PUT /api/users/profile`

**Description:** Update user profile information (bio, location, platform usernames).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "bio": "Updated bio text",
  "location": "New York, NY",
  "website": "https://mynewsite.com",
  "leetcode_username": "new_leetcode_handle",
  "codeforces_username": "new_cf_handle",
  "codechef_username": "new_cc_handle",
  "gfg_username": "new_gfg_handle"
}
```

**Note:** All fields are optional. Only send fields you want to update.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "avatar_url": "https://example.com/avatar.jpg",
    "is_verified": false,
    "created_at": "2026-01-24T10:30:00Z",
    "profile": {
      "bio": "Updated bio text",
      "location": "New York, NY",
      "website": "https://mynewsite.com",
      "leetcode_username": "new_leetcode_handle",
      "codeforces_username": "new_cf_handle",
      "codechef_username": "new_cc_handle",
      "gfg_username": "new_gfg_handle",
      "total_solved": 450,
      "easy_solved": 180,
      "medium_solved": 200,
      "hard_solved": 70,
      "platform_stats": []
    }
  }
}
```

---

### 2.3 Update User Account

**Endpoint:** `PATCH /api/users/account`

**Description:** Update user account details (username, avatar).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "username": "newusername",
  "avatar_url": "https://example.com/new-avatar.jpg"
}
```

**Validation Rules:**
- `username`: Optional, 3-50 characters, must be unique
- `avatar_url`: Optional, valid URL format

**Success Response (200):**
```json
{
  "success": true,
  "message": "Account updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "newusername",
    "avatar_url": "https://example.com/new-avatar.jpg",
    "is_verified": false,
    "created_at": "2026-01-24T10:30:00Z",
    "profile": { ... }
  }
}
```

**Error Responses:**

- **409 Conflict** - Username already taken
```json
{
  "success": false,
  "message": "Username already taken",
  "error": "username exists"
}
```

---

### 2.4 Change Password

**Endpoint:** `POST /api/users/change-password`

**Description:** Change user password (not available for OAuth-only accounts).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "old_password": "OldSecurePass123",
  "new_password": "NewSecurePass456"
}
```

**Validation Rules:**
- `old_password`: Required
- `new_password`: Required, minimum 8 characters

**Success Response (200):**
```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

**Error Responses:**

- **401 Unauthorized** - Invalid old password
```json
{
  "success": false,
  "message": "Invalid old password",
  "error": "password mismatch"
}
```

- **400 Bad Request** - OAuth-only account
```json
{
  "success": false,
  "message": "Cannot change password for OAuth-only accounts",
  "error": "no password set"
}
```

---

### 2.5 Sync Platform Statistics

**Endpoint:** `POST /api/users/sync-stats`

**Description:** Sync coding statistics from external platforms (LeetCode, Codeforces, CodeChef, GFG).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "platforms": ["leetcode", "codeforces", "codechef", "gfg"]
}
```

**Validation Rules:**
- `platforms`: Required, array of valid platform names
- Valid platforms: `leetcode`, `codeforces`, `codechef`, `gfg`

**Success Response (200):**
```json
{
  "success": true,
  "message": "Platform stats sync completed",
  "data": {
    "leetcode": {
      "status": "success"
    },
    "codeforces": {
      "status": "success"
    },
    "codechef": {
      "error": "CodeChef username not set"
    },
    "gfg": {
      "status": "success"
    }
  }
}
```

**Note:** Each platform sync is independent. Some may succeed while others fail. Check individual platform status in response.

**Common Errors Per Platform:**
- Username not set in profile
- Platform API unavailable
- Invalid username
- Rate limiting

---

## Error Handling

### Standard Error Response Format

All errors follow this consistent format:

```json
{
  "success": false,
  "message": "Human-readable error message",
  "error": "Technical error details"
}
```

### HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, PUT, PATCH, POST |
| 201 | Created | Successful resource creation (register) |
| 400 | Bad Request | Validation errors, malformed request |
| 401 | Unauthorized | Missing/invalid token, wrong password |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate email/username |
| 500 | Internal Server Error | Server-side errors |

### Common Error Scenarios

#### 1. Missing Authorization Header
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "missing or malformed token"
}
```

#### 2. Expired Access Token
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "token expired"
}
```

#### 3. Validation Error
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "email: must be a valid email address; password: must be at least 8 characters"
}
```

---

## Testing Guide

### Prerequisites

- **Postman** or **Thunder Client** (VS Code extension)
- Backend server running on `http://localhost:8080`
- PostgreSQL database running

### Testing Workflow

#### Step 1: Register a New User

```bash
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "email": "test@example.com",
  "username": "testuser",
  "password": "password123",
  "leetcode_username": "test_leetcode"
}
```

Save the `access_token` and `refresh_token` from the response.

#### Step 2: Test Protected Route (Get Profile)

```bash
GET http://localhost:8080/api/users/profile
Authorization: Bearer <your_access_token>
```

#### Step 3: Update Profile

```bash
PUT http://localhost:8080/api/users/profile
Authorization: Bearer <your_access_token>
Content-Type: application/json

{
  "bio": "Test bio",
  "location": "Test City",
  "leetcode_username": "my_leetcode"
}
```

#### Step 4: Sync Platform Stats

```bash
POST http://localhost:8080/api/users/sync-stats
Authorization: Bearer <your_access_token>
Content-Type: application/json

{
  "platforms": ["leetcode", "codeforces"]
}
```

#### Step 5: Change Password

```bash
POST http://localhost:8080/api/users/change-password
Authorization: Bearer <your_access_token>
Content-Type: application/json

{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

#### Step 6: Refresh Token

```bash
POST http://localhost:8080/api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "<your_refresh_token>"
}
```

#### Step 7: Logout

```bash
POST http://localhost:8080/api/auth/logout
Authorization: Bearer <your_access_token>
Content-Type: application/json

{
  "refresh_token": "<your_refresh_token>"
}
```

### OAuth Testing

#### Google OAuth
1. Open browser: `http://localhost:8080/api/auth/google/login`
2. Complete Google authentication
3. You'll receive tokens in the callback response

#### GitHub OAuth
1. Open browser: `http://localhost:8080/api/auth/github/login`
2. Complete GitHub authentication
3. You'll receive tokens in the callback response

---

## Project Structure

```
Backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loader
│   ├── dto/
│   │   ├── auth_dto.go            # Auth request/response DTOs
│   │   └── user_dto.go            # User request/response DTOs
│   ├── handler/
│   │   ├── auth_handler.go        # Auth HTTP handlers
│   │   └── user_handler.go        # User HTTP handlers
│   ├── middleware/
│   │   ├── auth.go                # JWT authentication middleware
│   │   ├── cors.go                # CORS middleware
│   │   ├── error_handler.go       # Global error handler
│   │   ├── logger.go              # Request logging
│   │   └── rate_limit.go          # Rate limiting
│   ├── models/
│   │   ├── user.go                # User model
│   │   ├── auth.go                # Auth models
│   │   └── ... (other models)
│   ├── repository/
│   │   ├── user_repository.go     # User database operations
│   │   └── auth_repository.go     # Auth database operations
│   ├── routes/
│   │   └── routes.go              # Route definitions
│   └── service/
│       ├── auth_service.go        # Auth business logic
│       └── user_service.go        # User business logic
├── pkg/
│   ├── oauth/
│   │   ├── google.go              # Google OAuth integration
│   │   └── github.go              # GitHub OAuth integration
│   ├── scraper/
│   │   └── platform_scraper.go    # Platform stats scrapers
│   └── utils/
│       ├── jwt.go                 # JWT utilities
│       ├── password.go            # Password hashing
│       ├── validator.go           # Request validation
│       ├── response.go            # Standard responses
│       └── errors.go              # Custom errors
├── .env                           # Environment variables
├── go.mod                         # Go module definition
└── go.sum                         # Go dependencies checksum
```

---

## Upcoming Modules

### Module 3: Social API
- Friend requests (send, accept, reject)
- Friend list management
- Block/unblock users
- Notifications system

### Module 4: Problem API
- Fetch problems from platforms
- Search and filter problems
- Personal notes on problems
- Problem difficulty tracking

### Module 5: Sheet API
- Create problem sheets
- Add/remove problems from sheets
- Share sheets with others
- Track sheet completion

### Module 6: Contest API
- Fetch upcoming contests
- Set contest reminders
- Track contest participation
- Contest history

### Module 7: Room API
- Create collaborative rooms (max 4 users)
- Invite participants
- Room settings and permissions
- Real-time participant management

### Module 8: WebSocket API
- Real-time code synchronization
- Collaborative whiteboard
- Video chat signaling
- Live cursor positions

---

## Notes

- All timestamps are in UTC format (ISO 8601)
- UUIDs are in standard format: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
- All API responses follow the standard format with `success`, `message`, and `data` fields
- Rate limiting is enforced: 100 requests per minute per IP
- CORS is enabled for `http://localhost:5173` (frontend)

---

**Last Updated:** January 24, 2026  
**Version:** 1.0.0  
**Status:** Modules 1 & 2 Complete
