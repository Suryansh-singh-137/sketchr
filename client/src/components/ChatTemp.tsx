"use client";
import { useState } from "react";

interface ChatMessage {
  username: string;
  message: string;
  type: "chat" | "system";
}

interface ChatProps {
  messages: ChatMessage[];
  onSend: (message: string) => void;
}

export default function ChatTemp({ messages, onSend }: ChatProps) {
  const [input, setInput] = useState("");

  const handleSend = () => {
    // agar input empty hai toh return
    if (!input.trim()) return;

    // onSend(input) call karo
    onSend(input);
    // input clear karo
    setInput("");
  };

  return (
    <div className="flex flex-col h-[500px] w-72 bg-zinc-900 rounded-xl">
      {/* Messages list */}
      <div className="flex-1 overflow-y-auto p-3 flex flex-col gap-2">
        {messages.map((msg, i) => (
          <div key={i} className="text-sm">
            <span className="text-zinc-400 font-semibold">
              {msg.username}:{" "}
            </span>
            <span className="text-white">{msg.message}</span>
          </div>
        ))}
      </div>

      {/* Input area */}
      <div className="flex gap-2 p-3 border-t border-zinc-700">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSend()}
          placeholder="Type a guess..."
          className="flex-1 bg-zinc-800 text-white px-3 py-2 rounded-lg outline-none"
        />
        <button
          onClick={handleSend}
          className="px-3 py-2 bg-white text-black rounded-lg font-semibold"
        >
          Send
        </button>
      </div>
    </div>
  );
}
