<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { can } from "$lib/auth";
  import type { PageData } from "./$types";

  // Component Imports
  import VehicleTypesTable from "./components/VehicleTypesTable.svelte";
  import VehicleTypeCreateDialog from "./components/VehicleTypeCreateDialog.svelte";
  import VehicleTypeEditSheet from "./components/VehicleTypeEditSheet.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedType = $state<any>(null);

  const presetColors = [
    "#3b82f6",
    "#ef4444",
    "#10b981",
    "#f59e0b",
    "#8b5cf6",
    "#ec4899",
    "#06b6d4",
    "#f97316",
    "#64748b",
    "#000000",
  ];

  let searchQuery = $state(page.url.searchParams.get("name") || "");

  function openEdit(type: any) {
    selectedType = { ...type };
    isEditOpen = true;
  }

  function openDelete(type: any) {
    selectedType = type;
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
    title="Типи ТЗ"
    description="Класифікація транспортних засобів для автоматичного розпізнавання."
  >
    {#snippet actions()}
      {#if can(data.user, "create:data")}
        <Button size="sm" class="h-8 px-3 text-xs gap-1.5" onclick={() => (isCreateOpen = true)}>
          <Plus class="h-3.5 w-3.5" /> Додати тип ТЗ
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Пошук за назвою..."
    bind:searchQuery
    paramName="name"
  />

  <VehicleTypesTable
    vehicleTypes={data.vehicleTypes}
    currentUser={data.user}
    onEdit={openEdit}
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

<VehicleTypeCreateDialog bind:open={isCreateOpen} {presetColors} />

<VehicleTypeEditSheet
  bind:open={isEditOpen}
  type={selectedType}
  {presetColors}
/>

<DeleteConfirmationDialog
  bind:open={isDeleteOpen}
  title="Видалити тип ТЗ?"
  itemName={selectedType?.name}
  action="?/delete"
  id={selectedType?.ID}
  successMessage="Тип ТЗ видалено"
/>
