<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import type { Role } from "$lib/server/auth-client";
    import { can as authCan } from "$lib/auth";
    import Pencil from "@lucide/svelte/icons/pencil";
    import Trash2 from "@lucide/svelte/icons/trash-2";
    import Shield from "@lucide/svelte/icons/shield";
    import Users from "@lucide/svelte/icons/users";

    let { roles, currentUser, flex = false, onPerms, onEdit, onDelete } = $props<{
        roles: Role[];
        currentUser: any;
        flex?: boolean;
        onPerms: (role: Role) => void;
        onEdit: (role: Role) => void;
        onDelete: (role: Role) => void;
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
            >Опис</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-xs uppercase tracking-wider text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(role: Role)}
    <Table.Cell class="py-2.5 px-4 font-medium text-sm">{role.name}</Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm text-muted-foreground"
        >{role.description || "-"}</Table.Cell
    >
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        {#if authCan(currentUser, "update:roles")}
            <Button
                variant="ghost"
                size="icon"
                class="size-8"
                onclick={() => onPerms(role)}
                title="Права доступу"
            >
                <Shield class="h-3.5 w-3.5 text-blue-600 dark:text-blue-400" />
            </Button>
            <Button
                variant="ghost"
                size="icon"
                class="size-8"
                onclick={() => onEdit(role)}
                title="Редагувати"
            >
                <Pencil class="h-3.5 w-3.5" />
            </Button>
        {/if}
        {#if authCan(currentUser, "delete:roles")}
            <Button
                variant="ghost"
                size="icon"
                class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                onclick={() => onDelete(role)}
                title="Видалити"
            >
                <Trash2 class="h-3.5 w-3.5" />
            </Button>
        {/if}
    </Table.Cell>
{/snippet}

<DataTable
    columns={3}
    items={roles}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Users}
/>
