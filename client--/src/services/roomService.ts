import api from '@/lib/api';
import type { Room, ApiResponse } from '@/types';

export const roomService = {
  async getRooms() {
    const response = await api.get<ApiResponse<Room[]>>('/rooms');
    return response.data.data;
  },

  async getRoom(id: string) {
    const response = await api.get<ApiResponse<Room>>(`/rooms/${id}`);
    return response.data.data;
  },

  async createRoom(data: {
    name: string;
    description?: string;
    max_participants: number;
  }) {
    const response = await api.post<ApiResponse<{ room: Room }>>('/rooms', data);
    return response.data.data.room;
  },

  async deleteRoom(id: string) {
    await api.delete(`/rooms/${id}`);
  },

  async joinRoom(id: string) {
    const response = await api.post<ApiResponse<any>>(`/rooms/${id}/join`);
    return response.data.data;
  },

  async joinRoomByCode(roomCode: string) {
    const response = await api.post<ApiResponse<{ room: Room }>>('/rooms/join', { room_code: roomCode });
    return response.data.data.room;
  },

  async leaveRoom(id: string) {
    await api.post(`/rooms/${id}/leave`);
  },

  getWebSocketUrl(roomId: string, token: string): string {
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const apiUrl = import.meta.env.VITE_API_URL?.replace('http://', '').replace('https://', '') || 'localhost:8080';
    const apiPath = apiUrl.includes('/api') ? apiUrl : `${apiUrl}/api`;
    return `${wsProtocol}//${apiPath}/rooms/${roomId}/ws?token=${token}`;
  },
};
