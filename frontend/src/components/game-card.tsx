import Link from "next/link";
import type { Game } from "@/lib/schemas";
import { GameCover } from "@/components/game-cover";
import { Rating } from "@/components/rating";

export function GameCard({ game }: { game: Game }) {
  return (
    <article className="game-card">
      <Link href={`/games/${game.id}`} aria-label={`Read reviews for ${game.title}`}>
        <GameCover title={game.title} genre={game.genre} />
        <div className="game-card-body">
          <div><p>{game.platforms.slice(0, 2).join(" · ")}</p><h3>{game.title}</h3></div>
          <Rating value={game.averageRating} count={game.reviewCount} compact />
        </div>
      </Link>
    </article>
  );
}

