import React from "react";
import { useFontsLoaded } from "./hooks/useFontsLoaded";

interface FontLoaderProps {
  children: React.ReactNode;
}

export function FontLoader({ children }: FontLoaderProps) {
  const fontsLoaded = useFontsLoaded();

  if (!fontsLoaded) {
    return <div>Loading...</div>; // Or a loading spinner
  }

  return <>{children}</>;
}
