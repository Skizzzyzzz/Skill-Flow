# SkillFlow API Documentation

## Base URL

```
https://api.skillflow.local/api/v1
```

## Authentication

SkillFlow API uses JWT (JSON Web Tokens) for authentication. Two types of tokens are issued:
- **Access Token**: Short-lived token for API requests (24 hours)
- **Refresh Token**: Long-lived token for obtaining new access tokens (7 days)

### Authentication Header

```
Authorization: Bearer <access_token>
```

## Endpoints

### Authentication

#### Register User

```http
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

#### Login

```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

#### Refresh Token

```http
POST /auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### OIDC Login

```http
GET /auth/oidc/login
```

Redirects to Keycloak for authentication.

#### OIDC Callback

```http
GET /auth/oidc/callback?code=<authorization_code>
```

### Users

#### Get Current User

```http
GET /users/me
```

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "is_active": true,
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z",
  "profile": {
    "first_name": "John",
    "last_name": "Doe",
    "display_name": "johndoe",
    "bio": "Software developer",
    "department": "Engineering",
    "position": "Senior Developer"
  }
}
```

#### Update Current User

```http
PUT /users/me
```

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "username": "newusername"
}
```

#### Get User by ID

```http
GET /users/{id}
```

#### Get User Profile

```http
GET /users/{id}/profile
```

#### Update User Profile

```http
PUT /users/{id}/profile
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "bio": "Software developer passionate about Go",
  "department": "Engineering",
  "position": "Senior Developer",
  "location": "San Francisco, CA",
  "phone": "+1234567890"
}
```

#### Search Users

```http
GET /users/search?q=john
```

### Posts

#### Create Post

```http
POST /posts
```

**Request Body:**
```json
{
  "content": "Hello, SkillFlow! This is my first post.",
  "visibility": "public",
  "media_urls": "[\"https://cdn.skillflow.local/image1.jpg\"]"
}
```

#### Get Feed

```http
GET /posts?page=1&limit=20
```

#### Get Post by ID

```http
GET /posts/{id}
```

#### Update Post

```http
PUT /posts/{id}
```

#### Delete Post

```http
DELETE /posts/{id}
```

#### Get User Posts

```http
GET /posts/user/{user_id}?page=1&limit=20
```

### Comments

#### Create Comment

```http
POST /posts/{id}/comments
```

**Request Body:**
```json
{
  "content": "Great post!",
  "parent_id": null
}
```

#### Get Post Comments

```http
GET /posts/{id}/comments
```

#### Update Comment

```http
PUT /comments/{id}
```

#### Delete Comment

```http
DELETE /comments/{id}
```

### Reactions

#### Add Reaction

```http
POST /posts/{id}/reactions
```

**Request Body:**
```json
{
  "type": "like"
}
```

Types: `like`, `love`, `celebrate`, `support`, `insightful`

#### Remove Reaction

```http
DELETE /posts/{id}/reactions
```

#### Get Reactions

```http
GET /posts/{id}/reactions
```

### Connections

#### Send Connection Request

```http
POST /connections
```

**Request Body:**
```json
{
  "target_id": 123
}
```

#### Get Connections

```http
GET /connections
```

#### Get Pending Requests

```http
GET /connections/pending
```

#### Accept Connection

```http
PUT /connections/{id}/accept
```

#### Reject Connection

```http
PUT /connections/{id}/reject
```

#### Remove Connection

```http
DELETE /connections/{id}
```

### Notifications

#### Get Notifications

```http
GET /notifications?page=1&limit=20
```

#### Mark as Read

```http
PUT /notifications/{id}/read
```

#### Mark All as Read

```http
PUT /notifications/read-all
```

#### Get Unread Count

```http
GET /notifications/unread/count
```

### Messages

#### Send Message

```http
POST /messages
```

**Request Body:**
```json
{
  "receiver_id": 123,
  "content": "Hello! How are you?"
}
```

#### Get Conversations

```http
GET /messages/conversations
```

#### Get Conversation

```http
GET /messages/conversation/{user_id}
```

#### Mark Message as Read

```http
PUT /messages/{id}/read
```

### Groups

#### Create Group

```http
POST /groups
```

**Request Body:**
```json
{
  "name": "Go Developers",
  "description": "A group for Go enthusiasts",
  "visibility": "public"
}
```

#### Get Groups

```http
GET /groups?page=1&limit=20
```

#### Get Group by ID

```http
GET /groups/{id}
```

#### Update Group

```http
PUT /groups/{id}
```

#### Delete Group

```http
DELETE /groups/{id}
```

#### Join Group

```http
POST /groups/{id}/join
```

#### Leave Group

```http
POST /groups/{id}/leave
```

#### Get Group Members

```http
GET /groups/{id}/members
```

#### Get Group Posts

```http
GET /groups/{id}/posts?page=1&limit=20
```

### Skills

#### Get Skills

```http
GET /skills
```

#### Create Skill

```http
POST /skills
```

**Request Body:**
```json
{
  "name": "Golang",
  "category": "Programming Languages",
  "description": "Go programming language"
}
```

#### Add User Skill

```http
POST /skills/user
```

**Request Body:**
```json
{
  "skill_id": 1,
  "level": "advanced",
  "years_of_experience": 5
}
```

Levels: `beginner`, `intermediate`, `advanced`, `expert`

#### Update User Skill

```http
PUT /skills/user/{id}
```

#### Remove User Skill

```http
DELETE /skills/user/{id}
```

#### Endorse Skill

```http
POST /skills/endorse/{user_skill_id}
```

**Request Body:**
```json
{
  "comment": "Great skills in Go!"
}
```

### Files

#### Upload File

```http
POST /files/upload
```

**Request:**
- Content-Type: `multipart/form-data`
- Form field: `file`

**Response:**
```json
{
  "id": 1,
  "name": "document.pdf",
  "size": 1024567,
  "mime_type": "application/pdf",
  "url": "https://cdn.skillflow.local/files/document.pdf",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Get File

```http
GET /files/{id}
```

#### Delete File

```http
DELETE /files/{id}
```

### WebSocket

Real-time communication endpoint.

```
wss://api.skillflow.local/api/v1/ws
```

**Authentication:**
- Send access token in query parameter: `?token=<access_token>`

**Message Format:**
```json
{
  "type": "notification|message|typing",
  "data": {}
}
```

### Admin

Admin endpoints require `admin` role.

#### Get All Users

```http
GET /admin/users
```

#### Activate User

```http
PUT /admin/users/{id}/activate
```

#### Deactivate User

```http
PUT /admin/users/{id}/deactivate
```

#### Delete Post (Admin)

```http
DELETE /admin/posts/{id}
```

#### Delete Comment (Admin)

```http
DELETE /admin/comments/{id}
```

#### Get Statistics

```http
GET /admin/stats
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### Common HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Rate Limiting

API requests are rate-limited to:
- 100 requests per minute (average)
- 50 burst requests

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

## Pagination

List endpoints support pagination:

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)

**Response Headers:**
```
X-Total-Count: 250
X-Page: 1
X-Per-Page: 20
```

## Versioning

API version is included in the URL path: `/api/v1/`

## Support

For API support:
- Documentation: https://docs.skillflow.local
- Issues: https://github.com/vern/skillflow/issues
