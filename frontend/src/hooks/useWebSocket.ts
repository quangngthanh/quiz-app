import { useEffect, useRef, useState } from 'react';
import type { WSMessage } from '../types/quiz';

interface UseWebSocketOptions {
  onMessage?: (message: WSMessage) => void;
  onError?: (error: Event) => void;
  onClose?: (event: CloseEvent) => void;
}

export const useWebSocket = (url: string, options: UseWebSocketOptions = {}) => {
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const connect = () => {
      try {
        const ws = new WebSocket(url);
        wsRef.current = ws;

        ws.onopen = () => {
          setIsConnected(true);
          setError(null);
          console.log('WebSocket connected');
        };

        ws.onmessage = (event) => {
          try {
            const message: WSMessage = JSON.parse(event.data);
            options.onMessage?.(message);
          } catch (err) {
            console.error('Failed to parse WebSocket message:', err);
          }
        };

        ws.onerror = (event) => {
          setError('WebSocket error occurred');
          options.onError?.(event);
        };

        ws.onclose = (event) => {
          setIsConnected(false);
          options.onClose?.(event);
          
          // Reconnect after 3 seconds if not a normal closure
          if (event.code !== 1000) {
            setTimeout(connect, 3000);
          }
        };
      } catch (err) {
        setError('Failed to connect to WebSocket');
      }
    };

    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close(1000);
      }
    };
  }, [url]);

  const sendMessage = (message: any) => {
    if (wsRef.current && isConnected) {
      wsRef.current.send(JSON.stringify(message));
    }
  };

  return { isConnected, error, sendMessage };
};