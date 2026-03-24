<script lang="ts">
    import { enhance } from '$app/forms';
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Dialog from "$lib/components/ui/dialog";
    import * as Select from "$lib/components/ui/select";
    import { Check, Copy } from "@lucide/svelte";
    import { toast } from "svelte-sonner";

    let { open = $bindable(false), posts = [], onCreated } = $props<{
        open: boolean;
        posts: any[];
        onCreated?: (camera: any, apiKey: string) => void;
    }>();

    let selectedType = $state("");
    let selectedPostId = $state("");
    let copied = $state(false);
    let isConfigOpen = $state(false);
    let generatedKey = $state("");
    let createdCamera = $state<any>(null);

    const typeLabel = $derived(selectedType === 'front' ? 'Передня (Front)' : selectedType === 'back' ? 'Задня (Back)' : 'Виберіть тип');
    const postLabel = $derived(posts.find((p: any) => p.ID.toString() === selectedPostId)?.name || 'Виберіть пост');

    const getEndpointUrl = (key: string) => {
        const host = typeof window !== 'undefined' ? window.location.origin : 'http://truckguard.local';
        return `${host}/ingest/plate?key=${key}`;
    };

    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
        copied = true;
        toast.success("Скопійовано");
        setTimeout(() => copied = false, 2000);
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-[500px] border-none shadow-2xl bg-background/95 backdrop-blur-xl">
        <Dialog.Header>
            <Dialog.Title class="text-2xl font-bold">Нова камера</Dialog.Title>
            <Dialog.Description>
                Система автоматично згенерує API-ключ для цього пристрою.
            </Dialog.Description>
        </Dialog.Header>
        <form
            method="POST"
            action="?/create"
            use:enhance={() => {
                return async ({ result, update }) => {
                    if (result.type === 'success' && result.data?.api_key) {
                        generatedKey = `${result.data.api_key}`;
                        createdCamera = result.data.camera;
                        open = false;
                        isConfigOpen = true;
                        onCreated?.(createdCamera, generatedKey);
                    } else {
                        await update();
                    }
                };
            }}
            class="space-y-5 pt-4"
        >
            <div class="space-y-2">
                <Label for="cam-name" class="text-xs font-semibold uppercase tracking-wider text-muted-foreground ml-1">Назва камери</Label>
                <Input id="cam-name" name="name" placeholder="Наприклад: В'їзд Основний" required class="h-11 bg-muted/20 border-none rounded-lg" />
            </div>
            <div class="space-y-2">
                <Label for="cam-description" class="text-xs font-semibold uppercase tracking-wider text-muted-foreground ml-1">Опис</Label>
                <Input id="cam-description" name="description" placeholder="Додаткова інформація" class="h-11 bg-muted/20 border-none rounded-lg" />
            </div>
            <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                    <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground ml-1">Тип розпізнавання</Label>
                    <Select.Root name="type" type="single" bind:value={selectedType} required>
                        <Select.Trigger class="h-11 bg-muted/20 border-none rounded-lg">
                            <span class={!selectedType ? 'text-muted-foreground' : ''}>{typeLabel}</span>
                        </Select.Trigger>
                        <Select.Content class="bg-background/95 backdrop-blur-xl">
                            <Select.Item value="front">Передня (Front)</Select.Item>
                            <Select.Item value="back">Задня (Back)</Select.Item>
                        </Select.Content>
                    </Select.Root>
                </div>
                <div class="space-y-2">
                    <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground ml-1">Митний пост</Label>
                    <Select.Root name="customs_post_id" type="single" bind:value={selectedPostId}>
                        <Select.Trigger class="h-11 bg-muted/20 border-none rounded-lg">
                            <span class={!selectedPostId ? 'text-muted-foreground' : ''}>{postLabel}</span>
                        </Select.Trigger>
                        <Select.Content class="bg-background/95 backdrop-blur-xl">
                            {#each posts as post}
                                <Select.Item value={post.ID.toString()}>{post.name}</Select.Item>
                            {/each}
                        </Select.Content>
                    </Select.Root>
                </div>
            </div>
            <div class="flex items-center gap-3 p-4 bg-primary/5 rounded-xl border border-primary/10">
                <input type="checkbox" id="cam-match_permit" name="match_permit" class="h-5 w-5 rounded border-gray-300 text-primary focus:ring-primary bg-muted/30 cursor-pointer" checked />
                <div class="cursor-pointer">
                    <Label for="cam-match_permit" class="font-bold text-sm cursor-pointer">Авто-прив'язка до перепусток</Label>
                    <p class="text-[10px] text-muted-foreground mt-0.5">Автоматично шукати активні перепустки за розпізнаним номером.</p>
                </div>
            </div>
            <Dialog.Footer class="pt-4">
                <Button type="button" variant="ghost" class="rounded-lg h-11" onclick={() => (open = false)}>Скасувати</Button>
                <Button type="submit" class="rounded-lg h-11 px-8 font-bold shadow-lg shadow-primary/20">Зберегти та отримати ключ</Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>

<!-- API Key Reveal Modal -->
<Dialog.Root bind:open={isConfigOpen}>
    <Dialog.Content class="sm:max-w-[600px] border-none shadow-2xl bg-slate-950 text-white overflow-hidden p-0">
        <div class="absolute inset-0 bg-linear-to-br from-primary/20 via-transparent to-purple-500/20 pointer-events-none"></div>
        <div class="p-8 relative space-y-8">
            <div class="flex items-center gap-4">
                <div class="h-12 w-12 rounded-2xl bg-primary/20 text-primary flex items-center justify-center border border-primary/30">
                    <Check class="h-7 w-7" />
                </div>
                <div>
                    <h2 class="text-2xl font-bold">Камеру налаштовано</h2>
                    <p class="text-slate-400">Використовуйте ці дані для конфігурації ANPR-камери.</p>
                </div>
            </div>
            <div class="space-y-6">
                <div class="space-y-2">
                    <div class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Device ID / Username</div>
                    <div class="bg-slate-900 border border-slate-800 rounded-xl px-4 h-12 flex items-center font-mono text-slate-200">
                        {createdCamera?.camera_id}
                    </div>
                </div>
                <div class="space-y-2">
                    <div class="text-[10px] font-bold uppercase tracking-widest text-slate-400">API Key / Password</div>
                    <div class="relative group">
                        <div class="bg-slate-900 border border-primary/30 rounded-xl px-4 py-4 font-mono text-lg text-white break-all pr-14 leading-relaxed group-hover:border-primary/50 transition-colors bg-linear-to-r from-primary/5 to-transparent">
                            {generatedKey}
                        </div>
                        <Button class="absolute right-2 top-2 h-10 w-10 rounded-lg bg-primary hover:bg-primary/90 text-white border-none" onclick={() => copyToClipboard(generatedKey)}>
                            {#if copied}<Check class="h-4 w-4" />{:else}<Copy class="h-4 w-4" />{/if}
                        </Button>
                    </div>
                </div>
                <div class="space-y-3 p-4 bg-slate-900 rounded-2xl border border-slate-800">
                    <div class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Platform Server (URL)</div>
                    <div class="flex items-center justify-between gap-3 font-mono text-sm text-slate-300 bg-black/40 p-3 rounded-lg flex-wrap">
                        <span class="break-all">{getEndpointUrl(generatedKey)}</span>
                        <Button variant="ghost" size="sm" class="text-slate-500 h-8 px-2 hover:bg-white/5" onclick={() => copyToClipboard(getEndpointUrl(generatedKey))}>
                            <Copy class="h-3.5 w-3.5" />
                        </Button>
                    </div>
                </div>
            </div>
            <div class="pt-4 border-t border-slate-800 flex justify-between items-center">
                <div class="text-[11px] text-slate-500">Збережіть цей ключ. Ви не зможете побачити його знову.</div>
                <Button class="rounded-xl px-10 h-12 font-bold bg-primary hover:bg-primary/90 border-none" onclick={() => (isConfigOpen = false)}>
                    Готово
                </Button>
            </div>
        </div>
    </Dialog.Content>
</Dialog.Root>
