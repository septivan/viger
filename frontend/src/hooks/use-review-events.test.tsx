import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, renderHook, waitFor } from "@testing-library/react";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { useReviewEvents } from "@/hooks/use-review-events";

class MockWebSocket {
  static instance: MockWebSocket;
  listeners = new Map<string, Array<(event: MessageEvent | Event) => void>>();
  constructor(public readonly url: string) { MockWebSocket.instance = this; }
  addEventListener(name: string, listener: (event: MessageEvent | Event) => void) {
    this.listeners.set(name, [...(this.listeners.get(name) ?? []), listener]);
  }
  close() { this.emit("close", new Event("close")); }
  emit(name: string, event: MessageEvent | Event) { for (const listener of this.listeners.get(name) ?? []) listener(event); }
}

describe("useReviewEvents", () => {
  beforeEach(() => { vi.stubGlobal("WebSocket", MockWebSocket); });

  it("tracks connection status and invalidates matching game data", async () => {
    const queryClient = new QueryClient();
    const invalidate = vi.spyOn(queryClient, "invalidateQueries").mockResolvedValue();
    const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
    const { result, unmount } = renderHook(() => useReviewEvents("game-1"), { wrapper });
    act(() => MockWebSocket.instance.emit("open", new Event("open")));
    expect(result.current.status).toBe("live");

    act(() => MockWebSocket.instance.emit("message", new MessageEvent("message", { data: JSON.stringify({
      type: "review.created", gameId: "game-1",
      review: { id: "review-live", gameId: "game-1", reviewerName: "Sam", rating: 5, text: "A thoughtful and detailed live review.", createdAt: "2026-08-21T12:00:00Z" },
    }) })));
    await waitFor(() => expect(result.current.latestReview?.id).toBe("review-live"));
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["game", "game-1"] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ["reviews", "game-1"] });
    unmount();
  });
});

