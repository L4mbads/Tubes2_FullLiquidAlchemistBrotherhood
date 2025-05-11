import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const LandingPage: React.FC = () => {
  const navigate = useNavigate();
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isStarting, setIsStarting] = useState(false);

  const handleStart = () => {
    setIsStarting(true);
    // Gradually speed up video playback
    if (videoRef.current) {
      let rate = videoRef.current.playbackRate;
      const ramp = setInterval(() => {
        if (!videoRef.current) return clearInterval(ramp);
        rate = Math.min(rate + 0.5, 10); // cap at 10x speed
        videoRef.current.playbackRate = rate;
        if (rate >= 3) clearInterval(ramp);
      }, 100);
    }
    // Navigate to the next page after a delay (in ms)
    setTimeout(() => navigate('/search'), 1500);
  };

  return (
    <div className="relative w-screen h-screen overflow-hidden bg-black">
      {/* Background video */}
      <video
        ref={videoRef}
        className="absolute top-0 left-0 w-full h-full object-cover"
        src="./Road.mp4" // Video file path
        autoPlay
        loop
        muted
        playsInline
      />

      {/* Overlay content 
        Please adjust this
      */}
      {!isStarting && (
        <div className="relative z-10 flex flex-col items-center justify-center w-full h-full text-center text-white px-4">
          <h1 className="text-4xl md:text-6xl font-bold mb-6">
            Full Liquid Alchemist
          </h1>
          <button
            onClick={handleStart}
            className="px-8 py-4 bg-blue-600 hover:bg-blue-700 rounded text-lg md:text-xl transition"
          >
            Start Searching
          </button>
        </div>
      )}
    </div>
  );
};

export default LandingPage;
