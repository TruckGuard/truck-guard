<script lang="ts">
  import RolesTable from "./components/RolesTable.svelte";
  import { Button } from "$lib/components/ui/button";
  import { can as authCan } from "$lib/auth";
  import { Plus, Users, GitBranch } from "@lucide/svelte";
  import Search from "@lucide/svelte/icons/search";
  import * as Dialog from "$lib/components/ui/dialog";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { Badge } from "$lib/components/ui/badge";
  import * as Tabs from "$lib/components/ui/tabs";
  import PermissionHierarchy from "./components/PermissionHierarchy.svelte";
  import type { PageData } from "./$types";
  import type { Role, Permission } from "$lib/server/auth-client";
  import PageLayout from "$lib/components/common/PageLayout.svelte";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";

  let { data }: { data: PageData } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isPermsOpen = $state(false);
  let isDeleteOpen = $state(false);

  let currentRole: Role | null = $state(null);
  let selectedPermissions: string[] = $state([]);
  let permSearch = $state("");
  let roleSearch = $state("");

  let filteredRoles = $derived(
    data.roles.filter((r: Role) => 
      r.name.toLowerCase().includes(roleSearch.toLowerCase()) ||
      r.description.toLowerCase().includes(roleSearch.toLowerCase())
    )
  );

  let filteredPermissions = $derived(
    data.permissions.filter(
      (p: Permission) =>
        p.id.toLowerCase().includes(permSearch.toLowerCase()) ||
        p.description.toLowerCase().includes(permSearch.toLowerCase()),
    ),
  );

  function openCreate() {
    currentRole = null;
    isCreateOpen = true;
  }

  function openEdit(role: Role) {
    currentRole = role;
    isEditOpen = true;
  }

  function openDelete(role: Role) {
    currentRole = role;
    isDeleteOpen = true;
  }

  function openPerms(role: Role) {
    currentRole = role;
    if (role.permissions) {
      selectedPermissions = role.permissions.map((p) => p.id);
    } else {
      selectedPermissions = [];
    }
    permSearch = "";
    isPermsOpen = true;
  }

  function togglePermission(id: string) {
    if (selectedPermissions.includes(id)) {
      selectedPermissions = selectedPermissions.filter((p) => p !== id);
    } else {
      selectedPermissions = [...selectedPermissions, id];
    }
  }

  let inheritedPermissions = $derived.by(() => {
    const inherited = new Set<string>();
    const check = (id: string, seen: Set<string>) => {
      if (seen.has(id)) return;
      seen.add(id);
      
      // Автоматично додаємо базове право для прав з суфіксом :all
      if (id.endsWith(":all")) {
        const base = id.slice(0, -4);
        if (!inherited.has(base)) {
          inherited.add(base);
          check(base, seen);
        }

        // Рекурсивно додаємо :all для залежних прав
        const children = data.hierarchy[base] || [];
        for (const child of children) {
          const childAll = child + ":all";
          if (data.permissions.some((p: { id: string }) => p.id === childAll)) {
            if (!inherited.has(childAll)) {
              inherited.add(childAll);
              check(childAll, seen);
            }
          }
        }
      }

      const children = data.hierarchy[id] || [];
      for (const child of children) {
        inherited.add(child);
        check(child, seen);
      }
    };
    for (const p of selectedPermissions) {
      check(p, new Set());
    }
    return Array.from(inherited);
  });
