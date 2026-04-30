<script lang="ts">
  import { Bell } from "@lucide/svelte";
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as Sheet from "$lib/components/ui/sheet";
  import { ScrollArea } from "$lib/components/ui/scroll-area";
  import { enhance } from "$app/forms";
  import { formatDistanceToNow } from "date-fns";
  import { uk } from "date-fns/locale";

  let { 
    notifications = [], 
    unreadCount = 0, 
    hasMore = false,
    loadingMore = false,
    onNotificationClick,
    onLoadMore
  } = $props<{
    notifications: any[];
    unreadCount: number;
    hasMore: boolean;
    loadingMore: boolean;
    onNotificationClick: (n: any) => void;
    onLoadMore: () => void;
  }>();

  let open = $state(false);
  let observerTarget = $state<HTMLElement | null>(null);

  function handleNotificationClick(n: any) {
    open = false; // Close drawer
    onNotificationClick(n);
  }

  $effect(() => {
    if (!observerTarget || !hasMore || loadingMore || !open) return;

    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        onLoadMore();
      }
    }, { threshold: 0.1 });

    observer.observe(observerTarget);
    return () => observer.disconnect();
  });
</script>

<Sheet.Root bind:open>
  <Sheet.Trigger class={buttonVariants({ variant: "ghost", size: "icon", className: "relative h-8 w-8 rounded-full" })}>
      <Bell class="h-4 w-4" />
      {#if unreadCount > 0}
        <span class="absolute top-1 right-1 h-4 w-4 rounded-full bg-red-500 text-[10px] font-medium text-white flex items-center justify-center -translate-y-1/4 translate-x-1/4">
          {unreadCount > 99 ? '99+' : unreadCount}
        </span>
      {/if}
      <span class="sr-only">Відкрити сповіщення</span>
  </Sheet.Trigger>
  <Sheet.Content class="w-[400px] sm:w-[540px] flex flex-col gap-0 p-0">
    <Sheet.Header class="px-6 py-4 border-b shrink-0 flex flex-row items-center justify-between">
      <div>
        <Sheet.Title class="text-lg">Сповіщення</Sheet.Title>
      </div>
      {#if unreadCount > 0}
        <form action="/notifications?/markAllRead" method="POST" use:enhance={() => {
          notifications.forEach((n: any) => n.is_read = true);
          return async ({ update }: any) => {
            await update({ reset: false });
          };
        }}>
          <Button type="submit" variant="outline" size="sm" class="h-8 text-xs">
            Позначити всі як прочитані
          </Button>
        </form>
      {/if}
    </Sheet.Header>

    <ScrollArea class="flex-1 h-full min-h-0">
      <div class="flex flex-col">
        {#if notifications.length === 0}
          <div class="flex flex-col items-center justify-center p-8 text-center text-muted-foreground">
            <Bell class="h-8 w-8 mb-2 opacity-20" />
            <p class="text-sm">Немає сповіщень</p>
          </div>
        {:else}
          {#each notifications as n}
            <form action="/notifications?/markRead" method="POST" use:enhance={() => {
              n.is_read = true;
              return async ({ update }) => {
                await update({ reset: false });
              };
            }}>
              <input type="hidden" name="id" value={n.ID} />
              <button 
                type="submit"
                class="w-full text-left flex flex-col gap-1 p-4 border-b cursor-pointer hover:bg-muted/50 transition-colors {n.is_read ? 'opacity-70' : 'bg-primary/5'}"
                onclick={() => handleNotificationClick(n)}
              >
                <div class="flex items-center justify-between w-full">
                  <span class="text-xs font-semibold text-primary">{n.type === 'fuzzy_match' ? 'Нечіткий збіг' : 'Нова перепустка'}</span>
                  <span class="text-[10px] text-muted-foreground">
                    {formatDistanceToNow(new Date(n.CreatedAt), { addSuffix: true, locale: uk })}
                  </span>
                </div>
                <p class="text-sm font-medium w-full">{n.message}</p>
              </button>
            </form>
          {/each}
          
          {#if hasMore}
            <div bind:this={observerTarget} class="p-4 flex items-center justify-center">
              {#if loadingMore}
                <div class="h-4 w-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
              {/if}
            </div>
          {/if}
        {/if}
      </div>
    </ScrollArea>
  </Sheet.Content>
</Sheet.Root>
