import { usePage } from "@inertiajs/vue3";

function fillPattern(pattern, params = {}) {
  let url = pattern ?? "/";

  Object.entries(params).forEach(([key, value]) => {
    url = url.replaceAll(`{${key}}`, String(value));
  });

  return url;
}

export function useDemoRoute(name, params = {}) {
  const page = usePage();
  const pattern = page.props.routes?.[name] ?? "/";

  return {
    url: fillPattern(pattern, params),
  };
}

export const appRoutes = {
  login: () => useDemoRoute("login"),
  logout: () => useDemoRoute("logout"),
  dashboard: () => useDemoRoute("dashboard"),
};

export const contactRoutes = {
  index: () => useDemoRoute("contacts.index"),
  create: () => useDemoRoute("contacts.create"),
  store: () => useDemoRoute("contacts.store"),
  show: (contact) => useDemoRoute("contacts.show", { contact }),
  edit: (contact) => useDemoRoute("contacts.edit", { contact }),
  update: (contact) => useDemoRoute("contacts.update", { contact }),
  favorite: (contact) => useDemoRoute("contacts.favorite", { contact }),
  storeNote: (contact) => useDemoRoute("contacts.notes.store", { contact }),
};

export const organizationRoutes = {
  index: () => useDemoRoute("organizations.index"),
  show: (organization) => useDemoRoute("organizations.show", { organization }),
  update: (organization) => useDemoRoute("organizations.update", { organization }),
};

export const featureRoutes = {
  formsUseForm: () => useDemoRoute("features.forms.use-form"),
  navigationLinks: () => useDemoRoute("features.navigation.links"),
  deferredProps: () => useDemoRoute("features.data-loading.deferred"),
  remember: () => useDemoRoute("features.state.remember"),
};
