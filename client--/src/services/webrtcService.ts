interface PeerConnection {
  connection: RTCPeerConnection;
  stream: MediaStream | null;
}

interface WebRTCServiceConfig {
  onRemoteStream: (userId: string, stream: MediaStream) => void;
  onRemoteStreamRemoved: (userId: string) => void;
  sendSignal: (signal: any) => void;
}

class WebRTCService {
  private peerConnections: Map<string, PeerConnection> = new Map();
  private localStream: MediaStream | null = null;
  private config: WebRTCServiceConfig | null = null;
  
  private iceServers = {
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' },
      { urls: 'stun:stun1.l.google.com:19302' },
    ],
  };

  initialize(config: WebRTCServiceConfig) {
    this.config = config;
  }

  async getLocalStream(audio: boolean = true, video: boolean = true): Promise<MediaStream> {
    if (this.localStream) {
      return this.localStream;
    }

    try {
      this.localStream = await navigator.mediaDevices.getUserMedia({
        video: video ? { width: 640, height: 480 } : false,
        audio: audio,
      });
      return this.localStream;
    } catch (error) {
      throw new Error('Failed to access camera/microphone');
    }
  }

  async createPeerConnection(userId: string, isInitiator: boolean = false): Promise<void> {
    if (this.peerConnections.has(userId)) {
      return;
    }

    const peerConnection = new RTCPeerConnection(this.iceServers);
    
    this.peerConnections.set(userId, {
      connection: peerConnection,
      stream: null,
    });

    // Add local stream tracks
    if (this.localStream) {
      this.localStream.getTracks().forEach(track => {
        if (this.localStream) {
          peerConnection.addTrack(track, this.localStream);
        }
      });
    }

    // Handle ICE candidates
    peerConnection.onicecandidate = (event) => {
      if (event.candidate && this.config) {
        this.config.sendSignal({
          type: 'ice-candidate',
          candidate: event.candidate,
          to: userId,
        });
      }
    };

    // Handle remote stream
    peerConnection.ontrack = (event) => {
      const peer = this.peerConnections.get(userId);
      if (peer && event.streams[0]) {
        peer.stream = event.streams[0];
        if (this.config) {
          this.config.onRemoteStream(userId, event.streams[0]);
        }
      }
    };

    // Handle connection state
    peerConnection.onconnectionstatechange = () => {
      if (peerConnection.connectionState === 'disconnected' || 
          peerConnection.connectionState === 'failed' ||
          peerConnection.connectionState === 'closed') {
        this.removePeer(userId);
      }
    };

    // If initiator, create and send offer
    if (isInitiator) {
      const offer = await peerConnection.createOffer();
      await peerConnection.setLocalDescription(offer);
      
      if (this.config) {
        this.config.sendSignal({
          type: 'offer',
          offer: offer,
          to: userId,
        });
      }
    }
  }

  async handleOffer(userId: string, offer: RTCSessionDescriptionInit): Promise<void> {
    await this.createPeerConnection(userId, false);
    
    const peer = this.peerConnections.get(userId);
    if (!peer) return;

    await peer.connection.setRemoteDescription(new RTCSessionDescription(offer));
    const answer = await peer.connection.createAnswer();
    await peer.connection.setLocalDescription(answer);

    if (this.config) {
      this.config.sendSignal({
        type: 'answer',
        answer: answer,
        to: userId,
      });
    }
  }

  async handleAnswer(userId: string, answer: RTCSessionDescriptionInit): Promise<void> {
    const peer = this.peerConnections.get(userId);
    if (!peer) return;

    await peer.connection.setRemoteDescription(new RTCSessionDescription(answer));
  }

  async handleIceCandidate(userId: string, candidate: RTCIceCandidateInit): Promise<void> {
    const peer = this.peerConnections.get(userId);
    if (!peer) return;

    try {
      await peer.connection.addIceCandidate(new RTCIceCandidate(candidate));
    } catch (error) {
      // Ignore ICE candidate errors
    }
  }

  toggleAudio(enabled: boolean): void {
    if (this.localStream) {
      this.localStream.getAudioTracks().forEach(track => {
        track.enabled = enabled;
      });
    }
  }

  toggleVideo(enabled: boolean): void {
    if (this.localStream) {
      this.localStream.getVideoTracks().forEach(track => {
        track.enabled = enabled;
      });
    }
  }

  removePeer(userId: string): void {
    const peer = this.peerConnections.get(userId);
    if (peer) {
      peer.connection.close();
      if (this.config) {
        this.config.onRemoteStreamRemoved(userId);
      }
      this.peerConnections.delete(userId);
    }
  }

  cleanup(): void {
    // Close all peer connections
    this.peerConnections.forEach((peer) => {
      peer.connection.close();
    });
    this.peerConnections.clear();

    // Stop local stream
    if (this.localStream) {
      this.localStream.getTracks().forEach(track => track.stop());
      this.localStream = null;
    }
  }

  getCurrentLocalStream(): MediaStream | null {
    return this.localStream;
  }
}

export const webrtcService = new WebRTCService();
