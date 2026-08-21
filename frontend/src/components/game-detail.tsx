"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { GameCover } from "@/components/game-cover";
import { Rating } from "@/components/rating";
import { ReviewForm } from "@/components/review-form";
import { ReviewList } from "@/components/review-list";
import { useReviewEvents } from "@/hooks/use-review-events";
import { getGame } from "@/lib/api";
import { queryKeys } from "@/lib/query";

export function GameDetail({ gameID }: { gameID: string }) {
  const game = useQuery({ queryKey: queryKeys.game(gameID), queryFn: () => getGame(gameID) });
  const live = useReviewEvents(gameID);

  if (game.isPending) return <main className="detail-page"><div className="detail-skeleton skeleton" /></main>;
  if (game.isError) return <main className="detail-page"><div className="state-panel"><span>Game unavailable</span><h1>We could not open this game.</h1><p>{game.error instanceof Error ? game.error.message : "Please try again."}</p><Link href="/">Return to the collection</Link></div></main>;

  const item = game.data;
  const maximumDistribution = Math.max(1, ...Object.values(item.ratingDistribution));
  return (
    <main className="detail-page">
      <Link className="back-link" href="/"><span aria-hidden="true">←</span> Back to all games</Link>
      <section className="game-hero">
        <GameCover genre={item.genre} large title={item.title} />
        <div className="game-hero-content">
          <div className="game-hero-topline"><span>{item.genre}</span><span className={`connection-status ${live.status}`}><i />{live.status === "live" ? "Live reviews" : live.status === "connecting" ? "Connecting" : "Reconnecting"}</span></div>
          <h1>{item.title}</h1>
          <p className="game-description">{item.description}</p>
          <dl className="game-facts"><div><dt>Developer</dt><dd>{item.developer}</dd></div><div><dt>Released</dt><dd>{new Intl.DateTimeFormat("en", { year: "numeric", month: "short", day: "numeric", timeZone: "UTC" }).format(new Date(item.releaseDate))}</dd></div><div><dt>Platforms</dt><dd>{item.platforms.join(", ")}</dd></div></dl>
          <div className="rating-summary"><Rating count={item.reviewCount} value={item.averageRating} /><div className="rating-bars">{[5, 4, 3, 2, 1].map((rating) => <div key={rating}><span>{rating}</span><i><b style={{ width: `${(item.ratingDistribution[String(rating)] ?? 0) / maximumDistribution * 100}%` }} /></i><small>{item.ratingDistribution[String(rating)] ?? 0}</small></div>)}</div></div>
          {live.latestReview && <p className="live-notice" role="status"><span aria-hidden="true">✦</span> A new review from {live.latestReview.reviewerName} just arrived.</p>}
        </div>
      </section>
      <div className="review-layout"><ReviewList gameID={gameID} /><ReviewForm gameID={gameID} gameTitle={item.title} /></div>
    </main>
  );
}

