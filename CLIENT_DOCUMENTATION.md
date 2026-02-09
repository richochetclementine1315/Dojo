# Dojo Client (Frontend) Documentation

<div align="center">
  <img src="https://github.com/user-attachments/assets/a6d002d6-c132-4c2c-99d3-5ab05f197116" alt="Dojo Logo" width="200"/>

<p align="center">
  <img src="https://img.shields.io/badge/React-19.2.0-61DAFB?logo=react" alt="React" />
  <img src="https://img.shields.io/badge/Vite-6.0.11-646CFF?logo=vite" alt="Vite" />
  <img src="https://img.shields.io/badge/TypeScript-5.9.3-3178C6?logo=typescript" alt="TypeScript" />
  <img src="https://img.shields.io/badge/TailwindCSS-3.4.19-38BDF8?logo=tailwindcss" alt="TailwindCSS" />
  <img src="https://img.shields.io/badge/Zustand-5.0.11-ff9100?logo=react" alt="Zustand" />
  <img src="https://img.shields.io/badge/Axios-1.13.4-5A29E4?logo=axios" alt="Axios" />
  <img src="https://img.shields.io/badge/React%20Router-7.13.0-CA4245?logo=react-router" alt="React Router" />
  <img src="https://img.shields.io/badge/Framer%20Motion-12.31.0-EF008F?logo=framer" alt="Framer Motion" />
  <img src="https://img.shields.io/badge/Three.js-0.182.0-000000?logo=three.js" alt="Three.js" />
  <img src="https://img.shields.io/badge/Lucide-0.563.0-000000?logo=lucide" alt="Lucide" />
  <img src="https://img.shields.io/badge/React%20Query-5.90.20-FF4154?logo=react-query" alt="React Query" />
  <img src="https://img.shields.io/badge/Vercel-Deploy-000000?logo=vercel" alt="Vercel" />
</p>

---

