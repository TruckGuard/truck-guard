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

<div class="space-y-4 p-6">
  {#if data.error}
    <Alert.Root variant="destructive">
      <CircleAlert class="h-4 w-4" />
      <Alert.Title>Помилка</Alert.Title>
      <Alert.Description>{data.error}</Alert.Description>
    </Alert.Root>
  {/if}

  <div class="flex items-center justify-between">
    <div class="flex items-center gap-2 max-w-sm w-full">
      <Input
        placeholder="Пошук за назвою..."
        bind:value={searchQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
      />
      <Button variant="outline" size="icon" onclick={handleSearch}>
        <Search class="h-4 w-4" />
      </Button>
    </div>
    <Button onclick={() => (isCreateOpen = true)}>
      <Plus class="mr-2 h-4 w-4" />
      Додати пост
    </Button>
  </div>

  <PostsTable posts={data.posts} onEdit={openEdit} onDelete={openDelete} />

  {#if data.pagination && data.pagination.total_pages > 1}
    <SimplePagination
      currentPage={data.pagination.current_page}
      totalPages={data.pagination.total_pages}
      onPageChange={handlePageChange}
    />
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
