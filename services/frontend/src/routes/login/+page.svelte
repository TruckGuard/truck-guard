<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { enhance } from "$app/forms";
  import type { ActionData } from "./$types";
  import { Eye, EyeOff } from "@lucide/svelte";

  let { form } = $props<ActionData>();

  let loading = $state(false);
  let showPassword = $state(false);
</script>

<div class="flex min-h-[100dvh] items-center justify-center bg-background p-4">
  <div class="w-full max-w-sm">
    <!-- Brand mark -->
    <div class="flex flex-col items-center gap-3 mb-8">
      <div class="w-11 h-11 rounded-lg bg-foreground text-background grid place-items-center font-bold text-base tracking-tight select-none">
        TG
      </div>
      <div class="text-center">
        <div class="text-base font-semibold tracking-tight">TruckGuard</div>
        <div class="text-xs text-muted-foreground uppercase tracking-widest mt-0.5">Митна система зважування</div>
      </div>
    </div>

    <Card.Root class="shadow-sm">
      <Card.Header class="pb-4">
        <Card.Title class="text-sm font-semibold">Вхід у систему</Card.Title>
        <Card.Description class="text-xs">
          Введіть ваші облікові дані для доступу.
        </Card.Description>
      </Card.Header>
      <Card.Content>
        <form
          method="POST"
          use:enhance={() => {
            loading = true;
            return async ({ update }) => {
              await update();
              loading = false;
            };
          }}
        >
          <div class="grid gap-3">
            <div class="flex flex-col gap-1.5">
              <Label for="username" class="text-xs font-medium">Логін</Label>
              <Input
                id="username"
                name="username"
                required
                disabled={loading}
                class="h-9 text-sm"
              />
            </div>
            <div class="flex flex-col gap-1.5">
              <Label for="password" class="text-xs font-medium">Пароль</Label>
              <div class="relative">
                <Input
                  id="password"
                  name="password"
                  type={showPassword ? "text" : "password"}
                  required
                  disabled={loading}
                  class="h-9 text-sm pr-9"
                />
                <Button
                  variant="ghost"
                  size="icon"
                  class="absolute right-0 top-0 h-9 w-9 text-muted-foreground hover:bg-transparent"
                  type="button"
                  onclick={() => (showPassword = !showPassword)}
                >
                  {#if showPassword}
                    <Eye class="h-3.5 w-3.5" />
                  {:else}
                    <EyeOff class="h-3.5 w-3.5" />
                  {/if}
                  <span class="sr-only">{showPassword ? "Сховати" : "Показати"} пароль</span>
                </Button>
              </div>
            </div>
          </div>

          {#if form?.status === 401}
            <div class="mt-3 rounded border border-destructive/30 bg-destructive/8 px-3 py-2 text-xs text-destructive">
              Невірний логін або пароль
            </div>
          {/if}

          {#if form?.status === 500}
            <div class="mt-3 rounded border border-destructive/30 bg-destructive/8 px-3 py-2 text-xs text-destructive">
              Виникла помилка на сервері
            </div>
          {/if}

          <Button class="mt-4 w-full h-9 text-sm" type="submit" disabled={loading}>
            {loading ? "Вхід..." : "Увійти"}
          </Button>
        </form>
      </Card.Content>
    </Card.Root>
  </div>
</div>
