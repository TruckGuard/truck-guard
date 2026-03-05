/**
 * Shared date formatting utilities for the frontend.
 * All functions use the "uk-UA" locale for Ukrainian formatting.
 */

/**
 * Formats a date+time string into a human-readable Ukrainian locale string.
 * Returns "-" if the input is empty/undefined.
 */
export function formatDate(dateStr?: string | null): string {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleString("uk-UA");
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
