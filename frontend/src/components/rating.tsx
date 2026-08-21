export function Rating({ value, count, compact = false }: { value: number; count?: number; compact?: boolean }) {
  return (
    <span className={`rating${compact ? " rating-compact" : ""}`} aria-label={`${value.toFixed(1)} out of 5${count === undefined ? "" : ` from ${count} reviews`}`}>
      <span aria-hidden="true">★</span>
      <strong>{value > 0 ? value.toFixed(1) : "New"}</strong>
      {count !== undefined && <small>{count} {count === 1 ? "review" : "reviews"}</small>}
    </span>
  );
}

