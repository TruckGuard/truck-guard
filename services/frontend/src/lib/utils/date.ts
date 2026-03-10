/**
 * Shared date formatting utilities for the frontend.
 * All functions use the "uk-UA" locale for Ukrainian formatting.
 */

/**
 * Formats a date+time into a human-readable Ukrainian locale string (DD.MM.YYYY, HH:MM).
 * Returns "-" if the input is empty/undefined.
 */
export function formatDate(date?: string | Date | null): string {
    if (!date) return "-";
    const d = typeof date === "string" ? new Date(date) : date;
    return d.toLocaleString("uk-UA", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}

/**
 * Formats a date into a Ukrainian locale date-only string (DD.MM.YYYY).
 * Returns "-" if the input is empty/undefined.
 */
export function formatDateOnly(date?: string | Date | null): string {
    if (!date) return "-";
    const d = typeof date === "string" ? new Date(date) : date;
    return d.toLocaleDateString("uk-UA", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
    });
}

/**
 * Formats a date string to show only the time portion (HH:MM:SS).
 * Returns "-" if the input is empty/undefined.
 */
export function formatTimeOnly(dateStr?: string | null): string {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleTimeString("uk-UA", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
    });
}
