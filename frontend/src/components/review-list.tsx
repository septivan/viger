"use client";

import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { Pagination } from "@/components/pagination";
import { listReviews, type ReviewQuery } from "@/lib/api";
import { queryKeys } from "@/lib/query";

export function ReviewList({ gameID }: { gameID: string }) {
  const [page, setPage] = useState(1);
  const [sort, setSort] = useState("newest");
  const query: ReviewQuery = { page, pageSize: 10, sort };
  const reviews = useQuery({ queryKey: queryKeys.reviews(gameID, query), queryFn: () => listReviews(gameID, query) });

  return (
    <section className="reviews-panel" aria-labelledby="reviews-title">
      <div className="reviews-heading"><div><p className="eyebrow">Player notes</p><h2 id="reviews-title">The conversation</h2></div><label><span className="sr-only">Sort reviews</span><select onChange={(event) => { setSort(event.target.value); setPage(1); }} value={sort}><option value="newest">Newest first</option><option value="oldest">Oldest first</option><option value="rating_desc">Highest rated</option><option value="rating_asc">Lowest rated</option></select></label></div>
      {reviews.isPending && <div className="review-stack" aria-label="Loading reviews">{Array.from({ length: 3 }, (_, index) => <div className="skeleton review-skeleton" key={index} />)}</div>}
      {reviews.isError && <div className="inline-state" role="alert"><p>The reviews could not be loaded.</p><button onClick={() => reviews.refetch()} type="button">Try again</button></div>}
      {reviews.data?.data.length === 0 && <div className="inline-state"><p>No reviews yet. Be the first to start the conversation.</p></div>}
      {reviews.data && reviews.data.data.length > 0 && <div className="review-stack">{reviews.data.data.map((review) => <article className="review-card" key={review.id}><div className="review-meta"><span className="review-avatar" aria-hidden="true">{review.reviewerName.slice(0, 1).toUpperCase()}</span><div><strong>{review.reviewerName}</strong><time dateTime={review.createdAt}>{new Intl.DateTimeFormat("en", { dateStyle: "medium" }).format(new Date(review.createdAt))}</time></div><span className="review-rating" aria-label={`${review.rating} out of 5`}><span aria-hidden="true">★</span> {review.rating}.0</span></div><p>{review.text}</p></article>)}</div>}
      {reviews.data && <Pagination label="Review pages" onPageChange={setPage} page={reviews.data.pagination.page} totalPages={reviews.data.pagination.totalPages} />}
    </section>
  );
}

