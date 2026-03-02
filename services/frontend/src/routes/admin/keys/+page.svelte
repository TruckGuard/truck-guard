<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import Plus from "@lucide/svelte/icons/plus";
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

  function canAccess(permission: string): boolean {
    return can(data.user, permission);
  }
</script>

<div class="p-6 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">API Ключі</h1>
      <p class="text-muted-foreground">
        Керування ключами доступу для зовнішніх інтеграцій.
      </p>
    </div>
    {#if canAccess("create:keys")}
      <Button onclick={openCreate}>
        <Plus class="mr-2 h-4 w-4" />
        Створити ключ
      </Button>
    {/if}
  </div>

  <KeysTable
    keys={data.keys}
    canUpdate={canAccess("update:keys")}
    canDelete={canAccess("delete:keys")}
    onOpenEdit={openEdit}
    onOpenDelete={openDelete}
    onOpenPerms={openPerms}
  />

  <KeysCreateDialog
    bind:open={isCreateOpen}
    permissions={data.permissions}
    {canAccess}
    onSuccess={handleCreateSuccess}
  />

  <KeysSecretDialog bind:open={isSecretOpen} secret={generatedSecret} />

  <KeysEditDialog bind:open={isEditOpen} key={currentKey} />

  <KeysPermissionsDialog
    bind:open={isPermsOpen}
    key={currentKey}
    permissions={data.permissions}
    {canAccess}
  />

  <KeysDeleteDialog bind:open={isDeleteOpen} key={currentKey} />
</div>
