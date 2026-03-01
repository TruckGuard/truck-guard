import type { Company, CustomsMode, CustomsPost, PaymentType, VehicleType } from "./data";
import type { RawPlateEvent, RawWeightEvent } from "./events";

export interface PermitCustomsData {
    ID: number;
    permit_id: number;
    declarant: string;
    goods: string;
    sender: string;
    receiver: string;
    vmd_number: string;
}

export interface PermitPayer {
    ID: number;
    permit_id: number;
    company_id: number;
    slot_index: number;
    company?: Company;
}

export interface PermitAudit {
    ID: number;
    permit_id: number;
    user_id?: number;
    user?: {
        id: number;
        first_name: string;
        last_name: string;
    };
    action: string;
    changes: any;
    comment: string;
    created_at: string;
}

export interface Permit {
    ID: number;
    code: string;
    is_closed: boolean;
    is_void: boolean;

    // Customs
    customs_post_id?: number;
    customs_post?: CustomsPost;
    declaration_number?: string;
    customs_mode_code?: string;
    customs_mode?: CustomsMode;
    notes?: string;

    // Vehicle
    vehicle_type_id?: number;
    vehicle_type?: VehicleType;
    plate_front?: string;
    plate_back?: string;
    total_weight?: number;

    // Financials
    payment_type_id?: number;
    payment_type?: PaymentType;
    entry_fee: number;
    exit_fee: number;
    total_sum: number;
    discount_amount: number;

    payers?: PermitPayer[];

    // Time
    entry_time: string;
    exit_time?: string;
    days_in_zone?: number;
    last_activity_at: string;

    // Users
    created_by?: number;
    verified_by?: number;
    verified_at?: string;
    responsible_user_id?: number;
    responsible_user?: {
        first_name: string;
        last_name: string;
    };
    verifier?: {
        first_name: string;
        last_name: string;
    };

    // Relations
    customs_data?: PermitCustomsData;
    plate_events?: RawPlateEvent[];
    weight_events?: RawWeightEvent[];
    audit_events?: PermitAudit[];

    CreatedAt: string;
    UpdatedAt: string;
}
