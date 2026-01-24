# Project Overview

**tons** is an AI-powered translation desktop app built with Wails v3 (Go backend + Svelte frontend). It supports multiple translation engines: built-in Yzma, CLI agents (Claude Code, Gemini CLI, Codex), and Ollama.

## Tech Stack

- **Backend**: Go 1.25, Wails v3 (alpha)
- **Frontend**: SvelteKit, Svelte 5 (runes), TypeScript, Tailwind CSS 4, bits-ui
- **Package Manager**: pnpm (frontend)

## Common Commands

```bash
# Build for production
wails3 build
# Generate Go->TypeScript bindings (auto-runs on build)
wails3 generate bindings -ts -d frontend/src/lib/bindings
```

## Structure

- `internal/`
  - `config/` App configuration (JSON persistence)
  - `engine/` Translation engine implementations
  - `services/` Wails services exposed to frontend
- `frontend/` Svelte kit

**Key Pattern**: Services in `internal/services/` are registered in `main.go` and auto-generate TypeScript bindings in `frontend/src/lib/bindings/`.
