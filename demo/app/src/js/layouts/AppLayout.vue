<script setup>
import { computed } from "vue";
import { Head, Link, usePage } from "@inertiajs/vue3";
import { IconDashboard, IconUsers } from "@tabler/icons-vue";
import LoadingBar from "@/js/components/ui/loading-bar/LoadingBar.vue";
import { Separator } from "@/js/components/ui/separator";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/js/components/ui/sidebar";
import { appRoutes, contactRoutes, organizationRoutes } from "@/js/lib/routes";

const props = defineProps({
  title: {
    type: String,
    default: "",
  },
  breadcrumbs: {
    type: Array,
    default: () => [],
  },
});

const page = usePage();
const currentPath = computed(() => String(page.url ?? "/").split("?")[0]);
const user = computed(() => page.props.auth?.user ?? null);

const groups = computed(() => [
  {
    label: "CRM",
    items: [
      { title: "Dashboard", href: appRoutes.dashboard().url, icon: IconDashboard },
      { title: "Contacts", href: contactRoutes.index().url, icon: IconUsers },
      { title: "Organizations", href: organizationRoutes.index().url, icon: IconUsers },
    ],
  },
]);
</script>

<template>
  <Head :title="title" />
  <LoadingBar />
  <SidebarProvider
    :default-open="true"
    :style="{
      '--sidebar-width': 'calc(var(--spacing) * 68)',
      '--header-height': 'calc(var(--spacing) * 12 + 1px)',
    }"
  >
    <Sidebar collapsible="icon" class="h-auto border-r">
      <SidebarHeader class="border-b">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton as-child class="data-[slot=sidebar-menu-button]:!p-1.5">
              <Link :href="appRoutes.dashboard().url">
                <IconDashboard class="!size-5" />
                <span class="text-base font-semibold">Inertia Kitchen Sink</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup v-for="group in groups" :key="group.label">
          <SidebarGroupContent>
            <SidebarGroupLabel>{{ group.label }}</SidebarGroupLabel>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in group.items" :key="item.title">
                <SidebarMenuButton
                  as-child
                  :is-active="currentPath === item.href"
                  :tooltip="item.title"
                >
                  <Link :href="item.href" view-transition>
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter v-if="user" class="border-t">
        <div class="flex items-center justify-between gap-3 px-2 py-1 text-sm">
          <div class="min-w-0">
            <p class="truncate font-medium">{{ user.name }}</p>
            <p class="text-muted-foreground truncate text-xs">{{ user.email }}</p>
          </div>
          <SidebarMenuButton as-child size="sm" tooltip="Log out">
            <Link :href="appRoutes.logout().url" method="post" as="button" class="cursor-pointer">
              <span>Log out</span>
            </Link>
          </SidebarMenuButton>
        </div>
      </SidebarFooter>
    </Sidebar>

    <SidebarInset>
      <header
        class="bg-background/90 sticky top-0 z-10 flex h-(--header-height) shrink-0 items-center gap-2 border-b"
      >
        <div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
          <div class="flex flex-col">
            <h1 class="text-base font-medium">{{ title }}</h1>
            <p v-if="breadcrumbs.length" class="text-muted-foreground text-xs">
              {{ breadcrumbs.map((item) => item.title).join(" / ") }}
            </p>
          </div>
        </div>
      </header>

      <div
        v-if="page.props.flash?.message"
        class="mx-4 mt-4 rounded-lg border p-4 lg:mx-6"
        :class="{
          'border-green-200 bg-green-50 text-green-800': page.props.flash.kind === 'success',
          'border-red-200 bg-red-50 text-red-800': page.props.flash.kind === 'error',
          'border-blue-200 bg-blue-50 text-blue-800':
            page.props.flash.kind === 'info' || !page.props.flash.kind,
        }"
      >
        <p class="text-sm font-medium">{{ page.props.flash.title || "Notice" }}</p>
        <p class="text-sm">{{ page.props.flash.message }}</p>
      </div>

      <div class="flex flex-1 flex-col">
        <slot />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
