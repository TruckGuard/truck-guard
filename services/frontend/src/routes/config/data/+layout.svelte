<script lang="ts">
  import { page } from "$app/state";
  import * as Tabs from "$lib/components/ui/tabs";
  import { goto } from "$app/navigation";

  let { children } = $props();

  const tabs = [
    { id: "companies", label: "Компанії", path: "/config/data/companies" },
    { id: "posts", label: "Митні пости", path: "/config/data/posts" },
    { id: "modes", label: "Режими митниці", path: "/config/data/modes" },
    {
      id: "vehicle-types",
      label: "Типи ТЗ",
      path: "/config/data/vehicle-types",
    },
    {
      id: "payment-types",
      label: "Типи оплати",
      path: "/config/data/payment-types",
    },
  ];

  let activeTab = $derived(
    tabs.find((t) => page.url.pathname.startsWith(t.path))?.id || "companies",
  );

  function handleTabChange(value: string) {
    const tab = tabs.find((t) => t.id === value);
    if (tab) {
      goto(tab.path);
    }
  }
</script>

<div class="flex flex-col gap-6">
  <div class="flex flex-col gap-1">
    <h1 class="text-3xl font-bold tracking-tight">Довідники</h1>
    <p class="text-muted-foreground">Керування основними даними системи.</p>
  </div>

  <Tabs.Root value={activeTab} onValueChange={handleTabChange} class="w-full">
    <Tabs.List class="grid w-full max-w-4xl grid-cols-5">
      {#each tabs as tab}
        <Tabs.Trigger value={tab.id}>{tab.label}</Tabs.Trigger>
      {/each}
    </Tabs.List>
    <div class="mt-6">
      {@render children()}
    </div>
  </Tabs.Root>
</div>
