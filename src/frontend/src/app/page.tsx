"use client";
import React, { useState, useRef, useEffect, useContext } from 'react';
import { useRouter } from 'next/navigation';
import '@/app/style.css';
import DarkModeToggleButton from '@/components/DarkModeToggleButton';
import localFont from 'next/font/local';
import { DarkModeContext } from '@/components/DarkModeProvider';

const sybreFont = localFont({
  src: '../fonts/Sybre.ttf',
  variable: '--font-sybre',
  display: 'swap',
});

const futronsFont = localFont({
  src: '../fonts/Futrons Demo.otf',
  variable: '--font-futrons',
  display: 'swap',
})

const Page: React.FC = () => {
  const context = useContext(DarkModeContext);

  if (!context) {
    throw new Error('No Context!');
  }

  const { darkMode } = context;  
  const darkVideoRef = useRef<HTMLVideoElement>(null);
  const lightVideoRef = useRef<HTMLVideoElement>(null);
  const [isStarting, setIsStarting] = useState(false);
  const router = useRouter();

  // Play the video once on mount
  useEffect(() => {
  const video = darkMode ? darkVideoRef.current : lightVideoRef.current;    if (video) {
      video.play().catch(error => {
        console.error("Background video playback failed:", error);
      });
    }
  }, []);

  // Handle acceleration
  useEffect(() => {
    const video = darkMode ? darkVideoRef.current : lightVideoRef.current;    let rate = 1;
    if (!isStarting || !video) return;
    video.playbackRate = rate;

    const rampInterval = setInterval(() => {
      rate = Math.min(rate + 0.5, 10);
      video.playbackRate = rate;

      if (rate >= 10) {
        clearInterval(rampInterval);
        setTimeout(() => {
          router.push('/search');
        }, 100); // 100ms delay
      }
    }, 100);

    return () => clearInterval(rampInterval);
  }, [isStarting, router]);

  const handleStart = () => {
    setIsStarting(true);
  };

  return (
      <div className="relative w-screen h-screen overflow-hidden bg-black" style={{color: 'black'}}>
        <video
          ref={darkVideoRef}
          className="absolute top-0 left-0 w-full h-full object-cover"
          autoPlay
          loop
          muted
          playsInline
          style={{backgroundColor: 'black', opacity: darkMode ? '1' : '0'}}
        >
          <source src="/Road.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
        <video
          ref={lightVideoRef}
          className="absolute top-0 left-0 w-full h-full object-cover"
          autoPlay
          loop
          muted
          playsInline
          style={{backgroundColor: 'black', opacity: darkMode ? '0' : '1'}}
        >
          <source src="/Road2.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>

        <div style={{color: 'white', paddingTop: 50}} className="relative z-10 flex flex-col items-center justify-center w-full h-full text-center px-4">
          {!isStarting ? (
            <>
              <div className="items-center" style={{background: '#1d1a2f', borderRadius: 8, paddingTop: 2}}><DarkModeToggleButton/></div>
              <div className={sybreFont.className}>
                <h1 className={darkMode ? 'title-label title-dark' : 'title-label title-light'}>
                  Full Liquid Alchemist
                </h1>
              </div>
              <button
                onClick={handleStart}
                className={darkMode ? "search-button search-dark" : "search-button search-light"} 
              >
                <div className={futronsFont.className}>Start Searching</div>
              </button>
            </>
          ) : (<></>
          )}
        </div>
      </div>
  );
};

export default Page;