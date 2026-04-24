<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Trash2, Car } from "@lucide/svelte";
    import { can } from "$lib/auth";

    let { plates, flex = false, currentUser = null, onDelete } = $props<{
        plates: any[];
        flex?: boolean;
        currentUser?: any;
        onDelete: (plate: any) => void;
    }>();

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Номер</Table.Head>
        <Table.Head class={th}>Коментар</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(plate: any)}
    <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
        <span class="plate">{plate.plate}</span>
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs text-muted-foreground">{plate.comment || '—'}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        {#if can(currentUser, "delete:settings")}
            <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(plate)} title="Видалити">
                <Trash2 class="h-3.5 w-3.5" />
            </Button>
        {/if}
    </Table.Cell>
{/snippet}

<DataTable columns={3} items={plates} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Car} />
