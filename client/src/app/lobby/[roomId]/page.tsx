"use client";
import { useParams, useSearchParams } from "next/navigation";
import { useState, useEffect, useRef } from "react";

export default function LobbyPage() {
  const params = useParams();
  const searchParams = useSearchParams();

  const roomId = params.roomId as string;
  const username = searchParams.get("username") || "Anonymous";
  const isHost = searchParams.get("host") === "true";

  const [players, setPlayers] = useState<string[]>([username]);
  const [drawTime, setDrawTime] = useState(80);
  const [rounds, setRounds] = useState(3);
  const wsRef = useRef<WebSocket | null>(null);
  useEffect(() => {
    fetch(`http://localhost:8080/room/players?room=${roomId}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.players) setPlayers(data.players);
      });
  }, [roomId]);

  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:8080/ws?room=${roomId}&username=${username}`,
    );
    wsRef.current = ws;

    ws.onopen = () => {
      ws.send(
        JSON.stringify({
          type: "player_joined",
          username: username,
        }),
      );
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.type === "player_joined") {
        setPlayers((prev) => {
          if (prev.includes(data.username)) return prev; // ← duplicate fix
          const updated = [...prev, data.username];
          // agar host hai toh updated list sabko bhejo
          if (isHost) {
            ws.send(
              JSON.stringify({
                type: "players_list",
                players: updated,
              }),
            );
          }
          return updated;
        });
      }

      // joiner ko existing players milenge
      if (data.type === "players_list") {
        setPlayers(data.players);
      }

      // game start hone pe game page pe jao
      if (data.type === "game_start") {
        window.location.href = `/game/${roomId}?username=${username}`;
      }
    };

    return () => ws.close();
  }, [roomId]);

  return (
    <main className="flex items-center justify-center min-h-screen bg-zinc-950 text-white">
      <div className="bg-zinc-900 rounded-xl p-8 w-96 flex flex-col gap-6">
        <h1 className="text-2xl font-bold">Room: {roomId}</h1>

        <button
          onClick={() =>
            navigator.clipboard.writeText(
              `${window.location.origin}/lobby/${roomId}?username=guest`,
            )
          }
          className="px-4 py-2 bg-zinc-700 rounded-lg text-sm"
        >
          🔗 Copy Invite Link
        </button>

        <div>
          <h2 className="text-zinc-400 text-sm mb-2">
            Players ({players.length}/8)
          </h2>
          {players.map((p, i) => (
            <div key={i} className="flex items-center gap-2 py-1">
              <span>👤</span>
              <span>{p}</span>
              {p === username && isHost && (
                <span className="text-xs text-yellow-400">host</span>
              )}
            </div>
          ))}
        </div>

        {isHost && (
          <div className="flex flex-col gap-3">
            <div className="flex justify-between items-center">
              <span>Draw Time</span>
              <select
                value={drawTime}
                onChange={(e) => setDrawTime(Number(e.target.value))}
                className="bg-zinc-800 rounded px-2 py-1"
              >
                <option value={30}>30s</option>
                <option value={60}>60s</option>
                <option value={80}>80s</option>
                <option value={120}>120s</option>
              </select>
            </div>
            <div className="flex justify-between items-center">
              <span>Rounds</span>
              <select
                value={rounds}
                onChange={(e) => setRounds(Number(e.target.value))}
                className="bg-zinc-800 rounded px-2 py-1"
              >
                <option value={2}>2</option>
                <option value={3}>3</option>
                <option value={5}>5</option>
              </select>
            </div>
          </div>
        )}

        {isHost && (
          <button
            onClick={() => {
              wsRef.current?.send(
                JSON.stringify({
                  type: "game_start",
                  drawTime: drawTime,
                  rounds: rounds,
                }),
              );
            }}
            className="px-6 py-3 bg-green-500 text-white rounded-xl font-bold text-lg"
          >
            Start Game! 🎮
          </button>
        )}
      </div>
    </main>
  );
}
