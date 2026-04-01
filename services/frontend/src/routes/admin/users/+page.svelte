<script lang="ts">
  import UsersTable from "./components/UsersTable.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { can } from "$lib/auth";
  import { Plus } from "@lucide/svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import PageHeader from "$lib/components/common/PageHeader.svelte";
  import SearchToolbar from "$lib/components/common/SearchToolbar.svelte";
  import PageLayout from "$lib/components/common/PageLayout.svelte";

  let { data } = $props();

  let isDeleteOpen = $state(false);
  let userToDelete: any = $state(null);
  let isResetOpen = $state(false);
  let userToReset: any = $state(null);
  let newPassword = $state("");
  let searchQuery = $state("");

  let filteredUsers = $derived(
    data.users.filter((u: any) =>
      u.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.role.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  function openDelete(user: any) {
    userToDelete = user;
    isDeleteOpen = true;
  }

  function openReset(user: any) {
    userToReset = user;
    newPassword = "";
    isResetOpen = true;
  }
</script>

<PageLayout>
  <PageHeader
    title="Користувачі"
    description="Керування обліковими записами та ролями персоналу."
  >
    {#snippet actions()}
      {#if can(data.user, "create:users")}
        <Button size="sm" class="h-10 shadow-sm" href="/admin/users/create">
          <Plus class="mr-2 h-4 w-4" />
          Створити користувача
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <SearchToolbar
    placeholder="Пошук за логіном або роллю..."
    bind:searchQuery
    onInput={() => {}}
  />

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <UsersTable
      users={filteredUsers}
      currentUser={data.user}
      posts={data.posts}
      flex={true}
      onDelete={openDelete}
      onResetPassword={openReset}
      onEdit={() => {}}
    />
  </div>

  <Dialog.Root bind:open={isDeleteOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Видалити користувача?</Dialog.Title>
        <Dialog.Description>
          Ви впевнені, що хочете видалити користувача <strong
            >{userToDelete?.username}</strong
          >? Цю дію неможливо скасувати.
        </Dialog.Description>
      </Dialog.Header>
      <Dialog.Footer>
        <Button variant="outline" onclick={() => (isDeleteOpen = false)}>
          Скасувати
        </Button>
        <form
          action="?/delete"
          method="POST"
          use:enhance={() => {
            isDeleteOpen = false;
            toast.loading("Видалення користувача...");
            return async ({ result, update }) => {
              if (result.type === "success") {
                toast.success("Користувача видалено");
                await update();
              } else {
                toast.error("Не вдалося видалити користувача");
              }
            };
          }}
        >
          <input type="hidden" name="id" value={userToDelete?.id} />
          <Button type="submit" variant="destructive">Видалити</Button>
        </form>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>

  <Dialog.Root bind:open={isResetOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Скинути пароль</Dialog.Title>
        <Dialog.Description>
          Введіть новий пароль для користувача <strong>{userToReset?.username}</strong>.
        </Dialog.Description>
      </Dialog.Header>
      <form
        action="?/resetPassword"
        method="POST"
        use:enhance={() => {
          isResetOpen = false;
          toast.loading("Скидання пароля...");
          return async ({ result, update }) => {
            if (result.type === "success") {
              toast.success("Пароль успішно змінено");
              await update();
            } else {
              toast.error("Не вдалося змінити пароль");
            }
          };
        }}
        class="space-y-4 py-4"
      >
        <input type="hidden" name="id" value={userToReset?.id} />
        <div class="space-y-2">
          <Input
            type="password"
            name="newPassword"
            placeholder="Новий пароль"
            bind:value={newPassword}
            required
          />
        </div>
        <Dialog.Footer>
          <Button variant="outline" type="button" onclick={() => (isResetOpen = false)}>
            Скасувати
          </Button>
          <Button type="submit">Зберегти</Button>
        </Dialog.Footer>
      </form>
    </Dialog.Content>
  </Dialog.Root>
</PageLayout>
