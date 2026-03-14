<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import * as Alert from "$lib/components/ui/alert";
  import { Plus, Search, CircleAlert } from "@lucide/svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  // Component Imports
  import ExcludedPlatesTable from "./components/ExcludedPlatesTable.svelte";
  import ExcludedPlateCreateDialog from "./components/ExcludedPlateCreateDialog.svelte";
  import DeleteConfirmationDialog from "$lib/components/common/DeleteConfirmationDialog.svelte";
  import SimplePagination from "$lib/components/common/SimplePagination.svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedPlate = $state<any>(null);

  let searchPlate = $state(page.url.searchParams.get("plate") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchPlate) url.searchParams.set("plate", searchPlate);
    else url.searchParams.delete("plate");

    url.searchParams.set("page", "1");
    goto(url);
  }

  function openDelete(plate: any) {
    selectedPlate = plate;
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
    <div class="flex items-center gap-3 max-w-sm w-full shrink-0">
      <div class="relative flex-1 group">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
        <Input
          placeholder="Номер..."
          class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
          bind:value={searchPlate}
          onkeydown={(e) => e.key === "Enter" && handleSearch()}
        />
      </div>
      <Button variant="outline" size="sm" class="h-10 px-4" onclick={handleSearch}>
        Пошук
      </Button>
    </div>
    <Button size="sm" class="h-10 shadow-sm" onclick={() => (isCreateOpen = true)}>
      <Plus class="mr-2 h-4 w-4" />
      Додати номер
    </Button>
  </div>


  <div class="flex-1 min-h-0 overflow-hidden flex flex-col mb-4">
    <ExcludedPlatesTable plates={data.plates} flex={true} onDelete={openDelete} />
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

  <ExcludedPlateCreateDialog bind:open={isCreateOpen} />

  <DeleteConfirmationDialog
    bind:open={isDeleteOpen}
    title="Видалити номер?"
    itemName={selectedPlate?.plate}
    action="?/delete"
    id={selectedPlate?.ID}
    successMessage="Номер видалено зі списку ігнорування"
  />
</div>
