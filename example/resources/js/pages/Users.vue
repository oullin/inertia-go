<script setup>
import { Link } from "@inertiajs/vue3";

defineProps({
  title: String,
  users: Array,
  stats: Object,
  metadata: Object,
  app: Object,
});
</script>

<template>
  <div class="container">
    <header>
      <h1>{{ app?.name }}</h1>
      <nav>
        <Link href="/">Home</Link>
        <Link href="/about">About</Link>
        <Link href="/users" class="active">Users</Link>
      </nav>
    </header>

    <main>
      <h2>{{ title }}</h2>

      <div v-if="stats" class="stats">
        <div class="stat">
          <span class="stat-value">{{ stats.total }}</span>
          <span class="stat-label">Total Users</span>
        </div>
        <div class="stat">
          <span class="stat-value">{{ stats.active }}</span>
          <span class="stat-label">Active</span>
        </div>
      </div>

      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>
              <span class="badge">{{ user.role }}</span>
            </td>
          </tr>
        </tbody>
      </table>

      <p v-if="metadata" class="meta">Last synced: {{ metadata.last_sync }}</p>
    </main>
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
  font-family: system-ui, sans-serif;
}
header {
  margin-bottom: 2rem;
}
h1 {
  color: #1a1a2e;
  margin-bottom: 1rem;
}
nav {
  display: flex;
  gap: 1rem;
}
nav a {
  color: #6366f1;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  transition: background 0.2s;
}
nav a:hover,
nav a.active {
  background: #eef2ff;
}
.stats {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}
.stat {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  text-align: center;
}
.stat-value {
  display: block;
  font-size: 2rem;
  font-weight: 700;
  color: #6366f1;
}
.stat-label {
  color: #64748b;
  font-size: 0.875rem;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th {
  text-align: left;
  padding: 0.75rem;
  border-bottom: 2px solid #e2e8f0;
  color: #64748b;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
td {
  padding: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
}
.badge {
  background: #eef2ff;
  color: #6366f1;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.875rem;
}
.meta {
  color: #94a3b8;
  font-size: 0.875rem;
  margin-top: 1rem;
}
</style>
