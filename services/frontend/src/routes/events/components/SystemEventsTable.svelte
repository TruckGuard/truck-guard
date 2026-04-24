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
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">ID</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Час</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Категорія</Table.Head>
        <Table.Head class="h-[var(--row-h)] px-3 py-0 text-xs font-medium uppercase tracking-wider text-muted-foreground border-b">Дані</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(item: any)}
    {@const data = parsePayload(item.payload)}
    <Table.Cell class="px-3 font-mono text-xs text-muted-foreground">
        <a
            href="/events/system/{item.ID}"
            class="text-primary hover:underline font-medium underline-offset-4"
        >
            #{item.ID}
        </a>
    </Table.Cell>
    <Table.Cell class="px-3 text-xs tabular-nums text-muted-foreground">{formatDate(item.timestamp)}</Table.Cell>
    <Table.Cell class="px-3">
        <div class="flex items-center gap-1.5">
            <Activity class="h-3 w-3 text-muted-foreground/50" />
            <span class="text-xs text-foreground/80">{item.type || "Системна подія"}</span>
        </div>
    </Table.Cell>
    <Table.Cell class="px-3">
        <div class="max-w-[400px]">
            {#if data && typeof data === "object"}
                <code class="text-xs font-mono bg-muted/50 px-1.5 py-0.5 rounded text-muted-foreground block truncate" title={JSON.stringify(data)}>
                    {JSON.stringify(data)}
                </code>
            {:else}
                <span class="text-muted-foreground/40">—</span>
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
