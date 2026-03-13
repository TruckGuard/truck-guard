<script lang="ts">
  import UsersTable from "./components/UsersTable.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { can } from "$lib/auth";
  import { Plus, Search } from "@lucide/svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";

  let { data } = $props();

  let isDeleteOpen = $state(false);
  let userToDelete: any = $state(null);
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
</script>

<div class="flex flex-col h-full overflow-hidden space-y-6">
  <div class="flex items-center justify-between gap-4 shrink-0">
    <div class="relative max-w-sm w-full group">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
      <Input
        placeholder="Пошук за логіном або роллю..."
        class="pl-9 h-10 bg-muted/50 border-none focus-visible:ring-1 focus-visible:ring-primary/20"
        bind:value={searchQuery}
      />
    </div>
    {#if can(data.user, "create:users")}
      <Button size="sm" class="h-10 shadow-sm" href="/admin/users/create">
        <Plus class="mr-2 h-4 w-4" />
        Створити користувача
      </Button>
    {/if}
  </div>

  <div class="flex-1 min-h-0 overflow-hidden flex flex-col">
    <UsersTable 
      users={filteredUsers} 
      currentUser={data.user} 
      posts={data.posts} 
      flex={true}
      onDelete={openDelete} 
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
</div>
