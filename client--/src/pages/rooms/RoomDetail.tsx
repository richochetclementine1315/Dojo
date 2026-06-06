import { useState, useEffect, useRef, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { roomService } from '@/services/roomService';
import { codeExecutionService } from '@/services/codeExecutionService';
import { webrtcService } from '@/services/webrtcService';
import { useAuthStore } from '@/store/authStore';
import Whiteboard, { type WBStroke } from '@/components/room/Whiteboard';
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
  Terminal,
  PenLine,
} from 'lucide-react';

// ─── Local types ──────────────────────────────────────────────────────────────

interface Message {
  id:        string;
  user:      string;
  content:   string;
  timestamp: string;
}

interface Participant {
  id:       string;
  username: string;
  color?:   string;
}

type MainTab = 'code' | 'whiteboard';

// ─── User colours (same pool as the backend) ─────────────────────────────────
const USER_COLORS = [
  '#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8',
  '#F7DC6F', '#BB8FCE', '#85C1E2', '#F8B195', '#C06C84',
];

// ─── Component ────────────────────────────────────────────────────────────────

export default function RoomDetail() {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate   = useNavigate();
  const { user, ensureFreshToken } = useAuthStore();

  const [room,      setRoom]      = useState<Room | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error,     setError]     = useState('');

  // WebSocket
  const wsRef           = useRef<WebSocket | null>(null);
  const connectionIdRef = useRef<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  // Chat
  const [messages,    setMessages]    = useState<Message[]>([]);
  const [newMessage,  setNewMessage]  = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const processedMessageIds = useRef<Set<string>>(new Set());

  // Code Editor
  const [code,        setCode]        = useState('// Start coding here...\n\n');
  const [language,    setLanguage]    = useState('javascript');
  const [output,      setOutput]      = useState('');
  const [isExecuting, setIsExecuting] = useState(false);

  // Participants
  const [participants, setParticipants] = useState<Participant[]>([]);

  // Video
  const [isMicOn,          setIsMicOn]          = useState(false);
  const [isVideoOn,        setIsVideoOn]         = useState(false);
  const [localStream,      setLocalStream]       = useState<MediaStream | null>(null);
  const [remoteStreams,     setRemoteStreams]     = useState<Map<string, MediaStream>>(new Map());
  const localVideoRef     = useRef<HTMLVideoElement>(null);
  const [isVideoCallActive, setIsVideoCallActive] = useState(false);

  // Tabs
  const [mainTab, setMainTab] = useState<MainTab>('code');

  // Whiteboard WS events
  const [remoteStroke, setRemoteStroke] = useState<WBStroke | null>(null);
  const [remoteClear,  setRemoteClear]  = useState<number>(0);
  const [remoteUndo,   setRemoteUndo]   = useState<{ userId: string; strokeId: string } | null>(null);

  // Local user colour (assigned from hub, fallback to deterministic)
  const localColor = useRef<string>(
    USER_COLORS[(user?.username?.charCodeAt(0) ?? 0) % USER_COLORS.length]
  );

  // ── Effects ─────────────────────────────────────────────────────────────────

  useEffect(() => {
    if (!roomId) return;
    fetchRoom();
  }, [roomId]);

  useEffect(() => {
    if (!roomId) return;
    const connectionId = `conn-${Date.now()}-${Math.random()}`;
    connectionIdRef.current = connectionId;
    if (wsRef.current) { wsRef.current.close(); wsRef.current = null; }
    connectWebSocket();
    return () => {
      if (wsRef.current && connectionIdRef.current === connectionId) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [roomId]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  useEffect(() => {
    webrtcService.initialize({
      onRemoteStream: (userId: string, stream: MediaStream) => {
        setRemoteStreams(prev => new Map(prev).set(userId, stream));
      },
      onRemoteStreamRemoved: (userId: string) => {
        setRemoteStreams(prev => { const m = new Map(prev); m.delete(userId); return m; });
      },
      sendSignal: (signal: any) => {
        if (wsRef.current && isConnected) {
          wsRef.current.send(JSON.stringify({ type: 'webrtc-signal', data: signal }));
        }
      },
    });
    return () => { webrtcService.cleanup(); };
  }, [isConnected]);

  useEffect(() => {
    if (localVideoRef.current && localStream) {
      localVideoRef.current.srcObject = localStream;
    }
  }, [localStream]);

  // ── Data fetching ────────────────────────────────────────────────────────

  const fetchRoom = async () => {
    if (!roomId) return;
    try {
      setIsLoading(true);
      const data = await roomService.getRoom(roomId);
      setRoom((data as any).room || data);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load room');
    } finally {
      setIsLoading(false);
    }
  };

  // ── WebSocket ─────────────────────────────────────────────────────────────

  const connectWebSocket = async () => {
    if (!roomId) return;
    try {
      const freshToken = await ensureFreshToken();
      if (!freshToken) { setError('Authentication failed. Please log in again.'); navigate('/login'); return; }

      const socket = new WebSocket(roomService.getWebSocketUrl(roomId, freshToken));

      socket.onopen = () => {
        setIsConnected(true);
        socket.send(JSON.stringify({ type: 'join', user: user?.username || 'Anonymous' }));
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);

          switch (data.type) {
            // ── chat ──────────────────────────────────────────────────────
            case 'chat': {
              const chatData = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              const content  = chatData?.message || '';
              const msgId = data.Timestamp
                ? `${data.UserID}-${data.Timestamp}`
                : `${data.UserID}-${content}-${Date.now()}`;
              if (processedMessageIds.current.has(msgId)) break;
              processedMessageIds.current.add(msgId);
              if (processedMessageIds.current.size > 100) {
                const arr = Array.from(processedMessageIds.current);
                processedMessageIds.current = new Set(arr.slice(-100));
              }
              setMessages(prev => [...prev, {
                id: `msg-${Date.now()}-${Math.random()}`,
                user: data.Username || 'Unknown',
                content,
                timestamp: new Date().toISOString(),
              }]);
              break;
            }

            // ── whiteboard ────────────────────────────────────────────────
            case 'whiteboard_draw': {
              const payload = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              if (payload && payload.userId !== user?.id) {
                setRemoteStroke({ ...payload });
              }
              break;
            }
            case 'whiteboard_clear': {
              setRemoteClear(prev => prev + 1);
              break;
            }
            case 'whiteboard_undo': {
              const payload = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              if (payload) setRemoteUndo(payload);
              break;
            }

            // ── WebRTC ────────────────────────────────────────────────────
            case 'rtc_offer': {
              const d = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              handleWebRTCSignal(data.UserID, { type: 'offer', offer: d });
              break;
            }
            case 'rtc_answer': {
              const d = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              handleWebRTCSignal(data.UserID, { type: 'answer', answer: d });
              break;
            }
            case 'rtc_candidate': {
              const d = typeof data.Data === 'string' ? JSON.parse(data.Data) : data.Data;
              handleWebRTCSignal(data.UserID, { type: 'ice-candidate', candidate: d });
              break;
            }

            // ── participants ──────────────────────────────────────────────
            case 'user_list': {
              if (Array.isArray(data.Data)) {
                const list: Participant[] = data.Data.map((p: any) => ({
                  id:       p.UserID || p.user_id || p.id,
                  username: p.Username || p.username || p.email,
                  color:    p.Color || p.color,
                }));
                const unique = list.filter(
                  (v, i, a) => a.findIndex(x => x.id === v.id) === i
                );
                setParticipants(unique);
                // Update local colour from hub assignment
                const me = unique.find(p => p.id === user?.id);
                if (me?.color) localColor.current = me.color;
              }
              break;
            }
            case 'user_joined': {
              const name = data.Username || data.Data?.username || 'Someone';
              setMessages(prev => [...prev, {
                id: `join-${Date.now()}`, user: 'System',
                content: `${name} joined the room`, timestamp: new Date().toISOString(),
              }]);
              if (isVideoCallActive && data.UserID && data.UserID !== user?.id) {
                setTimeout(() => webrtcService.createPeerConnection(data.UserID, true), 1000);
              }
              break;
            }
            case 'user_left': {
              const name = data.Username || data.Data?.username || 'Someone';
              setMessages(prev => [...prev, {
                id: `leave-${Date.now()}`, user: 'System',
                content: `${name} left the room`, timestamp: new Date().toISOString(),
              }]);
              break;
            }
            default: break;
          }
        } catch { /* silent */ }
      };

      socket.onerror  = () => setError('WebSocket connection error');
      socket.onclose  = () => setIsConnected(false);
      wsRef.current   = socket;
    } catch {
      setError('Failed to establish connection. Please try again.');
    }
  };

  // ── Handlers ────────────────────────────────────────────────────────────────

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !wsRef.current || !isConnected) return;
    wsRef.current.send(JSON.stringify({ type: 'chat', data: { message: newMessage } }));
    setNewMessage('');
  };

  const handleRunCode = async () => {
    if (!code.trim()) { setOutput('Error: No code to execute'); return; }
    setIsExecuting(true);
    setOutput('Running…');
    try {
      const result = await codeExecutionService.executeCode(language, code);
      let out = result.stdout || result.output || result.stderr || 'No output.';
      if (result.stderr && !result.stdout) out = `Error:\n${result.stderr}`;
      if (result.status && result.status !== 'Accepted') {
        out = `[${result.status}]\n${out}`;
      }
      if (result.time) out += `\n\n⏱ ${result.time}s`;
      setOutput(out);
    } catch (err: any) {
      setOutput(`Execution Error:\n${err.message}`);
    } finally {
      setIsExecuting(false);
    }
  };

  // Whiteboard stroke → WS
  const handleWhiteboardStroke = useCallback((stroke: WBStroke) => {
    if (!wsRef.current || !isConnected) return;

    if (stroke.id === 'clear') {
      wsRef.current.send(JSON.stringify({ type: 'whiteboard_clear', data: {} }));
    } else if (stroke.action === 'end' && stroke.points.length === 0) {
      // undo (empty points + non-clear id)
      wsRef.current.send(JSON.stringify({
        type: 'whiteboard_undo',
        data: { userId: stroke.userId, strokeId: stroke.id },
      }));
    } else {
      wsRef.current.send(JSON.stringify({ type: 'whiteboard_draw', data: stroke }));
    }
  }, [isConnected]);

  // WebRTC
  const handleWebRTCSignal = async (userId: string, signal: any) => {
    try {
      switch (signal.type) {
        case 'offer':
          if (!isVideoCallActive && !localStream) {
            const stream = await webrtcService.getLocalStream(true, true);
            setLocalStream(stream); setIsVideoCallActive(true); setIsMicOn(true); setIsVideoOn(true);
            await webrtcService.handleOffer(userId, signal.offer);
          } else {
            await webrtcService.handleOffer(userId, signal.offer);
          }
          break;
        case 'answer':    await webrtcService.handleAnswer(userId, signal.answer);      break;
        case 'ice-candidate': await webrtcService.handleIceCandidate(userId, signal.candidate); break;
      }
    } catch { /* ignore */ }
  };

  const handleStartVideoCall = async () => {
    try {
      const stream = await webrtcService.getLocalStream(true, true);
      setLocalStream(stream); setIsVideoCallActive(true); setIsMicOn(true); setIsVideoOn(true);
      setTimeout(() => {
        participants.forEach(p => {
          if (p.id && p.id !== user?.id) webrtcService.createPeerConnection(p.id, true);
        });
      }, 500);
    } catch {
      setError('Failed to start video call. Please allow camera and microphone access.');
    }
  };

  const handleStopVideoCall = () => {
    webrtcService.cleanup();
    setLocalStream(null); setRemoteStreams(new Map());
    setIsVideoCallActive(false); setIsMicOn(false); setIsVideoOn(false);
  };

  const handleLeaveRoom = async () => {
    if (wsRef.current) {
      wsRef.current.send(JSON.stringify({ type: 'leave', user: user?.username || 'Anonymous' }));
      wsRef.current.close(); wsRef.current = null;
    }
    if (roomId) await roomService.leaveRoom(roomId);
    navigate('/rooms');
  };

  // ── Loading / not found states ───────────────────────────────────────────

  if (isLoading) return (
    <div className="min-h-screen bg-dojo-black-900 flex items-center justify-center">
      <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
    </div>
  );

  if (!room) return (
    <div className="min-h-screen bg-dojo-black-900">
      <Navbar />
      <div className="container mx-auto px-4 py-8">
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-gray-400">Room not found</p>
            <Button onClick={() => navigate('/rooms')} className="mt-4">Back to Rooms</Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );

  // ── Main render ───────────────────────────────────────────────────────────

  return (
    <div className="min-h-screen bg-dojo-black-900">
      <Navbar />

      <div className="container mx-auto px-4 py-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-white mb-1">{room.name}</h1>
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
                  {room.room_code && room.room_code !== 'N/A' ? room.room_code : 'Loading…'}
                </span>
                <button
                  onClick={async () => {
                    if (room.room_code) await navigator.clipboard.writeText(room.room_code);
                  }}
                  className="text-gray-400 hover:text-white"
                  title="Copy room code"
                >📋</button>
              </div>
            </div>
          </div>
          <Button variant="outline" onClick={handleLeaveRoom}>
            <LogOut className="mr-2 h-4 w-4" /> Leave Room
          </Button>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
            {error}
          </div>
        )}

        {/* ── Main grid ─────────────────────────────────────────────────── */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-4" style={{ minHeight: 'calc(100vh - 200px)' }}>

          {/* Left: code + whiteboard tabs */}
          <div className="lg:col-span-3 flex flex-col gap-3">
            {/* Tab bar */}
            <div className="flex gap-2">
              <Button
                size="sm"
                variant={mainTab === 'code' ? 'default' : 'outline'}
                className="rounded-full gap-2"
                onClick={() => setMainTab('code')}
              >
                <Code className="h-4 w-4" /> Code Editor
              </Button>
              <Button
                size="sm"
                variant={mainTab === 'whiteboard' ? 'default' : 'outline'}
                className="rounded-full gap-2"
                onClick={() => setMainTab('whiteboard')}
              >
                <PenLine className="h-4 w-4" /> Whiteboard
              </Button>
            </div>

            {/* ── Code tab ───────────────────────────────────────────────── */}
            {mainTab === 'code' && (
              <>
                <Card className="flex flex-col" style={{ height: 'calc(100vh - 380px)', minHeight: 300 }}>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Code className="h-5 w-5 text-dojo-red-500" />
                        <span className="font-semibold text-white">Code Editor</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <select
                          value={language}
                          onChange={e => setLanguage(e.target.value)}
                          className="px-3 py-1 bg-dojo-black-800 border border-gray-700 rounded-lg text-white text-sm focus:outline-none focus:border-dojo-red-500"
                        >
                          <option value="javascript">JavaScript</option>
                          <option value="python">Python</option>
                          <option value="java">Java</option>
                          <option value="cpp">C++</option>
                          <option value="go">Go</option>
                          <option value="rust">Rust</option>
                          <option value="typescript">TypeScript</option>
                          <option value="csharp">C#</option>
                        </select>
                        <Button onClick={handleRunCode} disabled={isExecuting} size="sm" className="gap-1">
                          {isExecuting ? <Loader2 className="h-4 w-4 animate-spin" /> : <Play className="h-4 w-4" />}
                          {isExecuting ? 'Running…' : 'Run'}
                        </Button>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent className="flex-1 p-0 min-h-0">
                    <textarea
                      value={code}
                      onChange={e => setCode(e.target.value)}
                      className="w-full h-full p-4 bg-dojo-black-800 text-white font-mono text-sm focus:outline-none resize-none"
                      spellCheck={false}
                      placeholder="Write your code here…"
                      autoComplete="off"
                      autoCorrect="off"
                      autoCapitalize="off"
                    />
                  </CardContent>
                </Card>

                {/* Output */}
                <Card style={{ height: 180 }} className="flex flex-col">
                  <CardHeader>
                    <div className="flex items-center gap-2">
                      <Terminal className="h-5 w-5 text-dojo-red-500" />
                      <span className="font-semibold text-white">Output</span>
                      {output && (
                        <span className="ml-auto text-xs text-gray-500">
                          Powered by Judge0 CE
                        </span>
                      )}
                    </div>
                  </CardHeader>
                  <CardContent className="flex-1 overflow-y-auto">
                    <pre className="text-sm text-gray-300 font-mono whitespace-pre-wrap">
                      {output || 'Run your code to see output here…'}
                    </pre>
                  </CardContent>
                </Card>
              </>
            )}

            {/* ── Whiteboard tab ─────────────────────────────────────────── */}
            {mainTab === 'whiteboard' && (
              <div className="flex-1" style={{ height: 'calc(100vh - 280px)', minHeight: 480 }}>
                <Whiteboard
                  localColor={localColor.current}
                  localUserId={user?.id ?? 'anon'}
                  localUsername={user?.username ?? 'Anonymous'}
                  onStroke={handleWhiteboardStroke}
                  remoteStroke={remoteStroke}
                  remoteClear={remoteClear}
                  remoteUndo={remoteUndo}
                />
              </div>
            )}
          </div>

          {/* Right sidebar */}
          <div className="space-y-4 flex flex-col">
            {/* Video */}
            <Card>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Video className="h-5 w-5 text-dojo-red-500" />
                  <span className="font-semibold text-white">Video Call</span>
                </div>
              </CardHeader>
              <CardContent>
                {!isVideoCallActive ? (
                  <Button onClick={handleStartVideoCall} className="w-full" size="sm">
                    Start Video Call
                  </Button>
                ) : (
                  <>
                    <div className="mb-3">
                      <video ref={localVideoRef} autoPlay muted playsInline
                        className="w-full rounded-lg bg-dojo-black-900" style={{ maxHeight: 130 }} />
                      <p className="text-xs text-gray-500 text-center mt-1">You</p>
                    </div>
                    <div className="flex gap-2 mb-2">
                      <Button variant={isMicOn ? 'default' : 'outline'} size="sm" onClick={() => { const n = !isMicOn; setIsMicOn(n); webrtcService.toggleAudio(n); }} className="flex-1">
                        {isMicOn ? <Mic className="h-4 w-4" /> : <MicOff className="h-4 w-4" />}
                      </Button>
                      <Button variant={isVideoOn ? 'default' : 'outline'} size="sm" onClick={() => { const n = !isVideoOn; setIsVideoOn(n); webrtcService.toggleVideo(n); }} className="flex-1">
                        {isVideoOn ? <Video className="h-4 w-4" /> : <VideoOff className="h-4 w-4" />}
                      </Button>
                    </div>
                    <Button onClick={handleStopVideoCall} variant="outline" className="w-full" size="sm">
                      End Call
                    </Button>
                    {Array.from(remoteStreams.entries()).map(([userId, stream]) => {
                      const p = participants.find(x => x.id === userId);
                      return (
                        <div key={userId} className="mt-3">
                          <video autoPlay playsInline ref={el => { if (el) el.srcObject = stream; }}
                            className="w-full rounded-lg bg-dojo-black-900" style={{ maxHeight: 130 }} />
                          <p className="text-xs text-gray-500 text-center mt-1">
                            {p?.username || 'Participant'}
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
                  {participants.map(p => (
                    <div key={p.id} className="flex items-center gap-2 text-sm">
                      <div
                        className="h-2.5 w-2.5 rounded-full flex-shrink-0"
                        style={{ background: p.color || '#22c55e' }}
                      />
                      <span className="text-gray-300">{p.username}</span>
                      {p.color && (
                        <span
                          className="ml-auto text-xs px-1.5 rounded border"
                          style={{ borderColor: p.color + '60', color: p.color }}
                        >ink</span>
                      )}
                    </div>
                  ))}
                  {participants.length === 0 && (
                    <p className="text-gray-500 text-sm">No participants yet</p>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Chat */}
            <Card className="flex-1 flex flex-col" style={{ minHeight: 360 }}>
              <CardHeader>
                <div className="flex items-center gap-2">
                  <MessageCircle className="h-5 w-5 text-dojo-red-500" />
                  <span className="font-semibold text-white">Chat</span>
                </div>
              </CardHeader>
              <CardContent className="flex-1 flex flex-col p-0">
                <div className="flex-1 overflow-y-auto p-4 space-y-3" style={{ maxHeight: 300 }}>
                  {messages.length === 0
                    ? <p className="text-gray-500 text-sm text-center">No messages yet</p>
                    : messages.map(msg => {
                        const isOwn   = msg.user === user?.email || msg.user === user?.username;
                        const isSys   = msg.user === 'System';
                        if (isSys) return (
                          <div key={msg.id} className="text-center">
                            <span className="text-xs text-gray-500 italic">{msg.content}</span>
                          </div>
                        );
                        return (
                          <div key={msg.id} className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}>
                            <div className={`max-w-[75%] rounded-lg px-3 py-2 ${
                              isOwn ? 'bg-dojo-red-500 text-white' : 'bg-dojo-black-800 text-gray-200'
                            }`}>
                              {!isOwn && <div className="text-xs text-dojo-red-400 font-semibold mb-1">{msg.user}</div>}
                              <div className="text-sm break-words">{msg.content}</div>
                            </div>
                          </div>
                        );
                      })
                  }
                  <div ref={messagesEndRef} />
                </div>
                <form onSubmit={handleSendMessage} className="p-4 border-t border-gray-700">
                  <div className="flex gap-2">
                    <Input
                      value={newMessage}
                      onChange={e => setNewMessage(e.target.value)}
                      placeholder="Type a message…"
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
