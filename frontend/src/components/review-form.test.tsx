import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ReviewForm } from "@/components/review-form";

const createReview = vi.fn();
vi.mock("@/lib/api", async (importOriginal) => {
  const original = await importOriginal<typeof import("@/lib/api")>();
  return { ...original, createReview: (...arguments_: unknown[]) => createReview(...arguments_) };
});

function wrapper({ children }: { children: ReactNode }) {
  return <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>{children}</QueryClientProvider>;
}

describe("ReviewForm", () => {
  beforeEach(() => createReview.mockReset());

  it("shows accessible client validation", async () => {
    render(<ReviewForm gameID="game-1" gameTitle="Hades" />, { wrapper });
    await userEvent.click(screen.getByRole("button", { name: /publish review/i }));
    expect(await screen.findByText("Enter at least 2 characters.")).toBeVisible();
    expect(screen.getByText("Choose a rating.")).toBeVisible();
    expect(screen.getByText("Enter at least 10 characters.")).toBeVisible();
    expect(createReview).not.toHaveBeenCalled();
  });

  it("orders rating choices from one to five from left to right", () => {
    render(<ReviewForm gameID="game-1" gameTitle="Hades" />, { wrapper });
    expect(screen.getAllByRole("radio").map((radio) => radio.getAttribute("value"))).toEqual(["1", "2", "3", "4", "5"]);
  });

  it("submits a valid review and reports success", async () => {
    createReview.mockResolvedValue({ id: "review-1" });
    render(<ReviewForm gameID="game-1" gameTitle="Hades" />, { wrapper });
    await userEvent.click(screen.getByRole("radio", { name: "5 out of 5" }));
    await userEvent.type(screen.getByLabelText("Your name"), "Alex");
    await userEvent.type(screen.getByLabelText("Your review"), "A focused and beautifully paced experience.");
    await userEvent.click(screen.getByRole("button", { name: /publish review/i }));
    await waitFor(() => expect(createReview).toHaveBeenCalledWith("game-1", { reviewerName: "Alex", rating: 5, text: "A focused and beautifully paced experience." }));
    expect(await screen.findByRole("status")).toHaveTextContent("now part of the conversation");
  });
});
