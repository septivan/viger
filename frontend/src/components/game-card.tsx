import Link from "next/link";
import type { Game } from "@/lib/schemas";
import { GameCover } from "@/components/game-cover";
import { Rating } from "@/components/rating";

export function GameCard({ game }: { game: Game }) {
  const releaseYear = new Intl.DateTimeFormat("en", {
    year: "numeric",
    timeZone: "UTC",
  }).format(new Date(game.releaseDate));

  return (
    <article className="game-card">
      <Link href={`/games/${game.id}`} aria-label={`Read reviews for ${game.title}`}>
        <GameCover title={game.title} genre={game.genre} />
        <div className="game-card-body">
          <div className="game-card-copy">
            <h3>{game.title}</h3>
            <dl className="game-card-facts">
              <div><dt>Platforms</dt><dd>{game.platforms.slice(0, 2).join(", ")}</dd></div>
              <div><dt>Released</dt><dd><time dateTime={game.releaseDate}>{releaseYear}</time></dd></div>
            </dl>
          </div>
          <Rating value={game.averageRating} count={game.reviewCount} compact />
        </div>
      </Link>
    </article>
  );
}
