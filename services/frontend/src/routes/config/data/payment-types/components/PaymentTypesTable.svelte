<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Pencil, Trash2, Database } from "@lucide/svelte";
    import { can } from "$lib/auth";

    let { paymentTypes, flex = false, currentUser = null, onEdit, onDelete } = $props<{
        paymentTypes: any[];
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
        <Table.Head class={th}>Статус</Table.Head>
        <Table.Head class={th}>Опис</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(type: any)}
    <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
        <span class="inline-flex items-center h-5 px-1.5 rounded border border-border font-mono text-[10px] text-muted-foreground bg-muted/40">
            {type.code}
        </span>
    </Table.Cell>
    <Table.Cell class="px-3 py-0 font-medium text-sm">{type.name}</Table.Cell>
    <Table.Cell class="px-3 py-0">
        {#if type.is_active}
            <span class="badge-success inline-flex items-center h-5 px-1.5 rounded text-[10px] font-medium">Активний</span>
        {:else}
            <span class="inline-flex items-center h-5 px-1.5 rounded text-[10px] font-medium text-muted-foreground bg-muted/40 border border-border">Неактивний</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs text-muted-foreground">{type.description || "—"}</Table.Cell>
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

<DataTable columns={5} items={paymentTypes} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Database} />
