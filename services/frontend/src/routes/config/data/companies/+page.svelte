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
    <div class="flex items-center gap-2 max-w-xl w-full">
      <Input
        placeholder="Назва..."
        bind:value={searchQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
      />
      <Input
        placeholder="ЄДРПОУ..."
        bind:value={edrpouQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
        class="w-40"
      />
      <Button variant="outline" size="icon" onclick={handleSearch}>
        <Search class="h-4 w-4" />
      </Button>
    </div>
    <Button onclick={() => (isCreateOpen = true)}>
      <Plus class="mr-2 h-4 w-4" />
      Додати компанію
    </Button>
  </div>

  <CompaniesTable companies={data.companies} onDelete={openDelete} />

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
