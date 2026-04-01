<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import { ExternalLink, Trash2, Building2 } from "@lucide/svelte";

    let { companies, flex = false, onDelete } = $props<{
        companies: any[];
        flex?: boolean;
        onDelete: (company: any) => void;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Назва</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Код (ЄДРПОУ)</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(company: any)}
    <Table.Cell class="py-2.5 px-4 font-medium text-sm">{company.name}</Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm font-mono text-muted-foreground"
        >{company.edrpou}</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        <Button
            variant="ghost"
            size="icon"
            class="size-8"
            href="/config/data/companies/{company.ID}"
            title="Деталі"
        >
            <ExternalLink class="h-3.5 w-3.5" />
        </Button>
        <Button
            variant="ghost"
            size="icon"
            class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onclick={() => onDelete(company)}
            title="Видалити"
        >
            <Trash2 class="h-3.5 w-3.5" />
        </Button>
    </Table.Cell>
{/snippet}

<DataTable
    columns={3}
    items={companies}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Building2}
/>
