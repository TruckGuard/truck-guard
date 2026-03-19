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

<div class="flex items-center gap-6 mb-8">
    <Button
        variant="outline"
        size="icon"
        onclick={() => goto("/permits")}
        class="shrink-0 rounded-full h-10 w-10 border-slate-200 dark:border-slate-700 shadow-sm hover:bg-slate-50 dark:hover:bg-slate-800"
    >
        <ChevronLeft class="h-5 w-5" />
    </Button>
    <h1 class="text-3xl font-extrabold tracking-tight flex-1">
        Перепустка {permitCode || `#${permitId}`}
    </h1>
    <Button
        variant="ghost"
        size="sm"
        onclick={handleRefresh}
        disabled={refreshing}
        class="gap-2 text-muted-foreground hover:text-foreground"
    >
        <RefreshCw class="h-4 w-4 {refreshing ? 'animate-spin' : ''}" />
        Оновити
    </Button>
    {#if onclickLink && canUpdate}
        <Button
            variant="outline"
            size="sm"
            onclick={onclickLink}
            class="gap-2 border-blue-200 text-blue-700 hover:bg-blue-50 hover:text-blue-800 dark:border-blue-900/30 dark:text-blue-400 dark:hover:bg-blue-900/20"
        >
            <Link class="h-4 w-4" />
            Прив'язати
        </Button>
    {/if}
</div>
