import React, { useState, useEffect, useRef } from "react";
import "./App.css";

export default function App() {
  const [lines, setLines] = useState([]);
  const [input, setInput] = useState("");
  const wsRef = useRef(null);
  const bottomRef = useRef(null);

  // Connect ke backend WebSocket
  useEffect(() => {
    wsRef.current = new WebSocket("ws://localhost:9090/ws");

    wsRef.current.onopen = () => {
      setLines((prev) => [...prev, "Connected to backend"]);
    };

    wsRef.current.onmessage = (event) => {
      setLines((prev) => [...prev, event.data]);
    };

    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, []);

  // Auto scroll ke bawah
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [lines]);

  // Submit command
  const handleSubmit = (e) => {
    e.preventDefault();
    if (!input) return;

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(input + "\n");
    }

    setLines((prev) => [...prev, `> ${input}`]);
    setInput("");
  };

  return (
    <div className="App">
      <h2 style={{ color: "#fff", textAlign: "center", marginTop: "10px" }}>
        Pure React Terminal
      </h2>
      <div className="terminal-container">
        {lines.map((line, idx) => (
          <div key={idx} className="terminal-line">
            {line}
          </div>
        ))}
        <div ref={bottomRef}></div>
      </div>
      <form onSubmit={handleSubmit} className="terminal-input-form">
        <span>&gt; </span>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          autoFocus
          className="terminal-input"
        />
      </form>
    </div>
  );
}
