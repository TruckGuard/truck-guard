<script lang="ts">
    import { CalendarIcon } from "@lucide/svelte";
    import { Button, buttonVariants } from "$lib/components/ui/button";
    import * as Popover from "$lib/components/ui/popover";
    import { Input } from "$lib/components/ui/input";
    import { Calendar } from "$lib/components/ui/calendar";
    import { type DateValue, getLocalTimeZone } from "@internationalized/date";
    import { formatDate, formatDateOnly } from "$lib/utils/date";
    import { cn } from "$lib/utils";

    let {
        selectedOperator,
        inputValue = $bindable(),
        dateValue = $bindable(),
        timeHours = $bindable(),
        timeMinutes = $bindable(),
        includeTime = $bindable(),
        onApply,
    } = $props<{
        selectedOperator: string;
        inputValue: string;
        dateValue: DateValue | undefined;
        timeHours: string;
        timeMinutes: string;
        includeTime: boolean;
        onApply: () => void;
    }>();

    let calendarOpen = $state(false);
</script>

<div class="space-y-2">
    <label for="filterValue" class="text-xs text-muted-foreground block"
        >Значення</label
    >

    {#if selectedOperator === "last_days"}
        <Input
            type="number"
            min="1"
            id="filterValue"
            placeholder="Вкажіть кількість..."
            class="h-9 text-sm"
            bind:value={inputValue}
            onkeydown={(e) => {
                if (e.key === "Enter") onApply();
            }}
        />
    {:else}
        <Popover.Root bind:open={calendarOpen}>
            <Popover.Trigger
                class={cn(
                    buttonVariants({
                        variant: "outline",
                        class: "w-full justify-start text-left font-normal",
                    }),
                    !dateValue && "text-muted-foreground",
                )}
            >
                <CalendarIcon class="mr-2 h-4 w-4" />
                {dateValue
                    ? includeTime
                        ? formatDate(
                              `${dateValue.year}-${dateValue.month.toString().padStart(2, "0")}-${dateValue.day.toString().padStart(2, "0")}T${timeHours.padStart(2, "0")}:${timeMinutes.padStart(2, "0")}`,
                          )
                        : formatDateOnly(dateValue.toDate(getLocalTimeZone()))
                    : "Оберіть дату..."}
            </Popover.Trigger>
            <Popover.Content class="w-auto p-0 flex flex-col" align="start">
                <Calendar type="single" bind:value={dateValue} />

                <div class="p-3 border-t border-border">
                    <label
                        class="text-sm font-medium flex items-center gap-2 cursor-pointer"
                    >
                        <input
                            type="checkbox"
                            bind:checked={includeTime}
                            class="w-4 h-4 rounded border-gray-300 text-primary focus:ring-primary"
                        />
                        Враховувати час
                    </label>
                </div>

                {#if includeTime}
                    <div
                        class="p-3 border-t border-border flex items-center justify-between gap-2 bg-muted/30"
                    >
                        <div class="text-sm font-medium">Час:</div>
                        <div class="flex items-center gap-1">
                            <Input
                                type="number"
                                min="0"
                                max="23"
                                bind:value={timeHours}
                                class="w-16 h-8 text-center"
                                placeholder="ГГ"
                            />
                            <span>:</span>
                            <Input
                                type="number"
                                min="0"
                                max="59"
                                bind:value={timeMinutes}
                                class="w-16 h-8 text-center"
                                placeholder="ХХ"
                            />
                        </div>
                    </div>
                {/if}

                <div class="p-3 border-t border-border flex justify-end">
                    <Button
                        size="sm"
                        class="w-full"
                        onclick={() => {
                            calendarOpen = false;
                        }}
                    >
                        Застосувати
                    </Button>
                </div>
            </Popover.Content>
        </Popover.Root>
    {/if}
</div>
