<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Textarea } from "$lib/components/ui/textarea";
    import { enhance } from "$app/forms";
    import { toast } from "svelte-sonner";
    import { ChevronLeft, Save, RefreshCw, Building2, Info, StickyNote } from "@lucide/svelte";
    import type { PageData, ActionData } from "./$types";
    import { formatDate } from "$lib/utils/date";
    import PageLayout from "$lib/components/common/PageLayout.svelte";

    let { data, form }: { data: PageData; form: ActionData } = $props();

    // Local optimistic state
    let company = $state({ ...data.company });
    $effect(() => { company = { ...data.company }; });

    let isSaving = $state(false);
    let isSyncing = $state(false);

    // Known ЄДРПОУ detail labels — displayed as structured rows
    const KNOWN_LABELS: string[] = [
        'Повна назва', 'Скорочена назва', 'Директор', 'Адреса',
        'Дата реєстрації', 'КВЕД', 'Статус', 'IBAN', 'Банк', 'МФО',
        'Телефон', 'Email',
    ];

    function detailRows(details: Record<string, any> | null | undefined): { label: string; value: string }[] {
        if (!details) return [];
        const entries = Object.entries(details).filter(([, v]) => v !== null && v !== '');
        // Sort: known labels first in defined order, then the rest alphabetically
        const known = KNOWN_LABELS
            .filter(l => details[l] !== undefined && details[l] !== null && details[l] !== '')
            .map(l => ({ label: l, value: String(details[l]) }));
        const unknown = entries
            .filter(([k]) => !KNOWN_LABELS.includes(k))
            .sort(([a], [b]) => a.localeCompare(b))
            .map(([k, v]) => ({ label: k, value: String(v) }));
        return [...known, ...unknown];
    }

    let rows = $derived(detailRows(company.details));

    function handleSave() {
        return ({ result, update }: any) => {
            isSaving = false;
            if (result.type === 'success') {
                if (result.data?.company) company = { ...result.data.company };
                toast.success('Дані компанії збережено');
            } else {
                const msg = result.data?.message || 'Помилка при збереженні';
                toast.error(msg);
            }
            update({ reset: false });
        };
    }

    function handleSync() {
        return ({ result, update }: any) => {
            isSyncing = false;
            if (result.type === 'success' && result.data?.synced) {
                // Optimistic update details & last_synced_at
                company = {
                    ...company,
                    details: result.data.details ?? company.details,
                    last_synced_at: result.data.last_synced_at ?? company.last_synced_at,
                };
                toast.success('Дані синхронізовано з ЄДР');
            } else {
                const msg = result.data?.message || 'Помилка синхронізації';
                toast.error(msg);
            }
            update({ reset: false });
        };
    }
</script>

