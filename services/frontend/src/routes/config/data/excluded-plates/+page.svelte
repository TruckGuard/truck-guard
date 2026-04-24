<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { can } from "$lib/auth";
  import type { PageData } from "./$types";

  // Component Imports
  import ExcludedPlatesTable from "./components/ExcludedPlatesTable.svelte";
  import ExcludedPlateCreateDialog from "./components/ExcludedPlateCreateDialog.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedPlate = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("plate") || "");

  function openDelete(plate: any) {
    selectedPlate = plate;
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
  {#if data.error}
    <Alert.Root variant="destructive" class="shrink-0">
      <CircleAlert class="h-4 w-4" />
      <Alert.Title>Помилка</Alert.Title>
      <Alert.Description>{data.error}</Alert.Description>
    </Alert.Root>
  {/if}

  <PageHeader
    title="Ігноровані номери"
    description="Список номерних знаків, які система ігнорує при розпізнаванні (наприклад, службовий транспорт)."
  >
    {#snippet actions()}
      {#if can(data.user, "create:settings")}
        <Button size="sm" class="h-8 px-3 text-xs gap-1.5" onclick={() => (isCreateOpen = true)}>
          <Plus class="h-3.5 w-3.5" /> Додати номер
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Пошук номера..."
    bind:searchQuery
    paramName="plate"
  />

  <ExcludedPlatesTable
    plates={data.plates}
    currentUser={data.user}
    onDelete={openDelete}
  />

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

<ExcludedPlateCreateDialog bind:open={isCreateOpen} />

<DeleteConfirmationDialog
  bind:open={isDeleteOpen}
  title="Видалити номер?"
  itemName={selectedPlate?.plate}
  action="?/delete"
  id={selectedPlate?.ID}
  successMessage="Номер видалено зі списку ігнорування"
/>
