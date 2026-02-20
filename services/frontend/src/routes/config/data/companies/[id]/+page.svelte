<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Textarea } from "$lib/components/ui/textarea";
  import * as Alert from "$lib/components/ui/alert";
  import { ChevronLeft, Save, AlertCircle } from "@lucide/svelte";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import * as Card from "$lib/components/ui/card";
  import type { PageData } from "./$types";
  import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let detailsStr = $state(JSON.stringify(data.company.details || {}, null, 2));

  let isUpdating = $state(false);
</script>

<div class="space-y-6">
  {#if data.error}
    <Alert.Root variant="destructive">
      <AlertCircle class="h-4 w-4" />
      <Alert.Title>Помилка</Alert.Title>
      <Alert.Description>{data.error}</Alert.Description>
    </Alert.Root>
  {/if}

  {#if data.company}
    <div class="flex items-center gap-4">
      <Button variant="outline" size="icon" href="/config/data/companies">
        <ChevronLeft class="h-4 w-4" />
      </Button>
      <div>
        <h1 class="text-3xl font-bold tracking-tight">{data.company.name}</h1>
        <p class="text-muted-foreground">
          Керування детальною інформацією про компанію.
        </p>
      </div>
    </div>

    <div class="grid gap-6 md:grid-cols-2">
      <Card.Root>
        <Card.Header>
          <Card.Title>Основна інформація</Card.Title>
          <Card.Description>Назва та юридичні реквізити.</Card.Description>
        </Card.Header>
        <Card.Content>
          <form
            method="POST"
            action="?/update"
            use:enhance={() => {
              isUpdating = true;
              return async ({ result, update }) => {
                isUpdating = false;
                if (result.type === "success") {
                  toast.success("Дані компанії оновлено");
                  await update();
                } else {
                  const message = (result as { data?: { message?: string } })
                    .data?.message;
                  console.error("Update company failed:", result);
                  toast.error(mapErrorToFriendlyMessage(message));
                }
              };
            }}
            class="space-y-4"
          >
            <div class="space-y-2">
              <Label for="name">Назва компанії</Label>
              <Input
                id="name"
                name="name"
                bind:value={data.company.name}
                required
              />
            </div>
            <div class="space-y-2">
              <Label for="edrpou">ЄДРПОУ</Label>
              <Input
                id="edrpou"
                name="edrpou"
                bind:value={data.company.edrpou}
                required
                maxlength={10}
              />
            </div>

            <div class="space-y-2">
              <Label for="details">Додаткові деталі (JSON)</Label>
              <Textarea
                id="details"
                name="details"
                bind:value={detailsStr}
                rows={10}
                class="font-mono text-sm"
                placeholder={"{}"}
              />
            </div>

            <div class="flex justify-end gap-2 pt-4">
              <Button variant="outline" href="/config/data/companies"
                >Скасувати</Button
              >
              <Button type="submit" disabled={isUpdating}>
                <Save class="mr-2 h-4 w-4" />
                {isUpdating ? "Збереження..." : "Зберегти зміни"}
              </Button>
            </div>
          </form>
        </Card.Content>
      </Card.Root>

      <Card.Root>
        <Card.Header>
          <Card.Title>Системна інформація</Card.Title>
          <Card.Description
            >Технічні дані та статус синхронізації.</Card.Description
          >
        </Card.Header>
        <Card.Content class="space-y-4">
          <div class="grid grid-cols-2 gap-2 text-sm">
            <span class="text-muted-foreground">ID в базі:</span>
            <span class="font-mono">{data.company.ID}</span>

            <span class="text-muted-foreground">Створено:</span>
            <span
              >{new Date(data.company.CreatedAt).toLocaleString("uk-UA")}</span
            >

            <span class="text-muted-foreground">Останнє оновлення:</span>
            <span
              >{new Date(data.company.UpdatedAt).toLocaleString("uk-UA")}</span
            >

            <span class="text-muted-foreground">Остання синхронізація:</span>
            <span>
              {data.company.last_synced_at
                ? new Date(data.company.last_synced_at).toLocaleString("uk-UA")
                : "Ніколи"}
            </span>
          </div>
        </Card.Content>
      </Card.Root>
    </div>
  {:else if !data.error}
    <div class="flex flex-col items-center justify-center py-12 text-center">
      <AlertCircle class="h-12 w-12 text-muted-foreground mb-4" />
      <h2 class="text-xl font-semibold italic">Компанію не знайдено</h2>
      <Button variant="outline" href="/config/data/companies" class="mt-4">
        До списку компаній
      </Button>
    </div>
  {/if}
</div>
