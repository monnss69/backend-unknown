# Component Store Backend

This Go service exposes a REST API for saving and retrieving self-contained React components.

## API

- `POST /components` – create a component with fields `name` and `code`. Returns the stored component.
- `GET /components` – list all components.
- `GET /components/{id}` – fetch a component by id.

## Development

1. Set `DATABASE_URL` to a Postgres connection string.
2. Run database migrations from `db/schema.sql`.
3. Start the server:
   ```bash
   go run ./cmd/server
   ```

The server avoids executing user code and enforces simple validation rules.

