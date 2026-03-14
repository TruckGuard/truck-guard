<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { RefreshCcw, Plus } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button";
  import type { Permit } from "$lib/types/permits";
  import type { CustomsPost, VehicleType } from "$lib/types/data";
  import type { ApiResponse } from "$lib/types/events";

  import DataTable from "./data-table.svelte";
  import { columns } from "./columns";
  import * as AlertDialog from "$lib/components/ui/alert-dialog";
  import { enhance } from "$app/forms";
  import { toast } from "svelte-sonner";

  let { data } = $props<{
    data: {
      activePermits: ApiResponse<Permit>;
      customsPosts: CustomsPost[];
      vehicleTypes: VehicleType[];
      paymentTypes: any[];
      customsModes: any[];
      users: any[];
      hasAllPermitsAccess: boolean;
    };
  }>();

  let loading = $state(false);

  function refresh() {
    loading = true;
    const query = new URLSearchParams(page.url.searchParams.toString());
    query.set("page", "1");
    goto(`?${query.toString()}`, { invalidateAll: true }).then(
      () => (loading = false),
    );
  }

  function startCreatePermit(event?: any) {
    if (event) {
        pendingEvent = event;
        showCreateConfirm = true;
    } else {
        pendingEvent = null;
        showCreateConfirm = true;
    }
  }

  let pendingEvent = $state<any>(null);
  let showCreateConfirm = $state(false);
  let creating = $state(false);
</script>

<div
  class="flex flex-col h-full overflow-hidden space-y-6"
>
  <div class="flex items-center justify-between shrink-0">
    <div>
      <h1 class="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-2">
        Панель оператора перепусток
      </h1>
      <p class="text-muted-foreground text-sm mb-0">
        Моніторинг живих подій та керування активними перепустками в зоні
      </p>
    </div>
    <div class="flex items-center gap-3">
      <Button variant="outline" size="sm" class="h-10 px-4 shadow-sm" onclick={refresh} disabled={loading}>
        <RefreshCcw class="mr-2 h-4 w-4 {loading ? 'animate-spin' : ''}" />
        Оновити дані
      </Button>
      <Button size="sm" class="h-10 px-4 shadow-sm" onclick={() => startCreatePermit()}>
        <Plus class="mr-2 h-4 w-4" /> Створити вручну
      </Button>
    </div>
  </div>

  <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
    <DataTable
      data={data.activePermits?.data || []}
      {columns}
      customsPosts={data.customsPosts || []}
      vehicleTypes={data.vehicleTypes || []}
      paymentTypes={data.paymentTypes || []}
      customsModes={data.customsModes || []}
      users={data.users || []}
      hasAllPermitsAccess={data.hasAllPermitsAccess}
      metadata={data.activePermits?.metadata}
    />
  </div>
</div>

<AlertDialog.Root bind:open={showCreateConfirm}>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>Створити нову перепустку?</AlertDialog.Title>
            <AlertDialog.Description>
                Система створить нову порожню перепустку та призначить їй унікальний код. 
                Ви будете автоматично перенаправлені на сторінку редагування для заповнення даних.
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel>Скасувати</AlertDialog.Cancel>
            <form 
                action="?/create" 
                method="POST" 
                use:enhance={() => {
                    creating = true;
                    return async ({ result, update }) => {
                        creating = false;
                        if (result.type === 'error' || result.type === 'failure') {
                            toast.error("Не вдалося створити перепустку: " + (result as any).data?.error);
                        }
                        await update();
                    };
                }}
            >
                {#if pendingEvent}
                    {#if pendingEvent._type === 'plate'}
                        <input type="hidden" name="camera_event_id" value={pendingEvent.ID} />
                    {:else if pendingEvent._type === 'weight'}
                        <input type="hidden" name="scale_event_id" value={pendingEvent.ID} />
                    {/if}
                {/if}
                <Button type="submit" disabled={creating}>
                    {creating ? "Створення..." : "Так, створити"}
                </Button>
            </form>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