## Table of Contents
1. [Project Overview](#project-overview)
2. [Tech Stack](#tech-stack)
3. [Project Structure](#project-structure)
4. [Environment Setup](#environment-setup)
5. [Key Features](#key-features)
6. [Pages & Components](#pages--components)
7. [API Integration](#api-integration)
8. [State Management](#state-management)
9. [Styling](#styling)
10. [Testing & Linting](#testing--linting)
11. [Deployment](#deployment)

---

## Project Overview

**Dojo** is a modern, interactive competitive programming platform frontend, built with React, Vite, and TypeScript. It provides a seamless user experience for:
- Authentication (email/password, Google, GitHub)
- Profile management and platform integration
- Problem search, filter, and solve tracking
- Contest calendar and reminders
- Collaborative rooms and real-time features
- Problem sheets and progress tracking

---

## Tech Stack

- **React 19** ![React](https://img.shields.io/badge/-React-61DAFB?logo=react)
- **Vite** ![Vite](https://img.shields.io/badge/-Vite-646CFF?logo=vite)
- **TypeScript** ![TypeScript](https://img.shields.io/badge/-TypeScript-3178C6?logo=typescript)
- **TailwindCSS** ![TailwindCSS](https://img.shields.io/badge/-TailwindCSS-38BDF8?logo=tailwindcss)
- **Zustand** ![Zustand](https://img.shields.io/badge/-Zustand-ff9100?logo=react)
- **React Query** ![React Query](https://img.shields.io/badge/-React%20Query-FF4154?logo=react-query)
- **Axios** ![Axios](https://img.shields.io/badge/-Axios-5A29E4?logo=axios)
- **React Router** ![React Router](https://img.shields.io/badge/-React%20Router-CA4245?logo=react-router)
- **Framer Motion** ![Framer Motion](https://img.shields.io/badge/-Framer%20Motion-EF008F?logo=framer)
- **Three.js** ![Three.js](https://img.shields.io/badge/-Three.js-000000?logo=three.js)
- **Lucide** ![Lucide](https://img.shields.io/badge/-Lucide-000000?logo=lucide)
- **Vercel** ![Vercel](https://img.shields.io/badge/-Vercel-000000?logo=vercel)

---

## Project Structure

```
client--/
├── public/                # Static assets (logo, favicon, etc.)
├── src/
│   ├── components/        # Reusable UI and layout components
│   │   ├── effects/       # Visual/animation effects (Antigravity, Balatro, TargetCursor)
│   │   ├── layout/        # Navbar, layout wrappers
│   │   └── ui/            # Button, Card, Input, etc.
│   ├── lib/               # API utilities (axios instance)
│   ├── pages/             # Route-based pages (Landing, Auth, Problems, Contests, etc.)
│   ├── services/          # API service modules (problemsService, contestService, etc.)
│   ├── store/             # Zustand stores (authStore)
│   ├── types/             # TypeScript types
│   ├── App.tsx           # Main app component (routing)
│   └── main.tsx          # Entry point
├── .env                   # Environment variables
├── package.json           # Project metadata & scripts
├── tailwind.config.js     # TailwindCSS config
├── vite.config.ts         # Vite config
└── ...
```

---

## Environment Setup

1. **Install dependencies:**
   ```bash
   npm install
   # or
   yarn install
   ```
2. **Configure environment:**
   - Copy `.env.example` to `.env` and set API base URL, etc.
3. **Run development server:**
   ```bash
   npm run dev
   # or
   yarn dev
   ```
4. **Build for production:**
   ```bash
   npm run build
   # or
   yarn build
   ```

---

## Key Features

- 🔐 **Authentication:** Email/password, Google, GitHub OAuth
- 👤 **Profile:** View and edit user info, platform handles
- 🏆 **Problems:** Search, filter, mark as solved, sync from platforms
- 📅 **Contests:** Upcoming/past contests, reminders
- 📝 **Sheets:** Create, edit, share problem sheets
- 🧑‍🤝‍🧑 **Rooms:** Real-time collaboration, code sync, whiteboard
- 📊 **Dashboard:** Stats, recent activity, solved count
- 🎨 **Modern UI:** Responsive, animated, dark mode

---

## Pages & Components

### Main Pages
- **Landing:** Home page, intro, call-to-action
- **Login/Register:** Auth forms, OAuth buttons
- **Dashboard:** User stats, recent problems, upcoming contests
- **Problems:** List/search problems, mark as solved, sync
- **Contests:** List, filter, reminders
- **Sheets:** Manage problem sheets
- **Rooms:** List, join, create, real-time collab
- **Profile:** User info, platform stats
- **Settings:** Platform handle management

### Notable Components
- **Navbar:** Top navigation bar
- **Antigravity, Balatro, TargetCursor:** Visual effects
- **Button, Card, Input:** Custom UI primitives
- **Loader, Skeleton:** Loading states

---

## API Integration

- **API Base URL:** Set in `.env` (e.g., `VITE_API_URL`)
- **All API calls** use Axios via `lib/api.ts`
- **Service modules** in `src/services/` (e.g., `problemsService.ts`, `contestService.ts`)
- **Handles JWT tokens** (auto-refresh, error handling)
- **Error handling:** User-friendly messages, toast notifications

---

## State Management

- **Zustand** for global state (auth, user info)
- **React Query** (if enabled) for API caching, background updates
- **Local state** for page/component UI

---

## Styling

- **TailwindCSS** for utility-first styling
- **Custom themes:** Dark mode, accent colors
- **Responsive design:** Mobile-first, grid/flex layouts
- **Animated icons:** Lucide, Framer Motion

---

## Testing & Linting

- **ESLint** for code quality
- **Prettier** (if enabled) for formatting
- **Manual testing** via browser, Postman, Thunder Client
- **Unit tests:** (Add with Jest/React Testing Library as needed)

---

## Deployment

- **Vercel** for instant deploys (auto from GitHub)
- **Production build:** `npm run build` → `dist/`
- **Environment variables:** Set in Vercel dashboard
- **Preview URLs** for PRs/branches

---

## Useful Scripts

- `npm run dev` — Start dev server
- `npm run build` — Build for production
- `npm run lint` — Lint code
- `npm run preview` — Preview production build

---

## Contributing

- Fork the repo, create a feature branch, open a PR
- Follow code style and commit guidelines
- Add/Update documentation for new features

---

## Credits

- [React](https://react.dev/)
- [Vite](https://vitejs.dev/)
- [TailwindCSS](https://tailwindcss.com/)
- [Zustand](https://docs.pmnd.rs/zustand/getting-started/introduction)
- [Lucide](https://lucide.dev/)
- [Framer Motion](https://www.framer.com/motion/)
- [Three.js](https://threejs.org/)
- [Vercel](https://vercel.com/)

---

**Last Updated:** February 10, 2026  
**Version:** 1.0.0
