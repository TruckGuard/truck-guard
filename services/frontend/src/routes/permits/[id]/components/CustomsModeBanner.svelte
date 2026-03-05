<script lang="ts">
    import { ShieldCheck } from "@lucide/svelte";
    import * as Select from "$lib/components/ui/select";
    import type { CustomsMode } from "$lib/types/data";
    import type { Permit } from "$lib/types/permits";

    let {
        permit = $bindable(),
        customsModes,
        verified = false,
    } = $props<{
        permit: Partial<Permit>;
        customsModes: CustomsMode[];
        verified?: boolean;
    }>();

    const selectedMode = $derived(
        customsModes.find(
            (m: CustomsMode) => m.code === permit.customs_mode_code,
        ),
    );

    const hasMode = $derived(!!permit.customs_mode_code);
    const show = $derived(!hasMode && !verified);
</script>

<div
    class="overflow-hidden py-4"
    style="
        max-height: {show ? '200px' : '0px'};
        opacity: {show ? 1 : 0};
        margin-bottom: {show ? '0' : '-2rem'};
        transition: max-height 0.5s ease, opacity 0.4s ease, margin-bottom 0.5s ease;
        pointer-events: {show ? 'auto' : 'none'};
    "
>
    <div
        class="rounded-xl border-2 border-amber-200 dark:border-amber-700/50 bg-amber-50 dark:bg-amber-950/30 shadow-md overflow-hidden"
    >
        <div class="flex items-center gap-4 p-4 sm:p-5">
            <div
                class="shrink-0 h-12 w-12 rounded-xl flex items-center justify-center bg-amber-100 dark:bg-amber-900/50 text-amber-600 dark:text-amber-400"
            >
                <ShieldCheck class="h-6 w-6" />
            </div>

            <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2 flex-wrap">
                    <span
                        class="text-[10px] font-black uppercase tracking-widest text-amber-500"
                    >
                        Спочатку вкажіть
                    </span>
                    <span
                        class="text-[10px] bg-amber-200 dark:bg-amber-800 text-amber-800 dark:text-amber-300 font-bold px-2 py-0.5 rounded-full"
                    >
                        Обов'язково
                    </span>
                </div>
                <p class="text-xs text-muted-foreground mt-0.5">
                    Вибраний режим визначить, які поля потрібно заповнити.
                </p>
            </div>

            <div class="shrink-0 w-64 sm:w-80" id="mode-input">
                <Select.Root
                    type="single"
                    bind:value={permit.customs_mode_code}
                    name="customs_mode_top"
                >
                    <Select.Trigger
                        class="w-full h-11 font-bold text-sm border-amber-400 dark:border-amber-600 ring-2 ring-amber-200 dark:ring-amber-800 bg-white dark:bg-transparent"
                    >
                        <span
                            class="text-amber-700 dark:text-amber-400 font-bold"
                            >Виберіть режим митниці...</span
                        >
                    </Select.Trigger>
                    <Select.Content>
                        {#each customsModes as cm}
                            <Select.Item value={cm.code}>
                                <span
                                    class="font-mono text-xs bg-muted px-1.5 py-0.5 rounded mr-2"
                                    >{cm.code}</span
                                >
                                {cm.name}
                            </Select.Item>
                        {/each}
                    </Select.Content>
                </Select.Root>
            </div>
        </div>

        <div class="h-1 bg-amber-200 overflow-hidden relative">
            <div
                class="h-full bg-amber-400 w-2/5 absolute animate-shuttle"
            ></div>
        </div>
    </div>
</div>

<style>
    @keyframes shuttle {
        0% {
            left: -40%;
        }
        100% {
            left: 140%;
        }
    }
    .animate-shuttle {
        animation: shuttle 1.8s ease-in-out infinite;
    }
</style>
