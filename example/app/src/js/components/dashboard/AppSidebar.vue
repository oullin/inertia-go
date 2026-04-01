<script setup>
import { computed } from "vue";
import { usePage } from "@inertiajs/vue3";
import {
  IconChartBar,
  IconDashboard,
  IconFolder,
  IconInnerShadowTop,
  IconListDetails,
  IconSettings,
  IconUsers,
} from "@tabler/icons-vue";
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
import NavSecondary from "./NavSecondary.vue";
import NavUser from "./NavUser.vue";

const page = usePage();
const user = computed(
  () => page.props.auth?.user ?? { name: "Gus", email: "gus@example.com", initials: "GC" },
);

const navMain = [
  { title: "Dashboard", url: "/dashboard", icon: IconDashboard },
  { title: "Lifecycle", url: "/dashboard/navigation", icon: IconListDetails },
  { title: "Analytics", url: "/dashboard/data", icon: IconChartBar },
  { title: "Projects", url: "/dashboard/forms", icon: IconFolder },
  { title: "Team", url: "/dashboard/feed", icon: IconUsers },
];

const navSecondary = [{ title: "Settings", url: "/dashboard/state", icon: IconSettings }];
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
      <NavSecondary :items="navSecondary" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="user" />
    </SidebarFooter>
  </Sidebar>
</template>
