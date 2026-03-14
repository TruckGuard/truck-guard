import { renderComponent } from "$lib/components/ui/data-table";
import type { ColumnDef } from "@tanstack/table-core";
import type { Permit } from "$lib/types/permits";
import { formatDate } from "$lib/utils/date";

// Custom cell components
import CellPayers from "./components/cell-payers.svelte";
import CellVehicle from "./components/cell-vehicle.svelte";
import CellStatus from "./components/cell-status.svelte";
import CellVerifier from "./components/cell-verifier.svelte";

export const columns: ColumnDef<Permit>[] = [
    {
        id: "group_main",
        header: "Основна інформація",
        columns: [
            {
                accessorKey: "id",
                id: "id",
                header: "Код/Номер",
                cell: ({ row }) => {
                    const p = row.original;
                    return p.code || `PF-${p.ID}`;
                },
                meta: { mono: true, className: "text-base font-bold text-blue-600 dark:text-blue-500/90" },
                enableSorting: true,
            },
            {
                accessorKey: "entry_time",
                id: "entry_time",
                header: "Час заїзду",
                cell: ({ row }) => formatDate(row.getValue("entry_time")),
                meta: { className: "text-sm text-slate-500 dark:text-zinc-500" },
                enableSorting: true,
            },
            {
                accessorKey: "exit_time",
                id: "exit_time",
                header: "Час виїзду",
                cell: ({ row }) => {
                    const v = row.getValue("exit_time") as string;
                    return v ? formatDate(v) : "—";
                },
                meta: { className: "text-sm text-slate-500 dark:text-zinc-500" },
                enableSorting: true,
            },
            {
                accessorKey: "plate_front",
                id: "plate_front",
                header: "Номер передній",
                cell: ({ row }) => (row.getValue("plate_front") as string) || "—",
                meta: { mono: true, fontMedium: true, className: "text-slate-900 dark:text-zinc-100 text-base" },
                enableSorting: false,
            },
            {
                accessorKey: "plate_back",
                id: "plate_back",
                header: "Номер задній",
                cell: ({ row }) => (row.getValue("plate_back") as string) || "—",
                meta: { mono: true, fontMedium: true, className: "text-slate-900 dark:text-zinc-100 text-base" },
                enableSorting: false,
            },
            {
                accessorKey: "total_weight",
                id: "total_weight",
                header: "Вага",
                cell: ({ row }) => {
                    const w = row.getValue("total_weight") as number;
                    return w ? `${w.toLocaleString("uk-UA")} кг` : "—";
                },
                meta: { align: "right", tabular: true, className: "font-medium text-slate-900 dark:text-zinc-100 text-base" },
                enableSorting: true,
            },
            {
                accessorKey: "status",
                id: "is_closed",
                header: "Статус",
                cell: ({ row }) => {
                    return renderComponent(CellStatus, {
                        isClosed: row.original.is_closed,
                        isVoid: row.original.is_void,
                    });
                },
                enableSorting: true,
            },
            {
                accessorKey: "days_in_zone",
                id: "days_in_zone",
                header: "Дні в зоні",
                cell: ({ row }) => {
                    let days = row.getValue("days_in_zone") as number;
                    const p = row.original;

                    if (!days && !p.is_closed && p.verified_at) {
                        const entry = new Date(p.entry_time).getTime();
                        const now = Date.now();
                        days = Math.max(1, Math.ceil((now - entry) / (1000 * 60 * 60 * 24)));
                    }

                    if (!days) return "—";

                    const lastDigit = days % 10;
                    const lastTwo = days % 100;
                    let word = "днів";
                    if (lastDigit === 1 && lastTwo !== 11) {
                        word = "день";
                    } else if ([2, 3, 4].includes(lastDigit) && ![12, 13, 14].includes(lastTwo)) {
                        word = "дні";
                    }

                    return `${!p.is_closed ? "~ " : ""}${days} ${word}`;
                },
                meta: { align: "right", tabular: true, fontMedium: true, className: "text-slate-900 dark:text-zinc-100 text-base" },
                enableSorting: true,
            },
        ]
    },
    {
        id: "group_customs",
        header: "Оформлення в Єдиному вікні (Митниця)",
        columns: [
            {
                accessorKey: "customs_post",
                id: "customs_post_id",
                header: "Митний пост",
                cell: ({ row }) => {
                    const p = row.original.customs_post;
                    return p ? p.name : "—";
                },
                meta: { className: "text-slate-800 dark:text-zinc-300" },
                enableSorting: true,
            },
            {
                accessorKey: "customs_mode",
                id: "customs_mode",
                header: "Режим",
                cell: ({ row }) => {
                    const m = row.original.customs_mode;
                    return m ? m.name : "—";
                },
                meta: { className: "text-slate-800 dark:text-zinc-300" },
                enableSorting: true,
            },
            {
                accessorKey: "declaration_number",
                id: "declaration_number",
                header: "Номер декларації",
                cell: ({ row }) => {
                    return row.original.declaration_number || "—";
                },
                meta: { className: "text-slate-800 dark:text-zinc-300" },
                enableSorting: true,
            },
            {
                accessorKey: "customs_data.declarant",
                id: "customs_declarant_name",
                header: "Декларант",
                cell: ({ row }) => {
                    return row.original.customs_data?.declarant || "—";
                },
                meta: { className: "text-slate-600 dark:text-zinc-400 text-sm" },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.goods",
                id: "customs_commodity_description",
                header: "Опис товару",
                cell: ({ row }) => {
                    return row.original.customs_data?.goods || "—";
                },
                meta: { className: "text-slate-600 dark:text-zinc-400 text-sm" },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.vmd_number",
                id: "customs_vmd_number",
                header: "ВМД",
                cell: ({ row }) => {
                    return row.original.customs_data?.vmd_number || "—";
                },
                meta: { className: "text-slate-600 dark:text-zinc-400 text-sm" },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.sender",
                id: "customs_sender",
                header: "Відправник",
                cell: ({ row }) => {
                    return row.original.customs_data?.sender || "—";
                },
                meta: { className: "text-slate-600 dark:text-zinc-400 text-sm" },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.receiver",
                id: "customs_receiver",
                header: "Одержувач",
                cell: ({ row }) => {
                    return row.original.customs_data?.receiver || "—";
                },
                meta: { className: "text-slate-600 dark:text-zinc-400 text-sm" },
                enableSorting: false,
            },
        ]
    },
    {
        id: "group_finance",
        header: "Фінанси та транспорт",
        columns: [
            {
                accessorKey: "vehicle_type",
                id: "vehicle_type_id",
                header: "Тип авто",
                cell: ({ row }) => {
                    return renderComponent(CellVehicle, {
                        vehicleType: row.original.vehicle_type || null,
                    });
                },
                enableSorting: true,
            },
            {
                accessorKey: "payment_type",
                id: "payment_type_id",
                header: "Тип оплати",
                cell: ({ row }) => {
                    const payment = row.original.payment_type;
                    return payment ? payment.name : "—";
                },
                meta: { className: "text-slate-800 dark:text-zinc-300" },
                enableSorting: true,
            },
            {
                accessorKey: "entry_fee",
                id: "entry_fee",
                header: "Вхідна плата",
                cell: ({ row }) => {
                    const fee = row.getValue("entry_fee") as number;
                    return fee ? `${fee.toLocaleString("uk-UA")} грн` : "—";
                },
                meta: { align: "right", tabular: true, className: "text-slate-500 dark:text-zinc-500 text-sm" },
                enableSorting: true,
            },
            {
                accessorKey: "daily_fee",
                id: "daily_fee",
                header: "Денна плата",
                cell: ({ row }) => {
                    const fee = row.getValue("daily_fee") as number;
                    return fee ? `${fee.toLocaleString("uk-UA")} грн` : "—";
                },
                meta: { align: "right", tabular: true, className: "text-slate-500 dark:text-zinc-500 text-sm" },
                enableSorting: true,
            },
            {
                accessorKey: "total_sum",
                id: "total_sum",
                header: "Загальна плата",
                cell: ({ row }) => {
                    const p = row.original;
                    let fee = p.total_sum;

                    if (!p.is_closed && p.entry_time) {
                        const entry = new Date(p.entry_time).getTime();
                        const now = Date.now();
                        const hours = (now - entry) / (1000 * 60 * 60);

                        let days = Math.floor(hours / 24);
                        if (hours > days * 24 || (now - entry) <= 0) {
                            days++;
                        }
                        if (days < 1) days = 1;

                        const entryFee = p.entry_fee || p.vehicle_type?.entry_price || 0;
                        const dailyPrice = p.daily_fee || p.vehicle_type?.daily_price || 0;
                        fee = entryFee + days * dailyPrice;
                    }

                    return fee ? `${!p.is_closed ? "~ " : ""}${fee.toLocaleString("uk-UA")} грн` : "—";
                },
                meta: { align: "right", tabular: true, fontSemiBold: true, className: "text-slate-950 dark:text-zinc-50 text-base" },
                enableSorting: true,
            },
            {
                accessorKey: "payers",
                id: "payers",
                header: "Платники",
                cell: ({ row }) => {
                    return renderComponent(CellPayers, {
                        payers: row.original.payers || [],
                    });
                },
                enableSorting: false,
            },
        ]
    },
    {
        id: "group_validation",
        header: "Валідація",
        columns: [
            {
                accessorKey: "verifier",
                id: "verified_at",
                header: "Валідація",
                cell: ({ row }) => {
                    return renderComponent(CellVerifier, {
                        verifier: row.original.verifier || null,
                        verifiedAt: row.original.verified_at || null,
                    });
                },
                meta: { className: "text-slate-600 dark:text-zinc-400" },
                enableSorting: true,
            },
        ]
    },
];
