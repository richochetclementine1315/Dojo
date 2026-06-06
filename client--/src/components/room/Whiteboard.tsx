import { useRef, useEffect, useState, useCallback } from 'react';
import { Button } from '@/components/ui/Button';
import { Eraser, Trash2, Undo2, Pencil, Minus, Square, Circle as CircleIcon } from 'lucide-react';

// ─── Types ────────────────────────────────────────────────────────────────────

export type DrawTool = 'pen' | 'eraser' | 'line' | 'rect' | 'circle';

export interface WBPoint { x: number; y: number }

export interface WBStroke {
  id:        string;
  userId:    string;
  username:  string;
  color:     string;
  width:     number;
  tool:      DrawTool;
  points:    WBPoint[];
  action:    'start' | 'move' | 'end';
}

interface Props {
  /** Colour assigned to the local user (comes from WS hub) */
  localColor:   string;
  localUserId:  string;
  localUsername: string;
  /** Called when the user draws — parent should send over WS */
  onStroke: (stroke: WBStroke) => void;
  /** Remote stroke pushed in from WS */
  remoteStroke?: WBStroke | null;
  /** Remote clear event */
  remoteClear?: number;    // bump this value to trigger a clear
  /** Remote undo event */
  remoteUndo?: { userId: string; strokeId: string } | null;
}

// ─── Colour palette ───────────────────────────────────────────────────────────
const PALETTE = [
  '#FFFFFF', '#EF4444', '#F97316', '#EAB308',
  '#22C55E', '#3B82F6', '#A855F7', '#EC4899',
];

const BRUSH_SIZES = [2, 4, 8, 14];

