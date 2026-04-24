<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import DataTable from "$lib/components/common/DataTable.svelte";
    import { Button } from "$lib/components/ui/button";
    import { Pencil, Trash2, MapPin } from "@lucide/svelte";
    import { can } from "$lib/auth";

    let { posts, flex = false, currentUser = null, onEdit, onDelete } = $props<{
        posts: any[];
        flex?: boolean;
        currentUser?: any;
        onEdit: (post: any) => void;
        onDelete: (post: any) => void;
    }>();

    const th = "h-[var(--row-h)] px-3 py-0 font-medium text-xs uppercase tracking-wider text-muted-foreground border-b";
</script>

{#snippet header()}
    <Table.Row class="hover:bg-transparent">
        <Table.Head class={th}>Назва</Table.Head>
        <Table.Head class={th}>Опис</Table.Head>
        <Table.Head class={th}>Створено</Table.Head>
        <Table.Head class="{th} text-right">Дії</Table.Head>
    </Table.Row>
{/snippet}

{#snippet row(post: any)}
    <Table.Cell class="px-3 py-0 font-medium text-sm" style="height: var(--row-h);">{post.name}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs text-muted-foreground">{post.description || "—"}</Table.Cell>
    <Table.Cell class="px-3 py-0 text-xs tabular-nums font-mono text-muted-foreground">
        {new Date(post.CreatedAt).toLocaleDateString("uk-UA")}
    </Table.Cell>
    <Table.Cell class="px-3 py-0 text-right">
        <div class="flex items-center justify-end gap-0.5">
            {#if can(currentUser, "update:data")}
                <Button variant="ghost" size="icon" class="size-7" onclick={() => onEdit(post)} title="Редагувати">
                    <Pencil class="h-3.5 w-3.5" />
                </Button>
            {/if}
            {#if can(currentUser, "delete:data")}
                <Button variant="ghost" size="icon" class="size-7 text-destructive hover:bg-destructive/8 hover:text-destructive" onclick={() => onDelete(post)} title="Видалити">
                    <Trash2 class="h-3.5 w-3.5" />
                </Button>
            {/if}
        </div>
    </Table.Cell>
{/snippet}

<DataTable columns={4} items={posts} {flex} headerSnippet={header} rowSnippet={row} emptyStateIcon={MapPin} />
