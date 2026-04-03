<script setup>
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/js/components/ui/card";

defineProps({
  title: {
    type: String,
    default: "Inertia.js Kitchen Sink",
  },
  description: {
    type: String,
    default: "",
  },
});
</script>

<template>
  <div class="bg-muted/30 flex min-h-screen items-center justify-center p-4">
    <Card class="w-full max-w-md shadow-sm">
      <CardHeader class="space-y-2 text-center">
        <CardTitle class="text-2xl">{{ title }}</CardTitle>
        <CardDescription>{{ description }}</CardDescription>
      </CardHeader>
      <CardContent>
        <slot />
      </CardContent>
    </Card>
  </div>
</template>
