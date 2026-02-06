import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { roomService } from '@/services/roomService';
import { codeExecutionService } from '@/services/codeExecutionService';
import { webrtcService } from '@/services/webrtcService';
import { useAuthStore } from '@/store/authStore';
import type { Room } from '@/types';
import { 
  Users, 
  Loader2,
  Send,
  LogOut,
  MessageCircle,
  Code,
  Video,
  Mic,
  MicOff,
  VideoOff,
  Play,
  Terminal
} from 'lucide-react';

interface Message {
  id: string;
  user: string;
  content: string;
  timestamp: string;
}

interface Participant {
  id: string;
  username: string;
  cursor_position?: number;
}

export default function RoomDetail() {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();
  const { accessToken, user, ensureFreshToken } = useAuthStore();
  
  const [room, setRoom] = useState<Room | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  
  // WebSocket
  const wsRef = useRef<WebSocket | null>(null);
  const connectionIdRef = useRef<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  
  // Chat
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const processedMessageIds = useRef<Set<string>>(new Set());
  
  // Code Editor
  const [code, setCode] = useState('// Start coding here...\n\n');
  const [language, setLanguage] = useState('javascript');
  const [output, setOutput] = useState('');
  const [isExecuting, setIsExecuting] = useState(false);
  
  // Participants
  const [participants, setParticipants] = useState<Participant[]>([]);
  
  // Video
  const [isMicOn, setIsMicOn] = useState(false);
  const [isVideoOn, setIsVideoOn] = useState(false);
  const [localStream, setLocalStream] = useState<MediaStream | null>(null);
  const [remoteStreams, setRemoteStreams] = useState<Map<string, MediaStream>>(new Map());
  const localVideoRef = useRef<HTMLVideoElement>(null);
  const [isVideoCallActive, setIsVideoCallActive] = useState(false);

  useEffect(() => {
    if (!roomId) return;
    fetchRoom();
  }, [roomId]);

  useEffect(() => {
    if (!roomId) return;
    
    // Generate connection ID for this connection attempt
    const connectionId = `conn-${Date.now()}-${Math.random()}`;
    connectionIdRef.current = connectionId;
    
    // Close any existing connection before creating a new one
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    
    connectWebSocket();
    
    return () => {
      if (wsRef.current && connectionIdRef.current === connectionId) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [roomId]); // Only reconnect when roomId changes, not when room data updates

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  useEffect(() => {
    // Initialize WebRTC
    webrtcService.initialize({
      onRemoteStream: (userId: string, stream: MediaStream) => {
        setRemoteStreams(prev => new Map(prev).set(userId, stream));
      },
      onRemoteStreamRemoved: (userId: string) => {
        setRemoteStreams(prev => {
          const updated = new Map(prev);
          updated.delete(userId);
          return updated;
        });
      },
      sendSignal: (signal: any) => {
        if (wsRef.current && isConnected) {
          wsRef.current.send(JSON.stringify({
            type: 'webrtc-signal',
            data: signal,
          }));
        }
      },
    });

    return () => {
      webrtcService.cleanup();
    };
  }, [isConnected]);

  useEffect(() => {
    if (localVideoRef.current && localStream) {
      localVideoRef.current.srcObject = localStream;
    }
  }, [localStream]);

  const fetchRoom = async () => {
    if (!roomId) return;
    
    try {
      setIsLoading(true);
      const data = await roomService.getRoom(roomId);
      // Handle nested room property from backend
      const roomData = (data as any).room || data;
      setRoom(roomData);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load room');
    } finally {
      setIsLoading(false);
    }
  };

  const connectWebSocket = async () => {
    if (!roomId) return;

    try {
      // Ensure we have a fresh token before connecting
      const freshToken = await ensureFreshToken();
      if (!freshToken) {
        setError('Authentication failed. Please log in again.');
        navigate('/login');
        return;
      }

      const wsUrl = roomService.getWebSocketUrl(roomId, freshToken);
      const socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        setIsConnected(true);
        
        // Send join message - backend will send user_list with all participants
        socket.send(JSON.stringify({
          type: 'join',
          user: user?.username || 'Anonymous',
        }));
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          switch (data.type) {
            case 'chat':
              const chatData = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              const messageContent = chatData?.message || '';
              
              // Create a unique identifier for this message
              const messageId = data.Timestamp 
                ? `${data.UserID}-${data.Timestamp}` 
                : `${data.UserID}-${data.Username}-${messageContent}-${Date.now()}`;
              
              // Skip if we've already processed this message
              if (processedMessageIds.current.has(messageId)) {
                break;
              }
              
              processedMessageIds.current.add(messageId);
              
              // Clean up old message IDs (keep only last 100)
              if (processedMessageIds.current.size > 100) {
                const idsArray = Array.from(processedMessageIds.current);
                processedMessageIds.current = new Set(idsArray.slice(-100));
              }
              
              setMessages(prev => [...prev, {
                id: `msg-${Date.now()}-${Math.random()}`,
                user: data.Username || 'Unknown',
                content: messageContent,
                timestamp: new Date().toISOString(),
              }]);
              break;
              
            case 'code_update':
              // Ignore code updates to allow free typing
              break;

            case 'rtc_offer':
              const offerData = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              handleWebRTCSignal(data.UserID, { type: 'offer', offer: offerData });
              break;

            case 'rtc_answer':
              const answerData = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              handleWebRTCSignal(data.UserID, { type: 'answer', answer: answerData });
              break;

            case 'rtc_candidate':
              const candidateData = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              handleWebRTCSignal(data.UserID, { type: 'ice-candidate', candidate: candidateData });
              break;
              
            case 'user_list':
              if (Array.isArray(data.Data)) {
                const participantList = data.Data.map((p: any) => ({
                  id: p.UserID || p.user_id || p.id,
                  username: p.Username || p.username || p.email,
                }));
                // Deduplicate by user ID
                const uniqueParticipants = participantList.reduce((acc: Participant[], current: Participant) => {
                  if (!acc.find(p => p.id === current.id)) {
                    acc.push(current);
                  }
                  return acc;
                }, []);
                setParticipants(uniqueParticipants);
              }
              break;
              
            case 'participants':
              setParticipants(data.Data || data.participants || []);
              break;
              
            case 'user_joined':
              const joinedUsername = data.Username || data.Data?.username || 'Someone';
              setMessages(prev => [...prev, {
                id: `join-${Date.now()}-${Math.random()}`,
                user: 'System',
                content: `${joinedUsername} joined the room`,
                timestamp: new Date().toISOString(),
              }]);
              
              // If video call is active, create peer connection for new user
              if (isVideoCallActive && data.UserID && data.UserID !== user?.id) {
                setTimeout(() => {
                  webrtcService.createPeerConnection(data.UserID, true);
                }, 1000);
              }
              break;
              
            case 'user_left':
              const leftUsername = data.Username || data.Data?.username || 'Someone';
              setMessages(prev => [...prev, {
                id: `leave-${Date.now()}-${Math.random()}`,
                user: 'System',
                content: `${leftUsername} left the room`,
                timestamp: new Date().toISOString(),
              }]);
              break;
              
            case 'join':
            case 'get_participants':
              // Echo messages from backend, ignore
              break;
              
            default:
              break;
          }
        } catch (err) {
          // Silent error handling
        }
      };

      socket.onerror = () => {
        setError('WebSocket connection error');
      };

      socket.onclose = () => {
        setIsConnected(false);
      };

      wsRef.current = socket;
    } catch (err) {
      setError('Failed to establish connection. Please try again.');
    }
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !wsRef.current || !isConnected) return;

    wsRef.current.send(JSON.stringify({
      type: 'chat',
      data: {
        message: newMessage,
      },
    }));

    setNewMessage('');
  };

  const handleCodeChange = (newCode: string) => {
    setCode(newCode);
  };

  const handleWebRTCSignal = async (userId: string, signal: any) => {
    try {
      switch (signal.type) {
        case 'offer':
          // Auto-start video call when receiving an offer
          if (!isVideoCallActive && !localStream) {
            try {
              const stream = await webrtcService.getLocalStream(true, true);
              setLocalStream(stream);
              setIsVideoCallActive(true);
              setIsMicOn(true);
              setIsVideoOn(true);
              
              // Handle the offer after getting local stream
              await webrtcService.handleOffer(userId, signal.offer);
            } catch (error) {
              return;
            }
          } else {
            await webrtcService.handleOffer(userId, signal.offer);
          }
          break;
        case 'answer':
          await webrtcService.handleAnswer(userId, signal.answer);
          break;
        case 'ice-candidate':
          await webrtcService.handleIceCandidate(userId, signal.candidate);
          break;
      }
    } catch (error) {
      // Ignore WebRTC errors
    }
  };

  const handleStartVideoCall = async () => {
    try {
      // Start with both audio and video enabled
      const stream = await webrtcService.getLocalStream(true, true);
      setLocalStream(stream);
      setIsVideoCallActive(true);
      setIsMicOn(true);
      setIsVideoOn(true);

      // Wait a bit for the state to update, then create peer connections
      setTimeout(() => {
        participants.forEach(p => {
          if (p.id && p.id !== user?.id) {
            webrtcService.createPeerConnection(p.id, true);
          }
        });
      }, 500);
    } catch (error: any) {
      setError('Failed to start video call. Please allow camera and microphone access.');
    }
  };

  const handleStopVideoCall = () => {
    webrtcService.cleanup();
    setLocalStream(null);
    setRemoteStreams(new Map());
    setIsVideoCallActive(false);
    setIsMicOn(false);
    setIsVideoOn(false);
  };

  const handleToggleMic = () => {
    const newState = !isMicOn;
    setIsMicOn(newState);
    webrtcService.toggleAudio(newState);
  };

  const handleToggleVideo = () => {
    const newState = !isVideoOn;
    setIsVideoOn(newState);
    webrtcService.toggleVideo(newState);
  };

  const handleRunCode = async () => {
    if (!code.trim()) {
      setOutput('Error: No code to execute');
      return;
    }

    setIsExecuting(true);
    setOutput('Executing...');

    try {
      const result = await codeExecutionService.executeCode(language, code);
      
      if (result.stderr && !result.stdout) {
        setOutput(`Error:\n${result.stderr}`);
      } else if (result.output) {
        setOutput(result.output);
      } else if (result.stdout) {
        setOutput(result.stdout);
      } else {
        setOutput('Code executed successfully with no output.');
      }
    } catch (error: any) {
      setOutput(`Execution Error:\n${error.message}`);
    } finally {
      setIsExecuting(false);
    }
  };

  const handleLeaveRoom = async () => {
    if (wsRef.current) {
      wsRef.current.send(JSON.stringify({
        type: 'leave',
        user: user?.username || 'Anonymous',
      }));
      wsRef.current.close();
      wsRef.current = null;
    }
    
    if (roomId) {
      await roomService.leaveRoom(roomId);
    }
    
    navigate('/rooms');
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-dojo-black-900 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
      </div>
    );
  }

  if (!room) {
    return (
      <div className="min-h-screen bg-dojo-black-900">
        <Navbar />
        <div className="container mx-auto px-4 py-8">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-gray-400">Room not found</p>
              <Button onClick={() => navigate('/rooms')} className="mt-4">
                Back to Rooms
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-dojo-black-900">
      <Navbar />
      
      <div className="container mx-auto px-4 py-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-white mb-1">
              {room.name}
            </h1>
            <div className="flex items-center gap-4 text-sm text-gray-400">
              <div className="flex items-center gap-1">
                <Users className="h-4 w-4" />
                <span>{room.current_participants || participants.length} / {room.max_participants} participants</span>
              </div>
              <div className="flex items-center gap-1">
                <div className={`h-2 w-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                <span>{isConnected ? 'Connected' : 'Disconnected'}</span>
              </div>
              <div className="flex items-center gap-2 px-3 py-1 bg-dojo-black-800 rounded border border-gray-700">
                <span className="text-gray-500">Room Code:</span>
                <span className="font-mono font-bold text-dojo-red-400">
                  {room.room_code && room.room_code !== 'N/A' ? room.room_code : (
                    <span className="text-gray-500 text-xs" title={JSON.stringify(room)}>Loading...</span>
                  )}
                </span>
                <button
                  onClick={async () => {
                    if (room.room_code && room.room_code !== 'N/A') {
                      try {
                        await navigator.clipboard.writeText(room.room_code);
                        const btn = document.activeElement as HTMLButtonElement;
                        const originalText = btn.textContent;
                        btn.textContent = '✓';
                        setTimeout(() => {
                          btn.textContent = originalText || '📋';
                        }, 2000);
                      } catch (err) {
                        alert('Failed to copy room code');
                      }
                    } else {
                      alert('Room code not available');
                    }
                  }}
                  className="text-gray-400 hover:text-white transition-colors"
                  title="Copy room code"
                >
                  📋
                </button>
              </div>
            </div>
          </div>
          
          <Button variant="outline" onClick={handleLeaveRoom}>
            <LogOut className="mr-2 h-4 w-4" />
            Leave Room
          </Button>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
            {error}
          </div>
        )}

        {/* Main Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-4 h-[calc(100vh-200px)]">
          {/* Code Editor */}
          <div className="lg:col-span-3 flex flex-col gap-4">
            <Card className="flex flex-col" style={{ height: 'calc(100% - 13rem)' }}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Code className="h-5 w-5 text-dojo-red-500" />
                    <span className="font-semibold text-white">Code Editor</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <select
                      value={language}
                      onChange={(e) => setLanguage(e.target.value)}
                      className="px-3 py-1 bg-dojo-black-800 border border-gray-700 rounded-lg text-white text-sm focus:outline-none focus:border-dojo-red-500"
                    >
                      <option value="javascript">JavaScript</option>
                      <option value="python">Python</option>
                      <option value="java">Java</option>
                      <option value="cpp">C++</option>
                    </select>
                    <Button
                      onClick={handleRunCode}
                      disabled={isExecuting}
                      size="sm"
                      className="gap-1"
                    >
                      {isExecuting ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Play className="h-4 w-4" />
                      )}
                      {isExecuting ? 'Running...' : 'Run'}
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="flex-1 p-0 min-h-0">
                <textarea
                  value={code}
                  onChange={(e) => handleCodeChange(e.target.value)}
                  className="w-full h-full p-4 bg-dojo-black-800 text-white font-mono text-sm focus:outline-none resize-none"
                  spellCheck={false}
                  placeholder="Write your code here..."
                  style={{ minHeight: '300px' }}
                  autoComplete="off"
                  autoCorrect="off"
                  autoCapitalize="off"
                />
              </CardContent>
            </Card>

            {/* Output Panel */}
            <Card className="h-48 flex flex-col">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Terminal className="h-5 w-5 text-dojo-red-500" />
                  <span className="font-semibold text-white">Output</span>
                </div>
              </CardHeader>
              <CardContent className="flex-1 overflow-y-auto">
                <pre className="text-sm text-gray-300 font-mono whitespace-pre-wrap">
                  {output || 'Run your code to see the output here...'}
                </pre>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-4 flex flex-col">
            {/* Video Controls */}
            <Card>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Video className="h-5 w-5 text-dojo-red-500" />
                  <span className="font-semibold text-white">Video Call</span>
                </div>
              </CardHeader>
              <CardContent>
                {!isVideoCallActive ? (
                  <Button
                    onClick={handleStartVideoCall}
                    className="w-full"
                    size="sm"
                  >
                    Start Video Call
                  </Button>
                ) : (
                  <>
                    <div className="mb-3">
                      <video
                        ref={localVideoRef}
                        autoPlay
                        muted
                        playsInline
                        className="w-full rounded-lg bg-dojo-black-900"
                        style={{ maxHeight: '150px' }}
                      />
                      <p className="text-xs text-gray-500 text-center mt-1">You</p>
                    </div>
                    <div className="flex gap-2 mb-2">
                      <Button
                        variant={isMicOn ? 'default' : 'outline'}
                        size="sm"
                        onClick={handleToggleMic}
                        className="flex-1"
                        disabled={!isVideoCallActive}
                      >
                        {isMicOn ? <Mic className="h-4 w-4" /> : <MicOff className="h-4 w-4" />}
                      </Button>
                      <Button
                        variant={isVideoOn ? 'default' : 'outline'}
                        size="sm"
                        onClick={handleToggleVideo}
                        className="flex-1"
                        disabled={!isVideoCallActive}
                      >
                        {isVideoOn ? <Video className="h-4 w-4" /> : <VideoOff className="h-4 w-4" />}
                      </Button>
                    </div>
                    <Button
                      onClick={handleStopVideoCall}
                      variant="outline"
                      className="w-full"
                      size="sm"
                    >
                      End Call
                    </Button>
                    {Array.from(remoteStreams.entries()).map(([userId, stream]) => {
                      const participant = participants.find(p => p.id === userId);
                      return (
                        <div key={userId} className="mt-3">
                          <video
                            autoPlay
                            playsInline
                            ref={(el) => {
                              if (el) el.srcObject = stream;
                            }}
                            className="w-full rounded-lg bg-dojo-black-900"
                            style={{ maxHeight: '150px' }}
                          />
                          <p className="text-xs text-gray-500 text-center mt-1">
                            {participant?.username || 'Participant'}
                          </p>
                        </div>
                      );
                    })}
                  </>
                )}
              </CardContent>
            </Card>

            {/* Participants */}
            <Card>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Users className="h-5 w-5 text-dojo-red-500" />
                  <span className="font-semibold text-white">Participants</span>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {participants.map((p) => (
                    <div key={p.id} className="flex items-center gap-2 text-sm">
                      <div className="h-2 w-2 rounded-full bg-green-500" />
                      <span className="text-gray-300">{p.username}</span>
                    </div>
                  ))}
                  {participants.length === 0 && (
                    <p className="text-gray-500 text-sm">No participants yet</p>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Chat */}
            <Card className="flex-1 flex flex-col min-h-[400px]">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <MessageCircle className="h-5 w-5 text-dojo-red-500" />
                  <span className="font-semibold text-white">Chat</span>
                </div>
              </CardHeader>
              <CardContent className="flex-1 flex flex-col p-0">
                <div className="flex-1 overflow-y-auto p-4 space-y-3 max-h-[300px]">
                  {messages.length === 0 ? (
                    <p className="text-gray-500 text-sm text-center">No messages yet</p>
                  ) : (
                    messages.map((msg) => {
                      const isOwnMessage = msg.user === user?.email || msg.user === user?.username;
                      const isSystemMessage = msg.user === 'System';
                      
                      if (isSystemMessage) {
                        return (
                          <div key={msg.id} className="text-center">
                            <span className="text-xs text-gray-500 italic">{msg.content}</span>
                          </div>
                        );
                      }
                      
                      return (
                        <div key={msg.id} className={`flex ${isOwnMessage ? 'justify-end' : 'justify-start'}`}>
                          <div className={`max-w-[75%] rounded-lg px-3 py-2 ${
                            isOwnMessage 
                              ? 'bg-dojo-red-500 text-white' 
                              : 'bg-dojo-black-800 text-gray-200'
                          }`}>
                            {!isOwnMessage && (
                              <div className="text-xs text-dojo-red-400 font-semibold mb-1">
                                {msg.user}
                              </div>
                            )}
                            <div className="text-sm break-words">{msg.content}</div>
                          </div>
                        </div>
                      );
                    })
                  )}
                  <div ref={messagesEndRef} />
                </div>
                
                <form onSubmit={handleSendMessage} className="p-4 border-t border-gray-700">
                  <div className="flex gap-2">
                    <Input
                      value={newMessage}
                      onChange={(e) => setNewMessage(e.target.value)}
                      placeholder="Type a message..."
                      disabled={!isConnected}
                    />
                    <Button type="submit" disabled={!isConnected || !newMessage.trim()}>
                      <Send className="h-4 w-4" />
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
