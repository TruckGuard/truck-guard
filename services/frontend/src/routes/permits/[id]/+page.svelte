<script lang="ts">
  import { goto } from "$app/navigation";
  import { Clock, Activity } from "@lucide/svelte";
  import * as Tabs from "$lib/components/ui/tabs";
  import type { Permit } from "$lib/types/permits";
  import type { Company } from "$lib/types/data";
  import * as AlertDialog from "$lib/components/ui/alert-dialog";
  import { toast } from "svelte-sonner";
  import { invalidateAll } from "$app/navigation";

  // Component Imports
  import PermitHeader from "./components/PermitHeader.svelte";
  import VehicleDataSection from "./components/VehicleDataSection.svelte";
  import WeightSection from "./components/WeightSection.svelte";
  import CustomsDataSection from "./components/CustomsDataSection.svelte";
  import FinancialSection from "./components/FinancialSection.svelte";
  import PermitSidebar from "./components/PermitSidebar.svelte";
  import AuditLogsTable from "./components/AuditLogsTable.svelte";
  import EventsHistoryTable from "./components/EventsHistoryTable.svelte";
  import { formatDate } from "$lib/utils/date";
  import { onMount } from "svelte";
    import CustomsModeBanner from "./components/CustomsModeBanner.svelte";

  interface PayerState {
    company_id: number;
    slot_index: number;
    company: Company | null;
  }

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

  // Create a mutable copy of the permit to bind to the form
  let permit = $state<Partial<Permit>>({});

  // Payer state array
  let payers = $state<PayerState[]>([]);

  let activeTab = $state("events");
  let loadingAudits = $state(false);
  let auditLogs: any[] = $state([]);

  // Validation state
  let showValidationErrors = $state(false);

  // Field label map shared between ModeForm & permit checklist
  const FIELD_LABELS: Record<string, string> = {
    plate_front: "Номер (перед)",
    plate_back: "Номер (зад)",
    total_weight: "Вага",
    declaration_number: "Номер ПД",
    "customs_data.goods": "Вантаж",
    "customs_data.vmd_number": "Номер ВМД",
    "customs_data.declarant": "Декларант",
    "customs_data.sender": "Відправник",
    "customs_data.receiver": "Отримувач",
    vehicle_type_id: "Категорія авто",
    customs_mode_code: "Митний режим",
    payer: "Платник (мінімум один)",
  };

  // Target scroll IDs per field
  const FIELD_TARGETS: Record<string, string> = {
    plate_front: "plate-input",
    plate_back: "plate-input",
    total_weight: "weight-input",
    declaration_number: "decl-input",
    "customs_data.goods": "customs-data-section",
    "customs_data.vmd_number": "customs-data-section",
    "customs_data.declarant": "customs-data-section",
    "customs_data.sender": "customs-data-section",
    "customs_data.receiver": "customs-data-section",
    vehicle_type_id: "category-input",
    customs_mode_code: "mode-input",
    payer: "financial-section",
  };

  const selectedMode = $derived(
    data.customsModes.find((m: any) => m.code === permit.customs_mode_code),
  );

  const requiredFields = $derived(
    (selectedMode?.required_fields as string[]) ?? [],
  );

  function isFieldValid(key: string): boolean {
    switch (key) {
      case "plate_front":
        return !!permit.plate_front;
      case "plate_back":
        return !!permit.plate_back;
      case "total_weight":
        return !!(Number(permit.total_weight) || 0);
      case "declaration_number":
        return !!permit.declaration_number;
      case "customs_data.goods":
        return !!permit.customs_data?.goods;
      case "customs_data.vmd_number":
        return !!permit.customs_data?.vmd_number;
      case "customs_data.declarant":
        return !!permit.customs_data?.declarant;
      case "customs_data.sender":
        return !!permit.customs_data?.sender;
      case "customs_data.receiver":
        return !!permit.customs_data?.receiver;
      case "vehicle_type_id":
        return !!permit.vehicle_type_id;
      case "customs_mode_code":
        return !!permit.customs_mode_code;
      case "payer":
        return payers.some((p) => p.company && (p.company.ID || 0) > 0);
      default:
        return true;
    }
  }

  const validationItems = $derived(
    requiredFields.map((key: string) => ({
      key,
      label: FIELD_LABELS[key] ?? key,
      targetId: FIELD_TARGETS[key] ?? undefined,
      isValid: isFieldValid(key),
    })),
  );

  const isAllValid = $derived(
    validationItems.every((item: any) => item.isValid) &&
      permit.customs_mode_code,
  );

  const isValidPlate = $derived(!!permit.plate_front && !!permit.plate_back);
  const isValidWeight = $derived(!!(Number(permit.total_weight) || 0));
  const isValidPD = $derived(!!permit.declaration_number);
  const isValidVehicleType = $derived(!!permit.vehicle_type_id);
  const isValidCustomsMode = $derived(!!permit.customs_mode_code);

  onMount(() => {
    permit = { ...data.permit };

    // Sync payers from permit data
    if (permit.payers?.length) {
      payers = permit.payers.map((p) => ({
        company_id: p.company_id,
        slot_index: p.slot_index,
        company:
          p.company ||
          data.companies.find((c: any) => c.ID === p.company_id) ||
          null,
      }));
    } else {
      payers = [{ company_id: 0, slot_index: 1, company: null }];
    }

    if (!permit.customs_data) {
      permit.customs_data = {
        goods: "",
        declarant: "",
        vmd_number: "",
        sender: "",
        receiver: "",
      };
    }
  });

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

      // Map valid payers to payload
      const payloadPayers = payers
        .filter((p) => p.company) // Only slots with a selected company
        .map((p, i) => ({
          company_id: p.company?.ID,
          slot_index: p.slot_index || i + 1,
        }));

      // Submit updated payers
      payload.payers = payloadPayers;

      if (Object.keys(payload).length === 0) {
        if (!silent) toast.info("Немає змін для збереження.");
        return true;
      }

      // Add linking IDs if new
      // Event linking is now handled during initial creation via server action

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
        if (Object.keys(payload).length > 0) {
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

        // Just silently refresh missing data in background
        invalidateAll();
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

  async function handleRestorePermit() {
    loading = true;
    // Optimistic update
    const previousVoidStatus = permit.is_void;
    permit.is_void = false;

    try {
      const form = new FormData();
      const res = await fetch(`?/restore`, { method: "POST", body: form });
      const result = await res.json();
      if (result.type === "success") {
        toast.success("Перепустку відновлено.");
        await invalidateAll();
      } else {
        permit.is_void = previousVoidStatus; // Rollback
        toast.error("Помилка відновлення.", {
          description: result.error?.message,
        });
      }
    } catch (e: any) {
      permit.is_void = previousVoidStatus; // Rollback
      toast.error("Непередбачена помилка", { description: e.message });
    } finally {
      loading = false;
    }
  }

  async function handleDeletePermit() {
    loading = true;
    try {
      const form = new FormData();
      const res = await fetch(`?/delete`, { method: "POST", body: form });
      const result = await res.json();
      if (result.type === "success") {
        toast.success("Перепустку остаточно видалено.");
        goto("/permits");
      } else {
        toast.error("Помилка видалення.", {
          description: result.error?.message,
        });
      }
    } catch (e: any) {
      toast.error("Непередбачена помилка", { description: e.message });
    } finally {
      loading = false;
    }
  }

  async function handleVoidPermit() {
    loading = true;
    // Optimistic update
    const previousVoidStatus = permit.is_void;
    permit.is_void = true;
    
    try {
      const form = new FormData();
      const res = await fetch(`?/void`, { method: "POST", body: form });
      const result = await res.json();
      if (result.type === "success") {
        toast.success("Перепустку анульовано.");
        await invalidateAll();
      } else {
        permit.is_void = previousVoidStatus; // Rollback
        toast.error("Помилка анулювання.", {
          description: result.error?.message,
        });
      }
    } catch (e: any) {
      permit.is_void = previousVoidStatus; // Rollback
      toast.error("Непередбачена помилка", { description: e.message });
    } finally {
      loading = false;
    }
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
    if (activeTab === "audit") {
      fetchAuditLogs();
    }
  });
