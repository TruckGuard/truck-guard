<script lang="ts">
    import { enhance } from '$app/forms';
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Settings, Save, Info } from "@lucide/svelte";
    import { toast } from "svelte-sonner";
    import type { PageData, ActionData } from './$types';
    import PageHeader from "$lib/components/common/PageHeader.svelte";
    import PageLayout from "$lib/components/common/PageLayout.svelte";

    let { data, form }: { data: PageData, form: ActionData } = $props();

    let settings = $state<any[]>(data.settings || []);
    $effect(() => { settings = data.settings || []; });

    let selectedKey = $state<string | null>(null);
    $effect(() => {
        if (!selectedKey && settings.length > 0) selectedKey = settings[0].key;
    });

    let selectedSetting = $derived(settings.find((s) => s.key === selectedKey) ?? null);

    let editValue = $state('');
    $effect(() => { if (selectedSetting) editValue = selectedSetting.value; });

    let isSubmitting = $state(false);

    // Show form result toasts without re-rendering alert banners
    $effect(() => {
        if (form?.success) toast.success('Налаштування оновлено');
        if (form?.error) toast.error(form.error);
    });
</script>

<PageLayout>
    <PageHeader
        title="Налаштування системи"
        description="Конфігурація глобальних параметрів та системних лімітів."
    />

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 min-h-0">
        <!-- Settings list -->
        <div class="md:col-span-1 rounded-xl border bg-card overflow-hidden flex flex-col">
            <div class="px-4 py-3 border-b bg-muted/20 shrink-0">
                <p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Параметри</p>
            </div>
            <div class="flex-1 overflow-y-auto scrollbar-thin">
                {#each settings as setting}
                    <button
                        class="w-full flex flex-col items-start px-4 py-3 text-left transition-all border-l-2 hover:bg-muted/40
                               {selectedKey === setting.key
                                   ? 'border-primary bg-primary/5'
                                   : 'border-transparent opacity-70 hover:opacity-100'}"
                        onclick={() => { selectedKey = setting.key; editValue = setting.value; }}
                    >
                        <span class="text-sm font-medium">{setting.name}</span>
                        <span class="text-[11px] font-mono text-muted-foreground">{setting.key}</span>
                    </button>
                {/each}
            </div>
        </div>

        <!-- Edit panel -->
        <div class="md:col-span-2 rounded-xl border bg-card overflow-hidden flex flex-col">
            {#if selectedSetting}
                <div class="px-5 py-4 border-b shrink-0">
                    <p class="text-sm font-semibold text-foreground">{selectedSetting.name}</p>
                    <p class="text-[11px] font-mono text-muted-foreground mt-0.5">{selectedSetting.key}</p>
                </div>

                <div class="flex-1 p-5 space-y-5 overflow-y-auto">
                    <!-- Description block -->
                    <div class="flex gap-3 rounded-lg bg-muted/30 border border-muted/40 p-3">
                        <div class="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                            <Info class="h-4 w-4 text-primary" />
                        </div>
                        <div>
                            <p class="text-xs font-semibold mb-0.5">Опис</p>
                            <p class="text-xs text-muted-foreground leading-relaxed">{selectedSetting.description}</p>
                        </div>
                    </div>

                    <form
                        method="POST"
                        use:enhance={() => {
                            isSubmitting = true;
                            // Optimistic: immediately reflect value in list
                            const key = selectedKey;
                            const val = editValue;
                            const idx = settings.findIndex((s) => s.key === key);
                            if (idx !== -1) settings[idx] = { ...settings[idx], value: val };

                            return async ({ result, update }) => {
                                isSubmitting = false;
                                if (result.type === 'success') {
                                    toast.success('Налаштування оновлено');
                                } else {
                                    // Revert optimistic update on error
                                    await update({ reset: false });
                                    toast.error((result as any).data?.error || 'Помилка збереження');
                                }
                            };
                        }}
                        class="space-y-4"
                    >
                        <input type="hidden" name="key" value={selectedKey} />

                        <div class="space-y-2">
                            <Label for="value" class="text-xs font-medium">Значення параметра</Label>
                            <Input
                                id="value"
                                name="value"
                                bind:value={editValue}
                                placeholder="Введіть значення..."
                                class="h-10 bg-muted/20 border-muted focus:border-primary/50 focus:ring-primary/20"
                            />
                        </div>

                        <div class="flex items-center gap-2 px-3 py-2 rounded-lg bg-muted/10 border border-dotted border-muted text-xs text-muted-foreground">
                            <span class="font-semibold text-foreground">За замовчуванням:</span>
                            <code class="px-1.5 py-0.5 rounded bg-muted/30">{selectedSetting.default}</code>
                        </div>

                        <div class="flex justify-end pt-2">
                            <Button
                                type="submit"
                                size="sm"
                                class="h-8 px-4 text-xs gap-1.5"
                                disabled={isSubmitting}
                            >
                                {#if isSubmitting}
                                    <div class="h-3.5 w-3.5 border-2 border-white/30 border-t-white animate-spin rounded-full"></div>
                                    Збереження…
                                {:else}
                                    <Save class="h-3.5 w-3.5" />
                                    Зберегти зміни
                                {/if}
                            </Button>
                        </div>
                    </form>
                </div>
            {:else}
                <div class="flex flex-col items-center justify-center flex-1 text-muted-foreground p-12">
                    <div class="h-16 w-16 rounded-full bg-muted/20 flex items-center justify-center mb-4">
                        <Settings class="h-8 w-8 opacity-30 animate-[spin_12s_linear_infinite]" />
                    </div>
                    <p class="text-sm font-medium">Виберіть параметр зі списку</p>
                    <p class="text-xs text-center mt-1 opacity-60 max-w-xs">Виберіть налаштування зліва для редагування</p>
                </div>
            {/if}
        </div>
    </div>
</PageLayout>
