import { afterEach, describe, expect, it, vi } from "vitest";
import { createReview, listGames } from "@/lib/api";

afterEach(() => vi.restoreAllMocks());

describe("API client", () => {
  it("serializes game query parameters and validates the response", async () => {
    const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response(JSON.stringify({
      data: [{
        id: "game-1", title: "Hades", description: "A fast roguelike adventure.", genre: "Roguelike",
        platforms: ["PC"], developer: "Supergiant Games", releaseDate: "2020-09-17",
        averageRating: 4.8, reviewCount: 12, ratingDistribution: { "1": 0, "2": 0, "3": 1, "4": 3, "5": 8 },
      }],
      pagination: { page: 2, pageSize: 12, totalItems: 13, totalPages: 2 },
    }), { status: 200, headers: { "content-type": "application/json" } }));

    const result = await listGames({ q: "hades", genre: "Roguelike", platform: "PC", minRating: 4, sort: "rating_desc", page: 2, pageSize: 12 });
    expect(result.data[0]?.title).toBe("Hades");
    expect(fetchMock).toHaveBeenCalledWith(expect.stringContaining("q=hades"), expect.objectContaining({ cache: "no-store" }));
    expect(fetchMock).toHaveBeenCalledWith(expect.stringContaining("page=2"), expect.anything());
  });

  it("surfaces structured validation errors", async () => {
    vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response(JSON.stringify({
      error: { message: "The review contains invalid fields.", fields: { rating: "Rating must be between 1 and 5." } },
    }), { status: 422, headers: { "content-type": "application/json" } }));

    await expect(createReview("game-1", { reviewerName: "Alex", rating: 5, text: "A detailed and useful review." })).rejects.toMatchObject({
      status: 422,
      fields: { rating: "Rating must be between 1 and 5." },
    });
  });
});
