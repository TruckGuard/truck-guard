<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, Search, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  // Component Imports
  import ModesTable from "./components/ModesTable.svelte";
  import ModeCreateDialog from "./components/ModeCreateDialog.svelte";
  import ModeEditSheet from "./components/ModeEditSheet.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedMode = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchQuery) url.searchParams.set("name", searchQuery);
    else url.searchParams.delete("name");

    url.searchParams.set("page", "1");
    goto(url);
  }

  function openEdit(mode: any) {
    selectedMode = { ...mode };
    isEditOpen = true;
  }

  function openDelete(mode: any) {
    selectedMode = mode;
    isDeleteOpen = true;
  }

  function handlePageChange(newPage: number) {
    const url = new URL(page.url);
    url.searchParams.set("page", newPage.toString());
    goto(url);
  }
</script>

<div class="flex flex-col h-full overflow-hidden space-y-6">
  <div class="flex items-center justify-between gap-4 shrink-0">
    <div class="relative max-w-sm w-full group">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
      <Input
        placeholder="Пошук за назвою..."
        class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
        bind:value={searchQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
      />
    </div>
    <Button size="sm" class="h-10 shadow-sm" onclick={() => (isCreateOpen = true)}>
      <Plus class="mr-2 h-4 w-4" />
      Додати режим
    </Button>
  </div>

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <ModesTable modes={data.modes} flex={true} onEdit={openEdit} onDelete={openDelete} />
  </div>

  {#if data.pagination && data.pagination.total_pages > 1}
    <SimplePagination
      currentPage={data.pagination.current_page}
      totalPages={data.pagination.total_pages}
      onPageChange={handlePageChange}
    />
  {/if}

  <ModeCreateDialog bind:open={isCreateOpen} />

  <ModeEditSheet bind:open={isEditOpen} mode={selectedMode} />

  <DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    title="Видалити режим?"
    itemName={selectedMode?.name}
    description={`Ви впевнені, що хочете видалити режим ${selectedMode?.name} (${selectedMode?.code})? Цю дію неможливо скасувати.`}
    action="?/delete"
    id={selectedMode?.ID}
    successMessage="Режим видалено"
  />
</div>
