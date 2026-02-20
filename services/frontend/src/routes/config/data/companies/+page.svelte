<script lang="ts">
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Alert from "$lib/components/ui/alert";
  import {
    Pencil,
    Plus,
    Trash2,
    Search,
    ExternalLink,
    AlertCircle,
  } from "@lucide/svelte";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { Badge } from "$lib/components/ui/badge";
  import type { PageData } from "./$types";
  import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isCreateOpen = $state(false);
  let isDeleteOpen = $state(false);
  let selectedCompany = $state<any>(null);

  let searchQuery = $state(page.url.searchParams.get("name") || "");
  let edrpouQuery = $state(page.url.searchParams.get("edrpou") || "");

  function handleSearch() {
    const url = new URL(page.url);
    if (searchQuery) {
      url.searchParams.set("name", searchQuery);
    } else {
      url.searchParams.delete("name");
    }
    if (edrpouQuery) {
      url.searchParams.set("edrpou", edrpouQuery);
    } else {
      url.searchParams.delete("edrpou");
    }
    url.searchParams.set("page", "1");
    goto(url);
  }

  function openDelete(company: any) {
    selectedCompany = company;
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
    <div class="flex items-center gap-2 max-w-xl w-full">
      <Input
        placeholder="Назва..."
        bind:value={searchQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
      />
      <Input
        placeholder="ЄДРПОУ..."
        bind:value={edrpouQuery}
        onkeydown={(e) => e.key === "Enter" && handleSearch()}
        class="w-40"
      />
      <Button variant="outline" size="icon" onclick={handleSearch}>
        <Search class="h-4 w-4" />
      </Button>
    </div>
    <Button onclick={() => (isCreateOpen = true)}>
      <Plus class="mr-2 h-4 w-4" />
      Додати компанію
    </Button>
  </div>

  <div class="rounded-md border">
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Назва</Table.Head>
          <Table.Head>ЄДРПОУ</Table.Head>
          <Table.Head>Деталі</Table.Head>
          <Table.Head class="text-right">Дії</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#if data.companies && data.companies.length > 0}
          {#each data.companies as company (company.ID)}
            <Table.Row>
              <Table.Cell class="font-medium">{company.name}</Table.Cell>
              <Table.Cell>
                <code class="px-1 py-0.5 rounded bg-muted text-xs font-mono"
                  >{company.edrpou}</code
                >
              </Table.Cell>
              <Table.Cell>
                <span
                  class="text-xs text-muted-foreground italic truncate max-w-[200px] block"
                >
                  {Object.keys(company.details || {}).length > 0
                    ? "Є дані"
                    : "Немає даних"}
                </span>
              </Table.Cell>
              <Table.Cell class="text-right space-x-2">
                <Button
                  variant="ghost"
                  size="icon"
                  href={`/config/data/companies/${company.ID}`}
                  title="Деталі"
                >
                  <ExternalLink class="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  class="text-destructive hover:text-destructive"
                  onclick={() => openDelete(company)}
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
        <Dialog.Title>Додати компанію</Dialog.Title>
        <Dialog.Description
          >Введіть назву та код ЄДРПОУ для нової компанії.</Dialog.Description
        >
      </Dialog.Header>
      <form
        method="POST"
        action="?/create"
        use:enhance={() => {
          return async ({ result, update }) => {
            if (result.type === "success") {
              isCreateOpen = false;
              toast.success("Компанію успішно створено");
              await update();
            } else {
              const message = (result as { data?: { message?: string } }).data
                ?.message;
              console.error("Create company failed:", result);
              toast.error(mapErrorToFriendlyMessage(message));
            }
          };
        }}
        class="space-y-4 px-6 py-4"
      >
        <div class="space-y-2">
          <Label for="name">Назва</Label>
          <Input
            id="name"
            name="name"
            required
            placeholder="ТОВ 'Приклад'..."
          />
        </div>
        <div class="space-y-2">
          <Label for="edrpou">ЄДРПОУ</Label>
          <Input
            id="edrpou"
            name="edrpou"
            required
            placeholder="12345678"
            maxlength={10}
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

  <!-- Delete Dialog -->
  <Dialog.Root bind:open={isDeleteOpen}>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Видалити компанію?</Dialog.Title>
        <Dialog.Description>
          Ви впевнені, що хочете видалити компанію <strong
            >{selectedCompany?.name}</strong
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
                toast.success("Компанію видалено");
                await update();
              } else {
                const message = (result as { data?: { message?: string } }).data
                  ?.message;
                console.error("Delete company failed:", result);
                toast.error(mapErrorToFriendlyMessage(message));
              }
            };
          }}
        >
          <input
            type="hidden"
            name="id"
            value={selectedCompany?.ID ? String(selectedCompany.ID) : ""}
          />
          <Button type="submit" variant="destructive">Видалити</Button>
        </form>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
</div>
