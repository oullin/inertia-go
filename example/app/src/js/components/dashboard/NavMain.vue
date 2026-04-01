<script setup>
import { computed } from "vue";
import { Link, usePage } from "@inertiajs/vue3";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/js/components/ui/sidebar";

defineProps({
  items: {
    type: Array,
    required: true,
  },
});

const page = usePage();
const currentPath = computed(() => String(page.url ?? "/").split("?")[0]);
</script>

<template>
  <SidebarGroup>
    <SidebarGroupContent>
      <SidebarGroupLabel>Platform</SidebarGroupLabel>
      <SidebarMenu>
        <SidebarMenuItem v-for="item in items" :key="item.title">
          <SidebarMenuButton as-child :tooltip="item.title" :is-active="currentPath === item.url">
            <Link :href="item.url" view-transition>
              <component :is="item.icon" v-if="item.icon" />
              <span>{{ item.title }}</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroupContent>
  </SidebarGroup>
</template>
