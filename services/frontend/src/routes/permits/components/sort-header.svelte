<script lang="ts">
    import { ArrowDown, ArrowUp, ArrowUpDown } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";

    // In Svelte 5, $props() defines component inputs
    let { title, sortState = false } = $props<{
        title: string;
        sortState: "asc" | "desc" | false;
    }>();
</script>

<div class="flex items-center">
    <span class="mr-2 font-medium">{title}</span>
    {#if sortState === "asc"}
        <ArrowUp class="h-3.5 w-3.5" />
    {:else if sortState === "desc"}
        <ArrowDown class="h-3.5 w-3.5" />
    {:else}
        <ArrowUpDown
            class="h-3.5 w-3.5 opacity-30 group-hover:opacity-100 transition-opacity"
        />
    {/if}
</div>
