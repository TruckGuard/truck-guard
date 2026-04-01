<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Plus } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  // Component Imports
  import ModesTable from "./components/ModesTable.svelte";
  import ModeCreateDialog from "./components/ModeCreateDialog.svelte";
  import ModeEditSheet from "./components/ModeEditSheet.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedMode = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");

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
    goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
  }

  function handleLimitChange(newLimit: number) {
    const url = new URL(page.url);
    url.searchParams.set("limit", newLimit.toString());
    url.searchParams.set("page", "1");
    goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
  }
</script>

<PageLayout>
  <PageHeader
    title="Режими роботи"
    description="Налаштування логіки обробки дозволів та подій."
  >
    {#snippet actions()}
      <Button
        size="sm"
        class="h-10 shadow-sm"
        onclick={() => (isCreateOpen = true)}
      >
        <Plus class="mr-2 h-4 w-4" />
        Додати режим
      </Button>
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Пошук за назвою..."
    bind:searchQuery
    paramName="name"
  />

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <ModesTable
      modes={data.modes}
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
        itemsPerPage={data.pagination.limit}
        onPageChange={handlePageChange}
        onLimitChange={handleLimitChange}
      />
    </div>
  {/if}
</PageLayout>

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
