<script lang="ts">
    import * as Table from "$lib/components/ui/table";
    import { Button } from "$lib/components/ui/button";
    import { Pencil, Trash2, MapPin } from "@lucide/svelte";

    let { posts, onEdit, onDelete } = $props<{
        posts: any[];
        onEdit: (post: any) => void;
        onDelete: (post: any) => void;
    }>();
</script>

<div class="rounded-md border">
    <Table.Root>
        <Table.Header>
            <Table.Row>
                <Table.Head>Назва</Table.Head>
                <Table.Head>Опис</Table.Head>
                <Table.Head>Дата створення</Table.Head>
                <Table.Head class="text-right">Дії</Table.Head>
            </Table.Row>
        </Table.Header>
        <Table.Body>
            {#if posts && posts.length > 0}
                {#each posts as post (post.ID)}
                    <Table.Row>
                        <Table.Cell class="font-medium">{post.name}</Table.Cell>
                        <Table.Cell>{post.description || "-"}</Table.Cell>
                        <Table.Cell>
                            {new Date(post.CreatedAt).toLocaleDateString(
                                "uk-UA",
                            )}
                        </Table.Cell>
                        <Table.Cell class="text-right space-x-2">
                            <Button
                                variant="ghost"
                                size="icon"
                                onclick={() => onEdit(post)}
                            >
                                <Pencil class="h-4 w-4" />
                            </Button>
                            <Button
                                variant="ghost"
                                size="icon"
                                class="text-destructive hover:text-destructive"
                                onclick={() => onDelete(post)}
                            >
                                <Trash2 class="h-4 w-4" />
                            </Button>
                        </Table.Cell>
                    </Table.Row>
                {/each}
            {:else}
                <Table.Row>
                    <Table.Cell colspan={4} class="h-24 text-center">
                        <div
                            class="flex flex-col items-center justify-center gap-2"
                        >
                            <MapPin class="h-8 w-8 opacity-20" />
                            <span>Результатів не знайдено.</span>
                        </div>
                    </Table.Cell>
                </Table.Row>
            {/if}
        </Table.Body>
    </Table.Root>
</div>
