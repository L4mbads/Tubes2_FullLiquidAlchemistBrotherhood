import { Handle, NodeProps, Position, useNodeId, useStore } from 'reactflow';
import { CSSProperties } from 'react';

export type RecipeNodeData = {
  name: string;
};

export interface RecipeNodeType {
  name: string;
  recipes?: Recipe[] | null;
}

export interface Recipe {
  ingredient1: RecipeNodeType;
  ingredient2: RecipeNodeType;
}

export default function RecipeNode({ data }: NodeProps<RecipeNodeData>) {
  const nodeId = useNodeId();

  const connectedEdges = useStore((store) =>
    store.edges.filter((edge) => edge.source === nodeId || edge.target === nodeId)
  );

  const hasIncoming = connectedEdges.some((edge) => edge.target === nodeId);
  const hasOutgoing = connectedEdges.some((edge) => edge.source === nodeId);

  const style: CSSProperties = {
    padding: 10,
    border: '1px solid #ddd',
    borderRadius: 5,
    background: '#fff',
    pointerEvents: 'none',
    maxWidth: '100px',
  };

  return (
    <div style={style}>
      {hasIncoming && (
        <Handle type="target" position={Position.Top} style={{ pointerEvents: 'none' }} />
      )}
      <div className="nodrag">
        <strong>{data.name}</strong>
      </div>
      {hasOutgoing && (
        <Handle type="source" position={Position.Bottom} style={{ pointerEvents: 'none' }} />
      )}
    </div>
  );
}