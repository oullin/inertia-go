<script setup>
import { computed } from "vue";
import { usePage } from "@inertiajs/vue3";
import { IconDashboard, IconInnerShadowTop, IconListDetails } from "@tabler/icons-vue";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/js/components/ui/sidebar";
import NavMain from "./NavMain.vue";
import NavUser from "./NavUser.vue";

const page = usePage();
const user = computed(
  () => page.props.auth?.user ?? { name: "Gus", email: "gus@example.com", initials: "GC" },
);

const navMain = [
  { title: "Dashboard", url: "/dashboard", icon: IconDashboard },
  {
    title: "Lifecycle",
    icon: IconListDetails,
    children: [
      { title: "Navigation", url: "/dashboard/navigation" },
      { title: "Scroll", url: "/dashboard/scroll" },
      { title: "Redirects", url: "/dashboard/redirects" },
      { title: "Data Loading", url: "/dashboard/data" },
      { title: "Forms", url: "/dashboard/forms" },
      { title: "Feed", url: "/dashboard/feed" },
      { title: "State", url: "/dashboard/state" },
    ],
  },
];
</script>

<template>
  <Sidebar collapsible="icon" class="h-auto border-r">
    <SidebarHeader class="border-b">
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton as-child class="data-[slot=sidebar-menu-button]:!p-1.5">
            <a href="/dashboard">
              <IconInnerShadowTop class="!size-5" />
              <span class="text-base font-semibold">Progressive Oullin</span>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="navMain" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="user" />
    </SidebarFooter>
  </Sidebar>
</template>
