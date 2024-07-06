import { useState, useEffect } from "react";

export function useFontsLoaded() {
  const [fontsLoaded, setFontsLoaded] = useState(false);

  useEffect(() => {
    document.fonts.ready.then(() => {
      setFontsLoaded(true);
    });
  }, []);

  return fontsLoaded;
}

