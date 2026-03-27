"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function Home() {
  const [username, setUsername] = useState("");
  const [roomCode, setRoomCode] = useState("");
  const router = useRouter();

  const handleCreateRoom = async () => {
    const response = await fetch("http://localhost:8080/room/create", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
    const data = await response.json();
    const roomID = data.roomId;
    router.push(`/game/${roomID}?username=${username}`);
  };

  const handleJoinRoom = () => {
    // router.push to game page with roomCode and username
    if (!roomCode || !username) return;
    router.push(`/game/${roomCode}?username=${username}`);
  };

  return (
    <main className="flex flex-col items-center justify-center min-h-screen bg-zinc-950 text-white gap-6">
      <h1 className="text-5xl font-bold tracking-tight">Sketchr</h1>
      <p className="text-zinc-400">Draw. Guess. Win.</p>

      <input
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Enter your username..."
        className="px-4 py-2 rounded-lg bg-zinc-800 text-white w-72 outline-none"
      />

      <button
        onClick={handleCreateRoom}
        disabled={!username}
        className="px-6 py-2 bg-white text-black rounded-lg font-semibold disabled:opacity-40"
      >
        Create Room
      </button>

      <div className="flex gap-2">
        <input
          value={roomCode}
          onChange={(e) => setRoomCode(e.target.value)}
          placeholder="Room code..."
          className="px-4 py-2 rounded-lg bg-zinc-800 text-white outline-none"
        />
        <button
          onClick={handleJoinRoom}
          disabled={!username || !roomCode}
          className="px-4 py-2 bg-zinc-700 rounded-lg disabled:opacity-40"
        >
          Join
        </button>
      </div>
    </main>
  );
}
