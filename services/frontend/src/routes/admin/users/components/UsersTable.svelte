<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Badge } from "$lib/components/ui/badge";
    import { Button } from "$lib/components/ui/button";
    import { format } from "date-fns";
    import { can } from "$lib/auth";
    import { Pencil, Trash2, Users, Key } from "@lucide/svelte";

    let { users, currentUser, posts, flex = false, onEdit, onDelete, onResetPassword } = $props<{
        users: any[];
        currentUser: any;
        posts?: any[];
        flex?: boolean;
        onEdit: (user: any) => void;
        onDelete: (user: any) => void;
        onResetPassword: (user: any) => void;
    }>();
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Username</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Роль</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >ПІБ</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Email</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Телефон</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Пост</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b"
            >Останній вхід</Table.Head
        >
        <Table.Head
            class="h-12 px-4 py-2 font-bold text-sm uppercase tracking-wider text-muted-foreground border-b text-right"
            >Дії</Table.Head
        >
    </Table.Row>
{/snippet}

{#snippet row(user: any)}
    <Table.Cell class="py-2.5 px-4">
        <div class="flex items-center gap-2">
            <span class="font-medium">{user.username}</span>
            {#if user.id === currentUser.id}
                <Badge
                    variant="outline"
                    class="text-xs h-4.5 px-1.5 border-primary/30 text-primary"
                    >Ви</Badge
                >
            {/if}
        </div>
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4">
        <Badge
            variant="outline"
            class="font-mono text-xs uppercase tracking-wider bg-muted/50"
            >{user.role?.name || "No Role"}</Badge
        >
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm">
        {#if user.profile}
            {user.profile.last_name || ""}
            {user.profile.first_name || ""}
            {user.profile.third_name || ""}
        {:else}
            <span class="text-muted-foreground italic text-sm">Профіль не знайдено</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm text-muted-foreground">
        {user.profile?.email || "-"}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-sm text-muted-foreground">
        {user.profile?.phone_number || "-"}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4">
        {#if user.profile?.customs_post_id}
            {@const post = posts?.find((p: any) => p.ID === user.profile.customs_post_id)}
            <span
                class="text-sm font-medium px-2 py-0.5 rounded-full bg-muted border border-border"
            >
                {post?.name || user.profile.customs_post_id}
            </span>
        {:else}
            -
        {/if}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 whitespace-nowrap text-sm tabular-nums text-muted-foreground">
        {#if user.last_login}
            {format(new Date(user.last_login), "dd.MM.yyyy HH:mm")}
        {:else}
            <span class="text-muted-foreground opacity-50">Ніколи</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="py-2.5 px-4 text-right space-x-1">
        {#if can(currentUser, "update:users") && user.profile && user.id != currentUser.id}
            <Button
                variant="ghost"
                size="icon"
                class="size-8"
                href={`/admin/users/${user.id}`}
                title="Змінити"
            >
                <Pencil class="h-3.5 w-3.5" />
            </Button>
        {/if}
        {#if can(currentUser, "manage:users") && user.id != currentUser.id}
            <Button
                variant="ghost"
                size="icon"
                class="size-8 text-amber-500 hover:text-amber-600 hover:bg-amber-500/10"
                onclick={() => onResetPassword(user)}
                title="Скинути пароль"
            >
                <Key class="h-3.5 w-3.5" />
            </Button>
        {/if}
        {#if can(currentUser, "delete:users") && user.id != currentUser.id}
            <Button
                variant="ghost"
                size="icon"
                class="size-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                onclick={() => onDelete(user)}
                title="Видалити"
            >
                <Trash2 class="h-3.5 w-3.5" />
            </Button>
        {/if}
    </Table.Cell>
{/snippet}

<DataTable
    columns={8}
    items={users}
    {flex}
    headerSnippet={header}
    rowSnippet={row}
    emptyStateIcon={Users}
/>