</script>
<PageLayout>
  <PageHeader
    title="Ролі та права"
    description="Керування ролями користувачів та їх доступом до функцій системи."
  >
    {#snippet actions()}
      {#if authCan(data.user, "create:roles")}
        <Button size="sm" class="h-10 shadow-sm" onclick={openCreate}>
          <Plus class="mr-2 h-4 w-4" />
          Створити роль
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Пошук за назвою або описом..."
    bind:searchQuery={roleSearch}
    debounce={false}
  />

  <Tabs.Root value="roles" class="flex-1 flex flex-col min-h-0 overflow-hidden">
    <div class="flex items-center justify-between mb-4 shrink-0">
      <Tabs.List class="w-full justify-start grid-cols-2 lg:w-[400px] grid">
        <Tabs.Trigger value="roles">
          <Users class="mr-2 h-4 w-4" /> Список ролей
        </Tabs.Trigger>
        <Tabs.Trigger value="hierarchy">
          <GitBranch class="mr-2 h-4 w-4" /> Ієрархія прав
        </Tabs.Trigger>
      </Tabs.List>
    </div>

    <Tabs.Content value="roles" class="flex-1 min-h-0 overflow-hidden mt-0">
      <div class="h-full flex flex-col">
        <RolesTable 
          roles={filteredRoles} 
          currentUser={data.user} 
          flex={true}
          onPerms={openPerms} 
          onEdit={openEdit} 
          onDelete={openDelete} 
        />
      </div>
    </Tabs.Content>

    <Tabs.Content value="hierarchy" class="flex-1 min-h-0 overflow-y-auto mt-0">
      <PermissionHierarchy 
        hierarchy={data.hierarchy} 
        permissions={data.permissions} 
      />
    </Tabs.Content>
  </Tabs.Root>

  <!-- Create Role Dialog -->
  <Dialog.Root bind:open={isCreateOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Створити нову роль</Dialog.Title>
      </Dialog.Header>
      <form
        action="?/create"
        method="POST"
        use:enhance={({ formData }) => {
          const name = formData.get("name") as string;
          const description = formData.get("description") as string;

          // Optimistic UI
          const newRole: Role = {
            id: Date.now(), // Temporary id
            name,
            description,
            permissions: [],
          };

          const originalRoles = data.roles;
          data.roles = [...data.roles, newRole];
          isCreateOpen = false;
          toast.loading("Створення ролі...");

          return async ({ result, update }) => {
            if (result.type === "error" || result.type === "failure") {
              data.roles = originalRoles;
              toast.error("Не вдалося створити роль");
              isCreateOpen = true;
            } else {
              toast.success("Роль створено");
              await update();
            }
          };
        }}
      >
        <div class="space-y-4 py-4">
          <div class="grid gap-2">
            <Label for="name">Назва</Label>
            <Input id="name" name="name" required placeholder="admin" />
          </div>
          <div class="grid gap-2">
            <Label for="desc">Опис</Label>
            <Input
              id="desc"
              name="description"
              placeholder="Адміністратор системи"
            />
          </div>
        </div>
        <Dialog.Footer>
          <Button type="submit">Створити</Button>
        </Dialog.Footer>
      </form>
    </Dialog.Content>
  </Dialog.Root>

  <!-- Edit Role Dialog -->
  <Dialog.Root bind:open={isEditOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Редагувати роль</Dialog.Title>
      </Dialog.Header>
      <form
        action="?/update"
        method="POST"
        use:enhance={() => {
          toast.loading("Оновлення ролі...");
          return async ({ result, update }) => {
            if (result.type === "success") {
              isEditOpen = false;
              toast.success("Роль оновлено");
              await update();
            } else {
              toast.error("Не вдалося оновити роль");
            }
          };
        }}
      >
        <input type="hidden" name="id" value={currentRole?.id} />
        <div class="space-y-4 py-4">
          <div class="grid gap-2">
            <Label for="edit-name">Назва</Label>
            <Input
              id="edit-name"
              name="name"
              bind:value={currentRole!.name}
              required
            />
          </div>
          <div class="grid gap-2">
            <Label for="edit-desc">Опис</Label>
            <Input
              id="edit-desc"
              name="description"
              bind:value={currentRole!.description}
            />
          </div>
        </div>
        <Dialog.Footer>
          <Button type="submit">Зберегти</Button>
        </Dialog.Footer>
      </form>
    </Dialog.Content>
  </Dialog.Root>

  <!-- Permissions Dialog -->
  <Dialog.Root bind:open={isPermsOpen}>
    <Dialog.Content class="max-w-2xl">
      <Dialog.Header>
        <Dialog.Title
          >Налаштування прав доступу: {currentRole?.name}</Dialog.Title
        >
      </Dialog.Header>
      <form
        action="?/assignPermissions"
        method="POST"
        use:enhance={() => {
          toast.loading("Збереження прав...");
          return async ({ result, update }) => {
            if (result.type === "success") {
              isPermsOpen = false;
              toast.success("Права оновлено");
              await update();
            } else {
              toast.error("Не вдалося оновити права");
            }
          };
        }}
      >
        <input type="hidden" name="id" value={currentRole?.id} />
        {#each selectedPermissions as pId}
          <input type="hidden" name="permissions" value={pId} />
        {/each}
        <div class="px-1 py-4">
          <Input
            placeholder="Пошук прав..."
            bind:value={permSearch}
            class="mb-4"
          />
        </div>
        <div class="py-4 h-[50vh] overflow-y-auto">
          <div class="grid grid-cols-2 gap-4">
            {#each filteredPermissions as perm}
              <div class="flex items-start space-x-2">
                <Checkbox
                  id="perm-{perm.id}"
                  value={perm.id}
                  checked={selectedPermissions.includes(perm.id) || inheritedPermissions.includes(perm.id)}
                  onCheckedChange={() => togglePermission(perm.id)}
                  disabled={!authCan(data.user, perm.id) || inheritedPermissions.includes(perm.id)}
                />
                <div class="grid gap-1.5 leading-none">
                  <Label
                    for="perm-{perm.id}"
                    class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  >
                    <div class="flex flex-col">
                      <span class={inheritedPermissions.includes(perm.id) ? "text-primary/90 font-semibold" : ""}>
                        {perm.name}
                      </span>
                      <span class="text-xs text-muted-foreground flex items-center gap-1">
                        {perm.id}
                        {#if inheritedPermissions.includes(perm.id)}
                          <Badge variant="outline" class="h-4 px-1 text-[9px] uppercase tracking-tighter bg-primary/5 text-primary border-primary/20">успадковано</Badge>
                        {/if}
                      </span>
                    </div>
                  </Label>
                </div>
              </div>
            {/each}
          </div>
        </div>
        <Dialog.Footer>
          <Button type="submit">Зберегти права</Button>
        </Dialog.Footer>
      </form>
    </Dialog.Content>
  </Dialog.Root>

  <!-- Delete Confirmation Dialog -->
  <Dialog.Root bind:open={isDeleteOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Видалити роль?</Dialog.Title>
        <Dialog.Description>
          Ви впевнені, що хочете видалити роль <strong
            >{currentRole?.name}</strong
          >? Цю дію неможливо скасувати.
        </Dialog.Description>
      </Dialog.Header>
      <Dialog.Footer>
        <Button variant="outline" onclick={() => (isDeleteOpen = false)}
          >Скасувати</Button
        >
        <form
          action="?/delete"
          method="POST"
          use:enhance={() => {
            const roleId = currentRole?.id;
            const originalRoles = data.roles;
            if (roleId) {
              data.roles = data.roles.filter((r: Role) => r.id !== roleId);
            }
            isDeleteOpen = false;
            toast.info("Видалення ролі...");
            return async ({ result, update }) => {
              if (result.type === "error" || result.type === "failure") {
                data.roles = originalRoles;
                toast.error("Не вдалося видалити роль");
              } else {
                toast.success("Роль видалено");
                await update();
              }
            };
          }}
        >
          <input type="hidden" name="id" value={currentRole?.id} />
          <Button type="submit" variant="destructive">Видалити</Button>
        </form>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
</PageLayout>
