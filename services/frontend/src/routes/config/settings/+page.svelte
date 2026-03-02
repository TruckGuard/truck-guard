<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";
  import * as Card from "$lib/components/ui/card";
  import * as Alert from "$lib/components/ui/alert";
  import type { PageData } from "./$types";
  import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";
  import { CircleAlert, Info, Save } from "@lucide/svelte";

  let { data }: { data: PageData & { error?: string | null } } = $props();

  let isUpdating = $state(false);
</script>

<div class="space-y-6">
  {#if data.error}
    <Alert.Root variant="destructive">
      <CircleAlert class="h-4 w-4" />
      <Alert.Title>Помилка</Alert.Title>
      <Alert.Description>{data.error}</Alert.Description>
    </Alert.Root>
  {/if}

  <div>
    <h1 class="text-3xl font-bold tracking-tight">Системні налаштування</h1>
    <p class="text-muted-foreground">
      Керування глобальними параметрами системи.
    </p>
  </div>

  <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
    {#if data.settings && data.settings.length > 0}
      {#each data.settings as setting (setting.key)}
        <Card.Root>
          <Card.Header>
            <Card.Title>{setting.name || setting.key}</Card.Title>
            {#if setting.description}
              <Card.Description>{setting.description}</Card.Description>
            {/if}
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
                    toast.success("Налаштування збережено");
                    await update();
                  } else {
                    const message = (result as { data?: { message?: string } })
                      .data?.message;
                    console.error("Update setting failed:", result);
                    toast.error(mapErrorToFriendlyMessage(message));
                  }
                };
              }}
              class="space-y-4"
            >
              <input type="hidden" name="key" value={setting.key} />
              <div class="space-y-2">
                <Label for={`setting-${setting.key}`}>Значення</Label>
                <Input
                  id={`setting-${setting.key}`}
                  name="value"
                  bind:value={setting.value}
                />
              </div>
              <Button type="submit" class="w-full" disabled={isUpdating}>
                <Save class="mr-2 h-4 w-4" />
                Зберегти
              </Button>
            </form>
          </Card.Content>
        </Card.Root>
      {/each}
    {:else}
      <Card.Root class="col-span-full">
        <Card.Content
          class="flex flex-col items-center justify-center py-12 text-center"
        >
          <Info class="h-12 w-12 text-muted-foreground mb-4" />
          <h2 class="text-xl font-semibold italic">Налаштувань не знайдено</h2>
          <p class="text-muted-foreground">
            Система використовує параметри за замовчуванням.
          </p>
        </Card.Content>
      </Card.Root>
    {/if}
  </div>
</div>
