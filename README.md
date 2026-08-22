# 🚦 Tollbooth

> A fast, self-hostable URL shortener built with Go and MongoDB, with a modern Next.js frontend.

<!-- TODO: replace with your actual badge links or delete any you don't want -->
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-000000?style=flat&logo=next.js&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat&logo=typescript&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=flat&logo=mongodb&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue)

Tollbooth takes long, unwieldy URLs and returns short, shareable links, then tracks and manages them behind an authenticated API. The backend is a Go HTTP service backed by MongoDB Atlas; the frontend is a Next.js (App Router) single-page dashboard built with TypeScript, RTK Query, and Tailwind CSS.

## 🔗 Live Demo

- **Frontend:** https://www.toll-booth-ui.vercel.app
- **API:** https://toll-booth-production.up.railway.app/

## ✨ Features

- **Shorten URLs** — turn any long link into a compact short code
- **Redirects** — short codes resolve to their original destination
- **Authentication** — Basic auth and JWT-based auth for protected routes
- **Per-user link management** — create, list, and manage your own links
- **CORS-enabled API** — configured to talk to the Next.js frontend across origins
- **Cloud-ready** — deployable to Render / Railway with environment-based config

## 🛠️ Tech Stack

**Backend**
- Go (standard `net/http` server)
- MongoDB Atlas (official Go driver, BSON struct tags)
- Basic authentication
- Custom CORS middleware
- `godotenv` for local environment loading

**Frontend**
- Next.js (App Router) + TypeScript
- Redux Toolkit + RTK Query for data fetching and cache
- Tailwind CSS

**Infrastructure**
- MongoDB Atlas (managed database)
- Render / Railway (server hosting)
- Vercel (front end hosting)

## 📂 Project Structure

The backend follows standard Go project layout conventions:

```
tollbooth/
├── cmd/
│   └── server/          # main entrypoint
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # CORS, auth
│   ├── models/          # data models with BSON tags
│   └── db/              # MongoDB connection & queries
├── go.mod
└── go.sum
```


## 🚀 Getting Started

### Prerequisites

- Go 1.26+
- Node.js 18+
- A MongoDB Atlas cluster (or local MongoDB instance)

### Backend

```bash
# clone
git clone https://github.com/EggBurner/toll-booth.git
cd tollbooth

# set up environment (see below)
cp .env.example .env

# run
go run ./cmd/server
```

The API will start on `http://localhost:8080` <!-- TODO: your actual port -->.

**Environment variables** (`.env`):

```env
MONGODB_URI=your-atlas-connection-string
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
```
<!-- TODO: match these names to what your code actually reads -->

### Frontend

```bash
cd frontend   # or wherever your Next.js app lives
npm install
npm run dev
```

Set the API base URL for the frontend:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Then open `http://localhost:3000`.

## 📡 API Reference

<!-- TODO: adjust routes/methods to match your actual handlers -->

| Method | Endpoint            | Auth      | Description                    |
|--------|---------------------|-----------|--------------------------------|
| POST   | `/register`         | None      | Create a new account           |
| POST   | `/login`            | None      | Authenticate, returns a JWT    |
| POST   | `/shorten`          | JWT       | Create a short URL             |
| GET    | `/{shortCode}`      | None      | Redirect to the original URL   |
| GET    | `/urls`             | JWT       | List the current user's URLs   |

Example — shorten a URL:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/some/very/long/path"}'
```

## ☁️ Deployment

The service is designed to run on platforms like Render or Railway. A few notes from getting it live:

- Provide all environment variables through the platform dashboard rather than committing a `.env` file.
- Point `ALLOWED_ORIGINS` at your deployed frontend URL so CORS passes in production.
- MongoDB Atlas needs the hosting platform's outbound IPs allowlisted (or `0.0.0.0/0` for simplicity during testing).

## 🗺️ Roadmap

<!-- TODO: keep, edit, or delete -->
- [ ] Click analytics per short link
- [ ] Custom / vanity short codes
- [ ] Link expiration
- [ ] Rate limiting

## 📄 License

Released under the MIT License. <!-- TODO: confirm or change; add a LICENSE file -->
