<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Plus, Search } from "@lucide/svelte";
  import type { PageData } from "./$types";
  import type { APIKey } from "$lib/server/auth-client";
  import { can } from "$lib/auth";

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

<div class="flex flex-col h-full overflow-hidden space-y-6">
  <div class="shrink-0">
    <h1 class="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-2">Ключі API</h1>
    <p class="text-muted-foreground text-sm mb-0">
      Системні ідентифікатори для інтеграції з обладнанням
    </p>
  </div>

  <div class="flex items-center justify-between gap-4 shrink-0">
    <div class="relative max-w-sm w-full group">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
      <Input
        placeholder="Пошук власників ключів..."
        class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
        bind:value={searchQuery}
      />
    </div>
    {#if canAccessPermission("create:keys")}
      <Button size="sm" class="h-10 shadow-sm" onclick={openCreate}>
        <Plus class="mr-2 h-4 w-4" />
        Створити ключ
      </Button>
    {/if}
  </div>

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <KeysTable
      keys={filteredKeys}
      flex={true}
      onEdit={openEdit}
      onDelete={openDelete}
      onPermissions={openPerms}
    />
  </div>

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
</div>
