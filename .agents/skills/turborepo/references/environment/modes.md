# Environment Modes

Turborepo supports different modes for handling environment variables during task execution.

As of April 6, 2026, environment modes are generally available in supported Turborepo `1.x` and `2.x` releases starting with `1.10`. `Strict Mode` becomes the default in `2.0`. The `turbo.json` `envMode` key is newer and requires `2.1+`.

## Strict Mode

Default in Turborepo `2.0+`.

Only explicitly configured variables are available to tasks.

**Behavior:**

- Tasks only see vars listed in `env`, `globalEnv`, `passThroughEnv`, or `globalPassThroughEnv`
- Unlisted vars are filtered out
- Tasks fail if they require unlisted variables

**Benefits:**

- Guarantees cache correctness
- Prevents accidental dependencies on system vars
- Reproducible builds across machines

```bash
# Explicit in versions that support --env-mode
turbo run build --env-mode=strict
```

## Loose Mode

All system environment variables are available to tasks.

```bash
turbo run build --env-mode=loose
```

**Behavior:**

- Every system env var is passed through
- Only vars in `env`/`globalEnv` affect the hash
- Other vars are available but NOT hashed

**Risks:**

- Cache may restore incorrect results if unhashed vars changed
- "Works on my machine" bugs
- CI vs local environment mismatches

**Use case:** Migrating legacy projects or debugging strict mode issues.

## Framework Inference (Automatic)

Turborepo automatically detects frameworks and includes their conventional env vars.

### Inferred Variables by Framework

| Framework        | Pattern             |
| ---------------- | ------------------- |
| Next.js          | `NEXT_PUBLIC_*`     |
| Vite             | `VITE_*`            |
| Create React App | `REACT_APP_*`       |
| Gatsby           | `GATSBY_*`          |
| Nuxt             | `NUXT_*`, `NITRO_*`, `SERVER_*`, `AWS_APP_ID`, `INPUT_AZURE_STATIC_WEB_APPS_API_TOKEN`, `CLEAVR`, `CF_PAGES`, `FIREBASE_APP_HOSTING`, `NETLIFY`, `STORMKIT`, `NOW_BUILDER`, `ZEABUR`, `RENDER`, `LAUNCH_EDITOR` |
| Expo             | `EXPO_PUBLIC_*`     |
| Astro            | `PUBLIC_*`          |
| SvelteKit        | `PUBLIC_*`          |
| Remix            | `REMIX_*`           |
| Redwood          | `REDWOOD_ENV_*`     |
| Sanity           | `SANITY_STUDIO_*`   |
| Solid            | `VITE_*`            |

### Disabling Framework Inference

Globally via CLI:

```bash
turbo run build --framework-inference=false
```

Or exclude specific patterns in config:

```json
{
  "tasks": {
    "build": {
      "env": ["!NEXT_PUBLIC_*"]
    }
  }
}
```

### Why Disable?

- You want explicit control over all env vars
- Framework vars shouldn't bust the cache (e.g., analytics IDs)
- Debugging unexpected cache misses

## Checking Environment Mode

Use `--dry` to see which vars affect each task:

This example requires `jq` to be installed because it extracts environment variables from the JSON output.

```bash
turbo run build --dry=json | jq '.tasks[].environmentVariables'
```
