<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Badge } from "$lib/components/ui/badge";
    import { Pencil, Trash2, Database } from "@lucide/svelte";
    import { darkenColor } from "$lib/utils/colors";

    let { vehicleTypes, flex = false, onEdit, onDelete } = $props<{
        vehicleTypes: any[];
        flex?: boolean;
        onEdit: (type: any) => void;
        onDelete: (type: any) => void;
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
            >Ціна в'їзду</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Ціна доби</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(type: any)}
    <Table.Cell class="py-2.5 px-4">
        <Badge
            variant="outline"
            style="background-color: {type.color}22; color: {type.color}; border-color: {darkenColor(
                type.color,
                20,
            )}"
        >
            {type.code}
        </Badge>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 font-medium text-sm">{type.name}</Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm font-medium tabular-nums"
        >{type.entry_price?.toLocaleString("uk-UA")} ₴</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 text-sm font-medium tabular-nums"
        >{type.daily_price?.toLocaleString("uk-UA")} ₴</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        <Button
            variant="ghost"
            size="icon"
            class="size-8"
            onclick={() => onEdit(type)}
            title="Редагувати"
        >
            <Pencil class="h-3.5 w-3.5" />
        </Button>
        <Button
            variant="ghost"
            size="icon"
            class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onclick={() => onDelete(type)}
            title="Видалити"
        >
            <Trash2 class="h-3.5 w-3.5" />
        </Button>
    </Table.Cell>
{/snippet}

<DataTable
    columns={5}
    items={vehicleTypes}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Database}
/>
