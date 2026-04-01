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
        if (!date) return "-";
        return new Date(date).toLocaleDateString("uk-UA", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        });
    }
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Власник</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Статус</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b"
            >Створено</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(apiKey: APIKey)}
    <Table.Cell class="py-2.5 px-4">
        <div class="flex flex-col">
            <span class="font-medium text-sm">{apiKey.owner_name}</span>
            <span class="text-xs text-muted-foreground font-mono"
                >ID: {apiKey.id}</span
            >
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4">
        <span
            class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium {apiKey.is_active
                ? 'bg-emerald-500/10 text-emerald-500'
                : 'bg-muted text-muted-foreground'}"
        >
            {apiKey.is_active ? "Активний" : "Відкликаний"}
        </span>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-xs tabular-nums text-muted-foreground">
        {formatDate(apiKey.created_at)}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        <Button
            variant="ghost"
            size="icon"
            class="size-8"
            onclick={() => onPermissions(apiKey)}
            title="Дозволи"
        >
            <ShieldCheck class="h-3.5 w-3.5" />
        </Button>
        <Button
            variant="ghost"
            size="icon"
            class="size-8"
            onclick={() => onEdit(apiKey)}
            title="Редагувати"
        >
            <Pencil class="h-3.5 w-3.5" />
        </Button>
        <Button
            variant="ghost"
            size="icon"
            class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            onclick={() => onDelete(apiKey)}
            title="Видалити"
        >
            <Trash2 class="h-3.5 w-3.5" />
        </Button>
    </Table.Cell>
{/snippet}

<DataTable
    columns={4}
    items={keys}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Key}
/>
