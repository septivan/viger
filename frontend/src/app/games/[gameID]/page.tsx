import type { Metadata } from "next";
import { GameDetail } from "@/components/game-detail";

export const metadata: Metadata = { title: "Game reviews" };

export default async function GamePage({ params }: { params: Promise<{ gameID: string }> }) {
  const { gameID } = await params;
  return <GameDetail gameID={gameID} />;
}

