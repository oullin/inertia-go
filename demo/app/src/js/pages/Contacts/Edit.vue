<script setup>
import { useForm } from "@inertiajs/vue3";
import ContactForm from "./Partials/ContactForm.vue";

const props = defineProps({
  contact: {
    type: Object,
    required: true,
  },
  form: {
    type: Object,
    required: true,
  },
  organizations: {
    type: Array,
    default: () => [],
  },
});

const form = useForm({ ...props.form });
</script>

<template>
  <ContactForm mode="edit" :contact="contact" :form="form" :organizations="organizations" />
</template>
