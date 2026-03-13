<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { formatDate } from "$lib/utils/date";
    import { Scale, Database } from "@lucide/svelte";

    let { items, flex = false } = $props<{
        items: any[];
        flex?: boolean;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >ID</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Час</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Обладнання</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b text-right"
            >Вага</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(item: any)}
    <Table.Cell class="py-2.5 px-4 font-mono">
        <a 
            href="/events/weight/{item.ID}" 
            class="text-primary hover:underline font-bold decoration-primary/30 underline-offset-4"
        >
            #{item.ID}
        </a>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 tabular-nums"
        >{formatDate(item.timestamp)}</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4">
        <div class="flex items-center gap-2">
            <Scale class="h-5 w-5 text-muted-foreground/60" />
            <span class="font-medium text-foreground/80">{item.scale_source_name || item.scale_name || "Невідомо"}</span>
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-right">
        <span class="font-bold tabular-nums text-lg">{item.weight} <small class="text-muted-foreground font-normal text-xs">кг</small></span>
    </Table.Cell>
{/snippet}

<DataTable
    columns={4}
    items={items}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Database}
/>
