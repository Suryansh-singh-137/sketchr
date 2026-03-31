"use client";
import { useEffect, useRef, useState } from "react";
interface CanvasProps {
  roomId: string;
  username: string;
  onChatMessage: (msg: ChatMessage) => void;
}
interface StrokeData {
  type: string;
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  color: string;
  size: number;
}
interface ChatMessage {
  username: string;
  message: string;
  type: "chat" | "system";
}
export default function Canvas({
  roomId,
  username,
  onChatMessage,
}: CanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [color, setColor] = useState("#000000");
  const [brushSize, setBrushSize] = useState(4);
  const lastPos = useRef({ x: 0, y: 0 });
  const wsRef = useRef<WebSocket | null>(null);
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws?room=${roomId}`);
    wsRef.current = ws;

    ws.onopen = () => console.log("WebSocket connected!");
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === "draw") {
        drawStroke(data);
      }
      if (data.type === "clear") {
        const canvas = canvasRef.current;
        if (!canvas) return;
        const ctx = canvas.getContext("2d");
        if (!ctx) return;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
      }
      if (data.type === "chat") {
        onChatMessage(data);
      }
    };

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [roomId]);

  const startDrawing = (e: React.MouseEvent<HTMLCanvasElement>) => {
    setIsDrawing(true);
    lastPos.current = { x: e.nativeEvent.offsetX, y: e.nativeEvent.offsetY };
  };

  const draw = (e: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!isDrawing || !canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    // Step 1: current position save karo PEHLE
    const currentX = e.nativeEvent.offsetX;
    const currentY = e.nativeEvent.offsetY;

    // Step 2: draw karo
    ctx.strokeStyle = color;
    ctx.lineWidth = brushSize;
    ctx.lineCap = "round";
    ctx.beginPath();
    ctx.moveTo(lastPos.current.x, lastPos.current.y);
    ctx.lineTo(currentX, currentY);
    ctx.stroke();

    // Step 3: send karo — lastPos update se PEHLE
    wsRef.current?.send(
      JSON.stringify({
        type: "draw",
        x1: lastPos.current.x, // ← abhi bhi purani position hai ✅
        y1: lastPos.current.y,
        x2: currentX,
        y2: currentY,
        color: color,
        size: brushSize,
      }),
    );

    // Step 4: ab lastPos update karo
    lastPos.current = { x: currentX, y: currentY };
  };

  const stopDrawing = () => {
    setIsDrawing(false);
  };

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    wsRef.current?.send(
      JSON.stringify({
        type: "clear",
      }),
    );
  };

  const drawStroke = (data: StrokeData) => {
    const canavs = canvasRef.current;
    const ctx = canavs?.getContext("2d");
    if (!ctx) return;
    ctx.strokeStyle = data.color;
    ctx.lineWidth = data.size;
    ctx.lineCap = "round";
    ctx.beginPath();
    ctx.moveTo(data.x1, data.y1);
    ctx.lineTo(data.x2, data.y2);
    ctx.stroke();
  };

  return (
    <div className="flex flex-col items-center gap-4 bg-zinc-950 min-h-screen justify-center">
      <canvas
        ref={canvasRef}
        width={800}
        height={500}
        onMouseDown={startDrawing}
        onMouseMove={draw}
        onMouseUp={stopDrawing}
        onMouseLeave={stopDrawing}
        className="bg-white rounded-xl cursor-crosshair"
      />
      <div className="flex gap-4 items-center">
        {/* Color picker */}
        <input
          type="color"
          value={color}
          onChange={(e) => setColor(e.target.value)}
        />
        {/* Brush size */}
        <input
          type="range"
          min={1}
          max={20}
          value={brushSize}
          onChange={(e) => setBrushSize(Number(e.target.value))}
        />
        {/* Clear button */}
        <button
          onClick={clearCanvas}
          className="px-4 py-2 bg-red-500 text-white rounded-lg"
        >
          Clear
        </button>
      </div>
    </div>
  );
}
