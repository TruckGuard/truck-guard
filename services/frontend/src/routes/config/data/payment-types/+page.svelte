<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Plus } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { can } from "$lib/auth";
  import type { PageData } from "./$types";

  // Component Imports
  import PaymentTypesTable from "./components/PaymentTypesTable.svelte";
  import PaymentTypeCreateDialog from "./components/PaymentTypeCreateDialog.svelte";
  import PaymentTypeEditSheet from "./components/PaymentTypeEditSheet.svelte";
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
  <PageHeader
    title="Типи оплати"
    description="Керування доступними способами оплати послуг."
  >
    {#snippet actions()}
      {#if can(data.user, "create:data")}
        <Button size="sm" class="h-8 px-3 text-xs gap-1.5" onclick={() => (isCreateOpen = true)}>
          <Plus class="h-3.5 w-3.5" /> Додати тип оплати
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Пошук за назвою..."
    bind:searchQuery
    paramName="name"
  />

  <PaymentTypesTable
    paymentTypes={data.paymentTypes}
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
