'use client';

import { useState, useEffect } from 'react';
import dynamic from 'next/dynamic';
import '@/app/style.css'
import axios from 'axios';
import AutoCompleteInput from '@/components/AutoCompleteInput';
import { RecipeNodeType } from '@/components/RecipeNode';

const RecipeFlow = dynamic(() => import('../components/RecipeFlow'), {
  ssr: false,
});

export default function Page() {

  const [elements, setElements] = useState([])
  useEffect(() => {
    axios.get('http://localhost:8000/api/go/elements')
    .then(res => setElements(res.data))
    .catch(err => console.log(err))
  }, [])

  // console.log(elements)
  const [selectedElement, setSelectedElement] = useState('');
  const [strategy, setStrategy] = useState('dfs');

  const [showNumberInput, setShowNumberInput] = useState(false);
  const [recipeCount, setRecipeCount] = useState(1);

  const [recipeTree, setRecipeTree] = useState<RecipeNodeType | null>(null);

  const [loading, setLoading] = useState(false);
  const [loadTime, setLoadTime] = useState<number | null>(null);

  const handleToggleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setShowNumberInput(e.target.checked);
    if (!e.target.checked) {
      setRecipeCount(1);
    }
  };

  const handleSearch = () => {
    if (!selectedElement) {
      alert("Please select an element.");
      return;
    }

    setLoading(true);
    setLoadTime(null);
    const start = performance.now();

    const url = `http://localhost:8000/api/go/recipe?element=${encodeURIComponent(selectedElement)}&strategy=${strategy}&count=${recipeCount}`;
    axios.get(url)
      .then(res => {
        const end = performance.now();
        setLoadTime(end - start);
        setRecipeTree(res.data);
      })
      .catch(err => {
        console.error("Error fetching recipe:", err);
      })
      .finally(() => setLoading(false));
  };

  return (
    <div className="page-container">
      <div className="left-panel">
        <h1>Little Alchemy 2 Recipe</h1>
        <div className='panel-section'>
          <AutoCompleteInput options={elements}   onSelect={setSelectedElement}></AutoCompleteInput>
        </div>
        <div className="panel-section">
            <label className="section-title">Algorithm</label>
            <div className="radio-group">
                <label className="radio-option">
                    <input 
                    type="radio" 
                    name="option" 
                    value="dfs"
                    checked={strategy === 'dfs'}
                    onChange={() => setStrategy('dfs')}/>
                    DFS
                </label>
                <label className="radio-option">
                    <input 
                    type="radio" 
                    name="option"
                    value="bfs"
                    checked={strategy === 'bfs'}
                    onChange={() => setStrategy('bfs')}/>
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
        <button className="search-btn" onClick={handleSearch}>Search</button>
        {loading && (
          <div className="loading-indicator">
            <div className="spinner" />
            <p>Loading...</p>
          </div>
        )}

        {loadTime !== null && !loading && (
          <p className="load-time">Loaded in {(loadTime / 1000).toFixed(2)}s</p>
        )}
      </div>
      <div className="flow-container">
        <RecipeFlow tree={recipeTree}/>
      </div>
    </div>
  );
}