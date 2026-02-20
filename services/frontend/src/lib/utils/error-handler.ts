/**
 * Maps technical error messages from the backend to user-friendly Ukrainian strings.
 */
export function mapErrorToFriendlyMessage(error: string | any): string {
  if (!error) return "Ой, щось пішло не так";
  
  const errorStr = typeof error === 'string' ? error : JSON.stringify(error);
  
  // Unique constraint violations (Database specific patterns)
  if (errorStr.includes("unique constraint") || errorStr.includes("duplicate key")) {
    if (errorStr.includes("code")) {
      return "Ой, такий код вже існує";
    }
    if (errorStr.includes("name")) {
      return "Ой, такий запис з такою назвою вже існує";
    }
    if (errorStr.includes("edrpou")) {
      return "Ой, компанія з таким ЄДРПОУ вже існує";
    }
    return "Ой, такий запис вже існує";
  }

  // Not found
  if (errorStr.includes("not found")) {
    return "Запис не знайдено";
  }

  // Permission errors
  if (errorStr.includes("permission") || errorStr.includes("unauthorized") || errorStr.includes("forbidden")) {
    return "У вас недостатньо прав для цієї дії";
  }

  // Fallback for user-correctable errors if they are short enough
  if (errorStr.length < 100 && !errorStr.includes("{") && !errorStr.includes("[")) {
    return errorStr;
  }

  return "Ой, щось пішло не так";
}
