<script lang="ts">
    import {
        Activity,
        Clock,
        Coins,
        CircleCheck,
        Save,
        Send,
    } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";

    let {
        permit,
        data,
        loading,
        handleSave,
        handleValidatePermit,
        handleClosePermit,
        validationItems,
    } = $props<{
        permit: any;
        data: any;
        loading: boolean;
        handleSave: (silent?: boolean) => Promise<boolean>;
        handleValidatePermit: () => Promise<void>;
        handleClosePermit: () => Promise<void>;
        validationItems: {
            label: string;
            isValid: boolean;
            targetId?: string;
        }[];
    }>();

    function formatDate(dateStr?: string) {
        if (!dateStr) return "-";
        return new Date(dateStr).toLocaleString("uk-UA");
    }
</script>

<div class="lg:col-span-4 space-y-6 sticky top-8">
    {#if !data.isNew}
        <!-- Status & Info Card -->
        <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
            <div
                class="p-5 bg-muted/30 border-b flex items-center justify-between"
            >
                <h3
                    class="font-bold text-[10px] uppercase tracking-widest text-muted-foreground flex items-center gap-2"
                >
                    <Activity class="h-4 w-4" /> Статус перепустки
                </h3>
                {#if permit.is_closed}
                    <span
                        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-slate-100 text-slate-600 text-[10px] font-black uppercase border border-slate-200"
                    >
                        <Clock class="h-3 w-3" /> Закрита
                    </span>
                {:else}
                    <span
                        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 text-emerald-600 text-[10px] font-black uppercase border border-emerald-100"
                    >
                        <div
                            class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"
                        ></div>
                        Активна
                    </span>
                {/if}
            </div>
            <div class="p-5 space-y-4">
                {#if permit.verified_at}
                    <div
                        class="p-3 bg-blue-50/50 rounded-lg border border-blue-100/50"
                    >
                        <div
                            class="text-[10px] font-black uppercase tracking-widest text-blue-500 mb-1"
                        >
                            Валідація
                        </div>
                        <div
                            class="text-sm font-bold text-blue-900 flex items-center gap-2"
                        >
                            {permit.verifier
                                ? `${permit.verifier.first_name} ${permit.verifier.last_name}`
                                : "Оператор"}
                        </div>
                        <div class="text-[11px] text-blue-600/70 mt-0.5">
                            {formatDate(permit.verified_at)}
                        </div>
                    </div>
                {/if}

                <div class="space-y-3">
                    <div class="flex items-center justify-between text-sm">
                        <span
                            class="text-muted-foreground font-medium flex items-center gap-2"
                        >
                            <Activity class="h-4 w-4 text-indigo-400" /> Заїзд
                        </span>
                        <span class="font-bold"
                            >{permit.entry_time
                                ? formatDate(permit.entry_time)
                                : "—"}</span
                        >
                    </div>

                    {#if permit.is_closed && permit.exit_time}
                        <div class="flex items-center justify-between text-sm">
                            <span
                                class="text-muted-foreground font-medium flex items-center gap-2"
                            >
                                <Clock class="h-4 w-4 text-rose-400" /> Виїзд
                            </span>
                            <span class="font-bold"
                                >{formatDate(permit.exit_time)}</span
                            >
                        </div>
                    {/if}

                    {#if permit.total_sum !== undefined}
                        <div class="pt-3 border-t border-dashed">
                            <div class="flex items-center justify-between">
                                <span
                                    class="text-sm font-bold text-muted-foreground flex items-center gap-2"
                                >
                                    <Coins class="h-4 w-4 text-amber-500" /> До сплати
                                </span>
                                <span
                                    class="text-2xl font-black font-mono text-emerald-600"
                                    >₴{permit.total_sum.toFixed(2)}</span
                                >
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </div>

        <!-- Validation Checklist Card (if not verified) -->
        {#if !permit.verified_at && !permit.is_closed}
            <div
                class="bg-card border-2 border-amber-100 rounded-xl shadow-md overflow-hidden bg-amber-50/10 animate-in fade-in slide-in-from-right-4 duration-500"
            >
                <div
                    class="p-5 bg-amber-50/50 border-b border-amber-100 flex items-center gap-2"
                >
                    <CircleCheck class="h-4 w-4 text-amber-600" />
                    <h3
                        class="font-black text-[10px] uppercase tracking-widest text-amber-700"
                    >
                        Чек-лист валідації
                    </h3>
                </div>
                <div class="p-5 space-y-4">
                    {#each validationItems as item}
                        {@render validationItem(
                            item.label,
                            item.isValid,
                            item.targetId,
                        )}
                    {/each}
                    {#if validationItems.length === 0}
                        <p class="text-sm text-muted-foreground italic">
                            Оберіть митний режим, щоб побачити чек-лист.
                        </p>
                    {/if}
                </div>
            </div>
        {/if}
    {/if}

    <!-- Actions Card -->
    <div
        class="bg-card border rounded-xl shadow-lg p-5 space-y-3 bg-slate-50/50"
    >
        <Button
            onclick={() => handleSave()}
            disabled={loading}
            class="w-full gap-2.5 h-12 text-base font-bold shadow-md {data.isNew
                ? 'bg-indigo-600 hover:bg-indigo-700'
                : ''}"
        >
            <Save class="h-5 w-5 {loading ? 'animate-spin' : ''}" />
            {data.isNew ? "Зареєструвати заїзд" : "Зберегти зміни"}
        </Button>

        {#if !data.isNew && !permit.is_closed}
            {#if data.canValidate && !permit.verified_at}
                <Button
                    variant="outline"
                    onclick={handleValidatePermit}
                    disabled={loading}
                    class="w-full gap-2.5 h-12 text-base font-bold shadow-sm text-indigo-700 border-indigo-200 hover:bg-indigo-100 bg-white"
                >
                    <CircleCheck class="h-5 w-5" /> Валідувати
                </Button>
            {/if}

            <AlertDialog.Root>
                <AlertDialog.Trigger>
                    <Button
                        variant="secondary"
                        disabled={loading}
                        class="w-full gap-2.5 h-12 text-base font-bold border shadow-sm bg-white hover:bg-rose-50 hover:text-rose-600 hover:border-rose-200 transition-all"
                    >
                        <Send class="h-5 w-5" /> Завершити стоянку
                    </Button>
                </AlertDialog.Trigger>
                <AlertDialog.Content>
                    <AlertDialog.Header>
                        <AlertDialog.Title
                            >Підтвердження закриття</AlertDialog.Title
                        >
                        <AlertDialog.Description
                            >Ви впевнені, що хочете закрити цю перепустку?</AlertDialog.Description
                        >
                    </AlertDialog.Header>
                    <AlertDialog.Footer>
                        <AlertDialog.Cancel>Скасувати</AlertDialog.Cancel>
                        <AlertDialog.Action onclick={handleClosePermit}
                            >Підтвердити</AlertDialog.Action
                        >
                    </AlertDialog.Footer>
                </AlertDialog.Content>
            </AlertDialog.Root>
        {/if}
    </div>
</div>

{#snippet validationItem(label: string, isValid: boolean, targetId?: string)}
    <div class="flex items-center gap-3 py-0.5">
        {#if isValid}
            <div
                class="h-5 w-5 rounded-full bg-emerald-500 flex items-center justify-center shadow-sm shrink-0"
            >
                <CircleCheck class="h-3.5 w-3.5 text-white" />
            </div>
        {:else}
            <div
                class="h-5 w-5 rounded-full bg-amber-200 flex items-center justify-center shrink-0"
            >
                <div
                    class="h-1.5 w-1.5 rounded-full bg-amber-600 animate-pulse"
                ></div>
            </div>
        {/if}

        {#if targetId && !isValid}
            <button
                onclick={() =>
                    document.getElementById(targetId)?.scrollIntoView({
                        behavior: "smooth",
                        block: "center",
                    })}
                class="text-sm font-bold text-slate-500 hover:text-amber-700 hover:underline transition-colors text-left"
            >
                {label}
            </button>
        {:else}
            <span
                class="text-sm font-bold {isValid
                    ? 'text-emerald-700'
                    : 'text-slate-500'}"
            >
                {label}
            </span>
        {/if}
    </div>
{/snippet}
