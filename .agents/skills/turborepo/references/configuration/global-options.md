# Global Options Reference

Options that affect all tasks. Full docs: https://turborepo.dev/docs/reference/configuration

## globalEnv

Environment variables affecting all task hashes.

```json
{
  "globalEnv": ["CI", "NODE_ENV", "VERCEL_*"]
}
```

Use for variables that should invalidate all caches when changed.

## globalDependencies

Files that affect all task hashes.

```json
{
  "globalDependencies": ["tsconfig.json", ".env", "pnpm-lock.yaml"]
}
```

Lockfile is included by default. Add shared configs here.

## globalPassThroughEnv

Variables available to tasks but not included in hash.

```json
{
  "globalPassThroughEnv": ["AWS_SECRET_KEY", "GITHUB_TOKEN"]
}
```

Use for credentials that shouldn't affect cache keys.

## cacheDir

Custom cache location. Default: `node_modules/.cache/turbo`.

```json
{
  "cacheDir": ".turbo/cache"
}
```

## daemon

Background process for faster subsequent runs. Default: `true`.

```json
{
  "daemon": false
}
```

Disable in CI or when debugging.

## envMode

How unspecified env vars are handled in `turbo.json`. Requires Turborepo `2.1+`. In `2.0+`, the default is `"strict"`. On `1.10` through `2.0`, use the CLI `--env-mode` flag instead.

```json
{
  "envMode": "strict"  // Only specified vars available
  // or
  "envMode": "loose"   // All vars pass through
}
```

Strict mode catches missing env declarations.

## ui

Terminal UI mode. Default: `"stream"`.

```json
{
  "ui": "tui"     // Interactive terminal UI
  // or
  "ui": "stream"  // Traditional streaming logs
}
```

TUI provides better UX for parallel tasks.

## remoteCache

Configure remote caching with the global `remoteCache` block in `turbo.json`.

```json
{
  "remoteCache": {
    "enabled": true,
    "signature": true
  }
}
```

For the full option table, setup flow, signing, CI usage, and debugging, see [Remote Caching](../caching/remote-cache.md#configuration-in-turbojson).

## concurrency

Default: `"10"`

Limit parallel task execution.

```json
{
  "concurrency": "4"     // Max 4 tasks at once
  // or
  "concurrency": "50%"   // 50% of available CPUs
}
```

## futureFlags

Enable experimental features that will become default in future versions.

```json
{
  "futureFlags": {
    "errorsOnlyShowHash": true
  }
}
```

### `errorsOnlyShowHash`

When using `outputLogs: "errors-only"`, show task hashes on start/completion:

- Cache miss: `cache miss, executing <hash> (only logging errors)`
- Cache hit: `cache hit, replaying logs (no errors) <hash>`

## noUpdateNotifier

Disable update notifications when new turbo versions are available.

```json
{
  "noUpdateNotifier": true
}
```

## dangerouslyDisablePackageManagerCheck

Bypass the `packageManager` field requirement. Use for incremental migration.

```json
{
  "dangerouslyDisablePackageManagerCheck": true
}
```

**Warning**: Unstable lockfiles can cause unpredictable behavior.

## Git Worktree Cache Sharing

When working in Git worktrees, Turborepo automatically shares local cache between the main worktree and linked worktrees.

**How it works:**

- Detects worktree configuration
- Redirects cache to main worktree's `.turbo/cache`
- Works alongside Remote Cache

**Benefits:**

- Cache hits across branches
- Reduced disk usage
- Faster branch switching

**Disabled by**: Setting explicit `cacheDir` in turbo.json.
