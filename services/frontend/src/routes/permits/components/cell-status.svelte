<script lang="ts">
    import { CircleCheck, CircleX } from "@lucide/svelte";

    let { isClosed = false } = $props<{ isClosed: boolean }>();
</script>

<div class="flex justify-end">
    {#if isClosed}
        <div
            class="inline-flex items-center gap-1 text-slate-600 bg-slate-100 dark:bg-slate-800/50 px-2 py-1 rounded-full text-xs font-semibold"
        >
            <CircleX class="h-3.5 w-3.5" /> Виїхав
        </div>
    {:else}
        <div
            class="inline-flex items-center gap-1 text-emerald-600 bg-emerald-50 dark:bg-emerald-950/30 px-2 py-1 rounded-full text-xs font-semibold"
        >
            <CircleCheck class="h-3.5 w-3.5" /> В зоні
        </div>
    {/if}
</div>
