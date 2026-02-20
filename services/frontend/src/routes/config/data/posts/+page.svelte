<script lang="ts">
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Sheet from "$lib/components/ui/sheet";
  import * as Alert from "$lib/components/ui/alert";
  import { Pencil, Plus, Trash2, Search, AlertCircle } from "@lucide/svelte";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import type { PageData, ActionData } from "./$types";
  import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedPost = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchQuery) {
      url.searchParams.set("name", searchQuery);
    } else {
      url.searchParams.delete("name");
    }
    url.searchParams.set("page", "1");
    goto(url);
  }

  function openEdit(post: any) {
    selectedPost = { ...post };
    isEditOpen = true;
  }

  function openDelete(post: any) {
    selectedPost = post;
    isDeleteOpen = true;
  }
</script>

<div class="space-y-4">
  {#if data.error}
    <Alert.Root variant="destructive">
      <AlertCircle class="h-4 w-4" />
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
      Додати пост
    </Button>
  </div>

  <div class="rounded-md border">
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Назва</Table.Head>
          <Table.Head>Опис</Table.Head>
          <Table.Head>Дата створення</Table.Head>
          <Table.Head class="text-right">Дії</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#if data.posts && data.posts.length > 0}
          {#each data.posts as post (post.ID)}
            <Table.Row>
              <Table.Cell class="font-medium">{post.name}</Table.Cell>
              <Table.Cell>{post.description || "-"}</Table.Cell>
              <Table.Cell>
                {new Date(post.CreatedAt).toLocaleDateString("uk-UA")}
              </Table.Cell>
              <Table.Cell class="text-right space-x-2">
                <Button
                  variant="ghost"
                  size="icon"
                  onclick={() => openEdit(post)}
                >
                  <Pencil class="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  class="text-destructive hover:text-destructive"
                  onclick={() => openDelete(post)}
                >
                  <Trash2 class="h-4 w-4" />
                </Button>
              </Table.Cell>
            </Table.Row>
          {/each}
        {:else}
          <Table.Row>
            <Table.Cell colspan={4} class="h-24 text-center"
              >Результатів не знайдено.</Table.Cell
            >
          </Table.Row>
        {/if}
      </Table.Body>
    </Table.Root>
  </div>

  <!-- Create Dialog -->
  <Dialog.Root bind:open={isCreateOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Додати митний пост</Dialog.Title>
        <Dialog.Description
          >Введіть назву та опис для нового поста.</Dialog.Description
        >
      </Dialog.Header>
      <form
        method="POST"
        action="?/create"
        use:enhance={() => {
          return async ({ result, update }) => {
            if (result.type === "success") {
              isCreateOpen = false;
              toast.success("Пост успішно створено");
              await update();
            } else {
              const message = (result as { data?: { message?: string } }).data
                ?.message;
              console.error("Create post failed:", result);
              toast.error(mapErrorToFriendlyMessage(message));
            }
          };
        }}
        class="space-y-4 px-6 py-4"
      >
        <div class="space-y-2">
          <Label for="name">Назва</Label>
          <Input id="name" name="name" required placeholder="Назва поста..." />
        </div>
        <div class="space-y-2">
          <Label for="description">Опис</Label>
          <Input
            id="description"
            name="description"
            placeholder="Короткий опис..."
          />
        </div>
        <Dialog.Footer>
          <Button
            variant="outline"
            type="button"
            onclick={() => (isCreateOpen = false)}
          >
            Скасувати
          </Button>
          <Button type="submit">Створити</Button>
        </Dialog.Footer>
      </form>
    </Dialog.Content>
  </Dialog.Root>

  <!-- Edit Sheet -->
  <Sheet.Root bind:open={isEditOpen}>
    <Sheet.Content side="right" class="sm:max-w-md">
      <Sheet.Header>
        <Sheet.Title>Редагувати пост</Sheet.Title>
        <Sheet.Description>Змініть дані митного поста.</Sheet.Description>
      </Sheet.Header>
      <form
        method="POST"
        action="?/update"
        use:enhance={() => {
          return async ({ result, update }) => {
            if (result.type === "success") {
              isEditOpen = false;
              toast.success("Дані оновлено");
              await update();
            } else {
              const message = (result as { data?: { message?: string } }).data
                ?.message;
              console.error("Update post failed:", result);
              toast.error(mapErrorToFriendlyMessage(message));
            }
          };
        }}
        class="space-y-4 px-6 py-6"
      >
        <input type="hidden" name="id" value={selectedPost?.ID} />
        <div class="space-y-2">
          <Label for="edit-name">Назва</Label>
          <Input
            id="edit-name"
            name="name"
            bind:value={selectedPost.name}
            required
          />
        </div>
        <div class="space-y-2">
          <Label for="edit-description">Опис</Label>
          <Input
            id="edit-description"
            name="description"
            bind:value={selectedPost.description}
          />
        </div>
        <Sheet.Footer class="p-0">
          <Button
            variant="outline"
            type="button"
            onclick={() => (isEditOpen = false)}
          >
            Скасувати
          </Button>
          <Button type="submit">Зберегти</Button>
        </Sheet.Footer>
      </form>
    </Sheet.Content>
  </Sheet.Root>

  <!-- Delete Dialog -->
  <Dialog.Root bind:open={isDeleteOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Видалити пост?</Dialog.Title>
        <Dialog.Description>
          Ви впевнені, що хочете видалити пост <strong
            >{selectedPost?.name}</strong
          >? Цю дію неможливо скасувати.
        </Dialog.Description>
      </Dialog.Header>
      <Dialog.Footer>
        <Button
          variant="outline"
          type="button"
          onclick={() => (isDeleteOpen = false)}
        >
          Скасувати
        </Button>
        <form
          method="POST"
          action="?/delete"
          use:enhance={() => {
            return async ({ result, update }) => {
              if (result.type === "success") {
                isDeleteOpen = false;
                toast.success("Пост видалено");
                await update();
              } else {
                const message = (result as { data?: { message?: string } }).data
                  ?.message;
                console.error("Delete post failed:", result);
                toast.error(mapErrorToFriendlyMessage(message));
              }
            };
          }}
        >
          <input type="hidden" name="id" value={selectedPost?.ID} />
          <Button type="submit" variant="destructive">Видалити</Button>
        </form>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
</div>
