<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { ExternalLink, Trash2, Building2 } from "@lucide/svelte";
    import { can } from "$lib/auth";

    let { companies, flex = false, currentUser = null, onDelete } = $props<{
        companies: any[];
        flex?: boolean;
        currentUser?: any;
        onDelete: (company: any) => void;
    }>();

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Назва</Table.Head>
        <Table.Head class={th}>Код (ЄДРПОУ)</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(company: any)}
    <Table.Cell class="px-3 py-0 font-medium text-sm" style="height: var(--row-h);">{company.name}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs font-mono text-muted-foreground">{company.edrpou}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        <div class="flex items-center justify-end gap-0.5">
            {#if can(currentUser, "read:data")}
                <Button variant="ghost" size="icon" class="size-7" href="/config/data/companies/{company.ID}" title="Деталі">
                    <ExternalLink class="h-3.5 w-3.5" />
                </Button>
            {/if}
            {#if can(currentUser, "delete:data")}
                <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(company)} title="Видалити">
                    <Trash2 class="h-3.5 w-3.5" />
                </Button>
            {/if}
        </div>
    </Table.Cell>
{/snippet}

<DataTable columns={3} items={companies} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Building2} />
