import { useState, useEffect, useRef } from 'react';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [darkMode, setDarkMode] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [clientId] = useState(`User-${Math.floor(Math.random() * 1000)}`);
  const ws = useRef(null);
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    const connectWebSocket = () => {
      ws.current = new WebSocket('ws://localhost:80/chat');

      ws.current.onopen = () => {
        console.log('Connected to WebSocket server');
        setIsConnected(true);
      };

      ws.current.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          // Only add messages that aren't from this client
          if (data.sender !== clientId) {
            setMessages(prev => [...prev, {
              id: Date.now(),
              text: data.content,
              sender: data.sender,
              type: data.sender === 'System' ? 'system' : 'received',
              timestamp: data.timestamp || new Date().toLocaleTimeString()
            }]);
          }
        } catch (e) {
          console.error('Error parsing message:', e);
        }
      };

      ws.current.onclose = () => {
        console.log('WebSocket connection closed');
        setIsConnected(false);
        setTimeout(connectWebSocket, 5000);
      };

      ws.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        ws.current.close();
      };
    };

    connectWebSocket();
    return () => ws.current?.close();
  }, [clientId]);

  const sendMessage = (e) => {
    e.preventDefault();
    if (input.trim() && ws.current) {
      const message = {
        content: input,
        sender: clientId,
        timestamp: new Date().toLocaleTimeString()
      };

      // Add message locally first
      setMessages(prev => [...prev, {
        id: Date.now(),
        text: input,
        sender: clientId,
        type: 'sent',
        timestamp: new Date().toLocaleTimeString()
      }]);

      // Send to server
      ws.current.send(JSON.stringify(message));
      setInput('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage(e);
    }
  };

  const toggleDarkMode = () => {
    setDarkMode(prev => !prev);
  };

  return (
    <div className={`chat-app ${darkMode ? 'dark' : ''}`}>
      <div className="chat-container">
        <header className="chat-header">
          <h1>Chat Room</h1>
          <div className="header-controls">
            <span className={`connection-status ${isConnected ? 'connected' : 'disconnected'}`}>
              {isConnected ? '🟢 Connected' : '🔴 Disconnecting...'}
            </span>
            <button onClick={toggleDarkMode} className="theme-toggle">
              {darkMode ? '☀️ Light' : '🌙 Dark'}
            </button>
          </div>
        </header>

        <div className="messages-area">
          {messages.map((message) => (
            <div key={message.id} className={`message ${message.type}`}>
              <div className="message-bubble">
                {message.sender && <span className="sender">{message.sender}</span>}
                <p>{message.text}</p>
                <span className="timestamp">{message.timestamp}</span>
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>

        <form onSubmit={sendMessage} className="input-area">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type a message..."
            className="message-input"
            disabled={!isConnected}
          />
          <button
            type="submit"
            disabled={!input.trim() || !isConnected}
            className="send-button"
          >
            Send
          </button>
        </form>
      </div>
    </div>
  );
}

export default App;