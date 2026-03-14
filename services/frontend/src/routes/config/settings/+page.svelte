<script lang="ts">
    import { enhance } from '$app/forms';
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Card from "$lib/components/ui/card";
    import * as Table from "$lib/components/ui/table";
    import * as Alert from "$lib/components/ui/alert";
    import { Settings, Save, Info, CircleAlert, CheckCircle2 } from "@lucide/svelte";
    import type { PageData, ActionData } from './$types';

    let { data, form }: { data: PageData, form: ActionData } = $props();

    let settings = $state(data.settings || []);
    let selectedKey = $state<string | null>(settings.length > 0 ? settings[0].key : null);
    
    let selectedSetting = $derived(
        settings.find((s: any) => s.key === selectedKey) || null
    );

    let editValue = $state("");

    $effect(() => {
        if (selectedSetting) {
            editValue = selectedSetting.value;
        }
    });

    let isSubmitting = $state(false);
</script>

<div class="flex flex-col h-full space-y-6">
    <div class="flex items-center justify-between shrink-0">
        <h1 class="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Settings class="h-6 w-6 text-primary" />
            Системні налаштування
        </h1>
    </div>

    {#if data.error}
        <Alert.Root variant="destructive">
            <CircleAlert class="h-4 w-4" />
            <Alert.Title>Помилка</Alert.Title>
            <Alert.Description>{data.error}</Alert.Description>
        </Alert.Root>
    {/if}

    {#if form?.error}
        <Alert.Root variant="destructive">
            <CircleAlert class="h-4 w-4" />
            <Alert.Title>Помилка при оновленні</Alert.Title>
            <Alert.Description>{form.error}</Alert.Description>
        </Alert.Root>
    {/if}

    {#if form?.success}
        <Alert.Root class="border-green-500/50 bg-green-500/10 text-green-600 dark:text-green-400">
            <CheckCircle2 class="h-4 w-4" />
            <Alert.Title>Успішно</Alert.Title>
            <Alert.Description>Налаштування оновлено</Alert.Description>
        </Alert.Root>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 flex-1 overflow-hidden">
        <!-- Settings List -->
        <Card.Root class="md:col-span-1 flex flex-col overflow-hidden">
            <Card.Header>
                <Card.Title>Параметри</Card.Title>
                <Card.Description>Виберіть налаштування для редагування</Card.Description>
            </Card.Header>
            <Card.Content class="flex-1 overflow-y-auto p-0">
                <div class="flex flex-col">
                    {#each settings as setting}
                        <button 
                            class="flex flex-col items-start p-4 text-left transition-colors border-l-2 hover:bg-muted/50 {selectedKey === setting.key ? 'border-primary bg-primary/5' : 'border-transparent'}"
                            onclick={() => selectedKey = setting.key}
                        >
                            <span class="font-medium text-sm">{setting.name}</span>
                            <span class="text-xs text-muted-foreground line-clamp-1">{setting.key}</span>
                        </button>
                    {/each}
                </div>
            </Card.Content>
        </Card.Root>

        <!-- Setting Edit Form -->
        <Card.Root class="md:col-span-2 flex flex-col overflow-hidden">
            {#if selectedSetting}
                <Card.Header>
                    <div class="flex items-center justify-between">
                        <div>
                            <Card.Title>{selectedSetting.name}</Card.Title>
                            <Card.Description>{selectedSetting.key}</Card.Description>
                        </div>
                    </div>
                </Card.Header>
                <Card.Content class="space-y-6">
                    <div class="bg-muted/30 rounded-lg p-4 flex gap-3">
                        <Info class="h-5 w-5 text-primary shrink-0 mt-0.5" />
                        <div class="space-y-1">
                            <p class="text-sm font-medium">Опис</p>
                            <p class="text-sm text-muted-foreground">{selectedSetting.description}</p>
                        </div>
                    </div>

                    <form 
                        method="POST" 
                        use:enhance={() => {
                            isSubmitting = true;
                            return async ({ update }) => {
                                await update();
                                isSubmitting = false;
                                // Refresh local data after success
                                if (form?.success) {
                                    const idx = settings.findIndex((s: any) => s.key === selectedKey);
                                    if (idx !== -1) settings[idx].value = editValue;
                                }
                            };
                        }}
                        class="space-y-4"
                    >
                        <input type="hidden" name="key" value={selectedKey} />
                        
                        <div class="space-y-2">
                            <Label for="value">Поточне значення</Label>
                            <Input 
                                id="value" 
                                name="value" 
                                bind:value={editValue} 
                                placeholder="Введіть значення..."
                                class="h-11"
                            />
                        </div>

                        <div class="flex items-center gap-4 text-sm text-muted-foreground pt-2">
                            <span class="flex items-center gap-1.5">
                                <span class="font-medium text-foreground">За замовчуванням:</span>
                                {selectedSetting.default}
                            </span>
                        </div>

                        <div class="pt-4">
                            <Button type="submit" class="w-full md:w-auto" disabled={isSubmitting}>
                                {#if isSubmitting}
                                    Завантаження...
                                {:else}
                                    <Save class="mr-2 h-4 w-4" />
                                    Зберегти зміни
                                {/if}
                            </Button>
                        </div>
                    </form>
                </Card.Content>
            {:else}
                <div class="flex flex-col items-center justify-center flex-1 text-muted-foreground p-12">
                    <Settings class="h-12 w-12 opacity-20 mb-4" />
                    <p>Виберіть налаштування зі списку зліва</p>
                </div>
            {/if}
        </Card.Root>
    </div>
</div>
