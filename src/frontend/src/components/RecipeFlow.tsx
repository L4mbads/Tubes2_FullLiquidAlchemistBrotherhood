import React, { useMemo } from 'react';
import ReactFlow, {Background, BackgroundVariant, Node, Edge, Position } from 'reactflow';
import 'reactflow/dist/style.css';
import '../components/style.css'; 

import RecipeNode, { RecipeNodeType } from './RecipeNode';

const recipeTree: RecipeNodeType = {
  "name": "Atmosphere",
  "recipes": [
    {
      "ingredient1": {
        "name": "Air",
        "recipes": null
      },
      "ingredient2": {
        "name": "Planet",
        "recipes": [
          {
            "ingredient1": {
              "name": "Continent",
              "recipes": [
                {
                  "ingredient1": {
                    "name": "Land",
                    "recipes": [
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Earth",
                          "recipes": null
                        }
                      },
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Stone",
                          "recipes": [
                            {
                              "ingredient1": {
                                "name": "Earth",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Pressure",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Air",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Air",
                                      "recipes": null
                                    }
                                  }
                                ]
                              }
                            },
                            {
                              "ingredient1": {
                                "name": "Air",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Lava",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Fire",
                                      "recipes": null
                                    }
                                  },
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Heat",
                                      "recipes": [
                                        {
                                          "ingredient1": {
                                            "name": "Air",
                                            "recipes": null
                                          },
                                          "ingredient2": {
                                            "name": "Energy",
                                            "recipes": [
                                              {
                                                "ingredient1": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                },
                                                "ingredient2": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                }
                                              }
                                            ]
                                          }
                                        }
                                      ]
                                    }
                                  }
                                ]
                              }
                            }
                          ]
                        }
                      }
                    ]
                  },
                  "ingredient2": {
                    "name": "Land",
                    "recipes": [
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Earth",
                          "recipes": null
                        }
                      },
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Stone",
                          "recipes": [
                            {
                              "ingredient1": {
                                "name": "Earth",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Pressure",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Air",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Air",
                                      "recipes": null
                                    }
                                  }
                                ]
                              }
                            },
                            {
                              "ingredient1": {
                                "name": "Air",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Lava",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Fire",
                                      "recipes": null
                                    }
                                  },
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Heat",
                                      "recipes": [
                                        {
                                          "ingredient1": {
                                            "name": "Air",
                                            "recipes": null
                                          },
                                          "ingredient2": {
                                            "name": "Energy",
                                            "recipes": [
                                              {
                                                "ingredient1": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                },
                                                "ingredient2": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                }
                                              }
                                            ]
                                          }
                                        }
                                      ]
                                    }
                                  }
                                ]
                              }
                            }
                          ]
                        }
                      }
                    ]
                  }
                },
                {
                  "ingredient1": {
                    "name": "Land",
                    "recipes": [
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Earth",
                          "recipes": null
                        }
                      },
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Stone",
                          "recipes": [
                            {
                              "ingredient1": {
                                "name": "Earth",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Pressure",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Air",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Air",
                                      "recipes": null
                                    }
                                  }
                                ]
                              }
                            },
                            {
                              "ingredient1": {
                                "name": "Air",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Lava",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Fire",
                                      "recipes": null
                                    }
                                  },
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Heat",
                                      "recipes": [
                                        {
                                          "ingredient1": {
                                            "name": "Air",
                                            "recipes": null
                                          },
                                          "ingredient2": {
                                            "name": "Energy",
                                            "recipes": [
                                              {
                                                "ingredient1": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                },
                                                "ingredient2": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                }
                                              }
                                            ]
                                          }
                                        }
                                      ]
                                    }
                                  }
                                ]
                              }
                            }
                          ]
                        }
                      }
                    ]
                  },
                  "ingredient2": {
                    "name": "Earth",
                    "recipes": null
                  }
                }
              ]
            },
            "ingredient2": {
              "name": "Continent",
              "recipes": [
                {
                  "ingredient1": {
                    "name": "Land",
                    "recipes": [
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Earth",
                          "recipes": null
                        }
                      },
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Stone",
                          "recipes": [
                            {
                              "ingredient1": {
                                "name": "Earth",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Pressure",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Air",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Air",
                                      "recipes": null
                                    }
                                  }
                                ]
                              }
                            },
                            {
                              "ingredient1": {
                                "name": "Air",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Lava",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Fire",
                                      "recipes": null
                                    }
                                  },
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Heat",
                                      "recipes": [
                                        {
                                          "ingredient1": {
                                            "name": "Air",
                                            "recipes": null
                                          },
                                          "ingredient2": {
                                            "name": "Energy",
                                            "recipes": [
                                              {
                                                "ingredient1": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                },
                                                "ingredient2": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                }
                                              }
                                            ]
                                          }
                                        }
                                      ]
                                    }
                                  }
                                ]
                              }
                            }
                          ]
                        }
                      }
                    ]
                  },
                  "ingredient2": {
                    "name": "Land",
                    "recipes": [
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Earth",
                          "recipes": null
                        }
                      },
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Stone",
                          "recipes": [
                            {
                              "ingredient1": {
                                "name": "Earth",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Pressure",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Air",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Air",
                                      "recipes": null
                                    }
                                  }
                                ]
                              }
                            },
                            {
                              "ingredient1": {
                                "name": "Air",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Lava",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Fire",
                                      "recipes": null
                                    }
                                  },
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Heat",
                                      "recipes": [
                                        {
                                          "ingredient1": {
                                            "name": "Air",
                                            "recipes": null
                                          },
                                          "ingredient2": {
                                            "name": "Energy",
                                            "recipes": [
                                              {
                                                "ingredient1": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                },
                                                "ingredient2": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                }
                                              }
                                            ]
                                          }
                                        }
                                      ]
                                    }
                                  }
                                ]
                              }
                            }
                          ]
                        }
                      }
                    ]
                  }
                },
                {
                  "ingredient1": {
                    "name": "Land",
                    "recipes": [
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Earth",
                          "recipes": null
                        }
                      },
                      {
                        "ingredient1": {
                          "name": "Earth",
                          "recipes": null
                        },
                        "ingredient2": {
                          "name": "Stone",
                          "recipes": [
                            {
                              "ingredient1": {
                                "name": "Earth",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Pressure",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Air",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Air",
                                      "recipes": null
                                    }
                                  }
                                ]
                              }
                            },
                            {
                              "ingredient1": {
                                "name": "Air",
                                "recipes": null
                              },
                              "ingredient2": {
                                "name": "Lava",
                                "recipes": [
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Fire",
                                      "recipes": null
                                    }
                                  },
                                  {
                                    "ingredient1": {
                                      "name": "Earth",
                                      "recipes": null
                                    },
                                    "ingredient2": {
                                      "name": "Heat",
                                      "recipes": [
                                        {
                                          "ingredient1": {
                                            "name": "Air",
                                            "recipes": null
                                          },
                                          "ingredient2": {
                                            "name": "Energy",
                                            "recipes": [
                                              {
                                                "ingredient1": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                },
                                                "ingredient2": {
                                                  "name": "Fire",
                                                  "recipes": null
                                                }
                                              }
                                            ]
                                          }
                                        }
                                      ]
                                    }
                                  }
                                ]
                              }
                            }
                          ]
                        }
                      }
                    ]
                  },
                  "ingredient2": {
                    "name": "Earth",
                    "recipes": null
                  }
                }
              ]
            }
          }
        ]
      }
    }
  ]
};


