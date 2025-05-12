import { useEffect, useState } from 'react';
import Image from 'next/image';
import { Handle, NodeProps, Position, useNodeId, useStore } from 'reactflow';
import { CSSProperties } from 'react';

export interface RecipeNodeType {
  name: string;
  recipes?: Recipe[] | null;
}

export interface Recipe {
  ingredient1: RecipeNodeType;
  ingredient2: RecipeNodeType;
}

export type RecipeNodeData = {
  name: string;
};

type ElementInfo = {
  Name: string;
  Type: string;
  ImageUrl: string;
};

export default function RecipeNode({ data }: NodeProps<RecipeNodeData>) {
  const nodeId = useNodeId();
  const [elementInfo, setElementInfo] = useState<ElementInfo | null>(null);

  useEffect(() => {
    fetch(`http://localhost:8000/api/go/element/${data.name}`)
      .then((res) => res.json())
      .then(setElementInfo)
      .catch(console.error);
  }, [data.name]);

  const connectedEdges = useStore((store) =>
    store.edges.filter((edge) => edge.source === nodeId || edge.target === nodeId)
  );

  const hasIncoming = connectedEdges.some((edge) => edge.target === nodeId);
  const hasOutgoing = connectedEdges.some((edge) => edge.source === nodeId);

  const style: CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    padding: 10,
    border: '1px solid #ddd',
    borderRadius: 5,
    background: '#8bd450',
    pointerEvents: 'none',
    maxWidth: '150px',
    minHeight: '40px',
    gap: '10px',
    flexDirection: 'row',
    flexWrap: 'nowrap',
    overflow: 'hidden',
  };

  return (
    <div style={style}>
      {hasIncoming && (
        <Handle type="target" position={Position.Top} style={{ pointerEvents: 'none' }} />
      )}

      {elementInfo?.ImageUrl && (
        <div style={{ flexShrink: 0 }}>
          <Image
            src={elementInfo.ImageUrl}
            alt={data.name}
            width={40}
            height={40}
            style={{ objectFit: 'contain' }}
          />
        </div>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', overflowWrap:'break-word', wordBreak: 'break-word', whiteSpace: 'normal', flex: 1 }} className="nodrag">
        <strong>{data.name}</strong>
        {elementInfo?.Type && <small style={{ color: '#666' }}>Tier: {elementInfo.Type}</small>}
      </div>

      {hasOutgoing && (
        <Handle type="source" position={Position.Bottom} style={{ pointerEvents: 'none' }} />
      )}
    </div>
  );
}
