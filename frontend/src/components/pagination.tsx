type PaginationProps = {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  label: string;
};

export function Pagination({ page, totalPages, onPageChange, label }: PaginationProps) {
  if (totalPages <= 1) return null;
  const start = Math.max(1, Math.min(page - 2, totalPages - 4));
  const pages = Array.from({ length: Math.min(5, totalPages) }, (_, index) => start + index);
  return (
    <nav className="pagination" aria-label={label}>
      <button disabled={page === 1} onClick={() => onPageChange(page - 1)} type="button">Previous</button>
      <div className="page-numbers">
        {pages.map((value) => <button aria-current={value === page ? "page" : undefined} className={value === page ? "active" : ""} key={value} onClick={() => onPageChange(value)} type="button">{value}</button>)}
      </div>
      <button disabled={page === totalPages} onClick={() => onPageChange(page + 1)} type="button">Next</button>
    </nav>
  );
}

