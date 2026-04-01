<script setup>
import { SidebarInset, SidebarProvider } from "@/js/components/ui/sidebar";
import AppSidebar from "@/js/components/dashboard/AppSidebar.vue";
import SiteHeader from "@/js/components/dashboard/SiteHeader.vue";
import FlashBanner from "@/js/components/dashboard/FlashBanner.vue";
</script>

<template>
  <SidebarProvider
    :style="{
      '--sidebar-width': 'calc(var(--spacing) * 64)',
      '--header-height': 'calc(var(--spacing) * 12 + 1px)',
    }"
  >
    <AppSidebar />
    <SidebarInset>
      <SiteHeader />
      <FlashBanner />
      <div class="flex flex-1 flex-col">
        <div class="@container/main flex flex-1 flex-col gap-2">
          <slot />
        </div>
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
