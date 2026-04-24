<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { can } from "$lib/auth";
  import type { PageData } from "./$types";

  // Component Imports
  import CompaniesTable from "./components/CompaniesTable.svelte";
  import CompanyCreateDialog from "./components/CompanyCreateDialog.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

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
    goto(url.toString(), { keepFocus: true, noScroll: true, replaceState: true });
  }

  function openDelete(company: any) {
    selectedCompany = company;
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
    title="Компанії"
    description="Довідник компаній перевізників та контрагентів."
  >
    {#snippet actions()}
      {#if can(data.user, "create:data")}
        <Button size="sm" class="h-8 px-3 text-xs gap-1.5" onclick={() => (isCreateOpen = true)}>
          <Plus class="h-3.5 w-3.5" /> Додати компанію
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Назва..."
    bind:searchQuery
    paramName="name"
  >
    {#snippet extraFields()}
      <div class="relative w-40 group">
        <Input
          placeholder="ЄДРПОУ..."
          bind:value={edrpouQuery}
          onkeydown={(e) => e.key === "Enter" && handleSearch()}
          class="h-7 bg-transparent border-none focus-visible:ring-1 focus-visible:ring-primary/20 text-xs"
        />
      </div>
    {/snippet}
    <Button
      variant="outline"
      size="sm"
      class="h-7 px-3 text-xs border-none hover:bg-primary/8 hover:text-primary transition-colors"
      onclick={handleSearch}
    >
      Пошук
    </Button>
  </SearchToolbar>

  <CompaniesTable
    companies={data.companies}
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

<CompanyCreateDialog bind:open={isCreateOpen} />

<DeleteConfirmationDialog
  bind:open={isDeleteOpen}
  title="Видалити компанію?"
  itemName={selectedCompany?.name}
  action="?/delete"
  id={selectedCompany?.ID}
  successMessage="Компанію видалено"
/>
