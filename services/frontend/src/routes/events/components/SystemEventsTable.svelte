<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { formatDate } from "$lib/utils/date";
    import { Activity, Database } from "@lucide/svelte";

    let { items, flex = false } = $props<{
        items: any[];
        flex?: boolean;
    }>();

    function parsePayload(payload: any) {
        if (!payload) return null;
        if (typeof payload !== "string") return payload;
        try {
            return JSON.parse(payload);
        } catch (e) {
            return null;
        }
    }
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >ID</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Час</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Категорія</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Дані</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(item: any)}
    {@const data = parsePayload(item.payload)}
    <Table.Cell class="py-2.5 px-4 font-mono">
        <a
            href="/events/system/{item.ID}"
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
            <Activity class="h-3.5 w-3.5 text-muted-foreground/60" />
            <span class="font-medium text-foreground/80"
                >{item.type || "Системна подія"}</span
            >
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4">
        <div class="max-w-[400px]">
            {#if data && typeof data === "object"}
                <code
                    class="text-xs bg-muted px-2 py-1.5 rounded text-muted-foreground block line-clamp-3 break-all leading-relaxed"
                    title={JSON.stringify(data)}
                >
                    {JSON.stringify(data)}
                </code>
            {:else}
                <span class="text-muted-foreground/50 italic">-</span>
            {/if}
        </div>
    </Table.Cell>
{/snippet}

<DataTable
    columns={4}
    {items}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Database}
/>