</script>

<div class="container mx-auto py-8 max-w-7xl">
  <PermitHeader
    permitCode={permit.code}
    permitId={permit.ID}
  />

  <CustomsModeBanner
    bind:permit
    customsModes={data.customsModes}
    verified={!!permit.verified_at}
  />

  <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
    <div class="lg:col-span-8 space-y-8">
      <VehicleDataSection bind:permit {showValidationErrors} {isValidPlate} />
      <WeightSection bind:permit {showValidationErrors} {isValidWeight} />
      <CustomsDataSection bind:permit {showValidationErrors} {isValidPD} />
      <FinancialSection
        bind:permit
        {data}
        bind:payers
        {showValidationErrors}
        {isValidVehicleType}
        {isValidCustomsMode}
      />
    </div>

    <!-- Sidebar Column -->
    <PermitSidebar
      {permit}
      {data}
      {loading}
      {handleSave}
      {handleValidatePermit}
      {handleClosePermit}
      handleRestore={handleRestorePermit}
      handleVoid={handleVoidPermit}
      handleDelete={handleDeletePermit}
      {validationItems}
    />
  </div>

  <!-- Audit & Events History -->
  {#if permit.ID}
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

        <Tabs.Content value="audit">
          <AuditLogsTable {auditLogs} />
        </Tabs.Content>

        <Tabs.Content value="events">
          <EventsHistoryTable {permit} />
        </Tabs.Content>
      </Tabs.Root>
    </div>
  {/if}
</div>
