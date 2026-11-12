# PizzaHub frontend

## Local development

```sh
npm ci
npm run dev
```

By default, Vite proxies `/api/v1` to `http://localhost:8080`. Set
`API_PROXY_TARGET` if the local Go API listens elsewhere. For deployment, set
`VITE_API_URL` to the full backend API base URL, ending in `/api/v1` (no
trailing slash).

## Checks

```sh
npm test
npm run lint
npm run build
```

## Vercel

Set the Vercel project Root Directory to `frontend`, use the Vite preset,
`npm run build` as the build command, and `dist` as the output directory.
Configure `VITE_API_URL` as `https://YOUR-RENDER-SERVICE.onrender.com/api/v1`.
`vercel.json` sends client-side routes to the SPA entry point.