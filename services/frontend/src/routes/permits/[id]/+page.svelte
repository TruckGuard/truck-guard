<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import {
    ChevronLeft,
    Save,
    FileSearch,
    CheckCircle,
    XOctagon,
    Truck,
    Clock,
    Coins,
    UserCog,
    Send,
    Activity,
  } from "@lucide/svelte";
  import * as Tabs from "$lib/components/ui/tabs";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Textarea } from "$lib/components/ui/textarea";
  import * as Select from "$lib/components/ui/select";
  import type { Permit, PermitCustomsData } from "$lib/types/permits";
  import * as AlertDialog from "$lib/components/ui/alert-dialog";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  let { data } = $props<{
    data: {
      permit: Partial<Permit>;
      isNew: boolean;
      userRole: any;
      canValidate: boolean;
      canManageCustoms: boolean;
      vehicleTypes: any[];
      customsModes: any[];
      companies: any[];
      paymentTypes: any[];
    };
  }>();

  let loading = $state(false);
  let fetchingCustoms = $state(false);

  // Create a mutable copy of the permit to bind to the form
  let permit = $state({ ...data.permit });
  if (!permit.customs_data) {
    permit.customs_data = {
      goods: "",
      declarant: "",
      vmd_number: "",
      sender: "",
      receiver: "",
    };
  }

  // Payer state (for simplicity, handled via array)
  let primaryCompanyId = $state(
    permit.payers?.[0]?.company_id?.toString() || "",
  );

  let activeTab = $state("events");
  let loadingAudits = $state(false);
  let auditLogs: any[] = $state([]);

  // Validation state
  let showValidationErrors = $state(false);
  const isValidPlate = $derived(!!permit.plate_front && !!permit.plate_back);
  const isValidWeight = $derived(!!(permit.total_weight || 0));
  const isValidPD = $derived(!!permit.declaration_number);
  const isValidVehicleType = $derived(!!permit.vehicle_type_id);
  const isValidCustomsMode = $derived(!!permit.customs_mode_code);
  const isValidCustomsData = $derived(
    !!permit.customs_data &&
      !!permit.customs_data.goods &&
      !!permit.customs_data.vmd_number,
  );
  const isAllValid = $derived(
    isValidPlate &&
      isValidWeight &&
      isValidPD &&
      isValidVehicleType &&
      isValidCustomsMode &&
      isValidCustomsData,
  );

  async function handleSave(silent = false) {
    loading = true;
    try {
      // Build optimal payload by checking differences vs original `data.permit`
      const payload: any = {};

      const isDiff = (a: any, b: any) => {
        if (!a && !b) return false;
        return a !== b;
      };

      if (isDiff(permit.plate_front, data.permit.plate_front))
        payload.plate_front = permit.plate_front;
      if (isDiff(permit.plate_back, data.permit.plate_back))
        payload.plate_back = permit.plate_back;
      if (
        Number(permit.total_weight || 0) !==
        Number(data.permit.total_weight || 0)
      )
        payload.total_weight = Number(permit.total_weight) || 0;
      if (
        Number(permit.vehicle_type_id || 0) !==
        Number(data.permit.vehicle_type_id || 0)
      )
        payload.vehicle_type_id = Number(permit.vehicle_type_id) || undefined;
      if (isDiff(permit.customs_mode_code, data.permit.customs_mode_code))
        payload.customs_mode_code = permit.customs_mode_code || undefined;
      if (isDiff(permit.declaration_number, data.permit.declaration_number))
        payload.declaration_number = permit.declaration_number || undefined;
      if (
        Number(permit.payment_type_id || 0) !==
        Number(data.permit.payment_type_id || 0)
      )
        payload.payment_type_id = Number(permit.payment_type_id) || undefined;
      if (isDiff(permit.notes, data.permit.notes)) payload.notes = permit.notes;

      // Customs Data is nested, simplest is to just stringify and compare
      const isCustomsEmpty = (cd: any) =>
        !cd ||
        (!cd.vmd_number &&
          !cd.declarant &&
          !cd.goods &&
          !cd.sender &&
          !cd.receiver);
      const cdChanged =
        JSON.stringify(permit.customs_data || {}) !==
        JSON.stringify(data.permit.customs_data || {});

      if (
        cdChanged &&
        !(
          isCustomsEmpty(permit.customs_data) &&
          isCustomsEmpty(data.permit.customs_data)
        )
      ) {
        payload.customs_data = permit.customs_data;
      }

      // In real scenario we'd properly manage payers
      const newCompanyId = primaryCompanyId
        ? Number(primaryCompanyId)
        : undefined;
      const oldCompanyId = data.permit.payers?.[0]?.company_id;
      if (newCompanyId !== oldCompanyId) {
        if (newCompanyId) {
          payload.payers = [{ company_id: newCompanyId, slot_index: 0 }];
        } else {
          payload.payers = []; // Handle clearing payer if valid
        }
      }

      if (Object.keys(payload).length === 0 && !data.isNew) {
        if (!silent) toast.info("Немає змін для збереження.");
        return true;
      }

      // Add linking IDs if new
      if (data.isNew) {
        if ((permit as any)._initial_camera_event)
          payload.camera_event_id = Number(
            (permit as any)._initial_camera_event,
          );
        if ((permit as any)._initial_scale_event)
          payload.scale_event_id = Number((permit as any)._initial_scale_event);
      }

      // We make direct fetch calls since we can't easily access coreClient locally from svelte component.
      // Another approach is submitting native forms to actions in +page.server.ts. Let's use fetch to our backend APIs
      // Warning: Usually it's better to use actions, but to be able to use locals.coreClient we would need an endpoint.
      // For now let's just use regular fetch if we have an endpoint mapped, but we don't.

      // Let's implement an action call instead.
      const form = new FormData();
      form.append("data", JSON.stringify(payload));

      const res = await fetch(`?/save`, {
        method: "POST",
        body: form,
      });

      const result = await res.json();

      if (result.type === "success") {
        if (!silent) toast.success("Перепустку успішно збережено!");

        // Optimistic UI update for main payload:
        data.permit = {
          ...data.permit,
          ...payload,
          customs_data: permit.customs_data,
        };

        // Optimistic UI update for Audit History:
        if (!data.isNew && Object.keys(payload).length > 0) {
          const optimisticChanges: any = {};
          // To simulate a git diff, we store both old and new for the UI
          for (const key of Object.keys(payload)) {
            // @ts-ignore
            const oldVal = data.permit[key] || "взагалі відсутнє/порожньо";
            const newVal = payload[key] || "видалено";
            optimisticChanges[key] = { from: oldVal, to: newVal };
          }

          const newAudit = {
            ID: Math.random(),
            action: silent ? "validate/close" : "update (очікує бекенд)",
            CreatedAt: new Date().toISOString(),
            comment: "Оптимістичне збереження",
            user: { first_name: "Ви", last_name: "(зараз)" },
            changes: optimisticChanges,
          };
          if (activeTab === "audit") {
            auditLogs = [newAudit, ...auditLogs];
          }
        }

        if (data.isNew) {
          goto("/permits");
        } else {
          // Just silently refresh missing data in background
          invalidateAll();
        }
        return true;
      } else {
        toast.error("Помилка при збереженні.", {
          description: result.error?.message || "Unknown error",
        });
        return false;
      }
    } catch (e: any) {
      toast.error("Непередбачена помилка", { description: e.message });
      return false;
    } finally {
      loading = false;
    }
  }

  async function fetchCustomsData() {
    if (!permit.declaration_number) {
      toast.warning("Введіть номер попередньої декларації!");
      return;
    }

    fetchingCustoms = true;
    try {
      // Direct fetch to proxy
      const res = await fetch(
        `/api/data-parser/customs/declaration/${permit.declaration_number}`,
      );
      if (!res.ok) {
        throw new Error(await res.text());
      }
      const customsData: PermitCustomsData = await res.json();

      permit.customs_data = customsData;
      toast.success("Дані з митниці успішно завантажено!");
    } catch (e: any) {
      toast.error("Не вдалося завантажити дані", { description: e.message });
    } finally {
      fetchingCustoms = false;
    }
  }

  async function handleClosePermit() {
    // Need dedicated endpoint to close
    loading = true;
    try {
      const form = new FormData();
      const res = await fetch(`?/close`, { method: "POST", body: form });
      const result = await res.json();
      if (result.type === "success") {
        toast.success("Перепустку закрито.");
        goto("/permits");
      } else {
        toast.error("Помилка закриття.", {
          description: result.error?.message,
        });
      }
    } catch (e: any) {
      toast.error("Непередбачена помилка", { description: e.message });
    } finally {
      loading = false;
    }
  }

  async function handleValidatePermit() {
    showValidationErrors = true;
    if (!isAllValid) {
      toast.error("Неможливо провалідувати перепустку.", {
        description:
          "Будь ласка, заповніть всі обов'язкові поля, позначені червоним.",
      });
      return;
    }

    // 2. Save first
    const saved = await handleSave(true);
    if (!saved) return;

    loading = true;
    try {
      const form = new FormData();
      const res = await fetch(`?/validate`, { method: "POST", body: form });
      const result = await res.json();
      if (result.type === "success") {
        toast.success("Перепустку провалідовано.");
        // Optimistic update
        permit.verified_at = new Date().toISOString();
        await invalidateAll();
      } else {
        toast.error("Помилка при валідації.", {
          description: result.error?.message,
        });
      }
    } catch (e: any) {
      toast.error("Непередбачена помилка", { description: e.message });
    } finally {
      loading = false;
    }
  }

  function formatDate(dateStr?: string) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleString("uk-UA");
  }

  async function fetchAuditLogs() {
    if (auditLogs.length > 0 || loadingAudits) return;

    loadingAudits = true;
    try {
      const res = await fetch(`/api/permits/${permit.ID}/audit`);
      if (res.ok) {
        auditLogs = await res.json();
      } else {
        toast.error("Не вдалося завантажити аудит");
      }
    } catch (e: any) {
      toast.error("Помилка мережі", { description: e.message });
    } finally {
      loadingAudits = false;
    }
  }

  $effect(() => {
    if (activeTab === "audit" && !data.isNew) {
      fetchAuditLogs();
    }
  });