<PageLayout>
    <!-- Header row (back + title + actions) -->
    <div class="flex items-center gap-3 shrink-0 pb-4 border-b border-border mb-4">
        <Button
            variant="ghost"
            size="icon"
            href="/config/data/companies"
            class="h-8 w-8 rounded-full border bg-background/60"
        >
            <ChevronLeft class="h-4 w-4" />
        </Button>

        <div class="flex-1 min-w-0">
            <h1 class="text-lg font-semibold tracking-tight text-foreground leading-snug truncate">
                {company.name}
            </h1>
            <p class="mt-0.5 text-xs text-muted-foreground flex items-center gap-1.5">
                <span class="font-mono bg-muted/50 px-1.5 py-0.5 rounded text-[10px]">ЄДРПОУ: {company.edrpou}</span>
                {#if company.last_synced_at}
                    <span class="text-muted-foreground/40">·</span>
                    <span>Синхронізовано: {formatDate(company.last_synced_at)}</span>
                {/if}
            </p>
        </div>

        <form
            method="POST"
            action="?/syncEdrpou"
            use:enhance={() => {
                isSyncing = true;
                toast.loading('Синхронізація з ЄДР…', { id: 'sync' });
                return (ctx) => { toast.dismiss('sync'); return handleSync()(ctx); };
            }}
        >
            <Button
                type="submit"
                variant="outline"
                size="sm"
                class="h-8 px-3 text-xs gap-1.5"
                disabled={isSyncing}
            >
                <RefreshCw class="h-3.5 w-3.5 {isSyncing ? 'animate-spin' : ''}" />
                {isSyncing ? 'Синхронізація…' : 'Синхронізувати з ЄДР'}
            </Button>
        </form>
    </div>

    <!-- Main content -->
    <div class="grid grid-cols-1 lg:grid-cols-5 gap-5">
        <!-- Left: Edit form -->
        <div class="lg:col-span-2 space-y-4">
            <!-- Basic info form -->
            <form
                method="POST"
                action="?/update"
                use:enhance={() => {
                    isSaving = true;
                    return handleSave();
                }}
                class="space-y-4"
            >
                <input type="hidden" name="details" value={JSON.stringify(company.details || {})} />

                <section class="rounded-xl border bg-card p-4 space-y-3">
                    <div class="flex items-center gap-2 mb-1">
                        <Building2 class="h-3.5 w-3.5 text-muted-foreground" />
                        <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Основна інформація</span>
                    </div>

                    <div class="space-y-1.5">
                        <Label for="name" class="text-xs">Назва компанії</Label>
                        <Input
                            id="name"
                            name="name"
                            bind:value={company.name}
                            required
                            class="h-9 bg-muted/30 border-none focus-visible:ring-1"
                        />
                    </div>

                    <div class="space-y-1.5">
                        <Label for="edrpou" class="text-xs">ЄДРПОУ</Label>
                        <Input
                            id="edrpou"
                            name="edrpou"
                            bind:value={company.edrpou}
                            required
                            maxlength={10}
                            class="h-9 font-mono bg-muted/30 border-none focus-visible:ring-1"
                        />
                    </div>

                    <div class="grid grid-cols-2 gap-3">
                        <div class="space-y-1.5">
                            <Label for="discount_percentage" class="text-xs">Знижка %</Label>
                            <Input
                                id="discount_percentage"
                                name="discount_percentage"
                                type="number"
                                min="0"
                                max="100"
                                step="0.01"
                                bind:value={company.discount_percentage}
                                class="h-9 bg-muted/30 border-none focus-visible:ring-1"
                            />
                        </div>
                        <div class="space-y-1.5">
                            <Label for="discount_fixed" class="text-xs">Фікс. знижка</Label>
                            <Input
                                id="discount_fixed"
                                name="discount_fixed"
                                type="number"
                                min="0"
                                step="0.01"
                                bind:value={company.discount_fixed}
                                class="h-9 bg-muted/30 border-none focus-visible:ring-1"
                            />
                        </div>
                    </div>
                </section>

                <!-- Notes section -->
                <section class="rounded-xl border bg-card p-4 space-y-3">
                    <div class="flex items-center gap-2 mb-1">
                        <StickyNote class="h-3.5 w-3.5 text-muted-foreground" />
                        <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Примітки</span>
                    </div>
                    <Textarea
                        name="notes"
                        bind:value={company.notes}
                        rows={4}
                        placeholder="Внутрішні нотатки про компанію..."
                        class="text-sm bg-muted/30 border-none focus-visible:ring-1 resize-none"
                    />
                </section>

                <div class="flex justify-end">
                    <Button type="submit" size="sm" class="h-8 px-4 text-xs gap-1.5" disabled={isSaving}>
                        <Save class="h-3.5 w-3.5 {isSaving ? 'animate-pulse' : ''}" />
                        {isSaving ? 'Збереження…' : 'Зберегти'}
                    </Button>
                </div>
            </form>
        </div>

        <!-- Right: Details + system info -->
        <div class="lg:col-span-3 space-y-4">
            <!-- ЄДРПОУ details -->
            <section class="rounded-xl border bg-card p-4">
                <div class="flex items-center gap-2 mb-3">
                    <Info class="h-3.5 w-3.5 text-muted-foreground" />
                    <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Дані з ЄДР</span>
                </div>

                {#if rows.length > 0}
                    <dl class="divide-y divide-border/50">
                        {#each rows as row}
                            <div class="flex gap-3 py-2 text-sm">
                                <dt class="w-36 shrink-0 text-xs text-muted-foreground">{row.label}</dt>
                                <dd class="flex-1 text-xs font-medium text-foreground break-words">{row.value}</dd>
                            </div>
                        {/each}
                    </dl>
                {:else}
                    <div class="flex flex-col items-center py-8 text-center text-muted-foreground/50">
                        <RefreshCw class="h-8 w-8 mb-2 opacity-20" />
                        <p class="text-xs">Дані ще не синхронізовано.</p>
                        <p class="text-xs opacity-60 mt-0.5">Натисніть «Синхронізувати з ЄДР» вгорі.</p>
                    </div>
                {/if}
            </section>

            <!-- System info -->
            <section class="rounded-xl border bg-card p-4">
                <div class="flex items-center gap-2 mb-3">
                    <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Системна інформація</span>
                </div>
                <dl class="divide-y divide-border/50">
                    {#each [
                        ['ID', String(company.ID)],
                        ['Створено', formatDate(company.CreatedAt)],
                        ['Оновлено', formatDate(company.UpdatedAt)],
                        ['Остання синхронізація', company.last_synced_at ? formatDate(company.last_synced_at) : 'Ніколи'],
                    ] as [label, value]}
                        <div class="flex gap-3 py-2">
                            <dt class="w-36 shrink-0 text-xs text-muted-foreground">{label}</dt>
                            <dd class="text-xs font-mono text-foreground">{value}</dd>
                        </div>
                    {/each}
                </dl>
            </section>
        </div>
    </div>
</PageLayout>
