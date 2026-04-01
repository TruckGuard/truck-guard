<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Plus } from "@lucide/svelte";
  import type { PageData } from "./$types";
  import type { APIKey } from "$lib/server/auth-client";
  import { can } from "$lib/auth";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

  // Component Imports
  import KeysTable from "./components/KeysTable.svelte";
  import KeysCreateDialog from "./components/KeysCreateDialog.svelte";
  import KeysEditDialog from "./components/KeysEditDialog.svelte";
  import KeysPermissionsDialog from "./components/KeysPermissionsDialog.svelte";
  import KeysDeleteDialog from "./components/KeysDeleteDialog.svelte";
  import KeysSecretDialog from "./components/KeysSecretDialog.svelte";

  let { data }: { data: PageData } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isPermsOpen = $state(false);
  let isDeleteOpen = $state(false);
  let isSecretOpen = $state(false);

  let currentKey = $state<APIKey | null>(null);
  let generatedSecret = $state("");
  let searchQuery = $state("");

  let filteredKeys = $derived(
    data.keys.filter((k: APIKey) => 
      k.owner_name.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  function openCreate() {
    currentKey = null;
    isCreateOpen = true;
  }

  function openEdit(key: APIKey) {
    currentKey = { ...key };
    isEditOpen = true;
  }

  function openDelete(key: APIKey) {
    currentKey = key;
    isDeleteOpen = true;
  }

  function openPerms(key: APIKey) {
    currentKey = key;
    isPermsOpen = true;
  }

  function handleCreateSuccess(secret: string) {
    generatedSecret = secret;
    isSecretOpen = true;
  }

  function canAccessPermission(permission: string): boolean {
    return can(data.user, permission);
  }
</script>

<PageLayout>
  <PageHeader
    title="Ключі API"
    description="Системні ідентифікатори для інтеграції з обладнанням."
  >
    {#snippet actions()}
      {#if canAccessPermission("create:keys")}
        <Button size="sm" class="h-10 shadow-sm" onclick={openCreate}>
          <Plus class="mr-2 h-4 w-4" />
          Створити ключ
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar 
    placeholder="Пошук власників ключів..." 
    bind:searchQuery 
    onInput={() => {}} 
  />

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <KeysTable
      keys={filteredKeys}
      flex={true}
      onEdit={openEdit}
      onDelete={openDelete}
      onPermissions={openPerms}
    />
  </div>
</PageLayout>

<KeysCreateDialog
  bind:open={isCreateOpen}
  permissions={data.permissions}
  canAccess={canAccessPermission}
  onSuccess={handleCreateSuccess}
/>

<KeysEditDialog
  bind:open={isEditOpen}
  key={currentKey}
/>

<KeysPermissionsDialog
  bind:open={isPermsOpen}
  key={currentKey}
  permissions={data.permissions}
  canAccess={canAccessPermission}
/>

<KeysDeleteDialog
  bind:open={isDeleteOpen}
  key={currentKey}
/>

<KeysSecretDialog
  bind:open={isSecretOpen}
  secret={generatedSecret}
/>
