<script lang="ts">
    import { ShieldAlert } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import { enhance, applyAction } from "$app/forms";
    import { toast } from "svelte-sonner";

    let { sessions = [] } = $props();
    let sessionCount = $derived(sessions.length);
</script>

<div
    class="rounded-lg border border-destructive/20 bg-destructive/5 text-card-foreground shadow-sm p-6"
>
    <div class="flex items-center justify-between mb-6">
        <h2 class="text-lg font-semibold flex items-center gap-2">
            <ShieldAlert class="h-5 w-5 text-destructive" /> Безпека
        </h2>
        {#if sessionCount > 1}
            <span
                class="text-xs bg-destructive/10 text-destructive px-2 py-1 rounded-full font-medium border border-destructive/20"
            >
                {sessionCount} активних пристроїв
            </span>
        {/if}
    </div>

    <div class="space-y-4">
        <p class="text-sm text-muted-foreground">
            Якщо ви помітили підозрілу активність або хочете завершити всі
            активні сеанси на інших пристроях, скористайтеся кнопкою нижче.
        </p>

        <form
            method="POST"
            action="?/revokeAllSessions"
            use:enhance={() => {
                const loadingToast = toast.loading("Вихід з усіх пристроїв...");
                return async ({ result }) => {
                    toast.dismiss(loadingToast);
                    if (
                        result.type === "success" ||
                        result.type === "redirect"
                    ) {
                        toast.success(
                            "Всі сесії завершено. Вас буде перенаправлено...",
                        );
                        await applyAction(result);
                    } else {
                        toast.error("Не вдалося завершити сесії");
                    }
                };
            }}
        >
            <Button
                type="submit"
                variant="destructive"
                class="w-full md:w-auto"
            >
                Вийти з усіх пристроїв
            </Button>
        </form>
    </div>
</div>
