<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, Search, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  // Component Imports
  import PostsTable from "./components/PostsTable.svelte";
  import PostCreateDialog from "./components/PostCreateDialog.svelte";
  import PostEditSheet from "./components/PostEditSheet.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedPost = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchQuery) url.searchParams.set("name", searchQuery);
    else url.searchParams.delete("name");

    url.searchParams.set("page", "1");
    goto(url);
  }

  function openEdit(post: any) {
    selectedPost = { ...post };
    isEditOpen = true;
  }

  function openDelete(post: any) {
    selectedPost = post;
    isDeleteOpen = true;
  }

  function handlePageChange(newPage: number) {
    const url = new URL(page.url);
    url.searchParams.set("page", newPage.toString());
    goto(url);
  }
</script>

<div class="flex flex-col h-full overflow-hidden space-y-6">
  {#if data.error}
    <Alert.Root variant="destructive" class="shrink-0">
      <CircleAlert class="h-4 w-4" />
      <Alert.Title>Помилка</Alert.Title>
      <Alert.Description>{data.error}</Alert.Description>
    </Alert.Root>
  {/if}

  <div class="flex items-center justify-between gap-4 shrink-0">
    <div class="relative max-w-sm w-full group shrink-0">
      <Search
        class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors"
      />
      <Input
        placeholder="Пошук за назвою..."
        class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
        bind:value={searchQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
      />
    </div>
    <div class="flex items-center justify-between shrink-0">
      <Button
        size="sm"
        class="h-10 shadow-sm"
        onclick={() => (isCreateOpen = true)}
      >
        <Plus class="mr-2 h-4 w-4" />
        Додати пост
      </Button>
    </div>
  </div>

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col mb-4">
    <PostsTable
      posts={data.posts}
      flex={true}
      onEdit={openEdit}
      onDelete={openDelete}
    />
  </div>

  {#if data.pagination && data.pagination.total_pages > 1}
    <div class="shrink-0">
      <SimplePagination
        currentPage={data.pagination.current_page}
        totalPages={data.pagination.total_pages}
        onPageChange={handlePageChange}
      />
    </div>
  {/if}

  <PostCreateDialog bind:open={isCreateOpen} />

  <PostEditSheet bind:open={isEditOpen} post={selectedPost} />

  <DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    title="Видалити пост?"
    itemName={selectedPost?.name}
    action="?/delete"
    id={selectedPost?.ID}
    successMessage="Пост видалено"
  />
</div>
