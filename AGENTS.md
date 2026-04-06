# inertia-go Agent Guide

`AGENTS.md` is the canonical root entrypoint for agent guidance in this repo. `CLAUDE.md` and `GEMINI.md` should resolve to this same content.

## Repo Summary

`inertia-go` is a Go server-side adapter for [Inertia.js](https://inertiajs.com). The main maintained library surface is [`core/`](core), with a demo app under [`demo/`](demo).

## Quick Reference

Run core library commands from the repo root:

- Build: `cd core && go build ./...`
- Test: `cd core && go test ./...`
- Test (match CI): `cd core && go test -race ./...`
- Vet: `cd core && go vet ./...`
- Format: `make format`
- Tidy: `make tidy`

Formatting uses Docker Compose via [`go-fmt.compose.yaml`](go-fmt.compose.yaml). Do not use `gofmt` directly if you need CI-compatible formatting.

## Detailed Guidance

- [Core Architecture](.agents/core-architecture.md)
- [Core Props And Rendering](.agents/core-props-and-rendering.md)
- [Core Testing](.agents/core-testing.md)
- [Inertia Go Skill](.agents/skills/inertia-go/SKILL.md)
