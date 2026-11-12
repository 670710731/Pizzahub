# PizzaHub

Pizza review app with a React/Vite frontend and Go/PostgreSQL API.

- [Frontend local setup and Vercel deployment](./frontend/README.md)
- [Backend API, Supabase setup, and Render deployment](./backend/README.md)

## Deploy overview

1. Apply `backend/schema.sql` to the Supabase PostgreSQL database once.
2. Deploy `render.yaml` from Render and configure `DATABASE_URL` with the
   Supabase transaction-pooler URL (`sslmode=require`) and `ALLOWED_ORIGINS`
   with the Vercel site origin.
3. Deploy the Vercel project with Root Directory `frontend` and set
   `VITE_API_URL` to the Render API URL ending in `/api/v1`.
4. Register the initial account, then promote it to admin in Supabase SQL
   Editor as documented in [backend/README.md](./backend/README.md).