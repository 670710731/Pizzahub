# Pizzahub backend

## Database

The schema targets PostgreSQL. From `backend`, run
`psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f schema.sql` against a new database.

Set `DATABASE_URL`, `JWT_SECRET` (at least 32 characters of random data),
`HTTP_ADDR` (optional; defaults to `:8080`), and
`ALLOWED_ORIGINS` (optional comma-separated exact origins) before starting
with `go run .`.

The schema keeps a pizza's page content (`pizzas`) separate from its categories,
gallery images, and purchasable size/price options. Reviews belong to both a
registered user and a pizza; each user can leave one review per pizza and edit
it later. Only published pizzas and reviews should be exposed on public pages.
The application should calculate a pizza's average rating and review count from
published reviews rather than storing duplicated aggregate values.

Image columns store URLs, not binary files. The application can later point
these URLs at whichever object storage or image service is selected.

## Initial API contract

Use JSON and version public routes under `/api/v1`.

| Method | Route | Purpose |
| --- | --- | --- |
| `GET` | `/api/v1/pizzas` | List published pizzas; support `page`, `page_size`, `category`, and `sort` filters |
| `GET` | `/api/v1/pizzas/{slug}` | Pizza details, categories, images, available variants, rating summary |
| `GET` | `/api/v1/pizzas/{slug}/reviews` | Published reviews, paginated |
| `POST` | `/api/v1/pizzas/{slug}/reviews` | Create or update the signed-in user's review |
| `POST` | `/api/v1/auth/register` | Register a customer |
| `POST` | `/api/v1/auth/login` | Authenticate and issue an access token |
| `GET` | `/api/v1/categories` | List categories |
| `GET/POST` | `/api/v1/admin/categories` | List or create categories (admin) |
| `PUT/DELETE` | `/api/v1/admin/categories/{id}` | Replace or delete a category (admin) |
| `POST/PUT/DELETE` | `/api/v1/admin/pizzas[/{id}]` | Create, replace, or archive a pizza with category IDs, images, and variants (admin) |
| `GET` | `/api/v1/admin/reviews?status=pending` | List reviews for moderation (admin) |
| `PATCH` | `/api/v1/admin/reviews/{id}` | Set review status to `pending`, `published`, or `rejected` (admin) |

List responses should include pagination metadata. Public pizza details should
include a `rating_average` and `review_count`; review submission should accept a
rating from 1 to 5 and must take the user identity from authentication, never
from the request body. Admin routes require an admin role.

## Conventions

- Prices are exact decimal amounts in the currency specified by each variant
  (`THB` by default); do not represent money as floating-point values in Go.
- Use parameterized SQL for all database queries.
- Store password hashes, never plaintext passwords.
- Validate pagination and allowed sort values at the API boundary.
- Use `updated_at` from the application on updates, or add a shared database
  trigger if timestamp management is moved into PostgreSQL.
- Registration always creates a `customer`. Promote the initial administrator
  directly in PostgreSQL after registering: `UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';`
- `PUT` replaces all pizza fields and associated categories, images, and variants.
  `DELETE` archives the pizza rather than erasing its review history.
- Access tokens are signed, stateless bearer tokens and expire after 24 hours.
  Changing a user's role takes effect for new tokens; use short token lifetimes
  or add token revocation if immediate revocation becomes a requirement.
