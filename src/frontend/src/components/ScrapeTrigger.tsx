'use client';

import { useEffect } from 'react';
import axios from 'axios';

export default function ScrapeTrigger() {
  useEffect(() => {
    axios.get('http://localhost:8000/api/go/scrape')
      .then(res => console.log('Scraping triggered:', res.data))
      .catch(err => console.error('Scraping failed:', err));
  }, []);

  return null;
}
