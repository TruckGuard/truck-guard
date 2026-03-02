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
      Додати режим
    </Button>
  </div>

  <ModesTable modes={data.modes} onEdit={openEdit} onDelete={openDelete} />

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
