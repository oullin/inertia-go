<script setup>
import { SidebarInset, SidebarProvider } from "@/js/components/ui/sidebar";
import AppSidebar from "@/js/components/dashboard/AppSidebar.vue";
import SiteHeader from "@/js/components/dashboard/SiteHeader.vue";
import FlashBanner from "@/js/components/dashboard/FlashBanner.vue";
import LoadingBar from "@/js/components/ui/loading-bar/LoadingBar.vue";
import PageSkeleton from "@/js/components/dashboard/PageSkeleton.vue";
import { useInertiaLoading } from "@/js/composables/useInertiaLoading";

const { isLoading } = useInertiaLoading();
</script>

<template>
  <LoadingBar />
  <SidebarProvider
    :default-open="true"
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
          <PageSkeleton v-if="isLoading" />
          <slot v-else />
        </div>
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
