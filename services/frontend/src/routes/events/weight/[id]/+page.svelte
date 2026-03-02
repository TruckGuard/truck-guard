<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { Label } from "$lib/components/ui/label";
    import { ChevronLeft, Clock, Scale } from "@lucide/svelte";

    let { data } = $props();
    let event = data.event;

    function formatDate(dateStr: string) {
        if (!dateStr) return "-";
        return new Date(dateStr).toLocaleString("uk-UA");
    }
</script>

<div class="container mx-auto py-6 space-y-6 max-w-5xl">
    <!-- Header -->
    <div class="flex items-center gap-4">
        <Button variant="outline" size="icon" href="/events" class="h-8 w-8">
            <ChevronLeft class="h-4 w-4" />
        </Button>
        <div>
            <h1 class="text-2xl font-bold tracking-tight text-foreground">
                Подія зважування #{event.ID}
            </h1>
            <p class="text-muted-foreground text-sm">
                Детальна інформація про подію
            </p>
        </div>
        <div
            class="ml-auto bg-blue-50 dark:bg-blue-950/30 text-blue-700 dark:text-blue-400 px-3 py-1 rounded-full text-xs font-bold border border-blue-200 dark:border-blue-800"
        >
            WEIGHT
        </div>
    </div>

    <div class="grid gap-6 md:grid-cols-2">
        <!-- Main Info -->
        <Card.Root class="shadow-md md:col-span-1">
            <Card.Header class="pb-3">
                <Card.Title>Основні дані</Card.Title>
            </Card.Header>
            <Card.Content class="space-y-6">
                <div class="space-y-4">
                    <!-- Weight Block -->
                    <div
                        class="bg-muted/30 p-4 rounded-xl border border-border/50"
                    >
                        <div class="flex items-baseline gap-2">
                            <span
                                class="text-4xl font-black font-mono tracking-tight text-blue-700 dark:text-blue-400"
                            >
                                {event.weight}
                            </span>
                            <span
                                class="text-xl font-medium text-muted-foreground"
                                >кг</span
                            >
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
                                >Ваги</Label
                            >
                            <div class="flex items-center gap-2 font-medium">
                                <Button
                                    href={`/scales/${event.scale_id}`}
                                    variant="link"
                                >
                                <Scale class="h-4 w-4 text-muted-foreground" />
                                    {event.scale_source_name || event.scale_id}
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            </Card.Content>
        </Card.Root>

        <!-- Side Info -->
        <div class="md:col-span-1 space-y-6">
            <Card.Root class="border-0 shadow-sm bg-muted/10 h-full">
                <Card.Header>
                    <Card.Title
                        class="text-sm font-medium text-muted-foreground"
                        >Зв'язки</Card.Title
                    >
                </Card.Header>
                <Card.Content>
                    <dl class="grid grid-cols-2 gap-6">
                        <div>
                            <dt
                                class="text-xs font-medium text-muted-foreground"
                            >
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
                            <dt
                                class="text-xs font-medium text-muted-foreground"
                            >
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
