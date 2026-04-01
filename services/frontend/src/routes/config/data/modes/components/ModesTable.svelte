<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Badge } from "$lib/components/ui/badge";
    import { Button } from "$lib/components/ui/button";
    import { Pencil, Trash2, Database } from "@lucide/svelte";

    let { modes, flex = false, onEdit, onDelete } = $props<{
        modes: any[];
        flex?: boolean;
        onEdit: (mode: any) => void;
        onDelete: (mode: any) => void;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Код</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Назва</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Опис</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(mode: any)}
    <Table.Cell class="py-2.5 px-4 font-mono text-xs">
        <Badge variant="outline">{mode.code}</Badge>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 font-medium text-sm">{mode.name}</Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm text-muted-foreground"
        >{mode.description || "-"}</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        <Button
            variant="ghost"
            size="icon"
            class="size-8"
            onclick={() => onEdit(mode)}
            title="Редагувати"
        >
            <Pencil class="h-3.5 w-3.5" />
        </Button>
        <Button
            variant="ghost"
            size="icon"
            class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onclick={() => onDelete(mode)}
            title="Видалити"
        >
            <Trash2 class="h-3.5 w-3.5" />
        </Button>
    </Table.Cell>
{/snippet}

<DataTable
    columns={5}
    items={modes}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Database}
/>
