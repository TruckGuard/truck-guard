<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon.svg";
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import { ModeWatcher } from "mode-watcher";
  import { Toaster } from "$lib/components/ui/sonner";
  import { toast } from "svelte-sonner";
  import { goto } from "$app/navigation";
  import { onMount, onDestroy } from "svelte";

  let { children, data } = $props();

  let sse: EventSource | null = null;

  function connectNotifications() {
    sse = new EventSource("/api/notifications/stream");

    sse.onmessage = (e) => {
      try {
        const ev = JSON.parse(e.data);
        if (ev.type === "new_permit") {
          toast("Нова перепустка", {
            description: ev.plate
              ? `${ev.plate}${ev.code ? ` · ${ev.code}` : ""}`
              : ev.code || "Натисніть, щоб відкрити",
            duration: 10000,
            action: {
              label: "Відкрити →",
              onClick: () => goto(`/permits/${ev.id}`),
            },
            id: String(ev.id),
          });
        }
      } catch {}
    };
  }

  onMount(() => {
    if (data.user?.customs_post_id) {
      connectNotifications();
    }
  });

  onDestroy(() => {
    sse?.close();
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>TruckGuard</title>
</svelte:head>
<ModeWatcher defaultMode="light" />
<Toaster />

{#if data.user}
  <Sidebar.Provider>
    <AppSidebar user={data.user} />
    <Sidebar.Inset class="h-svh overflow-hidden flex flex-col">
      <header class="flex h-14 shrink-0 items-center gap-2 border-b px-4 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 z-30">
        <Sidebar.Trigger class="-ms-1" />
        <Separator orientation="vertical" class="me-2 h-4" />
      </header>
      <div class="flex-1 flex flex-col min-h-0 overflow-auto p-4 pt-0">
        {@render children()}
      </div>
    </Sidebar.Inset>
  </Sidebar.Provider>
{:else}
  {@render children()}
{/if}
