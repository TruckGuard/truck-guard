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

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Назва</Table.Head>
        <Table.Head class={th}>Опис</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(role: Role)}
    <Table.Cell class="px-3 py-0 font-medium text-sm" style="height: var(--row-h);">{role.name}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs text-muted-foreground">{role.description || "—"}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        <div class="flex items-center justify-end gap-0.5">
            {#if authCan(currentUser, "update:roles")}
                <Button variant="ghost" size="icon" class="size-7 text-[color:var(--status-info)] hover:bg-[color:var(--status-info-bg)]" onclick={() => onPerms(role)} title="Права доступу">
                    <Shield class="h-3.5 w-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="size-7" onclick={() => onEdit(role)} title="Редагувати">
                    <Pencil class="h-3.5 w-3.5" />
                </Button>
            {/if}
            {#if authCan(currentUser, "delete:roles")}
                <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(role)} title="Видалити">
                    <Trash2 class="h-3.5 w-3.5" />
                </Button>
            {/if}
        </div>
    </Table.Cell>
{/snippet}

<DataTable columns={3} items={roles} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Users} />
