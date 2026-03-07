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
                enableSorting: true,
            },
            {
                accessorKey: "entry_time",
                id: "entry_time",
                header: "Час заїзду",
                cell: ({ row }) => formatDate(row.getValue("entry_time")),
                enableSorting: true,
            },
            {
                accessorKey: "exit_time",
                id: "exit_time",
                header: "Час виїзду",
                cell: ({ row }) => {
                    const v = row.getValue("exit_time") as string;
                    return v ? formatDate(v) : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "plate_front",
                id: "plate_front",
                header: "Номер передній",
                cell: ({ row }) => {
                    const plate = row.getValue("plate_front") as string;
                    return plate || "-";
                },
                enableSorting: false,
            },
            {
                accessorKey: "plate_back",
                id: "plate_back",
                header: "Номер задній",
                cell: ({ row }) => {
                    const plate = row.getValue("plate_back") as string;
                    return plate || "-";
                },
                enableSorting: false,
            },
            {
                accessorKey: "total_weight",
                id: "total_weight",
                header: "Вага",
                cell: ({ row }) => {
                    const w = row.getValue("total_weight") as number;
                    return w ? `${w.toLocaleString("uk-UA")} кг` : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "status",
                id: "is_closed",
                header: "Статус",
                cell: ({ row }) => {
                    return renderComponent(CellStatus, {
                        isClosed: row.original.is_closed,
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

                    if (!days) return "-";

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
                    return p ? p.name : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "customs_mode",
                id: "customs_mode",
                header: "Режим",
                cell: ({ row }) => {
                    const m = row.original.customs_mode;
                    return m ? m.name : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "declaration_number",
                id: "declaration_number",
                header: "Номер декларації",
                cell: ({ row }) => {
                    return row.original.declaration_number || "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "customs_data.declarant",
                id: "customs_declarant_name",
                header: "Декларант",
                cell: ({ row }) => {
                    return row.original.customs_data?.declarant || "-";
                },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.goods",
                id: "customs_commodity_description",
                header: "Опис товару",
                cell: ({ row }) => {
                    return row.original.customs_data?.goods || "-";
                },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.vmd_number",
                id: "customs_vmd_number",
                header: "ВМД",
                cell: ({ row }) => {
                    return row.original.customs_data?.vmd_number || "-";
                },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.sender",
                id: "customs_sender",
                header: "Відправник",
                cell: ({ row }) => {
                    return row.original.customs_data?.sender || "-";
                },
                enableSorting: false,
            },
            {
                accessorKey: "customs_data.receiver",
                id: "customs_receiver",
                header: "Одержувач",
                cell: ({ row }) => {
                    return row.original.customs_data?.receiver || "-";
                },
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
                    return payment ? payment.name : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "entry_fee",
                id: "entry_fee",
                header: "Вхідна плата",
                cell: ({ row }) => {
                    const fee = row.getValue("entry_fee") as number;
                    return fee ? `${fee.toLocaleString("uk-UA")} грн` : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "exit_fee",
                id: "exit_fee",
                header: "Вихідна плата",
                cell: ({ row }) => {
                    const fee = row.getValue("exit_fee") as number;
                    return fee ? `${fee.toLocaleString("uk-UA")} грн` : "-";
                },
                enableSorting: true,
            },
            {
                accessorKey: "total_sum",
                id: "total_sum",
                header: "Загальна плата",
                cell: ({ row }) => {
                    const fee = row.getValue("total_sum") as number;
                    return fee ? `${fee.toLocaleString("uk-UA")} грн` : "-";
                },
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
                enableSorting: true,
            },
        ]
    },
];
