import type { Game, Pagination, Review, ReviewFormValues } from "@/lib/schemas";
import { gamePageSchema, gameResponseSchema, reviewPageSchema, reviewResponseSchema } from "@/lib/schemas";

const apiURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type GameQuery = {
  q?: string;
  genre?: string;
  platform?: string;
  minRating?: number;
  sort?: string;
  page: number;
  pageSize: number;
};

export type ReviewQuery = { sort: string; page: number; pageSize: number };
export type Page<T> = { data: T[]; pagination: Pagination };

export class APIError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly fields: Record<string, string> = {},
  ) {
    super(message);
    this.name = "APIError";
  }
}

async function request(path: string, init?: RequestInit): Promise<unknown> {
  const response = await fetch(`${apiURL}${path}`, {
    ...init,
    headers: { accept: "application/json", ...init?.headers },
    cache: "no-store",
  });
  const body = await response.json().catch(() => null) as { error?: { message?: string; fields?: Record<string, string> } } | null;
  if (!response.ok) {
    throw new APIError(body?.error?.message ?? "Viger could not complete the request.", response.status, body?.error?.fields);
  }
  return body;
}

export async function listGames(query: GameQuery): Promise<Page<Game>> {
  const parameters = new URLSearchParams();
  if (query.q) parameters.set("q", query.q);
  if (query.genre) parameters.set("genre", query.genre);
  if (query.platform) parameters.set("platform", query.platform);
  if (query.minRating) parameters.set("minRating", String(query.minRating));
  if (query.sort) parameters.set("sort", query.sort);
  parameters.set("page", String(query.page));
  parameters.set("pageSize", String(query.pageSize));
  return gamePageSchema.parse(await request(`/v1/games?${parameters}`));
}

export async function getGame(gameID: string): Promise<Game> {
  return gameResponseSchema.parse(await request(`/v1/games/${encodeURIComponent(gameID)}`)).data;
}

export async function listReviews(gameID: string, query: ReviewQuery): Promise<Page<Review>> {
  const parameters = new URLSearchParams({ sort: query.sort, page: String(query.page), pageSize: String(query.pageSize) });
  return reviewPageSchema.parse(await request(`/v1/games/${encodeURIComponent(gameID)}/reviews?${parameters}`));
}

export async function createReview(gameID: string, input: ReviewFormValues): Promise<Review> {
  return reviewResponseSchema.parse(await request(`/v1/games/${encodeURIComponent(gameID)}/reviews`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(input),
  })).data;
}

