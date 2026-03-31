"use client";
import Canvas from "@/components/Canvas";
import Chat from "@/components/ChatTemp";
import { useParams, useSearchParams } from "next/navigation";
import { useRef, useState, useEffect } from "react";
export default function GamePage() {
  const params = useParams();
  const searchParams = useSearchParams();
  const wsRef = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<
    { username: string; message: string; type: "chat" | "system" }[]
  >([]);
  const roomId = params.roomId as string;
  const username = searchParams.get("username") || "Anonymous";
  const sendChat = (message: string) => {
    wsRef.current?.send(
      JSON.stringify({
        type: "chat",
        username: username,
        message: message,
      }),
    );
  };
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws?room=${roomId}`);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === "chat") {
        setMessages((prev) => [...prev, data]);
      }
    };

    return () => ws.close();
  }, [roomId]);
  return (
    <main className="flex items-center justify-center gap-4 min-h-screen bg-zinc-950">
      <Canvas roomId={roomId} username={username} />
      <Chat messages={messages} onSend={sendChat} />
    </main>
  );
}
