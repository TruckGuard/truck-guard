<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import { Separator } from "$lib/components/ui/separator";
  import { ChevronLeft, Clock, Camera, CreditCard, Eye } from "@lucide/svelte";
  import { can } from "$lib/auth";
  import { toast } from "svelte-sonner";
  import { formatDate } from "$lib/utils/date";

  let { data } = $props();
  // Initialize as state for optimistic updates
  let event = data.event;
  const user = $derived(data.user);
</script>

<div class="container mx-auto py-6 space-y-6 max-w-5xl">
  <!-- Header -->
  <div class="flex items-center gap-4">
    <Button variant="outline" size="icon" href="/events" class="h-8 w-8">
      <ChevronLeft class="h-4 w-4" />
    </Button>
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-foreground">
        Подія розпізнавання #{event.ID}
      </h1>
      <p class="text-muted-foreground text-sm">Детальна інформація про подію</p>
    </div>
    <div
      class="ml-auto bg-blue-50 dark:bg-blue-950/30 text-blue-700 dark:text-blue-400 px-3 py-1 rounded-full text-xs font-bold border border-blue-200 dark:border-blue-800"
    >
      ANPR
    </div>
  </div>

  <div class="grid gap-6 md:grid-cols-2">
    <!-- Image Section -->
    <Card.Root
      class="overflow-hidden md:col-span-1 shadow-lg bg-card/50 gap-0 py-0"
    >
      <Card.Content
        class="p-0 relative aspect-video bg-black/5 flex items-center justify-center h-full min-h-[300px]"
      >
        {#if event.image_key}
          <img
            src={`/api/images/${event.image_key}`}
            alt="Vehicle Plate"
            class="w-full h-full object-cover"
          />
          <div
            class="absolute inset-0 opacity-20 hover:opacity-100 transition-opacity flex items-end justify-center p-4"
          >
            <Button
              variant="secondary"
              size="sm"
              class="gap-2"
              href={`/api/images/${event.image_key}`}
              target="_blank"
            >
              <Eye class="h-4 w-4" /> Відкрити оригінал
            </Button>
          </div>
        {:else}
          <div class="flex flex-col items-center gap-2 text-muted-foreground">
            <Camera class="h-12 w-12 opacity-20" />
            <span class="text-sm">Зображення відсутнє</span>
          </div>
        {/if}
      </Card.Content>
    </Card.Root>

    <!-- Details Section -->
    <div class="md:col-span-1 space-y-6">
      <Card.Root class="shadow-md">
        <Card.Header class="pb-3">
          <Card.Title>Основні дані</Card.Title>
        </Card.Header>
        <Card.Content class="space-y-6">
          <div class="space-y-4">
            <!-- Plate Number Block -->
            <div class="bg-muted/30 p-4 rounded-xl border border-border/50">
              <div class="flex items-baseline gap-3">
                <span
                  class="text-3xl font-black font-mono tracking-wider text-foreground"
                >
                  {event.plate}
                </span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1">
                <Label
                  class="text-xs text-muted-foreground uppercase tracking-wider"
                  >Час фіксації</Label
                >
                <div class="flex items-center gap-2 font-medium">
                  <Clock class="h-4 w-4 text-muted-foreground" />
                  {formatDate(event.timestamp)}
                </div>
              </div>
              <div class="space-y-1">
                <Label
                  class="text-xs text-muted-foreground uppercase tracking-wider"
                  >Камера</Label
                >
                <Button href={`/cameras/${event.camera_id}`} variant="link">
                  <Camera class="h-4 w-4 text-muted-foreground" />
                  {event.camera_source_name || event.camera_id}
                </Button>
              </div>
            </div>
          </div>
        </Card.Content>
      </Card.Root>

      <!-- Additional Info (JSON Payload or other) -->
      <Card.Root class="border-0 shadow-sm bg-muted/10">
        <Card.Header>
          <Card.Title class="text-sm font-medium text-muted-foreground"
            >Додаткова інформація</Card.Title
          >
        </Card.Header>
        <Card.Content>
          <dl class="grid grid-cols-2 gap-6">
            <div>
              <dt class="text-xs font-medium text-muted-foreground">
                Системна подія
              </dt>
              <dd class="mt-1 text-sm text-foreground font-mono">
                {#if event.system_event_id}
                  <a
                    href={`/events/system/${event.system_event_id}`}
                    class="underline hover:text-blue-500 transition-colors"
                    >#{event.system_event_id}</a
                  >
                {:else}
                  -
                {/if}
              </dd>
            </div>
            <div>
              <dt class="text-xs font-medium text-muted-foreground">
                Перепустка
              </dt>
              <dd class="mt-1 text-sm text-foreground font-mono">
                {#if event.permit_id}
                  <a
                    href={`/permits/${event.permit_id}`}
                    class="underline hover:text-blue-500 transition-colors"
                    >#{event.permit_id}</a
                  >
                {:else}
                  -
                {/if}
              </dd>
            </div>
          </dl>
        </Card.Content>
      </Card.Root>
    </div>
  </div>
</div>
