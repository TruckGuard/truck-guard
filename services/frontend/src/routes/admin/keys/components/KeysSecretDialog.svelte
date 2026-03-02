<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { Button } from "$lib/components/ui/button";
    import Copy from "@lucide/svelte/icons/copy";
    import { toast } from "svelte-sonner";
    import FormDialog from "$lib/components/common/FormDialog.svelte";

    let { open = $bindable(false), secret } = $props<{
        open: boolean;
        secret: string;
    }>();

    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
        toast.success("Скопійовано в буфер обміну");
    }
</script>

<FormDialog
    bind:open
    title="API Ключ Створено"
    description="Збережіть цей ключ зараз. Ви не зможете побачити його знову."
>
    <div class="flex items-center space-x-2 mt-4">
        <Input readonly value={secret} class="font-mono text-center bg-muted" />
        <Button
            variant="outline"
            size="icon"
            onclick={() => copyToClipboard(secret)}
        >
            <Copy class="h-4 w-4" />
        </Button>
    </div>
    <div class="flex justify-end pt-6">
        <Button onclick={() => (open = false)}>Закрити</Button>
    </div>
</FormDialog>
