"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { queryKeys } from "@/lib/query";
import { reviewCreatedEventSchema, type Review } from "@/lib/schemas";

export type ConnectionStatus = "connecting" | "live" | "offline";

const websocketURL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/v1/ws";

export function useReviewEvents(gameID: string) {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<ConnectionStatus>("connecting");
  const [latestReview, setLatestReview] = useState<Review | null>(null);

  useEffect(() => {
    let socket: WebSocket | undefined;
    let reconnectTimer: number | undefined;
    let active = true;

    function connect() {
      if (!active) return;
      setStatus("connecting");
      socket = new WebSocket(websocketURL);
      socket.addEventListener("open", () => setStatus("live"));
      socket.addEventListener("message", (message) => {
        try {
          const event = reviewCreatedEventSchema.parse(JSON.parse(String(message.data)));
          if (event.gameId !== gameID) return;
          setLatestReview(event.review);
          void queryClient.invalidateQueries({ queryKey: queryKeys.game(gameID) });
          void queryClient.invalidateQueries({ queryKey: queryKeys.reviewPrefix(gameID) });
        } catch {
          // Ignore unknown events; REST remains the source of truth.
        }
      });
      socket.addEventListener("close", () => {
        if (!active) return;
        setStatus("offline");
        reconnectTimer = window.setTimeout(connect, 1500);
      });
      socket.addEventListener("error", () => socket?.close());
    }

    connect();
    return () => {
      active = false;
      if (reconnectTimer !== undefined) window.clearTimeout(reconnectTimer);
      socket?.close(1000, "page closed");
    };
  }, [gameID, queryClient]);

  return { status, latestReview };
}

