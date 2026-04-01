<script setup>
import { Link } from '@inertiajs/vue3'

defineProps({
    title: String,
    description: String,
    features: Array,
    app: Object,
})
</script>

<template>
    <div class="container">
        <header>
            <h1>{{ app?.name }}</h1>
            <nav>
                <Link href="/">Home</Link>
                <Link href="/about" class="active">About</Link>
                <Link href="/users">Users</Link>
            </nav>
        </header>

        <main>
            <h2>{{ title }}</h2>
            <p class="description">{{ description }}</p>

            <div class="card">
                <h3>Features</h3>
                <ul>
                    <li v-for="feature in features" :key="feature">{{ feature }}</li>
                </ul>
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
.description { color: #666; font-size: 1.1rem; margin-bottom: 2rem; }
.card { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 1.5rem; }
.card h3 { margin-top: 0; color: #1a1a2e; }
ul { list-style: none; padding: 0; }
li { padding: 0.5rem 0; color: #64748b; border-bottom: 1px solid #e2e8f0; }
li:last-child { border-bottom: none; }
li::before { content: "->"; margin-right: 0.5rem; color: #6366f1; }
</style>
