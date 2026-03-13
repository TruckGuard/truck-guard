<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, Search, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  // Component Imports
  import CompaniesTable from "./components/CompaniesTable.svelte";
  import CompanyCreateDialog from "./components/CompanyCreateDialog.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedCompany = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");
  let edrpouQuery = $state(page.url.searchParams.get("edrpou") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchQuery) url.searchParams.set("name", searchQuery);
    else url.searchParams.delete("name");

    if (edrpouQuery) url.searchParams.set("edrpou", edrpouQuery);
    else url.searchParams.delete("edrpou");

    url.searchParams.set("page", "1");
    goto(url);
  }

  function openDelete(company: any) {
    selectedCompany = company;
    isDeleteOpen = true;
  }

  function handlePageChange(newPage: number) {
    const url = new URL(page.url);
    url.searchParams.set("page", newPage.toString());
    goto(url);
  }
</script>

<div class="flex flex-col overflow-hidden space-y-6">
  {#if data.error}
    <Alert.Root variant="destructive" class="shrink-0">
      <CircleAlert class="h-4 w-4" />
      <Alert.Title>Помилка</Alert.Title>
      <Alert.Description>{data.error}</Alert.Description>
    </Alert.Root>
  {/if}

  <div class="flex items-center justify-between shrink-0">
    <div class="flex items-center gap-3 max-w-xl w-full shrink-0">
      <div class="relative flex-1 group">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
        <Input
          placeholder="Назва..."
          class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
          bind:value={searchQuery}
          onkeydown={(e) => e.key === "Enter" && handleSearch()}
        />
      </div>
      <div class="relative w-40 group">
        <Input
          placeholder="ЄДРПОУ..."
          bind:value={edrpouQuery}
          onkeydown={(e) => e.key === "Enter" && handleSearch()}
          class="h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
        />
      </div>
      <Button variant="outline" size="sm" class="h-10 px-4" onclick={handleSearch}>
        Пошук
      </Button>
    </div>
    <Button size="sm" class="h-10 shadow-sm" onclick={() => (isCreateOpen = true)}>
      <Plus class="mr-2 h-4 w-4" />
      Додати компанію
    </Button>
  </div>


  <div class="flex-1 min-h-0 overflow-hidden flex flex-col mb-4">
    <CompaniesTable companies={data.companies} flex={true} onDelete={openDelete} />
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

  <CompanyCreateDialog bind:open={isCreateOpen} />

  <DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    title="Видалити компанію?"
    itemName={selectedCompany?.name}
    action="?/delete"
    id={selectedCompany?.ID}
    successMessage="Компанію видалено"
  />
</div>
