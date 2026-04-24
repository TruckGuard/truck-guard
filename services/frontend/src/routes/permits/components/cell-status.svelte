<script lang="ts">
    let { isClosed = false, isVoid = false } = $props<{ isClosed: boolean, isVoid?: boolean }>();
</script>

<div class="inline-flex items-center gap-1.5">
    {#if isVoid}
        <span class="h-1.5 w-1.5 rounded-full bg-danger shrink-0"></span>
        <span class="badge-danger inline-flex items-center rounded px-1.5 py-0.5 text-xs font-medium uppercase tracking-wide">Анульовано</span>
    {:else if isClosed}
        <span class="h-1.5 w-1.5 rounded-full bg-muted-foreground/50 shrink-0"></span>
        <span class="text-xs text-muted-foreground font-medium">Виїхав</span>
    {:else}
        <span class="relative flex h-1.5 w-1.5 shrink-0">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-success opacity-75"></span>
            <span class="relative inline-flex h-1.5 w-1.5 rounded-full bg-success"></span>
        </span>
        <span class="text-xs font-medium" style="color: var(--status-success)">В зоні</span>
    {/if}
</div>
