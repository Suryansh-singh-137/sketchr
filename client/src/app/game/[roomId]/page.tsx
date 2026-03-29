"use client";
import Canvas from "@/components/Canvas";
import { useParams, useSearchParams } from "next/navigation";

export default function GamePage() {
  const params = useParams();
  const searchParams = useSearchParams();

  const roomId = params.roomId as string;
  const username = searchParams.get("username") || "Anonymous";

  return (
    <main className="bg-zinc-950 min-h-screen">
      <Canvas roomId={roomId} username={username} />
    </main>
  );
}
