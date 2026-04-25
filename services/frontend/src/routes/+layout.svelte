<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon.svg";
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import { ModeWatcher } from "mode-watcher";
  import { Toaster } from "$lib/components/ui/sonner";
  import { toast } from "svelte-sonner";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { onMount, onDestroy } from "svelte";
  import FuzzyMatchModal from "$lib/components/FuzzyMatchModal.svelte";

  let { children, data } = $props();

  // ─── Route label map ────────────────────────────────────────────────────────
  const LABELS: Record<string, string> = {
    events:            "Події",
    permits:           "Журнал перепусток",
    config:            "Конфігурація",
    cameras:           "Камери",
    scales:            "Ваги",
    data:              "Довідники",
    companies:         "Компанії",
    posts:             "Пости",
    modes:             "Режими",
    "payment-types":   "Типи оплати",
    "vehicle-types":   "Типи авто",
    "excluded-plates": "Чорний список",
    settings:          "Налаштування",
    admin:             "Адміністрування",
    users:             "Користувачі",
    roles:             "Ролі та права",
    keys:              "API Ключі",
    audit:             "Аудит",
    create:            "Створення",
    login:             "Вхід",
  };

  function segmentLabel(seg: string): string {
    if (LABELS[seg]) return LABELS[seg];
    // Numeric or UUID-like segments → looks like an ID
    if (/^\d+$/.test(seg)) return `#${seg}`;
    return seg;
  }

  interface Crumb { label: string; href: string | null }

  const crumbs = $derived(buildCrumbs(page.url.pathname));

  function buildCrumbs(pathname: string): Crumb[] {
    const segments = pathname.split("/").filter(Boolean);
    const result: Crumb[] = [];
    let path = "";
    for (let i = 0; i < segments.length; i++) {
      path += "/" + segments[i];
      result.push({
        label: segmentLabel(segments[i]),
        href: i < segments.length - 1 ? path : null,
      });
    }
    return result;
  }

  // ─── SSE notifications ───────────────────────────────────────────────────────
  let sse: EventSource | null = null;

  let showFuzzyModal = $state(false);
  let fuzzyMatchData = $state<any>(null);

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
            action: { label: "Відкрити →", onClick: async () => await goto(`/permits/${ev.id}`) },
            id: String(ev.id),
          });
        } else if (ev.type === "fuzzy_match") {
          toast.warning("Знайдено нечіткий збіг номера", {
            description: `Розпізнано: ${ev.incoming_plate}. Знайдено кандидатів: ${ev.candidates?.length || 0}`,
            duration: 15000,
            action: { label: "Переглянути →", onClick: () => { 
                fuzzyMatchData = ev;
                showFuzzyModal = true; 
            } },
            id: `fuzzy-${ev.event_id}`,
          });
        }
      } catch {}
    };
  }

  onMount(() => {
    if (data.user?.customs_post_id) connectNotifications();
  });

  onDestroy(() => sse?.close());
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>TruckGuard</title>
</svelte:head>
<ModeWatcher defaultMode="light" />
<Toaster />

<FuzzyMatchModal bind:open={showFuzzyModal} data={fuzzyMatchData} onClose={() => { fuzzyMatchData = null; }} />

{#if data.user}
  <Sidebar.Provider>
    <AppSidebar user={data.user} />
    <Sidebar.Inset class="h-svh overflow-hidden flex flex-col">
      <header class="flex h-10 shrink-0 items-center gap-2 border-b px-3 bg-background z-30">
        <Sidebar.Trigger class="-ms-1 h-7 w-7" />
        <Separator orientation="vertical" class="me-1 h-4" />

        {#if crumbs.length > 0}
          <Breadcrumb.Root>
            <Breadcrumb.List class="flex-nowrap">
              {#each crumbs as crumb, i}
                <Breadcrumb.Item>
                  {#if crumb.href}
                    <Breadcrumb.Link
                      href={crumb.href}
                      class="text-xs text-muted-foreground hover:text-foreground transition-colors"
                    >
                      {crumb.label}
                    </Breadcrumb.Link>
                  {:else}
                    <Breadcrumb.Page class="text-xs font-medium text-foreground">
                      {crumb.label}
                    </Breadcrumb.Page>
                  {/if}
                </Breadcrumb.Item>
                {#if i < crumbs.length - 1}
                  <Breadcrumb.Separator class="text-muted-foreground/40" />
                {/if}
              {/each}
            </Breadcrumb.List>
          </Breadcrumb.Root>
        {/if}
      </header>

      <div class="flex-1 flex flex-col min-h-0 overflow-auto p-4 pt-3">
        {@render children()}
      </div>
    </Sidebar.Inset>
  </Sidebar.Provider>
{:else}
  {@render children()}
{/if}
