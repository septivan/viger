import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { GameCard } from "@/components/game-card";

describe("GameCard", () => {
  it("shows values that explain platform and release ordering results", () => {
    render(<GameCard game={{
      id: "signal-drift",
      title: "Signal Drift",
      description: "A test game description.",
      genre: "Adventure",
      platforms: ["PC", "PlayStation 5", "Mobile"],
      developer: "Northstar Works",
      releaseDate: "2025-04-18T00:00:00Z",
      averageRating: 4.4,
      reviewCount: 12,
      ratingDistribution: { "1": 0, "2": 0, "3": 2, "4": 3, "5": 7 },
    }} />);

    expect(screen.getByText("PC · PlayStation 5 · Mobile")).toBeInTheDocument();
    expect(screen.getByText("2025")).toBeInTheDocument();
    expect(screen.queryByText("Platforms")).not.toBeInTheDocument();
    expect(screen.queryByText("Released")).not.toBeInTheDocument();
  });
});
