<script lang="ts">
    import { enhance } from '$app/forms';
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Card from "$lib/components/ui/card";
    import * as Alert from "$lib/components/ui/alert";
    import { Settings, Save, Info, CircleAlert, CheckCircle2 } from "@lucide/svelte";
    import type { PageData, ActionData } from './$types';
    import PageHeader from "$lib/components/common/PageHeader.svelte";
    import PageLayout from "$lib/components/common/PageLayout.svelte";

    let { data, form }: { data: PageData, form: ActionData } = $props();

    let settings = $state(data.settings || []);
    $effect(() => { settings = data.settings || []; });
    let selectedKey = $state<string | null>(null);
    
    $effect(() => {
        if (!selectedKey && settings.length > 0) {
            selectedKey = settings[0].key;
        }
    });
    
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

<PageLayout>
    <PageHeader
        title="Налаштування системи"
        description="Конфігурація глобальних параметрів та системних лімітів."
    />

    {#if data.error}
        <Alert.Root variant="destructive" class="shrink-0">
            <CircleAlert class="h-4 w-4" />
            <Alert.Title>Помилка</Alert.Title>
            <Alert.Description>{data.error}</Alert.Description>
        </Alert.Root>
    {/if}

    {#if form?.error}
        <Alert.Root variant="destructive" class="shrink-0">
            <CircleAlert class="h-4 w-4" />
            <Alert.Title>Помилка при оновленні</Alert.Title>
            <Alert.Description>{form.error}</Alert.Description>
        </Alert.Root>
    {/if}

    {#if form?.success}
        <Alert.Root class="border-green-500/50 bg-green-500/10 text-green-600 dark:text-green-400 shrink-0">
            <CheckCircle2 class="h-4 w-4" />
            <Alert.Title>Успішно</Alert.Title>
            <Alert.Description>Налаштування оновлено</Alert.Description>
        </Alert.Root>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 flex-1 overflow-hidden min-h-0">
        <!-- Settings List -->
        <Card.Root class="md:col-span-1 flex flex-col overflow-hidden bg-card/50 border-none shadow-sm">
            <Card.Header>
                <Card.Title>Параметри</Card.Title>
                <Card.Description>Виберіть налаштування для редагування</Card.Description>
            </Card.Header>
            <Card.Content class="flex-1 overflow-y-auto p-0 scrollbar-thin scrollbar-thumb-muted">
                <div class="flex flex-col">
                    {#each settings as setting}
                        <button 
                            class="flex flex-col items-start p-4 text-left transition-all border-l-4 hover:bg-muted/50 {selectedKey === setting.key ? 'border-primary bg-primary/5' : 'border-transparent opacity-70 hover:opacity-100'}"
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
        <Card.Root class="md:col-span-2 flex flex-col overflow-hidden bg-white dark:bg-card border-none shadow-sm">
            {#if selectedSetting}
                <Card.Header>
                    <div class="flex items-center justify-between">
                        <div>
                            <Card.Title class="text-xl">{selectedSetting.name}</Card.Title>
                            <Card.Description class="font-mono text-xs">{selectedSetting.key}</Card.Description>
                        </div>
                    </div>
                </Card.Header>
                <Card.Content class="space-y-6 overflow-y-auto">
                    <div class="bg-muted/30 rounded-xl p-4 flex gap-4 border border-muted/50">
                        <div class="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                            <Info class="h-5 w-5 text-primary" />
                        </div>
                        <div class="space-y-1">
                            <p class="text-sm font-semibold">Опис</p>
                            <p class="text-sm text-muted-foreground leading-relaxed">{selectedSetting.description}</p>
                        </div>
                    </div>

                    <form 
                        method="POST" 
                        use:enhance={() => {
                            isSubmitting = true;
                            return async ({ update }) => {
                                await update();
                                isSubmitting = false;
                                if (form?.success) {
                                    const idx = settings.findIndex((s: any) => s.key === selectedKey);
                                    if (idx !== -1) settings[idx].value = editValue;
                                }
                            };
                        }}
                        class="space-y-6"
                    >
                        <input type="hidden" name="key" value={selectedKey} />
                        
                        <div class="space-y-3">
                            <Label for="value" class="text-sm font-medium">Значення параметра</Label>
                            <Input 
                                id="value" 
                                name="value" 
                                bind:value={editValue} 
                                placeholder="Введіть значення..."
                                class="h-12 bg-muted/20 border-muted focus:border-primary/50 focus:ring-primary/20 transition-all"
                            />
                        </div>

                        <div class="flex items-center gap-4 p-3 rounded-lg bg-muted/10 border border-dotted border-muted text-sm text-muted-foreground">
                            <div class="flex items-center gap-2">
                                <span class="font-semibold text-foreground">За замовчуванням:</span>
                                <code class="px-1.5 py-0.5 rounded bg-muted/30 text-xs">{selectedSetting.default}</code>
                            </div>
                        </div>

                        <div class="pt-4 flex justify-end">
                            <Button type="submit" class="h-11 px-8 shadow-md shadow-primary/20 transition-all hover:-translate-y-px active:translate-y-0" disabled={isSubmitting}>
                                {#if isSubmitting}
                                    <div class="h-4 w-4 border-2 border-white/30 border-t-white animate-spin rounded-full mr-2"></div>
                                    Збереження...
                                {:else}
                                    <Save class="mr-2 h-4 w-4" />
                                    Зберегти зміни
                                {/if}
                            </Button>
                        </div>
                    </form>
                </Card.Content>
            {:else}
                <div class="flex flex-col items-center justify-center flex-1 text-muted-foreground p-12 bg-muted/5">
                    <div class="h-20 w-20 rounded-full bg-muted/20 flex items-center justify-center mb-6">
                        <Settings class="h-10 w-10 opacity-40 animate-[spin_10s_linear_infinite]" />
                    </div>
                    <p class="font-medium">Виберіть налаштування зі списку зліва</p>
                    <p class="text-xs max-w-xs text-center mt-2 opacity-60">Використовуйте довідник параметрів для конфігурації системної логіки</p>
                </div>
            {/if}
        </Card.Root>
    </div>
</PageLayout>
