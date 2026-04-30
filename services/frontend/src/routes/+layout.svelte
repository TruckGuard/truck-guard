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
  import NotificationDrawer from "$lib/components/NotificationDrawer.svelte";

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
    profile:           "Профіль",
    plate:             "Номер",
    weight:            "Вага",
    system:            "Система"
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
  
  let notifications = $state<any[]>([]);
  let unreadCount = $derived(notifications.filter((n: any) => !n.is_read).length);
  let loadingMore = $state(false);
  let hasMore = $state(false);

  $effect(() => {
    if (data.notifications) {
      notifications = data.notifications;
      hasMore = data.notifications.length === 50;
    }
  });

  async function loadMoreNotifications() {
    if (loadingMore || !hasMore) return;
    loadingMore = true;
    try {
      const resp = await fetch(`/api/notifications?limit=50&offset=${notifications.length}`);
      if (!resp.ok) throw new Error("Failed to fetch");
      const res = await resp.json();
      const newNotifs = res.data || [];
      if (newNotifs.length < 50) {
        hasMore = false;
      }
      notifications = [...notifications, ...newNotifs];
    } catch (e) {
      console.error("Failed to load more notifications", e);
    } finally {
      loadingMore = false;
    }
  }

  function handleNotificationClick(n: any) {
    // Attempt to parse payload if it exists (for SSE events stored in DB)
    let ev = typeof n.payload === 'string' ? JSON.parse(n.payload) : n.payload;
    if (!ev) {
      // In case payload wasn't stored, but we have fields
      ev = { type: n.type };
    }
    
    if (n.type === "new_permit") {
      goto(`/permits/${ev.id || n.id}`);
    } else if (n.type === "fuzzy_match") {
      fuzzyMatchData = ev;
      showFuzzyModal = true;
    }
  }

  function connectNotifications() {
    sse = new EventSource("/api/notifications/stream");
    sse.onmessage = (e) => {
      try {
        const ev = JSON.parse(e.data);
        if (ev.type === "connected") return;
        
        // Add to persistent list
        const notif = {
          ID: ev.notification_id || Math.floor(Math.random() * 10000),
          type: ev.type,
          message: ev.type === 'new_permit' 
            ? `Нова перепустка: ${ev.plate || ev.code}` 
            : `Нечіткий збіг: ${ev.incoming_plate}`,
          payload: ev,
          is_read: false,
          CreatedAt: new Date().toISOString()
        };
        notifications = [notif, ...notifications];

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

        <div class="ml-auto flex items-center gap-2">
          <NotificationDrawer 
            {notifications} 
            {unreadCount} 
            {hasMore}
            {loadingMore}
            onNotificationClick={handleNotificationClick} 
            onLoadMore={loadMoreNotifications}
          />
        </div>
      </header>

      <div class="flex-1 flex flex-col min-h-0 overflow-auto p-4 pt-3">
        {@render children()}
      </div>
    </Sidebar.Inset>
  </Sidebar.Provider>
{:else}
  {@render children()}
{/if}
