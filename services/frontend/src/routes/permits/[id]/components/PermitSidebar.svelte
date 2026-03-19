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
    import { formatDate } from "$lib/utils/date";

    let {
        permit,
        data,
        loading,
        handleSave,
        handleValidatePermit,
        handleClosePermit,
        handleRestore,
        handleVoid,
        handleDelete,
        validationItems,
    } = $props<{
        permit: any;
        data: any;
        loading: boolean;
        handleSave: (silent?: boolean) => Promise<boolean>;
        handleValidatePermit: () => Promise<void>;
        handleClosePermit: () => Promise<void>;
        handleRestore: () => Promise<void>;
        handleVoid: () => Promise<void>;
        handleDelete: () => Promise<void>;
        validationItems: {
            label: string;
            isValid: boolean;
            targetId?: string;
        }[];
    }>();

    let estimatedSum = $derived.by(() => {
        if (permit.is_closed || !permit.entry_time) return null;

        const entry = new Date(permit.entry_time);
        const now = new Date();
        const durationMs = now.getTime() - entry.getTime();
        const hours = durationMs / (1000 * 60 * 60);

        let days = Math.floor(hours / 24);
        if (hours > days * 24 || durationMs <= 0) {
            days++;
        }
        if (days < 1) days = 1;

        const entryFee = permit.entry_fee || (permit.vehicle_type?.entry_price || 0);
        const dailyPrice = permit.daily_fee || (permit.vehicle_type?.daily_price || 0);

        return entryFee + days * dailyPrice;
    });
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
                        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 text-[10px] font-black uppercase border border-slate-200 dark:border-slate-700"
                    >
                        <Clock class="h-3 w-3" /> Закрита
                    </span>
                {:else}
                    <span
                        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400 text-[10px] font-black uppercase border border-emerald-100 dark:border-emerald-800"
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
                        class="p-3 bg-blue-50/50 dark:bg-blue-950/30 rounded-lg border border-blue-100/50 dark:border-blue-800/50"
                    >
                        <div
                            class="text-[10px] font-black uppercase tracking-widest text-blue-500 dark:text-blue-400 mb-1"
                        >
                            Валідація
                        </div>
                        <div
                            class="text-sm font-bold text-blue-900 dark:text-blue-100 flex items-center gap-2"
                        >
                            {permit.verifier
                                ? `${permit.verifier.first_name} ${permit.verifier.last_name}`
                                : "Оператор"}
                        </div>
                        <div
                            class="text-[11px] text-blue-600/70 dark:text-blue-400/70 mt-0.5"
                        >
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

                    {#if permit.creator}
                        <div class="flex items-center justify-between text-[11px] pt-1 border-t border-muted/20 mt-1">
                            <span class="text-muted-foreground italic">Реєстрація:</span>
                            <span class="font-medium">
                                {permit.creator.first_name} {permit.creator.last_name}
                            </span>
                        </div>
                    {/if}

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
                        {#if permit.closed_by}
                            <div class="flex items-center justify-between text-xs pt-1">
                                <span class="text-muted-foreground italic">Закрив:</span>
                                <span class="font-medium">{permit.closed_by.first_name} {permit.closed_by.last_name}</span>
                            </div>
                        {/if}
                    {/if}

                    {#if permit.is_void}
                        <div class="p-3 bg-red-50 dark:bg-red-950/30 rounded-lg border border-red-100 dark:border-red-800 flex items-center gap-3">
                            <div class="h-8 w-8 rounded-full bg-red-100 dark:bg-red-900/50 flex items-center justify-center shrink-0">
                                <Activity class="h-4 w-4 text-red-600" />
                            </div>
                            <div>
                                <div class="text-[10px] font-black uppercase text-red-600 dark:text-red-400">Анульовано</div>
                                <div class="text-[11px] text-red-700/70 dark:text-red-400/70">Ця перепустка більше не є дійсною</div>
                                {#if permit.voided_by}
                                    <div class="text-[10px] mt-1 italic text-red-600/80 dark:text-red-400/80">
                                        Анулював: {permit.voided_by.first_name} {permit.voided_by.last_name}
                                    </div>
                                {/if}
                            </div>
                        </div>
                    {/if}


                    {#if permit.is_closed && permit.total_sum !== undefined}
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
                    {:else if !permit.is_closed && estimatedSum !== null}
                        <div class="pt-3 border-t border-dashed animate-in fade-in slide-in-from-top-2 duration-300">
                            <div class="flex items-center justify-between">
                                <div class="flex flex-col">
                                    <span
                                        class="text-[10px] font-black uppercase tracking-widest text-amber-600 dark:text-amber-400 flex items-center gap-2"
                                    >
                                        <Coins class="h-3 w-3" /> Приблизна оплата
                                    </span>
                                    <span class="text-[9px] text-muted-foreground font-medium italic mt-0.5">
                                        (без урахування знижок)
                                    </span>
                                </div>
                                <span
                                    class="text-2xl font-black font-mono text-amber-600 dark:text-amber-400"
                                    >₴{estimatedSum.toFixed(2)}</span
                                >
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </div>

        <!-- Validation Checklist Card (if not verified and NOT void) -->
        {#if !permit.verified_at && !permit.is_closed && !permit.is_void}
            <div
                class="bg-card border-2 border-amber-100 dark:border-amber-800/50 rounded-xl shadow-md overflow-hidden bg-amber-50/10 dark:bg-amber-950/10 animate-in fade-in slide-in-from-right-4 duration-500"
            >
                <div
                    class="p-5 bg-amber-50/50 dark:bg-amber-950/30 border-b border-amber-100 dark:border-amber-800/50 flex items-center gap-2"
                >
                    <CircleCheck
                        class="h-4 w-4 text-amber-600 dark:text-amber-400"
                    />
                    <h3
                        class="font-black text-[10px] uppercase tracking-widest text-amber-700 dark:text-amber-400"
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
        class="bg-card border rounded-xl shadow-lg p-5 space-y-3 bg-slate-50/50 dark:bg-slate-900/30"
    >
        {#if data.canUpdate}
            <Button
                onclick={() => handleSave()}
                disabled={loading}
                class="w-full gap-2.5 h-12 text-base font-bold shadow-md {data.isNew
                    ? 'bg-indigo-600 hover:bg-indigo-700'
                    : ''}"
            >
                <Save class="h-5 w-5 {loading ? 'animate-spin' : ''}" />
                {!permit.verified_at ? "Зареєструвати заїзд" : "Зберегти зміни"}
            </Button>
        {/if}

        {#if !permit.is_closed && !permit.is_void}
            {#if data.canValidate && !permit.verified_at}
                <Button
                    variant="outline"
                    onclick={handleValidatePermit}
                    disabled={loading}
                    class="w-full gap-2.5 h-12 text-base font-bold shadow-sm text-indigo-700 dark:text-indigo-300 border-indigo-200 dark:border-indigo-700 hover:bg-indigo-100 dark:hover:bg-indigo-900/50 bg-white dark:bg-transparent"
                >
                    <CircleCheck class="h-5 w-5" /> Валідувати
                </Button>
            {/if}

            {#if !permit.is_void}
                <div class="grid grid-cols-2 gap-2">
                    {#if data.canDelete}
                        <AlertDialog.Root>
                            <AlertDialog.Trigger>
                                <Button
                                    variant="outline"
                                    disabled={loading}
                                    class="w-full gap-2 h-12 text-sm font-bold border shadow-sm bg-white dark:bg-transparent hover:bg-red-50 dark:hover:bg-red-950/30 hover:text-red-600 dark:hover:text-red-400 hover:border-red-200 dark:hover:border-red-800 transition-all"
                                >
                                    <Activity class="h-4 w-4" /> Анулювати
                                </Button>
                            </AlertDialog.Trigger>
                            <AlertDialog.Content>
                                <AlertDialog.Header>
                                    <AlertDialog.Title>Підтвердження анулювання</AlertDialog.Title>
                                    <AlertDialog.Description>
                                        Ви впевнені, що хочете анулювати цю перепустку? Це зробить її недійсною.
                                    </AlertDialog.Description>
                                </AlertDialog.Header>
                                <AlertDialog.Footer>
                                    <AlertDialog.Cancel>Скасувати</AlertDialog.Cancel>
                                    <AlertDialog.Action onclick={handleVoid} class="bg-red-600 hover:bg-red-700">Анулювати</AlertDialog.Action>
                                </AlertDialog.Footer>
                            </AlertDialog.Content>
                        </AlertDialog.Root>
                    {/if}

                    {#if data.canUpdate}
                        <AlertDialog.Root>
                            <AlertDialog.Trigger>
                                <Button
                                    variant="secondary"
                                    disabled={loading}
                                    class="w-full gap-2 h-12 text-sm font-bold border shadow-sm bg-white dark:bg-transparent hover:bg-slate-50 dark:hover:bg-slate-950/30 transition-all"
                                >
                                    <Send class="h-4 w-4" /> Закрити
                                </Button>
                            </AlertDialog.Trigger>
                            <AlertDialog.Content>
                                <AlertDialog.Header>
                                    <AlertDialog.Title>Підтвердження закриття</AlertDialog.Title>
                                    <AlertDialog.Description>Ви впевнені, що хочете закрити цю перепустку?</AlertDialog.Description>
                                </AlertDialog.Header>
                                <AlertDialog.Footer>
                                    <AlertDialog.Cancel>Скасувати</AlertDialog.Cancel>
                                    <AlertDialog.Action onclick={handleClosePermit}>Підтвердити</AlertDialog.Action>
                                </AlertDialog.Footer>
                            </AlertDialog.Content>
                        </AlertDialog.Root>
                    {/if}
                </div>
            {/if}
        {/if}

        {#if permit.is_void}
            {#if data.canUpdate}
                <Button
                    variant="outline"
                    onclick={handleRestore}
                    disabled={loading}
                    class="w-full gap-2.5 h-12 text-base font-bold border-emerald-200 dark:border-emerald-800 text-emerald-700 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-950/30"
                >
                    <Activity class="h-5 w-5" /> Відновити перепусту
                </Button>
            {/if}

            {#if data.canDelete}
                <AlertDialog.Root>
                    <AlertDialog.Trigger>
                        <Button
                            variant="destructive"
                            disabled={loading}
                            class="w-full gap-2.5 h-12 text-base font-bold shadow-md"
                        >
                            <Activity class="h-5 w-5" /> Видалити остаточно
                        </Button>
                    </AlertDialog.Trigger>
                    <AlertDialog.Content>
                        <AlertDialog.Header>
                            <AlertDialog.Title>Видалення перепустки</AlertDialog.Title>
                            <AlertDialog.Description>
                                Ви впевнені, що хочете остаточно видалити цю перепустку зі сховища? Це скасує всі пов'язані записи.
                            </AlertDialog.Description>
                        </AlertDialog.Header>
                        <AlertDialog.Footer>
                            <AlertDialog.Cancel>Скасувати</AlertDialog.Cancel>
                            <AlertDialog.Action onclick={handleDelete} class="bg-red-600 hover:bg-red-700">Видалити</AlertDialog.Action>
                        </AlertDialog.Footer>
                    </AlertDialog.Content>
                </AlertDialog.Root>
            {/if}
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
                class="text-sm font-bold text-slate-500 dark:text-slate-400 hover:text-amber-700 dark:hover:text-amber-400 hover:underline transition-colors text-left"
            >
                {label}
            </button>
        {:else}
            <span
                class="text-sm font-bold {isValid
                    ? 'text-emerald-700 dark:text-emerald-400'
                    : 'text-slate-500 dark:text-slate-400'}"
            >
                {label}
            </span>
        {/if}
    </div>
{/snippet}
