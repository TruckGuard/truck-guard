<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Pencil, Trash2, Database } from "@lucide/svelte";
    import { darkenColor } from "$lib/utils/colors";
    import { can } from "$lib/auth";

    let { vehicleTypes, flex = false, currentUser = null, onEdit, onDelete } = $props<{
        vehicleTypes: any[];
        flex?: boolean;
        currentUser?: any;
        onEdit: (type: any) => void;
        onDelete: (type: any) => void;
    }>();

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Код</Table.Head>
        <Table.Head class={th}>Назва</Table.Head>
        <Table.Head class="{th} text-right">Ціна в'їзду</Table.Head>
        <Table.Head class="{th} text-right">Ціна доби</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(type: any)}
    <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
        <span
            class="inline-flex items-center h-5 px-1.5 rounded border font-mono font-medium text-[10px]"
            style="background-color: {type.color}18; color: {type.color}; border-color: {darkenColor(type.color, 20)}40"
        >
            {type.code}
        </span>
    </Table.Cell>
    <Table.Cell class="px-3 py-0 font-medium text-sm">{type.name}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs font-mono tabular-nums text-right">
        {type.entry_price?.toLocaleString("uk-UA")} ₴
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs font-mono tabular-nums text-right">
        {type.daily_price?.toLocaleString("uk-UA")} ₴
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        <div class="flex items-center justify-end gap-0.5">
            {#if can(currentUser, "update:data")}
                <Button variant="ghost" size="icon" class="size-7" onclick={() => onEdit(type)} title="Редагувати">
                    <Pencil class="h-3.5 w-3.5" />
                </Button>
            {/if}
            {#if can(currentUser, "delete:data")}
                <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(type)} title="Видалити">
                    <Trash2 class="h-3.5 w-3.5" />
                </Button>
            {/if}
        </div>
    </Table.Cell>
{/snippet}

<DataTable columns={5} items={vehicleTypes} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Database} />
