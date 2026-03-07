<script lang="ts">
    import type { Company } from "$lib/types/data";

    let { payers = [] } = $props<{
        payers: { company?: Company }[];
    }>();
</script>

<div class="flex flex-wrap gap-1 max-w-[200px]">
    {#if payers && payers.length > 0}
        {#each payers as payer}
            {#if payer.company}
                <a
                    href={`/config/companies/${payer.company.ID}`}
                    class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-medium bg-secondary text-secondary-foreground border hover:bg-secondary/80 transition-colors cursor-pointer"
                    title="{payer.company.name} (ЄДРПОУ: {payer.company
                        .edrpou})"
                    onclick={(e) => e.stopPropagation()}
                >
                    {payer.company.name.length > 20
                        ? payer.company.name.substring(0, 20) + "..."
                        : payer.company.name}
                </a>
            {/if}
        {/each}
    {:else}
        <span class="text-muted-foreground italic text-xs">-</span>
    {/if}
</div>
