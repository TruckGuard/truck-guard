<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Trash2, Car } from "@lucide/svelte";

    let { plates, flex = false, onDelete } = $props<{
        plates: any[];
        flex?: boolean;
        onDelete: (plate: any) => void;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b"
            >Номер</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b"
            >Коментар</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-[10px] uppercase tracking-[0.2em] text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(plate: any)}
    <Table.Cell class="py-2.5 px-4 font-medium text-sm">
        <div class="flex items-center gap-2">
            <Car class="h-4 w-4 text-muted-foreground" />
            {plate.plate}
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm text-muted-foreground">
        {plate.comment || '-'}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-right">
        <Button
            variant="ghost"
            size="icon"
            class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onclick={() => onDelete(plate)}
            title="Видалити"
        >
            <Trash2 class="h-3.5 w-3.5" />
        </Button>
    </Table.Cell>
{/snippet}

<DataTable
    columns={3}
    items={plates}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Car}
/>
