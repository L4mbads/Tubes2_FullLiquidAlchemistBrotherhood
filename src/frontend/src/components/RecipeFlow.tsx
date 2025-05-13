import React, { useMemo, useEffect, useCallback, useContext } from 'react';
import ReactFlow, {
  Background,
  Node,
  Edge,
  Position,
  MiniMap,
  useReactFlow,
  ReactFlowProvider
} from 'reactflow';
import 'reactflow/dist/style.css';
import '../components/style.css';
import RecipeNode, { RecipeNodeType } from './RecipeNode';
import { DarkModeContext } from './DarkModeProvider';

type RecipeFlowProps = {
  tree: RecipeNodeType | null;
  isLive: boolean,
};

type Bounds = {
  minX: number;
  minY: number;
  maxX: number;
  maxY: number;
};

const WIDTH_SPACING = 175;
const HEIGHT_SPACING = 175;
const PADDING = 100;

function RecipeFlowInner({ tree, isLive }: RecipeFlowProps) {
  const context = useContext(DarkModeContext);

  if (!context) {
    throw new Error('No Context!');
  }

  const { darkMode } = context;

  const { setCenter } = useReactFlow();

  const buildTree = useCallback((recipe: RecipeNodeType): {
    nodes: Node[];
    edges: Edge[];
    rootId: string;
    bounds: Bounds;
  } => {
    const nodes: Node[] = [];
    const edges: Edge[] = [];
    let nodesIdCounter = 0;
    let edgesIdCounter = 0;
    const bounds: Bounds = { minX: Infinity, minY: Infinity, maxX: -Infinity, maxY: -Infinity };

    const widthCache = new Map<RecipeNodeType, number>();

    const calculateSubtreeWidth = (recipe: RecipeNodeType): number => {
      if (recipe == null) {
        return 0
      }
      if (widthCache.has(recipe)) {
        return widthCache.get(recipe)!;
      }

      if (!recipe.recipes || recipe.recipes.length === 0) {
        widthCache.set(recipe, 1);
        return 1;
      }

      let totalWidth = 0;
      recipe.recipes.forEach(r => {
        totalWidth += calculateSubtreeWidth(r.ingredient1) + calculateSubtreeWidth(r.ingredient2);
      });

      const width = Math.max(1, totalWidth);
      widthCache.set(recipe, width);
      return width;
    };

    const layoutTree = (
      recipe: RecipeNodeType,
      depth: number,
      xOffset: number,
      parentId: string | null = null
    ): { width: number; centerX: number } => {
      if(recipe == null) {
        return {width: 0, centerX: 0};
      }
      const recipeId = `node_${nodesIdCounter++}`;
      let subtreeCenterX = xOffset;

      const totalSubtreeWidth = widthCache.get(recipe) || 1;

      if (recipe.recipes && recipe.recipes.length > 0) {
        let currentX = xOffset;
        const childrenCenters: number[] = [];

        for (const subRecipe of recipe.recipes) {
          const w1 = widthCache.get(subRecipe.ingredient1) || 1;
          const w2 = widthCache.get(subRecipe.ingredient2) || 1;
          const totalWidth = w1 + w2;
          const midX = currentX + (totalWidth * WIDTH_SPACING) / 2;
          const stepId = `step_${nodesIdCounter++}`;
          const stepY = -((depth + 0.5) * HEIGHT_SPACING);

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
            style: {stroke: darkMode? 'white' : 'black', strokeWidth: '2px'}
          });

          layoutTree(subRecipe.ingredient1, depth + 1, currentX, stepId);
          layoutTree(
            subRecipe.ingredient2,
            depth + 1,
            currentX + w1 * WIDTH_SPACING,
            stepId
          );

          currentX += totalWidth * WIDTH_SPACING;
          childrenCenters.push(midX);

          bounds.minX = Math.min(bounds.minX, midX);
          bounds.maxX = Math.max(bounds.maxX, midX);
          bounds.minY = Math.min(bounds.minY, stepY);
          bounds.maxY = Math.max(bounds.maxY, stepY);
        }

        if (childrenCenters.length > 0) {
          subtreeCenterX = childrenCenters.reduce((sum, x) => sum + x, 0) / childrenCenters.length;
        }
      }

      const nodeY = -depth * HEIGHT_SPACING;

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
          style: {stroke: darkMode? 'white' : 'black', strokeWidth: '2px'}
        });
      }

      return { width: totalSubtreeWidth, centerX: subtreeCenterX };
    };

    calculateSubtreeWidth(recipe);

    layoutTree(recipe, 0, 0);
    const rootId = 'node_0';

    return { nodes, edges, rootId, bounds };
  }, [darkMode]);

  const { nodes, edges, rootId, bounds } = useMemo(() => {
    if (!tree) return { nodes: [], edges: [], rootId: '', bounds: { minX: 0, minY: 0, maxX: 0, maxY: 0 } };
    return buildTree(tree);
  }, [tree, buildTree]);

  useEffect(() => {
    if (!isLive) {
      if (rootId && bounds.minX !== Infinity) {
        const timer = setTimeout(() => {
          const viewportWidth = bounds.maxX - bounds.minX + (2 * PADDING);
          const centerX = bounds.minX + (viewportWidth / 2) - PADDING;

          setCenter(centerX, 0, {
            zoom: 1.5,
            duration: 800,
          });
        }, 100);

        return () => clearTimeout(timer);
      }
    }
  }, [rootId, setCenter, bounds, isLive]);

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={{ recipeNode: RecipeNode }}
      nodeOrigin={[0.5, 0.5]}
      onlyRenderVisibleElements
      fitView={isLive ? true : false}
      translateExtent={[
        [bounds.minX - 500, bounds.minY - 500],
        [bounds.maxX + 500, bounds.maxY + 500],
      ]}
    >
      <MiniMap
        position="bottom-left"
        pannable
        zoomable
        style={{ width: 200, height: 150, backgroundColor: darkMode ? '#734f9a': '#fbac4e' }}
        nodeColor="#1d1a2f"
      />
      <Background
        color="#734f9a"
        lineWidth={1}
      />
    </ReactFlow>
  );
}

export default function RecipeFlow({ tree, isLive }: RecipeFlowProps) {
  return (
    <div style={{ width: '100%', height: '100vh' }}>
      <ReactFlowProvider>
        <RecipeFlowInner tree={tree} isLive={isLive} />
      </ReactFlowProvider>
    </div>
  );
}