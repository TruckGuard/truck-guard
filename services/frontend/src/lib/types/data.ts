import type { PaginationMetadata } from "./events";

export interface GormModel {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
}

export interface CustomsPost extends GormModel {
    name: string;
    description: string;
}

export interface CustomsMode extends GormModel {
    name: string;
    code: string;
    description: string;
    required_fields: string[];
}

export interface Company extends GormModel {
    name: string;
    edrpou: string;
    details: Record<string, any>;
    last_synced_at: string | null;
}

export interface VehicleType extends GormModel {
    name: string;
    code: string;
    description: string;
    entry_price: number;
    daily_price: number;
    color: string;
}

export interface PaymentType extends GormModel {
    name: string;
    code: string;
    description: string;
    is_active: boolean;
    icon: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    metadata: PaginationMetadata;
}
