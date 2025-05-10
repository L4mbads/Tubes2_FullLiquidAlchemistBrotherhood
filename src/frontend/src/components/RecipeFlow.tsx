import React, { useMemo } from 'react';
import ReactFlow, { Background, BackgroundVariant, Node, Edge, Position } from 'reactflow';
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
  parentId: string | null = null
): { width: number; centerX: number } {
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
          border: 'none',
          width: '1px',
          height: '1px',
          padding: 0,
          margin: 0,
          pointerEvents:'none',
        },
        sourcePosition: Position.Bottom,
        targetPosition: Position.Top,
        className: 'node-with-hidden-handles',
        draggable: false,
        selectable: false,
        hidden: false,
        
      });

      edges.push({
        id: `edge_${edgesIdCounter++}`,
        source: stepId,
        target: recipeId,
        type: 'smoothstep',
        animated: true,
      });

      const left = layoutTree(subRecipe.ingredient1, depth + 1, currentX, nodes, edges, stepId);
      const right = layoutTree(
        subRecipe.ingredient2,
        depth + 1,
        currentX + w1 * widthSpacing,
        nodes,
        edges,
        stepId
      );

      currentX += totalWidth * widthSpacing;
      childrenWidths.push({ width: totalWidth, centerX: midX });
    }

    const avgX = childrenWidths.reduce((sum, w) => sum + w.centerX, 0) / childrenWidths.length;
    subtreeCenterX = avgX;
    totalSubtreeWidth = childrenWidths.reduce((sum, c) => sum + c.width, 0);
  }

  nodes.push({
    id: recipeId,
    type: 'recipeNode',
    data: { name: recipe.name },
    position: { x: subtreeCenterX, y: -depth * heightSpacing },
  });

  if (parentId) {
    edges.push({
      id: `edge_${edgesIdCounter++}`,
      source: recipeId,
      target: parentId,
      type: 'smoothstep',
      animated: true,
    });
  }

  return { width: totalSubtreeWidth, centerX: subtreeCenterX };
}


function buildTreeWrapper(recipe: RecipeNodeType): [Node[], Edge[]] {
  const nodes: Node[] = [];
  const edges: Edge[] = [];
  layoutTree(recipe, 0, 0, nodes, edges);
  return [nodes, edges];
}

export default function RecipeFlow({ tree }: RecipeFlowProps) {
  const [nodes, edges] = useMemo(() => {
    if (!tree) return [[], []];
    nodesIdCounter = 0;
    edgesIdCounter = 0;
    return buildTreeWrapper(tree);
  }, [tree]);

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        nodeTypes={{ recipeNode: RecipeNode }}
        nodeOrigin={[0.5, 0.5]}
        onlyRenderVisibleElements
      >
        <Background color="#ccc" variant={BackgroundVariant.Cross} lineWidth={1} />
      </ReactFlow>
    </div>
  );
}
