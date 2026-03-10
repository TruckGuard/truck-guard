<script lang="ts">
    let { isClosed = false } = $props<{ isClosed: boolean }>();
</script>

<div class="flex items-center gap-2">
    {#if isClosed}
        <span class="relative flex h-2 w-2">
            <span class="absolute inline-flex h-full w-full rounded-full bg-slate-300 dark:bg-slate-700 opacity-75"></span>
            <span class="relative inline-flex h-2 w-2 rounded-full bg-slate-400 dark:bg-slate-500"></span>
        </span>
        <span class="text-xs font-medium text-slate-500 dark:text-slate-400 leading-none">Виїхав</span>
    {:else}
        <span class="relative flex h-2 w-2">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex h-2 w-2 animate-pulse rounded-full bg-emerald-500"></span>
        </span>
        <span class="text-xs font-bold text-emerald-700 dark:text-emerald-400 leading-none">В зоні</span>
    {/if}
</div>
