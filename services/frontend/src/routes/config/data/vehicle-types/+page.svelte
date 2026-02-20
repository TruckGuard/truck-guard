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
  import { Badge } from "$lib/components/ui/badge";
  import type { PageData } from "./$types";
  import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";
  import { darkenColor, getContrastColor } from "$lib/utils/colors";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isEditOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedType = $state<any>(null);

  const presetColors = [
    "#3b82f6", // Blue
    "#ef4444", // Red
    "#10b981", // Emerald
    "#f59e0b", // Amber
    "#8b5cf6", // Violet
    "#ec4899", // Pink
    "#06b6d4", // Cyan
    "#f97316", // Orange
    "#64748b", // Slate
    "#000000", // Black
  ];

  let newColor = $state("#3b82f6");

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

  function openEdit(type: any) {
    selectedType = { ...type };
    isEditOpen = true;
  }

  function openDelete(type: any) {
    selectedType = type;
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
      Додати тип ТЗ
    </Button>
  </div>

  <div class="rounded-md border">
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Код</Table.Head>
          <Table.Head>Назва</Table.Head>
          <Table.Head>Ціна в'їзду</Table.Head>
          <Table.Head>Ціна доби</Table.Head>
          <Table.Head class="text-right">Дії</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#if data.vehicleTypes && data.vehicleTypes.length > 0}
          {#each data.vehicleTypes as type (type.ID)}
            <Table.Row>
              <Table.Cell>
                <Badge
                  variant="outline"
                  style="background-color: {type.color}22; color: {type.color}; border-color: {darkenColor(
                    type.color,
                    20,
                  )}"
                >
                  {type.code}
                </Badge>
              </Table.Cell>
              <Table.Cell class="font-medium">{type.name}</Table.Cell>
              <Table.Cell
                >{type.entry_price?.toLocaleString("uk-UA")} ₴</Table.Cell
              >
              <Table.Cell
                >{type.daily_price?.toLocaleString("uk-UA")} ₴</Table.Cell
              >
              <Table.Cell class="text-right space-x-2">
                <Button
                  variant="ghost"
                  size="icon"
                  onclick={() => openEdit(type)}
                >
                  <Pencil class="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  class="text-destructive hover:text-destructive"
                  onclick={() => openDelete(type)}
                >
                  <Trash2 class="h-4 w-4" />
                </Button>
              </Table.Cell>
            </Table.Row>
          {/each}
        {:else}
          <Table.Row>
            <Table.Cell colspan={6} class="h-24 text-center"
              >Результатів не знайдено.</Table.Cell
            >
          </Table.Row>
        {/if}
      </Table.Body>
    </Table.Root>
  </div>

  <!-- Create Dialog -->
  <Dialog.Root bind:open={isCreateOpen}>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <Dialog.Title>Додати тип ТЗ</Dialog.Title>
        <Dialog.Description
          >Введіть дані для нового типу транспортного засобу.</Dialog.Description
        >
      </Dialog.Header>
      <form
        method="POST"
        action="?/create"
        use:enhance={() => {
          return async ({ result, update }) => {
            if (result.type === "success") {
              isCreateOpen = false;
              toast.success("Тип ТЗ успішно створено");
              await update();
            } else {
              const message = (result as { data?: { message?: string } }).data
                ?.message;
              console.error("Create vehicle type failed:", result);
              toast.error(mapErrorToFriendlyMessage(message));
            }
          };
        }}
        class="space-y-4 px-6 py-4"
      >
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="code">Код</Label>
            <Input id="code" name="code" required placeholder="TRUCK..." />
          </div>
          <div class="space-y-2">
            <Label for="name">Назва</Label>
            <Input id="name" name="name" required placeholder="Вантажівка..." />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="entry_price">Ціна в'їзду (₴)</Label>
            <Input
              id="entry_price"
              name="entry_price"
              type="number"
              step="0.01"
              required
              value="0"
            />
          </div>
          <div class="space-y-2">
            <Label for="daily_price">Ціна доби (₴)</Label>
            <Input
              id="daily_price"
              name="daily_price"
              type="number"
              step="0.01"
              required
              value="0"
            />
          </div>
        </div>
        <div class="space-y-3">
          <Label for="color">Колір</Label>
          <div class="flex flex-wrap gap-2 mb-2">
            {#each presetColors as color}
              <button
                type="button"
                class="w-8 h-8 rounded-full border-2 transition-transform hover:scale-110 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                style="background-color: {color}; border-color: {newColor ===
                color
                  ? 'white'
                  : 'transparent'}; box-shadow: {newColor === color
                  ? '0 0 0 2px ' + color
                  : 'none'}"
                onclick={() => (newColor = color)}
                title={color}
              ></button>
            {/each}
          </div>
          <div class="flex gap-2">
            <div
              class="w-10 h-10 rounded border"
              style="background-color: {newColor}"
            ></div>
            <Input
              id="color"
              name="color"
              placeholder="#3b82f6"
              bind:value={newColor}
              class="flex-1"
            />
          </div>
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
        <Sheet.Title>Редагувати тип ТЗ</Sheet.Title>
        <Sheet.Description
          >Змініть параметри типу транспортного засобу.</Sheet.Description
        >
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
              console.error("Update vehicle type failed:", result);
              toast.error(mapErrorToFriendlyMessage(message));
            }
          };
        }}
        class="space-y-4 px-6 py-6"
      >
        <input type="hidden" name="id" value={selectedType?.ID} />
        <div class="space-y-2">
          <Label for="edit-code">Код</Label>
          <Input
            id="edit-code"
            name="code"
            bind:value={selectedType.code}
            required
          />
        </div>
        <div class="space-y-2">
          <Label for="edit-name">Назва</Label>
          <Input
            id="edit-name"
            name="name"
            bind:value={selectedType.name}
            required
          />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label for="edit-entry">Ціна в'їзду</Label>
            <Input
              id="edit-entry"
              name="entry_price"
              type="number"
              step="0.01"
              bind:value={selectedType.entry_price}
              required
            />
          </div>
          <div class="space-y-2">
            <Label for="edit-daily">Ціна доби</Label>
            <Input
              id="edit-daily"
              name="daily_price"
              type="number"
              step="0.01"
              bind:value={selectedType.daily_price}
              required
            />
          </div>
        </div>
        <div class="space-y-3">
          <Label for="edit-color">Колір</Label>
          <div class="flex flex-wrap gap-2 mb-2">
            {#each presetColors as color}
              <button
                type="button"
                class="w-8 h-8 rounded-full border-2 transition-transform hover:scale-110 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                style="background-color: {color}; border-color: {selectedType.color ===
                color
                  ? 'white'
                  : 'transparent'}; box-shadow: {selectedType.color === color
                  ? '0 0 0 2px ' + color
                  : 'none'}"
                onclick={() => (selectedType.color = color)}
                title={color}
              ></button>
            {/each}
          </div>
          <div class="flex gap-2">
            <div
              class="w-10 h-10 rounded border"
              style="background-color: {selectedType.color}"
            ></div>
            <Input
              id="edit-color"
              name="color"
              bind:value={selectedType.color}
              class="flex-1"
            />
          </div>
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
        <Dialog.Title>Видалити тип ТЗ?</Dialog.Title>
        <Dialog.Description>
          Ви впевнені, що хочете видалити тип ТЗ <strong
            >{selectedType?.name}</strong
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
                toast.success("Тип ТЗ видалено");
                await update();
              } else {
                const message = (result as { data?: { message?: string } }).data
                  ?.message;
                console.error("Delete vehicle type failed:", result);
                toast.error(mapErrorToFriendlyMessage(message));
              }
            };
          }}
        >
          <input type="hidden" name="id" value={selectedType?.ID} />
          <Button type="submit" variant="destructive">Видалити</Button>
        </form>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
</div>
