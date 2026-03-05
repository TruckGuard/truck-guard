<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, Search, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  // Component Imports
  import PaymentTypesTable from "./components/PaymentTypesTable.svelte";
  import PaymentTypeCreateDialog from "./components/PaymentTypeCreateDialog.svelte";
  import PaymentTypeEditSheet from "./components/PaymentTypeEditSheet.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedType = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchQuery) url.searchParams.set("name", searchQuery);
    else url.searchParams.delete("name");

    url.searchParams.set("page", "1");
    goto(url);
  }

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
      Додати тип оплати
    </Button>
  </div>

  <PaymentTypesTable
    paymentTypes={data.paymentTypes}
    onEdit={openEdit}
    onDelete={openDelete}
  />

  {#if data.pagination && data.pagination.total_pages > 1}
    <SimplePagination
      currentPage={data.pagination.current_page}
      totalPages={data.pagination.total_pages}
      onPageChange={handlePageChange}
    />
  {/if}

  <PaymentTypeCreateDialog bind:open={isCreateOpen} />

  <PaymentTypeEditSheet bind:open={isEditOpen} type={selectedType} />

  <DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    title="Видалити спосіб оплати?"
    itemName={selectedType?.name}
    action="?/delete"
    id={selectedType?.ID}
    successMessage="Спосіб оплати видалено"
  />
</div>
