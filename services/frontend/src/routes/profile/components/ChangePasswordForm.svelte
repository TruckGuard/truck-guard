<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { KeyRound, Loader2, CheckCircle2, AlertCircle } from "@lucide/svelte";
    import { enhance } from "$app/forms";
    import { toast } from "svelte-sonner";

    let loading = $state(false);
    let success = $state(false);
    let error = $state<string | null>(null);

    let currentPassword = $state("");
    let newPassword = $state("");
    let confirmPassword = $state("");

    function handleSubmit() {
        loading = true;
        error = null;
        success = false;
        
        return async ({ result }: { result: any }) => {
            loading = false;
            if (result.type === "success") {
                success = true;
                currentPassword = "";
                newPassword = "";
                confirmPassword = "";
                toast.success("Пароль успішно змінено");
            } else if (result.type === "failure") {
                error = result.data?.error || "Помилка при зміні пароля";
                if (error) toast.error(error);
            } else if (result.type === "error") {
                error = result.error?.message || "Непередбачена помилка";
                if (error) toast.error(error);
            }
        };
    }
</script>

<Card.Root>
    <Card.Header>
        <Card.Title class="flex items-center gap-2">
            <KeyRound class="size-5 text-primary" />
            Зміна пароля
        </Card.Title>
        <Card.Description>
            Введіть ваш поточний пароль та новий пароль для зміни.
        </Card.Description>
    </Card.Header>
    <Card.Content>
        <form 
            method="POST" 
            action="?/changePassword" 
            use:enhance={handleSubmit}
            class="space-y-4"
        >
            {#if error}
                <div class="bg-destructive/10 text-destructive text-sm p-3 rounded-md flex items-center gap-2">
                    <AlertCircle class="size-4" />
                    {error}
                </div>
            {/if}

            {#if success}
                <div class="bg-green-500/10 text-green-600 text-sm p-3 rounded-md flex items-center gap-2">
                    <CheckCircle2 class="size-4" />
                    Пароль успішно змінено
                </div>
            {/if}

            <div class="space-y-2">
                <Label for="current_password">Поточний пароль</Label>
                <Input 
                    id="current_password" 
                    name="current_password" 
                    type="password" 
                    required 
                    bind:value={currentPassword}
                />
            </div>

            <div class="space-y-2">
                <Label for="new_password">Новий пароль</Label>
                <Input 
                    id="new_password" 
                    name="new_password" 
                    type="password" 
                    required 
                    bind:value={newPassword}
                />
            </div>

            <div class="space-y-2">
                <Label for="confirm_password">Підтвердіть новий пароль</Label>
                <Input 
                    id="confirm_password" 
                    name="confirm_password" 
                    type="password" 
                    required 
                    bind:value={confirmPassword}
                />
            </div>

            <Button type="submit" class="w-full" disabled={loading || !newPassword || newPassword !== confirmPassword}>
                {#if loading}
                    <Loader2 class="mr-2 h-4 w-4 animate-spin" />
                    Зміна...
                {:else}
                    Змінити пароль
                {/if}
            </Button>
            
            {#if newPassword && confirmPassword && newPassword !== confirmPassword}
                <p class="text-[10px] text-destructive italic">Паролі не збігаються</p>
            {/if}
        </form>
    </Card.Content>
</Card.Root>
