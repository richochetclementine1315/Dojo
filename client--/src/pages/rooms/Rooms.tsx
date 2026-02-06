import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { roomService } from '@/services/roomService';
import type { Room } from '@/types';
import Antigravity from '@/components/effects/Antigravity';
import { 
  Users, 
  Loader2,
  Plus,
  Trash2,
  UserPlus,
  Circle,
  LogIn
} from 'lucide-react';

export default function Rooms() {
  const navigate = useNavigate();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showJoinModal, setShowJoinModal] = useState(false);
  const [roomCode, setRoomCode] = useState('');
  const [newRoom, setNewRoom] = useState({
    name: '',
    description: '',
    max_participants: 4,
  });
  const [isCreating, setIsCreating] = useState(false);
  const [isJoining, setIsJoining] = useState(false);

  useEffect(() => {
    fetchRooms();
  }, []);

  const fetchRooms = async () => {
    try {
      setIsLoading(true);
      const data = await roomService.getRooms();
      setRooms(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      console.error('Failed to fetch rooms:', err);
      setError(err.response?.data?.message || 'Failed to load rooms');
      setRooms([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateRoom = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!newRoom.name.trim()) return;

    try {
      setIsCreating(true);
      const room = await roomService.createRoom(newRoom);
      setNewRoom({ name: '', description: '', max_participants: 4 });
      setShowCreateModal(false);
      navigate(`/rooms/${room.id}`);
    } catch (err: any) {
      setError('Failed to create room');
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteRoom = async (id: string) => {
    if (!confirm('Are you sure you want to delete this room?')) return;

    try {
      await roomService.deleteRoom(id);
      await fetchRooms();
    } catch (err: any) {
      setError('Failed to delete room');
    }
  };

  const handleJoinRoom = (roomId: string) => {
    navigate(`/rooms/${roomId}`);
  };

  const handleJoinWithCode = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!roomCode.trim()) return;

    try {
      setIsJoining(true);
      setError('');
      const room = await roomService.joinRoomByCode(roomCode.trim());
      setRoomCode('');
      setShowJoinModal(false);
      navigate(`/rooms/${room.id}`);
    } catch (err: any) {
      console.error('Join room error:', err);
      setError(err.response?.data?.message || err.message || 'Failed to join room. Please check the room code.');
      setIsJoining(false);
    }
  };

  return (
    <div className="min-h-screen bg-dojo-black-900 relative">
      {/* Antigravity Background */}
      <div className="fixed inset-0 z-0">
        <Antigravity
          count={300}
          magnetRadius={6}
          ringRadius={7}
          waveSpeed={0.4}
          waveAmplitude={1}
          particleSize={1.5}
          lerpSpeed={0.05}
          color="#d40808"
          autoAnimate
          particleVariance={1}
          rotationSpeed={0}
          depthFactor={1}
          pulseSpeed={3}
          particleShape="capsule"
          fieldStrength={10}
        />
      </div>
      <div className="relative z-10">
        <Navbar />
      
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold text-white mb-2">
              Collaboration <span className="dojo-text-gradient">Rooms</span>
            </h1>
            <p className="text-gray-400">Code together in real-time with other developers</p>
          </div>
          <div className="flex gap-3">
            <Button variant="outline" onClick={() => setShowJoinModal(true)}>
              <UserPlus className="mr-2 h-4 w-4" />
              Join Room
            </Button>
            <Button onClick={() => setShowCreateModal(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Create Room
            </Button>
          </div>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400">
            {error}
          </div>
        )}

        {/* Create Modal */}
        {showCreateModal && (
          <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
            <Card className="w-full max-w-md">
              <CardHeader>
                <h2 className="text-2xl font-bold text-white">Create New Room</h2>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleCreateRoom} className="space-y-4">
                  <div>
                    <label className="text-sm text-gray-400 mb-2 block">Room Name</label>
                    <Input
                      value={newRoom.name}
                      onChange={(e) => setNewRoom({ ...newRoom, name: e.target.value })}
                      placeholder="e.g., Daily Coding Session"
                      required
                    />
                  </div>

                  <div>
                    <label className="text-sm text-gray-400 mb-2 block">Description (Optional)</label>
                    <textarea
                      value={newRoom.description}
                      onChange={(e) => setNewRoom({ ...newRoom, description: e.target.value })}
                      placeholder="Brief description of this room..."
                      className="w-full px-4 py-2 bg-dojo-black-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-dojo-red-500"
                      rows={3}
                    />
                  </div>

                  <div>
                    <label className="text-sm text-gray-400 mb-2 block">Max Participants</label>
                    <Input
                      type="number"
                      min="2"
                      max="50"
                      value={newRoom.max_participants}
                      onChange={(e) => setNewRoom({ ...newRoom, max_participants: parseInt(e.target.value) })}
                      required
                    />
                  </div>

                  <div className="flex gap-3">
                    <Button
                      type="submit"
                      disabled={isCreating}
                      className="flex-1"
                    >
                      {isCreating ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          Creating...
                        </>
                      ) : (
                        'Create Room'
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => {
                        setShowCreateModal(false);
                        setNewRoom({ name: '', description: '', max_participants: 4 });
                      }}
                    >
                      Cancel
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Join Room Modal */}
        {showJoinModal && (
          <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
            <Card className="w-full max-w-md">
              <CardHeader>
                <h2 className="text-2xl font-bold text-white">Join Room</h2>
                <p className="text-gray-400 text-sm mt-1">Enter the room code to join</p>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleJoinWithCode} className="space-y-4">
                  <div>
                    <label className="text-sm text-gray-400 mb-2 block">Room Code</label>
                    <Input
                      value={roomCode}
                      onChange={(e) => setRoomCode(e.target.value.toUpperCase())}
                      placeholder="e.g., ABC123"
                      className="font-mono text-lg tracking-wider"
                      required
                    />
                    <p className="text-xs text-gray-500 mt-2">
                      Ask the room creator for the room code
                    </p>
                  </div>

                  <div className="flex gap-3">
                    <Button
                      type="submit"
                      disabled={isJoining}
                      className="flex-1"
                    >
                      {isJoining ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          Joining...
                        </>
                      ) : (
                        <>
                          <UserPlus className="mr-2 h-4 w-4" />
                          Join Room
                        </>
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => {
                        setShowJoinModal(false);
                        setRoomCode('');
                      }}
                    >
                      Cancel
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Rooms Grid */}
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
          </div>
        ) : rooms.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Users className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400 text-lg">No active rooms</p>
              <p className="text-gray-500 text-sm mt-2">Create a room to start collaborating</p>
              <Button onClick={() => setShowCreateModal(true)} className="mt-4">
                <Plus className="mr-2 h-4 w-4" />
                Create Room
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {rooms.map((room) => {
              const isFull = room.current_participants >= room.max_participants;
              const participantsPercentage = (room.current_participants / room.max_participants) * 100;
              
              return (
                <Card key={room.id} className="group hover:border-dojo-red-500/50 transition-all">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <Circle 
                            className={`h-3 w-3 ${room.is_active ? 'fill-green-500 text-green-500' : 'fill-gray-500 text-gray-500'}`}
                          />
                          <span className="text-xs text-gray-400">
                            {room.is_active ? 'Active' : 'Inactive'}
                          </span>
                        </div>
                        <h3 className="text-xl font-semibold text-white group-hover:text-dojo-red-400 transition-colors">
                          {room.name}
                        </h3>
                      </div>
                      <button
                        onClick={() => handleDeleteRoom(room.id)}
                        className="text-gray-400 hover:text-red-400 transition-colors"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </CardHeader>
                  
                  <CardContent>
                    {room.description && (
                      <p className="text-gray-400 text-sm mb-4 line-clamp-2">
                        {room.description}
                      </p>
                    )}

                    {/* Participants Info */}
                    <div className="mb-4">
                      <div className="flex items-center justify-between text-sm mb-2">
                        <div className="flex items-center gap-2 text-gray-400">
                          <Users className="h-4 w-4" />
                          <span>Participants</span>
                        </div>
                        <span className={`font-medium ${isFull ? 'text-red-400' : 'text-gray-300'}`}>
                          {room.current_participants} / {room.max_participants}
                        </span>
                      </div>
                      <div className="w-full bg-dojo-black-800 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full transition-all ${
                            isFull 
                              ? 'bg-red-500' 
                              : 'bg-gradient-to-r from-dojo-red-500 to-orange-500'
                          }`}
                          style={{ width: `${participantsPercentage}%` }}
                        />
                      </div>
                    </div>

                    <Button
                      onClick={() => handleJoinRoom(room.id)}
                      disabled={isFull}
                      className="w-full"
                    >
                      {isFull ? (
                        <>
                          <Circle className="mr-2 h-4 w-4" />
                          Room Full
                        </>
                      ) : (
                        <>
                          <LogIn className="mr-2 h-4 w-4" />
                          Join Room
                        </>
                      )}
                    </Button>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </div>
      </div>
    </div>
  );
}
