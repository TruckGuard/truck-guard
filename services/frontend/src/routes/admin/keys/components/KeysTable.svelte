<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import type { APIKey } from "$lib/server/auth-client";
    import Pencil from "@lucide/svelte/icons/pencil";
    import Trash2 from "@lucide/svelte/icons/trash-2";
    import ShieldCheck from "@lucide/svelte/icons/shield-check";
    import Key from "@lucide/svelte/icons/key";

    let { keys, flex = false, onEdit, onDelete, onPermissions } = $props<{
        keys: APIKey[];
        flex?: boolean;
        onEdit: (key: APIKey) => void;
        onDelete: (key: APIKey) => void;
        onPermissions: (key: APIKey) => void;
    }>();

    function formatDate(date: string | Date | undefined) {
        if (!date) return "—";
        return new Date(date).toLocaleDateString("uk-UA", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        });
    }

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Власник</Table.Head>
        <Table.Head class={th}>Статус</Table.Head>
        <Table.Head class={th}>Створено</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(apiKey: APIKey)}
    <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
        <div class="flex flex-col justify-center gap-0.5">
            <span class="font-medium text-sm">{apiKey.owner_name}</span>
            <span class="text-[10px] text-muted-foreground font-mono">ID: {apiKey.id}</span>
        </div>
    </Table.Cell>
    <Table.Cell class="px-3 py-0">
        {#if apiKey.is_active}
            <span class="badge-success inline-flex items-center h-5 px-1.5 rounded text-[10px] font-medium">Активний</span>
        {:else}
            <span class="inline-flex items-center h-5 px-1.5 rounded border border-border text-[10px] font-medium text-muted-foreground bg-muted/40">Відкликаний</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs font-mono tabular-nums text-muted-foreground">
        {formatDate(apiKey.created_at)}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        <div class="flex items-center justify-end gap-0.5">
            <Button variant="ghost" size="icon" class="size-7" onclick={() => onPermissions(apiKey)} title="Дозволи">
                <ShieldCheck class="h-3.5 w-3.5" />
            </Button>
            <Button variant="ghost" size="icon" class="size-7" onclick={() => onEdit(apiKey)} title="Редагувати">
                <Pencil class="h-3.5 w-3.5" />
            </Button>
            <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(apiKey)} title="Видалити">
                <Trash2 class="h-3.5 w-3.5" />
            </Button>
        </div>
    </Table.Cell>
{/snippet}

<DataTable columns={4} items={keys} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Key} />
