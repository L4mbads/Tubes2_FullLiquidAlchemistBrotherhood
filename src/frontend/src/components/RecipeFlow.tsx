import React, { useMemo, useEffect } from 'react';
import ReactFlow, { Background, BackgroundVariant, Node, Edge, Position, MiniMap, useReactFlow, ReactFlowProvider } from 'reactflow';
import 'reactflow/dist/style.css';
import '../components/style.css';
import RecipeNode, { RecipeNodeType } from './RecipeNode';

type RecipeFlowProps = {
  tree: RecipeNodeType | null;
};

function calculateSubtreeWidth(recipe: RecipeNodeType): number {
  if (!recipe.recipes || recipe.recipes.length === 0) return 1;
  let totalWidth = 0;
  recipe.recipes.forEach(r => {
    totalWidth += calculateSubtreeWidth(r.ingredient1) + calculateSubtreeWidth(r.ingredient2);
  });
  return Math.max(1, totalWidth);
}

let nodesIdCounter = 0;
let edgesIdCounter = 0;

function layoutTree(
  recipe: RecipeNodeType,
  depth: number,
  xOffset: number,
  nodes: Node[],
  edges: Edge[],
  parentId: string | null = null,
  bounds = { minX: Infinity, minY: Infinity, maxX: -Infinity, maxY: -Infinity }
): { width: number; centerX: number; bounds: typeof bounds } {
  const widthSpacing = 100;
  const heightSpacing = 150;

  const recipeId = `node_${nodesIdCounter++}`;
  let subtreeCenterX = xOffset;
  let totalSubtreeWidth = 1;

  if (recipe.recipes && recipe.recipes.length > 0) {
    let currentX = xOffset;
    const childrenWidths: { centerX: number; width: number }[] = [];

    for (const subRecipe of recipe.recipes) {
      const w1 = calculateSubtreeWidth(subRecipe.ingredient1);
      const w2 = calculateSubtreeWidth(subRecipe.ingredient2);
      const totalWidth = w1 + w2;
      const midX = currentX + (totalWidth * widthSpacing) / 2;
      const stepId = `step_${nodesIdCounter++}`;
      const stepY = -((depth + 0.5) * heightSpacing);

      nodes.push({
        id: stepId,
        type: 'default',
        data: { label: '' },
        position: { x: midX, y: stepY },
        style: {
          background: 'black',
          width: '1px',
          height: '1px',
          pointerEvents: 'none',
          padding: 0,
          margin: 0,
        },
        sourcePosition: Position.Bottom,
        targetPosition: Position.Top,
        className: 'node-with-hidden-handles',
        draggable: false,
        selectable: false,
      });

      edges.push({
        id: `edge_${edgesIdCounter++}`,
        source: stepId,
        target: recipeId,
        type: 'smoothstep',
      });

      layoutTree(subRecipe.ingredient1, depth + 1, currentX, nodes, edges, stepId, bounds);
      layoutTree(subRecipe.ingredient2, depth + 1, currentX + w1 * widthSpacing, nodes, edges, stepId, bounds);

      currentX += totalWidth * widthSpacing;
      childrenWidths.push({ width: totalWidth, centerX: midX });

      bounds.minX = Math.min(bounds.minX, midX);
      bounds.maxX = Math.max(bounds.maxX, midX);
      bounds.minY = Math.min(bounds.minY, stepY);
      bounds.maxY = Math.max(bounds.maxY, stepY);
    }

    const avgX = childrenWidths.reduce((sum, w) => sum + w.centerX, 0) / childrenWidths.length;
    subtreeCenterX = avgX;
    totalSubtreeWidth = childrenWidths.reduce((sum, c) => sum + c.width, 0);
  }

  const nodeY = -depth * heightSpacing;
  nodes.push({
    id: recipeId,
    type: 'recipeNode',
    data: { name: recipe.name },
    position: { x: subtreeCenterX, y: nodeY },
  });

  bounds.minX = Math.min(bounds.minX, subtreeCenterX);
  bounds.maxX = Math.max(bounds.maxX, subtreeCenterX);
  bounds.minY = Math.min(bounds.minY, nodeY);
  bounds.maxY = Math.max(bounds.maxY, nodeY);

  if (parentId) {
    edges.push({
      id: `edge_${edgesIdCounter++}`,
      source: recipeId,
      target: parentId,
      type: 'smoothstep',
    });
  }

  return { width: totalSubtreeWidth, centerX: subtreeCenterX, bounds };
}


function buildTreeWrapper(recipe: RecipeNodeType): [Node[], Edge[], string, { minX: number; minY: number; maxX: number; maxY: number }] {
  const nodes: Node[] = [];
  const edges: Edge[] = [];
  const result = layoutTree(recipe, 0, 0, nodes, edges);
  const rootNode = nodes.find(n => n.id === 'node_0');
  return [nodes, edges, rootNode?.id || '', result.bounds];
}


function RecipeFlowInner({ tree }: RecipeFlowProps) {
  const { setCenter, getNode } = useReactFlow();

  const [nodes, edges, rootId, bounds] = useMemo(() => {
    if (!tree) return [[], [], '', { minX: 0, minY: 0, maxX: 0, maxY: 0 }];
    nodesIdCounter = 0;
    edgesIdCounter = 0;
    return buildTreeWrapper(tree);
  }, [tree]);

  useEffect(() => {
    if (rootId && bounds) {
      setTimeout(() => {
        const padding = 100;

        const minX = bounds.minX - padding;
        const minY = bounds.minY - padding;
        const maxX = bounds.maxX + padding;
        const maxY = bounds.maxY + padding;

        const viewportWidth = maxX - minX;
        const viewportHeight = maxY - minY;
        const centerX = minX + viewportWidth / 2;
        const centerY = minY + viewportHeight / 2;

        setCenter(centerX, centerY, {
          zoom: 1.5,
          duration: 800,
        });
      }, 100);
    }
  }, [rootId, getNode, setCenter, bounds]);

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={{ recipeNode: RecipeNode }}
      nodeOrigin={[0.5, 0.5]}
      onlyRenderVisibleElements
      fitView={false}
      translateExtent={[
        [bounds.minX - 200, bounds.minY - 200],
        [bounds.maxX + 200, bounds.maxY + 200],
      ]}
    >
      <MiniMap position="bottom-left" pannable zoomable style={{ width: 200, height: 150 }} nodeColor="#000" />
      <Background color="#ccc" variant={BackgroundVariant.Cross} lineWidth={1} />
    </ReactFlow>
  );
}

export default function RecipeFlow({ tree }: RecipeFlowProps) {
  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlowProvider>
        <RecipeFlowInner tree={tree} />
      </ReactFlowProvider>
    </div>
  );
}