function calculateSubtreeWidth(recipe: RecipeNodeType): number {
  if (!recipe.recipes || recipe.recipes.length === 0) {
    return 1;
  }
  
  let totalWidth = 0;
  recipe.recipes.forEach(r => {
    const ingredient1Width = calculateSubtreeWidth(r.ingredient1);
    const ingredient2Width = calculateSubtreeWidth(r.ingredient2);
    totalWidth += ingredient1Width + ingredient2Width;
  });
  
  return Math.max(1, totalWidth);
}

let nodesIdCounter = 0;
let edgesIdCounter = 0;

function buildTree(
  recipe: RecipeNodeType,
  depth = 0,
  x = 0,
  nodes: Node[] = [],
  edges: Edge[] = [],
  parentId: string | null = null
): [Node[], Edge[]] {
  const resultNodeId = `node_${nodesIdCounter++}`;
  const widthSpacing = 100;
  const heightSpacing = 150;

  nodes.push({
    id: resultNodeId,
    type: 'recipeNode',
    data: { name: recipe.name },
    position: { x, y: -depth * heightSpacing },
  });


  if (parentId) {
    edges.push({
      id: `edge_${edgesIdCounter++}`,
      source: resultNodeId,
      target: parentId,
      type: 'smoothstep',
      animated: true,
    });
  }
  
  if (recipe.recipes) {
    let currentX = x;

    recipe.recipes.forEach((r) => {
      const ingredient1Width = calculateSubtreeWidth(r.ingredient1);
      const ingredient2Width = calculateSubtreeWidth(r.ingredient2);
      const totalWidth = ingredient1Width + ingredient2Width;
      const ingredient1X = currentX - (ingredient1Width * widthSpacing) / 2;
      const ingredient2X = currentX + (ingredient2Width * widthSpacing) / 2;

      const stepNodeId = `step_${nodesIdCounter++}`;
      const stepX = currentX;
      const stepY = -(depth + 0.5) * heightSpacing;

      nodes.push({
        id: stepNodeId,
        type: 'default',
        data: { label: '' },
        position: { x: stepX, y: stepY },
        style: {
          background: 'black',
          border: 'none',
          borderRadius: '50%',
          width: '6px',
          height: '6px',
          padding: 0,
          margin: 0,
          boxShadow: 'none',
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
        source: stepNodeId,
        target: resultNodeId,
        type: 'smoothstep',
        animated: true,
      });

      buildTree(r.ingredient1, depth + 1, ingredient1X, nodes, edges, stepNodeId);

      buildTree(r.ingredient2, depth + 1, ingredient2X, nodes, edges, stepNodeId);

      currentX += totalWidth * widthSpacing;
    });
  }

  return [nodes, edges];
}



export default function RecipeFlow() {
  const [nodes, edges] = useMemo(() => buildTree(recipeTree), []);

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        nodeTypes={{ recipeNode: RecipeNode }}
        nodeOrigin={[0.5, 0.5]}
      >
        <Background color="#ccc" variant={BackgroundVariant.Cross} lineWidth={1} />
      </ReactFlow>
    </div>
  );
}