'use client';

import { useState, useEffect, useContext } from 'react';
import dynamic from 'next/dynamic';
import '@/app/search/style.css'
import axios from 'axios';
import AutoCompleteInput from '@/components/AutoCompleteInput';
import { RecipeNodeType } from '@/components/RecipeNode';
import DarkModeToggleButton from '@/components/DarkModeToggleButton';
import { DarkModeContext } from '@/components/DarkModeProvider';

const RecipeFlow = dynamic(() => import('../../components/RecipeFlow'), {
  ssr: false,
});

export default function Page() {
  const context = useContext(DarkModeContext);

  if (!context) {
    throw new Error('No Context!');
  }

  const { darkMode } = context;
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

  const [useLiveSearch, setUseLiveSearch] = useState(false);

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
        console.log(res.data);
        setRecipeTree(res.data);
      })
      .catch(err => {
        let errorMessage = "Error fetching recipe.";

        if (err.response && err.response.data) {
          errorMessage = err.response.data.message || errorMessage;
        }

        console.error("Error fetching recipe:", errorMessage);
        alert(errorMessage);
      })
      .finally(() => setLoading(false));
  };

  const handleSSEMessage = (data: string) => {
    try {
      const parsed = JSON.parse(data);
      setRecipeTree(parsed);
      console.log(parsed);
    } catch (error) {
      console.error("Invalid SSE data:", error);
    }
  };

  const handleLiveSearch = () => {
    if (!selectedElement) {
      alert("Please select an element.");
      return;
    }

    setLoading(true);
    setLoadTime(null);
    const start = performance.now();

    const url = `http://localhost:8000/api/go/liverecipe?element=${encodeURIComponent(selectedElement)}&strategy=${strategy}&count=${recipeCount}`;
    const eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      handleSSEMessage(event.data);
    };

    eventSource.addEventListener("done", () => {
      const end = performance.now();
      setLoadTime(end - start);
      setLoading(false);
      eventSource.close();
    });

    eventSource.onerror = () => {
      console.error("SSE connection error.");
      eventSource.close();
      setLoading(false);
    };
  };


  return (
      <div className="page-container">
        <div className= {darkMode? "left-panel left-panel-dark" : "left-panel left-panel-light"}>
          <DarkModeToggleButton />
          <h1 style={{color: darkMode? '#8bd450': '#de5857'}}>Little Alchemy 2 Recipe</h1>
          <div className='panel-section'>
            <AutoCompleteInput options={elements}   onSelect={setSelectedElement}></AutoCompleteInput>
          </div>
          <div className="panel-section">
              <label className={darkMode ? "section-title section-title-dark" : "section-title section-title-light"}>Algorithm</label>
              <div className="radio-group">
                  <label className="radio-option" style={{color: darkMode? 'white' : 'black'}}>
                      <input 
                      type="radio" 
                      name="option" 
                      value="dfs"
                      checked={strategy === 'dfs'}
                      onChange={() => setStrategy('dfs')}
                      className={darkMode? 'radio-dark': 'radio-light'}/>
                      DFS
                  </label>
                  <label className="radio-option" style={{color: darkMode? 'white' : 'black'}}>
                      <input 
                      type="radio" 
                      name="option"
                      value="bfs"
                      checked={strategy === 'bfs'}
                      onChange={() => setStrategy('bfs')}
                      className={darkMode? 'radio-dark': 'radio-light'}/>
                      BFS
                  </label>
              </div>
          </div>
          <div className='panel-section'>
            <div className="toggle-container">
              <span className="toggle-label" style={{color: darkMode? 'white' : 'black'}}><strong>Multiple Recipes</strong></span>
              <label className="toggle-switch">
                <input 
                  type="checkbox"
                  onChange={handleToggleChange}
                />
                <span className={darkMode ? 'toggle-slider-dark' : 'toggle-slider-light'}></span>
              </label>
            </div>
            {showNumberInput && (
              <div className="panel-section" style={{ marginTop: '10px' }}>
                <label className="section-title" style={{color: darkMode? 'white' : 'black'}}>Number of Recipes</label>
                <input
                  type="number"
                  className='text-input'
                  min="1"
                  value={recipeCount}
                  onChange={(e) => setRecipeCount(Math.max(1, parseInt(e.target.value) || 1))}
                />
              </div>
            )}
          </div>
          <div className="toggle-container">
            <span className='toggle-label'style={{color: darkMode ? 'white' : 'black'}}><strong>Live Visualization</strong></span>
            <label className="toggle-switch">
              <input
                type="checkbox"
                checked={useLiveSearch}
                onChange={() => setUseLiveSearch(prev => !prev)}
              />
              <span className={darkMode ? 'toggle-slider-dark' : 'toggle-slider-light'}></span>
            </label>
          </div>
          <button className={darkMode ? "search-btn search-dark" : "search-btn search-light"} onClick={useLiveSearch ? handleLiveSearch : handleSearch}>Search</button>
          {loading && (
            <div className="loading-indicator">
              <div className={darkMode ? "spinner spinner-dark" : "spinner spinner-light"} />
              <p style={{color: darkMode? 'white' : 'black'}}>Loading...</p>
            </div>
          )}

          {loadTime !== null && !loading && (
            <p className="load-time" style={{color: darkMode? 'white' : 'black'}}>Loaded in {(loadTime / 1000).toFixed(2)}s</p>
          )}
        </div>
        <div className={darkMode ? "flow-container flow-dark" : "flow-container flow-light"}>
          <RecipeFlow tree={recipeTree} isLive={useLiveSearch}/>
        </div>
      </div>
  );
}