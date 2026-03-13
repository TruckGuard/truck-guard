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
      Додати тип оплати
    </Button>
  </div>

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <PaymentTypesTable
      paymentTypes={data.paymentTypes}
      flex={true}
      onEdit={openEdit}
      onDelete={openDelete}
    />
  </div>

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
