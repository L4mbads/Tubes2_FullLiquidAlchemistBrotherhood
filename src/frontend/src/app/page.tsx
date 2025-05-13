"use client";
import React, { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import '@/app/style.css';
import DarkModeToggleButton from '@/components/DarkModeToggleButton';
import localFont from 'next/font/local';

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
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isStarting, setIsStarting] = useState(false);
  const router = useRouter();

  // Play the video once on mount
  useEffect(() => {
    const video = videoRef.current;
    if (video) {
      video.play().catch(error => {
        console.error("Background video playback failed:", error);
      });
    }
  }, []);

  // Handle acceleration
  useEffect(() => {
    if (!isStarting || !videoRef.current) return;

    const video = videoRef.current;
    let rate = 1;
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
          ref={videoRef}
          className="absolute top-0 left-0 w-full h-full object-cover"
          autoPlay
          loop
          muted
          playsInline
          style={{backgroundColor: 'black'}}
        >
          <source src="/Road.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>

        <div style={{color: 'white', paddingTop: 50}} className="relative z-10 flex flex-col items-center justify-center w-full h-full text-center px-4">
          {!isStarting ? (
            <>
              <div className="items-center" style={{background: '#1d1a2f', borderRadius: 8, paddingTop: 2}}><DarkModeToggleButton/></div>
              <div className={sybreFont.className}>
                <h1 style={{color: '#8de450', fontSize: 60, WebkitTextStroke: '3px', WebkitTextStrokeColor: '#1d1a2f'}} className='font-bold mb-6 select-none'>
                  Full Liquid Alchemist
                </h1>
              </div>
              <button
                onClick={handleStart}
                className="search-button"
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