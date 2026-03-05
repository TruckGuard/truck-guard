<script lang="ts">
    import { goto, invalidateAll } from "$app/navigation";
    import { ChevronLeft, RefreshCw } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    let { isNew, permitCode, permitId } = $props<{
        isNew: boolean;
        permitCode?: string;
        permitId?: number;
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
        class="shrink-0 rounded-full h-10 w-10 border-slate-200 shadow-sm hover:bg-slate-50"
    >
        <ChevronLeft class="h-5 w-5" />
    </Button>
    <h1 class="text-3xl font-extrabold tracking-tight flex-1">
        {#if isNew}
            Створення нової перепустки
        {:else}
            Перепустка {permitCode || `#${permitId}`}
        {/if}
    </h1>
    {#if !isNew}
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
    {/if}
</div>
