"use client";
import Canvas from "@/components/Canvas";
import Chat from "@/components/ChatTemp";
import { useParams, useSearchParams } from "next/navigation";
import { useRef, useState, useEffect } from "react";
// ye add karo imports ke neeche
interface StrokeData {
  type: string;
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  color: string;
  size: number;
}
export default function GamePage() {
  const params = useParams();
  const searchParams = useSearchParams();
  const [lastStroke, setLastStroke] = useState<StrokeData | null>(null);
  const [shouldClear, setShouldClear] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const [currentDrawer, setCurrentDrawer] = useState("");
  const [word, setWord] = useState(""); // sirf drawer ko dikhega
  const [hint, setHint] = useState(""); // guessers ko dikhega
  const [timer, setTimer] = useState(80);
  const [round, setRound] = useState(1);
  const [scores, setScores] = useState<{ [key: string]: number }>({});
  const [gameOver, setGameOver] = useState(false);
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
      if (data.type === "draw") {
        setLastStroke(data);
      }
      if (data.type === "clear") {
        setShouldClear(true);
      }
      if (data.type === "turn_start") {
        setCurrentDrawer(data.drawer);
        setRound(data.round);
        setWord(""); // word reset karo
        setHint(""); // hint reset karo
      }
      if (data.type === "your_word") {
        setWord(data.word); // sirf drawer ko milega
      }
      if (data.type === "word_hint") {
        setHint(data.hint); // guessers ko milega
      }
      if (data.type === "timer") {
        setTimer(data.seconds);
      }
      if (data.type === "correct_guess") {
        setMessages((prev) => [
          ...prev,
          {
            username: "🎉 System",
            message: `${data.username} ne sahi guess kar liya!`,
            type: "system",
          },
        ]);
      }
      if (data.type === "turn_end") {
        setMessages((prev) => [
          ...prev,
          {
            username: "⏰ System",
            message: `Word tha: ${data.word}`,
            type: "system",
          },
        ]);
      }
      if (data.type === "game_over") {
        setScores(data.scores);
        setGameOver(true);
      }
    };

    return () => ws.close();
  }, [roomId]);
  return (
    <main className="flex flex-col items-center justify-center gap-4 min-h-screen bg-zinc-950 text-white">
      {/* Game Info Bar */}
      {!gameOver && (
        <div className="flex gap-8 items-center bg-zinc-900 px-6 py-3 rounded-xl">
          <span>Round: {round}</span>
          <span
            className={`text-2xl font-bold ${timer <= 10 ? "text-red-500" : "text-white"}`}
          >
            {timer}s
          </span>
          <span>Drawing: {currentDrawer}</span>
          {/* word ya hint dikhao */}
          {word ? (
            <span className="text-green-400 font-bold">Word: {word}</span>
          ) : (
            <span className="text-zinc-400">{hint}</span>
          )}
        </div>
      )}

      {/* Game Over Screen */}
      {gameOver ? (
        <div className="bg-zinc-900 rounded-xl p-8 flex flex-col gap-4 items-center">
          <h1 className="text-3xl font-bold">🏆 Game Over!</h1>
          {Object.entries(scores)
            .sort(([, a], [, b]) => b - a) // score se sort karo
            .map(([player, score], i) => (
              <div key={i} className="flex gap-4 items-center">
                <span>{i === 0 ? "🥇" : i === 1 ? "🥈" : "🥉"}</span>
                <span>{player}</span>
                <span className="text-yellow-400">{score} pts</span>
              </div>
            ))}
          <button
            onClick={() => (window.location.href = "/")}
            className="px-6 py-2 bg-white text-black rounded-lg mt-4"
          >
            Play Again
          </button>
        </div>
      ) : (
        <div className="flex gap-4 items-center">
          <Canvas
            roomId={roomId}
            username={username}
            wsRef={wsRef}
            lastStroke={lastStroke}
            shouldClear={shouldClear}
            onClearDone={() => setShouldClear(false)}
          />
          <Chat messages={messages} onSend={sendChat} />
        </div>
      )}
    </main>
  );
}
