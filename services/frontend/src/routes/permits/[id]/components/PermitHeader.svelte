<script lang="ts">
    import { goto } from "$app/navigation";
    import { ChevronLeft } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    let { isNew, permitCode, permitId } = $props<{
        isNew: boolean;
        permitCode?: string;
        permitId?: number;
    }>();
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
    <h1 class="text-3xl font-extrabold tracking-tight">
        {#if isNew}
            Створення нової перепустки
        {:else}
            Перепустка {permitCode || `#${permitId}`}
        {/if}
    </h1>
</div>
