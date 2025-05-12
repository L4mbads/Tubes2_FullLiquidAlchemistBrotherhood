"use client";
import React, { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import '@/app/style.css';

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

      if (rate >= 3) {
        clearInterval(rampInterval);
        router.push('/search');
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
      >
        <source src="/Road.mp4" type="video/mp4" />
        Your browser does not support the video tag.
      </video>

      {/* Overlay content with forced white text */}
      <div style={{color: 'white', paddingTop: 50}} className="relative z-10 flex flex-col items-center w-full h-full text-center px-4">
        {!isStarting ? (
          <>
            <h1 style={{color: 'white'}} className="text-4xl md:text-6xl font-bold mb-6">
              Full Liquid Alchemist
            </h1>
            <button
              onClick={handleStart}
              className="search-button"
            >
              Start Searching
            </button>
          </>
        ) : (
          <h1 style={{color: 'white'}} className="text-4xl md:text-6xl font-bold mb-6">
            Accelerating...
          </h1>
        )}
      </div>
    </div>
  );
};

export default Page;