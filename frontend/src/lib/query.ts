import type { GameQuery, ReviewQuery } from "@/lib/api";

export const queryKeys = {
  games: (query: GameQuery) => ["games", query] as const,
  game: (gameID: string) => ["game", gameID] as const,
  reviews: (gameID: string, query: ReviewQuery) => ["reviews", gameID, query] as const,
  reviewPrefix: (gameID: string) => ["reviews", gameID] as const,
};

