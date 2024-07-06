import { Suspense } from "react";
import { useFontsLoaded } from "../hooks/useFontsLoaded";

import PortfolioItems, {
  type PortfolioItem,
} from "@components/Headless/Portfolio/PortfolioItems";

const attestoPortfolioItems: PortfolioItem[] = [
  {
    src: "/screenshots/attesto/attesto-01-checkin.png",
    alt: "Attesto Check-In",
  },
  {
    src: "/screenshots/attesto/attesto-02-welcome.png",
    alt: "Attesto Welcome",
  },
  {
    src: "/screenshots/attesto/attesto-03-role.png",
    alt: "Attesto Select Role",
  },
  {
    src: "/screenshots/attesto/attesto-04-questions.png",
    alt: "Attesto Questions",
  },
  {
    src: "/screenshots/attesto/attesto-05-add-question.png",
    alt: "Attesto Add Question",
  },
  {
    src: "/screenshots/attesto/attesto-06-select-candidates.png",
    alt: "Attesto Select Candidates",
  },
  {
    src: "/screenshots/attesto/attesto-07-upload-resume.png",
    alt: "Attesto Upload Resume",
  },
  {
    src: "/screenshots/attesto/attesto-08-check-assessments.png",
    alt: "Attesto Check Assessments",
  },
  {
    src: "/screenshots/attesto/attesto-09-assignment-popup.png",
    alt: "Attesto Assignment Popup",
  },
  {
    src: "/screenshots/attesto/attesto-10-questions.png",
    alt: "Attesto Assignment Questions",
  },
  {
    src: "/screenshots/attesto/attesto-11-questions-qa.png",
    alt: "Attesto Questions QA",
  },
  {
    src: "/screenshots/attesto/attesto-12-end-demo.png",
    alt: "Attesto End Demo",
  },
];

export default function Portfolio() {
  const fontsLoaded = useFontsLoaded();

  if (!fontsLoaded) {
    return <div>Loading...</div>; // Or a loading spinner
  }

  return (
    <div className="flex flex-col space-y-4">
      <h1 className="text-3xl text:xs-3xl md:text-6xl font-semibold text-white">
        Projects
      </h1>
      <h2 className="text-2xl md:text-4xl">Attesto - 2024</h2>
      <h3 className="text-xl md:text-2xl">Front End Demo for clients.</h3>
      <p>
        Converted Figma designs into React components with interactivity to
        mimic a real workflow.
      </p>
      <PortfolioItems items={attestoPortfolioItems} />
    </div>
  );
}
