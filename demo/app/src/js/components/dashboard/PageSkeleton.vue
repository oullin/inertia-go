<script setup>
import { Skeleton } from "@/js/components/ui/skeleton";
</script>

<template>
  <div class="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
    <div class="grid grid-cols-1 gap-4 px-4 lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <Skeleton class="h-24 rounded-xl" />
      <Skeleton class="h-24 rounded-xl" />
      <Skeleton class="h-24 rounded-xl" />
      <Skeleton class="h-24 rounded-xl" />
    </div>

    <div class="px-4 lg:px-6">
      <Skeleton class="h-64 w-full rounded-xl" />
    </div>

    <div class="flex flex-col gap-3 px-4 lg:px-6">
      <Skeleton class="h-10 w-full rounded-md" />
      <Skeleton class="h-4 w-full rounded-md" />
      <Skeleton class="h-4 w-5/6 rounded-md" />
      <Skeleton class="h-4 w-4/5 rounded-md" />
      <Skeleton class="h-4 w-full rounded-md" />
      <Skeleton class="h-4 w-3/4 rounded-md" />
    </div>
  </div>
</template>
