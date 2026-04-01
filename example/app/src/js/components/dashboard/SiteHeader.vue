<script setup>
import { IconCirclePlusFilled } from "@tabler/icons-vue";
import { Button } from "@/js/components/ui/button";
import { Separator } from "@/js/components/ui/separator";
import { SidebarTrigger } from "@/js/components/ui/sidebar";

defineProps({
  title: {
    type: String,
    default: "Documents",
  },
});
</script>

<template>
  <header
    class="bg-background/90 sticky top-0 z-10 flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)"
  >
    <div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
      <SidebarTrigger class="-ml-1" />
      <Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
      <h1 class="text-base font-medium">{{ title }}</h1>
      <div class="ml-auto flex items-center gap-2">
        <Button size="sm" class="hidden h-7 sm:flex">
          <IconCirclePlusFilled />
          <span>Quick Create</span>
        </Button>
      </div>
    </div>
  </header>
</template>