</script>

<div class="container mx-auto py-8 max-w-7xl">
  <!-- Header -->
  <div class="flex items-center gap-6 mb-8">
    <Button
      variant="outline"
      size="icon"
      onclick={() => goto("/permits")}
      class="shrink-0 rounded-full h-10 w-10 border-slate-200 shadow-sm hover:bg-slate-50"
    >
      <ChevronLeft class="h-5 w-5" />
    </Button>
    <h1 class="text-3xl font-extrabold tracking-tight">
      {#if data.isNew}
        Створення нової перепустки
      {:else}
        Перепустка {permit.code || `#${permit.ID}`}
      {/if}
    </h1>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
    <!-- Main Content Column -->
    <div class="lg:col-span-8 space-y-8">
      <!-- ANPR Дані -->
      <!-- ANPR Дані -->
      <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
        <div class="p-4 bg-muted/30 border-b flex items-center gap-2">
          <UserCog class="h-4 w-4 text-slate-500" />
          <h3 class="font-semibold text-sm">Дані автомобіля</h3>
        </div>
        <div class="p-4 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label
                for="plate_front"
                class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                >Номер (Тягач)</Label
              >
              <Input
                id="plate-input"
                bind:value={permit.plate_front}
                disabled={permit.is_closed}
                class="font-mono font-black text-2xl uppercase tracking-[0.2em] h-14 bg-slate-50/50 dark:bg-slate-900/50 focus:bg-white dark:focus:bg-slate-950 transition-colors {showValidationErrors &&
                !isValidPlate
                  ? 'border-rose-500 ring-2 ring-rose-500/10'
                  : ''}"
                placeholder="XX0000XX"
              />
            </div>
            <div class="space-y-2">
              <Label
                for="plate_back"
                class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                >Номер (Причіп)</Label
              >
              <Input
                id="plate_back"
                bind:value={permit.plate_back}
                disabled={permit.is_closed}
                class="font-mono font-black text-2xl uppercase tracking-[0.2em] h-14 bg-slate-50/50 dark:bg-slate-900/50 focus:bg-white dark:focus:bg-slate-950 transition-colors"
                placeholder="XX0000XX"
              />
            </div>
          </div>
          <p class="text-xs text-muted-foreground mt-2">
            Відредагуйте номери, якщо вони були розпізнані з помилкою.
          </p>

          {#if permit.plate_events && permit.plate_events.length > 0}
            <div class="mt-4 pt-4 border-t">
              <span class="text-xs font-medium text-muted-foreground mb-2 block"
                >Фото з камери заїзду:</span
              >
              <div class="grid grid-cols-2 gap-4">
                {#each permit.plate_events as plate_event}
                  <div
                    class="aspect-video bg-muted rounded overflow-hidden border"
                  >
                  <p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground px-3">{plate_event.plate} {plate_event.camera_source_name}</p>
                    {#if plate_event.image_key} 
                    <img
                      src="/api/images/{plate_event.image_key}"
                      alt="Vehicle"
                      class="w-full h-full object-cover"
                    />
                  {:else}
                    <div
                      class="h-full w-full flex items-center justify-center text-xs text-muted-foreground"
                    >
                      Фото недоступне
                    </div>
                  {/if}
                </div>
              {/each}
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Вага -->
      <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
        <div class="p-4 bg-muted/30 border-b flex items-center gap-2">
          <Truck class="h-4 w-4 text-slate-500" />
          <h3 class="font-semibold text-sm">Вимірювання ваги</h3>
        </div>
        <div class="p-4">
          <div class="space-y-2">
            <Label
              for="weight"
              class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
              >Показник ваги, кг</Label
            >
            <Input
              id="weight-input"
              type="number"
              disabled={permit.is_closed}
              bind:value={permit.total_weight}
              class="font-mono font-black text-3xl h-16 text-blue-600 dark:text-blue-400 bg-blue-50/20 dark:bg-blue-900/10 {showValidationErrors &&
              !isValidWeight
                ? 'border-rose-500 ring-2 ring-rose-500/10'
                : ''}"
              placeholder="0"
            />
          </div>
        </div>
      </div>
      <!-- Митні дані -->
      <div class="bg-card border rounded-xl shadow-sm overflow-hidden relative">
        <!-- decorative border top -->
        <div class="absolute top-0 left-0 right-0 h-1 bg-indigo-500"></div>

        <div class="p-5">
          <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
            <FileSearch class="h-5 w-5 text-indigo-500" /> Оформлення в Єдиному вікні
            (Митниця)
          </h3>

          <div class="flex items-end gap-3 mb-6">
            <div id="decl-input" class="flex-1 space-y-2">
              <Label for="decl">Номер попередньої декларації (ПД)</Label>
              <Input
                id="decl"
                bind:value={permit.declaration_number}
                placeholder="Введіть номер..."
                class="bg-stone-50 dark:bg-stone-950 font-mono tracking-wide {showValidationErrors &&
                !isValidPD
                  ? 'border-rose-500 ring-2 ring-rose-500/10'
                  : ''}"
              />
            </div>
            <Button
              onclick={fetchCustomsData}
              disabled={fetchingCustoms || !permit.declaration_number}
              class="gap-2 w-[180px] bg-indigo-50 text-indigo-700 hover:bg-indigo-100 border border-indigo-200 dark:bg-indigo-950/30 dark:text-indigo-300 dark:border-indigo-800"
            >
              <FileSearch
                class="h-4 w-4 {fetchingCustoms ? 'animate-pulse' : ''}"
              />
              {fetchingCustoms ? "Пошук..." : "Отримати дані"}
            </Button>
          </div>

          <!-- Editable Customs Data Form -->
          {#if permit.customs_data}
            <div
              id="customs-data-section"
              class="bg-slate-50 dark:bg-slate-900/50 border rounded-lg p-4 grid grid-cols-2 gap-x-6 gap-y-4 text-sm animate-in fade-in slide-in-from-top-4 duration-300"
            >
              <div class="col-span-2 space-y-1">
                <Label class="text-xs font-semibold text-slate-500 uppercase"
                  >Опис Товару</Label
                >
                <Textarea
                  bind:value={permit.customs_data.goods}
                  class="bg-white dark:bg-slate-950 min-h-[60px]"
                  placeholder="Опис товару..."
                />
              </div>
              <div class="space-y-1">
                <Label class="text-xs font-medium text-slate-500"
                  >Декларант</Label
                >
                <Input
                  bind:value={permit.customs_data.declarant}
                  class="bg-white dark:bg-slate-950 h-8"
                  placeholder="Декларант..."
                />
              </div>
              <div class="space-y-1">
                <Label class="text-xs font-medium text-slate-500">ВМД</Label>
                <Input
                  bind:value={permit.customs_data.vmd_number}
                  class="bg-white dark:bg-slate-950 h-8 font-mono"
                  placeholder="ВМД..."
                />
              </div>
              <div class="space-y-1">
                <Label class="text-xs font-medium text-slate-500"
                  >Відправник</Label
                >
                <Input
                  bind:value={permit.customs_data.sender}
                  class="bg-white dark:bg-slate-950 h-8"
                  placeholder="Відправник..."
                />
              </div>
              <div class="space-y-1">
                <Label class="text-xs font-medium text-slate-500"
                  >Одержувач</Label
                >
                <Input
                  bind:value={permit.customs_data.receiver}
                  class="bg-white dark:bg-slate-950 h-8"
                  placeholder="Одержувач..."
                />
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Фінансові деталі та категорізація -->
      <div
        class="bg-card border rounded-xl shadow-sm overflow-hidden grid grid-cols-2"
      >
        <div
          class="col-span-2 p-4 bg-muted/20 border-b flex justify-between items-center"
        >
          <h3 class="font-semibold flex items-center gap-2">
            <Coins class="h-4 w-4 text-amber-500" /> Розрахункові параметри
          </h3>
          {#if !data.isNew && permit.total_sum}
            <div class="flex items-center gap-3">
              {#if permit.discount_amount && permit.discount_amount > 0}
                <span
                  class="text-xs text-rose-500 bg-rose-50 px-2 py-1 rounded font-bold"
                  >Знижка: -{permit.discount_amount} ₴</span
                >
              {/if}
              <div
                class="text-2xl font-black font-mono text-emerald-600 bg-emerald-50 dark:bg-emerald-950/30 px-5 py-2.5 rounded-xl border border-emerald-100 shadow-sm"
              >
                ДО СПЛАТИ: <span class="text-3xl tracking-tighter"
                  >{permit.total_sum} ₴</span
                >
              </div>
            </div>
          {/if}
        </div>

        <div class="p-5 space-y-5 border-r border-dashed">
          <div class="space-y-2">
            <Label>Категорія авто (Тарифікація)</Label>
            <Select.Root
              type="single"
              bind:value={permit.vehicle_type_id}
              name="vehicle_type"
            >
              <Select.Trigger
                id="category-input"
                class="w-full text-base font-bold h-12 {showValidationErrors &&
                !isValidVehicleType
                  ? 'border-rose-500 ring-2 ring-rose-500/10'
                  : ''}"
              >
                {#if permit.vehicle_type_id}
                  {data.vehicleTypes.find((v) => v.ID == permit.vehicle_type_id)
                    ?.name || "Вибрати..."}
                {:else}
                  Виберіть категорію авто...
                {/if}
              </Select.Trigger>
              <Select.Content>
                {#each data.vehicleTypes as vt}
                  <Select.Item value={vt.ID.toString()}
                    >{vt.name}
                    <span class="text-xs text-muted-foreground ml-2"
                      >({vt.entry_price} ₴ / {vt.daily_price} ₴/день)</span
                    ></Select.Item
                  >
                {/each}
              </Select.Content>
            </Select.Root>
          </div>

          <div class="space-y-2">
            <Label>Митний режим</Label>
            <Select.Root
              type="single"
              bind:value={permit.customs_mode_code}
              name="customs_mode"
            >
              <Select.Trigger
                id="mode-input"
                class="w-full text-base font-bold h-12 {showValidationErrors &&
                !isValidCustomsMode
                  ? 'border-rose-500 ring-2 ring-rose-500/10'
                  : ''}"
              >
                {#if permit.customs_mode_code}
                  {data.customsModes.find(
                    (m) => m.code == permit.customs_mode_code,
                  )?.name || "Вибрати..."}
                {:else}
                  Виберіть режим...
                {/if}
              </Select.Trigger>
              <Select.Content>
                {#each data.customsModes as cm}
                  <Select.Item value={cm.code.toString()}
                    >{cm.name}
                    <span
                      class="bg-muted px-1 ml-1 rounded text-[10px] font-mono"
                      >{cm.code}</span
                    ></Select.Item
                  >
                {/each}
              </Select.Content>
            </Select.Root>
          </div>
        </div>

        <div class="p-5 space-y-5">
          <div class="space-y-2">
            <Label>Компанія - Платник</Label>
            <Select.Root
              type="single"
              bind:value={primaryCompanyId}
              name="company"
            >
              <Select.Trigger class="w-full">
                {#if primaryCompanyId}
                  {data.companies.find((c) => c.ID == primaryCompanyId)?.name ||
                    "Вибрати..."}
                {:else}
                  Виберіть компанію...
                {/if}
              </Select.Trigger>
              <Select.Content>
                {#each data.companies as c}
                  <Select.Item value={c.ID.toString()}
                    >{c.name}
                    <span class="text-xs ml-2 text-muted-foreground"
                      >{c.edrpou}</span
                    ></Select.Item
                  >
                {/each}
              </Select.Content>
            </Select.Root>
            <p class="text-xs text-muted-foreground">
              Від компанії залежить розрахунок знижок.
            </p>
          </div>

          <div class="space-y-2">
            <Label>Тип оплати / Договір</Label>
            <Select.Root
              type="single"
              bind:value={permit.payment_type_id}
              name="payment_type"
            >
              <Select.Trigger class="w-full">
                {#if permit.payment_type_id}
                  {data.paymentTypes.find((p) => p.ID == permit.payment_type_id)
                    ?.name || "Вибрати..."}
                {:else}
                  Виберіть спосіб оплати...
                {/if}
              </Select.Trigger>
              <Select.Content>
                {#each data.paymentTypes as p}
                  <Select.Item value={p.ID.toString()}>{p.name}</Select.Item>
                {/each}
              </Select.Content>
            </Select.Root>
          </div>
        </div>

        <div class="col-span-2 p-5 pt-0">
          <div class="space-y-2">
            <Label for="notes">Примітки</Label>
            <Textarea
              id="notes"
              bind:value={permit.notes}
              placeholder="Додаткова інформація..."
              class="min-h-[80px]"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Sidebar Column -->
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
                  <span class="font-bold">{formatDate(permit.exit_time)}</span>
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
                    >
                      ₴{permit.total_sum.toFixed(2)}
                    </span>
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
              <CheckCircle class="h-4 w-4 text-amber-600" />
              <h3
                class="font-black text-[10px] uppercase tracking-widest text-amber-700"
              >
                Чек-лист валідації
              </h3>
            </div>
            <div class="p-5 space-y-4">
              {@render validationItem(
                "Номери (перед/зад)",
                isValidPlate,
                "plate-input",
              )}
              {@render validationItem("Вага", isValidWeight, "weight-input")}
              {@render validationItem("Номер ПД", isValidPD, "decl-input")}
              {@render validationItem(
                "Дані митниці",
                isValidCustomsData,
                "customs-data-section",
              )}
              {@render validationItem(
                "Категорія авто",
                isValidVehicleType,
                "category-input",
              )}
              {@render validationItem(
                "Митний режим",
                isValidCustomsMode,
                "mode-input",
              )}
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
              <CheckCircle class="h-5 w-5" /> Валідувати
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
                <AlertDialog.Title>Підтвердження закриття</AlertDialog.Title>
                <AlertDialog.Description>
                  Ви впевнені, що хочете закрити цю перепустку?
                </AlertDialog.Description>
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
  </div>

  {#snippet validationItem(label: string, isValid: boolean, targetId?: string)}
    <div class="flex items-center gap-3 py-0.5">
      {#if isValid}
        <div
          class="h-5 w-5 rounded-full bg-emerald-500 flex items-center justify-center shadow-sm shrink-0"
        >
          <CheckCircle class="h-3.5 w-3.5 text-white" />
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
            document
              .getElementById(targetId)
              ?.scrollIntoView({ behavior: "smooth", block: "center" })}
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

  <!-- Audit & Events History -->
  {#if !data.isNew}
    <div class="mt-12 pt-12 border-t">
      <Tabs.Root bind:value={activeTab} class="w-full">
        <Tabs.List class="grid w-full grid-cols-2 mb-8 h-12">
          <Tabs.Trigger
            value="events"
            class="flex items-center gap-3 text-base font-bold"
          >
            <Activity class="h-5 w-5" />
            Події
          </Tabs.Trigger>
          <Tabs.Trigger
            value="audit"
            class="flex items-center gap-3 text-base font-bold"
          >
            <Clock class="h-5 w-5" />
            Історія змін
          </Tabs.Trigger>
        </Tabs.List>

        <Tabs.Content value="audit" class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-lg flex items-center gap-2">
              <Clock class="h-5 w-5 text-slate-500" />
              Логи аудиту
            </h3>
            <!-- <Button variant="outline" size="sm" onclick={toggleAudits}>
              {showAudits ? "Приховати аудит" : "Завантажити історію"}
            </Button> -->
          </div>

          {#if loadingAudits}
            <div
              class="flex items-center justify-center p-8 text-muted-foreground"
            >
              Завантаження...
            </div>
          {:else if auditLogs.length === 0}
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
                      <div class="font-bold flex items-center gap-3 text-lg">
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
                        <p class="text-muted-foreground text-sm font-medium">
                          {audit.comment}
                        </p>
                      {/if}
                      {#if audit.changes && Object.keys(audit.changes).length > 0}
                        <div
                          class="mt-4 text-sm border rounded-lg overflow-hidden bg-background shadow-inner"
                        >
                          <table class="w-full text-left">
                            <thead class="bg-muted/70 text-muted-foreground">
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
                                  <td class="px-4 py-3 font-mono text-sm">
                                    <span
                                      class="font-black inline-block text-slate-700 dark:text-slate-300"
                                    >
                                      {typeof val === "object"
                                        ? JSON.stringify(val)
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
        </Tabs.Content>

        <Tabs.Content value="events" class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-lg flex items-center gap-2">
              <Activity class="h-5 w-5 text-slate-500" />
              Події з Камери / Ваги
            </h3>
          </div>

          <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
            <table class="w-full text-sm text-left">
              <thead
                class="bg-muted/70 text-muted-foreground uppercase tracking-widest text-[11px] font-black"
              >
                <tr>
                  <th class="px-6 py-4">ID</th>
                  <th class="px-6 py-4">Час</th>
                  <th class="px-6 py-4">Тип</th>
                  <th class="px-6 py-4">Джерело</th>
                  <th class="px-6 py-4">Значення</th>
                </tr>
              </thead>
              <tbody class="divide-y">
                {#each [...(permit.plate_events || []), ...(permit.weight_events || [])].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()) as event}
                  {@const eventType = "plate" in event ? "plate" : "weight"}
                  <tr class="hover:bg-muted/30 transition-colors">
                    <td
                      class="px-6 py-4 text-muted-foreground whitespace-nowrap font-mono"
                    >
                      <Button
                        variant="link"
                        href={`/events/${eventType}/${event.ID}`}
                        class="p-0 h-auto font-bold"
                      >
                        #{event.ID}
                      </Button>
                    </td>
                    <td
                      class="px-6 py-4 text-muted-foreground whitespace-nowrap font-medium"
                    >
                      {formatDate(event.CreatedAt)}
                    </td>
                    <td class="px-6 py-4">
                      {#if "plate" in event}
                        <span
                          class="inline-flex items-center gap-2 bg-blue-50 text-blue-700 px-3 py-1 rounded-full text-[10px] font-black uppercase border border-blue-100"
                        >
                          <Truck class="h-3 w-3" /> Номер
                        </span>
                      {:else}
                        <span
                          class="inline-flex items-center gap-2 bg-amber-50 text-amber-700 px-3 py-1 rounded-full text-[10px] font-black uppercase border border-amber-100"
                        >
                          <Coins class="h-3 w-3" /> Ваги
                        </span>
                      {/if}
                    </td>
                    <td class="px-6 py-4 font-mono text-xs font-semibold">
                      {#if "camera_source_id" in event}
                        {event.camera_source_id || "Camera"}
                      {:else}
                        {event.scale_source_id || "Scale"}
                      {/if}
                    </td>
                    <td
                      class="px-6 py-4 font-black text-base text-slate-900 dark:text-slate-100"
                    >
                      {#if "plate" in event}
                        {event.plate}
                      {:else}
                        {event.weight} кг
                      {/if}
                    </td>
                  </tr>
                {/each}
                {#if !permit.plate_events?.length && !permit.weight_events?.length}
                  <tr>
                    <td
                      colspan="4"
                      class="px-4 py-8 text-center text-muted-foreground"
                    >
                      Подій не знайдено
                    </td>
                  </tr>
                {/if}
              </tbody>
            </table>
          </div>
        </Tabs.Content>
      </Tabs.Root>
    </div>
  {/if}
</div>
