export const buildObservabilityQuery = (params: Record<string, unknown>): string => {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null) {
      return;
    }
    const serialized = `${value}`.trim();
    if (serialized === '') {
      return;
    }
    query.set(key, serialized);
  });
  return query.toString();
};
