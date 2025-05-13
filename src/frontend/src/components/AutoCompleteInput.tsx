import React, { useState, useContext } from 'react';
import Image from 'next/image';
import './style.css';
import { DarkModeContext } from './DarkModeProvider';

type Option = {
  Name: string;
  Type: string;
  ImageUrl: string;
};

type ElementOptionsProps = {
  option: Option;
  onSelect: (value: string) => void;
};

type AutoCompleteInputProps = {
  options: Option[];
  onSelect: (value: string) => void;
};

export default function AutoCompleteInput({ options, onSelect }: AutoCompleteInputProps) {
  const [showOptions, setShowOptions] = useState<boolean>(false);
  const [element, setElement] = useState<string>('');

  const selectElementHandler = (newElementValue: string) => {
    setElement(newElementValue);
    setShowOptions(false);
    onSelect(newElementValue);
  }

  return (
    <div>
      <input
        type='text'
        className='text-input'
        placeholder='Enter elements...'
        onClick={() => setShowOptions(!showOptions)}
        value={element}
        onChange={(e) => 
          {
            setElement(e.target.value)
          }}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            selectElementHandler(element);
          }
        }}
      />
      {showOptions && (
        <ul className='options'>
            {options
            .filter((opt) => opt.Name.toLowerCase().includes(element.toLowerCase()))
            .map((opt, index) => (
                <ElementOptions key={index} option={opt} onSelect={selectElementHandler} />
            ))}
        </ul>
      )}
    </div>
  );
}

function ElementOptions({ option, onSelect }: ElementOptionsProps) {
  const context = useContext(DarkModeContext);

  if (!context) {
    throw new Error('No Context!');
  }

  const { darkMode } = context;
  return (
    <div className={darkMode ? "element-option element-option-dark" : "element-option element-option-light"} onClick={() => onSelect(option.Name)}>
    <Image
    src={option.ImageUrl}
    alt={option.Name}
    width={40}
    height={40}
    className="element-image"
    loading="lazy"
    />
      <div className="element-details">
        <div className="element-name" style={{color: darkMode ? 'white' : 'black'}}>{option.Name}</div>
        <div className="element-tier" style={{color: darkMode ? 'white' : 'black'}}>Tier: {option.Type}</div>
      </div>
    </div>
  );
}

