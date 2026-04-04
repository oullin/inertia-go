import { usePage } from "@inertiajs/vue3";
import type { DemoRoute, SharedPageProps } from "@/js/types";

function fillPattern(pattern: string, params: Record<string, string | number> = {}): string {
  let url = pattern ?? "#!wayfinder:unknown-route";

  Object.entries(params).forEach(([key, value]) => {
    url = url.replaceAll(`{${key}}`, encodeURIComponent(String(value)));
  });

  return url;
}

export function useDemoRoute(
  name: string,
  params: Record<string, string | number> = {},
): DemoRoute {
  const page = usePage<SharedPageProps>();
  const pattern = page.props.routes?.[name];

  if (!pattern) {
    console.warn(`[wayfinder] unknown route "${name}", returning fallback`);
  }

  return {
    url: fillPattern(pattern, params),
  };
}

export function featureRoute(name: string): string {
  return useDemoRoute(name).url;
}

export const appRoutes = {
  login: (): DemoRoute => useDemoRoute("login"),
  logout: (): DemoRoute => useDemoRoute("logout"),
  dashboard: (): DemoRoute => useDemoRoute("dashboard"),
};

export const contactRoutes = {
  index: (): DemoRoute => useDemoRoute("contacts.index"),
  create: (): DemoRoute => useDemoRoute("contacts.create"),
  store: (): DemoRoute => useDemoRoute("contacts.store"),
  show: (contact: string | number): DemoRoute => useDemoRoute("contacts.show", { contact }),
  edit: (contact: string | number): DemoRoute => useDemoRoute("contacts.edit", { contact }),
  update: (contact: string | number): DemoRoute => useDemoRoute("contacts.update", { contact }),
  destroy: (contact: string | number): DemoRoute => useDemoRoute("contacts.destroy", { contact }),
  favorite: (contact: string | number): DemoRoute => useDemoRoute("contacts.favorite", { contact }),
  storeNote: (contact: string | number): DemoRoute =>
    useDemoRoute("contacts.notes.store", { contact }),
};

export const organizationRoutes = {
  index: (): DemoRoute => useDemoRoute("organizations.index"),
  show: (organization: string | number): DemoRoute =>
    useDemoRoute("organizations.show", { organization }),
  update: (organization: string | number): DemoRoute =>
    useDemoRoute("organizations.update", { organization }),
};
