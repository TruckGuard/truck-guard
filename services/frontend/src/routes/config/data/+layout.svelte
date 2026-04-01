<script lang="ts">
  import { page } from "$app/state";
  import * as Tabs from "$lib/components/ui/tabs";
  import { goto } from "$app/navigation";

  let { children } = $props();

  const tabs = [
    { id: "companies", label: "Компанії", path: "/config/data/companies" },
    { id: "posts", label: "Митні пости", path: "/config/data/posts" },
    { id: "modes", label: "Режими", path: "/config/data/modes" },
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
    {
      id: "excluded-plates",
      label: "Ігнорування",
      path: "/config/data/excluded-plates",
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

<div class="flex flex-col h-full gap-8 overflow-hidden">
  <Tabs.Root value={activeTab} onValueChange={handleTabChange} class="w-full shrink-0">
    <Tabs.List class="flex w-full items-center justify-start border-b bg-transparent h-auto p-0 gap-8 rounded-none border-muted/60">
      {#each tabs as tab}
        <Tabs.Trigger
          value={tab.id}
          class="relative h-12 px-1 bg-transparent text-sm font-medium transition-all text-muted-foreground hover:text-foreground data-[state=active]:text-primary data-[state=active]:after:absolute data-[state=active]:after:-bottom-px data-[state=active]:after:left-0 data-[state=active]:after:right-0 data-[state=active]:after:h-[2px] data-[state=active]:after:bg-primary shadow-none border-none rounded-none"
        >
          {tab.label}
        </Tabs.Trigger>
      {/each}
    </Tabs.List>
  </Tabs.Root>

  <div class="flex-1 min-h-0 overflow-hidden">
    {@render children()}
  </div>
</div>