// ─── Component ────────────────────────────────────────────────────────────────
export default function Whiteboard({
  localColor,
  localUserId,
  localUsername,
  onStroke,
  remoteStroke,
  remoteClear,
  remoteUndo,
}: Props) {
  const canvasRef   = useRef<HTMLCanvasElement>(null);
  const overlayRef  = useRef<HTMLCanvasElement>(null); // live preview while dragging
  const isDrawing   = useRef(false);
  const currentPath = useRef<WBPoint[]>([]);
  const strokeId    = useRef('');

  const [tool,   setTool]  = useState<DrawTool>('pen');
  const [color,  setColor] = useState(localColor || '#FFFFFF');
  const [width,  setWidth] = useState(4);
  // local stroke history for undo
  const strokeHistory = useRef<Array<{ id: string; snapshot: ImageData }>>([]);

  // ── Helpers ─────────────────────────────────────────────────────────────────

  const ctx = useCallback(
    () => canvasRef.current?.getContext('2d') ?? null,
    []
  );
  const overlayCtx = useCallback(
    () => overlayRef.current?.getContext('2d') ?? null,
    []
  );

  const relPos = (e: React.MouseEvent | React.TouchEvent): WBPoint => {
    const canvas = canvasRef.current!;
    const rect   = canvas.getBoundingClientRect();
    const src = 'touches' in e ? e.touches[0] : e;
    return {
      x: (src.clientX - rect.left) * (canvas.width  / rect.width),
      y: (src.clientY - rect.top)  * (canvas.height / rect.height),
    };
  };

  // Draw a completed stroke onto the main canvas
  const renderStroke = useCallback((s: WBStroke) => {
    const c = ctx();
    if (!c || s.points.length === 0) return;

    c.save();
    c.globalCompositeOperation = s.tool === 'eraser' ? 'destination-out' : 'source-over';
    c.strokeStyle = s.color;
    c.lineWidth   = s.width;
    c.lineCap     = 'round';
    c.lineJoin    = 'round';

    if (s.tool === 'pen' || s.tool === 'eraser') {
      c.beginPath();
      c.moveTo(s.points[0].x, s.points[0].y);
      for (const p of s.points.slice(1)) c.lineTo(p.x, p.y);
      c.stroke();
    } else if (s.tool === 'line' && s.points.length >= 2) {
      const p0 = s.points[0], pN = s.points[s.points.length - 1];
      c.beginPath();
      c.moveTo(p0.x, p0.y);
      c.lineTo(pN.x, pN.y);
      c.stroke();
    } else if (s.tool === 'rect' && s.points.length >= 2) {
      const p0 = s.points[0], pN = s.points[s.points.length - 1];
      c.strokeRect(p0.x, p0.y, pN.x - p0.x, pN.y - p0.y);
    } else if (s.tool === 'circle' && s.points.length >= 2) {
      const p0 = s.points[0], pN = s.points[s.points.length - 1];
      const rx = Math.abs(pN.x - p0.x) / 2;
      const ry = Math.abs(pN.y - p0.y) / 2;
      const cx = p0.x + (pN.x - p0.x) / 2;
      const cy = p0.y + (pN.y - p0.y) / 2;
      c.beginPath();
      c.ellipse(cx, cy, rx, ry, 0, 0, Math.PI * 2);
      c.stroke();
    }

    c.restore();
  }, [ctx]);

  // ── React to remote events ────────────────────────────────────────────────

  useEffect(() => {
    if (!remoteStroke) return;
    if (remoteStroke.action === 'end') {
      renderStroke(remoteStroke);
    } else if (remoteStroke.action === 'move') {
      // For pen/eraser show live preview on overlay for remote users too
      const ov = overlayCtx();
      if (!ov) return;
      ov.clearRect(0, 0, overlayRef.current!.width, overlayRef.current!.height);
      if (remoteStroke.points.length >= 2) {
        ov.save();
        ov.strokeStyle = remoteStroke.color;
        ov.lineWidth   = remoteStroke.width;
        ov.lineCap     = 'round';
        ov.beginPath();
        ov.moveTo(remoteStroke.points[0].x, remoteStroke.points[0].y);
        for (const p of remoteStroke.points.slice(1)) ov.lineTo(p.x, p.y);
        ov.stroke();
        ov.restore();
      }
    }
  }, [remoteStroke, renderStroke, overlayCtx]);

  useEffect(() => {
    if (remoteClear === undefined) return;
    const c = ctx();
    if (!c) return;
    c.clearRect(0, 0, canvasRef.current!.width, canvasRef.current!.height);
    strokeHistory.current = [];
  }, [remoteClear, ctx]);

  useEffect(() => {
    if (!remoteUndo) return;
    // Find and restore snapshot before that stroke
    const idx = strokeHistory.current.findIndex(s => s.id === remoteUndo.strokeId);
    if (idx >= 0) {
      const c = ctx();
      if (!c) return;
      c.putImageData(strokeHistory.current[idx].snapshot, 0, 0);
      strokeHistory.current = strokeHistory.current.slice(0, idx);
    }
  }, [remoteUndo, ctx]);

  // ── Mouse / touch handlers ────────────────────────────────────────────────

  const handlePointerDown = (e: React.MouseEvent) => {
    const canvas = canvasRef.current!;
    isDrawing.current  = true;
    strokeId.current   = `${localUserId}-${Date.now()}`;
    currentPath.current = [relPos(e)];

    // Save snapshot before this stroke (for undo)
    const c = ctx();
    if (c) {
      strokeHistory.current.push({
        id:       strokeId.current,
        snapshot: c.getImageData(0, 0, canvas.width, canvas.height),
      });
      // Keep history bounded
      if (strokeHistory.current.length > 50) strokeHistory.current.shift();
    }

    onStroke({
      id: strokeId.current, userId: localUserId, username: localUsername,
      color: tool === 'eraser' ? '#000000' : color,
      width, tool, points: currentPath.current, action: 'start',
    });
  };

  const handlePointerMove = (e: React.MouseEvent) => {
    if (!isDrawing.current) return;
    currentPath.current.push(relPos(e));

    // Live preview on overlay canvas
    const ov = overlayCtx();
    if (ov) {
      ov.clearRect(0, 0, overlayRef.current!.width, overlayRef.current!.height);
      ov.save();
      ov.globalCompositeOperation = tool === 'eraser' ? 'destination-out' : 'source-over';
      ov.strokeStyle = tool === 'eraser' ? 'rgba(0,0,0,1)' : color;
      ov.lineWidth   = width;
      ov.lineCap     = 'round';
      ov.lineJoin    = 'round';

      if (tool === 'pen' || tool === 'eraser') {
        ov.beginPath();
        ov.moveTo(currentPath.current[0].x, currentPath.current[0].y);
        for (const p of currentPath.current.slice(1)) ov.lineTo(p.x, p.y);
        ov.stroke();
      } else {
        const p0 = currentPath.current[0];
        const pN = currentPath.current[currentPath.current.length - 1];
        ov.beginPath();
        if (tool === 'line') {
          ov.moveTo(p0.x, p0.y); ov.lineTo(pN.x, pN.y);
        } else if (tool === 'rect') {
          ov.rect(p0.x, p0.y, pN.x - p0.x, pN.y - p0.y);
        } else if (tool === 'circle') {
          const rx = Math.abs(pN.x - p0.x) / 2, ry = Math.abs(pN.y - p0.y) / 2;
          ov.ellipse(p0.x + (pN.x - p0.x) / 2, p0.y + (pN.y - p0.y) / 2, rx, ry, 0, 0, Math.PI * 2);
        }
        ov.stroke();
      }
      ov.restore();
    }

    onStroke({
      id: strokeId.current, userId: localUserId, username: localUsername,
      color: tool === 'eraser' ? '#000000' : color,
      width, tool, points: [...currentPath.current], action: 'move',
    });
  };

  const handlePointerUp = (e: React.MouseEvent) => {
    if (!isDrawing.current) return;
    isDrawing.current = false;
    currentPath.current.push(relPos(e));

    // Clear overlay, commit to main canvas
    const ov = overlayCtx();
    if (ov) ov.clearRect(0, 0, overlayRef.current!.width, overlayRef.current!.height);

    const finalStroke: WBStroke = {
      id: strokeId.current, userId: localUserId, username: localUsername,
      color: tool === 'eraser' ? '#000000' : color,
      width, tool, points: [...currentPath.current], action: 'end',
    };
    renderStroke(finalStroke);
    onStroke(finalStroke);
    currentPath.current = [];
  };

  const handleClear = () => {
    const c = ctx();
    if (!c) return;
    c.clearRect(0, 0, canvasRef.current!.width, canvasRef.current!.height);
    strokeHistory.current = [];
    onStroke({ id: 'clear', userId: localUserId, username: localUsername, color: '', width: 0, tool: 'pen', points: [], action: 'end' });
  };

  const handleUndo = () => {
    if (strokeHistory.current.length === 0) return;
    const last = strokeHistory.current[strokeHistory.current.length - 1];
    const c = ctx();
    if (!c) return;
    // Restore snapshot before this stroke
    const prevIdx = strokeHistory.current.length - 2;
    if (prevIdx >= 0) {
      c.putImageData(strokeHistory.current[prevIdx].snapshot, 0, 0);
    } else {
      c.clearRect(0, 0, canvasRef.current!.width, canvasRef.current!.height);
    }
    strokeHistory.current.pop();
    onStroke({ id: last.id, userId: localUserId, username: localUsername, color: '', width: 0, tool: 'pen', points: [], action: 'end' });
  };

  // ── Render ────────────────────────────────────────────────────────────────
  return (
    <div className="flex flex-col h-full gap-2">
      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-2 px-1 py-1 bg-dojo-black-800 rounded-lg border border-gray-700">
        {/* Tools */}
        <div className="flex gap-1">
          {([
            { t: 'pen',    Icon: Pencil },
            { t: 'eraser', Icon: Eraser },
            { t: 'line',   Icon: Minus  },
            { t: 'rect',   Icon: Square },
            { t: 'circle', Icon: CircleIcon },
          ] as const).map(({ t, Icon }) => (
            <button
              key={t}
              onClick={() => setTool(t as DrawTool)}
              title={t.charAt(0).toUpperCase() + t.slice(1)}
              className={`p-1.5 rounded transition-colors ${
                tool === t
                  ? 'bg-dojo-red-500 text-white'
                  : 'text-gray-400 hover:text-white hover:bg-gray-700'
              }`}
            >
              <Icon className="h-4 w-4" />
            </button>
          ))}
        </div>

        <div className="w-px h-5 bg-gray-600" />

        {/* Colour palette */}
        <div className="flex gap-1">
          {PALETTE.map(c => (
            <button
              key={c}
              onClick={() => { setColor(c); setTool('pen'); }}
              className={`w-5 h-5 rounded-full border-2 transition-transform ${
                color === c && tool !== 'eraser'
                  ? 'border-white scale-110'
                  : 'border-transparent hover:scale-110'
              }`}
              style={{ background: c }}
              title={c}
            />
          ))}
          {/* custom colour picker */}
          <label title="Custom colour" className="cursor-pointer">
            <span
              className="block w-5 h-5 rounded-full border-2 border-dashed border-gray-500 hover:border-white overflow-hidden"
              style={{ background: PALETTE.includes(color) ? 'transparent' : color }}
            />
            <input
              type="color"
              className="sr-only"
              value={color}
              onChange={e => { setColor(e.target.value); setTool('pen'); }}
            />
          </label>
        </div>

        <div className="w-px h-5 bg-gray-600" />

        {/* Brush sizes */}
        <div className="flex items-center gap-1">
          {BRUSH_SIZES.map(s => (
            <button
              key={s}
              onClick={() => setWidth(s)}
              className={`flex items-center justify-center w-6 h-6 rounded transition-colors ${
                width === s ? 'bg-dojo-red-500' : 'hover:bg-gray-700'
              }`}
              title={`${s}px`}
            >
              <span
                className="rounded-full bg-white"
                style={{ width: s, height: s }}
              />
            </button>
          ))}
        </div>

        <div className="ml-auto flex gap-1">
          <Button variant="outline" size="sm" onClick={handleUndo}   title="Undo">
            <Undo2 className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="sm" onClick={handleClear}  title="Clear all"
            className="text-red-400 border-red-400/40 hover:bg-red-500/10">
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Canvas area */}
      <div className="relative flex-1 rounded-lg overflow-hidden border border-gray-700 bg-dojo-black-900 cursor-crosshair"
        style={{ minHeight: 360 }}>
        {/* Main committed canvas */}
        <canvas
          ref={canvasRef}
          width={1400}
          height={900}
          className="absolute inset-0 w-full h-full"
        />
        {/* Overlay for live preview */}
        <canvas
          ref={overlayRef}
          width={1400}
          height={900}
          className="absolute inset-0 w-full h-full"
          onMouseDown={handlePointerDown}
          onMouseMove={handlePointerMove}
          onMouseUp={handlePointerUp}
          onMouseLeave={handlePointerUp}
        />
      </div>
    </div>
  );
}
