import React, { useState } from 'react';
import Image from 'next/image';
import './style.css';

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
        onChange={(e) => setElement(e.target.value)}
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
  return (
    <div className="element-option" onClick={() => onSelect(option.Name)}>
    <Image
    src={option.ImageUrl}
    alt={option.Name}
    width={40}
    height={40}
    className="element-image"
    loading="lazy"
    />
      <div className="element-details">
        <div className="element-name">{option.Name}</div>
        <div className="element-tier">Tier: {option.Type}</div>
      </div>
    </div>
  );
}

