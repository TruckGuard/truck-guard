<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
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

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Username</Table.Head>
        <Table.Head class={th}>Роль</Table.Head>
        <Table.Head class={th}>ПІБ</Table.Head>
        <Table.Head class={th}>Email</Table.Head>
        <Table.Head class={th}>Телефон</Table.Head>
        <Table.Head class={th}>Пост</Table.Head>
        <Table.Head class={th}>Вхід</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(user: any)}
    <Table.Cell class="px-3 py-0" style="height: var(--row-h);">
        <div class="flex items-center gap-2">
            <span class="font-medium text-sm">{user.username}</span>
            {#if user.id === currentUser.id}
                <span class="inline-flex items-center h-4.5 px-1.5 rounded border border-primary/30 text-primary text-[10px] font-medium">Ви</span>
            {/if}
        </div>
    </Table.Cell>
    <Table.Cell class="px-3 py-0">
        <span class="inline-flex items-center h-5 px-1.5 rounded border border-border font-mono text-[10px] uppercase tracking-wider text-muted-foreground bg-muted/40">
            {user.role?.name || "No Role"}
        </span>
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-sm">
        {#if user.profile}
            {user.profile.last_name || ""}
            {user.profile.first_name || ""}
            {user.profile.third_name || ""}
        {:else}
            <span class="text-muted-foreground italic text-xs">Профіль не знайдено</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs text-muted-foreground">
        {user.profile?.email || "—"}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs text-muted-foreground">
        {user.profile?.phone_number || "—"}
    </Table.Cell>
    <Table.Cell class="px-3 py-0">
        {#if user.profile?.customs_post_id}
            {@const post = posts?.find((p: any) => p.ID === user.profile.customs_post_id)}
            <span class="inline-flex items-center h-5 px-1.5 rounded border border-border text-xs font-medium bg-muted/40">
                {post?.name || user.profile.customs_post_id}
            </span>
        {:else}
            <span class="text-muted-foreground text-xs">—</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs font-mono tabular-nums text-muted-foreground whitespace-nowrap">
        {#if user.last_login}
            {format(new Date(user.last_login), "dd.MM.yyyy HH:mm")}
        {:else}
            <span class="opacity-50">Ніколи</span>
        {/if}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        <div class="flex items-center justify-end gap-0.5">
            {#if can(currentUser, "update:users") && user.profile && user.id != currentUser.id}
                <Button variant="ghost" size="icon" class="size-7" href={`/admin/users/${user.id}`} title="Змінити">
                    <Pencil class="h-3.5 w-3.5" />
                </Button>
            {/if}
            {#if can(currentUser, "manage:users") && user.id != currentUser.id}
                <Button variant="ghost" size="icon" class="size-7 text-[color:var(--status-warning)] hover:bg-[color:var(--status-warning-bg)]" onclick={() => onResetPassword(user)} title="Скинути пароль">
                    <Key class="h-3.5 w-3.5" />
                </Button>
            {/if}
            {#if can(currentUser, "delete:users") && user.id != currentUser.id}
                <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(user)} title="Видалити">
                    <Trash2 class="h-3.5 w-3.5" />
                </Button>
            {/if}
        </div>
    </Table.Cell>
{/snippet}

<DataTable columns={8} items={users} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={Users} />
