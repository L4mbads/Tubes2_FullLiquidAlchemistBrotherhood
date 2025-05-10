'use client';

import { useState } from 'react';
import dynamic from 'next/dynamic';
import '@/app/style.css'

const RecipeFlow = dynamic(() => import('../components/RecipeFlow'), {
  ssr: false,
});

export default function Page() {
  const [showNumberInput, setShowNumberInput] = useState(false);
  const [recipeCount, setRecipeCount] = useState(1);

  const handleToggleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setShowNumberInput(e.target.checked);
    if (!e.target.checked) {
      setRecipeCount(1);
    }
  };
  return (
    <div className="page-container">
      <div className="left-panel">
        <h1>Little Alchemy 2 Recipe</h1>
        <div className='panel-section'>
          <input type="text" className='text-input' placeholder='Enter Elements...'></input>
        </div>
        <div className="panel-section">
            <label className="section-title">Algorithm</label>
            <div className="radio-group">
                <label className="radio-option">
                    <input type="radio" name="option" defaultChecked/>
                    DFS
                </label>
                <label className="radio-option">
                    <input type="radio" name="option"/>
                    BFS
                </label>
            </div>
        </div>
        <div className='panel-section'>
          <div className="toggle-container">
            <span className="toggle-label"><strong>Multiple Recipes</strong></span>
            <label className="toggle-switch">
              <input 
                type="checkbox"
                onChange={handleToggleChange}
              />
              <span className="toggle-slider"></span>
            </label>
          </div>
          {showNumberInput && (
            <div className="panel-section" style={{ marginTop: '10px' }}>
              <label className="section-title">Number of Recipes</label>
              <input
                type="number"
                className='text-input'
                min="1"
                max="10"
                value={recipeCount}
                onChange={(e) => setRecipeCount(Math.max(1, parseInt(e.target.value) || 1))}
              />
            </div>
          )}
        </div>
        <button className="search-btn">Search</button>
      </div>
      <div className="flow-container">
        <RecipeFlow />
      </div>
    </div>
  );
}