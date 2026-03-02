<script lang="ts">
    import { Clock } from "@lucide/svelte";

    let { auditLogs } = $props<{ auditLogs: any[] }>();

    function formatDate(dateStr?: string) {
        if (!dateStr) return "-";
        return new Date(dateStr).toLocaleString("uk-UA");
    }
</script>

<div class="space-y-4">
    <div class="flex items-center justify-between">
        <h3 class="font-semibold text-lg flex items-center gap-2">
            <Clock class="h-5 w-5 text-slate-500" />
            Логи аудиту
        </h3>
    </div>

    {#if auditLogs.length === 0}
        <div
            class="bg-card border rounded-xl p-8 text-center text-muted-foreground shadow-sm"
        >
            Історія змін порожня
        </div>
    {:else}
        <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
            <ul class="divide-y text-sm">
                {#each auditLogs as audit}
                    <li
                        class="p-6 hover:bg-slate-50 dark:hover:bg-slate-900/50 flex flex-col sm:flex-row sm:items-start gap-4 sm:gap-8"
                    >
                        <div
                            class="w-40 shrink-0 text-muted-foreground whitespace-nowrap"
                        >
                            {formatDate(audit.CreatedAt)}
                        </div>
                        <div class="flex-1 space-y-2">
                            <div
                                class="font-bold flex items-center gap-3 text-lg"
                            >
                                <span
                                    class="capitalize px-3 py-1 rounded-md text-[10px] tracking-widest font-black uppercase
                  {audit.action === 'create'
                                        ? 'bg-blue-100 text-blue-700'
                                        : audit.action === 'validate'
                                          ? 'bg-indigo-100 text-indigo-700'
                                          : audit.action === 'close'
                                            ? 'bg-slate-100 text-slate-700'
                                            : 'bg-stone-100 text-stone-700'}"
                                >
                                    {audit.action}
                                </span>
                                {#if audit.user}
                                    <span class="text-foreground font-extrabold"
                                        >{audit.user.first_name}
                                        {audit.user.last_name}</span
                                    >
                                {:else if audit.user_id}
                                    <span class="text-foreground font-extrabold"
                                        >ID Користувача: {audit.user_id}</span
                                    >
                                {:else}
                                    <span class="text-foreground font-extrabold"
                                        >Система (Автоматично)</span
                                    >
                                {/if}
                            </div>
                            {#if audit.comment}
                                <p
                                    class="text-muted-foreground text-sm font-medium"
                                >
                                    {audit.comment}
                                </p>
                            {/if}
                            {#if audit.changes && Object.keys(audit.changes).length > 0}
                                <div
                                    class="mt-4 text-sm border rounded-lg overflow-hidden bg-background shadow-inner"
                                >
                                    <table class="w-full text-left">
                                        <thead
                                            class="bg-muted/70 text-muted-foreground"
                                        >
                                            <tr>
                                                <th
                                                    class="px-4 py-2.5 font-bold border-b w-1/4 uppercase tracking-tighter text-[11px]"
                                                    >Поле</th
                                                >
                                                <th
                                                    class="px-4 py-2.5 font-bold border-b w-3/4 uppercase tracking-tighter text-[11px]"
                                                    >Зміни</th
                                                >
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y">
                                            {#each Object.entries(audit.changes) as [key, val]}
                                                <tr>
                                                    <td
                                                        class="px-4 py-3 font-mono text-muted-foreground border-r font-bold bg-slate-50/30"
                                                        >{key}</td
                                                    >
                                                    <td
                                                        class="px-4 py-3 font-mono text-sm"
                                                    >
                                                        <span
                                                            class="font-black inline-block text-slate-700 dark:text-slate-300"
                                                        >
                                                            {typeof val ===
                                                            "object"
                                                                ? JSON.stringify(
                                                                      val,
                                                                  )
                                                                : val}
                                                        </span>
                                                    </td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            {/if}
                        </div>
                    </li>
                {/each}
            </ul>
        </div>
    {/if}
</div>
