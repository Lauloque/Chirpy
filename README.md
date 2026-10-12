# Chirpy

A simple HTTP server made with Go in the context of Boot.dev HTTP server class

[Learn HTTP Servers in Go: Interactive Course | Boot.dev](https://www.boot.dev/courses/learn-http-servers-golang)

<img width="900" height="550" alt="image" src="https://github.com/user-attachments/assets/62f351eb-dd6b-4ddd-90e5-9a0de738e0ad" />

## Requirements

- [Go](https://go.dev/doc/install), version 1.26.4 or compatible with the version declared in `go.mod`

- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

Clone the repository and build the executable:

```sh
git clone https://github.com/Lauloque/Chirpy.gitcd Chirpygo build -o chirpy .
```

## Setup

### 1. Set up PostgreSQL

Install PostgreSQL using the [official installation instructions](https://www.postgresql.org/download/), then start the PostgreSQL server.

Create a database for Chirpy and configure a PostgreSQL connection URL. For example:

```text
postgres://username:password@localhost:5432/chirpy?sslmode=disable
```

Replace `username` and `password` with your PostgreSQL credentials. The database named `chirpy` must exist before running the server. See [PostgreSQL: Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) for more.

### 2. Configure environment variables

Create a `.env` file in the project root:

```dotenv
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="replace-with-a-random-secret"
POLKA_KEY="replace-with-a-random-api-key"
```

> [!WARNING]
> Never commit the `.env` file or share its content: it contains database credentials and secret keys. Hence `.env` is in the `.gitignore` and needs to be generated.

- `DB_URL`: Your PostgreSQL connection URL for the Chirpy database.

- `PLATFORM`: Set to `dev` to enable the database reset endpoint. Any other value disables that endpoint.

- `SECRET`: Secret used to sign and validate authentication tokens. Keep it private and use a sufficiently long, randomly generated value.

- `POLKA_KEY`: Shared secret used to authenticate incoming simulated subscription webhooks.

> [!TIP]
> Generate `SECRET` with `openssl rand -base64 64`; Base64 is a compact format for opaque signing key material. Generate `POLKA_KEY` with `openssl rand -hex 32`; hex uses only letters and digits, making it convenient for environment variables and HTTP headers. Keep both values private.

### 3. Apply database migrations

Apply the SQL migrations in `sql/schema` to create the database tables before starting the server. The server does not apply these migrations automatically.

## Using the Makefile

The `Makefile` loads `DB_URL` from `.env` and provides helpers for SQL generation and database migrations. This way, all you need to do is run `make <command>` in the cli at the project root to perform the following:

```sh
make up         # generate SQL code and apply migrations
make down       # roll back the latest migration
make reset      # roll back and reapply the latest migration
make reset-all  # roll back and reapply all migrations
make build      # build the server executable as ./Chirpy
```

## Running the server

From the project root, run:

```sh
./Chirpy
```

The server listens on port `8080` by default.

The `/app/` route serves the static files in the project directory.

## API

The API uses JSON for request and response bodies unless otherwise specified. Replace placeholder IDs and tokens in the examples with actual values returned by the server.

### Authentication

Authenticated endpoints expect a bearer token in the `Authorization` header:

```http
Authorization: Bearer <token>
```

### Users

#### `POST /api/users`

Create a user.

Request:

```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

Success response (`201 Created`):

```json
{
  "id": "<user-uuid>",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### `PUT /api/users`

Update the authenticated user's email and password.

Requires a valid bearer token.

Request:

```json
{
  "email": "new@example.com",
  "password": "new-password"
}
```

Success response (`200 OK`): the updated user object.

### Authentication endpoints

#### `POST /api/login`

Log in with an existing account.

Request:

```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

Success response (`200 OK`): the user object, an access token, and a refresh token.

#### `POST /api/refresh`

Generate a new access token using a refresh token.

Header:

```http
Authorization: Bearer <refresh-token>
```

Success response (`200 OK`):

```json
{
  "token": "<new-access-token>"
}
```

#### `POST /api/revoke`

Revoke a refresh token.

Header:

```http
Authorization: Bearer <refresh-token>
```

Success response: `204 No Content`.

### Chirps

#### `GET /api/chirps`

List chirps. Authentication is not required.

Optional query parameters:

- `author_id`: UUID of the author whose chirps should be listed.

- `sort`: Set to `desc` to sort by creation date, newest first. The default order is oldest first.

Example:

```http
GET /api/chirps?sort=desc
```

Success response (`200 OK`): an array of chirp objects.

#### `GET /api/chirps/{chirpID}`

Retrieve a chirp by its UUID.

Success response (`200 OK`): a chirp object.

#### `POST /api/chirps`

Create a chirp. Requires a valid bearer token.

Request:

```json
{
  "body": "Hello, Chirpy!"
}
```

Chirps are limited to 140 characters. Certain words are automatically replaced with asterisks.

Success response (`201 Created`): the newly created chirp object.

#### `DELETE /api/chirps/{chirpID}`

Delete a chirp. Requires a valid bearer token and ownership of the chirp.

Success response: `204 No Content`.

### Webhooks

#### `POST /api/polka/webhooks`

Simulate a subscription provider notifying Chirpy that a user has upgraded to Chirpy Red.

Requires a valid API key in the `Authorization` header:

```http
Authorization: ApiKey <polka-key>
```

Request:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "<user-uuid>"
  }
}
```

Success response: `204 No Content`.

### Admin

#### `GET /admin/metrics`

Returns an HTML page displaying the number of visits to the static file server.

Success response: `200 OK`.

#### `POST /admin/reset`

Resets the static file server hit counter and clears the user database.

This endpoint only works when `PLATFORM=dev`. It is destructive and should not be exposed in a production environment.

Success response: `200 OK`.

#### `GET /api/healthz`

Readiness endpoint for checking whether the HTTP server responds.

Success response (`200 OK`):

```text
OK
```
