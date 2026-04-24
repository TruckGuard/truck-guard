<script lang="ts">
    import { goto, invalidateAll } from "$app/navigation";
    import { ChevronLeft, RefreshCw, Link } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    let { permitCode, permitId, onclickLink, canUpdate = false } = $props<{
        permitCode?: string;
        permitId?: number;
        onclickLink?: () => void;
        canUpdate?: boolean;
    }>();

    let refreshing = $state(false);

    async function handleRefresh() {
        refreshing = true;
        await invalidateAll();
        refreshing = false;
    }
</script>

<div class="flex items-center gap-3 pb-4 border-b border-border mb-4">
    <Button
        variant="ghost"
        size="icon"
        onclick={() => goto("/permits")}
        class="shrink-0 h-8 w-8"
    >
        <ChevronLeft class="h-4 w-4" />
    </Button>
    <h1 class="text-lg font-semibold tracking-tight flex-1">
        Перепустка <span class="font-mono text-primary">{permitCode || `#${permitId}`}</span>
    </h1>
    <Button
        variant="ghost"
        size="sm"
        onclick={handleRefresh}
        disabled={refreshing}
        class="h-8 px-3 text-xs gap-1.5 text-muted-foreground hover:text-foreground"
    >
        <RefreshCw class="h-3.5 w-3.5 {refreshing ? 'animate-spin' : ''}" />
        Оновити
    </Button>
    {#if onclickLink && canUpdate}
        <Button
            variant="outline"
            size="sm"
            onclick={onclickLink}
            class="h-8 px-3 text-xs gap-1.5"
        >
            <Link class="h-3.5 w-3.5" />
            Прив'язати
        </Button>
    {/if}
</div>
