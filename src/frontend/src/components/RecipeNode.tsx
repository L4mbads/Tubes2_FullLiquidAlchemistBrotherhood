import { useContext, useEffect, useState } from 'react';
import Image from 'next/image';
import { Handle, NodeProps, Position, useNodeId, useStore } from 'reactflow';
import { CSSProperties } from 'react';
import { DarkModeContext } from './DarkModeProvider';

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
  const context = useContext(DarkModeContext);

  if (!context) {
    throw new Error('No Context!');
  }

  const { darkMode } = context;
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
    border: darkMode ? '2px solid #fff' : '2px solid #000',
    borderRadius: 5,
    background: darkMode ? '#734f9a' : '#fbac4e',
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
        <label style={{color: darkMode? 'white' : 'black'}}><strong>{data.name}</strong></label>  
        {elementInfo?.Type && <small style={{color: darkMode? 'white' : 'black'}}>Tier: {elementInfo.Type}</small>}
      </div>

      {hasOutgoing && (
        <Handle type="source" position={Position.Bottom} style={{ pointerEvents: 'none' }} />
      )}
    </div>
  );
}
