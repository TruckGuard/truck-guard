<script lang="ts">
    import { enhance } from "$app/forms";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Save } from "@lucide/svelte";
    import type { Company } from "$lib/types/data";
    import { toast } from "svelte-sonner";
    import { mapErrorToFriendlyMessage } from "$lib/utils/error-handler";

    let {
        open = $bindable(false),
        initialQuery = "",
        onsuccess,
    } = $props<{
        open?: boolean;
        initialQuery?: string;
        onsuccess?: (company: Company) => void;
    }>();

    let name = $state("");
    let edrpou = $state("");
    let isSubmitting = $state(false);

    // When the dialog opens, initialize fields based on query
    $effect(() => {
        if (open) {
            if (/^\d+$/.test(initialQuery.trim())) {
                edrpou = initialQuery.trim();
                name = "";
            } else {
                name = initialQuery.trim();
                edrpou = "";
            }
        }
    });
</script>

<Dialog.Root bind:open>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Створити нову компанію</Dialog.Title>
            <Dialog.Description>
                Ця компанія буде відразу додана до довідника платників.
            </Dialog.Description>
        </Dialog.Header>

        <form
            method="POST"
            action="?/createCompany"
            use:enhance={() => {
                isSubmitting = true;
                return async ({ result }) => {
                    isSubmitting = false;
                    if (result.type === "success") {
                        toast.success("Компанію успішно створено!");
                        open = false;
                        // result.data should contain the created company from our server action
                        if (onsuccess && (result as any).data?.data) {
                            onsuccess((result as any).data.data as Company);
                        }
                    } else {
                        const errorMsg =
                            (result as any).data?.error?.message ||
                            "Не вдалося створити компанію";
                        toast.error(mapErrorToFriendlyMessage(errorMsg));
                    }
                };
            }}
            class="space-y-4 py-4"
        >
            <div class="space-y-2">
                <Label for="company-name"
                    >Назва компанії <span class="text-rose-500">*</span></Label
                >
                <Input
                    id="company-name"
                    name="name"
                    bind:value={name}
                    placeholder="Напр. ТОВ «Логістика Плюс»"
                    required
                    autocomplete="off"
                />
            </div>

            <div class="space-y-2">
                <Label for="company-edrpou"
                    >ЄДРПОУ / ІПН <span class="text-rose-500">*</span></Label
                >
                <Input
                    id="company-edrpou"
                    name="edrpou"
                    bind:value={edrpou}
                    placeholder="Код платника (8-10 цифр)"
                    required
                    autocomplete="off"
                />
            </div>

            <Dialog.Footer class="pt-4">
                <Button
                    variant="outline"
                    type="button"
                    onclick={() => (open = false)}
                    disabled={isSubmitting}
                >
                    Скасувати
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                    {#if isSubmitting}
                        Збереження...
                    {:else}
                        <Save class="w-4 h-4 mr-2" /> Зберегти
                    {/if}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
