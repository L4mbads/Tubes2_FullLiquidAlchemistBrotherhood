'use client';
import { useEffect, useState } from 'react';
import dynamic from 'next/dynamic';

import { RecipeNodeType } from '@/components/RecipeNode';

const RecipeFlow = dynamic(() => import('../../components/RecipeFlow'), {
  ssr: false,
});

const LiveRecipePage = () => {
  const [messages, setMessages] = useState<string[]>([]);

  const [recipeTree, setRecipeTree] = useState<RecipeNodeType | null>(null);
  useEffect(() => {
    // Establish an SSE connection
    const eventSource = new EventSource("http://localhost:8000/api/go/liverecipe?element=Bullet&count=5&strategy=bfs");

    // Listen for messages from the server
    eventSource.onmessage = (event) => {
      setMessages((prevMessages) => [event.data]);
      // setRecipeTree(event.data)
    };
    eventSource.addEventListener("done", () => {
      console.log("Recipe generation completed.");
      eventSource.close();
    });
    // Handle errors
    eventSource.onerror = () => {
      console.error("Error occurred with SSE connection.");
      eventSource.close();
    };

    // Cleanup the connection when the component unmounts
    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div>
      <h1>Live Recipe Updates</h1>
        <div className="absolute top-0">
        {messages.map((message, index) => (
          <p key={index}>{message}</p>
        ))}
      </div>

      {/* <div>
          <RecipeFlow tree={recipeTree}/>
        </div> */}
    </div>
  );
};

export default LiveRecipePage;