<script setup>
import { Link } from '@inertiajs/vue3'

defineProps({
    title: String,
    message: String,
    app: Object,
})
</script>

<template>
    <div class="container">
        <header>
            <h1>{{ app?.name }}</h1>
            <nav>
                <Link href="/" class="active">Home</Link>
                <Link href="/about">About</Link>
                <Link href="/users">Users</Link>
            </nav>
        </header>

        <main>
            <h2>{{ title }}</h2>
            <p class="message">{{ message }}</p>

            <div class="card">
                <h3>Getting Started</h3>
                <p>This example demonstrates the Inertia.js Go adapter with Vue 3 and Vite.</p>
                <p>Navigate between pages using the links above. Notice how the page transitions happen without a full page reload.</p>
            </div>
        </main>
    </div>
</template>

<style scoped>
.container { max-width: 800px; margin: 0 auto; padding: 2rem 1rem; font-family: system-ui, sans-serif; }
header { margin-bottom: 2rem; }
h1 { color: #1a1a2e; margin-bottom: 1rem; }
nav { display: flex; gap: 1rem; }
nav a { color: #6366f1; text-decoration: none; padding: 0.5rem 1rem; border-radius: 6px; transition: background 0.2s; }
nav a:hover, nav a.active { background: #eef2ff; }
main h2 { color: #333; margin-bottom: 0.5rem; }
.message { color: #666; font-size: 1.1rem; margin-bottom: 2rem; }
.card { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 1.5rem; }
.card h3 { margin-top: 0; color: #1a1a2e; }
.card p { color: #64748b; line-height: 1.6; }
</style>
